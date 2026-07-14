package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const playerIdentitySource = `package player
func (p *Player) Name() string { return "" }
func (p *Player) UUID() uuid.UUID { return uuid.Nil }
func (p *Player) XUID() string { return "" }
`

func TestPlayerIdentityUsesGoAST(t *testing.T) {
	path := filepath.Join(t.TempDir(), "player.go")
	if err := os.WriteFile(path, []byte(playerIdentitySource), 0o600); err != nil {
		t.Fatal(err)
	}
	methods, err := inspectPlayerIdentityMethods(path)
	if err != nil {
		t.Fatal(err)
	}
	generated := string(generatePlayerIdentityMethods(methods))
	for _, expected := range []string{
		"public string Name() => PlayerName;",
		"public Guid UUID() => PluginBridge.Host.PlayerUUID(Id);",
		"public string XUID() => PluginBridge.Host.PlayerXUID(_invocation, Id);",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated player identity surface missing %q:\n%s", expected, generated)
		}
	}
}

func TestPlayerIdentityRejectsSignatureDrift(t *testing.T) {
	tests := map[string][2]string{
		"Name": {"Name() string", "Name() []byte"},
		"UUID": {"UUID() uuid.UUID", "UUID() string"},
		"XUID": {"XUID() string", "XUID() int64"},
	}
	for name, replacement := range tests {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "player.go")
			source := strings.Replace(playerIdentitySource, replacement[0], replacement[1], 1)
			if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := inspectPlayerIdentityMethods(path); err == nil || !strings.Contains(err.Error(), "signature changed") {
				t.Fatalf("expected signature drift error, got %v", err)
			}
		})
	}
}

func TestPinnedDragonflyPlayerHasExactIdentitySurface(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	methods, err := inspectPlayerIdentityMethods(filepath.Join(string(bytes.TrimSpace(output)), "server", "player", "player.go"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.Join(methods, ","), strings.Join(selectedPlayerIdentityMethods, ","); got != want {
		t.Fatalf("player identity methods = %q, want %q", got, want)
	}
}
