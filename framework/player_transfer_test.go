package framework

import (
	"context"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

type transferFixture struct {
	manager        *WorldManager
	players        *host.Players
	source, target *world.World
	sourceID       native.WorldID
	targetID       native.WorldID
	handle         *world.EntityHandle
	playerID       native.PlayerID
}

type transferChangeRecorder struct {
	player.NopHandler
	calls         int
	before, after *world.World
}

type blockingPlayerSpawnHandler struct {
	*host.WorldHandler
	entered chan struct{}
	release chan struct{}
}

func (h *blockingPlayerSpawnHandler) HandleEntitySpawn(tx *world.Tx, entity world.Entity) {
	h.WorldHandler.HandleEntitySpawn(tx, entity)
	if _, ok := entity.(*player.Player); ok {
		close(h.entered)
		<-h.release
	}
}

func (r *transferChangeRecorder) HandleChangeWorld(_ *player.Player, before, after *world.World) {
	r.calls++
	r.before, r.after = before, after
}

func newTransferFixture(t *testing.T) *transferFixture {
	t.Helper()
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	source := world.Config{Synchronous: true}.New()
	if err := manager.RegisterCore(OverworldID, source); err != nil {
		t.Fatal(err)
	}
	target, err := manager.Create("example:target", world.Config{Synchronous: true})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = manager.CloseCustom()
		_ = source.Close()
	})
	playerUUID := uuid.MustParse("939d8cca-26be-4e86-a238-f60d95ee45f2")
	handle := world.EntitySpawnOpts{ID: playerUUID, Position: mgl64.Vec3{1, 65, 2}}.New(
		player.Type,
		player.Config{UUID: playerUUID, Name: "Traveller", Position: mgl64.Vec3{1, 65, 2}},
	)
	fixture := &transferFixture{manager: manager, players: players, source: source, target: target, handle: handle}
	fixture.sourceID, _ = manager.WorldByName(0, string(OverworldID))
	fixture.targetID, _ = manager.WorldByName(0, "example:target")
	if err := source.Do(func(tx *world.Tx) {
		connected := tx.AddEntity(handle).(*player.Player)
		fixture.playerID = players.Register(connected, 91)
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	return fixture
}

func TestTransferPlayerSameWorldUsesOrdinaryTeleport(t *testing.T) {
	fixture := newTransferFixture(t)
	want := native.Vec3{X: 8, Y: 72, Z: -3}
	if err := fixture.source.Do(func(tx *world.Tx) {
		invocation, end := fixture.players.BeginInvocation(tx)
		defer end()
		if !fixture.manager.TransferPlayer(invocation, fixture.playerID, fixture.sourceID, want) {
			t.Fatal("same-world transfer was rejected")
		}
		connected, ok := fixture.handle.Entity(tx)
		if !ok || connected.(*player.Player).Position() != (mgl64.Vec3{8, 72, -3}) {
			t.Fatalf("same-world position = %v, %v", connected, ok)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestTransferPlayerCrossWorldPreservesIdentityAndPosition(t *testing.T) {
	fixture := newTransferFixture(t)
	want := native.Vec3{X: 19.5, Y: 81, Z: -7.25}
	if err := fixture.source.Do(func(tx *world.Tx) {
		invocation, end := fixture.players.BeginInvocation(tx)
		defer end()
		if !fixture.manager.TransferPlayer(invocation, fixture.playerID, fixture.targetID, want) {
			t.Fatal("cross-world transfer was rejected")
		}
		// The transfer must not invalidate the callback's source transaction.
		connected, ok := fixture.handle.Entity(tx)
		if !ok {
			t.Fatal("player left source transaction before callback completed")
		}
		connected.(*player.Player).SetVelocity(mgl64.Vec3{0, 0.25, 0})
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if handle, ok := fixture.players.Handle(fixture.playerID); !ok || handle != fixture.handle {
		t.Fatalf("registered handle = %p, %v; want %p", handle, ok, fixture.handle)
	}
	if err := fixture.target.Do(func(tx *world.Tx) {
		entity, ok := fixture.handle.Entity(tx)
		if !ok {
			t.Fatal("player was not attached to destination")
		}
		connected := entity.(*player.Player)
		if connected.Position() != (mgl64.Vec3{19.5, 81, -7.25}) {
			t.Fatalf("destination position = %v", connected.Position())
		}
		id, ok := fixture.players.ID(connected)
		if !ok || id != fixture.playerID {
			t.Fatalf("destination player ID = %#v, %v", id, ok)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestTransferPlayerRejectsInvalidInputsWithoutMutation(t *testing.T) {
	fixture := newTransferFixture(t)
	if err := fixture.source.Do(func(tx *world.Tx) {
		invocation, end := fixture.players.BeginInvocation(tx)
		defer end()
		stale := fixture.playerID
		stale.Generation++
		for name, accepted := range map[string]bool{
			"stale player": fixture.manager.TransferPlayer(invocation, stale, fixture.targetID, native.Vec3{}),
			"stale world":  fixture.manager.TransferPlayer(invocation, fixture.playerID, fixture.targetID+100, native.Vec3{}),
			"non-finite":   fixture.manager.TransferPlayer(invocation, fixture.playerID, fixture.targetID, native.Vec3{X: math.NaN()}),
		} {
			if accepted {
				t.Fatalf("%s transfer was accepted", name)
			}
		}
		if entity, ok := fixture.handle.Entity(tx); !ok || entity.(*player.Player).Position() != (mgl64.Vec3{1, 65, 2}) {
			t.Fatalf("invalid calls mutated player: %v, %v", entity, ok)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestTransferPlayerOffCallbackFollowsHandleAndOrdersMutations(t *testing.T) {
	fixture := newTransferFixture(t)
	if !fixture.manager.TransferPlayer(0, fixture.playerID, fixture.targetID, native.Vec3{X: 4, Y: 70, Z: 9}) {
		t.Fatal("off-callback transfer was rejected")
	}
	if !fixture.players.TransformPlayer(0, fixture.playerID, native.PlayerTransformVelocity, native.Vec3{Y: 0.5}, 0, 0) {
		t.Fatal("mutation following transfer was rejected")
	}
	if err := fixture.target.Do(func(tx *world.Tx) {
		entity, ok := fixture.handle.Entity(tx)
		if !ok {
			t.Fatal("player was not attached to destination")
		}
		connected := entity.(*player.Player)
		if connected.Position() != (mgl64.Vec3{4, 70, 9}) || connected.Velocity() != (mgl64.Vec3{0, 0.5, 0}) {
			t.Fatalf("destination state position=%v velocity=%v", connected.Position(), connected.Velocity())
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestTransferPlayerEmitsOneNaturalChangeWorldEvent(t *testing.T) {
	fixture := newTransferFixture(t)
	recorder := new(transferChangeRecorder)
	if err := fixture.source.Do(func(tx *world.Tx) {
		entity, ok := fixture.handle.Entity(tx)
		if !ok {
			t.Fatal("source player missing")
		}
		connected := entity.(*player.Player)
		connected.Handle(recorder)
		connected.Tick(tx, 1) // Establish Dragonfly's previous world naturally.
		invocation, end := fixture.players.BeginInvocation(tx)
		defer end()
		if !fixture.manager.TransferPlayer(invocation, fixture.playerID, fixture.targetID, native.Vec3{Y: 70}) {
			t.Fatal("transfer was rejected")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if recorder.calls != 0 {
		t.Fatalf("synthetic change-world calls = %d", recorder.calls)
	}
	if err := fixture.target.Do(func(tx *world.Tx) {
		entity, ok := fixture.handle.Entity(tx)
		if !ok {
			t.Fatal("destination player missing")
		}
		connected := entity.(*player.Player)
		connected.Tick(tx, 2)
		connected.Tick(tx, 3)
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if recorder.calls != 1 || recorder.before != fixture.source || recorder.after != fixture.target {
		t.Fatalf("change-world calls=%d before=%p after=%p", recorder.calls, recorder.before, recorder.after)
	}
}

func TestTransferPlayerRestoresExactHandleWhenDestinationClosed(t *testing.T) {
	fixture := newTransferFixture(t)
	if err := fixture.target.Close(); err != nil {
		t.Fatal(err)
	}
	if err := fixture.source.Do(func(tx *world.Tx) {
		invocation, end := fixture.players.BeginInvocation(tx)
		defer end()
		if !fixture.manager.TransferPlayer(invocation, fixture.playerID, fixture.targetID, native.Vec3{Y: 80}) {
			t.Fatal("transfer to concurrently closed destination was not accepted")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := fixture.source.Do(func(tx *world.Tx) {
		entity, ok := fixture.handle.Entity(tx)
		if !ok || entity.H() != fixture.handle {
			t.Fatalf("restored entity = %v, %v", entity, ok)
		}
		if entity.(*player.Player).Position() != (mgl64.Vec3{1, 65, 2}) {
			t.Fatalf("restored source position = %v", entity.(*player.Player).Position())
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestTransferPlayerQuitBeforeDeferredHandoffFailsClosed(t *testing.T) {
	fixture := newTransferFixture(t)
	if err := fixture.source.Do(func(tx *world.Tx) {
		invocation, end := fixture.players.BeginInvocation(tx)
		defer end()
		if !fixture.manager.TransferPlayer(invocation, fixture.playerID, fixture.targetID, native.Vec3{Y: 80}) {
			t.Fatal("transfer was rejected")
		}
		entity, _ := fixture.handle.Entity(tx)
		_ = entity.(*player.Player).Close()
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := fixture.players.Handle(fixture.playerID); ok || !fixture.handle.Closed() {
		t.Fatal("quit player remained transferable")
	}
}

func TestTransferPlayerLeaseMakesUnloadWaitForHandoff(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	source := world.Config{Synchronous: true}.New()
	if err := manager.RegisterCore(OverworldID, source); err != nil {
		t.Fatal(err)
	}
	target, err := manager.Create("example:leased", world.Config{})
	if err != nil {
		t.Fatal(err)
	}
	blocking := &blockingPlayerSpawnHandler{
		WorldHandler: target.Handler().(*host.WorldHandler),
		entered:      make(chan struct{}),
		release:      make(chan struct{}),
	}
	target.Handle(blocking)
	t.Cleanup(func() {
		_ = manager.CloseCustom()
		_ = source.Close()
	})
	targetID, _ := manager.WorldByName(0, "example:leased")
	playerUUID := uuid.MustParse("2f5c0a74-5da9-47e1-a488-6ca3835ec4b8")
	handle := world.EntitySpawnOpts{ID: playerUUID}.New(player.Type, player.Config{UUID: playerUUID, Name: "Lease"})
	if err := source.Do(func(tx *world.Tx) {
		connected := tx.AddEntity(handle).(*player.Player)
		id := players.Register(connected, 92)
		invocation, end := players.BeginInvocation(tx)
		defer end()
		if !manager.TransferPlayer(invocation, id, targetID, native.Vec3{Y: 70}) {
			t.Fatal("transfer was rejected")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	select {
	case <-blocking.entered:
	case <-time.After(time.Second):
		t.Fatal("destination handoff did not start")
	}
	unloaded := make(chan error, 1)
	go func() { unloaded <- manager.Unload("example:leased") }()
	select {
	case err := <-unloaded:
		t.Fatalf("unload completed during handoff: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	close(blocking.release)
	select {
	case err := <-unloaded:
		if err == nil || !strings.Contains(err.Error(), "contains 1 player") {
			t.Fatalf("unload after handoff = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("unload did not resume after handoff")
	}
}
