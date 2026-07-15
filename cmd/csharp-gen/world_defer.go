package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

var selectedWorldDeferMethods = []string{"Defer", "DeferErr"}

type worldDeferMethod struct {
	Name      string
	Parameter string
}

func inspectWorldDeferMethods(path string) ([]worldDeferMethod, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, err
	}
	found := map[string]*ast.FuncDecl{}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && pointerReceiver(function, "Tx") {
			found[function.Name.Name] = function
		}
	}
	want := map[string]goSignature{
		"Defer":    {Parameters: "func(*Tx)", Results: "*Task"},
		"DeferErr": {Parameters: "func(*Tx) error", Results: "*Task"},
	}
	methods := make([]worldDeferMethod, 0, len(selectedWorldDeferMethods))
	for _, name := range selectedWorldDeferMethods {
		function := found[name]
		if function == nil {
			return nil, fmt.Errorf("Dragonfly world.Tx has no %s method", name)
		}
		if got := goFunctionSignature(function); got != want[name] {
			return nil, fmt.Errorf("Dragonfly world.Tx.%s signature changed: %+v", name, got)
		}
		parameters := function.Type.Params.List
		if len(parameters) != 1 || len(parameters[0].Names) != 1 {
			return nil, fmt.Errorf("Dragonfly world.Tx.%s parameter changed", name)
		}
		methods = append(methods, worldDeferMethod{Name: name, Parameter: parameters[0].Names[0].Name})
	}
	return methods, nil
}

func generateWorldDefer(methods []worldDeferMethod) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/tx.go Go AST. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\nusing System;\nusing Dragonfly.Native;\n\nnamespace Dragonfly;\n\n")
	output.WriteString("public sealed partial class World\n{\n    public partial class Tx\n    {\n")
	for _, method := range methods {
		switch method.Name {
		case "Defer":
			fmt.Fprintf(&output, "        public World.Task Defer(Action<Tx> %[1]s) =>\n            PluginBridge.Host.DeferWorld(Invocation, %[1]s, Abi.WorldDeferDefer);\n", method.Parameter)
		case "DeferErr":
			fmt.Fprintf(&output, "        public World.Task DeferErr(Func<Tx, Exception?> %[1]s) =>\n            PluginBridge.Host.DeferWorld(Invocation, %[1]s, Abi.WorldDeferDeferErr);\n", method.Parameter)
		}
	}
	output.WriteString("    }\n}\n")
	return output.Bytes()
}

func generateNativeWorldDefer(methods []worldDeferMethod) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/tx.go Go AST. DO NOT EDIT.\npackage native\n\n")
	output.WriteString("type WorldDeferKind uint32\n\nconst (\n")
	for index, method := range methods {
		fmt.Fprintf(&output, "\t%-28s WorldDeferKind = %d\n", "WorldDefer"+method.Name, index)
	}
	output.WriteString(")\n")
	return output.Bytes()
}

func generateCSharpWorldDefer(methods []worldDeferMethod) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/tx.go Go AST. DO NOT EDIT.\nnamespace Dragonfly.Native;\n\npublic static partial class Abi\n{\n")
	for index, method := range methods {
		fmt.Fprintf(&output, "    public const uint WorldDefer%s = %d;\n", method.Name, index)
	}
	output.WriteString("}\n")
	return output.Bytes()
}

func generateFrameworkWorldDefer(methods []worldDeferMethod) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/tx.go Go AST. DO NOT EDIT.\npackage framework\n\n")
	output.WriteString("import (\n\t\"github.com/bedrock-gophers/plugins/internal/native\"\n\t\"github.com/df-mc/dragonfly/server/world\"\n)\n\n")
	output.WriteString("func runExactWorldDefer(tx *world.Tx, callback func(*world.Tx) error, kind native.WorldDeferKind) (*world.Task, bool) {\n\tswitch kind {\n")
	for _, method := range methods {
		switch method.Name {
		case "Defer":
			output.WriteString("\tcase native.WorldDeferDefer:\n\t\treturn tx.Defer(func(tx *world.Tx) {\n\t\t\tif err := callback(tx); err != nil {\n\t\t\t\tpanic(err)\n\t\t\t}\n\t\t}), true\n")
		case "DeferErr":
			output.WriteString("\tcase native.WorldDeferDeferErr:\n\t\treturn tx.DeferErr(callback), true\n")
		}
	}
	output.WriteString("\tdefault:\n\t\treturn nil, false\n\t}\n}\n")
	return output.Bytes()
}
