package framework

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

func TestRegisterCustomItems(t *testing.T) {
	texture, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=")
	if err != nil {
		t.Fatal(err)
	}
	data, err := registerCustomItems([]native.CustomItemDefinition{{
		Identifier:        "framework:test_gem",
		Name:              "Framework Test Gem",
		TexturePNG:        texture,
		Category:          4,
		MaxCount:          16,
		ComponentDataJSON: `{"properties":{"foil":true},"components":{"minecraft:cooldown":{"duration":1.5,"category":"gem"}}}`,
	}})
	if err != nil {
		t.Fatal(err)
	}
	registered, ok := world.ItemByName("framework:test_gem", 0)
	if !ok {
		t.Fatal("custom item was not registered in Dragonfly")
	}
	custom, ok := registered.(world.CustomItem)
	if !ok || custom.Name() != "Framework Test Gem" || custom.Category().Uint8() != 4 {
		t.Fatalf("registered custom item = %#v", registered)
	}
	if counter, ok := registered.(item.MaxCounter); !ok || counter.MaxCount() != 16 {
		t.Fatalf("custom item maximum count = %#v", registered)
	}
	if foil, ok := data["framework:test_gem"].properties["foil"].(bool); !ok || !foil {
		t.Fatalf("custom item properties = %#v", data)
	}
	cooldown := data["framework:test_gem"].components["minecraft:cooldown"].(map[string]any)
	if cooldown["duration"] != float32(1.5) || cooldown["category"] != "gem" {
		t.Fatalf("custom item cooldown = %#v", cooldown)
	}
}

func TestRegisterCustomItemBehaviours(t *testing.T) {
	texture, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=")
	if err != nil {
		t.Fatal(err)
	}
	_, err = registerCustomItems([]native.CustomItemDefinition{
		{
			Identifier: "framework:component_food", Name: "Component Food", TexturePNG: texture, Category: 4, MaxCount: 16,
			ComponentDataJSON: `{"properties":{"use_duration":32,"use_animation":1},"components":{"minecraft:food":{"nutrition":6,"saturation_modifier":7.2,"can_always_eat":true},"minecraft:cooldown":{"duration":1,"category":"food"}}}`,
		},
		{
			Identifier: "framework:component_sword", Name: "Component Sword", TexturePNG: texture, Category: 3, MaxCount: 1,
			ComponentDataJSON: `{"properties":{"damage":8,"hand_equipped":true},"components":{"minecraft:durability":{"max_durability":850}}}`,
		},
		{
			Identifier: "framework:component_helmet", Name: "Component Helmet", TexturePNG: texture, Category: 3, MaxCount: 1,
			ComponentDataJSON: `{"properties":{},"components":{"minecraft:armor":{"protection":4,"toughness":3,"knockback_resistance":0.1},"minecraft:wearable":{"slot":"slot.armor.head","protection":4},"minecraft:durability":{"max_durability":600}}}`,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	food, _ := world.ItemByName("framework:component_food", 0)
	consumable, ok := food.(item.Consumable)
	if !ok || consumable.ConsumeDuration() != 1600*time.Millisecond || !consumable.AlwaysConsumable() {
		t.Fatalf("food behaviour = %#v", food)
	}
	if cooldown, ok := food.(item.Cooldown); !ok || cooldown.Cooldown() != time.Second {
		t.Fatalf("food cooldown = %#v", food)
	}
	sword, _ := world.ItemByName("framework:component_sword", 0)
	weapon, weaponOK := sword.(item.Weapon)
	durable, durableOK := sword.(item.Durable)
	if !weaponOK || weapon.AttackDamage() != 8 || !durableOK || durable.DurabilityInfo().MaxDurability != 850 {
		t.Fatalf("sword behaviour = %#v", sword)
	}
	helmet, _ := world.ItemByName("framework:component_helmet", 0)
	armour, armourOK := helmet.(item.Armour)
	_, helmetOK := helmet.(item.HelmetType)
	durable, durableOK = helmet.(item.Durable)
	if !armourOK || armour.DefencePoints() != 4 || !helmetOK || !durableOK || durable.DurabilityInfo().MaxDurability != 600 {
		t.Fatalf("helmet behaviour = %#v", helmet)
	}
}
