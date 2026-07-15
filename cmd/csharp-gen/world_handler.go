package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

var selectedWorldHandlerMethods = []string{
	"HandleLiquidFlow",
	"HandleLiquidDecay",
	"HandleLiquidHarden",
	"HandleSound",
	"HandleFireSpread",
	"HandleBlockBurn",
	"HandleCropTrample",
	"HandleLeavesDecay",
	"HandleEntitySpawn",
	"HandleEntityDespawn",
	"HandleExplosion",
	"HandleRedstoneUpdate",
	"HandleClose",
}

const firstWorldHandlerSubscriptionBit = 41

type redstoneUpdateSpec struct {
	Causes []string
	Fields []parameter
}

func inspectWorldHandler(path string) ([]method, error) {
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
			if !ok || typeSpec.Name.Name != "Handler" {
				continue
			}
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return nil, fmt.Errorf("world.Handler is not an interface")
			}
			found := make(map[string]method, len(selectedWorldHandlerMethods))
			for _, field := range interfaceType.Methods.List {
				if len(field.Names) != 1 {
					return nil, fmt.Errorf("world.Handler has an unnamed or multiply named method")
				}
				name := field.Names[0].Name
				if !selectedWorldHandlerMethod(name) {
					return nil, fmt.Errorf("unsupported world.Handler.%s method", name)
				}
				function, ok := field.Type.(*ast.FuncType)
				if !ok {
					return nil, fmt.Errorf("world.Handler.%s is not a method", name)
				}
				parameters, err := translateWorldHandlerParameters(name, function.Params)
				if err != nil {
					return nil, fmt.Errorf("world.Handler.%s: %w", name, err)
				}
				found[name] = method{Name: name, Parameters: parameters}
			}
			methods := make([]method, 0, len(selectedWorldHandlerMethods))
			for index, name := range selectedWorldHandlerMethods {
				translated, ok := found[name]
				if !ok {
					return nil, fmt.Errorf("Dragonfly world.Handler has no supported %s method", name)
				}
				translated.Subscription = uint64(1) << uint(firstWorldHandlerSubscriptionBit+index)
				methods = append(methods, translated)
			}
			return methods, nil
		}
	}
	return nil, fmt.Errorf("Dragonfly world.Handler interface not found")
}

func inspectRedstoneUpdate(path string) (redstoneUpdateSpec, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return redstoneUpdateSpec{}, err
	}
	var result redstoneUpdateSpec
	for _, declaration := range file.Decls {
		gen, ok := declaration.(*ast.GenDecl)
		if !ok {
			continue
		}
		switch gen.Tok {
		case token.TYPE:
			for _, raw := range gen.Specs {
				typeSpec, ok := raw.(*ast.TypeSpec)
				if !ok || typeSpec.Name.Name != "RedstoneUpdate" {
					continue
				}
				structure, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return redstoneUpdateSpec{}, fmt.Errorf("world.RedstoneUpdate is not a struct")
				}
				for _, field := range structure.Fields.List {
					if len(field.Names) == 0 {
						return redstoneUpdateSpec{}, fmt.Errorf("world.RedstoneUpdate has an embedded field")
					}
					for _, name := range field.Names {
						typeName, ok := redstoneUpdateCSharpType(field.Type)
						if !ok {
							return redstoneUpdateSpec{}, fmt.Errorf("world.RedstoneUpdate.%s has unsupported type %s", name.Name, formatGoExpression(field.Type))
						}
						if name.Name == "After" && typeName == "World.Block" {
							typeName += "?"
						}
						result.Fields = append(result.Fields, parameter{Name: name.Name, Type: typeName})
					}
				}
			}
		case token.CONST:
			for _, raw := range gen.Specs {
				value, ok := raw.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range value.Names {
					const prefix = "RedstoneUpdateCause"
					if len(name.Name) <= len(prefix) || name.Name[:len(prefix)] != prefix {
						continue
					}
					result.Causes = append(result.Causes, name.Name[len(prefix):])
				}
			}
		}
	}
	if len(result.Causes) == 0 {
		return redstoneUpdateSpec{}, fmt.Errorf("Dragonfly world.RedstoneUpdateCause constants not found")
	}
	if len(result.Fields) == 0 {
		return redstoneUpdateSpec{}, fmt.Errorf("Dragonfly world.RedstoneUpdate struct not found")
	}
	return result, nil
}

func redstoneUpdateCSharpType(expression ast.Expr) (string, bool) {
	if identifier, ok := expression.(*ast.Ident); ok {
		typeName, found := map[string]string{
			"bool":                "bool",
			"int":                 "int",
			"int64":               "long",
			"Block":               "World.Block",
			"RedstoneUpdateCause": "World.RedstoneUpdateCause",
		}[identifier.Name]
		if found {
			return typeName, true
		}
	}
	return worldHandlerCSharpType(expression)
}

func selectedWorldHandlerMethod(name string) bool {
	for _, selected := range selectedWorldHandlerMethods {
		if name == selected {
			return true
		}
	}
	return false
}

func translateWorldHandlerParameters(methodName string, fields *ast.FieldList) ([]parameter, error) {
	var parameters []parameter
	for _, field := range fields.List {
		if len(field.Names) == 0 {
			return nil, fmt.Errorf("unnamed parameter")
		}
		for _, name := range field.Names {
			typeName, ok := worldHandlerCSharpType(field.Type)
			if !ok {
				return nil, fmt.Errorf("unsupported parameter type %s", formatGoExpression(field.Type))
			}
			if methodName == "HandleLiquidDecay" && name.Name == "after" && typeName == "World.Liquid" {
				typeName += "?"
			}
			parameters = append(parameters, parameter{Name: name.Name, Type: typeName})
		}
	}
	return parameters, nil
}

func worldHandlerCSharpType(expression ast.Expr) (string, bool) {
	switch value := expression.(type) {
	case *ast.ArrayType:
		if value.Len != nil {
			return "", false
		}
		element, ok := worldHandlerCSharpType(value.Elt)
		if !ok {
			return "", false
		}
		return element + "[]", true
	case *ast.StarExpr:
		typeName, ok := worldHandlerCSharpType(value.X)
		if !ok {
			return "", false
		}
		if typeName == "World.Context" || typeName == "World.Tx" {
			return typeName, true
		}
		return "ref " + typeName, true
	case *ast.Ident:
		typeName, ok := map[string]string{
			"bool":           "bool",
			"float64":        "double",
			"Block":          "World.Block",
			"Context":        "World.Context",
			"Entity":         "World.Entity",
			"Liquid":         "World.Liquid",
			"RedstoneUpdate": "World.RedstoneUpdate",
			"Sound":          "World.Sound",
			"Tx":             "World.Tx",
		}[value.Name]
		return typeName, ok
	case *ast.SelectorExpr:
		packageName, ok := value.X.(*ast.Ident)
		if !ok {
			return "", false
		}
		typeName, ok := map[string]string{
			"cube.Pos":   "Cube.Pos",
			"mgl64.Vec3": "Vector3",
		}[packageName.Name+"."+value.Sel.Name]
		return typeName, ok
	default:
		return "", false
	}
}

func generateWorldHandler(methods []method, redstone redstoneUpdateSpec) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/handler.go. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\n")
	output.WriteString("namespace Dragonfly;\n\n")
	output.WriteString("public sealed partial class World\n{\n")
	output.WriteString("    public enum RedstoneUpdateCause\n    {\n")
	for index, cause := range redstone.Causes {
		fmt.Fprintf(&output, "        %s = %d,\n", cause, index)
	}
	output.WriteString("    }\n\n")
	output.WriteString("    public readonly record struct RedstoneUpdate(\n")
	for index, field := range redstone.Fields {
		comma := ","
		if index == len(redstone.Fields)-1 {
			comma = ""
		}
		fmt.Fprintf(&output, "        %s %s%s\n", field.Type, field.Name, comma)
	}
	output.WriteString("    );\n\n")
	output.WriteString("    public interface Handler\n    {\n")
	for _, method := range methods {
		fmt.Fprintf(&output, "        void %s(%s);\n", method.Name, formatParameters(method.Parameters))
	}
	output.WriteString("    }\n}\n\n")
	output.WriteString("public abstract partial class Plugin : World.Handler\n{\n")
	for _, method := range methods {
		fmt.Fprintf(&output, "    [HandlerSubscription(%dUL)]\n", method.Subscription)
		fmt.Fprintf(&output, "    public virtual void %s(%s) { }\n", method.Name, formatParameters(method.Parameters))
	}
	output.WriteString("}\n")
	return output.Bytes()
}
