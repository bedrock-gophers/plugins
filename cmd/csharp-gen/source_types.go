package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type sourceTypesSpec struct {
	Types []sourceTypeSpec
}

type sourceTypeSpec struct {
	Package             string
	Name                string
	Fields              []sourceFieldSpec
	Healing             bool
	ReducedByArmour     bool
	ReducedByResistance bool
	Fire                bool
	IgnoreTotem         bool
	AffectedEnchantment string
}

type sourceFieldSpec struct {
	Name       string
	GoType     string
	CSharpType string
}

type sourceTypeName struct {
	Package string
	Name    string
}

// This is an API ownership/order allowlist, not a schema. Every field, method,
// constant return and affected enchantment is read from Dragonfly's AST.
var sourceTypeNames = []sourceTypeName{
	{Package: "Entity", Name: "AttackDamageSource"},
	{Package: "Entity", Name: "VoidDamageSource"},
	{Package: "Entity", Name: "SuffocationDamageSource"},
	{Package: "Entity", Name: "DrowningDamageSource"},
	{Package: "Entity", Name: "FallDamageSource"},
	{Package: "Entity", Name: "GlideDamageSource"},
	{Package: "Entity", Name: "LightningDamageSource"},
	{Package: "Entity", Name: "ProjectileDamageSource"},
	{Package: "Entity", Name: "ExplosionDamageSource"},
	{Package: "Entity", Name: "FoodHealingSource"},
	{Package: "Effect", Name: "WitherDamageSource"},
	{Package: "Effect", Name: "InstantDamageSource"},
	{Package: "Effect", Name: "PoisonDamageSource"},
	{Package: "Effect", Name: "InstantHealingSource"},
	{Package: "Effect", Name: "RegenerationHealingSource"},
	{Package: "Player", Name: "StarvationDamageSource"},
	{Package: "Block", Name: "DamageSource"},
	{Package: "Block", Name: "MagmaDamageSource"},
	{Package: "Block", Name: "LavaDamageSource"},
	{Package: "Block", Name: "FireDamageSource"},
	{Package: "Enchantment", Name: "ThornsDamageSource"},
}

// inspectSourceTypes validates the complete exported Dragonfly damage/healing
// source surface. Directories must point at server/world, server/entity,
// server/entity/effect, server/player, server/block and server/item/enchantment.
func inspectSourceTypes(worldDir, entityDir, effectDir, playerDir, blockDir, enchantmentDir string) (sourceTypesSpec, error) {
	worldPackage, err := parseSourceTypePackage(worldDir, "world")
	if err != nil {
		return sourceTypesSpec{}, err
	}
	if err := validateSourceInterface(worldPackage, "DamageSource", nil, map[string]goSignature{
		"ReducedByArmour":     {Results: "bool"},
		"ReducedByResistance": {Results: "bool"},
		"Fire":                {Results: "bool"},
		"IgnoreTotem":         {Results: "bool"},
	}); err != nil {
		return sourceTypesSpec{}, err
	}
	if err := validateSourceInterface(worldPackage, "HealingSource", nil, map[string]goSignature{
		"HealingSource": {},
	}); err != nil {
		return sourceTypesSpec{}, err
	}

	packageDirectories := map[string]struct {
		GoName    string
		Directory string
	}{
		"Entity":      {GoName: "entity", Directory: entityDir},
		"Effect":      {GoName: "effect", Directory: effectDir},
		"Player":      {GoName: "player", Directory: playerDir},
		"Block":       {GoName: "block", Directory: blockDir},
		"Enchantment": {GoName: "enchantment", Directory: enchantmentDir},
	}
	packages := make(map[string]*ast.Package, len(packageDirectories))
	for name, input := range packageDirectories {
		pkg, parseErr := parseSourceTypePackage(input.Directory, input.GoName)
		if parseErr != nil {
			return sourceTypesSpec{}, parseErr
		}
		packages[name] = pkg
	}
	packages["World"] = worldPackage
	if err := validateSourceInterface(packages["Enchantment"], "AffectedDamageSource", []string{"world.DamageSource"}, map[string]goSignature{
		"AffectedByEnchantment": {Parameters: "item.EnchantmentType", Results: "bool"},
	}); err != nil {
		return sourceTypesSpec{}, err
	}

	actualSources := map[string]bool{}
	for packageName, pkg := range packages {
		for name, methods := range sourcePackageMethods(pkg) {
			if sourceHasDamageMethods(methods) || methods["HealingSource"] != nil {
				actualSources[packageName+"."+name] = true
			}
		}
	}
	expectedSources := make(map[string]bool, len(sourceTypeNames))
	for _, sourceType := range sourceTypeNames {
		expectedSources[sourceType.Package+"."+sourceType.Name] = true
	}
	var unknown []string
	for name := range actualSources {
		if !expectedSources[name] {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) != 0 {
		sort.Strings(unknown)
		return sourceTypesSpec{}, fmt.Errorf("unknown exported Dragonfly source types require API review: %s", strings.Join(unknown, ", "))
	}

	spec := sourceTypesSpec{Types: make([]sourceTypeSpec, 0, len(sourceTypeNames))}
	for _, sourceType := range sourceTypeNames {
		definition, inspectErr := inspectSourceType(packages[sourceType.Package], sourceType)
		if inspectErr != nil {
			return sourceTypesSpec{}, inspectErr
		}
		spec.Types = append(spec.Types, definition)
	}
	return spec, nil
}

func parseSourceTypePackage(directory, packageName string) (*ast.Package, error) {
	packages, err := parser.ParseDir(token.NewFileSet(), directory, func(info os.FileInfo) bool {
		return !info.IsDir() && filepath.Ext(info.Name()) == ".go" && !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		return nil, err
	}
	pkg := packages[packageName]
	if pkg == nil {
		return nil, fmt.Errorf("Dragonfly %s package not found", packageName)
	}
	return pkg, nil
}

func validateSourceInterface(pkg *ast.Package, name string, embeddings []string, methods map[string]goSignature) error {
	var found *ast.InterfaceType
	for _, file := range pkg.Files {
		for _, declaration := range file.Decls {
			general, ok := declaration.(*ast.GenDecl)
			if !ok || general.Tok != token.TYPE {
				continue
			}
			for _, raw := range general.Specs {
				typeSpec, ok := raw.(*ast.TypeSpec)
				if !ok || typeSpec.Name.Name != name {
					continue
				}
				found, _ = typeSpec.Type.(*ast.InterfaceType)
			}
		}
	}
	if found == nil {
		return fmt.Errorf("Dragonfly %s has no %s interface", pkg.Name, name)
	}
	actualEmbeddings := make([]string, 0)
	actualMethods := map[string]goSignature{}
	for _, field := range found.Methods.List {
		if len(field.Names) == 0 {
			actualEmbeddings = append(actualEmbeddings, formatGoExpression(field.Type))
			continue
		}
		function, ok := field.Type.(*ast.FuncType)
		if !ok || len(field.Names) != 1 {
			return fmt.Errorf("Dragonfly %s.%s interface shape changed", pkg.Name, name)
		}
		actualMethods[field.Names[0].Name] = goSignature{
			Parameters: formatFieldTypes(function.Params),
			Results:    formatFieldTypes(function.Results),
		}
	}
	if strings.Join(actualEmbeddings, ",") != strings.Join(embeddings, ",") || len(actualMethods) != len(methods) {
		return fmt.Errorf("Dragonfly %s.%s interface shape changed", pkg.Name, name)
	}
	for methodName, signature := range methods {
		if actual, ok := actualMethods[methodName]; !ok || actual != signature {
			return fmt.Errorf("Dragonfly %s.%s.%s signature changed: %+v", pkg.Name, name, methodName, actual)
		}
	}
	return nil
}

func sourcePackageMethods(pkg *ast.Package) map[string]map[string]*ast.FuncDecl {
	methods := map[string]map[string]*ast.FuncDecl{}
	for _, file := range pkg.Files {
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Recv == nil || !ast.IsExported(function.Name.Name) {
				continue
			}
			receiver := receiverTypeName(function)
			if receiver == "" || !ast.IsExported(receiver) {
				continue
			}
			if methods[receiver] == nil {
				methods[receiver] = map[string]*ast.FuncDecl{}
			}
			methods[receiver][function.Name.Name] = function
		}
	}
	return methods
}

func sourceHasDamageMethods(methods map[string]*ast.FuncDecl) bool {
	for _, name := range []string{"ReducedByArmour", "ReducedByResistance", "Fire", "IgnoreTotem"} {
		if methods[name] == nil {
			return false
		}
	}
	return true
}

func inspectSourceType(pkg *ast.Package, sourceType sourceTypeName) (sourceTypeSpec, error) {
	structure := sourceStruct(pkg, sourceType.Name)
	if structure == nil {
		return sourceTypeSpec{}, fmt.Errorf("Dragonfly %s.%s struct not found", pkg.Name, sourceType.Name)
	}
	fields, err := inspectSourceFields(structure)
	if err != nil {
		return sourceTypeSpec{}, fmt.Errorf("Dragonfly %s.%s %w", pkg.Name, sourceType.Name, err)
	}

	methods := sourcePackageMethods(pkg)[sourceType.Name]
	healing := methods["HealingSource"] != nil
	damage := sourceHasDamageMethods(methods)
	if healing == damage {
		return sourceTypeSpec{}, fmt.Errorf("Dragonfly %s.%s must implement exactly one source interface", pkg.Name, sourceType.Name)
	}
	if healing {
		if len(methods) != 1 {
			return sourceTypeSpec{}, fmt.Errorf("Dragonfly %s.%s exported methods changed", pkg.Name, sourceType.Name)
		}
		method := methods["HealingSource"]
		if !sourceValueReceiver(method, sourceType.Name) || goFunctionSignature(method) != (goSignature{}) || method.Body == nil || len(method.Body.List) != 0 {
			return sourceTypeSpec{}, fmt.Errorf("Dragonfly %s.%s.HealingSource changed", pkg.Name, sourceType.Name)
		}
	} else {
		wantedMethodCount := 4
		if methods["AffectedByEnchantment"] != nil {
			wantedMethodCount++
		}
		if len(methods) != wantedMethodCount {
			return sourceTypeSpec{}, fmt.Errorf("Dragonfly %s.%s exported methods changed", pkg.Name, sourceType.Name)
		}
		for _, name := range []string{"ReducedByArmour", "ReducedByResistance", "Fire", "IgnoreTotem"} {
			method := methods[name]
			if !sourceValueReceiver(method, sourceType.Name) || goFunctionSignature(method) != (goSignature{Results: "bool"}) {
				return sourceTypeSpec{}, fmt.Errorf("Dragonfly %s.%s.%s signature changed", pkg.Name, sourceType.Name, name)
			}
			if _, ok := sourceBooleanReturn(method); !ok {
				return sourceTypeSpec{}, fmt.Errorf("Dragonfly %s.%s.%s has unsupported non-constant behavior", pkg.Name, sourceType.Name, name)
			}
		}
		if methods["AffectedByEnchantment"] != nil {
			method := methods["AffectedByEnchantment"]
			if !sourceValueReceiver(method, sourceType.Name) || goFunctionSignature(method) != (goSignature{Parameters: "item.EnchantmentType", Results: "bool"}) {
				return sourceTypeSpec{}, fmt.Errorf("Dragonfly %s.%s.AffectedByEnchantment signature changed", pkg.Name, sourceType.Name)
			}
			if _, ok := sourceAffectedReturn(method); !ok {
				return sourceTypeSpec{}, fmt.Errorf("Dragonfly %s.%s.AffectedByEnchantment has unsupported behavior", pkg.Name, sourceType.Name)
			}
		}
	}
	values := map[string]bool{}
	for _, name := range []string{"ReducedByArmour", "ReducedByResistance", "Fire", "IgnoreTotem"} {
		if !healing {
			values[name], _ = sourceBooleanReturn(methods[name])
		}
	}
	affected := ""
	if method := methods["AffectedByEnchantment"]; method != nil {
		affected, _ = sourceAffectedReturn(method)
	}
	return sourceTypeSpec{
		Package: sourceType.Package, Name: sourceType.Name, Fields: fields, Healing: healing,
		ReducedByArmour: values["ReducedByArmour"], ReducedByResistance: values["ReducedByResistance"],
		Fire: values["Fire"], IgnoreTotem: values["IgnoreTotem"], AffectedEnchantment: affected,
	}, nil
}

func sourceStruct(pkg *ast.Package, name string) *ast.StructType {
	for _, file := range pkg.Files {
		for _, declaration := range file.Decls {
			general, ok := declaration.(*ast.GenDecl)
			if !ok || general.Tok != token.TYPE {
				continue
			}
			for _, raw := range general.Specs {
				typeSpec, ok := raw.(*ast.TypeSpec)
				if ok && typeSpec.Name.Name == name {
					structure, _ := typeSpec.Type.(*ast.StructType)
					return structure
				}
			}
		}
	}
	return nil
}

func inspectSourceFields(structure *ast.StructType) ([]sourceFieldSpec, error) {
	var fields []sourceFieldSpec
	for _, field := range structure.Fields.List {
		if len(field.Names) == 0 {
			return nil, fmt.Errorf("has embedded field %s", formatGoExpression(field.Type))
		}
		goType := formatGoExpression(field.Type)
		csharpType, ok := map[string]string{"bool": "bool", "world.Entity": "World.Entity?", "world.Block": "World.Block?"}[goType]
		if !ok {
			return nil, fmt.Errorf("has unsupported field type %s", goType)
		}
		for _, name := range field.Names {
			if !ast.IsExported(name.Name) {
				return nil, fmt.Errorf("has unexported field %s", name.Name)
			}
			fields = append(fields, sourceFieldSpec{Name: name.Name, GoType: goType, CSharpType: csharpType})
		}
	}
	return fields, nil
}

func sourceValueReceiver(function *ast.FuncDecl, name string) bool {
	if function.Recv == nil || len(function.Recv.List) != 1 {
		return false
	}
	receiver, ok := function.Recv.List[0].Type.(*ast.Ident)
	return ok && receiver.Name == name
}

func sourceBooleanReturn(function *ast.FuncDecl) (bool, bool) {
	if function.Body == nil || len(function.Body.List) != 1 {
		return false, false
	}
	statement, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(statement.Results) != 1 {
		return false, false
	}
	identifier, ok := statement.Results[0].(*ast.Ident)
	if !ok || (identifier.Name != "true" && identifier.Name != "false") {
		return false, false
	}
	return identifier.Name == "true", true
}

func sourceAffectedReturn(function *ast.FuncDecl) (string, bool) {
	if function.Body == nil || len(function.Body.List) != 1 || function.Type.Params == nil || len(function.Type.Params.List) != 1 || len(function.Type.Params.List[0].Names) != 1 {
		return "", false
	}
	statement, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(statement.Results) != 1 {
		return "", false
	}
	comparison, ok := statement.Results[0].(*ast.BinaryExpr)
	if !ok || comparison.Op != token.EQL {
		return "", false
	}
	parameter, ok := comparison.X.(*ast.Ident)
	selector, selectorOK := comparison.Y.(*ast.SelectorExpr)
	if !selectorOK {
		return "", false
	}
	packageName, packageOK := selector.X.(*ast.Ident)
	if !ok || !packageOK || parameter.Name != function.Type.Params.List[0].Names[0].Name || packageName.Name != "enchantment" || !ast.IsExported(selector.Sel.Name) {
		return "", false
	}
	return selector.Sel.Name, true
}

func generateSourceTypes(spec sourceTypesSpec) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly damage/healing source Go AST. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\n\nnamespace Dragonfly;\n\n")
	output.WriteString(`public sealed partial class World
{
    public interface DamageSource
    {
        bool ReducedByArmour();
        bool ReducedByResistance();
        bool Fire();
        bool IgnoreTotem();
    }

    // Go's HealingSource() marker cannot be named on this C# interface: C# reserves
    // a member matching its enclosing type name for constructors.
    public interface HealingSource { }
}
`)

	packageOrder := []string{"Entity", "Effect", "Player", "Block", "Enchantment"}
	for _, packageName := range packageOrder {
		if packageName == "Player" {
			output.WriteString("\npublic sealed partial class Player\n{\n")
		} else {
			fmt.Fprintf(&output, "\npublic static partial class %s\n{\n", packageName)
		}
		if packageName == "Enchantment" {
			output.WriteString("    public interface AffectedDamageSource : World.DamageSource\n    {\n        bool AffectedByEnchantment(Item.EnchantmentType e);\n    }\n\n")
		}
		first := true
		for _, definition := range spec.Types {
			if definition.Package != packageName {
				continue
			}
			if !first {
				output.WriteByte('\n')
			}
			first = false
			writeSourceType(&output, definition)
		}
		output.WriteString("}\n")
	}
	return output.Bytes()
}

func writeSourceType(output *bytes.Buffer, definition sourceTypeSpec) {
	interfaceName := "World.DamageSource"
	if definition.Healing {
		interfaceName = "World.HealingSource"
	} else if definition.AffectedEnchantment != "" {
		interfaceName = "Enchantment.AffectedDamageSource"
	}
	fmt.Fprintf(output, "    public readonly record struct %s", definition.Name)
	if len(definition.Fields) != 0 {
		output.WriteByte('(')
		for index, field := range definition.Fields {
			if index != 0 {
				output.WriteString(", ")
			}
			defaultValue := "false"
			if strings.HasSuffix(field.CSharpType, "?") {
				defaultValue = "null"
			}
			fmt.Fprintf(output, "%s %s = %s", field.CSharpType, field.Name, defaultValue)
		}
		output.WriteByte(')')
	}
	fmt.Fprintf(output, " : %s", interfaceName)
	if definition.Healing {
		output.WriteString(";\n")
		return
	}
	output.WriteString("\n    {\n")
	fmt.Fprintf(output, "        public bool ReducedByArmour() => %s;\n", strconv.FormatBool(definition.ReducedByArmour))
	fmt.Fprintf(output, "        public bool ReducedByResistance() => %s;\n", strconv.FormatBool(definition.ReducedByResistance))
	fmt.Fprintf(output, "        public bool Fire() => %s;\n", strconv.FormatBool(definition.Fire))
	fmt.Fprintf(output, "        public bool IgnoreTotem() => %s;\n", strconv.FormatBool(definition.IgnoreTotem))
	if definition.AffectedEnchantment != "" {
		fmt.Fprintf(output, "        public bool AffectedByEnchantment(Item.EnchantmentType e) => object.Equals(e, Item.%s);\n", definition.AffectedEnchantment)
	}
	output.WriteString("    }\n")
}
