package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorldDeferMatchesPinnedDragonfly(t *testing.T) {
	directory, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly").Output()
	if err != nil {
		t.Fatal(err)
	}
	methods, err := inspectWorldDeferMethods(filepath.Join(string(bytes.TrimSpace(directory)), "server", "world", "tx.go"))
	if err != nil {
		t.Fatal(err)
	}
	csharp := string(generateWorldDefer(methods))
	for _, expected := range []string{
		"public World.Task Defer(Action<Tx> f)",
		"public World.Task DeferErr(Func<Tx, Exception?> f)",
	} {
		if !strings.Contains(csharp, expected) {
			t.Fatalf("generated world defer missing %q:\n%s", expected, csharp)
		}
	}
	framework := string(generateFrameworkWorldDefer(methods))
	for _, expected := range []string{"tx.Defer(func(tx *world.Tx)", "tx.DeferErr(callback)"} {
		if !strings.Contains(framework, expected) {
			t.Fatalf("generated exact dispatch missing %q:\n%s", expected, framework)
		}
	}
}
