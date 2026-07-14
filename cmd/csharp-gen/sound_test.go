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
	generated := string(generateSounds(types))
	for _, expected := range []string{
		"record struct AnvilBreak : World.Sound",
		"record struct Attack(bool Damage) : World.Sound",
		"record struct Fall(double Distance) : World.Sound",
		"record struct BlockPlace(World.Block Block) : World.Sound",
		"record struct Note(Instrument Instrument, int Pitch) : World.Sound",
		"record struct EquipItem(World.Item Item) : World.Sound",
		"record struct BucketFill(World.Liquid Liquid) : World.Sound",
		"record struct CrossbowLoad(int Stage, bool QuickCharge) : World.Sound",
		"record struct GoatHorn(Horn Horn) : World.Sound",
		"86 => new GoatHorn(new Horn(checked((int)data)))",
		"83 => new BucketFill(block is World.Liquid liquid ? liquid",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated sound output missing %q:\n%s", expected, generated)
		}
	}
}
