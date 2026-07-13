package native

/*
#include "bridge.h"
*/
import "C"

// These transitional exports keep the Go bridge linkable until Task 4 wires
// detached entity operations. Task 4 must replace or delete this file.

//export bg_go_world_entity_remove
func bg_go_world_entity_remove(C.uint64_t, C.DfInvocationId, C.DfWorldId, C.DfEntityId, *C.DfDetachedEntityId) C.DfStatus {
	return C.DF_STATUS_ERROR
}

//export bg_go_world_entity_add
func bg_go_world_entity_add(C.uint64_t, C.DfInvocationId, C.DfWorldId, C.DfDetachedEntityId, *C.DfVec3, *C.DfEntityId) C.DfStatus {
	return C.DF_STATUS_ERROR
}

//export bg_go_detached_entity_drop
func bg_go_detached_entity_drop(C.uint64_t, C.DfDetachedEntityId) {}
