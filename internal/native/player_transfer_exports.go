package native

/*
#include "bridge.h"
*/
import "C"

import "math"

//export bg_go_player_transfer
func bg_go_player_transfer(context C.uint64_t, invocation C.DfInvocationId, player C.DfPlayerId, world C.DfWorldId, position C.DfVec3) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	value := nativeEntityVec3(position)
	if !ok || WorldID(world.value) == 0 || math.IsNaN(value.X) || math.IsNaN(value.Y) || math.IsNaN(value.Z) ||
		math.IsInf(value.X, 0) || math.IsInf(value.Y, 0) || math.IsInf(value.Z, 0) ||
		!host.TransferPlayer(InvocationID(invocation), playerID(player), WorldID(world.value), value) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}
