package native

/*
#include "bridge.h"
*/
import "C"

import "unicode/utf8"

//export bg_go_world_sound_play
func bg_go_world_sound_play(context C.uint64_t, invocation C.DfInvocationId, worldID C.DfWorldId, position C.DfVec3, view *C.DfSoundViewV1) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	value, valid := copySoundView(view)
	if !ok || !valid || !host.PlayWorldSound(InvocationID(invocation), WorldID(worldID.value), nativeEntityVec3(position), value) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}

//export bg_go_player_sound_play
func bg_go_player_sound_play(context C.uint64_t, invocation C.DfInvocationId, player C.DfPlayerId, view *C.DfSoundViewV1) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	value, valid := copySoundView(view)
	if !ok || !valid || !host.PlayPlayerSound(InvocationID(invocation), playerID(player), value) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}

func copySoundView(view *C.DfSoundViewV1) (WorldSound, bool) {
	if view == nil || SoundKind(view.kind) > SoundGoatHorn {
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
