package framework

import (
	"context"
	"image/color"
	"reflect"
	"testing"

	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
)

func TestDecodeTypedParticles(t *testing.T) {
	w := world.Config{Synchronous: true, Entities: entity.DefaultRegistry}.New()
	t.Cleanup(func() { _ = w.Close() })
	properties, ok := encodeBlockProperties(map[string]any{})
	if !ok {
		t.Fatal("encode empty block properties")
	}
	stone, ok := w.BlockRegistry().BlockByName("minecraft:stone", map[string]any{})
	if !ok {
		t.Fatal("resolve stone")
	}
	blockValue := &native.WorldBlock{Identifier: "minecraft:stone", PropertiesNBT: properties}
	for _, test := range []struct {
		name  string
		value native.WorldParticle
		want  world.Particle
	}{
		{"flame", native.WorldParticle{Kind: native.ParticleFlame, Colour: native.RGBA{R: 1, G: 2, B: 3, A: 4}}, particle.Flame{Colour: color.RGBA{R: 1, G: 2, B: 3, A: 4}}},
		{"dust", native.WorldParticle{Kind: native.ParticleDust, Colour: native.RGBA{R: 5, G: 6, B: 7, A: 8}}, particle.Dust{Colour: color.RGBA{R: 5, G: 6, B: 7, A: 8}}},
		{"block break", native.WorldParticle{Kind: native.ParticleBlockBreak, Block: blockValue}, particle.BlockBreak{Block: stone}},
		{"punch block", native.WorldParticle{Kind: native.ParticlePunchBlock, Data: uint32(cube.FaceEast), Block: blockValue}, particle.PunchBlock{Block: stone, Face: cube.FaceEast}},
		{"block force field", native.WorldParticle{Kind: native.ParticleBlockForceField}, particle.BlockForceField{}},
		{"bone meal", native.WorldParticle{Kind: native.ParticleBoneMeal, Data: 1}, particle.BoneMeal{Area: true}},
		{"note", native.WorldParticle{Kind: native.ParticleNote, Data: 14, Pitch: 24}, particle.Note{Instrument: sound.Banjo(), Pitch: 24}},
		{"dragon egg teleport", native.WorldParticle{Kind: native.ParticleDragonEggTeleport, Diff: native.BlockPos{X: -3, Y: 4, Z: 5}}, particle.DragonEggTeleport{Diff: cube.Pos{-3, 4, 5}}},
		{"evaporate", native.WorldParticle{Kind: native.ParticleEvaporate}, particle.Evaporate{}},
		{"water drip", native.WorldParticle{Kind: native.ParticleWaterDrip}, particle.WaterDrip{}},
		{"lava drip", native.WorldParticle{Kind: native.ParticleLavaDrip}, particle.LavaDrip{}},
		{"lava", native.WorldParticle{Kind: native.ParticleLava}, particle.Lava{}},
		{"dust plume", native.WorldParticle{Kind: native.ParticleDustPlume}, particle.DustPlume{}},
		{"huge explosion", native.WorldParticle{Kind: native.ParticleHugeExplosion}, particle.HugeExplosion{}},
		{"enderman teleport", native.WorldParticle{Kind: native.ParticleEndermanTeleport}, particle.EndermanTeleport{}},
		{"snowball poof", native.WorldParticle{Kind: native.ParticleSnowballPoof}, particle.SnowballPoof{}},
		{"egg smash", native.WorldParticle{Kind: native.ParticleEggSmash}, particle.EggSmash{}},
		{"splash", native.WorldParticle{Kind: native.ParticleSplash, Colour: native.RGBA{R: 9, G: 10, B: 11, A: 12}}, particle.Splash{Colour: color.RGBA{R: 9, G: 10, B: 11, A: 12}}},
		{"effect", native.WorldParticle{Kind: native.ParticleEffect, Colour: native.RGBA{R: 13, G: 14, B: 15, A: 16}}, particle.Effect{Colour: color.RGBA{R: 13, G: 14, B: 15, A: 16}}},
		{"entity flame", native.WorldParticle{Kind: native.ParticleEntityFlame}, particle.EntityFlame{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			decoded, ok := decodeParticle(w.BlockRegistry(), test.value)
			if !ok || !reflect.DeepEqual(decoded, test.want) {
				t.Fatalf("decoded = %#v, ok=%v, want %#v", decoded, ok, test.want)
			}
		})
	}
}

func TestValidateParticlePayload(t *testing.T) {
	properties, ok := encodeBlockProperties(map[string]any{})
	if !ok {
		t.Fatal("encode empty block properties")
	}
	for _, test := range []struct {
		name  string
		value native.WorldParticle
		valid bool
	}{
		{"block", native.WorldParticle{Kind: native.ParticleBlockBreak, Block: &native.WorldBlock{Identifier: "minecraft:stone", PropertiesNBT: properties}}, true},
		{"missing block", native.WorldParticle{Kind: native.ParticleBlockBreak}, false},
		{"invalid face", native.WorldParticle{Kind: native.ParticlePunchBlock, Data: 99, Block: &native.WorldBlock{Identifier: "minecraft:stone", PropertiesNBT: properties}}, false},
		{"invalid instrument", native.WorldParticle{Kind: native.ParticleNote, Data: 99}, false},
		{"dragonfly pitch", native.WorldParticle{Kind: native.ParticleNote, Pitch: 30}, true},
		{"unknown kind", native.WorldParticle{Kind: native.ParticleEntityFlame + 1}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if valid := validParticle(test.value); valid != test.valid {
				t.Fatalf("valid = %v, want %v", valid, test.valid)
			}
		})
	}
}

type particleRegistryTestBlock struct{}

func (particleRegistryTestBlock) EncodeBlock() (string, map[string]any) {
	return "example:particle", nil
}

func (particleRegistryTestBlock) Hash() (uint64, uint64)  { return 1 << 61, 0 }
func (particleRegistryTestBlock) Model() world.BlockModel { return nil }

func TestWorldParticleUsesInvocationAndTargetWorld(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	first, err := manager.Create("example:particles", world.Config{Synchronous: true, Entities: entity.DefaultRegistry})
	if err != nil {
		t.Fatal(err)
	}
	blocks := world.NewBlockRegistry()
	blocks.RegisterBlockState(world.BlockState{Name: "example:particle", Properties: map[string]any{}})
	blocks.RegisterBlock(particleRegistryTestBlock{})
	if _, err := manager.Create("example:particle-target", world.Config{Synchronous: true, Entities: entity.DefaultRegistry, Blocks: blocks}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	firstID, _ := manager.WorldByName(0, "example:particles")
	targetID, _ := manager.WorldByName(0, "example:particle-target")
	properties, ok := encodeBlockProperties(map[string]any{})
	if !ok {
		t.Fatal("encode empty block properties")
	}
	customBlock := &native.WorldBlock{Identifier: "example:particle", PropertiesNBT: properties}
	var stale native.InvocationID
	if err := first.Do(func(tx *world.Tx) {
		invocation, end := players.BeginInvocation(tx)
		stale = invocation
		if !manager.AddWorldParticle(invocation, 0, native.Vec3{Y: 64}, native.WorldParticle{Kind: native.ParticleHugeExplosion}) {
			t.Fatal("active invocation did not resolve world zero")
		}
		if !manager.AddWorldParticle(invocation, targetID, native.Vec3{Y: 64}, native.WorldParticle{Kind: native.ParticleBlockBreak, Block: customBlock}) {
			t.Fatal("cross-world particle was not decoded against the target registry and scheduled")
		}
		if manager.AddWorldParticle(invocation, firstID, native.Vec3{Y: 64}, native.WorldParticle{Kind: native.ParticleBlockBreak, Block: customBlock}) {
			t.Fatal("particle missing from target registry reported success")
		}
		end()
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if manager.AddWorldParticle(stale, 0, native.Vec3{Y: 64}, native.WorldParticle{Kind: native.ParticleHugeExplosion}) {
		t.Fatal("expired invocation resolved world zero")
	}
	if manager.AddWorldParticle(stale, firstID, native.Vec3{Y: 64}, native.WorldParticle{Kind: native.ParticleHugeExplosion}) {
		t.Fatal("expired invocation added a particle")
	}
}
