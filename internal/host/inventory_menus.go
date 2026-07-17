package host

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/bedrock-gophers/plugins/internal/native"
	dfunsafe "github.com/bedrock-gophers/unsafe"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const (
	inventoryMenuOpenDelay           = 500 * time.Millisecond
	inventoryMenuRejectedContainerID = byte(0xff)
)

// InventoryMenus owns the fake block inventories for one framework run.
type InventoryMenus struct {
	mu      sync.Mutex
	players *Players
	menus   map[*world.EntityHandle]*openInventoryMenu
}

type openInventoryMenu struct {
	id            uint64
	container     native.InventoryMenuContainer
	containerType byte
	inv           *inventory.Inventory
	pos           cube.Pos
	sessionPos    cube.Pos
	windowID      byte
	openAt        time.Time
	openSent      bool
	double        bool
}

type inventoryMenuSpec struct {
	block         world.Block
	blockActorID  string
	containerType byte
	size          int
	double        bool
}

func NewInventoryMenus(players *Players) *InventoryMenus {
	menus := &InventoryMenus{players: players, menus: map[*world.EntityHandle]*openInventoryMenu{}}
	if players != nil {
		players.attachInventoryMenus(menus)
	}
	return menus
}

func (m *InventoryMenus) SendPlayerInventoryMenu(invocation native.InvocationID, id native.PlayerID, value native.PlayerInventoryMenu) bool {
	if m == nil || m.players == nil || value.ID == 0 {
		return false
	}
	spec, ok := inventoryMenuSpecFor(value.Container)
	if !ok || len(value.Items) != spec.size {
		return false
	}
	stacks := make([]item.Stack, len(value.Items))
	for index, nativeStack := range value.Items {
		stack, valid := itemStackFromNative(nativeStack)
		if !valid {
			return false
		}
		stacks[index] = stack
	}

	accepted, ok := callPlayer(m.players, invocation, id, func(tx *world.Tx, connected *player.Player) bool {
		if dfunsafe.Session(connected) == session.Nop {
			return false
		}
		m.send(tx, connected, value, spec, stacks)
		return true
	})
	return ok && accepted
}

func (m *InventoryMenus) ClosePlayerInventoryMenu(invocation native.InvocationID, id native.PlayerID) bool {
	if m == nil || m.players == nil {
		return false
	}
	closed, ok := callPlayer(m.players, invocation, id, func(tx *world.Tx, connected *player.Player) bool {
		menu := m.take(connected.H())
		if menu == nil {
			return false
		}
		s := dfunsafe.Session(connected)
		connected.MoveItemsToInventory()
		closeInventoryMenuWindow(s, menu.windowID)
		m.restore(tx, s, menu)
		if invocation == 0 {
			m.closeCallback(tx, connected, menu)
		} else {
			m.closeCallbackWithInvocation(invocation, connected, menu)
		}
		return true
	})
	return ok && closed
}

// DiscardPlayerInventoryMenu rolls back one host-accepted menu without
// dispatching its close callback. The ID check avoids closing a newer menu.
func (m *InventoryMenus) DiscardPlayerInventoryMenu(id native.PlayerID, menuID uint64) bool {
	if m == nil || m.players == nil || menuID == 0 {
		return false
	}
	discarded, ok := callPlayer(m.players, 0, id, func(tx *world.Tx, connected *player.Player) bool {
		menu := m.current(connected.H())
		if menu == nil || menu.id != menuID || !m.remove(connected.H(), menu) {
			return false
		}
		s := dfunsafe.Session(connected)
		connected.MoveItemsToInventory()
		closeInventoryMenuWindow(s, menu.windowID)
		m.restore(tx, s, menu)
		return true
	})
	return ok && discarded
}

func (m *InventoryMenus) send(tx *world.Tx, connected *player.Player, value native.PlayerInventoryMenu, spec inventoryMenuSpec, stacks []item.Stack) {
	s := dfunsafe.Session(connected)
	if s == session.Nop {
		return
	}

	previous := m.current(connected.H())
	if previous == nil {
		connected.MoveItemsToInventory()
		inventoryMenuCloseCurrentContainer(s, tx, false)
	}
	update := value.Update && previous != nil && previous.container == value.Container
	var pos cube.Pos
	var sessionPos cube.Pos
	var windowID byte
	var openAt time.Time
	var openSent bool
	if update {
		pos, sessionPos, windowID = previous.pos, previous.sessionPos, previous.windowID
		openAt, openSent = previous.openAt, previous.openSent
	} else {
		if previous != nil {
			m.remove(connected.H(), previous)
			connected.MoveItemsToInventory()
			closeInventoryMenuWindow(s, previous.windowID)
			m.restore(tx, s, previous)
			native.CancelPlayerInventoryMenu(previous.id)
		}
		pos = inventoryMenuPosition(tx, connected)
		sessionPos = inventoryMenuSessionPosition(tx.Range(), pos)
		windowID = nextInventoryMenuWindowID(s)
		openAt = time.Now().Add(inventoryMenuOpenDelay)
	}

	inv := inventory.New(spec.size, nil)
	for index, stack := range stacks {
		_ = inv.SetItem(index, stack)
	}
	menu := &openInventoryMenu{
		id: value.ID, container: value.Container, containerType: spec.containerType,
		inv: inv, pos: pos, sessionPos: sessionPos, windowID: windowID, openAt: openAt, openSent: openSent, double: spec.double,
	}
	inv.Handle(inventoryMenuHandler{menus: m, menuID: value.ID})
	m.store(connected.H(), menu)
	setInventoryMenuSession(s, menu, spec.containerType)
	m.showFakeContainer(s, menu, spec, value.Name)

	if update && menu.openSent {
		inventoryMenuSendInventory(s, menu.inv, uint32(menu.windowID))
		return
	}
	m.sendContentAfter(connected.H(), menu)
}

func (m *InventoryMenus) sendContentAfter(handle *world.EntityHandle, menu *openInventoryMenu) {
	delay := time.Until(menu.openAt)
	if delay < 0 {
		delay = 0
	}
	handle.DoAfter(delay, func(_ *world.Tx, entity world.Entity) {
		connected, ok := entity.(*player.Player)
		if !ok || m.current(handle) != menu {
			return
		}
		s := dfunsafe.Session(connected)
		if s == session.Nop || !inventoryMenuWindowOpen(s, menu.windowID) {
			return
		}
		menu.openSent = true
		dfunsafe.WritePacket(s, &packet.ContainerOpen{
			WindowID:                menu.windowID,
			ContainerType:           menu.containerType,
			ContainerPosition:       inventoryMenuProtocolPos(menu.pos),
			ContainerEntityUniqueID: -1,
		})
		inventoryMenuSendInventory(s, menu.inv, uint32(menu.windowID))
	})
}

func (m *InventoryMenus) showFakeContainer(s *session.Session, menu *openInventoryMenu, spec inventoryMenuSpec, name string) {
	s.ViewBlockUpdate(menu.pos, spec.block, 0)
	s.ViewBlockUpdate(menu.pos.Add(cube.Pos{0, 1, 0}), block.Air{}, 0)
	data := inventoryMenuBlockActorData(spec.block, name, spec.blockActorID, menu.pos, nil, false)
	var pairedData map[string]any
	var paired cube.Pos
	if spec.double {
		paired = menu.pos.Add(cube.Pos{1, 0, 0})
		s.ViewBlockUpdate(paired, spec.block, 0)
		s.ViewBlockUpdate(paired.Add(cube.Pos{0, 1, 0}), block.Air{}, 0)
		data = inventoryMenuBlockActorData(spec.block, name, spec.blockActorID, menu.pos, &paired, true)
		pairedData = inventoryMenuBlockActorData(spec.block, name, spec.blockActorID, paired, &menu.pos, false)
	}
	dfunsafe.WritePacket(s, &packet.BlockActorData{Position: inventoryMenuProtocolPos(menu.pos), NBTData: data})
	if pairedData != nil {
		dfunsafe.WritePacket(s, &packet.BlockActorData{Position: inventoryMenuProtocolPos(paired), NBTData: pairedData})
	}
}

func inventoryMenuBlockActorData(backing world.Block, name, blockActorID string, pos cube.Pos, pair *cube.Pos, pairLead bool) map[string]any {
	data := map[string]any{}
	if nbter, ok := backing.(world.NBTer); ok {
		for key, value := range nbter.EncodeNBT() {
			data[key] = value
		}
	}
	data["CustomName"] = name
	data["id"] = blockActorID
	data["x"] = int32(pos[0])
	data["y"] = int32(pos[1])
	data["z"] = int32(pos[2])
	if blockActorID == "Dropper" {
		data["Items"] = []any{}
	}
	if pair != nil {
		data["pairx"] = int32((*pair)[0])
		data["pairz"] = int32((*pair)[2])
		data["pairlead"] = byte(0)
		if pairLead {
			data["pairlead"] = byte(1)
		}
	}
	return data
}

func (m *InventoryMenus) restore(tx *world.Tx, s *session.Session, menu *openInventoryMenu) {
	positions := []cube.Pos{menu.pos, menu.pos.Add(cube.Pos{0, 1, 0})}
	if menu.double {
		paired := menu.pos.Add(cube.Pos{1, 0, 0})
		positions = append(positions, paired, paired.Add(cube.Pos{0, 1, 0}))
	}
	for _, pos := range positions {
		s.ViewBlockUpdate(pos, tx.Block(pos), 0)
	}
}

// HandleIncomingPacket normalises fake-container stack requests and observes
// client closes before Dragonfly handles the packet.
func (m *InventoryMenus) HandleIncomingPacket(ctx *intercept.Context, value packet.Packet) {
	if m == nil || ctx == nil || ctx.Val() == nil {
		return
	}
	switch value.(type) {
	case *packet.ItemStackRequest, *packet.ContainerClose:
	default:
		return
	}
	handle, ok := ctx.Val().Handle()
	if !ok {
		return
	}
	_, _ = world.CallRef(context.Background(), world.NewEntityRef[*player.Player](handle), func(tx *world.Tx, connected *player.Player) (struct{}, error) {
		menu := m.current(handle)
		if menu == nil {
			return struct{}{}, nil
		}
		s := dfunsafe.Session(connected)
		switch packetValue := value.(type) {
		case *packet.ItemStackRequest:
			normaliseInventoryMenuRequests(packetValue.Requests)
		case *packet.ContainerClose:
			if packetValue.WindowID != menu.windowID || !inventoryMenuWindowOpen(s, menu.windowID) {
				ctx.Cancel()
				return struct{}{}, nil
			}
			// Consume the packet and reproduce Dragonfly's close handling before
			// invoking plugin code. The callback may synchronously open a new menu,
			// which the old packet must never be allowed to act on.
			ctx.Cancel()
			connected.MoveItemsToInventory()
			if !m.remove(handle, menu) {
				return struct{}{}, nil
			}
			acceptInventoryMenuWindowClose(s, menu.windowID)
			m.restore(tx, s, menu)
			m.closeCallback(tx, connected, menu)
		}
		return struct{}{}, nil
	})
}

// HandleOutgoingPacket observes a server-side close caused by Dragonfly, such
// as opening a real block container while a fake menu is active.
func (m *InventoryMenus) HandleOutgoingPacket(ctx *intercept.Context, value packet.Packet) {
	closed, ok := value.(*packet.ContainerClose)
	if !ok || m == nil || ctx == nil || ctx.Val() == nil {
		return
	}
	handle, ok := ctx.Val().Handle()
	if !ok {
		return
	}
	menu := m.current(handle)
	if menu == nil || closed.WindowID != menu.windowID || !m.remove(handle, menu) {
		return
	}
	handle.Do(func(tx *world.Tx, entity world.Entity) {
		connected, ok := entity.(*player.Player)
		if !ok {
			native.CancelPlayerInventoryMenu(menu.id)
			return
		}
		connected.MoveItemsToInventory()
		m.restore(tx, dfunsafe.Session(connected), menu)
		m.closeCallback(tx, connected, menu)
	})
}

func (m *InventoryMenus) closeCallback(tx *world.Tx, connected *player.Player, menu *openInventoryMenu) {
	m.players.WithInvocation(tx, func(invocation native.InvocationID) {
		m.closeCallbackWithInvocation(invocation, connected, menu)
	})
}

func (m *InventoryMenus) closeCallbackWithInvocation(invocation native.InvocationID, connected *player.Player, menu *openInventoryMenu) {
	snapshot, ok := m.playerSnapshot(connected)
	if !ok {
		native.CancelPlayerInventoryMenu(menu.id)
		return
	}
	native.ClosePlayerInventoryMenu(menu.id, invocation, snapshot)
}

func (m *InventoryMenus) submit(connected *player.Player, menuID uint64, stack item.Stack) {
	snapshot, ok := m.playerSnapshot(connected)
	value, valid := itemStackToNative(stack)
	if !ok || !valid {
		return
	}
	m.players.WithInvocation(connected.Tx(), func(invocation native.InvocationID) {
		native.SubmitPlayerInventoryMenu(menuID, invocation, snapshot, value)
	})
}

func (m *InventoryMenus) playerSnapshot(connected *player.Player) (native.PlayerSnapshot, bool) {
	id, ok := m.players.ID(connected)
	if !ok {
		return native.PlayerSnapshot{}, false
	}
	return native.PlayerSnapshot{
		Player: id, Name: connected.Name(),
		LatencyMilliseconds: uint64(max(connected.Latency().Milliseconds(), 0)),
		Position:            native.Vec3{X: connected.Position()[0], Y: connected.Position()[1], Z: connected.Position()[2]},
	}, true
}

func (m *InventoryMenus) disconnectHandle(handle *world.EntityHandle) {
	if m == nil || handle == nil {
		return
	}
	if menu := m.take(handle); menu != nil {
		native.CancelPlayerInventoryMenu(menu.id)
	}
}

// discardAll closes every client window and restores its fake blocks without
// invoking plugin callbacks. Native runtime draining owns callback teardown.
func (m *InventoryMenus) discardAll() {
	if m == nil {
		return
	}
	m.mu.Lock()
	menus := m.menus
	m.menus = map[*world.EntityHandle]*openInventoryMenu{}
	m.mu.Unlock()
	for handle, menu := range menus {
		_, _ = world.CallRef(context.Background(), world.NewEntityRef[*player.Player](handle), func(tx *world.Tx, connected *player.Player) (struct{}, error) {
			s := dfunsafe.Session(connected)
			connected.MoveItemsToInventory()
			closeInventoryMenuWindow(s, menu.windowID)
			m.restore(tx, s, menu)
			return struct{}{}, nil
		})
	}
}

func (m *InventoryMenus) current(handle *world.EntityHandle) *openInventoryMenu {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.menus[handle]
}

func (m *InventoryMenus) store(handle *world.EntityHandle, menu *openInventoryMenu) {
	m.mu.Lock()
	m.menus[handle] = menu
	m.mu.Unlock()
}

func (m *InventoryMenus) take(handle *world.EntityHandle) *openInventoryMenu {
	m.mu.Lock()
	menu := m.menus[handle]
	delete(m.menus, handle)
	m.mu.Unlock()
	return menu
}

func (m *InventoryMenus) remove(handle *world.EntityHandle, menu *openInventoryMenu) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.menus[handle] != menu {
		return false
	}
	delete(m.menus, handle)
	return true
}

type inventoryMenuHandler struct {
	inventory.NopHandler
	menus  *InventoryMenus
	menuID uint64
}

func (h inventoryMenuHandler) HandleTake(ctx *inventory.Context, _ int, stack item.Stack) {
	ctx.Cancel()
	if connected, ok := ctx.Val().(*player.Player); ok {
		h.menus.submit(connected, h.menuID, stack)
	}
}

func (inventoryMenuHandler) HandlePlace(ctx *inventory.Context, _ int, _ item.Stack) { ctx.Cancel() }

func (h inventoryMenuHandler) HandleDrop(ctx *inventory.Context, _ int, stack item.Stack) {
	ctx.Cancel()
	if connected, ok := ctx.Val().(*player.Player); ok {
		h.menus.submit(connected, h.menuID, stack)
	}
}

func inventoryMenuSpecFor(container native.InventoryMenuContainer) (inventoryMenuSpec, bool) {
	switch container {
	case native.InventoryMenuChest, native.InventoryMenuDoubleChest:
		chest := block.NewChest()
		chest.Facing = cube.North
		return inventoryMenuSpec{
			block: chest, blockActorID: "Chest", containerType: protocol.ContainerTypeContainer,
			size:   27 + 27*boolInt(container == native.InventoryMenuDoubleChest),
			double: container == native.InventoryMenuDoubleChest,
		}, true
	case native.InventoryMenuHopper:
		return inventoryMenuSpec{block: block.NewHopper(), blockActorID: "Hopper", containerType: protocol.ContainerTypeHopper, size: 5}, true
	case native.InventoryMenuDropper:
		dropper, ok := world.BlockByName("minecraft:dropper", map[string]any{"facing_direction": int32(0), "toggle_bit": false})
		return inventoryMenuSpec{block: dropper, blockActorID: "Dropper", containerType: protocol.ContainerTypeDropper, size: 9}, ok
	case native.InventoryMenuBarrel:
		return inventoryMenuSpec{block: block.NewBarrel(), blockActorID: "Barrel", containerType: protocol.ContainerTypeContainer, size: 27}, true
	case native.InventoryMenuEnderChest:
		enderChest := block.NewEnderChest()
		enderChest.Facing = cube.North
		return inventoryMenuSpec{block: enderChest, blockActorID: "EnderChest", containerType: protocol.ContainerTypeContainer, size: 27}, true
	default:
		return inventoryMenuSpec{}, false
	}
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func inventoryMenuPosition(tx *world.Tx, connected *player.Player) cube.Pos {
	pos := cube.PosFromVec3(connected.Rotation().Vec3().Mul(-2).Add(connected.Position())).Add(cube.Pos{0, 2, 0})
	height := tx.Range()
	pos[1] = min(max(pos[1], height.Min()), height.Max()-1)
	return pos
}

// inventoryMenuSessionPosition is deliberately outside the dimension height
// range, where Dragonfly resolves it as air. The client still receives pos as
// the visible fake block, but server-side close paths cannot remove a viewer
// from an unrelated real container at that position.
func inventoryMenuSessionPosition(height cube.Range, pos cube.Pos) cube.Pos {
	pos[1] = height.Max() + 1
	return pos
}

func inventoryMenuProtocolPos(pos cube.Pos) protocol.BlockPos {
	return protocol.BlockPos{int32(pos[0]), int32(pos[1]), int32(pos[2])}
}

func normaliseInventoryMenuRequests(requests []protocol.ItemStackRequest) {
	for _, request := range requests {
		for _, action := range request.Actions {
			switch value := action.(type) {
			case *protocol.TakeStackRequestAction:
				normaliseInventoryMenuTransfer(&value.Source, &value.Destination)
			case *protocol.PlaceStackRequestAction:
				normaliseInventoryMenuTransfer(&value.Source, &value.Destination)
			case *protocol.SwapStackRequestAction:
				normaliseInventoryMenuTransfer(&value.Source, &value.Destination)
			case *protocol.DropStackRequestAction:
				if !inventoryMenuPlayerContainer(value.Source.Container.ContainerID) {
					value.Source.Container.ContainerID = protocol.ContainerLevelEntity
				}
			case *protocol.DestroyStackRequestAction:
				if !inventoryMenuPlayerContainer(value.Source.Container.ContainerID) {
					// Destroy bypasses inventory handlers. Force Dragonfly to reject
					// the request instead of mutating this read-only inventory.
					value.Source.Container.ContainerID = inventoryMenuRejectedContainerID
				}
			}
		}
	}
}

func normaliseInventoryMenuTransfer(source, destination *protocol.StackRequestSlotInfo) {
	if !inventoryMenuPlayerContainer(source.Container.ContainerID) {
		source.Container.ContainerID = protocol.ContainerLevelEntity
	}
	if !inventoryMenuPlayerContainer(destination.Container.ContainerID) {
		destination.Container.ContainerID = protocol.ContainerLevelEntity
	}
}

func inventoryMenuPlayerContainer(id byte) bool {
	return id == protocol.ContainerCursor || id == protocol.ContainerHotBar || id == protocol.ContainerInventory || id == protocol.ContainerCombinedHotBarAndInventory
}

func nextInventoryMenuWindowID(s *session.Session) byte {
	opened := inventoryMenuSessionField[atomic.Uint32](s, "openedWindowID")
	if opened.CompareAndSwap(99, 1) {
		return 1
	}
	return byte(opened.Add(1))
}

func setInventoryMenuSession(s *session.Session, menu *openInventoryMenu, containerType byte) {
	inventoryMenuSessionField[atomic.Pointer[cube.Pos]](s, "openedPos").Store(&menu.sessionPos)
	inventoryMenuSessionField[atomic.Pointer[inventory.Inventory]](s, "openedWindow").Store(menu.inv)
	inventoryMenuSessionField[atomic.Uint32](s, "openedWindowID").Store(uint32(menu.windowID))
	inventoryMenuSessionField[atomic.Uint32](s, "openedContainerID").Store(uint32(containerType))
	inventoryMenuSessionField[atomic.Bool](s, "containerOpened").Store(true)
}

func inventoryMenuWindowOpen(s *session.Session, windowID byte) bool {
	return inventoryMenuSessionField[atomic.Bool](s, "containerOpened").Load() &&
		byte(inventoryMenuSessionField[atomic.Uint32](s, "openedWindowID").Load()) == windowID
}

func closeInventoryMenuWindow(s *session.Session, windowID byte) {
	containerType, ok := resetInventoryMenuWindow(s, windowID)
	if !ok {
		return
	}
	dfunsafe.WritePacket(s, &packet.ContainerClose{WindowID: windowID, ContainerType: containerType, ServerSide: true})
}

func acceptInventoryMenuWindowClose(s *session.Session, windowID byte) {
	containerType, ok := resetInventoryMenuWindow(s, windowID)
	if !ok {
		return
	}
	dfunsafe.WritePacket(s, &packet.ContainerClose{WindowID: windowID, ContainerType: containerType})
}

func resetInventoryMenuWindow(s *session.Session, windowID byte) (byte, bool) {
	if s == session.Nop || byte(inventoryMenuSessionField[atomic.Uint32](s, "openedWindowID").Load()) != windowID ||
		!inventoryMenuSessionField[atomic.Bool](s, "containerOpened").CompareAndSwap(true, false) {
		return 0, false
	}
	containerType := byte(inventoryMenuSessionField[atomic.Uint32](s, "openedContainerID").Load())
	inventoryMenuSessionField[atomic.Uint32](s, "openedContainerID").Store(0)
	inventoryMenuSessionField[atomic.Pointer[inventory.Inventory]](s, "openedWindow").Store(inventory.New(1, nil))
	return containerType, true
}

func inventoryMenuSessionField[T any](s *session.Session, name string) *T {
	field := reflect.ValueOf(s).Elem().FieldByName(name)
	if !field.IsValid() || !field.CanAddr() || field.Type() != reflect.TypeFor[T]() {
		panic("dragonfly session field is unavailable: " + name)
	}
	return (*T)(unsafe.Pointer(field.UnsafeAddr()))
}

//go:linkname inventoryMenuCloseCurrentContainer github.com/df-mc/dragonfly/server/session.(*Session).closeCurrentContainer
func inventoryMenuCloseCurrentContainer(*session.Session, *world.Tx, bool)

//go:linkname inventoryMenuSendInventory github.com/df-mc/dragonfly/server/session.(*Session).sendInv
func inventoryMenuSendInventory(*session.Session, *inventory.Inventory, uint32)
