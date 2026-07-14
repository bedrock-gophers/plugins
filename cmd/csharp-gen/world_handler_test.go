package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const worldHandlerSource = `package world
type Handler interface {
	HandleLiquidFlow(ctx *Context, from, into cube.Pos, liquid Liquid, replaced Block)
	HandleLiquidDecay(ctx *Context, pos cube.Pos, before, after Liquid)
	HandleLiquidHarden(ctx *Context, hardenedPos cube.Pos, liquidHardened, otherLiquid, newBlock Block)
	HandleSound(ctx *Context, s Sound, pos mgl64.Vec3)
	HandleFireSpread(ctx *Context, from, to cube.Pos)
	HandleBlockBurn(ctx *Context, pos cube.Pos)
	HandleCropTrample(ctx *Context, pos cube.Pos)
	HandleLeavesDecay(ctx *Context, pos cube.Pos)
	HandleEntitySpawn(tx *Tx, e Entity)
	HandleEntityDespawn(tx *Tx, e Entity)
	HandleExplosion(ctx *Context, position mgl64.Vec3, entities *[]Entity, blocks *[]cube.Pos, itemDropChance *float64, spawnFire *bool)
	HandleRedstoneUpdate(ctx *Context, update RedstoneUpdate)
	HandleClose(tx *Tx)
}`

const redstoneUpdateSource = `package world
type RedstoneUpdateCause uint8
const (
	RedstoneUpdateCauseBlockUpdate RedstoneUpdateCause = iota
	RedstoneUpdateCauseScheduledTick
	RedstoneUpdateCauseCompilerRebuild
)
type RedstoneUpdate struct {
	Pos cube.Pos
	ChangedNeighbour cube.Pos
	HasChangedNeighbour bool
	ChangedRedstoneRelevant bool
	Source cube.Pos
	HasSource bool
	Before Block
	After Block
	OldPower int
	NewPower int
	CurrentTick int64
	Cause RedstoneUpdateCause
}`

func TestWorldHandlerUsesGoAST(t *testing.T) {
	path := filepath.Join(t.TempDir(), "handler.go")
	if err := os.WriteFile(path, []byte(worldHandlerSource), 0o600); err != nil {
		t.Fatal(err)
	}
	methods, err := inspectWorldHandler(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(methods) != len(selectedWorldHandlerMethods) {
		t.Fatalf("generated %d methods, want %d", len(methods), len(selectedWorldHandlerMethods))
	}
	for index, method := range methods {
		if method.Name != selectedWorldHandlerMethods[index] {
			t.Fatalf("method %d = %s, want %s", index, method.Name, selectedWorldHandlerMethods[index])
		}
	}
	redstonePath := filepath.Join(t.TempDir(), "redstone.go")
	if err := os.WriteFile(redstonePath, []byte(redstoneUpdateSource), 0o600); err != nil {
		t.Fatal(err)
	}
	redstone, err := inspectRedstoneUpdate(redstonePath)
	if err != nil {
		t.Fatal(err)
	}
	generated := string(generateWorldHandler(methods, redstone))
	for _, expected := range []string{
		"public interface Handler",
		"public abstract partial class Plugin : World.Handler",
		"BlockUpdate = 0",
		"ScheduledTick = 1",
		"CompilerRebuild = 2",
		"Cube.Pos ChangedNeighbour",
		"World.Block? After",
		"long CurrentTick",
		"World.RedstoneUpdateCause Cause",
		"void HandleLiquidFlow(World.Context ctx, Cube.Pos from, Cube.Pos into, World.Liquid liquid, World.Block replaced);",
		"void HandleLiquidDecay(World.Context ctx, Cube.Pos pos, World.Liquid before, World.Liquid? after);",
		"void HandleLiquidHarden(World.Context ctx, Cube.Pos hardenedPos, World.Block liquidHardened, World.Block otherLiquid, World.Block newBlock);",
		"void HandleSound(World.Context ctx, World.Sound s, Vector3 pos);",
		"void HandleFireSpread(World.Context ctx, Cube.Pos from, Cube.Pos to);",
		"void HandleBlockBurn(World.Context ctx, Cube.Pos pos);",
		"void HandleCropTrample(World.Context ctx, Cube.Pos pos);",
		"void HandleLeavesDecay(World.Context ctx, Cube.Pos pos);",
		"void HandleEntitySpawn(World.Tx tx, World.Entity e);",
		"void HandleEntityDespawn(World.Tx tx, World.Entity e);",
		"void HandleExplosion(World.Context ctx, Vector3 position, ref World.Entity[] entities, ref Cube.Pos[] blocks, ref double itemDropChance, ref bool spawnFire);",
		"void HandleRedstoneUpdate(World.Context ctx, World.RedstoneUpdate update);",
		"void HandleClose(World.Tx tx);",
		"[HandlerSubscription(2199023255552UL)]",
		"[HandlerSubscription(9007199254740992UL)]",
		"public virtual void HandleClose(World.Tx tx) { }",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated output missing %q:\n%s", expected, generated)
		}
	}
}

func TestPinnedDragonflyWorldHandlerHasExactSurface(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	module, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	methods, err := inspectWorldHandler(filepath.Join(
		string(bytes.TrimSpace(module)), "server", "world", "handler.go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(methods) != len(selectedWorldHandlerMethods) {
		t.Fatalf("generated %d pinned Dragonfly handlers, want %d", len(methods), len(selectedWorldHandlerMethods))
	}
	for index, method := range methods {
		if method.Name != selectedWorldHandlerMethods[index] {
			t.Fatalf("pinned method %d = %s, want %s", index, method.Name, selectedWorldHandlerMethods[index])
		}
		wantSubscription := uint64(1) << uint(firstWorldHandlerSubscriptionBit+index)
		if method.Subscription != wantSubscription {
			t.Fatalf("%s subscription = %d, want %d", method.Name, method.Subscription, wantSubscription)
		}
	}
	redstone, err := inspectRedstoneUpdate(filepath.Join(
		string(bytes.TrimSpace(module)), "server", "world", "redstone.go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(redstone.Causes) != 3 || len(redstone.Fields) != 12 {
		t.Fatalf("pinned redstone surface has %d causes and %d fields, want 3 and 12", len(redstone.Causes), len(redstone.Fields))
	}
}

func TestWorldHandlerRejectsUnknownMethod(t *testing.T) {
	path := filepath.Join(t.TempDir(), "handler.go")
	if err := os.WriteFile(path, []byte(`package world
type Handler interface { HandleFuture(ctx *Context) }
`), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := inspectWorldHandler(path)
	if err == nil || !strings.Contains(err.Error(), "unsupported world.Handler.HandleFuture method") {
		t.Fatalf("inspectWorldHandler() error = %v, want unsupported-method error", err)
	}
}
