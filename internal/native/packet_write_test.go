package native

import "testing"

type recordingPacketWriteHost struct {
	noopHost
	invocation InvocationID
	player     PlayerID
	packet     PacketHandle
}

func (h *recordingPacketWriteHost) WritePlayerPacket(invocation InvocationID, player PlayerID, packet PacketHandle) bool {
	h.invocation, h.player, h.packet = invocation, player, packet
	return true
}

func TestWritePlayerPacketRoutesExactHandles(t *testing.T) {
	host := &recordingPacketWriteHost{}
	context := registerHost(host)
	t.Cleanup(func() { unregisterHost(context) })
	wantPlayer := PlayerID{UUID: [16]byte{1, 2, 3}, Generation: 4}
	if !writePlayerPacket(context, 7, wantPlayer, 9) {
		t.Fatal("valid packet write rejected")
	}
	if host.invocation != 7 || host.player != wantPlayer || host.packet != 9 {
		t.Fatalf("packet write = %d, %+v, %d", host.invocation, host.player, host.packet)
	}
	if writePlayerPacket(context, 7, wantPlayer, 0) {
		t.Fatal("zero packet handle accepted")
	}
	if writePlayerPacket(^uint64(0), 7, wantPlayer, 9) {
		t.Fatal("unknown host accepted")
	}
}
