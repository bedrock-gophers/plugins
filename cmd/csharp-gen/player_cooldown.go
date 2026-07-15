package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type playerCooldownSpec struct {
	Item, Duration string
}

func inspectPlayerCooldown(path string) (playerCooldownSpec, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return playerCooldownSpec{}, err
	}
	found := map[string]*ast.FuncDecl{}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && playerMethod(function) && (function.Name.Name == "HasCooldown" || function.Name.Name == "SetCooldown") {
			found[function.Name.Name] = function
		}
	}
	want := map[string]goSignature{
		"HasCooldown": {Parameters: "world.Item", Results: "bool"},
		"SetCooldown": {Parameters: "world.Item, time.Duration"},
	}
	for name, signature := range want {
		function := found[name]
		if function == nil {
			return playerCooldownSpec{}, fmt.Errorf("Dragonfly player.Player has no %s method", name)
		}
		if got := goFunctionSignature(function); got != signature {
			return playerCooldownSpec{}, fmt.Errorf("Dragonfly player.Player.%s signature changed: %+v", name, got)
		}
	}
	has := found["HasCooldown"].Type.Params.List
	set := found["SetCooldown"].Type.Params.List
	if len(has) != 1 || len(has[0].Names) != 1 || len(set) != 2 || len(set[0].Names) != 1 || len(set[1].Names) != 1 ||
		has[0].Names[0].Name != set[0].Names[0].Name {
		return playerCooldownSpec{}, fmt.Errorf("Dragonfly player cooldown parameter names changed")
	}
	return playerCooldownSpec{Item: has[0].Names[0].Name, Duration: set[1].Names[0].Name}, nil
}

func generatePlayerCooldown(spec playerCooldownSpec) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\nusing System;\n\nnamespace Dragonfly;\n\n")
	output.WriteString("public sealed partial class Player\n{\n")
	fmt.Fprintf(&output, "    public bool HasCooldown(World.Item? %s) => PluginBridge.Host.HasPlayerCooldown(_invocation, Id, %s);\n", spec.Item, spec.Item)
	fmt.Fprintf(&output, "    public void SetCooldown(World.Item? %s, TimeSpan %s) => PluginBridge.Host.SetPlayerCooldown(_invocation, Id, %s, %s);\n", spec.Item, spec.Duration, spec.Item, spec.Duration)
	output.WriteString("}\n")
	return output.Bytes()
}
