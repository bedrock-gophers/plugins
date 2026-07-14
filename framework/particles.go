package framework

import (
	"image/color"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
)

func (m *WorldManager) AddWorldParticle(invocation native.InvocationID, id native.WorldID, position native.Vec3, value native.WorldParticle) bool {
	entry, ok := m.entryForInvocation(invocation, id)
	if !ok || !finiteVec3(position) || !validParticle(value) {
		return false
	}
	entry.lifecycle.RLock()
	defer entry.lifecycle.RUnlock()
	if entry.closed {
		return false
	}
	decoded, ok := decodeParticle(entry.world.BlockRegistry(), value)
	if !ok {
		return false
	}
	return m.writeTx(invocation, entry, func(tx *world.Tx) {
		tx.AddParticle(vec3(position), decoded)
	})
}

func validParticle(value native.WorldParticle) bool {
	if value.Kind > native.ParticleEntityFlame {
		return false
	}
	switch value.Kind {
	case native.ParticleBlockBreak:
		return validParticleBlock(value.Block)
	case native.ParticlePunchBlock:
		return value.Data <= uint32(cube.FaceEast) && validParticleBlock(value.Block)
	case native.ParticleBoneMeal:
		return value.Data <= 1
	case native.ParticleNote:
		_, ok := particleInstrument(value.Data)
		return ok
	default:
		return true
	}
}

func validParticleBlock(value *native.WorldBlock) bool {
	if value == nil || value.Identifier == "" {
		return false
	}
	_, ok := decodeBlockProperties(value.PropertiesNBT)
	return ok
}

func decodeParticle(registry world.BlockRegistry, value native.WorldParticle) (world.Particle, bool) {
	colour := color.RGBA{R: value.Colour.R, G: value.Colour.G, B: value.Colour.B, A: value.Colour.A}
	switch value.Kind {
	case native.ParticleFlame:
		return particle.Flame{Colour: colour}, true
	case native.ParticleDust:
		return particle.Dust{Colour: colour}, true
	case native.ParticleBlockBreak:
		block, ok := particleBlock(registry, value.Block)
		return particle.BlockBreak{Block: block}, ok
	case native.ParticlePunchBlock:
		block, ok := particleBlock(registry, value.Block)
		if !ok || value.Data > uint32(cube.FaceEast) {
			return nil, false
		}
		return particle.PunchBlock{Block: block, Face: cube.Face(value.Data)}, true
	case native.ParticleBlockForceField:
		return particle.BlockForceField{}, true
	case native.ParticleBoneMeal:
		if value.Data > 1 {
			return nil, false
		}
		return particle.BoneMeal{Area: value.Data != 0}, true
	case native.ParticleNote:
		instrument, ok := particleInstrument(value.Data)
		if !ok {
			return nil, false
		}
		return particle.Note{Instrument: instrument, Pitch: int(value.Pitch)}, true
	case native.ParticleDragonEggTeleport:
		return particle.DragonEggTeleport{Diff: cube.Pos{int(value.Diff.X), int(value.Diff.Y), int(value.Diff.Z)}}, true
	case native.ParticleEvaporate:
		return particle.Evaporate{}, true
	case native.ParticleWaterDrip:
		return particle.WaterDrip{}, true
	case native.ParticleLavaDrip:
		return particle.LavaDrip{}, true
	case native.ParticleLava:
		return particle.Lava{}, true
	case native.ParticleDustPlume:
		return particle.DustPlume{}, true
	case native.ParticleHugeExplosion:
		return particle.HugeExplosion{}, true
	case native.ParticleEndermanTeleport:
		return particle.EndermanTeleport{}, true
	case native.ParticleSnowballPoof:
		return particle.SnowballPoof{}, true
	case native.ParticleEggSmash:
		return particle.EggSmash{}, true
	case native.ParticleSplash:
		return particle.Splash{Colour: colour}, true
	case native.ParticleEffect:
		return particle.Effect{Colour: colour}, true
	case native.ParticleEntityFlame:
		return particle.EntityFlame{}, true
	default:
		return nil, false
	}
}

func particleBlock(registry world.BlockRegistry, value *native.WorldBlock) (world.Block, bool) {
	if value == nil {
		return nil, false
	}
	properties, ok := decodeBlockProperties(value.PropertiesNBT)
	if !ok {
		return nil, false
	}
	return registry.BlockByName(value.Identifier, properties)
}

func particleInstrument(value uint32) (sound.Instrument, bool) {
	instruments := [...]sound.Instrument{
		sound.Piano(), sound.BassDrum(), sound.Snare(), sound.ClicksAndSticks(),
		sound.Bass(), sound.Flute(), sound.Bell(), sound.Guitar(), sound.Chimes(),
		sound.Xylophone(), sound.IronXylophone(), sound.CowBell(), sound.Didgeridoo(),
		sound.Bit(), sound.Banjo(), sound.Pling(),
	}
	if value >= uint32(len(instruments)) {
		return sound.Instrument{}, false
	}
	return instruments[value], true
}
