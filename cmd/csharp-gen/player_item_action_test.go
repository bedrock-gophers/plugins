package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlayerItemActionsMatchPinnedDragonfly(t *testing.T) {
	directory, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly").Output()
	if err != nil {
		t.Fatal(err)
	}
	methods, err := inspectPlayerItemActionMethods(filepath.Join(
		string(bytes.TrimSpace(directory)), "server", "player", "player.go"))
	if err != nil {
		t.Fatal(err)
	}
	generated := string(generatePlayerItemActions(methods))
	for _, expected := range []string{
		"public (int Added, bool Ok) Collect(Item.Stack s)",
		"public int Drop(Item.Stack s)",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated player item actions missing %q:\n%s", expected, generated)
		}
	}
	nativeGenerated := string(generateNativePlayerItemActions(methods))
	csharpNative := string(generateCSharpPlayerItemActions(methods))
	hostGenerated := string(generateHostPlayerItemActions(methods))
	for _, name := range selectedPlayerItemActionMethods {
		if !strings.Contains(nativeGenerated, "PlayerItemAction"+name) ||
			!strings.Contains(csharpNative, "PlayerItemAction"+name+" = ") ||
			!strings.Contains(hostGenerated, "connected."+name+"(stack)") {
			t.Fatalf("generated transport missing %s", name)
		}
	}
}
