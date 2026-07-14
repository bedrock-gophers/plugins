package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlayerItemMethodsFollowDragonflyAST(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	module, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	methods, err := inspectPlayerItemMethods(filepath.Join(string(bytes.TrimSpace(module)), "server", "player", "player.go"))
	if err != nil {
		t.Fatal(err)
	}
	generated := string(generatePlayerItemMethods(methods))
	for _, value := range []string{"Inventory.Value Inventory()", "Inventory.Armour Armour()", "HeldItems()", "SetHeldItems(", "SetHeldSlot("} {
		if !strings.Contains(generated, value) {
			t.Fatalf("generated player item methods missing %q", value)
		}
	}
}

func TestInventoriesFollowDragonflyAST(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	module, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	spec, err := inspectInventories(filepath.Join(string(bytes.TrimSpace(module)), "server", "item", "inventory"))
	if err != nil {
		t.Fatal(err)
	}
	generated := string(generateInventories(spec))
	for _, value := range []string{"int Size()", "Item.Stack Item(int slot)", "void SetItem", "int AddItem", "class Armour", "void SetHelmet", "Value Inventory()"} {
		if !strings.Contains(generated, value) {
			t.Fatalf("generated inventories missing %q", value)
		}
	}
	if strings.Contains(generated, "Clear()") {
		t.Fatal("lossy Clear method must remain absent until removed stacks can be returned")
	}
}
