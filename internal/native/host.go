package native

import (
	"sync"
	"sync/atomic"
	"time"
)

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
)

const maxSkinSnapshotsPerHost = 32
const maxItemSnapshotsPerHost = 64

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
