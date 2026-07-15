package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"reflect"
	"sort"
	"strconv"

	"github.com/df-mc/dragonfly/server/world"
)

var difficultyMethods = []difficultyMethod{
	{Name: "FoodRegenerates", GoResult: "bool", CSharpResult: "bool"},
	{Name: "StarvationHealthLimit", GoResult: "float64", CSharpResult: "double"},
	{Name: "FireSpreadIncrease", GoResult: "int", CSharpResult: "int"},
}

var difficultyNames = []string{
	"DifficultyPeaceful",
	"DifficultyEasy",
	"DifficultyNormal",
	"DifficultyHard",
}

type difficultyMethod struct {
	Name         string
	GoResult     string
	CSharpResult string
}

type difficultyValue struct {
	Name                  string
	PrivateType           string
	ID                    int
	FoodRegenerates       bool
	StarvationHealthLimit float64
	FireSpreadIncrease    int
}

type difficultySpec struct {
	Methods []difficultyMethod
	Values  []difficultyValue
}

func inspectDifficulties(path string) (difficultySpec, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return difficultySpec{}, err
	}
	types := map[string]*ast.TypeSpec{}
	variables := map[string]string{}
	functions := map[string]*ast.FuncDecl{}
	var registry *ast.CompositeLit
	for _, declaration := range file.Decls {
		switch value := declaration.(type) {
		case *ast.GenDecl:
			for _, raw := range value.Specs {
				switch spec := raw.(type) {
				case *ast.TypeSpec:
					types[spec.Name.Name] = spec
				case *ast.ValueSpec:
					if identifier, ok := spec.Type.(*ast.Ident); ok {
						for _, name := range spec.Names {
							if name.IsExported() {
								variables[name.Name] = identifier.Name
							}
						}
					}
					for index, name := range spec.Names {
						if name.Name != "difficultyReg" || index >= len(spec.Values) {
							continue
						}
						call, ok := spec.Values[index].(*ast.CallExpr)
						if ok && len(call.Args) == 1 {
							registry, _ = call.Args[0].(*ast.CompositeLit)
						}
					}
				}
			}
		case *ast.FuncDecl:
			if value.Recv == nil && value.Name.IsExported() {
				functions[value.Name.Name] = value
			}
		}
	}
	if err := inspectDifficultyInterface(types["Difficulty"]); err != nil {
		return difficultySpec{}, err
	}
	if function := functions["DifficultyByID"]; function == nil ||
		goFunctionSignature(function) != (goSignature{Parameters: "int", Results: "Difficulty, bool"}) {
		return difficultySpec{}, fmt.Errorf("Dragonfly world.DifficultyByID signature changed")
	}
	if function := functions["DifficultyID"]; function == nil ||
		goFunctionSignature(function) != (goSignature{Parameters: "Difficulty", Results: "int, bool"}) {
		return difficultySpec{}, fmt.Errorf("Dragonfly world.DifficultyID signature changed")
	}
	entries, err := inspectDifficultyRegistry(registry)
	if err != nil {
		return difficultySpec{}, err
	}
	if len(entries) != len(difficultyNames) {
		return difficultySpec{}, fmt.Errorf("Dragonfly difficulty registry has %d entries, want exactly %d", len(entries), len(difficultyNames))
	}
	live := map[string]world.Difficulty{
		"DifficultyPeaceful": world.DifficultyPeaceful,
		"DifficultyEasy":     world.DifficultyEasy,
		"DifficultyNormal":   world.DifficultyNormal,
		"DifficultyHard":     world.DifficultyHard,
	}
	for _, name := range difficultyNames {
		privateType := variables[name]
		if privateType == "" {
			return difficultySpec{}, fmt.Errorf("Dragonfly world.%s variable declaration not found", name)
		}
		structure, ok := types[privateType].Type.(*ast.StructType)
		if !ok || len(structure.Fields.List) != 0 {
			return difficultySpec{}, fmt.Errorf("Dragonfly world.%s concrete type is not an empty private struct", name)
		}
	}
	result := difficultySpec{Methods: append([]difficultyMethod(nil), difficultyMethods...)}
	for _, entry := range entries {
		value := live[entry.Name]
		if value == nil {
			return difficultySpec{}, fmt.Errorf("Dragonfly difficulty registry contains unexpected %s", entry.Name)
		}
		lookedUp, ok := world.DifficultyByID(entry.ID)
		if !ok || lookedUp != value {
			return difficultySpec{}, fmt.Errorf("Dragonfly live difficulty ID %d does not resolve to %s", entry.ID, entry.Name)
		}
		id, ok := world.DifficultyID(value)
		if !ok || id != entry.ID {
			return difficultySpec{}, fmt.Errorf("Dragonfly live %s reverse ID is %d, %v", entry.Name, id, ok)
		}
		if got := reflect.TypeOf(value).Name(); got != variables[entry.Name] {
			return difficultySpec{}, fmt.Errorf("Dragonfly live %s type is %s, want %s", entry.Name, got, variables[entry.Name])
		}
		result.Values = append(result.Values, difficultyValue{
			Name:                  entry.Name,
			PrivateType:           variables[entry.Name],
			ID:                    entry.ID,
			FoodRegenerates:       value.FoodRegenerates(),
			StarvationHealthLimit: value.StarvationHealthLimit(),
			FireSpreadIncrease:    value.FireSpreadIncrease(),
		})
	}
	unknown := math.MaxInt32
	for _, entry := range entries {
		if entry.ID == unknown {
			unknown--
		}
	}
	if fallback, ok := world.DifficultyByID(unknown); ok || fallback != world.DifficultyNormal {
		return difficultySpec{}, fmt.Errorf("Dragonfly unknown difficulty fallback changed")
	}
	return result, nil
}

func inspectDifficultyInterface(spec *ast.TypeSpec) error {
	if spec == nil {
		return fmt.Errorf("Dragonfly world.Difficulty interface not found")
	}
	interfaceType, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		return fmt.Errorf("Dragonfly world.Difficulty is not an interface")
	}
	found := make([]difficultyMethod, 0, len(interfaceType.Methods.List))
	for _, field := range interfaceType.Methods.List {
		function, ok := field.Type.(*ast.FuncType)
		if !ok || len(field.Names) != 1 || function.Params == nil || function.Params.NumFields() != 0 ||
			function.Results == nil || len(function.Results.List) != 1 {
			return fmt.Errorf("Dragonfly world.Difficulty contains an unsupported method")
		}
		result := formatGoExpression(function.Results.List[0].Type)
		method := difficultyMethod{Name: field.Names[0].Name, GoResult: result}
		for _, expected := range difficultyMethods {
			if expected.Name == method.Name {
				method.CSharpResult = expected.CSharpResult
			}
		}
		found = append(found, method)
	}
	if !reflect.DeepEqual(found, difficultyMethods) {
		return fmt.Errorf("Dragonfly world.Difficulty methods changed: got %+v", found)
	}
	return nil
}

type difficultyRegistryEntry struct {
	ID   int
	Name string
}

func inspectDifficultyRegistry(registry *ast.CompositeLit) ([]difficultyRegistryEntry, error) {
	if registry == nil || formatGoExpression(registry.Type) != "map[int]Difficulty" {
		return nil, fmt.Errorf("Dragonfly difficulty registry changed")
	}
	entries := make([]difficultyRegistryEntry, 0, len(registry.Elts))
	seenIDs := map[int]bool{}
	seenNames := map[string]bool{}
	for _, raw := range registry.Elts {
		pair, ok := raw.(*ast.KeyValueExpr)
		if !ok {
			return nil, fmt.Errorf("Dragonfly difficulty registry entry changed")
		}
		literal, ok := pair.Key.(*ast.BasicLit)
		if !ok || literal.Kind != token.INT {
			return nil, fmt.Errorf("Dragonfly difficulty registry ID changed")
		}
		id, err := strconv.Atoi(literal.Value)
		if err != nil {
			return nil, fmt.Errorf("Dragonfly difficulty registry ID changed: %w", err)
		}
		name, ok := pair.Value.(*ast.Ident)
		if !ok || !name.IsExported() || seenIDs[id] || seenNames[name.Name] {
			return nil, fmt.Errorf("Dragonfly difficulty registry entry changed")
		}
		seenIDs[id], seenNames[name.Name] = true, true
		entries = append(entries, difficultyRegistryEntry{ID: id, Name: name.Name})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })
	return entries, nil
}

func generateDifficulties(spec difficultySpec) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/difficulty.go AST and live registry. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\nusing System;\nusing Dragonfly.Native;\n\nnamespace Dragonfly;\n\n")
	output.WriteString("public sealed partial class World\n{\n")
	output.WriteString("    public interface Difficulty\n    {\n")
	for _, method := range spec.Methods {
		fmt.Fprintf(&output, "        %s %s();\n", method.CSharpResult, method.Name)
	}
	output.WriteString("    }\n\n")
	for _, value := range spec.Values {
		fmt.Fprintf(&output, "    public static readonly Difficulty %s = new BuiltinDifficulty(%d, %s, %s, %d);\n",
			value.Name, value.ID, strconv.FormatBool(value.FoodRegenerates), csharpDouble(value.StarvationHealthLimit), value.FireSpreadIncrease)
	}
	output.WriteString("\n    public static (Difficulty Difficulty, bool Ok) DifficultyByID(int id) => id switch\n    {\n")
	for _, value := range spec.Values {
		fmt.Fprintf(&output, "        %d => (%s, true),\n", value.ID, value.Name)
	}
	output.WriteString("        _ => (DifficultyNormal, false),\n    };\n\n")
	output.WriteString(`    public static (int ID, bool Ok) DifficultyID(Difficulty diff)
    {
        if (diff is BuiltinDifficulty builtin) return (builtin.ID, true);
        return (0, false);
    }

    internal static DifficultyView DifficultyView(Difficulty difficulty)
    {
        ArgumentNullException.ThrowIfNull(difficulty);
        return new DifficultyView
        {
            ID = difficulty is BuiltinDifficulty builtin ? checked((uint)builtin.ID) : 0,
            Builtin = difficulty is BuiltinDifficulty ? (byte)1 : (byte)0,
            FoodRegenerates = difficulty.FoodRegenerates() ? (byte)1 : (byte)0,
            StarvationHealthLimit = difficulty.StarvationHealthLimit(),
            FireSpreadIncrease = difficulty.FireSpreadIncrease(),
        };
    }

    internal static Difficulty DifficultyFromView(DifficultyView view)
    {
        if (view.Builtin > 1 || view.FoodRegenerates > 1)
            throw new InvalidOperationException("invalid difficulty view");
        if (view.Builtin == 1)
        {
            if (view.ID > int.MaxValue)
                throw new InvalidOperationException("invalid difficulty view");
            var (difficulty, ok) = DifficultyByID((int)view.ID);
            if (!ok) throw new InvalidOperationException("invalid difficulty view");
            return difficulty;
        }
        if (view.ID != 0)
            throw new InvalidOperationException("invalid difficulty view");
        return new CapabilityDifficulty(
            view.FoodRegenerates != 0,
            view.StarvationHealthLimit,
            view.FireSpreadIncrease);
    }

    private class CapabilityDifficulty(
        bool foodRegenerates,
        double starvationHealthLimit,
        int fireSpreadIncrease) : Difficulty
    {
        public bool FoodRegenerates() => foodRegenerates;
        public double StarvationHealthLimit() => starvationHealthLimit;
        public int FireSpreadIncrease() => fireSpreadIncrease;
    }

    private sealed class BuiltinDifficulty(
        int id,
        bool foodRegenerates,
        double starvationHealthLimit,
        int fireSpreadIncrease) : CapabilityDifficulty(foodRegenerates, starvationHealthLimit, fireSpreadIncrease)
    {
        internal int ID { get; } = id;
    }
}
`)
	return output.Bytes()
}
