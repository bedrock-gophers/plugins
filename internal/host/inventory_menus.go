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

// The fake-block container flow is adapted from github.com/bedrock-gophers/inv/inv.
const (
	inventoryMenuOpenDelay       = 500 * time.Millisecond
	inventoryMenuCloseAckTimeout = 3 * time.Second
)

type inventoryMenuCloseMode uint8

const (
	inventoryMenuCloseIgnore inventoryMenuCloseMode = iota
	inventoryMenuCloseAcknowledge
)

type InventoryMenus struct {
	players *Players
	mu      sync.Mutex
	active  map[native.PlayerID]*activeInventoryMenu
}

type activeInventoryMenu struct {
	id              uint64
	player          native.PlayerID
	container       native.InventoryMenuContainer
	position        cube.Pos
	windowID        byte
	packetContainer byte
	items           []item.Stack
	inventory       *inventory.Inventory
	opened          bool
	timer           *time.Timer
	closeTimer      *time.Timer
	pending         *pendingInventoryMenu
	double          bool
}

type pendingInventoryMenu struct {
	value native.PlayerInventoryMenu
	info  inventoryMenuContainerInfo
	items []item.Stack
}

type inventoryMenuContainerInfo struct {
	block           world.Block
	blockActorID    string
	packetContainer byte
	size            int
	double          bool
}

type inventoryMenuHandler struct {
	inventory.NopHandler
	menus *InventoryMenus
	menu  *activeInventoryMenu
}

func (handler inventoryMenuHandler) HandleTake(ctx *inventory.Context, slot int, _ item.Stack) {
	ctx.Cancel()
	handler.submit(ctx, slot)
}

func (handler inventoryMenuHandler) HandlePlace(ctx *inventory.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

func (handler inventoryMenuHandler) HandleDrop(ctx *inventory.Context, slot int, _ item.Stack) {
	ctx.Cancel()
	handler.submit(ctx, slot)
}

func (handler inventoryMenuHandler) submit(ctx *inventory.Context, slot int) {
	connected, ok := ctx.Val().(*player.Player)
	if !ok || handler.menus == nil || handler.menu == nil || slot < 0 || slot >= len(handler.menu.items) {
		return
	}
	handler.menus.mu.Lock()
	active := handler.menus.active[handler.menu.player] == handler.menu
	handler.menus.mu.Unlock()
	if !active {
		return
	}
	invocation, leave := handler.menus.players.BeginInvocation(connected.Tx())
	defer leave()
	native.CompletePlayerInventoryMenuClick(
		handler.menu.id,
		invocation,
		inventoryMenuPlayerSnapshot(connected, handler.menu.player),
		uint32(slot),
	)
	connected.Tx().Defer(func(*world.Tx) {
		syncInventoryMenuPlayerInventories(connected)
	})
}

func NewInventoryMenus(players *Players) *InventoryMenus {
	return &InventoryMenus{players: players, active: map[native.PlayerID]*activeInventoryMenu{}}
}

func (menus *InventoryMenus) SendPlayerInventoryMenu(invocation native.InvocationID, id native.PlayerID, value native.PlayerInventoryMenu) bool {
	if menus == nil || menus.players == nil || value.ID == 0 {
		return false
	}
	info, ok := resolveInventoryMenuContainer(value.Container)
	if !ok || value.Title == "" || len(value.Items) != info.size {
		return false
	}
	items := make([]item.Stack, len(value.Items))
	for index, encoded := range value.Items {
		stack, valid := ItemStackFromNative(encoded)
		if !valid {
			return false
		}
		items[index] = stack
	}
	if handled, accepted := menus.deferInventoryMenuReplacement(invocation, id, value, info, items); handled {
		return accepted
	}
	var replaced uint64
	var accepted bool
	resolved := menus.players.mutatePlayer(invocation, id, func(connected *player.Player) {
		replaced, accepted = menus.send(connected, id, value, info, items)
	})
	if replaced != 0 && replaced != value.ID {
		native.CancelPlayerInventoryMenu(replaced)
	}
	return resolved && accepted
}

func (menus *InventoryMenus) deferInventoryMenuReplacement(invocation native.InvocationID, id native.PlayerID, value native.PlayerInventoryMenu, info inventoryMenuContainerInfo, items []item.Stack) (bool, bool) {
	if invocation == 0 {
		return false, false
	}
	menus.mu.Lock()
	previous := menus.active[id]
	deferRequired := inventoryMenuDefersReplacement(invocation != 0, previous != nil && previous.opened)
	menus.mu.Unlock()
	if !deferRequired {
		return false, false
	}
	if !menus.players.OnInvocationEnd(invocation, func() {
		menus.sendDeferredInventoryMenu(id, value, info, items)
	}) {
		return true, false
	}
	return true, true
}

func (menus *InventoryMenus) sendDeferredInventoryMenu(id native.PlayerID, value native.PlayerInventoryMenu, info inventoryMenuContainerInfo, items []item.Stack) {
	entry, ok := menus.players.playerEntry(id)
	if !ok {
		native.CancelPlayerInventoryMenu(value.ID)
		return
	}
	task := world.NewEntityRef[*player.Player](entry.handle).Do(func(_ *world.Tx, connected *player.Player) {
		replaced, accepted := menus.send(connected, id, value, info, items)
		if replaced != 0 && replaced != value.ID {
			native.CancelPlayerInventoryMenu(replaced)
		}
		if !accepted {
			native.CancelPlayerInventoryMenu(value.ID)
		}
	})
	task.OnDone(func(err error) {
		if err != nil {
			native.CancelPlayerInventoryMenu(value.ID)
		}
	})
}

func (menus *InventoryMenus) ClosePlayerInventoryMenu(invocation native.InvocationID, id native.PlayerID) bool {
	if menus == nil || menus.players == nil {
		return false
	}
	closed := false
	return menus.players.mutatePlayer(invocation, id, func(connected *player.Player) {
		menu, callbackID := menus.remove(connected, id, true)
		if menu == nil {
			return
		}
		closed = true
		native.CompletePlayerInventoryMenuClose(callbackID, invocation, inventoryMenuPlayerSnapshot(connected, id))
	}) && closed
}

func (menus *InventoryMenus) send(connected *player.Player, id native.PlayerID, value native.PlayerInventoryMenu, info inventoryMenuContainerInfo, items []item.Stack) (uint64, bool) {
	s := dfunsafe.Session(connected)
	if s == session.Nop {
		return 0, false
	}
	inv := inventory.New(info.size, nil)
	for index, stack := range items {
		if err := inv.SetItem(index, stack); err != nil {
			return 0, false
		}
	}

	menus.mu.Lock()
	previous := menus.active[id]
	if previous != nil && previous.closing() {
		replaced := pendingInventoryMenuID(previous.pending)
		previous.pending = &pendingInventoryMenu{value: value, info: info, items: items}
		menus.mu.Unlock()
		return replaced, true
	}
	if previous != nil && inventoryMenuReusesWindow(value.Update, previous.opened, previous.container == value.Container) {
		replaced, replacedItems, replacedInventory := previous.id, previous.items, previous.inventory
		replacedPacketContainer := previous.packetContainer
		previous.id = value.ID
		previous.items = items
		previous.inventory = inv
		previous.packetContainer = info.packetContainer
		inv.Handle(inventoryMenuHandler{menus: menus, menu: previous})
		if !setInventoryMenuSession(s, previous) {
			previous.id, previous.items, previous.inventory = replaced, replacedItems, replacedInventory
			previous.packetContainer = replacedPacketContainer
			menus.mu.Unlock()
			return 0, false
		}
		menus.sendBlockActorData(s, previous, value.Title, info.blockActorID)
		if previous.opened {
			for slot, stack := range previous.items {
				s.ViewSlotChange(slot, stack)
			}
		}
		menus.mu.Unlock()
		return replaced, true
	}
	if previous != nil && previous.opened {
		previous.pending = &pendingInventoryMenu{value: value, info: info, items: items}
		menus.closeLocked(s, connected.Tx(), previous, true, true)
		previous.closeTimer = time.AfterFunc(inventoryMenuCloseAckTimeout, func() {
			menus.expirePendingInventoryMenu(id, previous)
		})
		menus.mu.Unlock()
		return previous.id, true
	}

	if previous != nil {
		menus.closeLocked(s, connected.Tx(), previous, true, true)
		delete(menus.active, id)
	} else {
		closeExistingSessionContainer(s, connected.Tx())
	}
	windowID, ok := nextInventoryMenuWindowID(s)
	if !ok {
		menus.mu.Unlock()
		return menuID(previous), false
	}
	position := inventoryMenuPosition(connected)
	menu := &activeInventoryMenu{
		id: value.ID, player: id, container: value.Container, position: position,
		windowID: windowID, packetContainer: info.packetContainer,
		items: items, inventory: inv, double: info.double,
	}
	inv.Handle(inventoryMenuHandler{menus: menus, menu: menu})
	if !setInventoryMenuSession(s, menu) {
		menus.mu.Unlock()
		return menuID(previous), false
	}
	menus.active[id] = menu
	menus.showFakeBlocks(s, menu, info, value.Title)
	menu.timer = time.AfterFunc(inventoryMenuOpenDelay, func() { menus.finishOpen(id, menu, s) })
	menus.mu.Unlock()
	return menuID(previous), true
}

func inventoryMenuPosition(connected *player.Player) cube.Pos {
	return cube.PosFromVec3(connected.Rotation().Vec3().Mul(-2).Add(connected.Position())).Add(cube.Pos{0, 2, 0})
}

func (menus *InventoryMenus) finishOpen(id native.PlayerID, menu *activeInventoryMenu, s *session.Session) {
	menus.mu.Lock()
	defer menus.mu.Unlock()
	if menus.active[id] != menu || menu.timer == nil {
		return
	}
	menu.timer = nil
	menu.opened = true
	menus.sendOpen(s, menu)
}

func (menus *InventoryMenus) showFakeBlocks(s *session.Session, menu *activeInventoryMenu, info inventoryMenuContainerInfo, title string) {
	s.ViewBlockUpdate(menu.position, info.block, 0)
	s.ViewBlockUpdate(menu.position.Add(cube.Pos{0, 1, 0}), block.Air{}, 0)
	if menu.double {
		s.ViewBlockUpdate(menu.position.Add(cube.Pos{1, 0, 0}), info.block, 0)
		s.ViewBlockUpdate(menu.position.Add(cube.Pos{1, 1, 0}), block.Air{}, 0)
	}
	menus.sendBlockActorData(s, menu, title, info.blockActorID)
}

func (menus *InventoryMenus) sendBlockActorData(s *session.Session, menu *activeInventoryMenu, title, blockActorID string) {
	data := map[string]any{"id": blockActorID, "CustomName": title}
	if menu.double {
		data["pairx"] = int32(menu.position[0] + 1)
		data["pairz"] = int32(menu.position[2])
	}
	dfunsafe.WritePacket(s, &packet.BlockActorData{
		Position: protocol.BlockPos{int32(menu.position[0]), int32(menu.position[1]), int32(menu.position[2])},
		NBTData:  data,
	})
}

func (menus *InventoryMenus) sendOpen(s *session.Session, menu *activeInventoryMenu) {
	dfunsafe.WritePacket(s, &packet.ContainerOpen{
		WindowID: menu.windowID, ContainerType: menu.packetContainer,
		ContainerPosition:       protocol.BlockPos{int32(menu.position[0]), int32(menu.position[1]), int32(menu.position[2])},
		ContainerEntityUniqueID: -1,
	})
	sendInventoryMenuContent(s, menu.inventory, uint32(menu.windowID))
}

//go:linkname sendInventoryMenuContent github.com/df-mc/dragonfly/server/session.(*Session).sendInv
func sendInventoryMenuContent(*session.Session, *inventory.Inventory, uint32)

func syncInventoryMenuPlayerInventories(connected *player.Player) {
	s := dfunsafe.Session(connected)
	if s == session.Nop {
		return
	}
	sendInventoryMenuContent(s, connected.Inventory(), protocol.WindowIDInventory)
	sendInventoryMenuContent(s, connected.Armour().Inventory(), protocol.WindowIDArmour)
	_, offhand := connected.HeldItems()
	offhandInventory := inventory.New(1, nil)
	_ = offhandInventory.SetItem(0, offhand)
	sendInventoryMenuContent(s, offhandInventory, protocol.WindowIDOffHand)
}

func (menus *InventoryMenus) remove(connected *player.Player, id native.PlayerID, sendClose bool) (*activeInventoryMenu, uint64) {
	s := dfunsafe.Session(connected)
	if s == session.Nop {
		return nil, 0
	}
	menus.mu.Lock()
	menu := menus.active[id]
	callbackID := menuID(menu)
	if menu != nil {
		delete(menus.active, id)
		if menu.pending != nil {
			callbackID = menu.pending.value.ID
			menu.pending = nil
		}
		menus.closeLocked(s, connected.Tx(), menu, sendClose, true)
	}
	menus.mu.Unlock()
	return menu, callbackID
}

func (menus *InventoryMenus) closeLocked(s *session.Session, tx *world.Tx, menu *activeInventoryMenu, sendClose, restoreBlocks bool) {
	if menu.timer != nil {
		menu.timer.Stop()
		menu.timer = nil
	}
	if menu.closeTimer != nil {
		menu.closeTimer.Stop()
		menu.closeTimer = nil
	}
	clearInventoryMenuSession(s, menu.windowID)
	if menu.opened && sendClose {
		dfunsafe.WritePacket(s, &packet.ContainerClose{
			WindowID: menu.windowID, ContainerType: menu.packetContainer, ServerSide: true,
		})
	}
	if restoreBlocks && tx != nil {
		s.ViewBlockUpdate(menu.position, tx.Block(menu.position), 0)
		above := menu.position.Add(cube.Pos{0, 1, 0})
		s.ViewBlockUpdate(above, tx.Block(above), 0)
		if menu.double {
			paired := menu.position.Add(cube.Pos{1, 0, 0})
			s.ViewBlockUpdate(paired, tx.Block(paired), 0)
			pairedAbove := menu.position.Add(cube.Pos{1, 1, 0})
			s.ViewBlockUpdate(pairedAbove, tx.Block(pairedAbove), 0)
		}
	}
}

func (menus *InventoryMenus) expirePendingInventoryMenu(id native.PlayerID, menu *activeInventoryMenu) {
	menus.mu.Lock()
	if menus.active[id] != menu || !menu.closing() {
		menus.mu.Unlock()
		return
	}
	pending := menu.pending
	menu.pending = nil
	menu.closeTimer = nil
	delete(menus.active, id)
	menus.mu.Unlock()
	if pending != nil {
		native.CancelPlayerInventoryMenu(pending.value.ID)
	}
}

func (menu *activeInventoryMenu) closing() bool {
	return menu != nil && menu.pending != nil
}

func inventoryMenuReusesWindow(update, previousOpened, sameContainer bool) bool {
	return sameContainer && (update || previousOpened)
}

func inventoryMenuDefersReplacement(hasInvocation, previousOpened bool) bool {
	return hasInvocation && previousOpened
}

func (menus *InventoryMenus) HandleIncomingPacket(ctx *intercept.Context, value packet.Packet) {
	if menus == nil || ctx == nil || value == nil || ctx.Val() == nil {
		return
	}
	if _, relevant := value.(*packet.ItemStackRequest); !relevant {
		if _, relevant = value.(*packet.ContainerClose); !relevant {
			return
		}
	}
	handle, ok := ctx.Val().Handle()
	if !ok || handle == nil {
		return
	}
	_, _ = world.CallRef(context.Background(), world.NewEntityRef[world.Entity](handle), func(tx *world.Tx, entity world.Entity) (struct{}, error) {
		connected, ok := entity.(*player.Player)
		if !ok {
			return struct{}{}, nil
		}
		id, ok := menus.players.ID(connected)
		if !ok {
			return struct{}{}, nil
		}
		menus.mu.Lock()
		menu := menus.active[id]
		menus.mu.Unlock()
		if menu == nil {
			return struct{}{}, nil
		}
		switch current := value.(type) {
		case *packet.ItemStackRequest:
			if !menu.closing() {
				rewriteInventoryMenuRequest(current, len(menu.items))
			}
		case *packet.ContainerClose:
			menus.handleContainerClose(ctx, tx, connected, menu, current)
		}
		return struct{}{}, nil
	})
}

func rewriteInventoryMenuRequest(request *packet.ItemStackRequest, size int) {
	for _, current := range request.Requests {
		for _, action := range current.Actions {
			switch value := action.(type) {
			case *protocol.TakeStackRequestAction:
				rewriteInventoryMenuSlot(&value.Source, size)
				rewriteInventoryMenuSlot(&value.Destination, size)
			case *protocol.PlaceStackRequestAction:
				rewriteInventoryMenuSlot(&value.Source, size)
				rewriteInventoryMenuSlot(&value.Destination, size)
			case *protocol.SwapStackRequestAction:
				rewriteInventoryMenuSlot(&value.Source, size)
				rewriteInventoryMenuSlot(&value.Destination, size)
			case *protocol.PlaceInContainerStackRequestAction:
				rewriteInventoryMenuSlot(&value.Source, size)
				rewriteInventoryMenuSlot(&value.Destination, size)
			case *protocol.TakeOutContainerStackRequestAction:
				rewriteInventoryMenuSlot(&value.Source, size)
				rewriteInventoryMenuSlot(&value.Destination, size)
			case *protocol.DropStackRequestAction:
				rewriteInventoryMenuSlot(&value.Source, size)
			case *protocol.DestroyStackRequestAction:
				rewriteInventoryMenuSlot(&value.Source, size)
			case *protocol.ConsumeStackRequestAction:
				rewriteInventoryMenuSlot(&value.Source, size)
			}
		}
	}
}

func rewriteInventoryMenuSlot(slot *protocol.StackRequestSlotInfo, size int) {
	if slot == nil || int(slot.Slot) >= size {
		return
	}
	switch slot.Container.ContainerID {
	case protocol.ContainerInventory, protocol.ContainerHotBar, protocol.ContainerCombinedHotBarAndInventory,
		protocol.ContainerCursor, protocol.ContainerOffhand, protocol.ContainerArmor:
		return
	default:
		slot.Container.ContainerID = protocol.ContainerLevelEntity
	}
}

func (menus *InventoryMenus) handleContainerClose(ctx *intercept.Context, tx *world.Tx, connected *player.Player, menu *activeInventoryMenu, request *packet.ContainerClose) {
	mode := classifyInventoryMenuClose(request.WindowID, menu.windowID, menu.closing())
	if mode == inventoryMenuCloseIgnore {
		ctx.Cancel()
		return
	}
	ctx.Cancel()
	menus.mu.Lock()
	if menus.active[menu.player] != menu {
		menus.mu.Unlock()
		return
	}
	if menu.closing() {
		pending := menu.pending
		menu.pending = nil
		if menu.closeTimer != nil {
			menu.closeTimer.Stop()
			menu.closeTimer = nil
		}
		delete(menus.active, menu.player)
		menus.mu.Unlock()
		dfunsafe.WritePacket(connected, &packet.ContainerClose{
			WindowID: menu.windowID, ContainerType: menu.packetContainer,
		})
		if pending != nil {
			replaced, accepted := menus.send(connected, menu.player, pending.value, pending.info, pending.items)
			if replaced != 0 && replaced != pending.value.ID {
				native.CancelPlayerInventoryMenu(replaced)
			}
			if !accepted {
				native.CancelPlayerInventoryMenu(pending.value.ID)
			}
		}
		return
	}
	delete(menus.active, menu.player)
	menus.closeLocked(dfunsafe.Session(connected), tx, menu, false, true)
	menus.mu.Unlock()
	if mode == inventoryMenuCloseAcknowledge {
		dfunsafe.WritePacket(connected, &packet.ContainerClose{WindowID: menu.windowID, ContainerType: menu.packetContainer})
	}
	invocation, leave := menus.players.BeginInvocation(tx)
	defer leave()
	native.CompletePlayerInventoryMenuClose(menu.id, invocation, inventoryMenuPlayerSnapshot(connected, menu.player))
}

func classifyInventoryMenuClose(requestWindow, activeWindow byte, awaitingAck bool) inventoryMenuCloseMode {
	if requestWindow == activeWindow {
		return inventoryMenuCloseAcknowledge
	}
	if awaitingAck && requestWindow == 0xff {
		return inventoryMenuCloseAcknowledge
	}
	// Bedrock uses 0xff while transitioning out of chat. It does not identify the
	// active menu and must not tear down its pending fake container.
	return inventoryMenuCloseIgnore
}

func (menus *InventoryMenus) Disconnect(tx *world.Tx, connected *player.Player) {
	if menus == nil || connected == nil {
		return
	}
	id, ok := menus.players.ID(connected)
	if !ok {
		return
	}
	menus.mu.Lock()
	menu := menus.active[id]
	callbackID := menuID(menu)
	if menu != nil {
		delete(menus.active, id)
		if menu.pending != nil {
			callbackID = menu.pending.value.ID
			menu.pending = nil
		}
		menus.closeLocked(dfunsafe.Session(connected), tx, menu, false, false)
	}
	menus.mu.Unlock()
	if menu != nil {
		invocation, leave := menus.players.BeginInvocation(tx)
		defer leave()
		native.CompletePlayerInventoryMenuClose(callbackID, invocation, inventoryMenuPlayerSnapshot(connected, id))
	}
}

func resolveInventoryMenuContainer(container native.InventoryMenuContainer) (inventoryMenuContainerInfo, bool) {
	switch container {
	case native.InventoryMenuChest, native.InventoryMenuDoubleChest:
		chest := block.NewChest()
		chest.Facing = 1
		size := 27
		double := container == native.InventoryMenuDoubleChest
		if double {
			size = 54
		}
		return inventoryMenuContainerInfo{
			block:           chest,
			blockActorID:    "Chest",
			packetContainer: protocol.ContainerTypeContainer,
			size:            size, double: double,
		}, true
	case native.InventoryMenuHopper:
		return inventoryMenuContainerInfo{
			block:           block.NewHopper(),
			blockActorID:    "Hopper",
			packetContainer: protocol.ContainerTypeHopper,
			size:            5,
		}, true
	case native.InventoryMenuDropper:
		dropper, ok := world.BlockByName("minecraft:dropper", map[string]any{"facing_direction": int32(0), "toggle_bit": false})
		if !ok {
			return inventoryMenuContainerInfo{}, false
		}
		return inventoryMenuContainerInfo{
			block:           dropper,
			blockActorID:    "Dropper",
			packetContainer: protocol.ContainerTypeDropper,
			size:            9,
		}, true
	case native.InventoryMenuBarrel:
		return inventoryMenuContainerInfo{
			block:           block.NewBarrel(),
			blockActorID:    "Barrel",
			packetContainer: protocol.ContainerTypeContainer,
			size:            27,
		}, true
	case native.InventoryMenuEnderChest:
		enderChest := block.NewEnderChest()
		enderChest.Facing = 1
		return inventoryMenuContainerInfo{
			block:           enderChest,
			blockActorID:    "EnderChest",
			packetContainer: protocol.ContainerTypeContainer,
			size:            27,
		}, true
	default:
		return inventoryMenuContainerInfo{}, false
	}
}

func inventoryMenuPlayerSnapshot(connected *player.Player, id native.PlayerID) native.PlayerSnapshot {
	position := connected.Position()
	return native.PlayerSnapshot{
		Player: id, Name: connected.Name(),
		LatencyMilliseconds: uint64(max(connected.Latency().Milliseconds(), 0)),
		Position:            native.Vec3{X: position.X(), Y: position.Y(), Z: position.Z()},
	}
}

func menuID(menu *activeInventoryMenu) uint64 {
	if menu == nil {
		return 0
	}
	return menu.id
}

func pendingInventoryMenuID(menu *pendingInventoryMenu) uint64 {
	if menu == nil {
		return 0
	}
	return menu.value.ID
}

func nextInventoryMenuWindowID(s *session.Session) (byte, bool) {
	window, ok := privateSessionField[atomic.Uint32](s, "openedWindowID")
	if !ok {
		return 0, false
	}
	if window.CompareAndSwap(99, 1) {
		return 1, true
	}
	return byte(window.Add(1)), true
}

func setInventoryMenuSession(s *session.Session, menu *activeInventoryMenu) bool {
	opened, openedOK := privateSessionField[atomic.Bool](s, "containerOpened")
	windowID, windowIDOK := privateSessionField[atomic.Uint32](s, "openedWindowID")
	containerID, containerIDOK := privateSessionField[atomic.Uint32](s, "openedContainerID")
	window, windowOK := privateSessionField[atomic.Pointer[inventory.Inventory]](s, "openedWindow")
	position, positionOK := privateSessionField[atomic.Pointer[cube.Pos]](s, "openedPos")
	if !openedOK || !windowIDOK || !containerIDOK || !windowOK || !positionOK {
		return false
	}
	window.Store(menu.inventory)
	position.Store(&menu.position)
	windowID.Store(uint32(menu.windowID))
	containerID.Store(uint32(menu.packetContainer))
	opened.Store(true)
	return true
}

func clearInventoryMenuSession(s *session.Session, expectedWindow byte) {
	opened, openedOK := privateSessionField[atomic.Bool](s, "containerOpened")
	windowID, windowIDOK := privateSessionField[atomic.Uint32](s, "openedWindowID")
	containerID, containerIDOK := privateSessionField[atomic.Uint32](s, "openedContainerID")
	window, windowOK := privateSessionField[atomic.Pointer[inventory.Inventory]](s, "openedWindow")
	if !openedOK || !windowIDOK || !containerIDOK || !windowOK || byte(windowID.Load()) != expectedWindow {
		return
	}
	opened.Store(false)
	containerID.Store(0)
	window.Store(inventory.New(1, nil))
}

func closeExistingSessionContainer(s *session.Session, tx *world.Tx) {
	opened, openedOK := privateSessionField[atomic.Bool](s, "containerOpened")
	windowID, windowIDOK := privateSessionField[atomic.Uint32](s, "openedWindowID")
	containerID, containerIDOK := privateSessionField[atomic.Uint32](s, "openedContainerID")
	window, windowOK := privateSessionField[atomic.Pointer[inventory.Inventory]](s, "openedWindow")
	position, positionOK := privateSessionField[atomic.Pointer[cube.Pos]](s, "openedPos")
	if !openedOK || !windowIDOK || !containerIDOK || !windowOK || !positionOK || !opened.CompareAndSwap(true, false) {
		return
	}
	containerType := byte(containerID.Load())
	currentWindow := byte(windowID.Load())
	containerID.Store(0)
	window.Store(inventory.New(1, nil))
	dfunsafe.WritePacket(s, &packet.ContainerClose{
		WindowID: currentWindow, ContainerType: containerType, ServerSide: true,
	})
	if tx == nil || position.Load() == nil {
		return
	}
	pos := *position.Load()
	switch current := tx.Block(pos).(type) {
	case block.Container:
		current.RemoveViewer(s, tx, pos)
	case block.EnderChest:
		current.RemoveViewer(tx, pos)
	}
}

func privateSessionField[T any](s *session.Session, name string) (field *T, ok bool) {
	if s == nil {
		return nil, false
	}
	defer func() {
		if recover() != nil {
			field, ok = nil, false
		}
	}()
	value := reflect.ValueOf(s).Elem().FieldByName(name)
	if !value.IsValid() || !value.CanAddr() || value.Type() != reflect.TypeOf(*new(T)) {
		return nil, false
	}
	return (*T)(unsafe.Pointer(value.UnsafeAddr())), true
}
