package framework

import (
	"context"
	"testing"

	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
)

func TestDecodeTypedParticles(t *testing.T) {
	w := world.Config{Synchronous: true, Entities: entity.DefaultRegistry}.New()
	t.Cleanup(func() { _ = w.Close() })
	properties, ok := encodeBlockProperties(map[string]any{})
	if !ok {
		t.Fatal("encode empty block properties")
	}
	if err := w.Do(func(tx *world.Tx) {
		for _, test := range []struct {
			name  string
			value native.WorldParticle
			check func(world.Particle) bool
		}{
			{"colour", native.WorldParticle{Kind: native.ParticleFlame, Colour: native.RGBA{R: 1, G: 2, B: 3, A: 4}}, func(value world.Particle) bool { _, ok := value.(particle.Flame); return ok }},
			{"block", native.WorldParticle{Kind: native.ParticleBlockBreak, Block: &native.WorldBlock{Identifier: "minecraft:stone", PropertiesNBT: properties}}, func(value world.Particle) bool { _, ok := value.(particle.BlockBreak); return ok }},
			{"note", native.WorldParticle{Kind: native.ParticleNote, Data: 14, Pitch: 24}, func(value world.Particle) bool { note, ok := value.(particle.Note); return ok && note.Pitch == 24 }},
			{"teleport", native.WorldParticle{Kind: native.ParticleDragonEggTeleport, Diff: native.BlockPos{X: -3, Y: 4, Z: 5}}, func(value world.Particle) bool {
				teleport, ok := value.(particle.DragonEggTeleport)
				return ok && teleport.Diff.X() == -3
			}},
		} {
			t.Run(test.name, func(t *testing.T) {
				decoded, ok := decodeParticle(tx, test.value)
				if !ok || !test.check(decoded) {
					t.Fatalf("decoded = %T, ok=%v", decoded, ok)
				}
			})
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
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

func TestWorldParticleRejectsExpiredInvocation(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	w, err := manager.Create("example:particles", world.Config{Synchronous: true, Entities: entity.DefaultRegistry})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	id, _ := manager.WorldByName(0, "example:particles")
	var stale native.InvocationID
	if err := w.Do(func(tx *world.Tx) {
		var end func()
		stale, end = players.BeginInvocation(tx)
		end()
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if manager.AddWorldParticle(stale, id, native.Vec3{Y: 64}, native.WorldParticle{Kind: native.ParticleHugeExplosion}) {
		t.Fatal("expired invocation added a particle")
	}
}
