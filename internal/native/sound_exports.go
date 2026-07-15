package native

/*
#include "bridge.h"
*/
import "C"

import "unicode/utf8"

//export bg_go_world_sound_play
func bg_go_world_sound_play(context C.uint64_t, invocation C.DfInvocationId, worldID C.DfWorldId, position C.DfVec3, view *C.DfSoundViewV2) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	value, valid := copySoundView(view)
	if !ok || !valid || !host.PlayWorldSound(InvocationID(invocation), WorldID(worldID.value), nativeEntityVec3(position), value) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}

// CallWorldSound invokes one borrowed C# World.Sound.Play callback.
func CallWorldSound(callback WorldSoundCallback, world WorldID, position Vec3) bool {
	if callback.Function == 0 || callback.Context == 0 || world == 0 {
		return false
	}
	return C.bg_call_world_sound(
		C.uintptr_t(callback.Function), C.uintptr_t(callback.Context),
		C.DfWorldId{value: C.uint64_t(world)},
		C.DfVec3{x: C.double(position.X), y: C.double(position.Y), z: C.double(position.Z)},
	) == C.DF_STATUS_OK
}

//export bg_go_player_sound_play
func bg_go_player_sound_play(context C.uint64_t, invocation C.DfInvocationId, player C.DfPlayerId, view *C.DfSoundViewV2) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	value, valid := copySoundView(view)
	if !ok || !valid || !host.PlayPlayerSound(InvocationID(invocation), playerID(player), value) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}

func copySoundView(view *C.DfSoundViewV2) (WorldSound, bool) {
	if view == nil || view.callback == 0 != (view.callback_context == 0) {
		return WorldSound{}, false
	}
	if view.callback != 0 {
		if view.kind != 0 || view.data != 0 || view.integer != 0 || view.flags != 0 || view.scalar != 0 ||
			view.block != nil || view.item != nil {
			return WorldSound{}, false
		}
		callback := WorldSoundCallback{Function: uint64(view.callback), Context: uint64(view.callback_context)}
		return WorldSound{Callback: &callback}, true
	}
	if SoundKind(view.kind) > SoundGoatHorn {
		return WorldSound{}, false
	}
	value := WorldSound{
		Kind: SoundKind(view.kind), Data: uint32(view.data), Integer: int32(view.integer),
		Flags: uint32(view.flags), Scalar: float64(view.scalar),
	}
	if view.block != nil {
		identifier, validIdentifier := copyWorldBytes(view.block.identifier, maxBlockIdentifierBytes)
		properties, validProperties := copyWorldBytes(view.block.properties_nbt, maxBlockPropertiesBytes)
		if !validIdentifier || !validProperties || len(identifier) == 0 || !utf8.Valid(identifier) {
			return WorldSound{}, false
		}
		value.Block = &WorldBlock{Identifier: string(identifier), PropertiesNBT: properties}
	}
	if view.item != nil {
		item, valid := copyItemStackView(view.item)
		if !valid {
			return WorldSound{}, false
		}
		value.Item = &item
	}
	return value, true
}
