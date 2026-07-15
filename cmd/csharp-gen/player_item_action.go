package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

var selectedPlayerItemActionMethods = []string{
	"Collect",
	"Drop",
}

type playerItemActionMethod struct {
	Name      string
	Parameter string
}

func inspectPlayerItemActionMethods(path string) ([]playerItemActionMethod, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, err
	}
	found := map[string]*ast.FuncDecl{}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && playerMethod(function) {
			found[function.Name.Name] = function
		}
	}
	want := map[string]goSignature{
		"Collect": {Parameters: "item.Stack", Results: "int, bool"},
		"Drop":    {Parameters: "item.Stack", Results: "int"},
	}
	methods := make([]playerItemActionMethod, 0, len(selectedPlayerItemActionMethods))
	for _, name := range selectedPlayerItemActionMethods {
		function := found[name]
		if function == nil {
			return nil, fmt.Errorf("Dragonfly player.Player has no %s method", name)
		}
		if got := goFunctionSignature(function); got != want[name] {
			return nil, fmt.Errorf("Dragonfly player.Player.%s signature changed: %+v", name, got)
		}
		parameters := function.Type.Params.List
		if len(parameters) != 1 || len(parameters[0].Names) != 1 {
			return nil, fmt.Errorf("Dragonfly player.Player.%s parameter changed", name)
		}
		methods = append(methods, playerItemActionMethod{Name: name, Parameter: parameters[0].Names[0].Name})
	}
	return methods, nil
}

func generatePlayerItemActions(methods []playerItemActionMethod) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.\n")
	output.WriteString("using Dragonfly.Native;\n\nnamespace Dragonfly;\n\npublic sealed partial class Player\n{\n")
	for _, method := range methods {
		switch method.Name {
		case "Collect":
			fmt.Fprintf(&output, "    public (int Added, bool Ok) Collect(Item.Stack %[1]s) =>\n        PluginBridge.Host.RunPlayerItemAction(_invocation, Id, %[1]s, Abi.PlayerItemActionCollect);\n", method.Parameter)
		case "Drop":
			fmt.Fprintf(&output, "    public int Drop(Item.Stack %[1]s) =>\n        PluginBridge.Host.RunPlayerItemAction(_invocation, Id, %[1]s, Abi.PlayerItemActionDrop).Count;\n", method.Parameter)
		}
	}
	output.WriteString("}\n")
	return output.Bytes()
}

func generateNativePlayerItemActions(methods []playerItemActionMethod) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.\npackage native\n\n")
	output.WriteString("type PlayerItemActionKind uint32\n\nconst (\n")
	for index, method := range methods {
		fmt.Fprintf(&output, "\t%-40s PlayerItemActionKind = %d\n", "PlayerItemAction"+method.Name, index)
	}
	output.WriteString(")\n")
	return output.Bytes()
}

func generateCSharpPlayerItemActions(methods []playerItemActionMethod) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.\nnamespace Dragonfly.Native;\n\npublic static partial class Abi\n{\n")
	for index, method := range methods {
		fmt.Fprintf(&output, "    public const uint PlayerItemAction%s = %d;\n", method.Name, index)
	}
	output.WriteString("}\n")
	return output.Bytes()
}

func generateHostPlayerItemActions(methods []playerItemActionMethod) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.\npackage host\n\n")
	output.WriteString("import (\n\t\"github.com/bedrock-gophers/plugins/internal/native\"\n\t\"github.com/df-mc/dragonfly/server/item\"\n\t\"github.com/df-mc/dragonfly/server/player\"\n)\n\n")
	output.WriteString("func runExactPlayerItemAction(connected *player.Player, stack item.Stack, kind native.PlayerItemActionKind) (int, bool, bool) {\n\tswitch kind {\n")
	for _, method := range methods {
		switch method.Name {
		case "Collect":
			output.WriteString("\tcase native.PlayerItemActionCollect:\n\t\tcount, ok := connected.Collect(stack)\n\t\treturn count, ok, true\n")
		case "Drop":
			output.WriteString("\tcase native.PlayerItemActionDrop:\n\t\treturn connected.Drop(stack), true, true\n")
		}
	}
	output.WriteString("\tdefault:\n\t\treturn 0, false, false\n\t}\n}\n")
	return output.Bytes()
}
