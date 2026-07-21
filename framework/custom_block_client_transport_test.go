package framework

import (
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestApplyCustomBlockClientDataConvertsDragonflyComponents(t *testing.T) {
	components := map[string]any{
		"minecraft:block_light_emission": map[string]any{"emission": float32(1)},
		"minecraft:block_light_filter":   map[string]any{"lightLevel": int32(0)},
		"minecraft:friction":             map[string]any{"value": float32(0.6)},
		"minecraft:unit_cube":            map[string]any{},
	}
	entries := []protocol.BlockEntry{{
		Name:       "test:lamp",
		Properties: map[string]any{"components": components},
	}}

	applyCustomBlockClientData(entries, map[string]customBlockClientData{
		"test:lamp": {emission: 15, dampening: 0, friction: 0.6},
	})

	if _, ok := components["minecraft:block_light_emission"]; ok {
		t.Fatal("legacy light emission component was not removed")
	}
	if _, ok := components["minecraft:block_light_filter"]; ok {
		t.Fatal("legacy light filter component was not removed")
	}
	emission := components["minecraft:light_emission"].(map[string]any)
	if got := emission["emission"]; got != uint8(15) {
		t.Fatalf("light emission = %v, want 15", got)
	}
	dampening := components["minecraft:light_dampening"].(map[string]any)
	if got := dampening["lightLevel"]; got != uint8(0) {
		t.Fatalf("light dampening = %v, want 0", got)
	}
	friction := components["minecraft:friction"].(map[string]any)
	if got := friction["value"]; got != float32(0.4) {
		t.Fatalf("friction = %v, want 0.4", got)
	}
	if _, ok := components["minecraft:unit_cube"]; !ok {
		t.Fatal("unrelated native component was modified")
	}
	encoded, err := nbt.MarshalEncoding(entries[0].Properties, nbt.NetworkLittleEndian)
	if err != nil {
		t.Fatalf("marshal block properties as network NBT: %v", err)
	}
	var decoded map[string]any
	if err := nbt.UnmarshalEncoding(encoded, &decoded, nbt.NetworkLittleEndian); err != nil {
		t.Fatalf("unmarshal block properties as network NBT: %v", err)
	}
	decodedComponents := decoded["components"].(map[string]any)
	decodedEmission := decodedComponents["minecraft:light_emission"].(map[string]any)
	if got := decodedEmission["emission"]; got != uint8(15) {
		t.Fatalf("network NBT emission = %v (%T), want byte 15", got, got)
	}
}

func TestApplyCustomBlockClientDataOmitsDefaultLightValues(t *testing.T) {
	components := map[string]any{
		"minecraft:block_light_emission": map[string]any{"emission": float32(0)},
		"minecraft:block_light_filter":   map[string]any{"lightLevel": int32(15)},
	}
	entries := []protocol.BlockEntry{{
		Name:       "test:solid",
		Properties: map[string]any{"components": components},
	}}

	applyCustomBlockClientData(entries, map[string]customBlockClientData{
		"test:solid": {emission: 0, dampening: 15, friction: 0.6},
	})

	if _, ok := components["minecraft:light_emission"]; ok {
		t.Fatal("default light emission should be omitted")
	}
	if _, ok := components["minecraft:light_dampening"]; ok {
		t.Fatal("default light dampening should be omitted")
	}
}

func TestApplyCustomBlockClientDataConvertsSlipperyFriction(t *testing.T) {
	components := map[string]any{}
	entries := []protocol.BlockEntry{{
		Name:       "test:ice",
		Properties: map[string]any{"components": components},
	}}

	applyCustomBlockClientData(entries, map[string]customBlockClientData{
		"test:ice": {dampening: 15, friction: 0.98},
	})

	friction := components["minecraft:friction"].(map[string]any)
	if got := friction["value"].(float32); got < 0.019 || got > 0.021 {
		t.Fatalf("friction = %v, want approximately 0.02", got)
	}
}
