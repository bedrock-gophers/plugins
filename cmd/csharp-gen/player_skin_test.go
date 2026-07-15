package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestPlayerSkinMatchesPinnedDragonfly(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	directory, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	spec, err := inspectPlayerSkin(filepath.Join(
		string(bytes.TrimSpace(directory)), "server", "player", "player.go"))
	if err != nil {
		t.Fatal(err)
	}
	generated := generatePlayerSkin(spec)
	for _, expected := range [][]byte{
		[]byte("public Skin Skin()"),
		[]byte("public void SetSkin(Skin skin)"),
	} {
		if !bytes.Contains(generated, expected) {
			t.Fatalf("generated player skin API missing %q", expected)
		}
	}
}
