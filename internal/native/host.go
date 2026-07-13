package native

import (
	"sync"
	"sync/atomic"
	"time"
)

type PlayerForm struct {
	ID          uint64
	RequestJSON []byte
}

type PlayerTransformKind uint32

const (
	PlayerTransformTeleport PlayerTransformKind = iota
	PlayerTransformMove
	PlayerTransformVelocity
)

type PlayerTitle struct {
	Text       string
	Subtitle   string
	ActionText string
	FadeIn     time.Duration
	Duration   time.Duration
	FadeOut    time.Duration
}

type PlayerScoreboard struct {
	Name       string
	Lines      []string
	Padding    bool
	Descending bool
}

type SkinAnimation struct {
	Width, Height uint32
	Type          uint32
	FrameCount    int64
	Expression    int64
	Pixels        []byte
}

type PlayerSkin struct {
	Width, Height                   uint32
	Persona                         bool
	PlayFabID, FullID               string
	Pixels                          []byte
	ModelDefault, ModelAnimatedFace string
	Model                           []byte
	CapeWidth, CapeHeight           uint32
	CapePixels                      []byte
	Animations                      []SkinAnimation
}

type InventoryKind uint32

const (
	InventoryMain InventoryKind = iota
	InventoryArmour
	InventoryOffhand
)

type InventoryID struct {
	Player PlayerID
	Kind   InventoryKind
}

type ItemEnchantment struct {
	ID, Level uint32
}

type ItemStack struct {
	Identifier     string
	Metadata       int32
	Count, Damage  uint32
	Unbreakable    bool
	AnvilCost      int32
	CustomName     string
	Lore           []string
	NBT, ValuesNBT []byte
	Enchantments   []ItemEnchantment
}

// Host executes synchronous actions requested by native plugins.
type Host interface {
	SendPlayerText(PlayerID, PlayerTextKind, string) bool
	SendPlayerTitle(PlayerID, PlayerTitle) bool
	SendPlayerScoreboard(PlayerID, PlayerScoreboard) bool
	RemovePlayerScoreboard(PlayerID) bool
	SendPlayerForm(PlayerID, PlayerForm) bool
	ClosePlayerForm(PlayerID) bool
	TransformPlayer(PlayerID, PlayerTransformKind, Vec3, float64, float64) bool
	PlayerRotation(PlayerID) (Rotation, bool)
	SetPlayerState(PlayerID, PlayerStateKind, PlayerStateValue) bool
	PlayerState(PlayerID, PlayerStateKind) (PlayerStateValue, bool)
	ChangePlayerEffect(PlayerID, PlayerEffectOperation, PlayerEffect) bool
	SetPlayerEntityVisible(PlayerID, EntityID, bool) bool
	PlayerSkin(PlayerID) (PlayerSkin, bool)
	SetPlayerSkin(PlayerID, PlayerSkin) bool
	InventorySize(InventoryID) (uint32, bool)
	InventoryItem(InventoryID, uint32) (ItemStack, bool)
	SetInventoryItem(InventoryID, uint32, ItemStack) bool
	AddInventoryItem(InventoryID, ItemStack) (uint32, bool)
	ClearInventory(InventoryID) bool
	HeldItem(PlayerID, uint32) (ItemStack, bool)
	SetHeldItems(PlayerID, ItemStack, ItemStack) bool
	SetHeldSlot(PlayerID, uint32) bool
}

type noopHost struct{}

func (noopHost) SendPlayerText(PlayerID, PlayerTextKind, string) bool { return false }
func (noopHost) SendPlayerTitle(PlayerID, PlayerTitle) bool           { return false }
func (noopHost) SendPlayerScoreboard(PlayerID, PlayerScoreboard) bool { return false }
func (noopHost) RemovePlayerScoreboard(PlayerID) bool                 { return false }
func (noopHost) SendPlayerForm(PlayerID, PlayerForm) bool             { return false }
func (noopHost) ClosePlayerForm(PlayerID) bool                        { return false }
func (noopHost) TransformPlayer(PlayerID, PlayerTransformKind, Vec3, float64, float64) bool {
	return false
}
func (noopHost) PlayerRotation(PlayerID) (Rotation, bool)                        { return Rotation{}, false }
func (noopHost) SetPlayerState(PlayerID, PlayerStateKind, PlayerStateValue) bool { return false }
func (noopHost) PlayerState(PlayerID, PlayerStateKind) (PlayerStateValue, bool) {
	return PlayerStateValue{}, false
}
func (noopHost) ChangePlayerEffect(PlayerID, PlayerEffectOperation, PlayerEffect) bool { return false }
func (noopHost) SetPlayerEntityVisible(PlayerID, EntityID, bool) bool                  { return false }
func (noopHost) PlayerSkin(PlayerID) (PlayerSkin, bool)                                { return PlayerSkin{}, false }
func (noopHost) SetPlayerSkin(PlayerID, PlayerSkin) bool                               { return false }
func (noopHost) InventorySize(InventoryID) (uint32, bool)                              { return 0, false }
func (noopHost) InventoryItem(InventoryID, uint32) (ItemStack, bool)                   { return ItemStack{}, false }
func (noopHost) SetInventoryItem(InventoryID, uint32, ItemStack) bool                  { return false }
func (noopHost) AddInventoryItem(InventoryID, ItemStack) (uint32, bool)                { return 0, false }
func (noopHost) ClearInventory(InventoryID) bool                                       { return false }
func (noopHost) HeldItem(PlayerID, uint32) (ItemStack, bool)                           { return ItemStack{}, false }
func (noopHost) SetHeldItems(PlayerID, ItemStack, ItemStack) bool                      { return false }
func (noopHost) SetHeldSlot(PlayerID, uint32) bool                                     { return false }

var (
	hostSequence         atomic.Uint64
	hosts                sync.Map
	skinSnapshotSequence atomic.Uint64
	skinSnapshotMu       sync.Mutex
	skinSnapshots        = map[uint64]skinSnapshot{}
	skinSnapshotCounts   = map[uint64]int{}
	itemSnapshotSequence atomic.Uint64
	itemSnapshotMu       sync.Mutex
	itemSnapshots        = map[uint64]itemSnapshot{}
	itemSnapshotCounts   = map[uint64]int{}
	formSequence         atomic.Uint64
	formMu               sync.Mutex
	formCond             = sync.NewCond(&formMu)
	forms                = map[uint64]formRegistration{}
	formHostState        = map[uint64]*formState{}
)

const maxSkinSnapshotsPerHost = 32
const maxItemSnapshotsPerHost = 64
const maxFormsPerHost = 128
const maxFormsPerPlayer = 16

type formRegistration struct {
	host    uint64
	player  PlayerID
	respond func(PlayerID, bool, []byte) bool
	drop    func()
}
type formState struct {
	closing, draining bool
	inflight, count   int
	players           map[PlayerID]int
}

func registerForm(host uint64, player PlayerID, respond func(PlayerID, bool, []byte) bool, drop func()) (uint64, bool) {
	formMu.Lock()
	defer formMu.Unlock()
	state := formHostState[host]
	if state == nil {
		state = &formState{players: map[PlayerID]int{}}
		formHostState[host] = state
	}
	if state.closing || state.draining || state.count >= maxFormsPerHost || state.players[player] >= maxFormsPerPlayer {
		return 0, false
	}
	id := formSequence.Add(1)
	forms[id] = formRegistration{host: host, player: player, respond: respond, drop: drop}
	state.count++
	state.players[player]++
	return id, true
}

func CompletePlayerForm(id uint64, submitter PlayerID, closed bool, response []byte) bool {
	formMu.Lock()
	registration, ok := forms[id]
	if !ok {
		formMu.Unlock()
		return true
	}
	delete(forms, id)
	state := formHostState[registration.host]
	state.count--
	state.players[registration.player]--
	if state.players[registration.player] == 0 {
		delete(state.players, registration.player)
	}
	state.inflight++
	formMu.Unlock()
	if submitter != registration.player || len(response) > maxFormJSONBytes {
		registration.drop()
		formMu.Lock()
		state.inflight--
		formCond.Broadcast()
		formMu.Unlock()
		return false
	}
	ok = registration.respond(submitter, closed, response)
	formMu.Lock()
	state.inflight--
	formCond.Broadcast()
	formMu.Unlock()
	return ok
}

func abandonForm(id uint64) {
	formMu.Lock()
	defer formMu.Unlock()
	r, ok := forms[id]
	if !ok {
		return
	}
	delete(forms, id)
	s := formHostState[r.host]
	s.count--
	s.players[r.player]--
	if s.players[r.player] == 0 {
		delete(s.players, r.player)
	}
}

func CancelPlayerForms(player PlayerID) {
	cancelMatchingForms(func(r formRegistration) bool { return r.player == player })
}
func CancelPlayerForm(id uint64) {
	var callback func()
	var state *formState
	formMu.Lock()
	if r, ok := forms[id]; ok {
		delete(forms, id)
		state = formHostState[r.host]
		state.count--
		state.players[r.player]--
		if state.players[r.player] == 0 {
			delete(state.players, r.player)
		}
		state.inflight++
		callback = r.drop
	}
	formMu.Unlock()
	if callback != nil {
		callback()
		formMu.Lock()
		state.inflight--
		formCond.Broadcast()
		formMu.Unlock()
	}
}
func cancelMatchingForms(match func(formRegistration) bool) {
	type pendingDrop struct {
		callback func()
		state    *formState
	}
	var callbacks []pendingDrop
	formMu.Lock()
	for id, r := range forms {
		if match(r) {
			delete(forms, id)
			s := formHostState[r.host]
			s.count--
			s.players[r.player]--
			if s.players[r.player] == 0 {
				delete(s.players, r.player)
			}
			s.inflight++
			callbacks = append(callbacks, pendingDrop{callback: r.drop, state: s})
		}
	}
	formMu.Unlock()
	for _, callback := range callbacks {
		callback.callback()
		formMu.Lock()
		callback.state.inflight--
		formCond.Broadcast()
		formMu.Unlock()
	}
}

func drainHostForms(host uint64, closing bool) {
	formMu.Lock()
	state := formHostState[host]
	if state == nil {
		if closing {
			formHostState[host] = &formState{closing: true, draining: true, players: map[PlayerID]int{}}
		}
		formMu.Unlock()
		return
	}
	state.draining = true
	if closing {
		state.closing = true
	}
	var callbacks []func()
	for id, r := range forms {
		if r.host == host {
			delete(forms, id)
			callbacks = append(callbacks, r.drop)
		}
	}
	state.count = 0
	state.players = map[PlayerID]int{}
	for state.inflight != 0 {
		formCond.Wait()
	}
	if !closing {
		state.draining = false
	}
	formMu.Unlock()
	for _, callback := range callbacks {
		callback()
	}
}

type skinSnapshot struct {
	host uint64
	skin PlayerSkin
}

type itemSnapshot struct {
	host uint64
	item ItemStack
}

func registerHost(host Host) uint64 {
	if host == nil {
		host = noopHost{}
	}
	id := hostSequence.Add(1)
	hosts.Store(id, host)
	return id
}

func unregisterHost(id uint64) {
	if id != 0 {
		drainHostForms(id, true)
		formMu.Lock()
		delete(formHostState, id)
		formMu.Unlock()
		hosts.Delete(id)
		skinSnapshotMu.Lock()
		for snapshotID, snapshot := range skinSnapshots {
			if snapshot.host == id {
				delete(skinSnapshots, snapshotID)
			}
		}
		delete(skinSnapshotCounts, id)
		skinSnapshotMu.Unlock()
		itemSnapshotMu.Lock()
		for snapshotID, snapshot := range itemSnapshots {
			if snapshot.host == id {
				delete(itemSnapshots, snapshotID)
			}
		}
		delete(itemSnapshotCounts, id)
		itemSnapshotMu.Unlock()
	}
}

func resolveHost(id uint64) (Host, bool) {
	host, ok := hosts.Load(id)
	if !ok {
		return nil, false
	}
	return host.(Host), true
}

func registerSkinSnapshot(host uint64, skin PlayerSkin) (uint64, bool) {
	skinSnapshotMu.Lock()
	defer skinSnapshotMu.Unlock()
	if skinSnapshotCounts[host] >= maxSkinSnapshotsPerHost {
		return 0, false
	}
	id := skinSnapshotSequence.Add(1)
	skinSnapshots[id] = skinSnapshot{host: host, skin: clonePlayerSkin(skin)}
	skinSnapshotCounts[host]++
	return id, true
}

func resolveSkinSnapshot(host, id uint64) (PlayerSkin, bool) {
	skinSnapshotMu.Lock()
	defer skinSnapshotMu.Unlock()
	value, ok := skinSnapshots[id]
	if !ok || value.host != host {
		return PlayerSkin{}, false
	}
	return value.skin, true
}

func unregisterSkinSnapshot(host, id uint64) {
	skinSnapshotMu.Lock()
	defer skinSnapshotMu.Unlock()
	value, ok := skinSnapshots[id]
	if ok && value.host == host {
		delete(skinSnapshots, id)
		skinSnapshotCounts[host]--
		if skinSnapshotCounts[host] == 0 {
			delete(skinSnapshotCounts, host)
		}
	}
}

func clonePlayerSkin(value PlayerSkin) PlayerSkin {
	value.Pixels = append([]byte(nil), value.Pixels...)
	value.Model = append([]byte(nil), value.Model...)
	value.CapePixels = append([]byte(nil), value.CapePixels...)
	value.Animations = append([]SkinAnimation(nil), value.Animations...)
	for index := range value.Animations {
		value.Animations[index].Pixels = append([]byte(nil), value.Animations[index].Pixels...)
	}
	return value
}

func registerItemSnapshot(host uint64, item ItemStack) (uint64, bool) {
	itemSnapshotMu.Lock()
	defer itemSnapshotMu.Unlock()
	if itemSnapshotCounts[host] >= maxItemSnapshotsPerHost {
		return 0, false
	}
	id := itemSnapshotSequence.Add(1)
	itemSnapshots[id] = itemSnapshot{host: host, item: cloneItemStack(item)}
	itemSnapshotCounts[host]++
	return id, true
}

func resolveItemSnapshot(host, id uint64) (ItemStack, bool) {
	itemSnapshotMu.Lock()
	defer itemSnapshotMu.Unlock()
	value, ok := itemSnapshots[id]
	if !ok || value.host != host {
		return ItemStack{}, false
	}
	return value.item, true
}

func unregisterItemSnapshot(host, id uint64) {
	itemSnapshotMu.Lock()
	defer itemSnapshotMu.Unlock()
	value, ok := itemSnapshots[id]
	if ok && value.host == host {
		delete(itemSnapshots, id)
		itemSnapshotCounts[host]--
		if itemSnapshotCounts[host] == 0 {
			delete(itemSnapshotCounts, host)
		}
	}
}

func cloneItemStack(value ItemStack) ItemStack {
	value.Lore = append([]string(nil), value.Lore...)
	value.NBT = append([]byte(nil), value.NBT...)
	value.ValuesNBT = append([]byte(nil), value.ValuesNBT...)
	value.Enchantments = append([]ItemEnchantment(nil), value.Enchantments...)
	return value
}
