package host

import (
	"testing"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestPacketsBorrowReadMutateAndRelease(t *testing.T) {
	registry := NewPackets(nil)
	value := &packet.Text{Message: "before", Parameters: []string{"one", "two"}}
	handle, release, ok := registry.Borrow(value, true)
	if !ok || handle == 0 {
		t.Fatal("Borrow failed")
	}
	message, ok := registry.PacketField(handle, 3)
	if !ok || message.Kind != native.PacketFieldString || string(message.Data) != "before" {
		t.Fatalf("message = %#v, %v", message, ok)
	}
	if !registry.SetPacketField(handle, 3, native.PacketFieldValue{Kind: native.PacketFieldString, Data: []byte("after")}) {
		t.Fatal("SetPacketField failed")
	}
	if value.Message != "after" {
		t.Fatalf("Message = %q, want after", value.Message)
	}
	parameters, ok := registry.PacketField(handle, 4)
	if !ok || parameters.Kind != native.PacketFieldJSON || string(parameters.Data) != `["one","two"]` {
		t.Fatalf("parameters = %#v, %v", parameters, ok)
	}
	release()
	if _, ok := registry.PacketField(handle, 3); ok {
		t.Fatal("released packet handle remained valid")
	}
}

func TestPacketsRejectWrongFieldKindsAndOverflow(t *testing.T) {
	registry := NewPackets(nil)
	value := &packet.Text{TextType: 1}
	handle, release, ok := registry.Borrow(value, true)
	if !ok {
		t.Fatal("Borrow failed")
	}
	defer release()
	if registry.SetPacketField(handle, 0, native.PacketFieldValue{Kind: native.PacketFieldUnsigned, Unsigned: 256}) {
		t.Fatal("accepted byte overflow")
	}
	if registry.SetPacketField(handle, 0, native.PacketFieldValue{Kind: native.PacketFieldString, Data: []byte("bad")}) {
		t.Fatal("accepted wrong field kind")
	}
}

func TestPacketsRejectOutgoingMutation(t *testing.T) {
	registry := NewPackets(nil)
	value := &packet.Text{Message: "shared"}
	handle, release, ok := registry.Borrow(value, false)
	if !ok {
		t.Fatal("Borrow failed")
	}
	defer release()
	if registry.SetPacketField(handle, 3, native.PacketFieldValue{Kind: native.PacketFieldString, Data: []byte("changed")}) {
		t.Fatal("outgoing packet mutation was accepted")
	}
	if value.Message != "shared" {
		t.Fatalf("Message = %q, want shared", value.Message)
	}
}

func TestPacketsWriteBorrowedPacketToPlayer(t *testing.T) {
	withPlayerTx(t, func(tx *world.Tx, connected *player.Player) {
		players := NewPlayers()
		id := players.Register(connected, 91)
		invocation, end := players.BeginInvocation(tx)
		defer end()
		registry := NewPackets(players)
		handle, release, ok := registry.Borrow(&packet.Text{Message: "hello"}, true)
		if !ok || !registry.WritePlayerPacket(invocation, id, handle) {
			t.Fatal("live borrowed packet was not written")
		}
		release()
		if registry.WritePlayerPacket(invocation, id, handle) {
			t.Fatal("released packet handle was written")
		}
	})
}
