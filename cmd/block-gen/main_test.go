package main

import (
	"strings"
	"testing"

	dfblock "github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
)

func TestGenerateTypedBlocks(t *testing.T) {
	states := []blockState{
		{name: "minecraft:sand"},
		{name: "minecraft:oak_log", properties: map[string]any{"pillar_axis": "x"}},
		{name: "minecraft:oak_log", properties: map[string]any{"pillar_axis": "y"}},
		{name: "minecraft:oak_log", properties: map[string]any{"pillar_axis": "z"}},
		{name: "minecraft:oak_door", properties: map[string]any{"direction": int32(0), "open_bit": uint8(0)}},
		{name: "minecraft:oak_door", properties: map[string]any{"direction": int32(1), "open_bit": uint8(1)}},
	}

	generated, err := generate(states)
	if err != nil {
		t.Fatal(err)
	}
	source := string(generated)
	for _, expected := range []string{
		"pub struct Sand;",
		"impl From<Sand> for Block",
		"fn from(_value: Sand)",
		"pub enum PillarAxis",
		"pub struct OakLog",
		"pub const fn new(pillar_axis: PillarAxis) -> Self",
		"pub const fn new(direction: Direction, open_bit: bool) -> Self",
		"Property::Uint8(u8::from(value.open_bit))",
	} {
		if !strings.Contains(source, expected) {
			t.Errorf("generated source missing %q", expected)
		}
	}
}

func TestGenerateRejectsMixedPropertyKinds(t *testing.T) {
	_, err := generate([]blockState{
		{name: "minecraft:test", properties: map[string]any{"state": "one"}},
		{name: "minecraft:test", properties: map[string]any{"state": int32(1)}},
	})
	if err == nil || !strings.Contains(err.Error(), "mixed property types") {
		t.Fatalf("expected mixed property type error, got %v", err)
	}
}

func TestGenerateNormalisesBooleanByteProperties(t *testing.T) {
	generated, err := generate([]blockState{
		{name: "minecraft:tnt", properties: map[string]any{"explode_bit": uint8(0)}},
		{name: "minecraft:tnt", properties: map[string]any{"explode_bit": true}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if source := string(generated); !strings.Contains(source, "pub const fn new(explode_bit: bool)") {
		t.Fatalf("generated source did not normalise boolean byte property:\n%s", source)
	}
}

func TestGenerateFiniteNumericStateEnum(t *testing.T) {
	generated, err := generate([]blockState{
		{name: "minecraft:piston", properties: map[string]any{"facing_direction": int32(0)}},
		{name: "minecraft:piston", properties: map[string]any{"facing_direction": int32(1)}},
		{name: "minecraft:piston", properties: map[string]any{"facing_direction": int32(2)}},
	})
	if err != nil {
		t.Fatal(err)
	}
	source := string(generated)
	for _, expected := range []string{
		"pub enum FacingDirection",
		"Value0 = 0",
		"Value2 = 2",
		"pub const fn new(facing_direction: FacingDirection)",
		"Property::Int32(value.facing_direction.as_i32())",
	} {
		if !strings.Contains(source, expected) {
			t.Fatalf("generated source missing %q:\n%s", expected, source)
		}
	}
}

func TestRustNameHandlesNamespacesAndDigits(t *testing.T) {
	if got := rustTypeName("example:3d_block"); got != "Example3dBlock" {
		t.Fatalf("rustTypeName() = %q", got)
	}
	if got := rustTypeName("minecraft:3d_block"); got != "Block3dBlock" {
		t.Fatalf("rustTypeName() = %q", got)
	}
	if got := rustTypeName("minecraft:self"); got != "BlockSelf" {
		t.Fatalf("rustTypeName() = %q", got)
	}
	if got := rustTypeName("minecraft:block"); got != "BuiltinBlock" {
		t.Fatalf("rustTypeName() = %q", got)
	}
	if got := rustFieldName("3d-state"); got != "state_3d_state" {
		t.Fatalf("rustFieldName() = %q", got)
	}
	if got := rustFieldName("self"); got != "state_self" {
		t.Fatalf("rustFieldName() = %q", got)
	}
}

func TestEnumNameDoesNotCollideWithBlockType(t *testing.T) {
	generated, err := generate([]blockState{
		{name: "minecraft:pillar_axis"},
		{name: "minecraft:oak_log", properties: map[string]any{"pillar_axis": "x"}},
		{name: "minecraft:oak_log", properties: map[string]any{"pillar_axis": "y"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	source := string(generated)
	for _, expected := range []string{"pub struct PillarAxis;", "pub enum PillarAxisState"} {
		if !strings.Contains(source, expected) {
			t.Fatalf("generated source missing %q:\n%s", expected, source)
		}
	}
}

func TestNamespacedPropertyUsesRustFieldName(t *testing.T) {
	generated, err := generate([]blockState{
		{name: "minecraft:door", properties: map[string]any{"minecraft:cardinal_direction": "north"}},
		{name: "minecraft:door", properties: map[string]any{"minecraft:cardinal_direction": "south"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	source := string(generated)
	if strings.Contains(source, "minecraft:cardinal_direction:") {
		t.Fatalf("generated invalid Rust field:\n%s", source)
	}
	for _, expected := range []string{"pub enum CardinalDirection", "cardinal_direction: CardinalDirection"} {
		if !strings.Contains(source, expected) {
			t.Fatalf("generated source missing %q:\n%s", expected, source)
		}
	}
}

func TestStatefulBlockWithoutSemanticDefaultRequiresExplicitState(t *testing.T) {
	generated, err := generate([]blockState{
		{name: "minecraft:door", properties: map[string]any{"direction": int32(0), "hinge": uint8(0), "open": uint8(0)}},
		{name: "minecraft:door", properties: map[string]any{"direction": int32(1), "hinge": uint8(1), "open": uint8(1)}},
	})
	if err != nil {
		t.Fatal(err)
	}
	source := string(generated)
	for _, expected := range []string{
		"pub const fn new(direction: Direction, hinge: bool, open: bool) -> Self",
		"pub const fn with_direction(mut self, direction: Direction) -> Self",
		"pub const fn with_open(mut self, open: bool) -> Self",
	} {
		if !strings.Contains(source, expected) {
			t.Fatalf("generated source missing %q:\n%s", expected, source)
		}
	}
	if strings.Contains(source, "impl Default for Door") || strings.Contains(source, "pub const DEFAULT") {
		t.Fatalf("generated an invented default state:\n%s", source)
	}
}

func TestSemanticDefaultUsesDragonflyZeroState(t *testing.T) {
	generated, err := generate([]blockState{
		{name: "minecraft:oak_log", properties: map[string]any{"pillar_axis": "x"}},
		{name: "minecraft:oak_log", properties: map[string]any{"pillar_axis": "y"}, semanticDefault: true},
		{name: "minecraft:oak_log", properties: map[string]any{"pillar_axis": "z"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	source := string(generated)
	for _, expected := range []string{
		"pub const DEFAULT: Self = Self { pillar_axis: PillarAxis::Y }",
		"impl Default for OakLog",
		"pub const fn new(pillar_axis: PillarAxis) -> Self",
	} {
		if !strings.Contains(source, expected) {
			t.Fatalf("generated source missing %q:\n%s", expected, source)
		}
	}
}

func TestSemanticZeroStateUsesDragonflyEncoding(t *testing.T) {
	state, ok := semanticZeroState(dfblock.Log{})
	if !ok {
		t.Fatal("Dragonfly Log zero state was not encodable")
	}
	if state.name != "minecraft:oak_log" || state.properties["pillar_axis"] != "y" {
		t.Fatalf("semanticZeroState(Log{}) = %#v", state)
	}
}

func TestRegistryNumericStatesBeyondTypedLimitAreExplicit(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	registered := world.DefaultBlockRegistry.Blocks()
	states := make([]blockState, 0, len(registered))
	for _, block := range registered {
		name, properties := block.EncodeBlock()
		states = append(states, blockState{name: name, properties: properties})
	}
	blocks, err := inspect(states)
	if err != nil {
		t.Fatal(err)
	}
	deferred := map[string]int{}
	for _, block := range blocks {
		for _, property := range block.properties {
			if property.variable && (property.kind == propertyInt32 || property.kind == propertyUint8) &&
				!isBinary(property.values) && !hasGeneratedStateType(property) {
				deferred[property.name] = len(property.values)
			}
		}
	}
	if len(deferred) != 0 {
		t.Fatalf("untyped numeric registry states: %#v", deferred)
	}
}
