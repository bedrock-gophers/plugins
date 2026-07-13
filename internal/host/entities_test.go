package host

import (
	"context"
	"testing"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

func TestEntitiesKeepStableGenerationAndResolveExactTransaction(t *testing.T) {
	other := world.Config{Synchronous: true}.New()
	t.Cleanup(func() { _ = other.Close() })
	withPlayerTx(t, func(tx *world.Tx, connected *player.Player) {
		entities := NewEntities()
		first := entities.Register(connected)
		if first.Generation == 0 {
			t.Fatal("zero entity generation")
		}
		if same := entities.Register(connected); same != first {
			t.Fatalf("same handle changed identity: %+v to %+v", first, same)
		}
		resolved, ok := entities.Resolve(first, tx)
		if !ok || resolved.H() != connected.H() {
			t.Fatalf("entity did not resolve in owner transaction: %T %v", resolved, ok)
		}
		if err := other.Do(func(otherTx *world.Tx) {
			if _, ok := entities.Resolve(first, otherTx); ok {
				t.Fatal("entity resolved in a different world's transaction")
			}
		}).Wait(context.Background()); err != nil {
			t.Fatal(err)
		}

		entities.deactivateHandle(connected.H())
		if _, ok := entities.Handle(first); ok {
			t.Fatal("worldless entity handle remained active")
		}
		if reactivated := entities.Register(connected); reactivated != first {
			t.Fatalf("world transfer changed identity: %+v to %+v", first, reactivated)
		}

		entities.Unregister(connected)
		if _, ok := entities.Resolve(first, tx); ok {
			t.Fatal("expired entity ID resolved")
		}
		second := entities.Register(connected)
		if second.UUID != first.UUID || second.Generation <= first.Generation {
			t.Fatalf("re-registration did not advance generation: %+v to %+v", first, second)
		}
		if _, ok := entities.Resolve(first, tx); ok {
			t.Fatal("old generation aliased re-registered entity")
		}
		if _, ok := entities.Resolve(second, tx); !ok {
			t.Fatal("new generation did not resolve")
		}
	})
}

func TestPlayersShareIdentityWithEntityRegistry(t *testing.T) {
	withPlayerTx(t, func(tx *world.Tx, connected *player.Player) {
		players := NewPlayers()
		playerID := players.Register(connected, 41)
		entityID, ok := players.EntityRegistry().ID(connected)
		if !ok || entityID.UUID != playerID.UUID || entityID.Generation != playerID.Generation {
			t.Fatalf("player ID %+v and entity ID %+v differ (ok=%v)", playerID, entityID, ok)
		}
		invocation, end := players.BeginInvocation(tx)
		defer end()
		resolved, ok := players.ResolveEntityID(entityID, invocation)
		if !ok || resolved.H() != connected.H() {
			t.Fatalf("shared entity ID did not resolve: %T %v", resolved, ok)
		}
		players.Unregister(connected)
		if _, ok := players.EntityRegistry().Handle(entityID); ok {
			t.Fatal("unregistered player remained in entity registry")
		}
	})
}

func TestEntitiesReserveFreshGeneration(t *testing.T) {
	withPlayerTx(t, func(_ *world.Tx, connected *player.Player) {
		entities := NewEntities()
		handle := connected.H()
		first := entities.registerHandle(handle, 0)
		entities.expire(first)
		reserved := entities.reserve(handle)
		if first == reserved || reserved.Generation == 0 {
			t.Fatalf("reserved ID = %#v after %#v", reserved, first)
		}
		if _, ok := entities.Handle(reserved); ok {
			t.Fatal("reserved ID resolved before activation")
		}
		entities.activate(handle)
		if _, ok := entities.Handle(reserved); !ok {
			t.Fatal("activated ID did not resolve")
		}
	})
}

func TestEntitiesReserveActiveHandleFreshGeneration(t *testing.T) {
	withPlayerTx(t, func(_ *world.Tx, connected *player.Player) {
		entities := NewEntities()
		handle := connected.H()
		first := entities.registerHandle(handle, 0)
		reserved := entities.reserve(handle)
		if reserved == first || reserved.Generation <= first.Generation {
			t.Fatalf("active ID was reused: %#v then %#v", first, reserved)
		}
		if _, ok := entities.Handle(first); ok {
			t.Fatal("prior active ID remained resolvable after reservation")
		}
		if _, ok := entities.Handle(reserved); ok {
			t.Fatal("reserved ID resolved before activation")
		}
		entities.activate(handle)
		if _, ok := entities.Handle(first); ok {
			t.Fatal("prior active ID revived after activation")
		}
		if resolved, ok := entities.Handle(reserved); !ok || resolved != handle {
			t.Fatalf("activated ID resolved handle = %p, ok=%v", resolved, ok)
		}
	})
}
