package native

/*
#include <stdlib.h>
#include "bridge.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

const (
	PacketClientEvent uint32 = 55 + iota
	PacketServerEvent
)

const (
	PacketClientSubscription uint64 = 1 << (PacketClientEvent - 1 + iota)
	PacketServerSubscription
)

func (r *Runtime) HandlePacket(event uint32, handle PacketHandle, packetID uint32, xuid string, cancelled bool) (bool, error) {
	if r == nil || r.ptr == nil {
		return cancelled, errors.New("native runtime is closed")
	}
	if event != PacketClientEvent && event != PacketServerEvent || handle == 0 {
		return cancelled, errors.New("invalid packet callback")
	}
	xuidData := C.CBytes([]byte(xuid))
	defer C.free(xuidData)
	input := C.DfPacketInput{
		packet: C.uint64_t(handle), packet_id: C.uint32_t(packetID),
		xuid: C.DfStringView{data: (*C.uint8_t)(xuidData), len: C.uint64_t(len(xuid))},
	}
	state := C.DfPacketState{cancelled: C.uint8_t(boolByte(cancelled))}
	status := C.bg_runtime_handle_event(r.ptr, C.DfEventId(event), unsafe.Pointer(&input), unsafe.Pointer(&state))
	runtime.KeepAlive(xuid)
	result := stickyCancellation(cancelled, state.cancelled != 0)
	if status != C.DF_STATUS_OK {
		return result, fmt.Errorf("native packet handler failed with status %d", int32(status))
	}
	return result, nil
}
