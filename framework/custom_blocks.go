package framework

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"maps"
	"slices"
	"strings"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/customblock"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/category"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type customBlockDataDocument struct {
	Server       map[string]any                   `json:"server"`
	States       map[string][]any                 `json:"states"`
	Permutations []customBlockPermutationDocument `json:"permutations"`
}

type customBlockPermutationDocument struct {
	Condition  string         `json:"condition"`
	Properties map[string]any `json:"properties"`
}

type decodedCustomBlockData struct {
	server       map[string]any
	states       map[string][]any
	permutations []customBlockPermutation
}

type customBlockPermutation struct {
	condition  string
	properties map[string]any
}

type pluginCustomBlock struct {
	identifier, name                    string
	textureKey                          string
	texture                             image.Image
	geometry                            []byte
	category                            category.Category
	maxCount                            int
	hash                                uint64
	stateHash                           uint64
	state                               map[string]any
	properties                          customblock.Properties
	collision                           cube.BBox
	solidFaces                          bool
	hardness, blastResistance, friction float64
	lightEmission, lightDiffusion       uint8
	encouragement, flammability         int
	lavaFlammable, replaceable          bool
	allowedFaces                        map[cube.Face]bool
	placementFilter                     map[string]bool
	displaceWater, displaceLava         bool
	states                              map[string][]any
	permutations                        []customblock.Permutation
}

func (b pluginCustomBlock) EncodeBlock() (string, map[string]any) { return b.identifier, b.state }
func (b pluginCustomBlock) EncodeItem() (string, int16)           { return b.identifier, 0 }
func (b pluginCustomBlock) Hash() (uint64, uint64)                { return b.hash, b.stateHash }
func (b pluginCustomBlock) Model() world.BlockModel {
	return customBlockModel{box: b.collision, solidFaces: b.solidFaces}
}
func (b pluginCustomBlock) Properties() customblock.Properties { return b.properties }
func (b pluginCustomBlock) States() map[string][]any           { return b.states }
func (b pluginCustomBlock) Permutations() []customblock.Permutation {
	return b.permutations
}
func (b pluginCustomBlock) Name() string     { return b.name }
func (b pluginCustomBlock) Geometry() []byte { return b.geometry }
func (b pluginCustomBlock) Textures() map[string]image.Image {
	return map[string]image.Image{b.textureKey: b.texture}
}
func (b pluginCustomBlock) Texture() image.Image           { return b.texture }
func (b pluginCustomBlock) Category() category.Category    { return b.category }
func (b pluginCustomBlock) MaxCount() int                  { return b.maxCount }
func (b pluginCustomBlock) LightEmissionLevel() uint8      { return b.lightEmission }
func (b pluginCustomBlock) LightDiffusionLevel() uint8     { return b.lightDiffusion }
func (b pluginCustomBlock) Friction() float64              { return b.friction }
func (b pluginCustomBlock) ReplaceableBy(world.Block) bool { return b.replaceable }
func (b pluginCustomBlock) CanDisplace(liquid world.Liquid) bool {
	name, _ := liquid.EncodeBlock()
	return name == "minecraft:water" && b.displaceWater || name == "minecraft:lava" && b.displaceLava
}
func (b pluginCustomBlock) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool { return b.solidFaces }
func (b pluginCustomBlock) FlammabilityInfo() block.FlammabilityInfo {
	return block.FlammabilityInfo{Encouragement: b.encouragement, Flammability: b.flammability, LavaFlammable: b.lavaFlammable}
}
func (b pluginCustomBlock) BreakInfo() block.BreakInfo {
	return block.BreakInfo{
		Hardness: b.hardness, BlastResistance: b.blastResistance,
		Harvestable: func(item.Tool) bool { return true }, Effective: func(item.Tool) bool { return false },
		Drops: func(item.Tool, []item.Enchantment) []item.Stack { return []item.Stack{item.NewStack(b, 1)} },
	}
}
func (b pluginCustomBlock) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	if len(b.allowedFaces) != 0 && !b.allowedFaces[face] {
		return false
	}
	if len(b.placementFilter) != 0 {
		name, _ := tx.Block(pos).EncodeBlock()
		if !b.placementFilter[name] {
			return false
		}
	}
	target := pos.Side(face)
	if replaceable, ok := tx.Block(pos).(block.Replaceable); ok && replaceable.ReplaceableBy(b) {
		target = pos
	}
	if current := tx.Block(target); current != nil {
		if replaceable, ok := current.(block.Replaceable); !ok || !replaceable.ReplaceableBy(b) {
			return false
		}
	}
	if placer, ok := user.(block.Placer); ok {
		placer.PlaceBlock(target, b, ctx)
	} else {
		tx.SetBlock(target, b, nil)
		tx.PlaySound(target.Vec3(), sound.BlockPlace{Block: b})
		ctx.SubtractFromCount(1)
	}
	return true
}

type customBlockModel struct {
	box        cube.BBox
	solidFaces bool
}

func (m customBlockModel) BBox(cube.Pos, world.BlockSource) []cube.BBox          { return []cube.BBox{m.box} }
func (m customBlockModel) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool { return m.solidFaces }

type customBlockClientData struct {
	emission, dampening uint8
	friction            float64
}

func registerCustomBlocks(definitions []native.CustomBlockDefinition) (map[string]customBlockClientData, error) {
	seen := make(map[string]struct{}, len(definitions))
	clientData := make(map[string]customBlockClientData, len(definitions))
	for index, definition := range definitions {
		identifier := strings.ToLower(definition.Identifier)
		parts := strings.Split(identifier, ":")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" || parts[0] == "minecraft" {
			return nil, fmt.Errorf("custom block %d has invalid identifier %q", index, definition.Identifier)
		}
		if _, exists := seen[identifier]; exists {
			return nil, fmt.Errorf("duplicate custom block %q", identifier)
		}
		if definition.Name == "" || definition.MaxCount < 1 || definition.MaxCount > 64 {
			return nil, fmt.Errorf("custom block %q has invalid name or maximum count", identifier)
		}
		seen[identifier] = struct{}{}
		texture, err := png.Decode(bytes.NewReader(definition.TexturePNG))
		if err != nil {
			return nil, fmt.Errorf("decode custom block %q texture: %w", identifier, err)
		}
		data, err := decodeCustomBlockData(definition.ComponentDataJSON)
		if err != nil {
			return nil, fmt.Errorf("decode custom block %q data: %w", identifier, err)
		}
		itemCategory, err := customBlockCategory(definition.Category)
		if err != nil {
			return nil, fmt.Errorf("custom block %q: %w", identifier, err)
		}
		if definition.Group != "" {
			itemCategory = itemCategory.WithGroup(definition.Group)
		}
		box := blockBox(data.server["collision_box"], cube.Box(0, 0, 0, 1, 1, 1))
		geometryName := ""
		cubeGeometry := len(definition.GeometryJSON) == 0
		geometry := customBlockGeometry(definition.GeometryJSON)
		if !cubeGeometry {
			geometryName = stringValue(data.server["geometry_identifier"])
			if geometryName == "" {
				geometryName = "geometry." + strings.ReplaceAll(identifier, ":", ".")
			}
		}
		value := pluginCustomBlock{
			identifier: identifier, name: definition.Name, texture: texture, geometry: geometry,
			textureKey: strings.NewReplacer(":", "_", "/", "_").Replace(identifier),
			category:   itemCategory, maxCount: int(definition.MaxCount), hash: block.NextHash(), collision: box,
			solidFaces: boolValueDefault(data.server["solid_faces"], cubeGeometry), hardness: floatValueDefault(data.server["hardness"], 1.5),
			blastResistance: floatValueDefault(data.server["blast_resistance"], 7.5), friction: floatValueDefault(data.server["friction"], .6),
			lightEmission: uint8Value(data.server["light_emission"], 0), lightDiffusion: uint8Value(data.server["light_dampening"], 15),
			encouragement: intValue(data.server["fire_encouragement"]), flammability: intValue(data.server["flammability"]),
			lavaFlammable: boolValue(data.server["lava_flammable"]), replaceable: boolValue(data.server["replaceable"]),
			states: data.states,
		}
		value.permutations = make([]customblock.Permutation, 0, len(data.permutations))
		for _, permutation := range data.permutations {
			value.permutations = append(value.permutations, customblock.Permutation{
				Condition: permutation.condition,
				Properties: customBlockProperties(permutation.properties, value.textureKey, false,
					stringValue(permutation.properties["geometry_identifier"]), false),
			})
		}
		value.allowedFaces, value.placementFilter = customBlockPlacementFilter(data.server)
		value.displaceWater = boolValue(data.server["displace_water"])
		value.displaceLava = boolValue(data.server["displace_lava"])
		value.properties = customBlockProperties(data.server, value.textureKey, cubeGeometry, geometryName, true)
		clientData[identifier] = customBlockClientData{
			emission: value.lightEmission, dampening: value.lightDiffusion, friction: value.friction,
		}
		states, err := customBlockStates(data.states)
		if err != nil {
			return nil, fmt.Errorf("custom block %q states: %w", identifier, err)
		}
		value.state, value.stateHash = states[0], 0
		for stateIndex, state := range states {
			variant := value
			variant.state, variant.stateHash = state, uint64(stateIndex)
			world.RegisterBlock(variant)
		}
		world.RegisterItem(value)
	}
	return clientData, nil
}

func customBlockCategory(value uint32) (category.Category, error) {
	switch value {
	case 1:
		return category.Construction(), nil
	case 2:
		return category.Nature(), nil
	case 3:
		return category.Equipment(), nil
	case 4:
		return category.Items(), nil
	}
	return category.Category{}, fmt.Errorf("invalid category %d", value)
}

func customBlockGeometry(value []byte) []byte {
	if len(value) == 0 {
		return nil
	}
	return bytes.Clone(value)
}

func customBlockProperties(values map[string]any, textureKey string, cubeGeometry bool, geometryName string, base bool) customblock.Properties {
	properties := customblock.Properties{
		Cube:        cubeGeometry,
		Geometry:    geometryName,
		MapColour:   stringValue(values["map_color"]),
		Scale:       blockVector(values["scale"], mgl64.Vec3{}),
		Translation: blockVector(values["translation"], mgl64.Vec3{}),
		Rotation:    blockRotation(values["rotation"]),
	}
	if base {
		properties.CollisionBox = blockBox(values["collision_box"], cube.Box(0, 0, 0, 1, 1, 1))
		properties.SelectionBox = blockBox(values["selection_box"], properties.CollisionBox)
	} else {
		properties.CollisionBox = blockBox(values["collision_box"], cube.BBox{})
		properties.SelectionBox = blockBox(values["selection_box"], cube.BBox{})
	}
	if base || values["render_method"] != nil {
		material := customblock.NewMaterial(textureKey, customBlockRenderMethod(values["render_method"]))
		if faceDimming, ok := values["face_dimming"].(bool); ok && !faceDimming {
			material = material.WithoutFaceDimming()
		}
		if ambientOcclusion, ok := values["ambient_occlusion"].(bool); ok {
			if ambientOcclusion {
				material = material.WithAmbientOcclusion()
			} else {
				material = material.WithoutAmbientOcclusion()
			}
		}
		properties.Textures = make(map[string]customblock.Material)
		targets, _ := values["texture_targets"].([]any)
		if len(targets) == 0 {
			properties.Textures["*"] = material
		} else {
			for _, target := range targets {
				if name := stringValue(target); name != "" {
					properties.Textures[name] = material
				}
			}
		}
	}
	return properties
}

func customBlockRenderMethod(value any) customblock.Method {
	switch intValue(value) {
	case 1:
		return customblock.AlphaTestRenderMethod()
	case 2:
		return customblock.BlendRenderMethod()
	case 3:
		return customblock.DoubleSidedRenderMethod()
	default:
		return customblock.OpaqueRenderMethod()
	}
}

func blockBox(value any, fallback cube.BBox) cube.BBox {
	v, ok := value.(map[string]any)
	if !ok {
		return fallback
	}
	return cube.Box(floatValue(v["min_x"]), floatValue(v["min_y"]), floatValue(v["min_z"]), floatValue(v["max_x"]), floatValue(v["max_y"]), floatValue(v["max_z"]))
}
func blockVector(value any, fallback mgl64.Vec3) mgl64.Vec3 {
	v, ok := value.(map[string]any)
	if !ok {
		return fallback
	}
	return mgl64.Vec3{floatValue(v["x"]), floatValue(v["y"]), floatValue(v["z"])}
}
func blockRotation(value any) cube.Pos {
	v, ok := value.(map[string]any)
	if !ok {
		return cube.Pos{}
	}
	return cube.Pos{intValue(v["x"]), intValue(v["y"]), intValue(v["z"])}
}
func stringValue(value any) string { v, _ := value.(string); return v }
func floatValueDefault(value any, fallback float64) float64 {
	if value == nil {
		return fallback
	}
	return floatValue(value)
}
func boolValueDefault(value any, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return boolValue(value)
}
func uint8Value(value any, fallback uint8) uint8 {
	if value == nil {
		return fallback
	}
	v := intValue(value)
	if v < 0 {
		return 0
	}
	if v > 15 {
		return 15
	}
	return uint8(v)
}

func customBlockPlacementFilter(server map[string]any) (map[cube.Face]bool, map[string]bool) {
	faces := map[cube.Face]bool{}
	if values, ok := server["allowed_faces"].([]any); ok {
		for _, value := range values {
			face := intValue(value)
			if face >= 0 && face <= 5 {
				faces[cube.Face(face)] = true
			}
		}
	}
	filter := map[string]bool{}
	if values, ok := server["placement_filter"].([]any); ok {
		for _, value := range values {
			if name := stringValue(value); name != "" {
				filter[name] = true
			}
		}
	}
	return faces, filter
}

func decodeCustomBlockData(value string) (decodedCustomBlockData, error) {
	if value == "" {
		return decodedCustomBlockData{}, nil
	}
	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.UseNumber()
	decoder.DisallowUnknownFields()
	var document customBlockDataDocument
	if err := decoder.Decode(&document); err != nil {
		return decodedCustomBlockData{}, err
	}
	if err := ensureJSONEnd(decoder); err != nil {
		return decodedCustomBlockData{}, err
	}
	server, err := normaliseComponentMap(document.Server)
	if err != nil {
		return decodedCustomBlockData{}, fmt.Errorf("server behaviours: %w", err)
	}
	states := make(map[string][]any, len(document.States))
	for name, values := range document.States {
		if name == "" || len(name) > 128 || len(values) == 0 || len(values) > 256 {
			return decodedCustomBlockData{}, fmt.Errorf("invalid state %q", name)
		}
		normalised, err := normaliseComponentValue(values)
		if err != nil {
			return decodedCustomBlockData{}, fmt.Errorf("state %q: %w", name, err)
		}
		states[name] = normalised.([]any)
	}
	permutations := make([]customBlockPermutation, 0, len(document.Permutations))
	if len(document.Permutations) > 256 {
		return decodedCustomBlockData{}, fmt.Errorf("too many permutations")
	}
	for index, permutation := range document.Permutations {
		if permutation.Condition == "" || len(permutation.Condition) > 1024 {
			return decodedCustomBlockData{}, fmt.Errorf("permutation %d has invalid condition", index)
		}
		values, err := normaliseComponentMap(permutation.Properties)
		if err != nil {
			return decodedCustomBlockData{}, fmt.Errorf("permutation %d: %w", index, err)
		}
		permutations = append(permutations, customBlockPermutation{condition: permutation.Condition, properties: values})
	}
	return decodedCustomBlockData{server: server, states: states, permutations: permutations}, nil
}

func customBlockStates(definitions map[string][]any) ([]map[string]any, error) {
	if len(definitions) == 0 {
		return []map[string]any{nil}, nil
	}
	names := make([]string, 0, len(definitions))
	for name := range definitions {
		names = append(names, name)
	}
	slices.Sort(names)
	states := []map[string]any{{}}
	for _, name := range names {
		values := definitions[name]
		if len(states) > 4096/len(values) {
			return nil, fmt.Errorf("state combinations exceed 4096")
		}
		next := make([]map[string]any, 0, len(states)*len(values))
		for _, state := range states {
			for _, value := range values {
				copyState := maps.Clone(state)
				copyState[name] = value
				next = append(next, copyState)
			}
		}
		states = next
	}
	return states, nil
}
