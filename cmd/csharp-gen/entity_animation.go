package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
)

type entityAnimationSpec struct {
	ConstructorParameter string
	Methods              []entityAnimationMethod
}

type entityAnimationMethod struct {
	Name      string
	Parameter string
}

var entityAnimationFields = []string{"name", "nextState", "controller", "stopCondition"}
var entityAnimationMethods = []entityAnimationMethod{
	{Name: "Name"},
	{Name: "Controller"},
	{Name: "WithController", Parameter: "controller"},
	{Name: "NextState"},
	{Name: "WithNextState", Parameter: "state"},
	{Name: "StopCondition"},
	{Name: "WithStopCondition", Parameter: "condition"},
}

func inspectEntityAnimation(path string) (entityAnimationSpec, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return entityAnimationSpec{}, err
	}
	var declaration *ast.TypeSpec
	functions := map[string]*ast.FuncDecl{}
	for _, raw := range file.Decls {
		switch value := raw.(type) {
		case *ast.GenDecl:
			for _, rawSpec := range value.Specs {
				spec, ok := rawSpec.(*ast.TypeSpec)
				if ok && spec.Name.Name == "EntityAnimation" {
					declaration = spec
				}
			}
		case *ast.FuncDecl:
			if value.Name.Name == "NewEntityAnimation" || valueReceiver(value, "EntityAnimation") {
				functions[value.Name.Name] = value
			}
		}
	}
	if err := validateEntityAnimationStruct(declaration); err != nil {
		return entityAnimationSpec{}, err
	}
	constructor := functions["NewEntityAnimation"]
	if constructor == nil || goFunctionSignature(constructor) != (goSignature{Parameters: "string", Results: "EntityAnimation"}) ||
		!returnsEntityAnimationField(constructor, "name", "name") {
		return entityAnimationSpec{}, fmt.Errorf("Dragonfly world.NewEntityAnimation changed")
	}
	parameter, ok := soleParameterName(constructor)
	if !ok || parameter != "name" {
		return entityAnimationSpec{}, fmt.Errorf("Dragonfly world.NewEntityAnimation parameter changed")
	}
	for _, method := range entityAnimationMethods {
		function := functions[method.Name]
		if function == nil || !valueReceiver(function, "EntityAnimation") {
			return entityAnimationSpec{}, fmt.Errorf("Dragonfly world.EntityAnimation has no %s method", method.Name)
		}
		if method.Parameter == "" {
			if goFunctionSignature(function) != (goSignature{Results: "string"}) ||
				!returnsReceiverField(function, entityAnimationFieldForMethod(method.Name)) {
				return entityAnimationSpec{}, fmt.Errorf("Dragonfly world.EntityAnimation.%s changed", method.Name)
			}
			continue
		}
		if goFunctionSignature(function) != (goSignature{Parameters: "string", Results: "EntityAnimation"}) ||
			!setsReceiverFieldAndReturns(function, entityAnimationFieldForMethod(method.Name), method.Parameter) {
			return entityAnimationSpec{}, fmt.Errorf("Dragonfly world.EntityAnimation.%s changed", method.Name)
		}
		name, ok := soleParameterName(function)
		if !ok || name != method.Parameter {
			return entityAnimationSpec{}, fmt.Errorf("Dragonfly world.EntityAnimation.%s parameter changed", method.Name)
		}
	}
	return entityAnimationSpec{ConstructorParameter: parameter, Methods: append([]entityAnimationMethod(nil), entityAnimationMethods...)}, nil
}

func validateEntityAnimationStruct(spec *ast.TypeSpec) error {
	if spec == nil {
		return fmt.Errorf("Dragonfly world.EntityAnimation declaration not found")
	}
	structure, ok := spec.Type.(*ast.StructType)
	if !ok {
		return fmt.Errorf("Dragonfly world.EntityAnimation is not a struct")
	}
	var fields []string
	for _, field := range structure.Fields.List {
		typeName, ok := field.Type.(*ast.Ident)
		if !ok || typeName.Name != "string" {
			return fmt.Errorf("Dragonfly world.EntityAnimation contains non-string fields")
		}
		for _, name := range field.Names {
			fields = append(fields, name.Name)
		}
	}
	if !reflect.DeepEqual(fields, entityAnimationFields) {
		return fmt.Errorf("Dragonfly world.EntityAnimation fields changed: %v", fields)
	}
	return nil
}

func soleParameterName(function *ast.FuncDecl) (string, bool) {
	if function == nil || function.Type.Params == nil || len(function.Type.Params.List) != 1 || len(function.Type.Params.List[0].Names) != 1 {
		return "", false
	}
	return function.Type.Params.List[0].Names[0].Name, true
}

func returnsEntityAnimationField(function *ast.FuncDecl, field, parameter string) bool {
	if function == nil || function.Body == nil || len(function.Body.List) != 1 {
		return false
	}
	statement, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(statement.Results) != 1 {
		return false
	}
	literal, ok := statement.Results[0].(*ast.CompositeLit)
	if !ok || formatGoExpression(literal.Type) != "EntityAnimation" || len(literal.Elts) != 1 {
		return false
	}
	pair, ok := literal.Elts[0].(*ast.KeyValueExpr)
	return ok && formatGoExpression(pair.Key) == field && formatGoExpression(pair.Value) == parameter
}

func returnsReceiverField(function *ast.FuncDecl, field string) bool {
	if function.Body == nil || len(function.Body.List) != 1 {
		return false
	}
	statement, ok := function.Body.List[0].(*ast.ReturnStmt)
	return ok && len(statement.Results) == 1 && formatGoExpression(statement.Results[0]) == "a."+field
}

func setsReceiverFieldAndReturns(function *ast.FuncDecl, field, parameter string) bool {
	if function.Body == nil || len(function.Body.List) != 2 {
		return false
	}
	assignment, ok := function.Body.List[0].(*ast.AssignStmt)
	if !ok || assignment.Tok != token.ASSIGN || len(assignment.Lhs) != 1 || len(assignment.Rhs) != 1 ||
		formatGoExpression(assignment.Lhs[0]) != "a."+field || formatGoExpression(assignment.Rhs[0]) != parameter {
		return false
	}
	statement, ok := function.Body.List[1].(*ast.ReturnStmt)
	return ok && len(statement.Results) == 1 && formatGoExpression(statement.Results[0]) == "a"
}

func entityAnimationFieldForMethod(method string) string {
	switch method {
	case "Name":
		return "name"
	case "Controller", "WithController":
		return "controller"
	case "NextState", "WithNextState":
		return "nextState"
	case "StopCondition", "WithStopCondition":
		return "stopCondition"
	default:
		panic("unsupported entity animation method: " + method)
	}
}

func generateEntityAnimation(spec entityAnimationSpec) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/entity_animation.go Go AST. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\n\nnamespace Dragonfly;\n\n")
	output.WriteString("public sealed partial class World\n{\n")
	fmt.Fprintf(&output, "    public static EntityAnimation NewEntityAnimation(string %s) => new(%s);\n\n", spec.ConstructorParameter, spec.ConstructorParameter)
	output.WriteString("    public readonly struct EntityAnimation\n    {\n")
	output.WriteString("        private readonly string? _name;\n        private readonly string? _nextState;\n")
	output.WriteString("        private readonly string? _controller;\n        private readonly string? _stopCondition;\n\n")
	output.WriteString("        internal EntityAnimation(string name, string nextState = \"\", string controller = \"\", string stopCondition = \"\") =>\n")
	output.WriteString("            (_name, _nextState, _controller, _stopCondition) =\n")
	output.WriteString("                (name ?? throw new ArgumentNullException(nameof(name)), nextState, controller, stopCondition);\n")
	for _, method := range spec.Methods {
		output.WriteByte('\n')
		field := entityAnimationFieldForMethod(method.Name)
		if method.Parameter == "" {
			fmt.Fprintf(&output, "        public string %s() => _%s ?? string.Empty;\n", method.Name, field)
			continue
		}
		fmt.Fprintf(&output, "        public EntityAnimation %s(string %s) =>\n", method.Name, method.Parameter)
		value := method.Parameter + " ?? throw new ArgumentNullException(nameof(" + method.Parameter + "))"
		switch field {
		case "controller":
			fmt.Fprintf(&output, "            new(Name(), NextState(), %s, StopCondition());\n", value)
		case "nextState":
			fmt.Fprintf(&output, "            new(Name(), %s, Controller(), StopCondition());\n", value)
		case "stopCondition":
			fmt.Fprintf(&output, "            new(Name(), NextState(), Controller(), %s);\n", value)
		default:
			panic("unsupported entity animation setter field: " + field)
		}
	}
	output.WriteString("    }\n}\n")
	return output.Bytes()
}
