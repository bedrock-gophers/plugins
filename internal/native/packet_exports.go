package native

/*
#include "bridge.h"
*/
import "C"

import "unsafe"

const maxPacketFieldBytes = 16 << 20

type packetHost interface {
	PacketField(PacketHandle, uint32) (PacketFieldValue, bool)
	SetPacketField(PacketHandle, uint32, PacketFieldValue) bool
}

//export bg_go_packet_field_get
func bg_go_packet_field_get(context C.uint64_t, packet C.uint64_t, field C.uint32_t, output *C.DfPacketFieldValue) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	packets, supported := host.(packetHost)
	if !ok || !supported || output == nil {
		return C.DF_STATUS_ERROR
	}
	buffer := output.data
	value, ok := packets.PacketField(PacketHandle(packet), uint32(field))
	if !ok || len(value.Data) > maxPacketFieldBytes {
		return C.DF_STATUS_ERROR
	}
	*output = C.DfPacketFieldValue{
		kind: C.uint32_t(value.Kind), signed_value: C.int64_t(value.Signed),
		unsigned_value: C.uint64_t(value.Unsigned), number: C.double(value.Number),
		x: C.double(value.X), y: C.double(value.Y), z: C.double(value.Z), data: buffer,
	}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(&output.uuid.bytes[0])), len(value.UUID)), value.UUID[:])
	output.data.len = C.uint64_t(len(value.Data))
	if len(value.Data) != 0 && uint64(buffer.capacity) >= uint64(len(value.Data)) && buffer.data != nil {
		copy(unsafe.Slice((*byte)(unsafe.Pointer(buffer.data)), len(value.Data)), value.Data)
	}
	return C.DF_STATUS_OK
}

//export bg_go_packet_field_set
func bg_go_packet_field_set(context C.uint64_t, packet C.uint64_t, field C.uint32_t, input *C.DfPacketFieldValue) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	packets, supported := host.(packetHost)
	if !ok || !supported || input == nil || uint64(input.data.len) > maxPacketFieldBytes || input.data.len != 0 && input.data.data == nil {
		return C.DF_STATUS_ERROR
	}
	value := PacketFieldValue{
		Kind: PacketFieldKind(input.kind), Signed: int64(input.signed_value),
		Unsigned: uint64(input.unsigned_value), Number: float64(input.number),
		X: float64(input.x), Y: float64(input.y), Z: float64(input.z),
	}
	copy(value.UUID[:], unsafe.Slice((*byte)(unsafe.Pointer(&input.uuid.bytes[0])), len(value.UUID)))
	if input.data.len != 0 {
		value.Data = append([]byte(nil), unsafe.Slice((*byte)(unsafe.Pointer(input.data.data)), int(input.data.len))...)
	}
	if !packets.SetPacketField(PacketHandle(packet), uint32(field), value) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}
