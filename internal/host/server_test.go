package host

import (
	"context"
	"iter"
	"testing"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

type serverPlayerSourceStub struct {
	handles []*world.EntityHandle
	byUUID  map[uuid.UUID]*world.EntityHandle
}

func (s *serverPlayerSourceStub) Players(tx *world.Tx) iter.Seq[*player.Player] {
	return func(yield func(*player.Player) bool) {
		for _, handle := range s.handles {
			if tx != nil {
				if current, ok := handle.Entity(tx); ok {
					if !yield(current.(*player.Player)) {
						return
					}
					continue
				}
			}
			stop, err := player.Call(context.Background(), handle, func(_ *world.Tx, connected *player.Player) (bool, error) {
				return !yield(connected), nil
			})
			if err == nil && stop {
				return
			}
		}
	}
}

func (s *serverPlayerSourceStub) Player(id uuid.UUID) (*world.EntityHandle, bool) {
	handle, ok := s.byUUID[id]
	return handle, ok
}

func TestServerPlayersRotateBorrowedTransactions(t *testing.T) {
	players, source, closeWorlds := serverPlayersFixture(t)
	defer closeWorlds()
	host := NewServer(players)
	host.source = source

	iterator, ok := host.OpenServerPlayerIterator(0)
	if !ok {
		t.Fatal("open server players")
	}
	firstInvocation, first, found, valid := host.NextServerPlayer(0, iterator)
	if !valid || !found || firstInvocation == 0 || first.Name != "Alpha" {
		t.Fatalf("first = %d, %#v, %v, %v", firstInvocation, first, found, valid)
	}
	if _, ok := players.ResolveID(first.Player, firstInvocation); !ok {
		t.Fatal("first player is not valid inside its iteration body")
	}

	secondInvocation, second, found, valid := host.NextServerPlayer(0, iterator)
	if !valid || !found || secondInvocation == 0 || second.Name != "Bravo" {
		t.Fatalf("second = %d, %#v, %v, %v", secondInvocation, second, found, valid)
	}
	if _, ok := players.InvocationTx(firstInvocation); ok {
		t.Fatal("first player remained valid after iterator advanced")
	}
	if _, ok := players.ResolveID(second.Player, secondInvocation); !ok {
		t.Fatal("second player is not valid inside its iteration body")
	}

	host.CloseServerPlayers(0, iterator)
	if _, ok := players.InvocationTx(secondInvocation); ok {
		t.Fatal("current player remained valid after iterator disposal")
	}
	if _, _, _, valid := host.NextServerPlayer(0, iterator); valid {
		t.Fatal("disposed iterator remained valid")
	}
}

func TestServerPlayersCloseWithOuterInvocation(t *testing.T) {
	players, source, closeWorlds := serverPlayersFixture(t)
	defer closeWorlds()
	host := NewServer(players)
	host.source = source

	var outer native.InvocationID
	var iterator native.PlayerIteratorID
	var nested native.InvocationID
	if err := sourceWorld(source.handles[0]).Do(func(tx *world.Tx) {
		outer, end := players.BeginInvocation(tx)
		defer end()
		var ok bool
		iterator, ok = host.OpenServerPlayerIterator(outer)
		if !ok {
			t.Fatal("open transaction-bound server players")
		}
		nested, _, _, ok = host.NextServerPlayer(outer, iterator)
		if !ok || nested == 0 {
			t.Fatal("advance transaction-bound server players")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := players.InvocationTx(nested); ok {
		t.Fatal("nested player survived outer invocation")
	}
	if _, _, _, valid := host.NextServerPlayer(outer, iterator); valid {
		t.Fatal("iterator survived outer invocation")
	}
}

func TestServerPlayerLookupsUseExactDragonflyHandles(t *testing.T) {
	players, source, closeWorlds := serverPlayersFixture(t)
	defer closeWorlds()
	host := NewServer(players)
	host.source = source
	alpha := uuid.UUID(source.handles[0].UUID())

	byUUID, found, valid := host.ServerPlayer([16]byte(alpha))
	if !valid || !found || !byUUID.Valid() {
		t.Fatalf("UUID lookup = %#v, %v, %v", byUUID, found, valid)
	}
	byName, found, valid := host.ServerPlayerByName("Alpha")
	if !valid || !found || byName != byUUID {
		t.Fatalf("name lookup = %#v, %v, %v; want %#v", byName, found, valid, byUUID)
	}
	again, found, valid := host.ServerPlayer([16]byte(alpha))
	if !valid || !found || again != byUUID {
		t.Fatalf("repeated UUID lookup = %#v, %v, %v; want %#v", again, found, valid, byUUID)
	}
	if _, found, valid := host.ServerPlayerByName("alpha"); !valid || found {
		t.Fatalf("lowercase lookup = found %v, valid %v", found, valid)
	}
	if _, found, valid := host.ServerPlayer([16]byte(uuid.New())); !valid || found {
		t.Fatalf("missing UUID lookup = found %v, valid %v", found, valid)
	}
}

func serverPlayersFixture(t *testing.T) (*Players, *serverPlayerSourceStub, func()) {
	t.Helper()
	players := NewPlayers()
	worlds := []*world.World{
		world.Config{Synchronous: true}.New(),
		world.Config{Synchronous: true}.New(),
	}
	source := &serverPlayerSourceStub{
		byUUID: map[uuid.UUID]*world.EntityHandle{},
	}
	for index, name := range []string{"Alpha", "Bravo"} {
		id := uuid.New()
		handle := world.EntitySpawnOpts{ID: id, Position: mgl64.Vec3{float64(index), 64, 0}}.New(
			player.Type,
			player.Config{UUID: id, Name: name, Position: mgl64.Vec3{float64(index), 64, 0}},
		)
		if err := worlds[index].Do(func(tx *world.Tx) {
			players.Register(tx.AddEntity(handle).(*player.Player), uint64(index+1))
		}).Wait(context.Background()); err != nil {
			t.Fatal(err)
		}
		source.handles = append(source.handles, handle)
		source.byUUID[id] = handle
	}
	return players, source, func() {
		for _, current := range worlds {
			_ = current.Close()
		}
	}
}

func sourceWorld(handle *world.EntityHandle) *world.World {
	var owner *world.World
	_, _ = player.Call(context.Background(), handle, func(tx *world.Tx, _ *player.Player) (struct{}, error) {
		owner = tx.World()
		return struct{}{}, nil
	})
	return owner
}
