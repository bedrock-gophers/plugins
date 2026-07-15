package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestPlayerCooldownMatchesPinnedDragonfly(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	directory, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	spec, err := inspectPlayerCooldown(filepath.Join(string(bytes.TrimSpace(directory)), "server", "player", "player.go"))
	if err != nil {
		t.Fatal(err)
	}
	generated := generatePlayerCooldown(spec)
	for _, expected := range [][]byte{
		[]byte("public bool HasCooldown(World.Item? item)"),
		[]byte("public void SetCooldown(World.Item? item, TimeSpan cooldown)"),
	} {
		if !bytes.Contains(generated, expected) {
			t.Fatalf("generated player cooldown API missing %q", expected)
		}
	}
}
