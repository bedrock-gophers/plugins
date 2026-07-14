// Command csharp-gen emits the supported C# surface directly from Dragonfly's Go AST.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type method struct {
	Name         string
	Parameters   []parameter
	Subscription uint64
}

type parameter struct {
	Name string
	Type string
}

type commandInterface struct {
	Name       string
	Embeddings []string
	Methods    []commandMethod
}

type commandMethod struct {
	Name       string
	Parameters []parameter
	ReturnType string
}

type generatedFile struct {
	Path    string
	Content []byte
}

var supportedPlayerHandlers = map[string]uint64{
	"HandleChat":         1 << 1,
	"HandleFoodLoss":     1 << 8,
	"HandleJump":         1 << 14,
	"HandleMove":         1,
	"HandlePunchAir":     1 << 17,
	"HandleQuit":         1 << 3,
	"HandleTeleport":     1 << 15,
	"HandleToggleSneak":  1 << 13,
	"HandleToggleSprint": 1 << 12,
}

var selectedCommandInterfaces = []string{
	"Runnable",
	"Allower",
	"Target",
	"NamedTarget",
	"Source",
	"Enum",
}

func main() {
	root := flag.String("root", ".", "repository root")
	dragonfly := flag.String("dragonfly", "", "Dragonfly module directory")
	check := flag.Bool("check", false, "fail if generated output differs")
	flag.Parse()

	directory := *dragonfly
	if directory == "" {
		command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
		command.Dir = *root
		output, err := command.Output()
		if err != nil {
			fatal(err)
		}
		directory = string(bytes.TrimSpace(output))
	}
	methods, err := playerHandlerMethods(filepath.Join(directory, "server", "player", "handler.go"))
	if err != nil {
		fatal(err)
	}
	interfaces, err := commandInterfaces(filepath.Join(directory, "server", "cmd"))
	if err != nil {
		fatal(err)
	}
	files := []generatedFile{
		{
			Path:    filepath.Join(*root, "csharp", "Dragonfly", "Generated", "Player.Handler.g.cs"),
			Content: generatePlayerHandler(methods),
		},
		{
			Path:    filepath.Join(*root, "csharp", "Dragonfly", "Generated", "Cmd.Interfaces.g.cs"),
			Content: generateCommandInterfaces(interfaces),
		},
	}
	if err := syncGeneratedFiles(files, *check); err != nil {
		fatal(err)
	}
}

func syncGeneratedFiles(files []generatedFile, check bool) error {
	for _, file := range files {
		if check {
			current, err := os.ReadFile(file.Path)
			if err != nil || !bytes.Equal(current, file.Content) {
				return fmt.Errorf("%s is stale; run make generate", file.Path)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(file.Path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(file.Path, file.Content, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func playerHandlerMethods(path string) ([]method, error) {
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
				return nil, fmt.Errorf("player.Handler is not an interface")
			}
			var methods []method
			for _, field := range interfaceType.Methods.List {
				if len(field.Names) != 1 {
					continue
				}
				subscription, supported := supportedPlayerHandlers[field.Names[0].Name]
				if !supported {
					continue
				}
				function, ok := field.Type.(*ast.FuncType)
				if !ok {
					return nil, fmt.Errorf("player.Handler.%s is not a method", field.Names[0].Name)
				}
				translated, err := translateParameters(function.Params)
				if err != nil {
					return nil, fmt.Errorf("player.Handler.%s: %w", field.Names[0].Name, err)
				}
				methods = append(methods, method{Name: field.Names[0].Name, Parameters: translated, Subscription: subscription})
			}
			for name := range supportedPlayerHandlers {
				found := false
				for _, method := range methods {
					found = found || method.Name == name
				}
				if !found {
					return nil, fmt.Errorf("Dragonfly player.Handler has no supported %s method", name)
				}
			}
			return methods, nil
		}
	}
	return nil, fmt.Errorf("Dragonfly player.Handler interface not found")
}

func translateParameters(fields *ast.FieldList) ([]parameter, error) {
	var parameters []parameter
	for _, field := range fields.List {
		typeName, ok := csharpType(field.Type)
		if !ok {
			return nil, fmt.Errorf("unsupported parameter type %T", field.Type)
		}
		for _, name := range field.Names {
			parameters = append(parameters, parameter{Name: name.Name, Type: typeName})
		}
	}
	return parameters, nil
}

func csharpType(expression ast.Expr) (string, bool) {
	switch value := expression.(type) {
	case *ast.StarExpr:
		typeName, ok := csharpType(value.X)
		if !ok {
			return "", false
		}
		if typeName == "string" || typeName == "int" || typeName == "double" || typeName == "bool" {
			return "ref " + typeName, true
		}
		return typeName, true
	case *ast.Ident:
		typeName, ok := map[string]string{
			"bool":    "bool",
			"Context": "Player.Context",
			"int":     "int",
			"Player":  "Player",
			"string":  "string",
		}[value.Name]
		return typeName, ok
	case *ast.SelectorExpr:
		packageName, ok := value.X.(*ast.Ident)
		if !ok {
			return "", false
		}
		typeName, ok := map[string]string{
			"mgl64.Vec3":    "Vector3",
			"cube.Rotation": "Rotation",
		}[packageName.Name+"."+value.Sel.Name]
		return typeName, ok
	default:
		return "", false
	}
}

func commandInterfaces(directory string) ([]commandInterface, error) {
	packages, err := parser.ParseDir(token.NewFileSet(), directory, nil, 0)
	if err != nil {
		return nil, err
	}
	pkg, ok := packages["cmd"]
	if !ok {
		return nil, fmt.Errorf("Dragonfly cmd package not found")
	}
	found := map[string]commandInterface{}
	for _, file := range pkg.Files {
		for _, declaration := range file.Decls {
			gen, ok := declaration.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gen.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || !selectedCommandInterface(typeSpec.Name.Name) {
					continue
				}
				interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
				if !ok {
					return nil, fmt.Errorf("cmd.%s is not an interface", typeSpec.Name.Name)
				}
				translated, err := translateCommandInterface(typeSpec.Name.Name, interfaceType)
				if err != nil {
					return nil, err
				}
				found[typeSpec.Name.Name] = translated
			}
		}
	}
	interfaces := make([]commandInterface, 0, len(selectedCommandInterfaces))
	for _, name := range selectedCommandInterfaces {
		definition, ok := found[name]
		if !ok {
			return nil, fmt.Errorf("Dragonfly cmd.%s interface not found", name)
		}
		interfaces = append(interfaces, definition)
	}
	return interfaces, nil
}

func selectedCommandInterface(name string) bool {
	for _, selected := range selectedCommandInterfaces {
		if name == selected {
			return true
		}
	}
	return false
}

func translateCommandInterface(name string, interfaceType *ast.InterfaceType) (commandInterface, error) {
	definition := commandInterface{Name: name}
	for _, field := range interfaceType.Methods.List {
		if len(field.Names) == 0 {
			embedding, ok := field.Type.(*ast.Ident)
			if !ok || !selectedCommandInterface(embedding.Name) {
				return commandInterface{}, fmt.Errorf("cmd.%s has unsupported embedded interface", name)
			}
			definition.Embeddings = append(definition.Embeddings, embedding.Name)
			continue
		}
		if len(field.Names) != 1 {
			return commandInterface{}, fmt.Errorf("cmd.%s has unnamed method", name)
		}
		function, ok := field.Type.(*ast.FuncType)
		if !ok {
			return commandInterface{}, fmt.Errorf("cmd.%s.%s is not a method", name, field.Names[0].Name)
		}
		parameters, err := translateCommandParameters(function.Params)
		if err != nil {
			return commandInterface{}, fmt.Errorf("cmd.%s.%s: %w", name, field.Names[0].Name, err)
		}
		returnType, err := translateCommandResult(function.Results)
		if err != nil {
			return commandInterface{}, fmt.Errorf("cmd.%s.%s: %w", name, field.Names[0].Name, err)
		}
		definition.Methods = append(definition.Methods, commandMethod{
			Name:       field.Names[0].Name,
			Parameters: parameters,
			ReturnType: returnType,
		})
	}
	return definition, nil
}

func translateCommandParameters(fields *ast.FieldList) ([]parameter, error) {
	var parameters []parameter
	if fields == nil {
		return parameters, nil
	}
	for _, field := range fields.List {
		typeName, ok := commandCSharpType(field.Type)
		if !ok {
			return nil, fmt.Errorf("unsupported parameter type %s", formatGoExpression(field.Type))
		}
		for _, name := range field.Names {
			parameters = append(parameters, parameter{Name: name.Name, Type: typeName})
		}
	}
	return parameters, nil
}

func translateCommandResult(fields *ast.FieldList) (string, error) {
	if fields == nil || len(fields.List) == 0 {
		return "void", nil
	}
	if len(fields.List) != 1 || len(fields.List[0].Names) > 1 {
		return "", fmt.Errorf("multiple return values are unsupported")
	}
	typeName, ok := commandCSharpType(fields.List[0].Type)
	if !ok {
		return "", fmt.Errorf("unsupported return type %s", formatGoExpression(fields.List[0].Type))
	}
	return typeName, nil
}

func commandCSharpType(expression ast.Expr) (string, bool) {
	switch value := expression.(type) {
	case *ast.StarExpr:
		typeName, ok := commandCSharpType(value.X)
		if !ok {
			return "", false
		}
		if typeName == "World.Tx" {
			return typeName + "?", true
		}
		return typeName, true
	case *ast.ArrayType:
		if value.Len != nil {
			return "", false
		}
		element, ok := commandCSharpType(value.Elt)
		if !ok {
			return "", false
		}
		return "IReadOnlyList<" + element + ">", true
	case *ast.Ident:
		typeName, ok := map[string]string{
			"bool":   "bool",
			"Output": "Output",
			"Source": "Source",
			"string": "string",
		}[value.Name]
		return typeName, ok
	case *ast.SelectorExpr:
		packageName, ok := value.X.(*ast.Ident)
		if !ok {
			return "", false
		}
		typeName, ok := map[string]string{
			"mgl64.Vec3": "Vector3",
			"world.Tx":   "World.Tx",
		}[packageName.Name+"."+value.Sel.Name]
		return typeName, ok
	default:
		return "", false
	}
}

func formatGoExpression(expression ast.Expr) string {
	switch value := expression.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.SelectorExpr:
		return formatGoExpression(value.X) + "." + value.Sel.Name
	case *ast.StarExpr:
		return "*" + formatGoExpression(value.X)
	case *ast.ArrayType:
		return "[]" + formatGoExpression(value.Elt)
	default:
		return fmt.Sprintf("%T", expression)
	}
}

func generatePlayerHandler(methods []method) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/player/handler.go. DO NOT EDIT.\n")
	output.WriteString("namespace Dragonfly;\n\n")
	output.WriteString("public sealed partial class Player\n{\n    public interface Handler\n    {\n")
	for _, method := range methods {
		fmt.Fprintf(&output, "        void %s(%s);\n", method.Name, formatParameters(method.Parameters))
	}
	output.WriteString("    }\n}\n\n")
	output.WriteString("public abstract partial class Plugin\n{\n")
	for _, method := range methods {
		fmt.Fprintf(&output, "    [HandlerSubscription(%dUL)]\n", method.Subscription)
		fmt.Fprintf(&output, "    public virtual void %s(%s) { }\n", method.Name, formatParameters(method.Parameters))
	}
	output.WriteString("}\n")
	return output.Bytes()
}

func generateCommandInterfaces(interfaces []commandInterface) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/cmd Go AST. DO NOT EDIT.\n#nullable enable\n")
	output.WriteString("namespace Dragonfly;\n\n")
	output.WriteString("public static partial class Cmd\n{\n")
	for index, definition := range interfaces {
		fmt.Fprintf(&output, "    public interface %s", definition.Name)
		if len(definition.Embeddings) != 0 {
			fmt.Fprintf(&output, " : %s", strings.Join(definition.Embeddings, ", "))
		}
		output.WriteString("\n    {\n")
		for _, method := range definition.Methods {
			fmt.Fprintf(&output, "        %s %s(%s);\n", method.ReturnType, method.Name, formatParameters(method.Parameters))
		}
		output.WriteString("    }\n")
		if index != len(interfaces)-1 {
			output.WriteString("\n")
		}
	}
	output.WriteString("}\n")
	return output.Bytes()
}

func formatParameters(parameters []parameter) string {
	values := make([]string, len(parameters))
	for index, parameter := range parameters {
		values[index] = parameter.Type + " " + parameter.Name
	}
	return strings.Join(values, ", ")
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
