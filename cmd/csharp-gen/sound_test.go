package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInspectSoundsRejectsUnknownConcreteSound(t *testing.T) {
	directory := t.TempDir()
	source := `package sound

type sound struct{}
func (sound) Play(*world.World, mgl64.Vec3) {}
type NewUpstreamSound struct{ sound }
`
	if err := os.WriteFile(filepath.Join(directory, "sound.go"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := inspectSounds(directory)
	if err == nil || !strings.Contains(err.Error(), "NewUpstreamSound") || !strings.Contains(err.Error(), "ABI review") {
		t.Fatalf("inspectSounds() error = %v, want unknown concrete sound ABI review", err)
	}
}

func TestInspectSoundsRejectsPlayImplementationDrift(t *testing.T) {
	directory := t.TempDir()
	source := `package sound

type sound struct{}
func (sound) Play(*world.World, mgl64.Vec3) { panic("changed") }
`
	if err := os.WriteFile(filepath.Join(directory, "sound.go"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := inspectSounds(directory); err == nil || !strings.Contains(err.Error(), "implementation changed") {
		t.Fatalf("inspectSounds() error = %v, want Play implementation drift", err)
	}
}

func TestPinnedDragonflySoundsUseGoAST(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	module, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	types, err := inspectSounds(filepath.Join(
		string(bytes.TrimSpace(module)), "server", "world", "sound"))
	if err != nil {
		t.Fatal(err)
	}
	soundMethod, err := inspectSoundInterface(filepath.Join(
		string(bytes.TrimSpace(module)), "server", "world", "sound.go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 87 || len(selectedSoundTypes) != 87 {
		t.Fatalf("generated %d sounds from %d selected types, want 87", len(types), len(selectedSoundTypes))
	}
	seen := make(map[string]struct{}, len(types))
	for index, definition := range types {
		if definition.Name != selectedSoundTypes[index] {
			t.Fatalf("sound %d = %s, want %s", index, definition.Name, selectedSoundTypes[index])
		}
		if _, ok := seen[definition.Name]; ok {
			t.Fatalf("sound %s appears more than once in the pinned ABI order", definition.Name)
		}
		seen[definition.Name] = struct{}{}
	}
	generated := string(generateSounds(soundMethod, types))
	for _, expected := range []string{
		"public interface Sound { void Play(World w, Vector3 pos); }",
		"record struct AnvilBreak : World.Sound",
		"record struct Attack(bool Damage) : World.Sound",
		"record struct Fall(double Distance) : World.Sound",
		"record struct BlockPlace(World.Block Block) : World.Sound",
		"record struct Note(Instrument Instrument, int Pitch) : World.Sound",
		"record struct EquipItem(World.Item Item) : World.Sound",
		"record struct BucketFill(World.Liquid Liquid) : World.Sound",
		"record struct CrossbowLoad(int Stage, bool QuickCharge) : World.Sound",
		"record struct GoatHorn(Horn Horn) : World.Sound",
		"public void Play(World w, Vector3 pos) { }",
		"86 => new GoatHorn(new Horn(checked((int)data)))",
		"83 => new BucketFill(block is World.Liquid liquid ? liquid",
		"internal static class SoundCodec",
		"case Sound.Attack value:",
		"encoded = new(68u, 0u, 0, value.Damage ? 1u : 0u, 0d, null, null); return true;",
		"case Sound.Note value:",
		"encoded = new(78u, value.Instrument.Id, value.Pitch, 0u, 0d, null, null); return true;",
		"case Sound.EquipItem value:",
		"encoded = new(82u, 0u, 0, 0u, 0d, null, value.Item); return true;",
		"case Sound.BucketFill value when value.Liquid is Block.Water or Block.Lava:",
		"value.Liquid is Block.Water ? 0u : 1u",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated sound output missing %q:\n%s", expected, generated)
		}
	}
	playerMethod, err := inspectPlayerPlaySound(filepath.Join(
		string(bytes.TrimSpace(module)), "server", "player", "player.go"))
	if err != nil {
		t.Fatal(err)
	}
	playerGenerated := string(generatePlayerPlaySound(playerMethod))
	if !strings.Contains(playerGenerated, "void PlaySound(World.Sound sound)") ||
		!strings.Contains(playerGenerated, "PluginBridge.Host.PlayPlayerSound(_invocation, Id, sound)") {
		t.Fatalf("generated Player.PlaySound is incomplete:\n%s", playerGenerated)
	}
}

func TestInspectPlayerPlaySoundRejectsDrift(t *testing.T) {
	for name, source := range map[string]string{
		"parameter type": "package player\nfunc (p *Player) PlaySound(sound world.Particle) {}",
		"result":         "package player\nfunc (p *Player) PlaySound(sound world.Sound) bool { return false }",
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "player.go")
			if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
				t.Fatal(err)
			}
			_, err := inspectPlayerPlaySound(path)
			if err == nil || !strings.Contains(err.Error(), "signature changed") &&
				!strings.Contains(err.Error(), "parameter shape changed") {
				t.Fatalf("expected Player.PlaySound drift error, got %v", err)
			}
		})
	}
}

func TestInspectSoundInterfaceRejectsDrift(t *testing.T) {
	for name, source := range map[string]string{
		"world":    "package world\ntype Sound interface { Play(w *Tx, pos mgl64.Vec3) }",
		"position": "package world\ntype Sound interface { Play(w *World, pos cube.Pos) }",
		"result":   "package world\ntype Sound interface { Play(w *World, pos mgl64.Vec3) bool }",
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "sound.go")
			if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := inspectSoundInterface(path); err == nil || !strings.Contains(err.Error(), "signature changed") {
				t.Fatalf("expected Sound.Play drift error, got %v", err)
			}
		})
	}
}
