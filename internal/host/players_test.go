package host

import (
	"testing"

	"github.com/df-mc/dragonfly/server/player"
)

func TestPlayersTracksStableGenerationAndNames(t *testing.T) {
	withPlayer(t, func(player *player.Player) {
		players := NewPlayers()
		id := players.Register(player, 42)
		if id.Generation != 42 || len(players.Names()) != 1 || players.Names()[0] != "TestPlayer" {
			t.Fatalf("id=%+v names=%v", id, players.Names())
		}
		resolved, ok := players.ResolveName("testplayer")
		if !ok || resolved != id {
			t.Fatalf("resolved=%+v ok=%v", resolved, ok)
		}
		connected, ok := players.ResolveID(id)
		if !ok || connected != player {
			t.Fatalf("connected=%p ok=%v", connected, ok)
		}
		players.Unregister(player)
		if len(players.Names()) != 0 {
			t.Fatalf("names after unregister = %v", players.Names())
		}
	})
}
