package native

/*
#include <stdlib.h>
#include "bridge.h"

static inline void bg_native_call_event_drop(DfEventDropFn callback, void *context) {
    if (callback != NULL) callback(context);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
	"unsafe"
)

const (
	WorldLiquidFlowEvent uint32 = 42 + iota
	WorldLiquidDecayEvent
	WorldLiquidHardenEvent
	WorldSoundEvent
	WorldFireSpreadEvent
	WorldBlockBurnEvent
	WorldCropTrampleEvent
	WorldLeavesDecayEvent
	WorldEntitySpawnEvent
	WorldEntityDespawnEvent
	WorldExplosionEvent
	WorldRedstoneUpdateEvent
	WorldCloseEvent
)

const (
	WorldLiquidFlowSubscription uint64 = 1 << (WorldLiquidFlowEvent - 1 + iota)
	WorldLiquidDecaySubscription
	WorldLiquidHardenSubscription
	WorldSoundSubscription
	WorldFireSpreadSubscription
	WorldBlockBurnSubscription
	WorldCropTrampleSubscription
	WorldLeavesDecaySubscription
	WorldEntitySpawnSubscription
	WorldEntityDespawnSubscription
	WorldExplosionSubscription
	WorldRedstoneUpdateSubscription
	WorldCloseSubscription
)

const maxWorldExplosionValues = 1 << 20

// HandleWorldScheduled dispatches or drops one managed callback registered by
// a plugin. Execution receives a callback-scoped Dragonfly transaction.
func (r *Runtime) HandleWorldScheduled(plugin, callback uint64, invocation InvocationID, execute bool) error {
	if r == nil || r.ptr == nil || plugin == 0 || callback == 0 || execute && invocation == 0 {
		return errors.New("invalid scheduled world callback")
	}
	var executeValue C.uint8_t
	if execute {
		executeValue = 1
	}
	if status := C.bg_runtime_handle_scheduled(r.ptr, C.uint64_t(plugin), C.uint64_t(callback), C.DfInvocationId(invocation), executeValue); status != C.DF_STATUS_OK {
		return errors.New("scheduled world callback failed")
	}
	return nil
}

type WorldLiquidFlowInput struct {
	From, Into       BlockPos
	Liquid, Replaced WorldBlock
}

type WorldLiquidDecayInput struct {
	Position BlockPos
	Before   WorldBlock
	After    *WorldBlock
}

type WorldLiquidHardenInput struct {
	Position                    BlockPos
	LiquidHardened, OtherLiquid WorldBlock
	NewBlock                    WorldBlock
}

type WorldSoundInput struct {
	Sound    WorldSound
	Position Vec3
}

type WorldFireSpreadInput struct {
	From, To BlockPos
}

type WorldPositionInput struct {
	Position BlockPos
}

type WorldEntityInput struct {
	Entity EntityID
}

type WorldExplosionInput struct {
	Position Vec3
	Entities []EntityID
	Blocks   []BlockPos
}

type WorldExplosionOutput struct {
	Cancelled      bool
	Entities       []EntityID
	Blocks         []BlockPos
	ItemDropChance float64
	SpawnFire      bool
}

type RedstoneUpdateCause uint8

const (
	RedstoneUpdateCauseBlockUpdate RedstoneUpdateCause = iota
	RedstoneUpdateCauseScheduledTick
	RedstoneUpdateCauseCompilerRebuild
)

type WorldRedstoneUpdateInput struct {
	Position                BlockPos
	ChangedNeighbour        BlockPos
	HasChangedNeighbour     bool
	ChangedRedstoneRelevant bool
	Source                  BlockPos
	HasSource               bool
	Before                  WorldBlock
	After                   *WorldBlock
	OldPower, NewPower      int
	CurrentTick             int64
	Cause                   RedstoneUpdateCause
}

func (r *Runtime) HandleWorldLiquidFlow(invocation InvocationID, input WorldLiquidFlowInput, cancelled bool) (bool, error) {
	arena := &nativeViewArena{}
	defer arena.release()
	liquid, ok := nativeBlockView(input.Liquid, arena)
	if !ok {
		return cancelled, errors.New("encode world liquid-flow liquid")
	}
	replaced, ok := nativeBlockView(input.Replaced, arena)
	if !ok {
		return cancelled, errors.New("encode world liquid-flow replaced block")
	}
	nativeInput := C.DfWorldLiquidFlowInput{
		invocation: C.DfInvocationId(invocation), from: nativeBlockPos(input.From), into: nativeBlockPos(input.Into),
		liquid: liquid, replaced: replaced,
	}
	return r.handleWorldCancellable(WorldLiquidFlowEvent, unsafe.Pointer(&nativeInput), cancelled, "liquid-flow")
}

func (r *Runtime) HandleWorldLiquidDecay(invocation InvocationID, input WorldLiquidDecayInput, cancelled bool) (bool, error) {
	arena := &nativeViewArena{}
	defer arena.release()
	before, ok := nativeBlockView(input.Before, arena)
	if !ok {
		return cancelled, errors.New("encode world liquid-decay before liquid")
	}
	nativeInput := C.DfWorldLiquidDecayInput{
		invocation: C.DfInvocationId(invocation), position: nativeBlockPos(input.Position), before: before,
	}
	if input.After != nil {
		after, ok := nativeOptionalBlockView(input.After, arena)
		if !ok {
			return cancelled, errors.New("encode world liquid-decay after liquid")
		}
		nativeInput.after = after
	}
	return r.handleWorldCancellable(WorldLiquidDecayEvent, unsafe.Pointer(&nativeInput), cancelled, "liquid-decay")
}

func (r *Runtime) HandleWorldLiquidHarden(invocation InvocationID, input WorldLiquidHardenInput, cancelled bool) (bool, error) {
	arena := &nativeViewArena{}
	defer arena.release()
	hardened, ok := nativeBlockView(input.LiquidHardened, arena)
	if !ok {
		return cancelled, errors.New("encode world liquid-harden liquid")
	}
	other, ok := nativeBlockView(input.OtherLiquid, arena)
	if !ok {
		return cancelled, errors.New("encode world liquid-harden other liquid")
	}
	newBlock, ok := nativeBlockView(input.NewBlock, arena)
	if !ok {
		return cancelled, errors.New("encode world liquid-harden block")
	}
	nativeInput := C.DfWorldLiquidHardenInput{
		invocation: C.DfInvocationId(invocation), position: nativeBlockPos(input.Position),
		liquid_hardened: hardened, other_liquid: other, new_block: newBlock,
	}
	return r.handleWorldCancellable(WorldLiquidHardenEvent, unsafe.Pointer(&nativeInput), cancelled, "liquid-harden")
}

func (r *Runtime) HandleWorldSound(invocation InvocationID, input WorldSoundInput, cancelled bool) (bool, error) {
	arena := &nativeViewArena{}
	defer arena.release()
	sound, ok := nativeWorldSoundView(input.Sound, arena)
	if !ok {
		return cancelled, errors.New("encode world sound")
	}
	nativeInput := C.DfWorldSoundInput{
		invocation: C.DfInvocationId(invocation), sound: sound,
		position: C.DfVec3{x: C.double(input.Position.X), y: C.double(input.Position.Y), z: C.double(input.Position.Z)},
	}
	return r.handleWorldCancellable(WorldSoundEvent, unsafe.Pointer(&nativeInput), cancelled, "sound")
}

func (r *Runtime) HandleWorldFireSpread(invocation InvocationID, input WorldFireSpreadInput, cancelled bool) (bool, error) {
	nativeInput := C.DfWorldFireSpreadInput{
		invocation: C.DfInvocationId(invocation), from: nativeBlockPos(input.From), to: nativeBlockPos(input.To),
	}
	return r.handleWorldCancellable(WorldFireSpreadEvent, unsafe.Pointer(&nativeInput), cancelled, "fire-spread")
}

func (r *Runtime) HandleWorldBlockBurn(invocation InvocationID, input WorldPositionInput, cancelled bool) (bool, error) {
	return r.handleWorldPosition(WorldBlockBurnEvent, invocation, input, cancelled, "block-burn")
}

func (r *Runtime) HandleWorldCropTrample(invocation InvocationID, input WorldPositionInput, cancelled bool) (bool, error) {
	return r.handleWorldPosition(WorldCropTrampleEvent, invocation, input, cancelled, "crop-trample")
}

func (r *Runtime) HandleWorldLeavesDecay(invocation InvocationID, input WorldPositionInput, cancelled bool) (bool, error) {
	return r.handleWorldPosition(WorldLeavesDecayEvent, invocation, input, cancelled, "leaves-decay")
}

func (r *Runtime) HandleWorldEntitySpawn(invocation InvocationID, input WorldEntityInput) error {
	return r.handleWorldEntityNotification(WorldEntitySpawnEvent, invocation, input, "entity-spawn")
}

func (r *Runtime) HandleWorldEntityDespawn(invocation InvocationID, input WorldEntityInput) error {
	return r.handleWorldEntityNotification(WorldEntityDespawnEvent, invocation, input, "entity-despawn")
}

func (r *Runtime) HandleWorldExplosion(invocation InvocationID, input WorldExplosionInput, itemDropChance float64, spawnFire, cancelled bool) (WorldExplosionOutput, error) {
	output := WorldExplosionOutput{
		Cancelled: cancelled, Entities: append([]EntityID(nil), input.Entities...), Blocks: append([]BlockPos(nil), input.Blocks...),
		ItemDropChance: itemDropChance, SpawnFire: spawnFire,
	}
	if r == nil || r.ptr == nil {
		return output, errors.New("native runtime is closed")
	}
	if len(input.Entities) > maxWorldExplosionValues || len(input.Blocks) > maxWorldExplosionValues {
		return output, errors.New("world explosion input exceeds limits")
	}
	entities, releaseEntities, ok := nativeEntityIDs(input.Entities)
	defer releaseEntities()
	if !ok {
		return output, errors.New("encode world explosion entities")
	}
	blocks, releaseBlocks, ok := nativeBlockPositions(input.Blocks)
	defer releaseBlocks()
	if !ok {
		return output, errors.New("encode world explosion blocks")
	}
	nativeInput := C.DfWorldExplosionInput{
		invocation: C.DfInvocationId(invocation),
		position:   C.DfVec3{x: C.double(input.Position.X), y: C.double(input.Position.Y), z: C.double(input.Position.Z)},
		entities:   entities, entity_count: C.uint64_t(len(input.Entities)),
		blocks: blocks, block_count: C.uint64_t(len(input.Blocks)),
	}
	state := C.DfWorldExplosionState{
		cancelled: C.uint8_t(boolByte(cancelled)), spawn_fire: C.uint8_t(boolByte(spawnFire)), item_drop_chance: C.double(itemDropChance),
	}
	status := C.bg_runtime_handle_event(r.ptr, C.DfEventId(WorldExplosionEvent), unsafe.Pointer(&nativeInput), unsafe.Pointer(&state))
	if state.replacement_drop != nil {
		defer C.bg_native_call_event_drop(state.replacement_drop, state.replacement_context)
	}
	output.Cancelled = stickyCancellation(cancelled, state.cancelled != 0)
	if status != C.DF_STATUS_OK {
		return output, fmt.Errorf("native world explosion handler failed with status %d", int32(status))
	}
	if state.spawn_fire > 1 {
		return output, errors.New("native world explosion handler returned invalid spawn-fire state")
	}
	hasReplacement := state.replacement_drop != nil
	replacementFields := state.replacement_entities != nil || state.replacement_entity_count != 0 ||
		state.replacement_blocks != nil || state.replacement_block_count != 0 || state.replacement_context != nil
	if !hasReplacement && replacementFields {
		return output, errors.New("native world explosion handler returned unowned replacement values")
	}
	if hasReplacement {
		replacementEntities, ok := copyNativeEntityIDs(state.replacement_entities, uint64(state.replacement_entity_count))
		if !ok {
			return output, errors.New("native world explosion handler returned invalid entities")
		}
		replacementBlocks, ok := copyNativeBlockPositions(state.replacement_blocks, uint64(state.replacement_block_count))
		if !ok {
			return output, errors.New("native world explosion handler returned invalid blocks")
		}
		output.Entities, output.Blocks = replacementEntities, replacementBlocks
	}
	output.ItemDropChance = float64(state.item_drop_chance)
	output.SpawnFire = state.spawn_fire != 0
	return output, nil
}

func (r *Runtime) HandleWorldRedstoneUpdate(invocation InvocationID, input WorldRedstoneUpdateInput, cancelled bool) (bool, error) {
	if input.OldPower < math.MinInt32 || input.OldPower > math.MaxInt32 || input.NewPower < math.MinInt32 || input.NewPower > math.MaxInt32 ||
		input.Cause > RedstoneUpdateCauseCompilerRebuild {
		return cancelled, errors.New("world redstone-update input exceeds limits")
	}
	arena := &nativeViewArena{}
	defer arena.release()
	before, ok := nativeBlockView(input.Before, arena)
	if !ok {
		return cancelled, errors.New("encode world redstone-update before block")
	}
	nativeInput := C.DfWorldRedstoneUpdateInput{
		invocation: C.DfInvocationId(invocation), position: nativeBlockPos(input.Position),
		changed_neighbour: nativeBlockPos(input.ChangedNeighbour), has_changed_neighbour: C.uint8_t(boolByte(input.HasChangedNeighbour)),
		changed_redstone_relevant: C.uint8_t(boolByte(input.ChangedRedstoneRelevant)), source: nativeBlockPos(input.Source),
		has_source: C.uint8_t(boolByte(input.HasSource)), before: before, old_power: C.int32_t(input.OldPower),
		new_power: C.int32_t(input.NewPower), current_tick: C.int64_t(input.CurrentTick), cause: C.uint32_t(input.Cause),
	}
	if input.After != nil {
		after, ok := nativeOptionalBlockView(input.After, arena)
		if !ok {
			return cancelled, errors.New("encode world redstone-update after block")
		}
		nativeInput.after = after
	}
	return r.handleWorldCancellable(WorldRedstoneUpdateEvent, unsafe.Pointer(&nativeInput), cancelled, "redstone-update")
}

func (r *Runtime) HandleWorldClose(invocation InvocationID) error {
	if r == nil || r.ptr == nil {
		return errors.New("native runtime is closed")
	}
	input := C.DfWorldCloseInput{invocation: C.DfInvocationId(invocation)}
	var state C.DfWorldNotificationState
	if status := C.bg_runtime_handle_event(r.ptr, C.DfEventId(WorldCloseEvent), unsafe.Pointer(&input), unsafe.Pointer(&state)); status != C.DF_STATUS_OK {
		return fmt.Errorf("native world close handler failed with status %d", int32(status))
	}
	return nil
}

func (r *Runtime) handleWorldPosition(event uint32, invocation InvocationID, input WorldPositionInput, cancelled bool, name string) (bool, error) {
	nativeInput := C.DfWorldPositionInput{invocation: C.DfInvocationId(invocation), position: nativeBlockPos(input.Position)}
	return r.handleWorldCancellable(event, unsafe.Pointer(&nativeInput), cancelled, name)
}

func (r *Runtime) handleWorldEntityNotification(event uint32, invocation InvocationID, input WorldEntityInput, name string) error {
	if r == nil || r.ptr == nil {
		return errors.New("native runtime is closed")
	}
	if input.Entity.Generation == 0 {
		return fmt.Errorf("world %s entity is invalid", name)
	}
	var nativeInput C.DfWorldEntityInput
	nativeInput.invocation = C.DfInvocationId(invocation)
	fillEntityID(&nativeInput.entity, input.Entity)
	var state C.DfWorldNotificationState
	if status := C.bg_runtime_handle_event(r.ptr, C.DfEventId(event), unsafe.Pointer(&nativeInput), unsafe.Pointer(&state)); status != C.DF_STATUS_OK {
		return fmt.Errorf("native world %s handler failed with status %d", name, int32(status))
	}
	return nil
}

func (r *Runtime) handleWorldCancellable(event uint32, input unsafe.Pointer, cancelled bool, name string) (bool, error) {
	if r == nil || r.ptr == nil {
		return cancelled, errors.New("native runtime is closed")
	}
	state := C.DfWorldCancellableState{cancelled: C.uint8_t(boolByte(cancelled))}
	status := C.bg_runtime_handle_event(r.ptr, C.DfEventId(event), input, unsafe.Pointer(&state))
	result := stickyCancellation(cancelled, state.cancelled != 0)
	if status != C.DF_STATUS_OK {
		return result, fmt.Errorf("native world %s handler failed with status %d", name, int32(status))
	}
	return result, nil
}

func nativeWorldSoundView(value WorldSound, arena *nativeViewArena) (C.DfSoundViewV1, bool) {
	if value.Kind > SoundGoatHorn || !finiteWorldEventFloat(value.Scalar) {
		return C.DfSoundViewV1{}, false
	}
	view := C.DfSoundViewV1{
		kind: C.uint32_t(value.Kind), data: C.uint32_t(value.Data), integer: C.int32_t(value.Integer),
		flags: C.uint32_t(value.Flags), scalar: C.double(value.Scalar),
	}
	if value.Block != nil {
		block, ok := nativeBlockView(*value.Block, arena)
		if !ok {
			return C.DfSoundViewV1{}, false
		}
		pointer, ok := arena.allocate(C.sizeof_DfBlockView)
		if !ok {
			return C.DfSoundViewV1{}, false
		}
		*(*C.DfBlockView)(pointer) = block
		view.block = (*C.DfBlockView)(pointer)
	}
	if value.Item != nil {
		item, ok := nativeItemStackView(*value.Item, arena)
		if !ok {
			return C.DfSoundViewV1{}, false
		}
		pointer, ok := arena.allocate(C.sizeof_DfItemStackViewV3)
		if !ok {
			return C.DfSoundViewV1{}, false
		}
		*(*C.DfItemStackViewV3)(pointer) = item
		view.item = (*C.DfItemStackViewV3)(pointer)
	}
	return view, true
}

func nativeOptionalBlockView(value *WorldBlock, arena *nativeViewArena) (*C.DfBlockView, bool) {
	if value == nil {
		return nil, false
	}
	view, ok := nativeBlockView(*value, arena)
	if !ok {
		return nil, false
	}
	pointer, ok := arena.allocate(C.sizeof_DfBlockView)
	if !ok {
		return nil, false
	}
	*(*C.DfBlockView)(pointer) = view
	return (*C.DfBlockView)(pointer), true
}

func nativeEntityIDs(values []EntityID) (*C.DfEntityId, func(), bool) {
	if len(values) == 0 {
		return nil, func() {}, true
	}
	pointer := C.malloc(C.size_t(len(values)) * C.size_t(C.sizeof_DfEntityId))
	if pointer == nil {
		return nil, func() {}, false
	}
	result := (*C.DfEntityId)(pointer)
	for index, value := range values {
		if value.Generation == 0 {
			C.free(pointer)
			return nil, func() {}, false
		}
		fillEntityID(&unsafe.Slice(result, len(values))[index], value)
	}
	return result, func() { C.free(pointer) }, true
}

func nativeBlockPositions(values []BlockPos) (*C.DfBlockPos, func(), bool) {
	if len(values) == 0 {
		return nil, func() {}, true
	}
	pointer := C.malloc(C.size_t(len(values)) * C.size_t(C.sizeof_DfBlockPos))
	if pointer == nil {
		return nil, func() {}, false
	}
	result := (*C.DfBlockPos)(pointer)
	positions := unsafe.Slice(result, len(values))
	for index, value := range values {
		positions[index] = nativeBlockPos(value)
	}
	return result, func() { C.free(pointer) }, true
}

func copyNativeEntityIDs(values *C.DfEntityId, count uint64) ([]EntityID, bool) {
	if count > maxWorldExplosionValues || count > uint64(^uint(0)>>1) || count != 0 && values == nil {
		return nil, false
	}
	result := make([]EntityID, int(count))
	for index, value := range unsafe.Slice(values, int(count)) {
		result[index] = entityID(value)
		if result[index].Generation == 0 {
			return nil, false
		}
	}
	return result, true
}

func copyNativeBlockPositions(values *C.DfBlockPos, count uint64) ([]BlockPos, bool) {
	if count > maxWorldExplosionValues || count > uint64(^uint(0)>>1) || count != 0 && values == nil {
		return nil, false
	}
	result := make([]BlockPos, int(count))
	for index, value := range unsafe.Slice(values, int(count)) {
		result[index] = BlockPos{X: int32(value.x), Y: int32(value.y), Z: int32(value.z)}
	}
	return result, true
}

func finiteWorldEventFloat(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
