package native

/*
#include "bridge.h"
*/
import "C"

import (
	"unsafe"
)

//export bg_go_player_effects
func bg_go_player_effects(context C.uint64_t, invocation C.DfInvocationId, player C.DfPlayerId, output *C.DfEffectBuffer) C.DfStatus {
	if output == nil {
		return C.DF_STATUS_ERROR
	}
	output.len = 0
	if !validPlayerEffectBuffer(unsafe.Pointer(output.data), uint64(output.capacity)) {
		return C.DF_STATUS_ERROR
	}
	host, ok := resolveHost(uint64(context))
	if !ok {
		return C.DF_STATUS_ERROR
	}
	values, ok := host.PlayerEffects(InvocationID(invocation), playerID(player))
	if !ok {
		return C.DF_STATUS_ERROR
	}
	encoded, ok := playerEffectSnapshotViews(values)
	if !ok {
		return C.DF_STATUS_ERROR
	}
	required, fits := writeBoundedSnapshot(encoded, uint64(output.capacity), func(values []C.DfEffectView) {
		if len(values) != 0 {
			copy(unsafe.Slice(output.data, len(values)), values)
		}
	})
	output.len = C.uint64_t(required)
	if !fits {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}

func validPlayerEffectBuffer(data unsafe.Pointer, capacity uint64) bool {
	return capacity == 0 || data != nil
}

func playerEffectSnapshotViews(values []PlayerEffect) ([]C.DfEffectView, bool) {
	encoded := make([]C.DfEffectView, len(values))
	for index, value := range values {
		view, valid := playerEffectSnapshotView(value)
		if !valid {
			return nil, false
		}
		encoded[index] = view
	}
	return encoded, true
}

func writeBoundedSnapshot[T any](values []T, capacity uint64, write func([]T)) (uint64, bool) {
	required := uint64(len(values))
	if capacity < required {
		return required, false
	}
	write(values)
	return required, true
}

func playerEffectSnapshotView(value PlayerEffect) (C.DfEffectView, bool) {
	return C.DfEffectView{
		effect_type: C.int32_t(value.Type), level: C.int32_t(value.Level),
		duration_nanoseconds: C.int64_t(value.Duration), potency: C.double(value.Potency),
		ambient: C.uint8_t(boolByte(value.Ambient)), particles_hidden: C.uint8_t(boolByte(value.ParticlesHidden)),
		infinite: C.uint8_t(boolByte(value.Infinite)), tick: C.int64_t(value.Tick),
	}, true
}

//export bg_go_player_effects_clear
func bg_go_player_effects_clear(context C.uint64_t, invocation C.DfInvocationId, player C.DfPlayerId) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	if !ok || !host.ClearPlayerEffects(InvocationID(invocation), playerID(player)) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}
