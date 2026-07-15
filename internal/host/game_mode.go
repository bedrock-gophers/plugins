package host

import "github.com/df-mc/dragonfly/server/world"

const (
	gameModeBuiltinFlag uint64 = 1 << 63
	gameModeCustomMask  uint64 = 1<<8 - 1
)

// descriptorGameMode is the host representation of an unregistered Dragonfly
// game mode. Each bit maps to one GameMode capability in interface order.
type descriptorGameMode uint8

func (m descriptorGameMode) AllowsEditing() bool       { return m&(1<<0) != 0 }
func (m descriptorGameMode) AllowsTakingDamage() bool  { return m&(1<<1) != 0 }
func (m descriptorGameMode) CreativeInventory() bool   { return m&(1<<2) != 0 }
func (m descriptorGameMode) HasCollision() bool        { return m&(1<<3) != 0 }
func (m descriptorGameMode) AllowsFlying() bool        { return m&(1<<4) != 0 }
func (m descriptorGameMode) AllowsInteraction() bool   { return m&(1<<5) != 0 }
func (m descriptorGameMode) Visible() bool             { return m&(1<<6) != 0 }
func (m descriptorGameMode) InstantPortalTravel() bool { return m&(1<<7) != 0 }

func decodeGameModeDescriptor(value int64) (world.GameMode, bool) {
	descriptor := uint64(value)
	if descriptor&gameModeBuiltinFlag != 0 {
		id := descriptor &^ gameModeBuiltinFlag
		if id > uint64(^uint(0)>>1) {
			return nil, false
		}
		mode, ok := world.GameModeByID(int(id))
		if !ok {
			return nil, false
		}
		return mode, true
	}
	if descriptor&^gameModeCustomMask != 0 {
		return nil, false
	}
	return descriptorGameMode(uint8(descriptor)), true
}

// DecodeGameModeDescriptor decodes the private native representation used by
// both player and world game-mode operations.
func DecodeGameModeDescriptor(value int64) (world.GameMode, bool) {
	return decodeGameModeDescriptor(value)
}

func encodeGameModeDescriptor(mode world.GameMode) (value int64, ok bool) {
	if mode == nil {
		return 0, false
	}
	if id, registered := registeredGameModeID(mode); registered {
		if id < 0 || uint64(id) >= gameModeBuiltinFlag {
			return 0, false
		}
		return int64(gameModeBuiltinFlag | uint64(id)), true
	}

	defer func() {
		if recover() != nil {
			value, ok = 0, false
		}
	}()
	var descriptor uint64
	capabilities := [...]bool{
		mode.AllowsEditing(),
		mode.AllowsTakingDamage(),
		mode.CreativeInventory(),
		mode.HasCollision(),
		mode.AllowsFlying(),
		mode.AllowsInteraction(),
		mode.Visible(),
		mode.InstantPortalTravel(),
	}
	for bit, enabled := range capabilities {
		if enabled {
			descriptor |= 1 << bit
		}
	}
	return int64(descriptor), true
}

// EncodeGameModeDescriptor encodes a registered or structural Dragonfly game mode.
func EncodeGameModeDescriptor(mode world.GameMode) (int64, bool) {
	return encodeGameModeDescriptor(mode)
}

func registeredGameModeID(mode world.GameMode) (id int, ok bool) {
	defer func() {
		if recover() != nil {
			id, ok = 0, false
		}
	}()
	return world.GameModeID(mode)
}
