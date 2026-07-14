package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const serverSource = `package server
func (srv *Server) MaxPlayerCount() int { return 0 }
func (srv *Server) PlayerCount() int { return 0 }
func (srv *Server) Players(tx *world.Tx) iter.Seq[*player.Player] { return nil }
func (srv *Server) Player(uuid uuid.UUID) (*world.EntityHandle, bool) { return nil, false }
func (srv *Server) PlayerByName(name string) (*world.EntityHandle, bool) { return nil, false }
func (srv *Server) PlayerByXUID(xuid string) (*world.EntityHandle, bool) { return nil, false }
`

func TestServerUsesGoAST(t *testing.T) {
	path := filepath.Join(t.TempDir(), "server.go")
	if err := os.WriteFile(path, []byte(serverSource), 0o600); err != nil {
		t.Fatal(err)
	}
	methods, err := inspectServerMethods(path)
	if err != nil {
		t.Fatal(err)
	}
	generated := string(generateServer(methods))
	for _, expected := range []string{
		"public int MaxPlayerCount() => PluginBridge.Host.ServerMaxPlayerCount()",
		"public int PlayerCount() => PluginBridge.Host.ServerPlayerCount()",
		"public IEnumerable<Player> Players(World.Tx? tx = null)",
		"PluginBridge.Host.ServerPlayers(tx?.Invocation ?? 0)",
		"public (World.EntityHandle? Player, bool Ok) Player(Guid uuid)",
		"PluginBridge.Host.ServerPlayer(uuid)",
		"public (World.EntityHandle? Player, bool Ok) PlayerByName(string name)",
		"PluginBridge.Host.ServerPlayerByName(name)",
		"public (World.EntityHandle? Player, bool Ok) PlayerByXUID(string xuid)",
		"PluginBridge.Host.ServerPlayerByXUID(xuid)",
		"public Server Server() => new();",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated server surface missing %q:\n%s", expected, generated)
		}
	}
}

func TestServerRejectsSignatureDrift(t *testing.T) {
	tests := map[string][2]string{
		"MaxPlayerCount": {"MaxPlayerCount() int", "MaxPlayerCount() int64"},
		"PlayerCount":    {"PlayerCount() int", "PlayerCount(extra int) int"},
		"Players":        {"Players(tx *world.Tx)", "Players(tx world.Tx)"},
		"Player":         {"Player(uuid uuid.UUID)", "Player(uuid string)"},
		"PlayerByName":   {"PlayerByName(name string)", "PlayerByName(name []byte)"},
		"PlayerByXUID":   {"PlayerByXUID(xuid string)", "PlayerByXUID(xuid []byte)"},
	}
	for name, replacement := range tests {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "server.go")
			source := strings.Replace(serverSource, replacement[0], replacement[1], 1)
			if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := inspectServerMethods(path); err == nil || !strings.Contains(err.Error(), "signature changed") {
				t.Fatalf("expected signature drift error, got %v", err)
			}
		})
	}
}

func TestPinnedDragonflyServerHasExactSurface(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	methods, err := inspectServerMethods(filepath.Join(string(bytes.TrimSpace(output)), "server", "server.go"))
	if err != nil {
		t.Fatal(err)
	}
	want := []commandMethod{
		{Name: "MaxPlayerCount", ReturnType: "int"},
		{Name: "PlayerCount", ReturnType: "int"},
		{Name: "Players", Parameters: []parameter{{Name: "tx", Type: "World.Tx?"}}, ReturnType: "IEnumerable<Player>"},
		{Name: "Player", Parameters: []parameter{{Name: "uuid", Type: "Guid"}}, ReturnType: "(World.EntityHandle? Player, bool Ok)"},
		{Name: "PlayerByName", Parameters: []parameter{{Name: "name", Type: "string"}}, ReturnType: "(World.EntityHandle? Player, bool Ok)"},
		{Name: "PlayerByXUID", Parameters: []parameter{{Name: "xuid", Type: "string"}}, ReturnType: "(World.EntityHandle? Player, bool Ok)"},
	}
	if len(methods) != len(want) {
		t.Fatalf("server methods = %d, want %d", len(methods), len(want))
	}
	for index := range want {
		if methods[index].Name != want[index].Name || methods[index].ReturnType != want[index].ReturnType ||
			!equalParameters(methods[index].Parameters, want[index].Parameters) {
			t.Fatalf("server method %d = %+v, want %+v", index, methods[index], want[index])
		}
	}
}
