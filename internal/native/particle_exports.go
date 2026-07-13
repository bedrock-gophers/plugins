package native

/*
#include "bridge.h"
*/
import "C"

import (
	"unicode/utf8"
)

//export bg_go_world_particle_add
func bg_go_world_particle_add(context C.uint64_t, invocation C.DfInvocationId, worldID C.DfWorldId, position C.DfVec3, view *C.DfParticleViewV1) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	if !ok || view == nil {
		return C.DF_STATUS_ERROR
	}
	kind := ParticleKind(view.kind)
	if kind > ParticleEntityFlame {
		return C.DF_STATUS_ERROR
	}
	value := WorldParticle{
		Kind: kind, Data: uint32(view.data), Pitch: int32(view.pitch),
		Colour: RGBA{R: uint8(view.colour.r), G: uint8(view.colour.g), B: uint8(view.colour.b), A: uint8(view.colour.a)},
		Diff:   BlockPos{X: int32(view.diff.x), Y: int32(view.diff.y), Z: int32(view.diff.z)},
	}
	if view.block != nil {
		identifier, validIdentifier := copyWorldBytes(view.block.identifier, maxBlockIdentifierBytes)
		properties, validProperties := copyWorldBytes(view.block.properties_nbt, maxBlockPropertiesBytes)
		if !validIdentifier || !validProperties || len(identifier) == 0 || !utf8.Valid(identifier) {
			return C.DF_STATUS_ERROR
		}
		value.Block = &WorldBlock{Identifier: string(identifier), PropertiesNBT: properties}
	}
	if !host.AddWorldParticle(InvocationID(invocation), WorldID(worldID.value), nativeEntityVec3(position), value) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}
