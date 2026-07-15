package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestInspectDifficulties(t *testing.T) {
	directory := difficultyDragonflyDirectory(t)
	spec, err := inspectDifficulties(filepath.Join(directory, "server", "world", "difficulty.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(spec.Methods, difficultyMethods) {
		t.Fatalf("methods = %+v", spec.Methods)
	}
	want := []difficultyValue{
		{Name: "DifficultyPeaceful", PrivateType: "difficultyPeaceful", ID: 0, FoodRegenerates: true, StarvationHealthLimit: 20, FireSpreadIncrease: 0},
		{Name: "DifficultyEasy", PrivateType: "difficultyEasy", ID: 1, StarvationHealthLimit: 10, FireSpreadIncrease: 7},
		{Name: "DifficultyNormal", PrivateType: "difficultyNormal", ID: 2, StarvationHealthLimit: 2, FireSpreadIncrease: 14},
		{Name: "DifficultyHard", PrivateType: "difficultyHard", ID: 3, StarvationHealthLimit: -1, FireSpreadIncrease: 21},
	}
	if !reflect.DeepEqual(spec.Values, want) {
		t.Fatalf("values = %+v, want %+v", spec.Values, want)
	}
}

func TestGenerateDifficulties(t *testing.T) {
	directory := difficultyDragonflyDirectory(t)
	spec, err := inspectDifficulties(filepath.Join(directory, "server", "world", "difficulty.go"))
	if err != nil {
		t.Fatal(err)
	}
	generated := string(generateDifficulties(spec))
	for _, expected := range []string{
		"public interface Difficulty",
		"bool FoodRegenerates();",
		"double StarvationHealthLimit();",
		"int FireSpreadIncrease();",
		"public static readonly Difficulty DifficultyPeaceful = new BuiltinDifficulty(0, true, 20d, 0);",
		"public static readonly Difficulty DifficultyHard = new BuiltinDifficulty(3, false, -1d, 21);",
		"public static (Difficulty Difficulty, bool Ok) DifficultyByID(int id)",
		"_ => (DifficultyNormal, false)",
		"public static (int ID, bool Ok) DifficultyID(Difficulty diff)",
		"internal static DifficultyView DifficultyView(Difficulty difficulty)",
		"internal static Difficulty DifficultyFromView(DifficultyView view)",
		"FoodRegenerates = difficulty.FoodRegenerates() ? (byte)1 : (byte)0",
		"return new CapabilityDifficulty(",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated difficulty output missing %q:\n%s", expected, generated)
		}
	}
}

func TestInspectDifficultiesRejectsDrift(t *testing.T) {
	directory := difficultyDragonflyDirectory(t)
	path := filepath.Join(directory, "server", "world", "difficulty.go")
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	tests := map[string]struct {
		old, new string
		want     string
	}{
		"interface result": {
			old:  "FoodRegenerates() bool",
			new:  "FoodRegenerates() int",
			want: "Difficulty methods changed",
		},
		"lookup input": {
			old:  "func DifficultyByID(id int)",
			new:  "func DifficultyByID(id int64)",
			want: "DifficultyByID signature changed",
		},
		"registry id": {
			old:  "0: DifficultyPeaceful",
			new:  "4: DifficultyPeaceful",
			want: "live difficulty ID 4",
		},
		"builtin variable": {
			old:  "DifficultyHard difficultyHard",
			new:  "DifficultyHard difficultyNormal",
			want: "live DifficultyHard type",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			modified := strings.Replace(string(source), test.old, test.new, 1)
			if modified == string(source) {
				t.Fatalf("source does not contain %q", test.old)
			}
			file := filepath.Join(t.TempDir(), "difficulty.go")
			if err := os.WriteFile(file, []byte(modified), 0o600); err != nil {
				t.Fatal(err)
			}
			_, err := inspectDifficulties(file)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func difficultyDragonflyDirectory(t *testing.T) string {
	t.Helper()
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes.TrimSpace(output))
}
