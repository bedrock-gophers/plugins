package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"image/color"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	dfeffect "github.com/df-mc/dragonfly/server/entity/effect"
)

type effectTypeSpec struct {
	Name    string
	ID      int
	Lasting bool
	Colour  color.RGBA
}

type effectSpec struct {
	Types []effectTypeSpec
}

var selectedPlayerEffectMethods = []string{
	"AddEffect",
	"RemoveEffect",
	"Effect",
	"Effects",
}

func inspectEffects(directory string) (effectSpec, error) {
	packages, err := parser.ParseDir(token.NewFileSet(), directory, func(info os.FileInfo) bool {
		return !info.IsDir() && filepath.Ext(info.Name()) == ".go" && !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		return effectSpec{}, err
	}
	pkg := packages["effect"]
	if pkg == nil {
		return effectSpec{}, fmt.Errorf("Dragonfly effect package not found")
	}
	types := map[string]*ast.TypeSpec{}
	functions := map[string]*ast.FuncDecl{}
	methods := map[string]map[string]*ast.FuncDecl{}
	variables := map[string]bool{}
	var initFunctions []*ast.FuncDecl
	for _, file := range pkg.Files {
		for _, declaration := range file.Decls {
			switch value := declaration.(type) {
			case *ast.GenDecl:
				for _, raw := range value.Specs {
					switch spec := raw.(type) {
					case *ast.TypeSpec:
						types[spec.Name.Name] = spec
					case *ast.ValueSpec:
						for _, name := range spec.Names {
							variables[name.Name] = true
						}
					}
				}
			case *ast.FuncDecl:
				if value.Recv == nil {
					if value.Name.Name == "init" {
						initFunctions = append(initFunctions, value)
					} else {
						functions[value.Name.Name] = value
					}
					continue
				}
				receiver := receiverTypeName(value)
				if methods[receiver] == nil {
					methods[receiver] = map[string]*ast.FuncDecl{}
				}
				methods[receiver][value.Name.Name] = value
			}
		}
	}
	if err := validateEffectAST(types, functions, methods); err != nil {
		return effectSpec{}, err
	}
	registered := map[int]string{}
	for _, function := range initFunctions {
		ast.Inspect(function.Body, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			name, ok := call.Fun.(*ast.Ident)
			if !ok || name.Name != "Register" || len(call.Args) != 2 {
				return true
			}
			idLiteral, idOK := call.Args[0].(*ast.BasicLit)
			value, valueOK := call.Args[1].(*ast.Ident)
			if !idOK || idLiteral.Kind != token.INT || !valueOK || !variables[value.Name] {
				registered[-1] = ""
				return true
			}
			id, parseErr := strconv.Atoi(idLiteral.Value)
			if parseErr != nil || id < 0 || registered[id] != "" {
				registered[-1] = ""
				return true
			}
			registered[id] = value.Name
			return true
		})
	}
	if _, malformed := registered[-1]; malformed {
		return effectSpec{}, fmt.Errorf("Dragonfly effect registration AST changed")
	}
	if len(registered) == 0 {
		return effectSpec{}, fmt.Errorf("Dragonfly has no registered effects")
	}
	result := effectSpec{Types: make([]effectTypeSpec, 0, len(registered))}
	for id, name := range registered {
		value, ok := dfeffect.ByID(id)
		if !ok {
			return effectSpec{}, fmt.Errorf("Dragonfly effect.%s ID %d missing from live registry", name, id)
		}
		if liveID, found := dfeffect.ID(value); !found || liveID != id {
			return effectSpec{}, fmt.Errorf("Dragonfly effect.%s live ID changed", name)
		}
		result.Types = append(result.Types, effectTypeSpec{
			Name: name, ID: id, Colour: value.RGBA(),
			Lasting: reflect.TypeOf(value).Implements(reflect.TypeFor[dfeffect.LastingType]()),
		})
	}
	sort.Slice(result.Types, func(i, j int) bool { return result.Types[i].ID < result.Types[j].ID })
	return result, nil
}

func validateEffectAST(types map[string]*ast.TypeSpec, functions map[string]*ast.FuncDecl, methods map[string]map[string]*ast.FuncDecl) error {
	if err := validateItemInterface(types["Type"], "Type", []string{"RGBA()->color.RGBA", "Apply(world.Entity,Effect)->"}); err != nil {
		return err
	}
	if err := validateItemInterface(types["LastingType"], "LastingType", []string{
		"embed:Type", "Start(world.Entity,int)->", "End(world.Entity,int)->",
	}); err != nil {
		return err
	}
	if fields := allItemFields(types["Effect"]); !reflect.DeepEqual(fields, []string{
		"t Type", "d time.Duration", "lvl int", "potency float64", "ambient bool", "particlesHidden bool", "infinite bool", "tick int",
	}) {
		return fmt.Errorf("Dragonfly effect.Effect fields changed: %v", fields)
	}
	for name, signature := range map[string]goSignature{
		"NewInstant":            {Parameters: "Type, int", Results: "Effect"},
		"NewInstantWithPotency": {Parameters: "Type, int, float64", Results: "Effect"},
		"New":                   {Parameters: "LastingType, int, time.Duration", Results: "Effect"},
		"NewAmbient":            {Parameters: "LastingType, int, time.Duration", Results: "Effect"},
		"NewInfinite":           {Parameters: "LastingType, int", Results: "Effect"},
		"ResultingColour":       {Parameters: "[]Effect", Results: "color.RGBA, bool"},
	} {
		function := functions[name]
		if function == nil || goFunctionSignature(function) != signature {
			return fmt.Errorf("Dragonfly effect.%s signature changed", name)
		}
	}
	for name, signature := range map[string]goSignature{
		"WithoutParticles": {Results: "Effect"},
		"ParticlesHidden":  {Results: "bool"},
		"Level":            {Results: "int"},
		"Duration":         {Results: "time.Duration"},
		"Ambient":          {Results: "bool"},
		"Infinite":         {Results: "bool"},
		"Type":             {Results: "Type"},
		"TickDuration":     {Results: "Effect"},
		"Tick":             {Results: "int"},
	} {
		method := methods["Effect"][name]
		if method == nil || !valueReceiver(method, "Effect") || goFunctionSignature(method) != signature {
			return fmt.Errorf("Dragonfly effect.Effect.%s signature changed", name)
		}
	}
	return nil
}

func inspectPlayerEffectMethods(path string) ([]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, err
	}
	found := map[string]*ast.FuncDecl{}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && playerMethod(function) {
			found[function.Name.Name] = function
		}
	}
	want := map[string]goSignature{
		"AddEffect":    {Parameters: "effect.Effect"},
		"RemoveEffect": {Parameters: "effect.Type"},
		"Effect":       {Parameters: "effect.Type", Results: "effect.Effect, bool"},
		"Effects":      {Results: "[]effect.Effect"},
	}
	for _, name := range selectedPlayerEffectMethods {
		if found[name] == nil || goFunctionSignature(found[name]) != want[name] {
			return nil, fmt.Errorf("Dragonfly player.Player.%s signature changed", name)
		}
	}
	return append([]string(nil), selectedPlayerEffectMethods...), nil
}

func generateEffects(spec effectSpec) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/entity/effect Go AST and live registry. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\nusing System;\nusing System.Collections.Generic;\n\nnamespace Dragonfly;\n\n")
	output.WriteString(`public static partial class Effect
{
    public interface Type
    {
        Color.RGBA RGBA();
    }

    public interface LastingType : Type { }

    public readonly record struct Value
    {
        private Type? _type { get; init; }
        private TimeSpan _duration { get; init; }
        private int _level { get; init; }
        private double _potency { get; init; }
        private bool _ambient { get; init; }
        private bool _particlesHidden { get; init; }
        private bool _infinite { get; init; }
        private int _tick { get; init; }

        internal Value(Type type, TimeSpan duration, int level, double potency, bool ambient, bool particlesHidden, bool infinite, int tick)
        {
            _type = type;
            _duration = duration;
            _level = level;
            _potency = potency;
            _ambient = ambient;
            _particlesHidden = particlesHidden;
            _infinite = infinite;
            _tick = tick;
        }

        public Value WithoutParticles() => this with { _particlesHidden = true };
        public bool ParticlesHidden() => _particlesHidden;
        public int Level() => _level;
        public TimeSpan Duration() => _duration;
        public bool Ambient() => _ambient;
        public bool Infinite() => _infinite;
        public Type? Type() => _type;

        public Value TickDuration()
        {
            if (_type is not LastingType) return this;
            return this with
            {
                _duration = _infinite ? _duration : _duration - TimeSpan.FromMilliseconds(50),
                _tick = _tick + 1,
            };
        }

        public int Tick() => _tick;

        internal double Potency => _potency;
    }

    public static Value NewInstant(Type type, int level) => NewInstantWithPotency(type, level, 1d);

    public static Value NewInstantWithPotency(Type type, int level, double potency)
    {
        ArgumentNullException.ThrowIfNull(type);
        return new Value(type, default, level, potency, false, false, false, 0);
    }

    public static Value New(LastingType type, int level, TimeSpan duration)
    {
        ArgumentNullException.ThrowIfNull(type);
        return new Value(type, duration, level, 0d, false, false, false, 0);
    }

    public static Value NewAmbient(LastingType type, int level, TimeSpan duration)
    {
        ArgumentNullException.ThrowIfNull(type);
        return new Value(type, duration, level, 0d, true, false, false, 0);
    }

    public static Value NewInfinite(LastingType type, int level)
    {
        ArgumentNullException.ThrowIfNull(type);
        return new Value(type, default, level, 0d, false, false, true, 0);
    }

    public static (Color.RGBA Colour, bool Ambient) ResultingColour(IReadOnlyList<Value> effects)
    {
        ArgumentNullException.ThrowIfNull(effects);
        long red = 0;
        long green = 0;
        long blue = 0;
        long alpha = 0;
        var count = 0;
        var ambient = true;
        foreach (var effect in effects)
        {
            if (effect.ParticlesHidden()) continue;
            var type = effect.Type() ?? throw new InvalidOperationException("Effect has no type.");
            var colour = type.RGBA();
            red += colour.R;
            green += colour.G;
            blue += colour.B;
            alpha += colour.A;
            count++;
            if (!effect.Ambient()) ambient = false;
        }
        if (count == 0) return (new Color.RGBA(0x38, 0x5d, 0xc6, 0xff), false);
        return (new Color.RGBA((byte)(red / count), (byte)(green / count), (byte)(blue / count), (byte)(alpha / count)), ambient);
    }

`)
	for _, effect := range spec.Types {
		typeName := "Type"
		implementation := "BuiltinInstantType"
		if effect.Lasting {
			typeName = "LastingType"
			implementation = "BuiltinLastingType"
		}
		fmt.Fprintf(&output, "    public static readonly %s %s = new %s(%d, new Color.RGBA(%d, %d, %d, %d));\n",
			typeName, effect.Name, implementation, effect.ID, effect.Colour.R, effect.Colour.G, effect.Colour.B, effect.Colour.A)
	}
	output.WriteString(`
    internal static bool TryID(Type? type, out int id)
    {
        switch (type)
        {
            case BuiltinInstantType instant:
                id = instant.ID;
                return true;
            case BuiltinLastingType lasting:
                id = lasting.ID;
                return true;
            default:
                id = 0;
                return false;
        }
    }

    public static (Type? Type, bool Ok) ByID(int id)
    {
        var type = TypeByID(id);
        return (type, type is not null);
    }

    public static (int ID, bool Ok) ID(Type type)
    {
        ArgumentNullException.ThrowIfNull(type);
        return TryID(type, out var id) ? (id, true) : (0, false);
    }

    internal static Type? TypeByID(int id) => id switch
    {
`)
	for _, effect := range spec.Types {
		fmt.Fprintf(&output, "        %d => %s,\n", effect.ID, effect.Name)
	}
	output.WriteString(`        _ => null,
    };

    private sealed record BuiltinInstantType(int ID, Color.RGBA Colour) : Type
    {
        public Color.RGBA RGBA() => Colour;
    }

    private sealed record BuiltinLastingType(int ID, Color.RGBA Colour) : LastingType
    {
        public Color.RGBA RGBA() => Colour;
    }
}
`)
	return output.Bytes()
}

func generatePlayerEffects(methods []string) []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\nusing System.Collections.Generic;\n\nnamespace Dragonfly;\n\n")
	output.WriteString("public sealed partial class Player\n{\n")
	for _, method := range methods {
		switch method {
		case "AddEffect":
			output.WriteString("    public void AddEffect(Effect.Value e) => PluginBridge.Host.AddPlayerEffect(_invocation, Id, e);\n")
		case "RemoveEffect":
			output.WriteString("    public void RemoveEffect(Effect.Type e) => PluginBridge.Host.RemovePlayerEffect(_invocation, Id, e);\n")
		case "Effect":
			output.WriteString("    public (Effect.Value Effect, bool Ok) Effect(Effect.Type e) => PluginBridge.Host.PlayerEffect(_invocation, Id, e);\n")
		case "Effects":
			output.WriteString("    public IReadOnlyList<Effect.Value> Effects() => PluginBridge.Host.PlayerEffects(_invocation, Id);\n")
		default:
			panic("unsupported player effect method: " + method)
		}
	}
	output.WriteString("}\n")
	return output.Bytes()
}

func csharpEffectValue(value dfeffect.Effect, spec effectSpec) (string, error) {
	id, ok := dfeffect.ID(value.Type())
	if !ok {
		return "", fmt.Errorf("unregistered effect type %T", value.Type())
	}
	name := ""
	lasting := false
	for _, candidate := range spec.Types {
		if candidate.ID == id {
			name, lasting = candidate.Name, candidate.Lasting
			break
		}
	}
	if name == "" {
		return "", fmt.Errorf("effect ID %d missing from AST registry", id)
	}
	var result string
	if lasting {
		ticks := csharpDurationTicks(value.Duration())
		switch {
		case value.Infinite():
			result = fmt.Sprintf("Effect.NewInfinite(Effect.%s, %d)", name, value.Level())
		case value.Ambient():
			result = fmt.Sprintf("Effect.NewAmbient(Effect.%s, %d, TimeSpan.FromTicks(%d))", name, value.Level(), ticks)
		default:
			result = fmt.Sprintf("Effect.New(Effect.%s, %d, TimeSpan.FromTicks(%d))", name, value.Level(), ticks)
		}
	} else {
		result = fmt.Sprintf("Effect.NewInstant(Effect.%s, %d)", name, value.Level())
	}
	if value.ParticlesHidden() {
		result += ".WithoutParticles()"
	}
	return result, nil
}

func csharpEffectList(values []dfeffect.Effect, spec effectSpec) (string, error) {
	parts := make([]string, len(values))
	for index, value := range values {
		formatted, err := csharpEffectValue(value, spec)
		if err != nil {
			return "", err
		}
		parts[index] = formatted
	}
	if len(parts) == 0 {
		return "Array.Empty<Effect.Value>()", nil
	}
	var output bytes.Buffer
	output.WriteString("new Effect.Value[] { ")
	for index, part := range parts {
		if index != 0 {
			output.WriteString(", ")
		}
		output.WriteString(part)
	}
	output.WriteString(" }")
	return output.String(), nil
}

func csharpEffectResult(value any, spec effectSpec) (string, error) {
	values, ok := value.([]dfeffect.Effect)
	if !ok {
		return "", fmt.Errorf("result is %T, not []effect.Effect", value)
	}
	return csharpEffectList(values, spec)
}
