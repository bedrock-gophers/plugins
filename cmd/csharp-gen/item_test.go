package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInspectItemsUsesASTAndRegistry(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	spec, err := inspectItems(filepath.Join(string(bytes.TrimSpace(output)), "server", "item"))
	if err != nil {
		t.Fatal(err)
	}
	if spec.AirIdentifier != "minecraft:air" || len(spec.ToolTiers) != 7 || len(spec.Types) < 100 {
		t.Fatalf("item spec has air=%q tiers=%d types=%d", spec.AirIdentifier, len(spec.ToolTiers), len(spec.Types))
	}
	generated := string(generateItems(spec))
	for _, value := range []string{
		"public readonly record struct Sword(ToolTier Tier) : World.Item",
		"public readonly record struct Beef(bool Cooked) : World.Item",
		"Item.ToolTierDiamond",
		`identifier = "minecraft:diamond_sword"; metadata = 0`,
		`identifier = "minecraft:cooked_beef"; metadata = 0`,
		"private sealed record EncodedItem",
	} {
		if !strings.Contains(generated, value) {
			t.Fatalf("generated items missing %q", value)
		}
	}
}
