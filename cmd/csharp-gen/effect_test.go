package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	dfeffect "github.com/df-mc/dragonfly/server/entity/effect"
)

func TestInspectEffectsUsesASTAndLiveRegistry(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	spec, err := inspectEffects(filepath.Join(string(bytes.TrimSpace(output)), "server", "entity", "effect"))
	if err != nil {
		t.Fatal(err)
	}
	if len(spec.Types) != 28 {
		t.Fatalf("effect count=%d, want 28", len(spec.Types))
	}
	want := map[int]string{
		1: "Speed", 2: "Slowness", 3: "Haste", 4: "MiningFatigue", 5: "Strength",
		6: "InstantHealth", 7: "InstantDamage", 8: "JumpBoost", 9: "Nausea", 10: "Regeneration",
		11: "Resistance", 12: "FireResistance", 13: "WaterBreathing", 14: "Invisibility",
		15: "Blindness", 16: "NightVision", 17: "Hunger", 18: "Weakness", 19: "Poison",
		20: "Wither", 21: "HealthBoost", 22: "Absorption", 23: "Saturation", 24: "Levitation",
		25: "FatalPoison", 26: "ConduitPower", 27: "SlowFalling", 30: "Darkness",
	}
	for _, value := range spec.Types {
		if want[value.ID] != value.Name {
			t.Fatalf("effect ID %d=%q, want %q", value.ID, value.Name, want[value.ID])
		}
		if value.Colour.A != 0xff {
			t.Fatalf("effect %s alpha=%d", value.Name, value.Colour.A)
		}
		if value.Lasting == (value.ID == 6 || value.ID == 7) {
			t.Fatalf("effect %s lasting=%v", value.Name, value.Lasting)
		}
	}
	generated := string(generateEffects(spec))
	for _, expected := range []string{
		"public static partial class Effect",
		"public interface Type",
		"public interface LastingType : Type",
		"public readonly record struct Value",
		"public static Value NewInstant(Type type, int level)",
		"public static Value NewInstantWithPotency(Type type, int level, double potency)",
		"public static Value New(LastingType type, int level, TimeSpan duration)",
		"public static Value NewAmbient(LastingType type, int level, TimeSpan duration)",
		"public static Value NewInfinite(LastingType type, int level)",
		"public Value WithoutParticles()",
		"public Value TickDuration()",
		"public static (Color.RGBA Colour, bool Ambient) ResultingColour",
		"public static readonly LastingType Speed = new BuiltinLastingType(1",
		"public static readonly Type InstantHealth = new BuiltinInstantType(6",
		"public static readonly LastingType Darkness = new BuiltinLastingType(30",
		"internal static bool TryID(Type? type, out int id)",
		"public static (Type? Type, bool Ok) ByID(int id)",
		"public static (int ID, bool Ok) ID(Type type)",
		"internal static Type? TypeByID(int id)",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated effects missing %q", expected)
		}
	}
	if strings.Contains(generated, "effect_type") || strings.Contains(generated, "minecraft:") {
		t.Fatalf("public effect API exposes transport details")
	}
}

func TestEffectValuesGenerateDragonflySemantics(t *testing.T) {
	spec := effectSpec{Types: []effectTypeSpec{
		{Name: "Speed", ID: 1, Lasting: true},
		{Name: "InstantHealth", ID: 6},
	}}
	values := []dfeffect.Effect{
		dfeffect.New(dfeffect.Speed, 2, time.Second),
		dfeffect.NewInstant(dfeffect.InstantHealth, 1).WithoutParticles(),
	}
	generated, err := csharpEffectList(values, spec)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Effect.New(Effect.Speed, 2, TimeSpan.FromTicks(10000000))",
		"Effect.NewInstant(Effect.InstantHealth, 1).WithoutParticles()",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("effect list missing %q: %s", expected, generated)
		}
	}
}

func TestPlayerEffectMethodsUseGoAST(t *testing.T) {
	path := filepath.Join(t.TempDir(), "player.go")
	source := `package player
func (p *Player) AddEffect(e effect.Effect) {}
func (p *Player) RemoveEffect(e effect.Type) {}
func (p *Player) Effect(e effect.Type) (effect.Effect, bool) { return effect.Effect{}, false }
func (p *Player) Effects() []effect.Effect { return nil }`
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	methods, err := inspectPlayerEffectMethods(path)
	if err != nil {
		t.Fatal(err)
	}
	generated := string(generatePlayerEffects(methods))
	for _, expected := range []string{
		"public void AddEffect(Effect.Value e)",
		"public void RemoveEffect(Effect.Type e)",
		"public (Effect.Value Effect, bool Ok) Effect(Effect.Type e)",
		"public IReadOnlyList<Effect.Value> Effects()",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated player effects missing %q", expected)
		}
	}
}
