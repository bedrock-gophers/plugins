package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const worldScheduleSource = `package world
func (w *World) Do(f func(tx *Tx)) *Task { return nil }
`

func TestWorldScheduleUsesGoAST(t *testing.T) {
	path := filepath.Join(t.TempDir(), "task.go")
	if err := os.WriteFile(path, []byte(worldScheduleSource), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := inspectWorldSchedule(path); err != nil {
		t.Fatal(err)
	}
	generated := string(generateWorldSchedule())
	for _, expected := range []string{
		"using System;",
		"public void Schedule(Action<World.Tx> callback)",
		"PluginBridge.Host.ScheduleWorld(this, callback);",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated world schedule surface missing %q:\n%s", expected, generated)
		}
	}
}

func TestWorldScheduleRejectsSignatureDrift(t *testing.T) {
	tests := map[string][2]string{
		"callback": {"func(tx *Tx)", "func(tx Tx)"},
		"result":   {"*Task", "Task"},
	}
	for name, replacement := range tests {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "task.go")
			source := strings.Replace(worldScheduleSource, replacement[0], replacement[1], 1)
			if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := inspectWorldSchedule(path); err == nil || !strings.Contains(err.Error(), "signature changed") {
				t.Fatalf("expected signature drift error, got %v", err)
			}
		})
	}
}

func TestPinnedDragonflyWorldHasExactScheduleSurface(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	if err := inspectWorldSchedule(filepath.Join(
		string(bytes.TrimSpace(output)), "server", "world", "task.go")); err != nil {
		t.Fatal(err)
	}
}
