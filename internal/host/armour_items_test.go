package host

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/item"
)

func TestArmourItemsToNative(t *testing.T) {
	dye := color.RGBA{R: 0x12, G: 0x34, B: 0x56, A: 0xff}
	trim := item.ArmourTrim{Template: item.TemplateSentry(), Material: item.AmethystShard{}}
	tests := []struct {
		name       string
		identifier string
		armour     item.Stack
	}{
		{"helmet", "minecraft:leather_helmet", item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{Colour: dye}, Trim: trim}, 1)},
		{"chestplate", "minecraft:leather_chestplate", item.NewStack(item.Chestplate{Tier: item.ArmourTierLeather{Colour: dye}, Trim: trim}, 1)},
		{"leggings", "minecraft:leather_leggings", item.NewStack(item.Leggings{Tier: item.ArmourTierLeather{Colour: dye}, Trim: trim}, 1)},
		{"boots", "minecraft:leather_boots", item.NewStack(item.Boots{Tier: item.ArmourTierLeather{Colour: dye}, Trim: trim}, 1)},
	}
	wantNBT := map[string]any{
		"customColor": armourCustomColour(dye),
		"Trim": map[string]any{
			"Material": "amethyst",
			"Pattern":  "sentry",
		},
	}
	if wantNBT["customColor"].(int32) >= 0 {
		t.Fatal("test colour must exercise signed ARGB")
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := itemStackToNative(test.armour)
			if !ok {
				t.Fatal("encode armour")
			}
			if got.Identifier != test.identifier || got.Metadata != 0 || got.Count != 1 {
				t.Fatalf("armour stack = %#v", got)
			}
			assertArmourNBT(t, got.NBT, wantNBT)
		})
	}
}

func TestArmourItemsFromNative(t *testing.T) {
	dye := color.RGBA{R: 0x12, G: 0x34, B: 0x56, A: 0xff}
	stack, ok := itemStackFromNative(native.ItemStack{
		Identifier: "minecraft:leather_helmet",
		Count:      1,
		NBT: mustMarshalArmourNBT(t, map[string]any{
			"customColor": armourCustomColour(dye),
			"Trim": map[string]any{
				"Material": "amethyst",
				"Pattern":  "sentry",
			},
		}),
	})
	if !ok {
		t.Fatal("decode leather helmet")
	}
	helmet, ok := stack.Item().(item.Helmet)
	if !ok {
		t.Fatalf("helmet item type = %T", stack.Item())
	}
	leather, ok := helmet.Tier.(item.ArmourTierLeather)
	if !ok || leather.Colour != dye {
		t.Fatalf("helmet tier = %#v", helmet.Tier)
	}
	if helmet.Trim.Template != item.TemplateSentry() {
		t.Fatalf("helmet trim template = %#v", helmet.Trim.Template)
	}
	if _, ok := helmet.Trim.Material.(item.AmethystShard); !ok {
		t.Fatalf("helmet trim material = %T", helmet.Trim.Material)
	}
	if stack.Count() != 1 {
		t.Fatalf("helmet count = %d", stack.Count())
	}
}

func TestArmourItemsRejectInvalidTrim(t *testing.T) {
	for _, test := range []struct {
		name     string
		material string
		pattern  string
	}{
		{"material", "missing", "sentry"},
		{"pattern", "amethyst", "missing"},
	} {
		t.Run(test.name, func(t *testing.T) {
			stack, ok := itemStackFromNative(native.ItemStack{
				Identifier: "minecraft:leather_helmet",
				Count:      1,
				NBT: mustMarshalArmourNBT(t, map[string]any{
					"Trim": map[string]any{
						"Material": test.material,
						"Pattern":  test.pattern,
					},
				}),
			})
			if !ok {
				t.Fatal("decode armour")
			}
			helmet := stack.Item().(item.Helmet)
			if !reflect.DeepEqual(helmet.Trim, item.ArmourTrim{}) {
				t.Fatalf("invalid trim = %#v", helmet.Trim)
			}
		})
	}
}

func TestArmourItemsNonLeatherIgnoreCustomColour(t *testing.T) {
	stack, ok := itemStackFromNative(native.ItemStack{
		Identifier: "minecraft:diamond_helmet",
		Count:      1,
		NBT: mustMarshalArmourNBT(t, map[string]any{
			"customColor": int32(-1),
			"Trim": map[string]any{
				"Material": "gold",
				"Pattern":  "ward",
			},
		}),
	})
	if !ok {
		t.Fatal("decode diamond helmet")
	}
	helmet := stack.Item().(item.Helmet)
	if _, ok := helmet.Tier.(item.ArmourTierDiamond); !ok {
		t.Fatalf("helmet tier = %T", helmet.Tier)
	}
	if helmet.Trim.Template != item.TemplateWard() {
		t.Fatalf("helmet trim template = %#v", helmet.Trim.Template)
	}
	if _, ok := helmet.Trim.Material.(item.GoldIngot); !ok {
		t.Fatalf("helmet trim material = %T", helmet.Trim.Material)
	}

	encoded, ok := itemStackToNative(stack)
	if !ok {
		t.Fatal("re-encode diamond helmet")
	}
	assertArmourNBT(t, encoded.NBT, map[string]any{
		"Trim": map[string]any{
			"Material": "gold",
			"Pattern":  "ward",
		},
	})
}

func TestArmourItemsNetheriteUpgradeIsZeroTrim(t *testing.T) {
	stack, ok := itemStackFromNative(native.ItemStack{
		Identifier: "minecraft:iron_helmet",
		Count:      1,
		NBT: mustMarshalArmourNBT(t, map[string]any{
			"Trim": map[string]any{
				"Material": "diamond",
				"Pattern":  "netherite_upgrade",
			},
		}),
	})
	if !ok {
		t.Fatal("decode netherite-upgrade trim")
	}
	helmet := stack.Item().(item.Helmet)
	if !helmet.Trim.Zero() || helmet.Trim.Template != item.TemplateNetheriteUpgrade() {
		t.Fatalf("netherite-upgrade trim = %#v", helmet.Trim)
	}

	encoded, ok := itemStackToNative(stack)
	if !ok {
		t.Fatal("re-encode netherite-upgrade trim")
	}
	assertArmourNBT(t, encoded.NBT, map[string]any{})
}

func armourCustomColour(value color.RGBA) int32 {
	if value.R == 0 && value.G == 0 && value.B == 0 {
		return -0x1000000
	}
	return int32(uint32(value.A)<<24 | uint32(value.R)<<16 | uint32(value.G)<<8 | uint32(value.B))
}

func mustMarshalArmourNBT(t *testing.T, value map[string]any) []byte {
	t.Helper()
	encoded, ok := marshalItemNBT(value)
	if !ok {
		t.Fatal("marshal armour NBT")
	}
	return encoded
}

func assertArmourNBT(t *testing.T, encoded []byte, want map[string]any) {
	t.Helper()
	got, ok := unmarshalItemNBT(encoded)
	if !ok {
		t.Fatal("unmarshal armour NBT")
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("armour NBT = %#v, want %#v", got, want)
	}
}
