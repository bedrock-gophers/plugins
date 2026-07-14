package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type entityMethod struct {
	Name       string
	ReturnType string
}

func inspectWorldEntity(path string) ([]entityMethod, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, err
	}
	for _, declaration := range file.Decls {
		gen, ok := declaration.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != "Entity" {
				continue
			}
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return nil, fmt.Errorf("world.Entity is not an interface")
			}
			methods := make([]entityMethod, 0, 4)
			for _, field := range interfaceType.Methods.List {
				if len(field.Names) == 0 {
					selector, ok := field.Type.(*ast.SelectorExpr)
					if !ok {
						return nil, fmt.Errorf("world.Entity has unsupported embedded interface")
					}
					pkg, pkgOK := selector.X.(*ast.Ident)
					if !pkgOK || pkg.Name != "io" || selector.Sel.Name != "Closer" {
						return nil, fmt.Errorf("world.Entity has unsupported embedded interface")
					}
					methods = append(methods, entityMethod{Name: "Close", ReturnType: "void"})
					continue
				}
				if len(field.Names) != 1 {
					return nil, fmt.Errorf("world.Entity has multiply named method")
				}
				function, ok := field.Type.(*ast.FuncType)
				if !ok || function.Params.NumFields() != 0 || function.Results == nil || function.Results.NumFields() != 1 {
					return nil, fmt.Errorf("world.Entity.%s has unsupported signature", field.Names[0].Name)
				}
				name := field.Names[0].Name
				returnType, ok := entityReturnType(function.Results.List[0].Type)
				if !ok {
					return nil, fmt.Errorf("world.Entity.%s has unsupported return type", name)
				}
				methods = append(methods, entityMethod{Name: name, ReturnType: returnType})
			}
			if len(methods) != 4 {
				return nil, fmt.Errorf("world.Entity has %d methods, want 4", len(methods))
			}
			return methods, nil
		}
	}
	return nil, fmt.Errorf("world.Entity interface not found")
}

func entityReturnType(expression ast.Expr) (string, bool) {
	switch value := expression.(type) {
	case *ast.StarExpr:
		identifier, ok := value.X.(*ast.Ident)
		if !ok || identifier.Name != "EntityHandle" {
			return "", false
		}
		return "EntityHandle", true
	case *ast.SelectorExpr:
		pkg, ok := value.X.(*ast.Ident)
		if !ok {
			return "", false
		}
		translated, ok := map[string]string{
			"mgl64.Vec3":    "Vector3",
			"cube.Rotation": "Rotation",
		}[pkg.Name+"."+value.Sel.Name]
		return translated, ok
	default:
		return "", false
	}
}

func generateWorldEntity(methods []entityMethod) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/entity.go Go AST. DO NOT EDIT.\n")
	output.WriteString("namespace Dragonfly;\n\npublic sealed partial class World\n{\n    public interface Entity\n    {\n")
	for _, method := range methods {
		fmt.Fprintf(&output, "        %s %s();\n", method.ReturnType, method.Name)
	}
	output.WriteString("    }\n}\n")
	return output.Bytes()
}
