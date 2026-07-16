package native

/*
#include "bridge.h"
*/
import "C"

import "unicode/utf8"

const maxEntityAnimationStringBytes = 64 << 10

//export bg_go_world_entity_animation
func bg_go_world_entity_animation(context C.uint64_t, invocation C.DfInvocationId, entity C.DfEntityId, input *C.DfEntityAnimationView) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	if !ok || input == nil {
		return C.DF_STATUS_ERROR
	}
	name, nameOK := copyNativeBytes(input.name, maxEntityAnimationStringBytes)
	nextState, nextStateOK := copyNativeBytes(input.next_state, maxEntityAnimationStringBytes)
	controller, controllerOK := copyNativeBytes(input.controller, maxEntityAnimationStringBytes)
	stopCondition, stopConditionOK := copyNativeBytes(input.stop_condition, maxEntityAnimationStringBytes)
	if !nameOK || !nextStateOK || !controllerOK || !stopConditionOK ||
		!utf8.Valid(name) || !utf8.Valid(nextState) || !utf8.Valid(controller) || !utf8.Valid(stopCondition) {
		return C.DF_STATUS_ERROR
	}
	if !host.PlayEntityAnimation(InvocationID(invocation), entityID(entity), WorldEntityAnimation{
		Name: string(name), NextState: string(nextState), Controller: string(controller), StopCondition: string(stopCondition),
	}) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}
