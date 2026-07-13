package host

import (
	"testing"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/player"
)

func TestPlayersHandleRequiresExactLiveGeneration(t *testing.T) {
	withPlayer(t, func(connected *player.Player) {
		players := NewPlayers()
		id := players.Register(connected, 47)
		if handle, ok := players.Handle(id); !ok || handle != connected.H() {
			t.Fatalf("Handle(%#v) = %p, %v", id, handle, ok)
		}
		stale := native.PlayerID{UUID: id.UUID, Generation: id.Generation + 1}
		if _, ok := players.Handle(stale); ok {
			t.Fatal("Handle accepted a stale generation")
		}
		players.Unregister(connected)
		if _, ok := players.Handle(id); ok {
			t.Fatal("Handle accepted an unregistered player")
		}
	})
}

func TestPlayersUnregistersWorldlessHandle(t *testing.T) {
	withPlayer(t, func(connected *player.Player) {
		players := NewPlayers()
		id := players.Register(connected, 48)
		players.UnregisterHandle(connected.H())
		if _, ok := players.Handle(id); ok {
			t.Fatal("worldless handle remained registered")
		}
	})
}
