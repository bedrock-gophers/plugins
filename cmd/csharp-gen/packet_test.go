package main

import (
	"bytes"
	"path/filepath"
	"runtime"
	"testing"
)

func TestInspectPacketsGeneratesRegisteredTypesAndTypedFields(t *testing.T) {
	directory := moduleDirectoryForTest(t, "github.com/sandertv/gophertunnel")
	types, err := inspectPackets(filepath.Join(directory, "minecraft", "protocol", "packet"))
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 233 {
		t.Fatalf("packet type count = %d, want 233", len(types))
	}
	var textPacket *packetTypeSpec
	for index := range types {
		if types[index].Name == "Text" {
			textPacket = &types[index]
			break
		}
	}
	if textPacket == nil || textPacket.ID != 9 {
		t.Fatalf("Text packet = %#v", textPacket)
	}
	if got := textPacket.Fields[3]; got.Name != "Message" || got.Type != "string" || got.Kind != packetFieldString {
		t.Fatalf("Text.Message = %#v", got)
	}
	if got := textPacket.Fields[4]; got.Name != "Parameters" || got.Kind != packetFieldValue {
		t.Fatalf("Text.Parameters = %#v", got)
	}
}

func TestGeneratedPacketHandlersUseExactInterceptNames(t *testing.T) {
	output := generatePacketHandler()
	for _, expected := range [][]byte{[]byte("HandleClientPacket"), []byte("HandleServerPacket")} {
		if !bytes.Contains(output, expected) {
			t.Fatalf("generated handler missing %q", expected)
		}
	}
	if bytes.Contains(output, []byte("Incoming")) || bytes.Contains(output, []byte("Outgoing")) {
		t.Fatal("generated handler invented direction aliases")
	}
}

func TestGeneratedPacketsKeepHostHandlesInternal(t *testing.T) {
	types := []packetTypeSpec{{Name: "Text", ID: 9}}
	output := generatePacketTypes(types)
	if !bytes.Contains(output, []byte("ulong Packet.HostHandle() => _handle;")) {
		t.Fatal("generated packet does not implement hidden host handle")
	}
	if bytes.Contains(output, []byte("public ulong HostHandle")) {
		t.Fatal("generated packet exposes host handle publicly")
	}
}

func TestInterceptHandlerASTMatchesPinnedLibrary(t *testing.T) {
	directory := moduleDirectoryForTest(t, "github.com/bedrock-gophers/intercept")
	if err := inspectInterceptHandler(filepath.Join(directory, "intercept", "handler.go")); err != nil {
		t.Fatal(err)
	}
}

func moduleDirectoryForTest(t *testing.T, module string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test source")
	}
	return moduleDirectory(filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..")), module)
}
