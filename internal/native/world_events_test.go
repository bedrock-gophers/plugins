package native

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func openWorldEventRuntime(t *testing.T) *Runtime {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("world-event C fixture does not support Windows")
	}
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	extension := ".so"
	flags := []string{"-shared", "-fPIC", "-std=c11"}
	if runtime.GOOS == "darwin" {
		extension = ".dylib"
		flags[0] = "-dynamiclib"
	}
	runtimeLibrary := filepath.Join(root, "build", "lib", "libdragonfly_plugin_runtime"+extension)
	if _, err := os.Stat(runtimeLibrary); err != nil {
		t.Skipf("C# runtime not built: run make build-native (%v)", err)
	}
	pluginDirectory := t.TempDir()
	plugin := filepath.Join(pluginDirectory, "world-event-test"+extension)
	arguments := append(flags,
		"-I", filepath.Join(root, "abi", "include"),
		filepath.Join(root, "internal", "native", "testdata", "world_event_plugin.c"),
		"-o", plugin,
	)
	command := exec.Command("cc", arguments...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("compile world-event fixture: %v\n%s", err, output)
	}
	pluginRuntime, err := Open(runtimeLibrary, pluginDirectory)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pluginRuntime.Close)
	if err := pluginRuntime.Enable(); err != nil {
		t.Fatal(err)
	}
	return pluginRuntime
}

func TestWorldEventIDsAndSubscriptionsAreContiguous(t *testing.T) {
	events := []uint32{
		WorldLiquidFlowEvent, WorldLiquidDecayEvent, WorldLiquidHardenEvent, WorldSoundEvent,
		WorldFireSpreadEvent, WorldBlockBurnEvent, WorldCropTrampleEvent, WorldLeavesDecayEvent,
		WorldEntitySpawnEvent, WorldEntityDespawnEvent, WorldExplosionEvent, WorldRedstoneUpdateEvent,
		WorldCloseEvent,
	}
	subscriptions := []uint64{
		WorldLiquidFlowSubscription, WorldLiquidDecaySubscription, WorldLiquidHardenSubscription, WorldSoundSubscription,
		WorldFireSpreadSubscription, WorldBlockBurnSubscription, WorldCropTrampleSubscription, WorldLeavesDecaySubscription,
		WorldEntitySpawnSubscription, WorldEntityDespawnSubscription, WorldExplosionSubscription, WorldRedstoneUpdateSubscription,
		WorldCloseSubscription,
	}
	for index, event := range events {
		want := uint32(42 + index)
		if event != want {
			t.Fatalf("event %d = %d, want %d", index, event, want)
		}
		if subscriptions[index] != uint64(1)<<(event-1) {
			t.Fatalf("subscription %d = %d, want bit %d", index, subscriptions[index], event-1)
		}
	}
}

func TestWorldEventRuntimeRoundTrip(t *testing.T) {
	runtime := openWorldEventRuntime(t)
	if got := runtime.Subscriptions(); got != uint64(0x003ffe0000000000) {
		t.Fatalf("subscriptions = %#x", got)
	}
	water := WorldBlock{Identifier: "minecraft:water"}
	lava := WorldBlock{Identifier: "minecraft:lava"}
	stone := WorldBlock{Identifier: "minecraft:stone"}
	air := WorldBlock{Identifier: "minecraft:air"}
	from, to := BlockPos{X: 1, Y: 2, Z: 3}, BlockPos{X: 4, Y: 5, Z: 6}
	cancellables := []struct {
		name string
		call func() (bool, error)
	}{
		{name: "liquid flow", call: func() (bool, error) {
			return runtime.HandleWorldLiquidFlow(73, WorldLiquidFlowInput{From: from, Into: to, Liquid: water, Replaced: air}, false)
		}},
		{name: "liquid decay", call: func() (bool, error) {
			return runtime.HandleWorldLiquidDecay(73, WorldLiquidDecayInput{Before: water}, false)
		}},
		{name: "liquid harden", call: func() (bool, error) {
			return runtime.HandleWorldLiquidHarden(73, WorldLiquidHardenInput{LiquidHardened: water, OtherLiquid: lava, NewBlock: stone}, false)
		}},
		{name: "sound", call: func() (bool, error) {
			return runtime.HandleWorldSound(73, WorldSoundInput{Position: Vec3{Y: 64}}, false)
		}},
		{name: "fire spread", call: func() (bool, error) {
			return runtime.HandleWorldFireSpread(73, WorldFireSpreadInput{From: from, To: to}, false)
		}},
		{name: "block burn", call: func() (bool, error) {
			return runtime.HandleWorldBlockBurn(73, WorldPositionInput{Position: BlockPos{X: 7, Y: 8, Z: 9}}, false)
		}},
		{name: "crop trample", call: func() (bool, error) {
			return runtime.HandleWorldCropTrample(73, WorldPositionInput{Position: BlockPos{X: 7, Y: 8, Z: 9}}, false)
		}},
		{name: "leaves decay", call: func() (bool, error) {
			return runtime.HandleWorldLeavesDecay(73, WorldPositionInput{Position: BlockPos{X: 7, Y: 8, Z: 9}}, false)
		}},
	}
	for _, test := range cancellables {
		t.Run(test.name, func(t *testing.T) {
			cancelled, err := test.call()
			if err != nil || !cancelled {
				t.Fatalf("cancelled=%v error=%v", cancelled, err)
			}
		})
	}
	customCancelled, err := runtime.HandleWorldSound(73, WorldSoundInput{
		Sound:    WorldSound{Callback: &WorldSoundCallback{Function: 41, Context: 73}},
		Position: Vec3{X: 9, Y: 64},
	}, false)
	if err != nil || !customCancelled {
		t.Fatalf("custom sound cancelled=%v error=%v", customCancelled, err)
	}

	entity := WorldEntityInput{Entity: EntityID{UUID: [16]byte{3}, Generation: 4}}
	if err := runtime.HandleWorldEntitySpawn(73, entity); err != nil {
		t.Fatal(err)
	}
	if err := runtime.HandleWorldEntityDespawn(73, entity); err != nil {
		t.Fatal(err)
	}
	explosion, err := runtime.HandleWorldExplosion(73, WorldExplosionInput{
		Position: Vec3{X: 1}, Entities: []EntityID{entity.Entity}, Blocks: []BlockPos{from, to},
	}, 0.25, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if !explosion.Cancelled || !explosion.SpawnFire || explosion.ItemDropChance != 0.75 ||
		!reflect.DeepEqual(explosion.Entities, []EntityID{{UUID: [16]byte{9}, Generation: 99}}) ||
		!reflect.DeepEqual(explosion.Blocks, []BlockPos{{X: -9, Y: 90, Z: 9}}) {
		t.Fatalf("explosion = %#v", explosion)
	}
	after := WorldBlock{Identifier: "minecraft:redstone_wire", PropertiesNBT: []byte{1}}
	redstoneCancelled, err := runtime.HandleWorldRedstoneUpdate(73, WorldRedstoneUpdateInput{
		Position: from, ChangedNeighbour: to, HasChangedNeighbour: true, ChangedRedstoneRelevant: true,
		Source: BlockPos{X: 7, Y: 8, Z: 9}, HasSource: true,
		Before: WorldBlock{Identifier: "minecraft:redstone_wire"}, After: &after,
		OldPower: 2, NewPower: 13, CurrentTick: 1234, Cause: RedstoneUpdateCauseCompilerRebuild,
	}, false)
	if err != nil || !redstoneCancelled {
		t.Fatalf("redstone cancelled=%v error=%v", redstoneCancelled, err)
	}
	if err := runtime.HandleWorldClose(73); err != nil {
		t.Fatal(err)
	}
}

func TestCancellableWorldCallbacksKeepIncomingCancellationOnRuntimeError(t *testing.T) {
	runtime := &Runtime{}
	block := WorldBlock{Identifier: "minecraft:stone"}
	tests := []struct {
		name string
		call func() (bool, error)
	}{
		{name: "liquid flow", call: func() (bool, error) {
			return runtime.HandleWorldLiquidFlow(1, WorldLiquidFlowInput{Liquid: block, Replaced: block}, true)
		}},
		{name: "liquid decay", call: func() (bool, error) {
			return runtime.HandleWorldLiquidDecay(1, WorldLiquidDecayInput{Before: block}, true)
		}},
		{name: "liquid harden", call: func() (bool, error) {
			return runtime.HandleWorldLiquidHarden(1, WorldLiquidHardenInput{LiquidHardened: block, OtherLiquid: block, NewBlock: block}, true)
		}},
		{name: "sound", call: func() (bool, error) {
			return runtime.HandleWorldSound(1, WorldSoundInput{}, true)
		}},
		{name: "fire spread", call: func() (bool, error) {
			return runtime.HandleWorldFireSpread(1, WorldFireSpreadInput{}, true)
		}},
		{name: "block burn", call: func() (bool, error) {
			return runtime.HandleWorldBlockBurn(1, WorldPositionInput{}, true)
		}},
		{name: "crop trample", call: func() (bool, error) {
			return runtime.HandleWorldCropTrample(1, WorldPositionInput{}, true)
		}},
		{name: "leaves decay", call: func() (bool, error) {
			return runtime.HandleWorldLeavesDecay(1, WorldPositionInput{}, true)
		}},
		{name: "redstone update", call: func() (bool, error) {
			return runtime.HandleWorldRedstoneUpdate(1, WorldRedstoneUpdateInput{Before: block}, true)
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cancelled, err := test.call()
			if err == nil {
				t.Fatal("closed runtime call returned no error")
			}
			if !cancelled {
				t.Fatal("incoming cancellation was cleared")
			}
		})
	}
}

func TestClosedRuntimeWorldExplosionPreservesIndependentInput(t *testing.T) {
	entities := []EntityID{{UUID: [16]byte{1}, Generation: 2}}
	blocks := []BlockPos{{X: 3, Y: 4, Z: 5}}
	output, err := (&Runtime{}).HandleWorldExplosion(7, WorldExplosionInput{
		Position: Vec3{X: 1, Y: 2, Z: 3}, Entities: entities, Blocks: blocks,
	}, 0.25, true, true)
	if err == nil {
		t.Fatal("closed runtime call returned no error")
	}
	if !output.Cancelled || output.ItemDropChance != 0.25 || !output.SpawnFire ||
		!reflect.DeepEqual(output.Entities, entities) || !reflect.DeepEqual(output.Blocks, blocks) {
		t.Fatalf("output = %#v", output)
	}
	entities[0].Generation = 99
	blocks[0].X = 99
	if output.Entities[0].Generation != 2 || output.Blocks[0].X != 3 {
		t.Fatal("output aliases input slices")
	}
}

func TestWorldEventNativeSliceRoundTrip(t *testing.T) {
	entities := []EntityID{
		{UUID: [16]byte{1, 2, 3}, Generation: 4},
		{UUID: [16]byte{5, 6, 7}, Generation: 8},
	}
	nativeEntities, releaseEntities, ok := nativeEntityIDs(entities)
	if !ok {
		t.Fatal("encode entities")
	}
	defer releaseEntities()
	decodedEntities, ok := copyNativeEntityIDs(nativeEntities, uint64(len(entities)))
	if !ok || !reflect.DeepEqual(decodedEntities, entities) {
		t.Fatalf("entities = %#v, %v", decodedEntities, ok)
	}

	blocks := []BlockPos{{X: -1, Y: 64, Z: 2}, {X: 3, Y: -64, Z: 5}}
	nativeBlocks, releaseBlocks, ok := nativeBlockPositions(blocks)
	if !ok {
		t.Fatal("encode blocks")
	}
	defer releaseBlocks()
	decodedBlocks, ok := copyNativeBlockPositions(nativeBlocks, uint64(len(blocks)))
	if !ok || !reflect.DeepEqual(decodedBlocks, blocks) {
		t.Fatalf("blocks = %#v, %v", decodedBlocks, ok)
	}
}

func TestWorldEventNativeSlicesRejectInvalidValues(t *testing.T) {
	if _, release, ok := nativeEntityIDs([]EntityID{{}}); ok {
		release()
		t.Fatal("encoded zero-generation entity")
	}
	if _, ok := copyNativeEntityIDs(nil, 1); ok {
		t.Fatal("decoded null entity slice")
	}
	if _, ok := copyNativeBlockPositions(nil, 1); ok {
		t.Fatal("decoded null block slice")
	}
}

func TestWorldNotificationCallbacksRejectClosedRuntime(t *testing.T) {
	runtime := &Runtime{}
	entity := WorldEntityInput{Entity: EntityID{Generation: 1}}
	for name, call := range map[string]func() error{
		"spawn":   func() error { return runtime.HandleWorldEntitySpawn(1, entity) },
		"despawn": func() error { return runtime.HandleWorldEntityDespawn(1, entity) },
		"close":   func() error { return runtime.HandleWorldClose(1) },
	} {
		t.Run(name, func(t *testing.T) {
			if err := call(); err == nil {
				t.Fatal("closed runtime call returned no error")
			}
		})
	}
}
