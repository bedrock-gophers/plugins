package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

type inventorySpec struct {
	Inventory []string
	Armour    []string
}

var selectedPlayerItemMethods = []string{
	"Inventory",
	"EnderChestInventory",
	"Armour",
	"HeldItems",
	"SetHeldItems",
	"SetHeldSlot",
}

func inspectPlayerItemMethods(path string) ([]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, err
	}
	found := map[string]*ast.FuncDecl{}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || !playerMethod(function) {
			continue
		}
		found[function.Name.Name] = function
	}
	want := map[string]goSignature{
		"Inventory":           {Results: "*inventory.Inventory"},
		"EnderChestInventory": {Results: "*inventory.Inventory"},
		"Armour":              {Results: "*inventory.Armour"},
		"HeldItems":           {Results: "item.Stack, item.Stack"},
		"SetHeldItems":        {Parameters: "item.Stack, item.Stack"},
		"SetHeldSlot":         {Parameters: "int", Results: "error"},
	}
	for _, name := range selectedPlayerItemMethods {
		function := found[name]
		if function == nil {
			return nil, fmt.Errorf("Dragonfly player.Player has no %s method", name)
		}
		signature := goFunctionSignature(function)
		if signature != want[name] {
			return nil, fmt.Errorf("Dragonfly player.Player.%s signature changed: %+v", name, signature)
		}
	}
	return append([]string(nil), selectedPlayerItemMethods...), nil
}

func goFunctionSignature(function *ast.FuncDecl) goSignature {
	return goSignature{
		Parameters: formatFieldTypes(function.Type.Params),
		Results:    formatFieldTypes(function.Type.Results),
	}
}

func formatFieldTypes(fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}
	var result []string
	for _, field := range fields.List {
		count := len(field.Names)
		if count == 0 {
			count = 1
		}
		for range count {
			result = append(result, formatGoExpression(field.Type))
		}
	}
	return joinComma(result)
}

func joinComma(values []string) string {
	var output bytes.Buffer
	for index, value := range values {
		if index != 0 {
			output.WriteString(", ")
		}
		output.WriteString(value)
	}
	return output.String()
}

func generatePlayerItemMethods(methods []string) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\nusing Dragonfly.Native;\n\nnamespace Dragonfly;\n\n")
	output.WriteString("public sealed partial class Player\n{\n")
	for _, method := range methods {
		switch method {
		case "Inventory":
			output.WriteString("    public Inventory.Value Inventory() => new(_invocation, Id, Abi.InventoryMain);\n")
		case "EnderChestInventory":
			output.WriteString("    public Inventory.Value EnderChestInventory() => new(_invocation, Id, Abi.InventoryEnderChest);\n")
		case "Armour":
			output.WriteString("    public Inventory.Armour Armour() => new(_invocation, Id);\n")
		case "HeldItems":
			output.WriteString("    public (Item.Stack MainHand, Item.Stack OffHand) HeldItems() =>\n        PluginBridge.Host.HeldItems(_invocation, Id);\n")
		case "SetHeldItems":
			output.WriteString("    public void SetHeldItems(Item.Stack mainHand, Item.Stack offHand) =>\n        PluginBridge.Host.SetHeldItems(_invocation, Id, mainHand, offHand);\n")
		case "SetHeldSlot":
			output.WriteString("    public void SetHeldSlot(int to) => PluginBridge.Host.SetHeldSlot(_invocation, Id, to);\n")
		default:
			panic("unsupported player item method: " + method)
		}
	}
	output.WriteString("}\n")
	return output.Bytes()
}

func inspectInventories(directory string) (inventorySpec, error) {
	inventoryMethods := map[string]goSignature{
		"Size":    {Results: "int"},
		"Item":    {Parameters: "int", Results: "item.Stack, error"},
		"SetItem": {Parameters: "int, item.Stack", Results: "error"},
		"AddItem": {Parameters: "item.Stack", Results: "int, error"},
	}
	armourMethods := map[string]goSignature{
		"Helmet":        {Results: "item.Stack"},
		"Chestplate":    {Results: "item.Stack"},
		"Leggings":      {Results: "item.Stack"},
		"Boots":         {Results: "item.Stack"},
		"SetHelmet":     {Parameters: "item.Stack"},
		"SetChestplate": {Parameters: "item.Stack"},
		"SetLeggings":   {Parameters: "item.Stack"},
		"SetBoots":      {Parameters: "item.Stack"},
		"Set":           {Parameters: "item.Stack, item.Stack, item.Stack, item.Stack"},
		"Inventory":     {Results: "*Inventory"},
	}
	found := map[string]map[string]*ast.FuncDecl{"Inventory": {}, "Armour": {}}
	for _, name := range []string{"inventory.go", "armour.go"} {
		file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(directory, name), nil, 0)
		if err != nil {
			return inventorySpec{}, err
		}
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok {
				continue
			}
			receiver := receiverTypeName(function)
			if methods := found[receiver]; methods != nil {
				methods[function.Name.Name] = function
			}
		}
	}
	validate := func(receiver string, methods map[string]goSignature) ([]string, error) {
		names := make([]string, 0, len(methods))
		for name, want := range methods {
			function := found[receiver][name]
			if function == nil {
				return nil, fmt.Errorf("Dragonfly inventory.%s has no %s method", receiver, name)
			}
			if got := goFunctionSignature(function); got != want {
				return nil, fmt.Errorf("Dragonfly inventory.%s.%s signature changed: %+v", receiver, name, got)
			}
			names = append(names, name)
		}
		return names, nil
	}
	inventory, err := validate("Inventory", inventoryMethods)
	if err != nil {
		return inventorySpec{}, err
	}
	armour, err := validate("Armour", armourMethods)
	if err != nil {
		return inventorySpec{}, err
	}
	return inventorySpec{Inventory: inventory, Armour: armour}, nil
}

func receiverTypeName(function *ast.FuncDecl) string {
	if function.Recv == nil || len(function.Recv.List) != 1 {
		return ""
	}
	typeExpression := function.Recv.List[0].Type
	if pointer, ok := typeExpression.(*ast.StarExpr); ok {
		typeExpression = pointer.X
	}
	identifier, _ := typeExpression.(*ast.Ident)
	if identifier == nil {
		return ""
	}
	return identifier.Name
}

func generateInventories(spec inventorySpec) []byte {
	requireInventoryMethod := func(name string) {
		for _, method := range spec.Inventory {
			if method == name {
				return
			}
		}
		panic("missing inventory method: " + name)
	}
	requireArmourMethod := func(name string) {
		for _, method := range spec.Armour {
			if method == name {
				return
			}
		}
		panic("missing armour method: " + name)
	}
	for _, name := range []string{"Size", "Item", "SetItem", "AddItem"} {
		requireInventoryMethod(name)
	}
	for _, name := range []string{"Helmet", "Chestplate", "Leggings", "Boots", "SetHelmet", "SetChestplate", "SetLeggings", "SetBoots", "Set", "Inventory"} {
		requireArmourMethod(name)
	}
	return []byte(`// Code generated from Dragonfly server/item/inventory Go AST. DO NOT EDIT.
#nullable enable
using Dragonfly.Native;

namespace Dragonfly;

public static class Inventory
{
    public sealed class Value
    {
        private readonly ulong _invocation;
        private readonly InventoryId _id;

        internal Value(ulong invocation, PlayerId player, uint kind)
        {
            _invocation = invocation;
            _id = new InventoryId { Player = player, Kind = kind };
        }

        public int Size() => PluginBridge.Host.InventorySize(_invocation, _id);

        public Item.Stack Item(int slot)
        {
            CheckSlot(slot);
            return PluginBridge.Host.InventoryItem(_invocation, _id, slot);
        }

        public void SetItem(int slot, Item.Stack item)
        {
            CheckSlot(slot);
            PluginBridge.Host.SetInventoryItem(_invocation, _id, slot, item);
        }

        public int AddItem(Item.Stack item) =>
            PluginBridge.Host.AddInventoryItem(_invocation, _id, item);

        private void CheckSlot(int slot)
        {
            if (slot < 0 || slot >= Size()) throw new ArgumentOutOfRangeException(nameof(slot));
        }
    }

    public sealed class Armour
    {
        private readonly Value _inventory;

        internal Armour(ulong invocation, PlayerId player) =>
            _inventory = new Value(invocation, player, Abi.InventoryArmour);

        public Item.Stack Helmet() => _inventory.Item(0);
        public Item.Stack Chestplate() => _inventory.Item(1);
        public Item.Stack Leggings() => _inventory.Item(2);
        public Item.Stack Boots() => _inventory.Item(3);

        public void SetHelmet(Item.Stack helmet) => _inventory.SetItem(0, helmet);
        public void SetChestplate(Item.Stack chestplate) => _inventory.SetItem(1, chestplate);
        public void SetLeggings(Item.Stack leggings) => _inventory.SetItem(2, leggings);
        public void SetBoots(Item.Stack boots) => _inventory.SetItem(3, boots);

        public void Set(Item.Stack helmet, Item.Stack chestplate, Item.Stack leggings, Item.Stack boots)
        {
            SetHelmet(helmet);
            SetChestplate(chestplate);
            SetLeggings(leggings);
            SetBoots(boots);
        }

        public Value Inventory() => _inventory;
    }
}
`)
}
