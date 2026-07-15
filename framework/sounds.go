package framework

import (
	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type customWorldSound struct {
	callback native.WorldSoundCallback
	world    native.WorldID
	called   bool
	ok       bool
}

func (s *customWorldSound) Play(_ *world.World, position mgl64.Vec3) {
	s.called = true
	s.ok = native.CallWorldSound(s.callback, s.world, native.Vec3{X: position.X(), Y: position.Y(), Z: position.Z()})
}

func (m *WorldManager) PlayWorldSound(invocation native.InvocationID, id native.WorldID, position native.Vec3, value native.WorldSound) bool {
	entry, ok := m.entryForInvocation(invocation, id)
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

// PlayCustomWorldSound invokes a C# implementation of world.Sound through the
// exact synchronous transaction that received Tx.PlaySound.
func (m *WorldManager) PlayCustomWorldSound(invocation native.InvocationID, id native.WorldID, position native.Vec3, callback native.WorldSoundCallback) bool {
	entry, ok := m.entryForInvocation(invocation, id)
	if !ok || !finiteVec3(position) || callback.Function == 0 || callback.Context == 0 {
		return false
	}
	entry.lifecycle.RLock()
	closed := entry.closed
	entry.lifecycle.RUnlock()
	if closed {
		return false
	}
	tx := m.currentTx(invocation, entry.world)
	if tx == nil {
		return false
	}
	sound := &customWorldSound{callback: callback, world: entry.id, ok: true}
	tx.PlaySound(vec3(position), sound)
	return !sound.called || sound.ok
}
