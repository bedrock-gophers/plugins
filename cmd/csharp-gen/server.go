package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

var selectedServerMethods = []string{"Players", "Player", "PlayerByName"}

func inspectServerMethods(path string) ([]commandMethod, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, err
	}
	found := map[string]commandMethod{}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || !pointerReceiver(function, "Server") || !selectedServerMethod(function.Name.Name) {
			continue
		}
		method, err := translateServerMethod(function)
		if err != nil {
			return nil, fmt.Errorf("server.Server.%s: %w", function.Name.Name, err)
		}
		found[function.Name.Name] = method
	}
	methods := make([]commandMethod, 0, len(selectedServerMethods))
	for _, name := range selectedServerMethods {
		method, ok := found[name]
		if !ok {
			return nil, fmt.Errorf("Dragonfly server.Server has no supported %s method", name)
		}
		methods = append(methods, method)
	}
	return methods, nil
}

func selectedServerMethod(name string) bool {
	for _, selected := range selectedServerMethods {
		if name == selected {
			return true
		}
	}
	return false
}

func translateServerMethod(function *ast.FuncDecl) (commandMethod, error) {
	method := commandMethod{Name: function.Name.Name}
	switch function.Name.Name {
	case "Players":
		if !serverPlayersParameter(function.Type.Params) || !serverPlayersResult(function.Type.Results) {
			return method, fmt.Errorf("signature changed")
		}
		method.Parameters = []parameter{{Name: "tx", Type: "World.Tx?"}}
		method.ReturnType = "IEnumerable<Player>"
	case "Player":
		if !singleSelectorParameter(function.Type.Params, "uuid", "uuid", "UUID") ||
			!serverPlayerLookupResults(function.Type.Results) {
			return method, fmt.Errorf("signature changed")
		}
		method.Parameters = []parameter{{Name: "uuid", Type: "Guid"}}
		method.ReturnType = "(World.EntityHandle? Player, bool Ok)"
	case "PlayerByName":
		if !singleIdentifierParameter(function.Type.Params, "name", "string") ||
			!serverPlayerLookupResults(function.Type.Results) {
			return method, fmt.Errorf("signature changed")
		}
		method.Parameters = []parameter{{Name: "name", Type: "string"}}
		method.ReturnType = "(World.EntityHandle? Player, bool Ok)"
	default:
		return method, fmt.Errorf("unsupported method")
	}
	return method, nil
}

func serverPlayersParameter(fields *ast.FieldList) bool {
	if fields == nil || len(fields.List) != 1 || len(fields.List[0].Names) != 1 || fields.List[0].Names[0].Name != "tx" {
		return false
	}
	pointer, ok := fields.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	selector, ok := pointer.X.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Tx" {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	return ok && pkg.Name == "world"
}

func serverPlayersResult(fields *ast.FieldList) bool {
	if fields == nil || len(fields.List) != 1 {
		return false
	}
	index, ok := fields.List[0].Type.(*ast.IndexExpr)
	if !ok {
		return false
	}
	selector, ok := index.X.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Seq" {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	if !ok || pkg.Name != "iter" {
		return false
	}
	pointer, ok := index.Index.(*ast.StarExpr)
	if !ok {
		return false
	}
	playerType, ok := pointer.X.(*ast.SelectorExpr)
	if !ok || playerType.Sel.Name != "Player" {
		return false
	}
	playerPkg, ok := playerType.X.(*ast.Ident)
	return ok && playerPkg.Name == "player"
}

func singleSelectorParameter(fields *ast.FieldList, name, packageName, typeName string) bool {
	if fields == nil || len(fields.List) != 1 || len(fields.List[0].Names) != 1 || fields.List[0].Names[0].Name != name {
		return false
	}
	selector, ok := fields.List[0].Type.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != typeName {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	return ok && pkg.Name == packageName
}

func singleIdentifierParameter(fields *ast.FieldList, name, typeName string) bool {
	return fields != nil && len(fields.List) == 1 && len(fields.List[0].Names) == 1 &&
		fields.List[0].Names[0].Name == name && singleResultType(fields.List[0].Type, typeName)
}

func serverPlayerLookupResults(fields *ast.FieldList) bool {
	if fields == nil || len(fields.List) != 2 || !singleResultType(fields.List[1].Type, "bool") {
		return false
	}
	pointer, ok := fields.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	selector, ok := pointer.X.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "EntityHandle" {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	return ok && pkg.Name == "world"
}

func generateServer(methods []commandMethod) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/server.go Go AST. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\nusing System;\nusing System.Collections.Generic;\n\nnamespace Dragonfly;\n\n")
	output.WriteString("public sealed partial class Server\n{\n    internal Server() { }\n\n")
	for index, method := range methods {
		switch method.Name {
		case "Players":
			output.WriteString("    public IEnumerable<Player> Players(World.Tx? tx = null) =>\n        PluginBridge.Host.ServerPlayers(tx?.Invocation ?? 0);\n")
		case "Player":
			output.WriteString("    public (World.EntityHandle? Player, bool Ok) Player(Guid uuid) =>\n        PluginBridge.Host.ServerPlayer(uuid);\n")
		case "PlayerByName":
			output.WriteString("    public (World.EntityHandle? Player, bool Ok) PlayerByName(string name) =>\n        PluginBridge.Host.ServerPlayerByName(name);\n")
		default:
			panic("unsupported server method: " + method.Name)
		}
		if index != len(methods)-1 {
			output.WriteByte('\n')
		}
	}
	output.WriteString("}\n\npublic abstract partial class Plugin\n{\n    public Server Server() => new();\n}\n")
	return output.Bytes()
}
