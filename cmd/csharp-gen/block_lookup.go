package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// inspectWorldBlockByName keeps the generated C# lookup pinned to Dragonfly's
// package-level world.BlockByName function instead of a hand-maintained schema.
func inspectWorldBlockByName(path string) error {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return err
	}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Recv != nil || function.Name.Name != "BlockByName" {
			continue
		}
		if function.Type.Params == nil || len(function.Type.Params.List) != 2 ||
			function.Type.Results == nil || len(function.Type.Results.List) != 2 {
			return fmt.Errorf("Dragonfly world.BlockByName shape changed")
		}
		name, properties := function.Type.Params.List[0], function.Type.Params.List[1]
		if len(name.Names) != 1 || name.Names[0].Name != "name" || formatGoExpression(name.Type) != "string" ||
			len(properties.Names) != 1 || properties.Names[0].Name != "properties" ||
			formatGoExpression(properties.Type) != "map[string]any" ||
			formatGoExpression(function.Type.Results.List[0].Type) != "Block" ||
			formatGoExpression(function.Type.Results.List[1].Type) != "bool" {
			return fmt.Errorf("Dragonfly world.BlockByName signature changed")
		}
		return nil
	}
	return fmt.Errorf("Dragonfly world.BlockByName is missing")
}
