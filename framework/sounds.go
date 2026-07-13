package framework

import (
	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/world"
)

func (m *WorldManager) PlayWorldSound(invocation native.InvocationID, id native.WorldID, position native.Vec3, value native.WorldSound) bool {
	entry, ok := m.entryByHandle(id)
	if !ok || !finiteVec3(position) || !host.ValidSound(value) {
		return false
	}
	entry.lifecycle.RLock()
	defer entry.lifecycle.RUnlock()
	if entry.closed {
		return false
	}
	play := func(tx *world.Tx) bool {
		decoded, ok := host.SoundFromNative(tx, value)
		if !ok {
			return false
		}
		tx.PlaySound(vec3(position), decoded)
		return true
	}
	if tx := m.currentTx(invocation, entry.world); tx != nil {
		return play(tx)
	}
	if invocation != 0 {
		if _, ok := m.invocationTx(invocation); !ok {
			return false
		}
	}
	// Cross-world and off-callback writes are fire-and-forget. Success means the
	// validated sound was accepted for scheduling.
	entry.world.Do(func(tx *world.Tx) { play(tx) })
	return true
}
