package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
)

var selectedSoundTypes = []string{
	"AnvilBreak", "AnvilLand", "AnvilUse", "ArrowHit", "BarrelClose", "BarrelOpen",
	"BlastFurnaceCrackle", "BowShoot", "Burning", "Burp", "CampfireCrackle", "ChestClose",
	"ChestOpen", "Click", "ComposterEmpty", "ComposterFill", "ComposterFillLayer", "ComposterReady",
	"CopperScraped", "CrossbowShoot", "DecoratedPotInsertFailed", "Deny", "DoorCrash", "Drowning",
	"EnderChestClose", "EnderChestOpen", "Experience", "Explosion", "FireCharge", "FireExtinguish",
	"FireworkBlast", "FireworkHugeBlast", "FireworkLaunch", "FireworkTwinkle", "Fizz", "FurnaceCrackle",
	"GhastShoot", "GhastWarning", "GlassBreak", "Ignite", "ItemAdd", "ItemBreak", "ItemFrameRemove",
	"ItemFrameRotate", "ItemThrow", "LecternBookPlace", "LevelUp", "LightningExplode", "LightningThunder",
	"MusicDiscEnd", "Pop", "PotionBrewed", "PowerOff", "PowerOn", "SignWaxed", "SmokerCrackle",
	"StopUsingSpyglass", "TNT", "Teleport", "Thunder", "Totem", "UseSpyglass", "WaxRemoved",
	"WaxedSignFailedInteraction", "ShulkerBoxOpen", "ShulkerBoxClose", "EnderEyePlaced", "EndPortalCreated",
	"Attack", "Fall", "BlockPlace", "BlockBreaking", "DoorOpen", "DoorClose", "TrapdoorOpen",
	"TrapdoorClose", "FenceGateOpen", "FenceGateClose", "Note", "MusicDiscPlay", "DecoratedPotInserted",
	"ItemUseOn", "EquipItem", "BucketFill", "BucketEmpty", "CrossbowLoad", "GoatHorn",
}

type soundTypeSpec struct {
	Name   string
	Fields []parameter
}

func inspectPlayerPlaySound(path string) (commandMethod, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return commandMethod{}, err
	}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Name.Name != "PlaySound" || !pointerReceiver(function, "Player") {
			continue
		}
		if signature := goFunctionSignature(function); signature != (goSignature{Parameters: "world.Sound"}) {
			return commandMethod{}, fmt.Errorf("Dragonfly player.Player.PlaySound signature changed: %+v", signature)
		}
		if function.Type.Params == nil || len(function.Type.Params.List) != 1 ||
			len(function.Type.Params.List[0].Names) != 1 {
			return commandMethod{}, fmt.Errorf("Dragonfly player.Player.PlaySound parameter shape changed")
		}
		return commandMethod{
			Name: "PlaySound", ReturnType: "void",
			Parameters: []parameter{{Name: function.Type.Params.List[0].Names[0].Name, Type: "World.Sound"}},
		}, nil
	}
	return commandMethod{}, fmt.Errorf("Dragonfly player.Player has no PlaySound method")
}

func inspectSounds(directory string) ([]soundTypeSpec, error) {
	packages, err := parser.ParseDir(token.NewFileSet(), directory, nil, 0)
	if err != nil {
		return nil, err
	}
	pkg, ok := packages["sound"]
	if !ok {
		return nil, fmt.Errorf("Dragonfly sound package not found")
	}
	concrete := map[string]*ast.StructType{}
	for _, file := range pkg.Files {
		for _, declaration := range file.Decls {
			gen, ok := declaration.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, raw := range gen.Specs {
				typeSpec, ok := raw.(*ast.TypeSpec)
				if !ok || !ast.IsExported(typeSpec.Name.Name) {
					continue
				}
				structure, ok := typeSpec.Type.(*ast.StructType)
				if !ok || !embedsPrivateSound(structure) {
					continue
				}
				concrete[typeSpec.Name.Name] = structure
			}
		}
	}
	unknown := make([]string, 0)
	for name := range concrete {
		if !selectedSoundType(name) {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) != 0 {
		sort.Strings(unknown)
		return nil, fmt.Errorf("unknown concrete Dragonfly sound types require ABI review: %s", strings.Join(unknown, ", "))
	}
	result := make([]soundTypeSpec, 0, len(selectedSoundTypes))
	for _, name := range selectedSoundTypes {
		structure, ok := concrete[name]
		if !ok {
			return nil, fmt.Errorf("Dragonfly sound.%s type not found", name)
		}
		definition := soundTypeSpec{Name: name}
		for _, field := range structure.Fields.List {
			if len(field.Names) == 0 {
				continue
			}
			fieldType, ok := soundCSharpType(field.Type)
			if !ok {
				return nil, fmt.Errorf("sound.%s has unsupported field type %s", name, formatGoExpression(field.Type))
			}
			for _, fieldName := range field.Names {
				if ast.IsExported(fieldName.Name) {
					definition.Fields = append(definition.Fields, parameter{Name: fieldName.Name, Type: fieldType})
				}
			}
		}
		result = append(result, definition)
	}
	return result, nil
}

func embedsPrivateSound(structure *ast.StructType) bool {
	for _, field := range structure.Fields.List {
		if len(field.Names) != 0 {
			continue
		}
		expression := field.Type
		if pointer, ok := expression.(*ast.StarExpr); ok {
			expression = pointer.X
		}
		if name, ok := expression.(*ast.Ident); ok && name.Name == "sound" {
			return true
		}
	}
	return false
}

func selectedSoundType(name string) bool {
	for _, selected := range selectedSoundTypes {
		if name == selected {
			return true
		}
	}
	return false
}

func soundCSharpType(expression ast.Expr) (string, bool) {
	switch value := expression.(type) {
	case *ast.Ident:
		typeName, ok := map[string]string{
			"bool":       "bool",
			"float64":    "double",
			"int":        "int",
			"DiscType":   "DiscType",
			"Horn":       "Horn",
			"Instrument": "Instrument",
		}[value.Name]
		return typeName, ok
	case *ast.SelectorExpr:
		pkg, ok := value.X.(*ast.Ident)
		if !ok || pkg.Name != "world" {
			return "", false
		}
		typeName, ok := map[string]string{
			"Block":  "World.Block",
			"Item":   "World.Item",
			"Liquid": "World.Liquid",
		}[value.Sel.Name]
		return typeName, ok
	default:
		return "", false
	}
}

func generateSounds(types []soundTypeSpec) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/sound Go AST. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\n\nnamespace Dragonfly;\n\n")
	output.WriteString("public static partial class Sound\n{\n")
	for _, definition := range types {
		if len(definition.Fields) == 0 {
			fmt.Fprintf(&output, "    public readonly record struct %s : World.Sound;\n", definition.Name)
			continue
		}
		fmt.Fprintf(&output, "    public readonly record struct %s(%s) : World.Sound;\n",
			definition.Name, formatParameters(definition.Fields))
	}
	output.WriteString("\n    internal static World.Sound DecodeEvent(\n")
	output.WriteString("        uint kind, uint data, int integer, uint flags, double scalar, World.Block? block, World.Item? item) =>\n")
	output.WriteString("        kind switch\n        {\n")
	for index, definition := range types {
		fmt.Fprintf(&output, "            %d => new %s(%s),\n", index, definition.Name, soundDecodeArguments(definition))
	}
	output.WriteString("            _ => throw new InvalidOperationException(\"Invalid sound kind.\"),\n")
	output.WriteString("        };\n")
	output.WriteString("}\n\n")
	output.WriteString("internal readonly record struct EncodedSound(\n")
	output.WriteString("    uint Kind, uint Data, int Integer, uint Flags, double Scalar, World.Block? Block, World.Item? Item);\n\n")
	output.WriteString("internal static class SoundCodec\n{\n")
	output.WriteString("    internal static bool TryEncode(World.Sound sound, out EncodedSound encoded)\n    {\n")
	output.WriteString("        switch (sound)\n        {\n")
	for index, definition := range types {
		binding := "_"
		if len(definition.Fields) != 0 {
			binding = "value"
		}
		guard := soundEncodeGuard(definition)
		fmt.Fprintf(&output, "            case Sound.%s %s%s:\n", definition.Name, binding, guard)
		fmt.Fprintf(&output, "                encoded = new(%du, %s); return true;\n", index, soundEncodeArguments(definition))
	}
	output.WriteString("            default:\n                encoded = default; return false;\n")
	output.WriteString("        }\n    }\n")
	output.WriteString("}\n")
	return output.Bytes()
}

func generatePlayerPlaySound(method commandMethod) []byte {
	if method.Name != "PlaySound" || method.ReturnType != "void" || len(method.Parameters) != 1 ||
		method.Parameters[0].Type != "World.Sound" {
		panic("unsupported Player.PlaySound method")
	}
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.\n")
	output.WriteString("namespace Dragonfly;\n\npublic sealed partial class Player\n{\n")
	fmt.Fprintf(&output, "    public void PlaySound(World.Sound %s) =>\n", method.Parameters[0].Name)
	fmt.Fprintf(&output, "        PluginBridge.Host.PlayPlayerSound(_invocation, Id, %s);\n", method.Parameters[0].Name)
	output.WriteString("}\n")
	return output.Bytes()
}

func soundEncodeGuard(definition soundTypeSpec) string {
	for _, field := range definition.Fields {
		if field.Type == "World.Liquid" {
			return " when value." + field.Name + " is Block.Water or Block.Lava"
		}
	}
	return ""
}

func soundEncodeArguments(definition soundTypeSpec) string {
	values := map[string]string{
		"data": "0u", "integer": "0", "flags": "0u", "scalar": "0d", "block": "null", "item": "null",
	}
	for _, field := range definition.Fields {
		expression := "value." + field.Name
		switch field.Type {
		case "bool":
			values["flags"] = expression + " ? 1u : 0u"
		case "double":
			values["scalar"] = expression
		case "int":
			values["integer"] = expression
		case "World.Block":
			values["block"] = expression
		case "World.Item":
			values["item"] = expression
		case "World.Liquid":
			values["block"] = expression
			values["data"] = expression + " is Block.Water ? 0u : 1u"
		case "Instrument":
			values["data"] = expression + ".Id"
		case "DiscType", "Horn":
			values["data"] = expression + ".Uint8()"
		default:
			panic("unsupported sound encode field " + definition.Name + "." + field.Name)
		}
	}
	return strings.Join([]string{
		values["data"], values["integer"], values["flags"], values["scalar"], values["block"], values["item"],
	}, ", ")
}

func soundDecodeArguments(definition soundTypeSpec) string {
	arguments := make([]string, 0, len(definition.Fields))
	for _, field := range definition.Fields {
		value := map[string]string{
			"bool:Damage":           "flags != 0",
			"bool:QuickCharge":      "flags != 0",
			"double:Distance":       "scalar",
			"double:Progress":       "scalar",
			"int:Pitch":             "integer",
			"int:Stage":             "integer",
			"World.Block:Block":     "block ?? throw new InvalidOperationException(\"Sound requires a block.\")",
			"World.Item:Item":       "item ?? throw new InvalidOperationException(\"Sound requires an item.\")",
			"World.Liquid:Liquid":   "block is World.Liquid liquid ? liquid : throw new InvalidOperationException(\"Sound requires a liquid.\")",
			"Instrument:Instrument": "new Instrument(data)",
			"DiscType:DiscType":     "new DiscType(checked((int)data))",
			"Horn:Horn":             "new Horn(checked((int)data))",
		}[field.Type+":"+field.Name]
		if value == "" {
			panic("unsupported sound decode field " + definition.Name + "." + field.Name)
		}
		arguments = append(arguments, value)
	}
	return strings.Join(arguments, ", ")
}
