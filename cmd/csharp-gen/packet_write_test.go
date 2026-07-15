package main

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestUnsafeWritePacketASTMatchesPinnedLibrary(t *testing.T) {
	directory := moduleDirectoryForTest(t, "github.com/bedrock-gophers/unsafe")
	if err := inspectUnsafeWritePacket(filepath.Join(directory, "unsafe.go")); err != nil {
		t.Fatal(err)
	}
}

func TestGeneratedPlayerWritePacketUsesHiddenHandle(t *testing.T) {
	output := generatePlayerWritePacket()
	for _, expected := range [][]byte{
		[]byte("public void WritePacket(Packet.Packet pk)"),
		[]byte("pk.HostHandle()"),
	} {
		if !bytes.Contains(output, expected) {
			t.Fatalf("generated Player.WritePacket missing %q", expected)
		}
	}
	if bytes.Contains(output, []byte("ulong packet")) {
		t.Fatal("generated public method exposes raw packet handle")
	}
}
