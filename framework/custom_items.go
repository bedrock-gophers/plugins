package framework

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/category"
	"github.com/df-mc/dragonfly/server/world"
)

type pluginCustomItem struct {
	identifier string
	name       string
	texture    image.Image
	category   category.Category
	maxCount   int
}

func (i pluginCustomItem) EncodeItem() (string, int16) { return i.identifier, 0 }
func (i pluginCustomItem) Name() string                { return i.name }
func (i pluginCustomItem) Texture() image.Image        { return i.texture }
func (i pluginCustomItem) Category() category.Category { return i.category }
func (i pluginCustomItem) MaxCount() int               { return i.maxCount }

type durablePluginCustomItem struct {
	pluginCustomItem
	durability item.DurabilityInfo
}

func (i durablePluginCustomItem) DurabilityInfo() item.DurabilityInfo {
	return i.durability
}

type weaponPluginCustomItem struct {
	pluginCustomItem
	damage float64
}

func (i weaponPluginCustomItem) AttackDamage() float64 { return i.damage }

type durableWeaponPluginCustomItem struct {
	durablePluginCustomItem
	damage float64
}

func (i durableWeaponPluginCustomItem) AttackDamage() float64 { return i.damage }

type offHandPluginCustomItem struct{ pluginCustomItem }

func (offHandPluginCustomItem) OffHand() bool { return true }

type armourPluginCustomItem struct {
	pluginCustomItem
	defence             float64
	toughness           float64
	knockBackResistance float64
}

func (i armourPluginCustomItem) DefencePoints() float64       { return i.defence }
func (i armourPluginCustomItem) Toughness() float64           { return i.toughness }
func (i armourPluginCustomItem) KnockBackResistance() float64 { return i.knockBackResistance }

type helmetPluginCustomItem struct{ armourPluginCustomItem }
type chestplatePluginCustomItem struct{ armourPluginCustomItem }
type leggingsPluginCustomItem struct{ armourPluginCustomItem }
type bootsPluginCustomItem struct{ armourPluginCustomItem }

func (helmetPluginCustomItem) Helmet() bool         { return true }
func (chestplatePluginCustomItem) Chestplate() bool { return true }
func (leggingsPluginCustomItem) Leggings() bool     { return true }
func (bootsPluginCustomItem) Boots() bool           { return true }
func (helmetPluginCustomItem) Use(_ *world.Tx, _ item.User, context *item.UseContext) bool {
	context.SwapHeldWithArmour(0)
	return false
}
func (chestplatePluginCustomItem) Use(_ *world.Tx, _ item.User, context *item.UseContext) bool {
	context.SwapHeldWithArmour(1)
	return false
}
func (leggingsPluginCustomItem) Use(_ *world.Tx, _ item.User, context *item.UseContext) bool {
	context.SwapHeldWithArmour(2)
	return false
}
func (bootsPluginCustomItem) Use(_ *world.Tx, _ item.User, context *item.UseContext) bool {
	context.SwapHeldWithArmour(3)
	return false
}

type durableHelmetPluginCustomItem struct {
	helmetPluginCustomItem
	durability item.DurabilityInfo
}
type durableChestplatePluginCustomItem struct {
	chestplatePluginCustomItem
	durability item.DurabilityInfo
}
type durableLeggingsPluginCustomItem struct {
	leggingsPluginCustomItem
	durability item.DurabilityInfo
}
type durableBootsPluginCustomItem struct {
	bootsPluginCustomItem
	durability item.DurabilityInfo
}

func (i durableHelmetPluginCustomItem) DurabilityInfo() item.DurabilityInfo {
	return i.durability
}
func (i durableChestplatePluginCustomItem) DurabilityInfo() item.DurabilityInfo {
	return i.durability
}
func (i durableLeggingsPluginCustomItem) DurabilityInfo() item.DurabilityInfo {
	return i.durability
}
func (i durableBootsPluginCustomItem) DurabilityInfo() item.DurabilityInfo {
	return i.durability
}

type foodPluginCustomItem struct {
	pluginCustomItem
	nutrition  int
	saturation float64
	always     bool
	duration   time.Duration
	drink      bool
}

func (i foodPluginCustomItem) AlwaysConsumable() bool         { return i.always }
func (i foodPluginCustomItem) ConsumeDuration() time.Duration { return i.duration }
func (i foodPluginCustomItem) Drinkable() bool                { return i.drink }
func (i foodPluginCustomItem) Consume(_ *world.Tx, consumer item.Consumer) item.Stack {
	consumer.Saturate(i.nutrition, i.saturation)
	return item.Stack{}
}

type cooldownFoodPluginCustomItem struct {
	foodPluginCustomItem
	cooldown time.Duration
}

func (i cooldownFoodPluginCustomItem) Cooldown() time.Duration { return i.cooldown }

type customItemData struct {
	properties map[string]any
	components map[string]any
	server     map[string]any
}

type customItemDataDocument struct {
	Properties map[string]any `json:"properties"`
	Components map[string]any `json:"components"`
	Server     map[string]any `json:"server"`
}

func registerCustomItems(definitions []native.CustomItemDefinition) (map[string]customItemData, error) {
	seen := make(map[string]struct{}, len(definitions))
	data := make(map[string]customItemData, len(definitions))
	for index, definition := range definitions {
		identifier := strings.ToLower(definition.Identifier)
		parts := strings.Split(identifier, ":")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" || strings.HasPrefix(identifier, "minecraft:") {
			return nil, fmt.Errorf("custom item %d has invalid identifier %q", index, definition.Identifier)
		}
		if definition.Name == "" || definition.MaxCount < 1 || definition.MaxCount > 64 {
			return nil, fmt.Errorf("custom item %q has invalid name or maximum count", identifier)
		}
		if _, exists := seen[identifier]; exists {
			return nil, fmt.Errorf("duplicate custom item %q", identifier)
		}
		seen[identifier] = struct{}{}
		texture, err := png.Decode(bytes.NewReader(definition.TexturePNG))
		if err != nil {
			return nil, fmt.Errorf("decode custom item %q texture: %w", identifier, err)
		}
		var itemCategory category.Category
		switch definition.Category {
		case 1:
			itemCategory = category.Construction()
		case 2:
			itemCategory = category.Nature()
		case 3:
			itemCategory = category.Equipment()
		case 4:
			itemCategory = category.Items()
		default:
			return nil, fmt.Errorf("custom item %q has invalid category %d", identifier, definition.Category)
		}
		itemData, err := decodeCustomItemData(definition.ComponentDataJSON)
		if err != nil {
			return nil, fmt.Errorf("decode custom item %q components: %w", identifier, err)
		}
		if definition.Group != "" {
			itemCategory = itemCategory.WithGroup(definition.Group)
		}
		base := pluginCustomItem{
			identifier: identifier,
			name:       definition.Name,
			texture:    texture,
			category:   itemCategory,
			maxCount:   int(definition.MaxCount),
		}
		world.RegisterItem(customItemWithBehaviours(base, itemData))
		if len(itemData.properties) != 0 || len(itemData.components) != 0 {
			data[identifier] = itemData
		}
	}
	return data, nil
}

func customItemWithBehaviours(base pluginCustomItem, data customItemData) world.Item {
	if armour, ok := componentMap(data.components, "minecraft:armor"); ok {
		wearable, _ := componentMap(data.components, "minecraft:wearable")
		value := armourPluginCustomItem{
			pluginCustomItem:    base,
			defence:             floatValue(armour["protection"]),
			toughness:           floatValue(data.server["armour_toughness"]),
			knockBackResistance: floatValue(data.server["armour_knockback_resistance"]),
		}
		durability, durable := componentMap(data.components, "minecraft:durability")
		durabilityInfo := customDurabilityInfo(intValue(durability["max_durability"]), data.server)
		switch wearable["slot"] {
		case "slot.armor.head":
			slotted := helmetPluginCustomItem{armourPluginCustomItem: value}
			if durable {
				return durableHelmetPluginCustomItem{helmetPluginCustomItem: slotted, durability: durabilityInfo}
			}
			return slotted
		case "slot.armor.chest":
			slotted := chestplatePluginCustomItem{armourPluginCustomItem: value}
			if durable {
				return durableChestplatePluginCustomItem{chestplatePluginCustomItem: slotted, durability: durabilityInfo}
			}
			return slotted
		case "slot.armor.legs":
			slotted := leggingsPluginCustomItem{armourPluginCustomItem: value}
			if durable {
				return durableLeggingsPluginCustomItem{leggingsPluginCustomItem: slotted, durability: durabilityInfo}
			}
			return slotted
		case "slot.armor.feet":
			slotted := bootsPluginCustomItem{armourPluginCustomItem: value}
			if durable {
				return durableBootsPluginCustomItem{bootsPluginCustomItem: slotted, durability: durabilityInfo}
			}
			return slotted
		}
	}
	if food, ok := componentMap(data.components, "minecraft:food"); ok {
		value := foodPluginCustomItem{
			pluginCustomItem: base,
			nutrition:        intValue(food["nutrition"]),
			saturation:       floatValue(food["saturation_modifier"]),
			always:           boolValue(food["can_always_eat"]),
			duration:         ticksDuration(data.properties["use_duration"]),
			drink:            intValue(data.properties["use_animation"]) == 2,
		}
		if value.duration <= 0 {
			value.duration = item.DefaultConsumeDuration
		}
		if cooldown, exists := componentMap(data.components, "minecraft:cooldown"); exists {
			return cooldownFoodPluginCustomItem{foodPluginCustomItem: value, cooldown: secondsDuration(cooldown["duration"])}
		}
		return value
	}
	durability, durable := componentMap(data.components, "minecraft:durability")
	damage, weapon := data.properties["damage"]
	if durable && weapon {
		return durableWeaponPluginCustomItem{
			durablePluginCustomItem: durablePluginCustomItem{pluginCustomItem: base, durability: customDurabilityInfo(intValue(durability["max_durability"]), data.server)},
			damage:                  floatValue(damage),
		}
	}
	if durable {
		return durablePluginCustomItem{pluginCustomItem: base, durability: customDurabilityInfo(intValue(durability["max_durability"]), data.server)}
	}
	if weapon {
		return weaponPluginCustomItem{pluginCustomItem: base, damage: floatValue(damage)}
	}
	if boolValue(data.properties["allow_off_hand"]) {
		return offHandPluginCustomItem{pluginCustomItem: base}
	}
	return base
}

func customDurabilityInfo(maximum int, server map[string]any) item.DurabilityInfo {
	attackDamage, breakDamage := 1, 1
	if value, ok := server["durability_attack_damage"]; ok {
		attackDamage = intValue(value)
	}
	if value, ok := server["durability_break_damage"]; ok {
		breakDamage = intValue(value)
	}
	return item.DurabilityInfo{
		MaxDurability:    maximum,
		BrokenItem:       func() item.Stack { return item.Stack{} },
		AttackDurability: attackDamage,
		BreakDurability:  breakDamage,
		Persistent:       boolValue(server["durability_persistent"]),
	}
}

func componentMap(components map[string]any, name string) (map[string]any, bool) {
	value, ok := components[name].(map[string]any)
	return value, ok
}

func intValue(value any) int {
	switch value := value.(type) {
	case int32:
		return int(value)
	case int64:
		return int(value)
	case float32:
		return int(value)
	default:
		return 0
	}
}

func floatValue(value any) float64 {
	switch value := value.(type) {
	case int32:
		return float64(value)
	case int64:
		return float64(value)
	case float32:
		return float64(value)
	default:
		return 0
	}
}

func boolValue(value any) bool {
	result, _ := value.(bool)
	return result
}

func ticksDuration(value any) time.Duration {
	return time.Duration(floatValue(value) * float64(time.Second) / 20)
}

func secondsDuration(value any) time.Duration {
	return time.Duration(floatValue(value) * float64(time.Second))
}

func decodeCustomItemData(value string) (customItemData, error) {
	if value == "" {
		return customItemData{}, nil
	}
	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.UseNumber()
	decoder.DisallowUnknownFields()
	var document customItemDataDocument
	if err := decoder.Decode(&document); err != nil {
		return customItemData{}, err
	}
	if err := ensureJSONEnd(decoder); err != nil {
		return customItemData{}, err
	}
	for name := range document.Properties {
		if name == "minecraft:icon" || name == "creative_group" || name == "creative_category" || name == "max_stack_size" {
			return customItemData{}, fmt.Errorf("property %q is managed by the custom item definition", name)
		}
	}
	for name := range document.Components {
		if name == "item_properties" || name == "minecraft:display_name" {
			return customItemData{}, fmt.Errorf("component %q is managed by the custom item definition", name)
		}
	}
	properties, err := normaliseComponentMap(document.Properties)
	if err != nil {
		return customItemData{}, fmt.Errorf("properties: %w", err)
	}
	components, err := normaliseComponentMap(document.Components)
	if err != nil {
		return customItemData{}, fmt.Errorf("components: %w", err)
	}
	server, err := normaliseComponentMap(document.Server)
	if err != nil {
		return customItemData{}, fmt.Errorf("server behaviours: %w", err)
	}
	return customItemData{properties: properties, components: components, server: server}, nil
}

func ensureJSONEnd(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("component data contains multiple JSON values")
		}
		return err
	}
	return nil
}

func normaliseComponentMap(values map[string]any) (map[string]any, error) {
	if len(values) == 0 {
		return nil, nil
	}
	result := make(map[string]any, len(values))
	for key, value := range values {
		if key == "" || len(key) > 128 {
			return nil, fmt.Errorf("invalid key %q", key)
		}
		normalised, err := normaliseComponentValue(value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}
		result[key] = normalised
	}
	return result, nil
}

func normaliseComponentValue(value any) (any, error) {
	switch value := value.(type) {
	case nil:
		return nil, errors.New("null values are not supported")
	case bool, string:
		return value, nil
	case json.Number:
		if integer, err := strconv.ParseInt(string(value), 10, 64); err == nil {
			if integer >= -1<<31 && integer <= 1<<31-1 {
				return int32(integer), nil
			}
			return integer, nil
		}
		decimal, err := strconv.ParseFloat(string(value), 32)
		if err != nil || math.IsInf(decimal, 0) || math.IsNaN(decimal) {
			return nil, fmt.Errorf("invalid number %q", value)
		}
		return float32(decimal), nil
	case map[string]any:
		return normaliseComponentMap(value)
	case []any:
		result := make([]any, len(value))
		for index, entry := range value {
			normalised, err := normaliseComponentValue(entry)
			if err != nil {
				return nil, fmt.Errorf("index %d: %w", index, err)
			}
			result[index] = normalised
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported value type %T", value)
	}
}
