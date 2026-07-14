package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorldEntityUsesGoAST(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "entity.go")
	source := `package world
import (
    "io"
    "github.com/go-gl/mathgl/mgl64"
    "github.com/df-mc/dragonfly/server/block/cube"
)
type Entity interface {
    io.Closer
    H() *EntityHandle
    Position() mgl64.Vec3
    Rotation() cube.Rotation
}`
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	methods, err := inspectWorldEntity(path)
	if err != nil {
		t.Fatal(err)
	}
	generated := string(generateWorldEntity(methods))
	for _, expected := range []string{
		"void Close();",
		"EntityHandle H();",
		"Vector3 Position();",
		"Rotation Rotation();",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated entity interface missing %q:\n%s", expected, generated)
		}
	}
}

func TestPinnedDragonflyWorldEntityHasExactSurface(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	directory := string(bytes.TrimSpace(output))
	methods, err := inspectWorldEntity(filepath.Join(directory, "server", "world", "entity.go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(methods) != 4 {
		t.Fatalf("world.Entity has %d methods, want 4", len(methods))
	}
	want := []entityMethod{
		{Name: "Close", ReturnType: "void"},
		{Name: "H", ReturnType: "EntityHandle"},
		{Name: "Position", ReturnType: "Vector3"},
		{Name: "Rotation", ReturnType: "Rotation"},
	}
	for index := range want {
		if methods[index] != want[index] {
			t.Fatalf("world.Entity method %d = %+v, want %+v", index, methods[index], want[index])
		}
	}
}
