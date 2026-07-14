package host

import (
	"log/slog"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// WorldRuntime is the native event surface consumed by WorldHandler.
type WorldRuntime interface {
	Subscriptions() uint64
	HandleWorldScheduled(uint64, uint64, native.InvocationID, native.WorldTaskPhase, native.WorldTaskResult) error
	HandleWorldLiquidFlow(native.InvocationID, native.WorldLiquidFlowInput, bool) (bool, error)
	HandleWorldLiquidDecay(native.InvocationID, native.WorldLiquidDecayInput, bool) (bool, error)
	HandleWorldLiquidHarden(native.InvocationID, native.WorldLiquidHardenInput, bool) (bool, error)
	HandleWorldSound(native.InvocationID, native.WorldSoundInput, bool) (bool, error)
	HandleWorldFireSpread(native.InvocationID, native.WorldFireSpreadInput, bool) (bool, error)
	HandleWorldBlockBurn(native.InvocationID, native.WorldPositionInput, bool) (bool, error)
	HandleWorldCropTrample(native.InvocationID, native.WorldPositionInput, bool) (bool, error)
	HandleWorldLeavesDecay(native.InvocationID, native.WorldPositionInput, bool) (bool, error)
	HandleWorldEntitySpawn(native.InvocationID, native.WorldEntityInput) error
	HandleWorldEntityDespawn(native.InvocationID, native.WorldEntityInput) error
	HandleWorldExplosion(native.InvocationID, native.WorldExplosionInput, float64, bool, bool) (native.WorldExplosionOutput, error)
	HandleWorldRedstoneUpdate(native.InvocationID, native.WorldRedstoneUpdateInput, bool) (bool, error)
	HandleWorldClose(native.InvocationID) error
}

// WorldHandler is installed on every framework-owned world before the server starts.
type WorldHandler struct {
	world.NopHandler
	entities *Entities
	players  *Players
	world    native.WorldID
	runtime  WorldRuntime
	log      *slog.Logger
}

var _ world.Handler = (*WorldHandler)(nil)

func NewWorldHandler(entities *Entities, players *Players, worldID native.WorldID) *WorldHandler {
	return &WorldHandler{entities: entities, players: players, world: worldID, log: slog.Default()}
}

// AttachRuntime installs the runtime used for subscribed world events.
func (h *WorldHandler) AttachRuntime(runtime WorldRuntime, log *slog.Logger) {
	h.runtime = runtime
	if log != nil {
		h.log = log
	}
}

func (h *WorldHandler) HandleLiquidFlow(ctx *world.Context, from, into cube.Pos, liquid world.Liquid, replaced world.Block) {
	if !h.subscribed(native.WorldLiquidFlowSubscription) {
		return
	}
	liquidValue, liquidOK := nativeEventBlock(liquid)
	replacedValue, replacedOK := nativeEventBlock(replaced)
	if !liquidOK || !replacedOK {
		h.log.Error("encode native world liquid-flow event")
		return
	}
	h.cancellable(ctx, "liquid-flow", func(invocation native.InvocationID, cancelled bool) (bool, error) {
		return h.runtime.HandleWorldLiquidFlow(invocation, native.WorldLiquidFlowInput{
			From: nativeBlockPos(from), Into: nativeBlockPos(into), Liquid: liquidValue, Replaced: replacedValue,
		}, cancelled)
	})
}

func (h *WorldHandler) HandleLiquidDecay(ctx *world.Context, position cube.Pos, before, after world.Liquid) {
	if !h.subscribed(native.WorldLiquidDecaySubscription) {
		return
	}
	beforeValue, ok := nativeEventBlock(before)
	if !ok {
		h.log.Error("encode native world liquid-decay event")
		return
	}
	var afterValue *native.WorldBlock
	if after != nil {
		encoded, ok := nativeEventBlock(after)
		if !ok {
			h.log.Error("encode native world liquid-decay event")
			return
		}
		afterValue = &encoded
	}
	h.cancellable(ctx, "liquid-decay", func(invocation native.InvocationID, cancelled bool) (bool, error) {
		return h.runtime.HandleWorldLiquidDecay(invocation, native.WorldLiquidDecayInput{
			Position: nativeBlockPos(position), Before: beforeValue, After: afterValue,
		}, cancelled)
	})
}

func (h *WorldHandler) HandleLiquidHarden(ctx *world.Context, position cube.Pos, hardened, other, newBlock world.Block) {
	if !h.subscribed(native.WorldLiquidHardenSubscription) {
		return
	}
	hardenedValue, hardenedOK := nativeEventBlock(hardened)
	otherValue, otherOK := nativeEventBlock(other)
	newValue, newOK := nativeEventBlock(newBlock)
	if !hardenedOK || !otherOK || !newOK {
		h.log.Error("encode native world liquid-harden event")
		return
	}
	h.cancellable(ctx, "liquid-harden", func(invocation native.InvocationID, cancelled bool) (bool, error) {
		return h.runtime.HandleWorldLiquidHarden(invocation, native.WorldLiquidHardenInput{
			Position: nativeBlockPos(position), LiquidHardened: hardenedValue, OtherLiquid: otherValue, NewBlock: newValue,
		}, cancelled)
	})
}

func (h *WorldHandler) HandleSound(ctx *world.Context, sound world.Sound, position mgl64.Vec3) {
	if !h.subscribed(native.WorldSoundSubscription) {
		return
	}
	value, ok := SoundToNative(ctx.Tx, sound)
	if !ok {
		h.log.Error("encode native world sound event", "sound", sound)
		return
	}
	h.cancellable(ctx, "sound", func(invocation native.InvocationID, cancelled bool) (bool, error) {
		return h.runtime.HandleWorldSound(invocation, native.WorldSoundInput{
			Sound: value, Position: native.Vec3{X: position.X(), Y: position.Y(), Z: position.Z()},
		}, cancelled)
	})
}

func (h *WorldHandler) HandleFireSpread(ctx *world.Context, from, to cube.Pos) {
	h.positionPair(ctx, native.WorldFireSpreadSubscription, "fire-spread", func(invocation native.InvocationID, cancelled bool) (bool, error) {
		return h.runtime.HandleWorldFireSpread(invocation, native.WorldFireSpreadInput{From: nativeBlockPos(from), To: nativeBlockPos(to)}, cancelled)
	})
}

func (h *WorldHandler) HandleBlockBurn(ctx *world.Context, position cube.Pos) {
	h.position(ctx, position, native.WorldBlockBurnSubscription, "block-burn", h.runtimeHandleBlockBurn)
}

func (h *WorldHandler) HandleCropTrample(ctx *world.Context, position cube.Pos) {
	h.position(ctx, position, native.WorldCropTrampleSubscription, "crop-trample", h.runtimeHandleCropTrample)
}

func (h *WorldHandler) HandleLeavesDecay(ctx *world.Context, position cube.Pos) {
	h.position(ctx, position, native.WorldLeavesDecaySubscription, "leaves-decay", h.runtimeHandleLeavesDecay)
}

func (h *WorldHandler) HandleEntitySpawn(tx *world.Tx, entity world.Entity) {
	if h.entities == nil || entity == nil {
		return
	}
	id := h.entities.Register(entity)
	if !h.subscribed(native.WorldEntitySpawnSubscription) || id.Generation == 0 {
		return
	}
	cleanup := func() {}
	if connected, ok := entity.(*player.Player); ok {
		cleanup = h.players.registerTemporary(connected, id)
	}
	defer cleanup()
	h.invoke(tx, "entity-spawn", func(invocation native.InvocationID) error {
		return h.runtime.HandleWorldEntitySpawn(invocation, native.WorldEntityInput{Entity: id})
	})
}

func (h *WorldHandler) HandleEntityDespawn(tx *world.Tx, entity world.Entity) {
	if h.entities == nil || entity == nil {
		return
	}
	handle := entity.H()
	if id, ok := h.entities.ID(entity); ok && h.subscribed(native.WorldEntityDespawnSubscription) {
		h.invoke(tx, "entity-despawn", func(invocation native.InvocationID) error {
			return h.runtime.HandleWorldEntityDespawn(invocation, native.WorldEntityInput{Entity: id})
		})
	}
	if connected, ok := entity.(*player.Player); ok && h.players != nil {
		h.players.recordWorldDeparture(connected, h.world)
	}
	h.entities.deactivateHandle(handle)
	tx.Defer(func(*world.Tx) {
		if handle.Closed() {
			h.entities.unregisterHandle(handle)
			if h.players != nil {
				h.players.forgetWorldDeparture(handle)
			}
		}
	})
}

func (h *WorldHandler) HandleExplosion(ctx *world.Context, position mgl64.Vec3, entities *[]world.Entity, blocks *[]cube.Pos, itemDropChance *float64, spawnFire *bool) {
	if h.entities == nil || !h.subscribed(native.WorldExplosionSubscription) || entities == nil || blocks == nil || itemDropChance == nil || spawnFire == nil {
		return
	}
	entityIDs := make([]native.EntityID, len(*entities))
	for index, entity := range *entities {
		entityIDs[index] = h.entities.Register(entity)
		if entityIDs[index].Generation == 0 {
			h.log.Error("encode native world explosion event", "entity", index)
			return
		}
	}
	blockPositions := make([]native.BlockPos, len(*blocks))
	for index, position := range *blocks {
		blockPositions[index] = nativeBlockPos(position)
	}
	invocation, leave := h.players.BeginInvocation(ctx.Tx)
	defer leave()
	output, err := h.runtime.HandleWorldExplosion(invocation, native.WorldExplosionInput{
		Position: native.Vec3{X: position.X(), Y: position.Y(), Z: position.Z()}, Entities: entityIDs, Blocks: blockPositions,
	}, *itemDropChance, *spawnFire, ctx.Cancelled())
	if err != nil {
		if output.Cancelled {
			ctx.Cancel()
		}
		h.log.Error("native plugin world handler failed", "event", "explosion", "error", err)
		return
	}
	resolved := make([]world.Entity, len(output.Entities))
	for index, id := range output.Entities {
		var ok bool
		resolved[index], ok = h.entities.Resolve(id, ctx.Tx)
		if !ok {
			h.log.Error("native plugin world explosion handler returned an invalid entity", "entity", id)
			return
		}
	}
	replacementBlocks := make([]cube.Pos, len(output.Blocks))
	for index, position := range output.Blocks {
		replacementBlocks[index] = cube.Pos{int(position.X), int(position.Y), int(position.Z)}
	}
	*entities, *blocks = resolved, replacementBlocks
	*itemDropChance, *spawnFire = output.ItemDropChance, output.SpawnFire
	if output.Cancelled {
		ctx.Cancel()
	}
}

func (h *WorldHandler) HandleRedstoneUpdate(ctx *world.Context, update world.RedstoneUpdate) {
	if !h.subscribed(native.WorldRedstoneUpdateSubscription) {
		return
	}
	before, beforeOK := nativeEventBlock(update.Before)
	var after *native.WorldBlock
	if update.After != nil {
		encoded, ok := nativeEventBlock(update.After)
		if !ok {
			h.log.Error("encode native world redstone-update event")
			return
		}
		after = &encoded
	}
	if !beforeOK {
		h.log.Error("encode native world redstone-update event")
		return
	}
	h.cancellable(ctx, "redstone-update", func(invocation native.InvocationID, cancelled bool) (bool, error) {
		return h.runtime.HandleWorldRedstoneUpdate(invocation, native.WorldRedstoneUpdateInput{
			Position: nativeBlockPos(update.Pos), ChangedNeighbour: nativeBlockPos(update.ChangedNeighbour),
			HasChangedNeighbour: update.HasChangedNeighbour, ChangedRedstoneRelevant: update.ChangedRedstoneRelevant,
			Source: nativeBlockPos(update.Source), HasSource: update.HasSource, Before: before, After: after,
			OldPower: update.OldPower, NewPower: update.NewPower, CurrentTick: update.CurrentTick,
			Cause: native.RedstoneUpdateCause(update.Cause),
		}, cancelled)
	})
}

func (h *WorldHandler) HandleClose(tx *world.Tx) {
	if h.subscribed(native.WorldCloseSubscription) {
		h.invoke(tx, "close", h.runtime.HandleWorldClose)
	}
	// Dragonfly replaces the handler immediately after this callback, so chunk
	// entity closure will not emit despawn callbacks. Expire every remaining
	// handle after the plugin has had its final transaction-scoped view.
	if h.entities != nil {
		for entity := range tx.Entities() {
			if connected, ok := entity.(*player.Player); ok && h.players != nil {
				if _, accepted := h.players.ID(connected); accepted {
					// HandleQuit owns accepted-player cleanup and still needs a
					// valid snapshot after the world begins closing.
					continue
				}
			}
			h.entities.unregisterHandle(entity.H())
		}
	}
}

func (h *WorldHandler) subscribed(subscription uint64) bool {
	return h.runtime != nil && h.players != nil && h.runtime.Subscriptions()&subscription != 0
}

func (h *WorldHandler) invoke(tx *world.Tx, event string, handle func(native.InvocationID) error) {
	invocation, leave := h.players.BeginInvocation(tx)
	defer leave()
	if err := handle(invocation); err != nil {
		h.log.Error("native plugin world handler failed", "event", event, "error", err)
	}
}

func (h *WorldHandler) cancellable(ctx *world.Context, event string, handle func(native.InvocationID, bool) (bool, error)) {
	invocation, leave := h.players.BeginInvocation(ctx.Tx)
	defer leave()
	cancelled, err := handle(invocation, ctx.Cancelled())
	if cancelled {
		ctx.Cancel()
	}
	if err != nil {
		h.log.Error("native plugin world handler failed", "event", event, "error", err)
		return
	}
}

func (h *WorldHandler) positionPair(ctx *world.Context, subscription uint64, event string, handle func(native.InvocationID, bool) (bool, error)) {
	if h.subscribed(subscription) {
		h.cancellable(ctx, event, handle)
	}
}

func (h *WorldHandler) position(ctx *world.Context, position cube.Pos, subscription uint64, event string, handle func(native.InvocationID, native.WorldPositionInput, bool) (bool, error)) {
	if !h.subscribed(subscription) {
		return
	}
	h.cancellable(ctx, event, func(invocation native.InvocationID, cancelled bool) (bool, error) {
		return handle(invocation, native.WorldPositionInput{Position: nativeBlockPos(position)}, cancelled)
	})
}

func (h *WorldHandler) runtimeHandleBlockBurn(invocation native.InvocationID, input native.WorldPositionInput, cancelled bool) (bool, error) {
	return h.runtime.HandleWorldBlockBurn(invocation, input, cancelled)
}

func (h *WorldHandler) runtimeHandleCropTrample(invocation native.InvocationID, input native.WorldPositionInput, cancelled bool) (bool, error) {
	return h.runtime.HandleWorldCropTrample(invocation, input, cancelled)
}

func (h *WorldHandler) runtimeHandleLeavesDecay(invocation native.InvocationID, input native.WorldPositionInput, cancelled bool) (bool, error) {
	return h.runtime.HandleWorldLeavesDecay(invocation, input, cancelled)
}
