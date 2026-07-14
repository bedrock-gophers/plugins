package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

var selectedPlayerIdentityMethods = []string{"Name", "UUID", "XUID"}

func inspectPlayerIdentityMethods(path string) ([]string, error) {
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
		"Name": {Results: "string"},
		"UUID": {Results: "uuid.UUID"},
		"XUID": {Results: "string"},
	}
	for _, name := range selectedPlayerIdentityMethods {
		function := found[name]
		if function == nil {
			return nil, fmt.Errorf("Dragonfly player.Player has no %s method", name)
		}
		if got := goFunctionSignature(function); got != want[name] {
			return nil, fmt.Errorf("Dragonfly player.Player.%s signature changed: %+v", name, got)
		}
	}
	return append([]string(nil), selectedPlayerIdentityMethods...), nil
}

func generatePlayerIdentityMethods(methods []string) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.\n")
	output.WriteString("using System;\n\nnamespace Dragonfly;\n\npublic sealed partial class Player\n{\n")
	for _, method := range methods {
		switch method {
		case "Name":
			output.WriteString("    public string Name() => PlayerName;\n")
		case "UUID":
			output.WriteString("    public Guid UUID() => PluginBridge.Host.PlayerUUID(Id);\n")
		case "XUID":
			output.WriteString("    public string XUID() => PluginBridge.Host.PlayerXUID(_invocation, Id);\n")
		default:
			panic("unsupported player identity method: " + method)
		}
	}
	output.WriteString("}\n")
	return output.Bytes()
}
