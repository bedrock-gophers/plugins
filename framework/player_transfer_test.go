package framework

import (
	"context"
	"errors"
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

type cancellingTeleportHandler struct {
	player.NopHandler
	calls int
}

func (h *cancellingTeleportHandler) HandleTeleport(ctx *player.Context, _ mgl64.Vec3) {
	h.calls++
	ctx.Cancel()
}

type blockingPlayerSpawnHandler struct {
	*host.WorldHandler
	entered chan struct{}
	release chan struct{}
}

type blockingPlayerDespawnHandler struct {
	*host.WorldHandler
	entered chan struct{}
	release chan struct{}
}

func (h *blockingPlayerDespawnHandler) HandleEntityDespawn(tx *world.Tx, entity world.Entity) {
	h.WorldHandler.HandleEntityDespawn(tx, entity)
	if _, ok := entity.(*player.Player); ok {
		close(h.entered)
		<-h.release
	}
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
	handler := new(cancellingTeleportHandler)
	if err := fixture.source.Do(func(tx *world.Tx) {
		invocation, end := fixture.players.BeginInvocation(tx)
		defer end()
		entity, _ := fixture.handle.Entity(tx)
		entity.(*player.Player).Handle(handler)
		if !fixture.manager.TransferPlayer(invocation, fixture.playerID, fixture.sourceID, want) {
			t.Fatal("same-world transfer was rejected")
		}
		connected, ok := fixture.handle.Entity(tx)
		if !ok || connected.(*player.Player).Position() != (mgl64.Vec3{1, 65, 2}) {
			t.Fatalf("same-world position = %v, %v", connected, ok)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if handler.calls != 1 {
		t.Fatalf("teleport handler calls = %d", handler.calls)
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

func TestTransferPlayerOffCallbackMutationWaitsWhileHandleWorldless(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	source := world.Config{Synchronous: true}.New()
	if err := manager.RegisterCore(OverworldID, source); err != nil {
		t.Fatal(err)
	}
	target, err := manager.Create("example:async-target", world.Config{})
	if err != nil {
		t.Fatal(err)
	}
	entered := make(chan struct{})
	release := make(chan struct{})
	barrier := target.Do(func(*world.Tx) {
		close(entered)
		<-release
	})
	t.Cleanup(func() {
		select {
		case <-release:
		default:
			close(release)
		}
		_ = manager.CloseCustom()
		_ = source.Close()
	})
	targetID, _ := manager.WorldByName(0, "example:async-target")
	playerUUID := uuid.MustParse("2c0f8607-c1ad-4c19-9ab4-3d8a6f5fa994")
	handle := world.EntitySpawnOpts{ID: playerUUID}.New(player.Type, player.Config{UUID: playerUUID, Name: "Ordered"})
	var id native.PlayerID
	if err := source.Do(func(tx *world.Tx) {
		id = players.Register(tx.AddEntity(handle).(*player.Player), 93)
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("destination barrier did not start")
	}
	if !manager.TransferPlayer(0, id, targetID, native.Vec3{Y: 70}) {
		t.Fatal("off-callback transfer was rejected")
	}
	if !players.TransformPlayer(0, id, native.PlayerTransformVelocity, native.Vec3{Y: 0.75}, 0, 0) {
		t.Fatal("worldless mutation was rejected")
	}
	close(release)
	if err := barrier.Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for {
		velocity, err := world.Call(context.Background(), target, func(tx *world.Tx) (mgl64.Vec3, error) {
			entity, ok := handle.Entity(tx)
			if !ok {
				return mgl64.Vec3{}, nil
			}
			return entity.(*player.Player).Velocity(), nil
		})
		if err == nil && velocity == (mgl64.Vec3{0, 0.75, 0}) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("queued mutation velocity = %v, error = %v", velocity, err)
		}
		time.Sleep(time.Millisecond)
	}
}

func TestTransferPlayerOffCallbackWaitsForAsyncSourceSubmission(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	source := world.Config{}.New()
	if err := manager.RegisterCore(OverworldID, source); err != nil {
		t.Fatal(err)
	}
	despawn := &blockingPlayerDespawnHandler{
		WorldHandler: source.Handler().(*host.WorldHandler),
		entered:      make(chan struct{}),
		release:      make(chan struct{}),
	}
	source.Handle(despawn)
	target, err := manager.Create("example:async-source-target", world.Config{Synchronous: true})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		select {
		case <-despawn.release:
		default:
			close(despawn.release)
		}
		_ = manager.CloseCustom()
		_ = source.Close()
	})
	targetID, _ := manager.WorldByName(0, "example:async-source-target")
	playerUUID := uuid.MustParse("0beb95eb-c78d-44fe-a5c5-69685873d2e0")
	handle := world.EntitySpawnOpts{ID: playerUUID}.New(player.Type, player.Config{UUID: playerUUID, Name: "AsyncSource"})
	var id native.PlayerID
	if err := source.Do(func(tx *world.Tx) {
		id = players.Register(tx.AddEntity(handle).(*player.Player), 94)
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	result := make(chan bool, 1)
	go func() {
		result <- manager.TransferPlayer(0, id, targetID, native.Vec3{Y: 70})
	}()
	select {
	case <-despawn.entered:
	case <-time.After(time.Second):
		t.Fatal("source despawn did not reach barrier")
	}
	select {
	case accepted := <-result:
		t.Fatalf("transfer returned before destination submission: %v", accepted)
	case <-time.After(50 * time.Millisecond):
	}
	close(despawn.release)
	select {
	case accepted := <-result:
		if !accepted {
			t.Fatal("transfer was rejected after source submission")
		}
	case <-time.After(time.Second):
		t.Fatal("transfer did not return after source submission")
	}
	if err := target.Do(func(tx *world.Tx) {
		if _, ok := handle.Entity(tx); !ok {
			t.Fatal("player did not reach destination")
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

func TestTransferPlayerNopSessionCloseBeforeDeferredHandoffFailsClosed(t *testing.T) {
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

func TestTransferPlayerRejectsBusyDestinationWithoutBlockingOwner(t *testing.T) {
	fixture := newTransferFixture(t)
	destination, ok := fixture.manager.entryByHandle(fixture.targetID)
	if !ok {
		t.Fatal("destination missing")
	}
	destination.lifecycle.Lock()
	result := make(chan bool, 1)
	go func() {
		_ = fixture.source.Do(func(tx *world.Tx) {
			invocation, end := fixture.players.BeginInvocation(tx)
			defer end()
			result <- fixture.manager.TransferPlayer(invocation, fixture.playerID, fixture.targetID, native.Vec3{Y: 70})
		}).Wait(context.Background())
	}()
	select {
	case accepted := <-result:
		destination.lifecycle.Unlock()
		if accepted {
			t.Fatal("transfer acquired a destination with an active lifecycle writer")
		}
	case <-time.After(50 * time.Millisecond):
		destination.lifecycle.Unlock()
		<-result
		t.Fatal("transfer blocked the source owner on destination lifecycle")
	}
}

func TestTransferPlayerOffCallbackReportsRejectedBusyDestination(t *testing.T) {
	fixture := newTransferFixture(t)
	destination, ok := fixture.manager.entryByHandle(fixture.targetID)
	if !ok {
		t.Fatal("destination missing")
	}
	destination.lifecycle.Lock()
	defer destination.lifecycle.Unlock()

	if fixture.manager.TransferPlayer(0, fixture.playerID, fixture.targetID, native.Vec3{Y: 70}) {
		t.Fatal("off-callback transfer reported acceptance after its destination precheck became invalid")
	}
	if err := fixture.source.Do(func(tx *world.Tx) {
		if _, ok := fixture.handle.Entity(tx); !ok {
			t.Fatal("rejected transfer removed player from source")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
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

func TestTransferPlayerEvictsHandleWhenBothWorldsClose(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	source := world.Config{}.New()
	if err := manager.RegisterCore(OverworldID, source); err != nil {
		t.Fatal(err)
	}
	target, err := manager.Create("example:closing-target", world.Config{})
	if err != nil {
		t.Fatal(err)
	}
	entered := make(chan struct{})
	release := make(chan struct{})
	barrier := target.Do(func(*world.Tx) {
		close(entered)
		<-release
	})
	t.Cleanup(func() {
		select {
		case <-release:
		default:
			close(release)
		}
		_ = manager.CloseCustom()
		_ = source.Close()
	})
	targetID, _ := manager.WorldByName(0, "example:closing-target")
	playerUUID := uuid.MustParse("015606a4-46a4-493f-9eab-952126fe8800")
	handle := world.EntitySpawnOpts{ID: playerUUID}.New(player.Type, player.Config{UUID: playerUUID, Name: "Orphan"})
	var id native.PlayerID
	if err := source.Do(func(tx *world.Tx) {
		id = players.Register(tx.AddEntity(handle).(*player.Player), 95)
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("destination barrier did not start")
	}
	if !manager.TransferPlayer(0, id, targetID, native.Vec3{Y: 70}) {
		t.Fatal("transfer was rejected")
	}
	targetClosed := make(chan error, 1)
	go func() { targetClosed <- target.Close() }()
	if err := source.Close(); err != nil {
		t.Fatal(err)
	}
	close(release)
	if err := barrier.Wait(context.Background()); err != nil && !errors.Is(err, world.ErrWorldClosed) {
		t.Fatal(err)
	}
	select {
	case err := <-targetClosed:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("destination did not close")
	}
	deadline := time.Now().Add(time.Second)
	for {
		_, registered := players.Handle(id)
		if !registered && handle.Closed() {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("terminal handle registered=%v closed=%v", registered, handle.Closed())
		}
		time.Sleep(time.Millisecond)
	}
}
