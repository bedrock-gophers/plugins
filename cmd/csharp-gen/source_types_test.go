package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const sourceWorldFixture = `package world
type Entity interface{}
type Block interface{}
type DamageSource interface {
	ReducedByArmour() bool
	ReducedByResistance() bool
	Fire() bool
	IgnoreTotem() bool
}
type HealingSource interface { HealingSource() }
`

const sourceEntityFixture = `package entity
import (
	"example/item"
	"example/item/enchantment"
	"example/world"
)
type (
	AttackDamageSource struct { Attacker world.Entity }
	VoidDamageSource struct{}
	SuffocationDamageSource struct{}
	DrowningDamageSource struct{}
	FallDamageSource struct{}
	GlideDamageSource struct{}
	LightningDamageSource struct{}
	ProjectileDamageSource struct { Projectile, Owner world.Entity }
	ExplosionDamageSource struct{}
	FoodHealingSource struct { QuickRegeneration bool }
)
func (AttackDamageSource) ReducedByArmour() bool { return true }
func (AttackDamageSource) ReducedByResistance() bool { return true }
func (AttackDamageSource) Fire() bool { return false }
func (AttackDamageSource) IgnoreTotem() bool { return false }
func (VoidDamageSource) ReducedByArmour() bool { return false }
func (VoidDamageSource) ReducedByResistance() bool { return false }
func (VoidDamageSource) Fire() bool { return false }
func (VoidDamageSource) IgnoreTotem() bool { return true }
func (SuffocationDamageSource) ReducedByArmour() bool { return false }
func (SuffocationDamageSource) ReducedByResistance() bool { return false }
func (SuffocationDamageSource) Fire() bool { return false }
func (SuffocationDamageSource) IgnoreTotem() bool { return false }
func (DrowningDamageSource) ReducedByArmour() bool { return false }
func (DrowningDamageSource) ReducedByResistance() bool { return false }
func (DrowningDamageSource) Fire() bool { return false }
func (DrowningDamageSource) IgnoreTotem() bool { return false }
func (FallDamageSource) ReducedByArmour() bool { return false }
func (FallDamageSource) ReducedByResistance() bool { return true }
func (FallDamageSource) Fire() bool { return false }
func (FallDamageSource) IgnoreTotem() bool { return false }
func (FallDamageSource) AffectedByEnchantment(e item.EnchantmentType) bool { return e == enchantment.FeatherFalling }
func (GlideDamageSource) ReducedByArmour() bool { return false }
func (GlideDamageSource) ReducedByResistance() bool { return true }
func (GlideDamageSource) Fire() bool { return false }
func (GlideDamageSource) IgnoreTotem() bool { return false }
func (LightningDamageSource) ReducedByArmour() bool { return true }
func (LightningDamageSource) ReducedByResistance() bool { return true }
func (LightningDamageSource) Fire() bool { return false }
func (LightningDamageSource) IgnoreTotem() bool { return false }
func (ProjectileDamageSource) ReducedByArmour() bool { return true }
func (ProjectileDamageSource) ReducedByResistance() bool { return true }
func (ProjectileDamageSource) Fire() bool { return false }
func (ProjectileDamageSource) IgnoreTotem() bool { return false }
func (ProjectileDamageSource) AffectedByEnchantment(e item.EnchantmentType) bool { return e == enchantment.ProjectileProtection }
func (ExplosionDamageSource) ReducedByArmour() bool { return true }
func (ExplosionDamageSource) ReducedByResistance() bool { return true }
func (ExplosionDamageSource) Fire() bool { return false }
func (ExplosionDamageSource) IgnoreTotem() bool { return false }
func (ExplosionDamageSource) AffectedByEnchantment(e item.EnchantmentType) bool { return e == enchantment.BlastProtection }
func (FoodHealingSource) HealingSource() {}
`

const sourceEffectFixture = `package effect
type WitherDamageSource struct{}
type InstantDamageSource struct{}
type PoisonDamageSource struct { Fatal bool }
type InstantHealingSource struct{}
type RegenerationHealingSource struct{}
func (WitherDamageSource) ReducedByArmour() bool { return false }
func (WitherDamageSource) ReducedByResistance() bool { return true }
func (WitherDamageSource) Fire() bool { return false }
func (WitherDamageSource) IgnoreTotem() bool { return false }
func (InstantDamageSource) ReducedByArmour() bool { return false }
func (InstantDamageSource) ReducedByResistance() bool { return true }
func (InstantDamageSource) Fire() bool { return false }
func (InstantDamageSource) IgnoreTotem() bool { return false }
func (PoisonDamageSource) ReducedByArmour() bool { return false }
func (PoisonDamageSource) ReducedByResistance() bool { return true }
func (PoisonDamageSource) Fire() bool { return false }
func (PoisonDamageSource) IgnoreTotem() bool { return false }
func (InstantHealingSource) HealingSource() {}
func (RegenerationHealingSource) HealingSource() {}
`

const sourcePlayerFixture = `package player
type StarvationDamageSource struct{}
func (StarvationDamageSource) ReducedByArmour() bool { return false }
func (StarvationDamageSource) ReducedByResistance() bool { return false }
func (StarvationDamageSource) Fire() bool { return false }
func (StarvationDamageSource) IgnoreTotem() bool { return false }
`

const sourceBlockFixture = `package block
import (
	"example/item"
	"example/item/enchantment"
	"example/world"
)
type DamageSource struct { Block world.Block }
type MagmaDamageSource struct{}
type LavaDamageSource struct{}
type FireDamageSource struct{}
func (DamageSource) ReducedByArmour() bool { return true }
func (DamageSource) ReducedByResistance() bool { return true }
func (DamageSource) Fire() bool { return false }
func (DamageSource) IgnoreTotem() bool { return false }
func (MagmaDamageSource) ReducedByArmour() bool { return true }
func (MagmaDamageSource) ReducedByResistance() bool { return true }
func (MagmaDamageSource) Fire() bool { return true }
func (MagmaDamageSource) IgnoreTotem() bool { return false }
func (MagmaDamageSource) AffectedByEnchantment(e item.EnchantmentType) bool { return e == enchantment.FireProtection }
func (LavaDamageSource) ReducedByArmour() bool { return true }
func (LavaDamageSource) ReducedByResistance() bool { return true }
func (LavaDamageSource) Fire() bool { return true }
func (LavaDamageSource) IgnoreTotem() bool { return false }
func (FireDamageSource) ReducedByArmour() bool { return true }
func (FireDamageSource) ReducedByResistance() bool { return true }
func (FireDamageSource) Fire() bool { return true }
func (FireDamageSource) IgnoreTotem() bool { return false }
func (FireDamageSource) AffectedByEnchantment(e item.EnchantmentType) bool { return e == enchantment.FireProtection }
`

const sourceEnchantmentFixture = `package enchantment
import (
	"example/item"
	"example/world"
)
var FeatherFalling, ProjectileProtection, BlastProtection, FireProtection item.EnchantmentType
type AffectedDamageSource interface {
	world.DamageSource
	AffectedByEnchantment(e item.EnchantmentType) bool
}
type ThornsDamageSource struct { Owner world.Entity }
func (ThornsDamageSource) ReducedByArmour() bool { return false }
func (ThornsDamageSource) ReducedByResistance() bool { return true }
func (ThornsDamageSource) Fire() bool { return false }
func (ThornsDamageSource) IgnoreTotem() bool { return false }
`

type sourceFixtureDirs struct {
	world       string
	entity      string
	effect      string
	player      string
	block       string
	enchantment string
}

func writeSourceFixtures(t *testing.T, replacement ...string) sourceFixtureDirs {
	t.Helper()
	root := t.TempDir()
	values := map[string]string{
		"world": sourceWorldFixture, "entity": sourceEntityFixture, "effect": sourceEffectFixture,
		"player": sourcePlayerFixture, "block": sourceBlockFixture, "enchantment": sourceEnchantmentFixture,
	}
	if len(replacement) != 0 {
		values[replacement[0]] = strings.Replace(values[replacement[0]], replacement[1], replacement[2], 1)
	}
	directories := sourceFixtureDirs{}
	for name, source := range values {
		directory := filepath.Join(root, name)
		if err := os.MkdirAll(directory, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, name+".go"), []byte(source), 0o600); err != nil {
			t.Fatal(err)
		}
		switch name {
		case "world":
			directories.world = directory
		case "entity":
			directories.entity = directory
		case "effect":
			directories.effect = directory
		case "player":
			directories.player = directory
		case "block":
			directories.block = directory
		case "enchantment":
			directories.enchantment = directory
		}
	}
	return directories
}

func inspectSourceFixtures(directories sourceFixtureDirs) (sourceTypesSpec, error) {
	return inspectSourceTypes(directories.world, directories.entity, directories.effect, directories.player, directories.block, directories.enchantment)
}

func TestSourceTypesUseExactGoAST(t *testing.T) {
	spec, err := inspectSourceFixtures(writeSourceFixtures(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(spec.Types) != 21 {
		t.Fatalf("got %d source types, want 21", len(spec.Types))
	}
	healing := 0
	for _, definition := range spec.Types {
		if definition.Healing {
			healing++
		}
	}
	if healing != 3 {
		t.Fatalf("got %d healing source types, want 3", healing)
	}

	generated := string(generateSourceTypes(spec))
	for _, expected := range []string{
		"public interface DamageSource",
		"public interface HealingSource { }",
		"public static partial class Entity",
		"public readonly record struct AttackDamageSource(World.Entity? Attacker = null) : World.DamageSource",
		"public readonly record struct ProjectileDamageSource(World.Entity? Projectile = null, World.Entity? Owner = null) : Enchantment.AffectedDamageSource",
		"public static partial class Effect",
		"public readonly record struct PoisonDamageSource(bool Fatal = false) : World.DamageSource",
		"public sealed partial class Player",
		"public readonly record struct StarvationDamageSource : World.DamageSource",
		"public static partial class Block",
		"public readonly record struct DamageSource(World.Block? Block = null) : World.DamageSource",
		"public static partial class Enchantment",
		"public interface AffectedDamageSource : World.DamageSource",
		"public readonly record struct ThornsDamageSource(World.Entity? Owner = null) : World.DamageSource",
		"public bool AffectedByEnchantment(Item.EnchantmentType e) => object.Equals(e, Item.FeatherFalling);",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated source types missing %q:\n%s", expected, generated)
		}
	}
	for _, invented := range []string{"CustomDamageSource", "CustomHealingSource", "BlockDamageSource", "World.AttackDamageSource"} {
		if strings.Contains(generated, invented) {
			t.Fatalf("generated source types contain invented type %q:\n%s", invented, generated)
		}
	}
}

func TestSourceTypesRejectDrift(t *testing.T) {
	tests := []struct {
		name  string
		pkg   string
		old   string
		new   string
		error string
	}{
		{name: "world interface", pkg: "world", old: "Fire() bool", new: "Fire() int", error: "signature changed"},
		{name: "affected interface", pkg: "enchantment", old: "AffectedByEnchantment(e item.EnchantmentType) bool", new: "AffectedByEnchantment(e string) bool", error: "signature changed"},
		{name: "pointer receiver", pkg: "entity", old: "func (AttackDamageSource) ReducedByArmour()", new: "func (*AttackDamageSource) ReducedByArmour()", error: "signature changed"},
		{name: "field", pkg: "effect", old: "type PoisonDamageSource struct { Fatal bool }", new: "type PoisonDamageSource struct { Fatal int }", error: "unsupported field type"},
		{name: "missing source", pkg: "player", old: "type StarvationDamageSource struct{}", new: "type starvationDamageSource struct{}", error: "struct not found"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directories := writeSourceFixtures(t, test.pkg, test.old, test.new)
			if _, err := inspectSourceFixtures(directories); err == nil || !strings.Contains(err.Error(), test.error) {
				t.Fatalf("expected %q error, got %v", test.error, err)
			}
		})
	}
}

func TestSourceTypesASTControlsFieldsAndConstants(t *testing.T) {
	directories := writeSourceFixtures(t)
	entityPath := filepath.Join(directories.entity, "entity.go")
	source, err := os.ReadFile(entityPath)
	if err != nil {
		t.Fatal(err)
	}
	source = bytes.Replace(source,
		[]byte("func (VoidDamageSource) IgnoreTotem() bool { return true }"),
		[]byte("func (VoidDamageSource) IgnoreTotem() bool { return false }"), 1)
	source = bytes.Replace(source,
		[]byte("return e == enchantment.FeatherFalling"),
		[]byte("return e == enchantment.FireProtection"), 1)
	source = bytes.Replace(source,
		[]byte("FoodHealingSource struct { QuickRegeneration bool }"),
		[]byte("FoodHealingSource struct { Fast bool }"), 1)
	if err := os.WriteFile(entityPath, source, 0o600); err != nil {
		t.Fatal(err)
	}
	spec, err := inspectSourceFixtures(directories)
	if err != nil {
		t.Fatal(err)
	}
	generated := string(generateSourceTypes(spec))
	start := strings.Index(generated, "record struct VoidDamageSource")
	if start == -1 {
		t.Fatal("generated VoidDamageSource missing")
	}
	end := strings.Index(generated[start:], "    }\n")
	if end == -1 || !strings.Contains(generated[start:start+end], "IgnoreTotem() => false") {
		t.Fatalf("AST bool constant did not control output:\n%s", generated)
	}
	if !strings.Contains(generated, "AffectedByEnchantment(Item.EnchantmentType e) => object.Equals(e, Item.FireProtection)") {
		t.Fatalf("AST enchantment selector did not control output:\n%s", generated)
	}
	if !strings.Contains(generated, "record struct FoodHealingSource(bool Fast = false)") || strings.Contains(generated, "QuickRegeneration") {
		t.Fatalf("AST fields did not control output:\n%s", generated)
	}
}

func TestSourceTypesRejectUnknownConcreteSource(t *testing.T) {
	tests := []struct {
		pkg, marker, addition string
	}{
		{"entity", "type (", `type NewDamageSource struct{}
func (NewDamageSource) ReducedByArmour() bool { return false }
func (NewDamageSource) ReducedByResistance() bool { return false }
func (NewDamageSource) Fire() bool { return false }
func (NewDamageSource) IgnoreTotem() bool { return false }

type (`},
		{"world", "type Entity interface{}", `type NewDamageSource struct{}
func (NewDamageSource) ReducedByArmour() bool { return false }
func (NewDamageSource) ReducedByResistance() bool { return false }
func (NewDamageSource) Fire() bool { return false }
func (NewDamageSource) IgnoreTotem() bool { return false }

type Entity interface{}`},
	}
	for _, test := range tests {
		t.Run(test.pkg, func(t *testing.T) {
			directories := writeSourceFixtures(t, test.pkg, test.marker, test.addition)
			if _, err := inspectSourceFixtures(directories); err == nil || !strings.Contains(err.Error(), "unknown exported Dragonfly source types") {
				t.Fatalf("expected unknown source error, got %v", err)
			}
		})
	}
}

func TestPinnedDragonflyHasExactSourceTypes(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	root := string(bytes.TrimSpace(output))
	if _, err := inspectSourceTypes(
		filepath.Join(root, "server", "world"),
		filepath.Join(root, "server", "entity"),
		filepath.Join(root, "server", "entity", "effect"),
		filepath.Join(root, "server", "player"),
		filepath.Join(root, "server", "block"),
		filepath.Join(root, "server", "item", "enchantment"),
	); err != nil {
		t.Fatal(err)
	}
}
