package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type playerSkinSpec struct {
	Parameter string
}

func inspectPlayerSkin(path string) (playerSkinSpec, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return playerSkinSpec{}, err
	}
	found := map[string]*ast.FuncDecl{}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && playerMethod(function) && (function.Name.Name == "Skin" || function.Name.Name == "SetSkin") {
			found[function.Name.Name] = function
		}
	}
	want := map[string]goSignature{
		"Skin":    {Results: "skin.Skin"},
		"SetSkin": {Parameters: "skin.Skin"},
	}
	for name, signature := range want {
		function := found[name]
		if function == nil {
			return playerSkinSpec{}, fmt.Errorf("Dragonfly player.Player has no %s method", name)
		}
		if got := goFunctionSignature(function); got != signature {
			return playerSkinSpec{}, fmt.Errorf("Dragonfly player.Player.%s signature changed: %+v", name, got)
		}
	}
	parameters := found["SetSkin"].Type.Params.List
	if len(parameters) != 1 || len(parameters[0].Names) != 1 {
		return playerSkinSpec{}, fmt.Errorf("Dragonfly player.Player.SetSkin parameter changed")
	}
	return playerSkinSpec{Parameter: parameters[0].Names[0].Name}, nil
}

func generatePlayerSkin(spec playerSkinSpec) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.\n")
	output.WriteString("namespace Dragonfly;\n\npublic sealed partial class Player\n{\n")
	output.WriteString("    public Skin Skin() => PluginBridge.Host.PlayerSkin(_invocation, Id);\n")
	fmt.Fprintf(&output, "    public void SetSkin(Skin %s) => PluginBridge.Host.SetPlayerSkin(_invocation, Id, %s);\n", spec.Parameter, spec.Parameter)
	output.WriteString("}\n")
	return output.Bytes()
}
