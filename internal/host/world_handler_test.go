package host

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"slices"
	"strings"
	"testing"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

func TestWorldHandlerRecordsPlayerDeparture(t *testing.T) {
	players := NewPlayers()
	w := world.Config{Synchronous: true}.New()
	w.Handle(NewWorldHandler(players.EntityRegistry(), players, native.WorldID(44)))
	t.Cleanup(func() { _ = w.Close() })
	id := uuid.MustParse("65c2514b-f068-4da5-a403-1d5890c8f2a7")
	handle := world.EntitySpawnOpts{ID: id, Position: mgl64.Vec3{1, 2, 3}}.New(
		player.Type,
		player.Config{UUID: id, Name: "Traveller", Position: mgl64.Vec3{1, 2, 3}},
	)
	if err := w.Do(func(tx *world.Tx) {
		connected := tx.AddEntity(handle).(*player.Player)
		players.Register(connected, 7)
		if tx.RemoveEntity(connected) == nil {
			t.Fatal("remove player returned no handle")
		}
		departure, ok := players.takeWorldDeparture(connected)
		if !ok || departure != 44 {
			t.Fatalf("departure = %d, %v; want 44, true", departure, ok)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

type worldRuntimeStub struct {
	subscriptions       uint64
	players             *Players
	calls               map[string]int
	flow                native.WorldLiquidFlowInput
	decay               native.WorldLiquidDecayInput
	harden              native.WorldLiquidHardenInput
	sound               native.WorldSoundInput
	redstone            native.WorldRedstoneUpdateInput
	spawnEntity         native.EntityID
	spawnResolved       bool
	spawnPlayerResolved bool
	despawnResolved     bool
	explosionResolved   bool
	explosionOutput     native.WorldExplosionOutput
	explosionErr        error
	flowErr             error
}

type callbackWorldSound struct{ callback native.WorldSoundCallback }

func (s callbackWorldSound) Play(*world.World, mgl64.Vec3)            {}
func (s callbackWorldSound) SoundCallback() native.WorldSoundCallback { return s.callback }

func (r *worldRuntimeStub) Subscriptions() uint64 { return r.subscriptions }
func (r *worldRuntimeStub) HandleWorldScheduled(
	uint64, uint64, native.InvocationID, native.WorldTaskPhase, native.WorldTaskResult,
) error {
	return nil
}
func (r *worldRuntimeStub) called(name string, invocation native.InvocationID) {
	r.calls[name]++
	if _, ok := r.players.InvocationTx(invocation); !ok {
		panic("world callback has no active transaction")
	}
}
func (r *worldRuntimeStub) HandleWorldLiquidFlow(i native.InvocationID, input native.WorldLiquidFlowInput, _ bool) (bool, error) {
	r.called("flow", i)
	r.flow = input
	return true, r.flowErr
}
func (r *worldRuntimeStub) HandleWorldLiquidDecay(i native.InvocationID, input native.WorldLiquidDecayInput, _ bool) (bool, error) {
	r.called("decay", i)
	r.decay = input
	return true, nil
}
func (r *worldRuntimeStub) HandleWorldLiquidHarden(i native.InvocationID, input native.WorldLiquidHardenInput, _ bool) (bool, error) {
	r.called("harden", i)
	r.harden = input
	return true, nil
}
func (r *worldRuntimeStub) HandleWorldSound(i native.InvocationID, input native.WorldSoundInput, _ bool) (bool, error) {
	r.called("sound", i)
	r.sound = input
	return true, nil
}
func (r *worldRuntimeStub) HandleWorldFireSpread(i native.InvocationID, _ native.WorldFireSpreadInput, _ bool) (bool, error) {
	r.called("fire-spread", i)
	return true, nil
}
func (r *worldRuntimeStub) HandleWorldBlockBurn(i native.InvocationID, _ native.WorldPositionInput, _ bool) (bool, error) {
	r.called("block-burn", i)
	return true, nil
}
func (r *worldRuntimeStub) HandleWorldCropTrample(i native.InvocationID, _ native.WorldPositionInput, _ bool) (bool, error) {
	r.called("crop-trample", i)
	return true, nil
}
func (r *worldRuntimeStub) HandleWorldLeavesDecay(i native.InvocationID, _ native.WorldPositionInput, _ bool) (bool, error) {
	r.called("leaves-decay", i)
	return true, nil
}
func (r *worldRuntimeStub) HandleWorldEntitySpawn(i native.InvocationID, input native.WorldEntityInput) error {
	r.called("entity-spawn", i)
	r.spawnEntity = input.Entity
	_, r.spawnResolved = r.players.ResolveEntityID(input.Entity, i)
	_, r.spawnPlayerResolved = r.players.EntityPlayer(i, input.Entity)
	return nil
}
func (r *worldRuntimeStub) HandleWorldEntityDespawn(i native.InvocationID, input native.WorldEntityInput) error {
	r.called("entity-despawn", i)
	_, r.despawnResolved = r.players.ResolveEntityID(input.Entity, i)
	return nil
}

func TestWorldHandlerRegistersPlayerBeforeSpawnCallback(t *testing.T) {
	players := NewPlayers()
	runtime := &worldRuntimeStub{
		players: players, calls: map[string]int{}, subscriptions: native.WorldEntitySpawnSubscription,
	}
	handler := NewWorldHandler(players.EntityRegistry(), players, 1)
	handler.AttachRuntime(runtime, nil)
	w := world.Config{Synchronous: true}.New()
	w.Handle(handler)
	var connected *player.Player
	if err := w.Do(func(tx *world.Tx) {
		connected = addWorldHandlerTestPlayer(tx, "Early", "ed77bd87-bbe1-4b3e-8065-00491194fa28")
		if !runtime.spawnPlayerResolved {
			t.Fatal("spawn callback did not materialize player")
		}
		if _, ok := players.ID(connected); ok {
			t.Fatal("temporary player registration leaked past spawn callback")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	entityID, ok := players.EntityRegistry().ID(connected)
	if !ok || entityID != runtime.spawnEntity {
		t.Fatalf("generic entity ID changed after spawn: event=%#v current=%#v", runtime.spawnEntity, entityID)
	}
	later := players.Register(connected, runtime.spawnEntity.Generation+100)
	want := native.PlayerID{UUID: runtime.spawnEntity.UUID, Generation: runtime.spawnEntity.Generation}
	if later != want {
		t.Fatalf("later registration changed player ID: event=%#v later=%#v", runtime.spawnEntity, later)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(players.Names(), "Early") {
		t.Fatal("world close removed accepted player before quit cleanup")
	}
	players.Unregister(connected)
	if _, ok := players.ID(connected); ok {
		t.Fatal("quit cleanup left accepted player registered")
	}
}

func TestWorldHandlerExplosionWithoutEntityRegistryIsIgnored(t *testing.T) {
	players := NewPlayers()
	runtime := &worldRuntimeStub{players: players, calls: map[string]int{}, subscriptions: native.WorldExplosionSubscription}
	handler := NewWorldHandler(nil, players, 1)
	handler.AttachRuntime(runtime, nil)
	w := world.Config{Synchronous: true}.New()
	t.Cleanup(func() { _ = w.Close() })
	if err := w.Do(func(tx *world.Tx) {
		connected := addWorldHandlerTestPlayer(tx, "NoRegistry", "20d82244-0d0e-42e4-86a2-444833e582b3")
		entities := []world.Entity{connected}
		blocks := []cube.Pos{{1, 2, 3}}
		dropChance, spawnFire := .5, false
		handler.HandleExplosion(tx.Event(), mgl64.Vec3{}, &entities, &blocks, &dropChance, &spawnFire)
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if runtime.calls["explosion"] != 0 {
		t.Fatal("explosion callback ran without entity registry")
	}
}
func (r *worldRuntimeStub) HandleWorldExplosion(i native.InvocationID, input native.WorldExplosionInput, _ float64, _ bool, _ bool) (native.WorldExplosionOutput, error) {
	r.called("explosion", i)
	if len(input.Entities) != 0 {
		_, r.explosionResolved = r.players.ResolveEntityID(input.Entities[0], i)
	}
	return r.explosionOutput, r.explosionErr
}
func (r *worldRuntimeStub) HandleWorldRedstoneUpdate(i native.InvocationID, input native.WorldRedstoneUpdateInput, _ bool) (bool, error) {
	r.called("redstone", i)
	r.redstone = input
	return true, nil
}
func (r *worldRuntimeStub) HandleWorldClose(i native.InvocationID) error {
	r.called("close", i)
	return nil
}

func TestWorldHandlerForwardsAllCallbacks(t *testing.T) {
	players := NewPlayers()
	runtime := &worldRuntimeStub{players: players, calls: map[string]int{}, subscriptions: native.WorldLiquidFlowSubscription | native.WorldLiquidDecaySubscription | native.WorldLiquidHardenSubscription |
		native.WorldSoundSubscription | native.WorldFireSpreadSubscription | native.WorldBlockBurnSubscription |
		native.WorldCropTrampleSubscription | native.WorldLeavesDecaySubscription | native.WorldEntitySpawnSubscription |
		native.WorldEntityDespawnSubscription | native.WorldExplosionSubscription | native.WorldRedstoneUpdateSubscription |
		native.WorldCloseSubscription}
	handler := NewWorldHandler(players.EntityRegistry(), players, 91)
	handler.AttachRuntime(runtime, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	w := world.Config{Synchronous: true}.New()
	t.Cleanup(func() { _ = w.Close() })
	if err := w.Do(func(tx *world.Tx) {
		first := addWorldHandlerTestPlayer(tx, "First", "c36d1736-c9c6-4d03-beec-71535682521d")
		second := addWorldHandlerTestPlayer(tx, "Second", "a8080d2a-f10b-4494-87ca-6f20da6173ea")
		players.Register(first, 1)
		players.Register(second, 2)
		handler.HandleEntitySpawn(tx, first)
		handler.HandleEntitySpawn(tx, second)

		water := block.Water{Depth: 8, Still: true}
		lava := block.Lava{Depth: 8, Still: true}
		stone := block.Stone{}
		assertWorldCancelled(t, tx, func(ctx *world.Context) {
			handler.HandleLiquidFlow(ctx, cube.Pos{1, 2, 3}, cube.Pos{4, 5, 6}, water, stone)
		})
		assertWorldCancelled(t, tx, func(ctx *world.Context) { handler.HandleLiquidDecay(ctx, cube.Pos{7, 8, 9}, water, nil) })
		assertWorldCancelled(t, tx, func(ctx *world.Context) { handler.HandleLiquidHarden(ctx, cube.Pos{10, 11, 12}, water, lava, stone) })
		assertWorldCancelled(t, tx, func(ctx *world.Context) {
			handler.HandleSound(ctx, sound.Note{Instrument: sound.Pling(), Pitch: 30}, mgl64.Vec3{1.5, 2.5, 3.5})
		})
		assertWorldCancelled(t, tx, func(ctx *world.Context) { handler.HandleFireSpread(ctx, cube.Pos{1, 1, 1}, cube.Pos{2, 2, 2}) })
		assertWorldCancelled(t, tx, func(ctx *world.Context) { handler.HandleBlockBurn(ctx, cube.Pos{3, 3, 3}) })
		assertWorldCancelled(t, tx, func(ctx *world.Context) { handler.HandleCropTrample(ctx, cube.Pos{4, 4, 4}) })
		assertWorldCancelled(t, tx, func(ctx *world.Context) { handler.HandleLeavesDecay(ctx, cube.Pos{5, 5, 5}) })

		firstID, _ := players.EntityRegistry().ID(first)
		secondID, _ := players.EntityRegistry().ID(second)
		runtime.explosionOutput = native.WorldExplosionOutput{
			Cancelled: true, Entities: []native.EntityID{secondID, firstID},
			Blocks: []native.BlockPos{{X: 9, Y: 8, Z: 7}}, ItemDropChance: .25, SpawnFire: true,
		}
		entities := []world.Entity{first, second}
		blocks := []cube.Pos{{1, 2, 3}, {4, 5, 6}}
		dropChance, spawnFire := .75, false
		ctx := tx.Event()
		handler.HandleExplosion(ctx, mgl64.Vec3{6, 5, 4}, &entities, &blocks, &dropChance, &spawnFire)
		if !ctx.Cancelled() || len(entities) != 2 || entities[0].H() != second.H() || entities[1].H() != first.H() ||
			len(blocks) != 1 || blocks[0] != (cube.Pos{9, 8, 7}) || dropChance != .25 || !spawnFire {
			t.Fatalf("explosion mutation entities=%v blocks=%v drop=%v fire=%v cancelled=%v", entities, blocks, dropChance, spawnFire, ctx.Cancelled())
		}

		redstone := world.RedstoneUpdate{Pos: cube.Pos{2, 3, 4}, ChangedNeighbour: cube.Pos{3, 3, 4}, HasChangedNeighbour: true,
			ChangedRedstoneRelevant: true, Source: cube.Pos{1, 3, 4}, HasSource: true, Before: stone, After: block.RedstoneWire{},
			OldPower: 4, NewPower: 12, CurrentTick: 99, Cause: world.RedstoneUpdateCauseScheduledTick}
		assertWorldCancelled(t, tx, func(ctx *world.Context) { handler.HandleRedstoneUpdate(ctx, redstone) })
		handler.HandleEntityDespawn(tx, first)
		handler.HandleClose(tx)
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, event := range []string{"flow", "decay", "harden", "sound", "fire-spread", "block-burn", "crop-trample", "leaves-decay", "explosion", "redstone", "entity-despawn", "close"} {
		if runtime.calls[event] != 1 {
			t.Fatalf("%s calls = %d", event, runtime.calls[event])
		}
	}
	if runtime.calls["entity-spawn"] != 2 || !runtime.spawnResolved || !runtime.despawnResolved || !runtime.explosionResolved {
		t.Fatalf("entity callbacks spawn=%d resolved=%v/%v/%v", runtime.calls["entity-spawn"], runtime.spawnResolved, runtime.despawnResolved, runtime.explosionResolved)
	}
	if runtime.flow.Liquid.Identifier != "minecraft:water" || runtime.decay.After != nil || runtime.harden.NewBlock.Identifier != "minecraft:stone" ||
		runtime.sound.Sound.Kind != native.SoundNote || runtime.sound.Sound.Data != 15 || runtime.sound.Sound.Integer != 30 ||
		runtime.redstone.Cause != native.RedstoneUpdateCauseScheduledTick || runtime.redstone.NewPower != 12 {
		t.Fatalf("converted callbacks flow=%#v decay=%#v harden=%#v sound=%#v redstone=%#v", runtime.flow, runtime.decay, runtime.harden, runtime.sound, runtime.redstone)
	}
}

func TestWorldHandlerForwardsCustomSound(t *testing.T) {
	players := NewPlayers()
	runtime := &worldRuntimeStub{
		players: players, calls: map[string]int{}, subscriptions: native.WorldSoundSubscription,
	}
	handler := NewWorldHandler(players.EntityRegistry(), players, 1)
	handler.AttachRuntime(runtime, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	w := world.Config{Synchronous: true}.New()
	t.Cleanup(func() { _ = w.Close() })
	if err := w.Do(func(tx *world.Tx) {
		ctx := tx.Event()
		handler.HandleSound(ctx, callbackWorldSound{callback: native.WorldSoundCallback{
			Function: 41, Context: 73,
		}}, mgl64.Vec3{1, 2, 3})
		if !ctx.Cancelled() {
			t.Fatal("custom sound was not cancellable")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if runtime.calls["sound"] != 1 || runtime.sound.Sound.Callback == nil ||
		runtime.sound.Sound.Callback.Function != 41 || runtime.sound.Sound.Callback.Context != 73 {
		t.Fatalf("custom sound callback = %#v", runtime.sound)
	}
}

func TestWorldHandlerInvalidExplosionMutationPreservesDefaults(t *testing.T) {
	players := NewPlayers()
	var logs bytes.Buffer
	runtime := &worldRuntimeStub{players: players, calls: map[string]int{}, subscriptions: native.WorldExplosionSubscription,
		explosionOutput: native.WorldExplosionOutput{Cancelled: true, Entities: []native.EntityID{{Generation: 999}},
			Blocks: []native.BlockPos{{X: 9}}, ItemDropChance: .1, SpawnFire: true}}
	handler := NewWorldHandler(players.EntityRegistry(), players, 1)
	handler.AttachRuntime(runtime, slog.New(slog.NewTextHandler(&logs, nil)))
	w := world.Config{Synchronous: true}.New()
	t.Cleanup(func() { _ = w.Close() })
	if err := w.Do(func(tx *world.Tx) {
		connected := addWorldHandlerTestPlayer(tx, "Default", "d831f992-4cff-4417-bf33-bc92cb9a7901")
		players.Register(connected, 1)
		entities := []world.Entity{connected}
		blocks := []cube.Pos{{1, 2, 3}}
		dropChance, spawnFire := .8, false
		ctx := tx.Event()
		handler.HandleExplosion(ctx, mgl64.Vec3{}, &entities, &blocks, &dropChance, &spawnFire)
		if ctx.Cancelled() || len(entities) != 1 || entities[0] != connected || len(blocks) != 1 || blocks[0] != (cube.Pos{1, 2, 3}) || dropChance != .8 || spawnFire {
			t.Fatalf("invalid mutation changed defaults")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(logs.String(), "invalid entity") {
		t.Fatalf("missing invalid mutation log: %s", logs.String())
	}
}

func TestWorldHandlerKeepsCancellationOnRuntimeError(t *testing.T) {
	players := NewPlayers()
	runtime := &worldRuntimeStub{
		players: players, calls: map[string]int{},
		subscriptions: native.WorldLiquidFlowSubscription | native.WorldExplosionSubscription,
		flowErr:       errors.New("flow failed"),
		explosionErr:  errors.New("explosion failed"),
		explosionOutput: native.WorldExplosionOutput{
			Cancelled: true, Blocks: []native.BlockPos{{X: 9, Y: 8, Z: 7}},
			ItemDropChance: .1, SpawnFire: true,
		},
	}
	handler := NewWorldHandler(players.EntityRegistry(), players, 1)
	handler.AttachRuntime(runtime, slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)))
	w := world.Config{Synchronous: true}.New()
	t.Cleanup(func() { _ = w.Close() })
	if err := w.Do(func(tx *world.Tx) {
		flow := tx.Event()
		handler.HandleLiquidFlow(flow, cube.Pos{}, cube.Pos{1, 2, 3}, block.Water{Depth: 8, Still: true}, block.Stone{})
		if !flow.Cancelled() {
			t.Fatal("runtime error cleared liquid-flow cancellation")
		}

		entities := []world.Entity{}
		blocks := []cube.Pos{{1, 2, 3}}
		dropChance, spawnFire := .8, false
		explosion := tx.Event()
		handler.HandleExplosion(explosion, mgl64.Vec3{}, &entities, &blocks, &dropChance, &spawnFire)
		if !explosion.Cancelled() {
			t.Fatal("runtime error cleared explosion cancellation")
		}
		if len(entities) != 0 || len(blocks) != 1 || blocks[0] != (cube.Pos{1, 2, 3}) || dropChance != .8 || spawnFire {
			t.Fatalf("runtime error applied explosion mutations: entities=%v blocks=%v drop=%v fire=%v", entities, blocks, dropChance, spawnFire)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func addWorldHandlerTestPlayer(tx *world.Tx, name, rawID string) *player.Player {
	id := uuid.MustParse(rawID)
	handle := world.EntitySpawnOpts{ID: id, Position: mgl64.Vec3{1, 2, 3}}.New(player.Type, player.Config{UUID: id, Name: name, Position: mgl64.Vec3{1, 2, 3}})
	return tx.AddEntity(handle).(*player.Player)
}

func assertWorldCancelled(t *testing.T, tx *world.Tx, call func(*world.Context)) {
	t.Helper()
	ctx := tx.Event()
	call(ctx)
	if !ctx.Cancelled() {
		t.Fatal("world callback was not cancelled")
	}
}
