package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

var selectedWorldLifecycleMethods = []string{
	"Name",
	"Range",
	"HighestLightBlocker",
	"Time",
	"SetTime",
	"Spawn",
	"SetSpawn",
	"Save",
	"Close",
}

type worldLifecycleMethod struct {
	Name       string
	Parameters []string
}

func inspectWorldLifecycleMethods(path string) ([]worldLifecycleMethod, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, err
	}
	found := map[string]*ast.FuncDecl{}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && receiverTypeName(function) == "World" {
			found[function.Name.Name] = function
		}
	}
	want := map[string]goSignature{
		"Name":                {Results: "string"},
		"Range":               {Results: "cube.Range"},
		"HighestLightBlocker": {Parameters: "int, int", Results: "int"},
		"Time":                {Results: "int"},
		"SetTime":             {Parameters: "int"},
		"Spawn":               {Results: "cube.Pos"},
		"SetSpawn":            {Parameters: "cube.Pos"},
		"Save":                {},
		"Close":               {Results: "error"},
	}
	wantParameterCount := map[string]int{
		"Name": 0, "Range": 0, "HighestLightBlocker": 2, "Time": 0, "SetTime": 1,
		"Spawn": 0, "SetSpawn": 1, "Save": 0, "Close": 0,
	}
	methods := make([]worldLifecycleMethod, 0, len(selectedWorldLifecycleMethods))
	for _, name := range selectedWorldLifecycleMethods {
		function := found[name]
		if function == nil {
			return nil, fmt.Errorf("Dragonfly world.World has no %s method", name)
		}
		if got := goFunctionSignature(function); got != want[name] {
			return nil, fmt.Errorf("Dragonfly world.World.%s signature changed: %+v", name, got)
		}
		method := worldLifecycleMethod{Name: name}
		if function.Type.Params != nil {
			for _, field := range function.Type.Params.List {
				for _, parameter := range field.Names {
					method.Parameters = append(method.Parameters, parameter.Name)
				}
			}
		}
		if len(method.Parameters) != wantParameterCount[name] {
			return nil, fmt.Errorf("Dragonfly world.World.%s parameter names changed", name)
		}
		methods = append(methods, method)
	}
	return methods, nil
}

func generateWorldLifecycleMethods(methods []worldLifecycleMethod) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/world.go Go AST. DO NOT EDIT.\n")
	output.WriteString("namespace Dragonfly;\n\npublic sealed partial class World\n{\n")
	for _, method := range methods {
		switch method.Name {
		case "Name":
			output.WriteString("    public string Name() => PluginBridge.Host.WorldName(_invocation, Id) ?? string.Empty;\n")
		case "Range":
			output.WriteString("    public Cube.Range Range() => PluginBridge.Host.WorldRange(_invocation, Id);\n")
		case "HighestLightBlocker":
			fmt.Fprintf(&output, "    public int HighestLightBlocker(int %s, int %s) =>\n        PluginBridge.Host.WorldHighestLightBlocker(_invocation, Id, %s, %s);\n", method.Parameters[0], method.Parameters[1], method.Parameters[0], method.Parameters[1])
		case "Time":
			output.WriteString("    public int Time() => PluginBridge.Host.WorldTime(_invocation, Id);\n")
		case "SetTime":
			parameter := csharpIdentifier(method.Parameters[0])
			fmt.Fprintf(&output, "    public void SetTime(int %s) => PluginBridge.Host.SetWorldTime(_invocation, Id, %s);\n", parameter, parameter)
		case "Spawn":
			output.WriteString("    public Cube.Pos Spawn() => PluginBridge.Host.WorldSpawn(_invocation, Id);\n")
		case "SetSpawn":
			fmt.Fprintf(&output, "    public void SetSpawn(Cube.Pos %s) =>\n        PluginBridge.Host.SetWorldSpawn(_invocation, Id, %s);\n", method.Parameters[0], method.Parameters[0])
		case "Save":
			output.WriteString("    public void Save() => PluginBridge.Host.SaveWorld(_invocation, Id);\n")
		case "Close":
			output.WriteString("    public void Close() => PluginBridge.Host.CloseWorld(_invocation, Id);\n")
		default:
			panic("unsupported world lifecycle method: " + method.Name)
		}
	}
	output.WriteString("}\n")
	return output.Bytes()
}
