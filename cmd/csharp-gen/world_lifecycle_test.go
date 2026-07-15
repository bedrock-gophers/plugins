package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPinnedDragonflyWorldLifecycleUsesGoAST(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	module, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	methods, err := inspectWorldLifecycleMethods(filepath.Join(
		string(bytes.TrimSpace(module)), "server", "world", "world.go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(methods) != 9 {
		t.Fatalf("generated %d world lifecycle methods, want 9", len(methods))
	}
	generated := string(generateWorldLifecycleMethods(methods))
	for _, expected := range []string{
		"string Name()",
		"Cube.Range Range()",
		"int HighestLightBlocker(int x, int z)",
		"int Time()",
		"void SetTime(int @new)",
		"Cube.Pos Spawn()",
		"void SetSpawn(Cube.Pos pos)",
		"void Save()",
		"void Close()",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated world lifecycle missing %q:\n%s", expected, generated)
		}
	}
}
