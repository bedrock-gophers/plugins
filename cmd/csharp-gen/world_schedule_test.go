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
import (
	"context"
	"time"
)
type World struct{}
type Tx struct{}
type Task struct{}
func (w *World) Do(f func(tx *Tx)) *Task { return nil }
func (w *World) DoAfter(delay time.Duration, f func(tx *Tx)) *Task { return nil }
func (t *Task) Done() <-chan struct{} { return nil }
func (t *Task) Err() error { return nil }
func (t *Task) Wait(ctx context.Context) error { return nil }
func (t *Task) OnDone(f func(err error)) {}
func (t *Task) Cancel() bool { return false }
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
		"#nullable enable",
		"public Task Do(Action<Tx> callback)",
		"PluginBridge.Host.ScheduleWorld(this, TimeSpan.Zero, callback);",
		"public Task DoAfter(TimeSpan delay, Action<Tx> callback)",
		"PluginBridge.Host.ScheduleWorld(this, delay, callback);",
		"public sealed class Task",
		"internal Task(ulong callbackId)",
		"internal void CallbackFailed(Exception error)",
		"internal void Complete(uint result)",
		"public System.Threading.Tasks.Task Done()",
		"public Exception? Err()",
		"public Exception? Wait(CancellationToken cancellationToken = default)",
		"public void OnDone(Action<Exception?> callback)",
		"public bool Cancel()",
		"PluginBridge.Host.CancelWorldTask(_callbackId)",
		"System.Threading.Tasks.Task.Run(() => callback(Err()))",
		"new TaskCanceledException(\"world task was cancelled\")",
		"new InvalidOperationException(\"world closed before task ran\")",
		"new InvalidOperationException(\"world task failed\")",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated world task surface missing %q:\n%s", expected, generated)
		}
	}
	if strings.Contains(generated, " Schedule(") {
		t.Fatalf("generated obsolete Schedule compatibility method:\n%s", generated)
	}
	if strings.Contains(generated, "unknown result") {
		t.Fatalf("generated ABI result detail leaked into public errors:\n%s", generated)
	}
}

func TestWorldScheduleRejectsSignatureDrift(t *testing.T) {
	tests := map[string][2]string{
		"do callback":    {"func(tx *Tx)", "func(tx Tx)"},
		"do result":      {"func (w *World) Do(f func(tx *Tx)) *Task", "func (w *World) Do(f func(tx *Tx)) Task"},
		"do after delay": {"delay time.Duration", "delay int64"},
		"done":           {"<-chan struct{}", "chan struct{}"},
		"err":            {"func (t *Task) Err() error", "func (t *Task) Err() bool"},
		"wait":           {"ctx context.Context", "ctx any"},
		"on done":        {"func(err error)", "func()"},
		"cancel":         {"Cancel() bool", "Cancel() error"},
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

func TestWorldScheduleRejectsMissingTaskStruct(t *testing.T) {
	path := filepath.Join(t.TempDir(), "task.go")
	source := strings.Replace(worldScheduleSource, "type Task struct{}", "type Task interface{}", 1)
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := inspectWorldSchedule(path); err == nil || !strings.Contains(err.Error(), "no Task struct") {
		t.Fatalf("expected missing Task struct error, got %v", err)
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
