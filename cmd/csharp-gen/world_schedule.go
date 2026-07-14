package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func inspectWorldSchedule(path string) error {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return err
	}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Name.Name != "Do" || !pointerReceiver(function, "World") {
			continue
		}
		want := goSignature{Parameters: "func(*Tx)", Results: "*Task"}
		if got := goFunctionSignature(function); got != want {
			return fmt.Errorf("Dragonfly world.World.Do signature changed: %+v", got)
		}
		return nil
	}
	return fmt.Errorf("Dragonfly world.World has no Do method")
}

func generateWorldSchedule() []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/task.go Go AST. DO NOT EDIT.\n")
	output.WriteString("using System;\n\nnamespace Dragonfly;\n\n")
	output.WriteString("public sealed partial class World\n{\n")
	output.WriteString("    public void Schedule(Action<World.Tx> callback) =>\n")
	output.WriteString("        PluginBridge.Host.ScheduleWorld(this, callback);\n")
	output.WriteString("}\n")
	return output.Bytes()
}
