use crate::{BlockFace, BlockPos, block, sound::Instrument};

/// An RGBA colour used by coloured particles.
#[derive(Clone, Copy, Debug, Default, Eq, Hash, PartialEq)]
pub struct Colour {
    pub red: u8,
    pub green: u8,
    pub blue: u8,
    pub alpha: u8,
}

impl Colour {
    pub const fn new(red: u8, green: u8, blue: u8, alpha: u8) -> Self {
        Self {
            red,
            green,
            blue,
            alpha,
        }
    }

    pub const fn rgb(red: u8, green: u8, blue: u8) -> Self {
        Self::new(red, green, blue, u8::MAX)
    }
}

/// A particle that Dragonfly can add to a world.
///
/// This trait is sealed because the host ABI can only represent Dragonfly's
/// concrete built-in particle types.
pub trait Particle: private::Sealed {
    #[doc(hidden)]
    fn encode(&self) -> EncodedParticle<'_>;
}

/// The host-independent representation of a typed particle.
#[doc(hidden)]
#[derive(Clone, Copy, Debug)]
pub struct EncodedParticle<'a> {
    kind: u32,
    data: u32,
    pitch: i32,
    colour: Colour,
    diff: BlockPos,
    block: Option<&'a block::Block>,
}

impl EncodedParticle<'_> {
    pub(crate) fn with_raw<R>(
        &self,
        callback: impl FnOnce(&dragonfly_plugin_sys::DfParticleViewV1) -> R,
    ) -> Option<R> {
        let properties = match self.block {
            Some(block) => Some(block.properties_nbt()?),
            None => None,
        };
        let raw_block = self
            .block
            .zip(properties.as_ref())
            .map(|(block, properties)| dragonfly_plugin_sys::DfBlockView {
                identifier: crate::string_view_from_str(block.identifier()),
                properties_nbt: crate::bytes_view(properties),
            });
        let view = dragonfly_plugin_sys::DfParticleViewV1 {
            kind: self.kind,
            data: self.data,
            pitch: self.pitch,
            colour: dragonfly_plugin_sys::DfRgba {
                r: self.colour.red,
                g: self.colour.green,
                b: self.colour.blue,
                a: self.colour.alpha,
            },
            diff: self.diff.into(),
            block: raw_block
                .as_ref()
                .map_or(core::ptr::null(), core::ptr::from_ref),
        };
        Some(callback(&view))
    }
}

fn encoded(kind: u32) -> EncodedParticle<'static> {
    EncodedParticle {
        kind,
        data: 0,
        pitch: 0,
        colour: Colour::default(),
        diff: BlockPos::default(),
        block: None,
    }
}

macro_rules! simple_particles {
    ($($name:ident => $kind:ident),+ $(,)?) => {$ (
        #[derive(Clone, Copy, Debug, Default, Eq, Hash, PartialEq)]
        pub struct $name;

        impl private::Sealed for $name {}

        impl Particle for $name {
            fn encode(&self) -> EncodedParticle<'_> {
                encoded(dragonfly_plugin_sys::$kind)
            }
        }
    )+ };
}

simple_particles! {
    BlockForceField => DF_PARTICLE_BLOCK_FORCE_FIELD,
    Evaporate => DF_PARTICLE_EVAPORATE,
    WaterDrip => DF_PARTICLE_WATER_DRIP,
    LavaDrip => DF_PARTICLE_LAVA_DRIP,
    Lava => DF_PARTICLE_LAVA,
    DustPlume => DF_PARTICLE_DUST_PLUME,
    HugeExplosion => DF_PARTICLE_HUGE_EXPLOSION,
    EndermanTeleport => DF_PARTICLE_ENDERMAN_TELEPORT,
    SnowballPoof => DF_PARTICLE_SNOWBALL_POOF,
    EggSmash => DF_PARTICLE_EGG_SMASH,
    EntityFlame => DF_PARTICLE_ENTITY_FLAME,
}

#[derive(Clone, Copy, Debug, Default, Eq, Hash, PartialEq)]
pub struct Flame {
    colour: Option<Colour>,
}

impl Flame {
    /// Creates the ordinary, uncoloured flame particle.
    pub const fn new() -> Self {
        Self { colour: None }
    }

    pub const fn coloured(colour: Colour) -> Self {
        Self {
            colour: Some(colour),
        }
    }

    pub const fn colour(self) -> Option<Colour> {
        self.colour
    }
}

impl private::Sealed for Flame {}

impl Particle for Flame {
    fn encode(&self) -> EncodedParticle<'_> {
        EncodedParticle {
            colour: self.colour.unwrap_or_default(),
            ..encoded(dragonfly_plugin_sys::DF_PARTICLE_FLAME)
        }
    }
}

#[derive(Clone, Copy, Debug, Default, Eq, Hash, PartialEq)]
pub struct Dust {
    colour: Colour,
}

impl Dust {
    pub const fn new(colour: Colour) -> Self {
        Self { colour }
    }

    pub const fn colour(self) -> Colour {
        self.colour
    }
}

impl private::Sealed for Dust {}

impl Particle for Dust {
    fn encode(&self) -> EncodedParticle<'_> {
        EncodedParticle {
            colour: self.colour,
            ..encoded(dragonfly_plugin_sys::DF_PARTICLE_DUST)
        }
    }
}

#[derive(Clone, Debug, Eq, PartialEq)]
pub struct BlockBreak {
    block: block::Block,
}

impl BlockBreak {
    pub fn new(block: impl Into<block::Block>) -> Self {
        Self {
            block: block.into(),
        }
    }

    pub const fn block(&self) -> &block::Block {
        &self.block
    }
}

impl private::Sealed for BlockBreak {}

impl Particle for BlockBreak {
    fn encode(&self) -> EncodedParticle<'_> {
        EncodedParticle {
            block: Some(&self.block),
            ..encoded(dragonfly_plugin_sys::DF_PARTICLE_BLOCK_BREAK)
        }
    }
}

#[derive(Clone, Debug, Eq, PartialEq)]
pub struct PunchBlock {
    block: block::Block,
    face: BlockFace,
}

impl PunchBlock {
    pub fn new(block: impl Into<block::Block>, face: BlockFace) -> Self {
        Self {
            block: block.into(),
            face,
        }
    }

    pub const fn block(&self) -> &block::Block {
        &self.block
    }

    pub const fn face(&self) -> BlockFace {
        self.face
    }
}

impl private::Sealed for PunchBlock {}

impl Particle for PunchBlock {
    fn encode(&self) -> EncodedParticle<'_> {
        EncodedParticle {
            data: self.face as u32,
            block: Some(&self.block),
            ..encoded(dragonfly_plugin_sys::DF_PARTICLE_PUNCH_BLOCK)
        }
    }
}

#[derive(Clone, Copy, Debug, Default, Eq, Hash, PartialEq)]
pub struct BoneMeal {
    area: bool,
}

impl BoneMeal {
    pub const fn new(area: bool) -> Self {
        Self { area }
    }

    pub const fn area(self) -> bool {
        self.area
    }
}

impl private::Sealed for BoneMeal {}

impl Particle for BoneMeal {
    fn encode(&self) -> EncodedParticle<'_> {
        EncodedParticle {
            data: self.area as u32,
            ..encoded(dragonfly_plugin_sys::DF_PARTICLE_BONE_MEAL)
        }
    }
}

#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct Note {
    instrument: Instrument,
    pitch: i32,
}

impl Note {
    pub const fn new(instrument: Instrument, pitch: i32) -> Self {
        Self { instrument, pitch }
    }

    pub const fn instrument(self) -> Instrument {
        self.instrument
    }

    pub const fn pitch(self) -> i32 {
        self.pitch
    }
}

impl private::Sealed for Note {}

impl Particle for Note {
    fn encode(&self) -> EncodedParticle<'_> {
        EncodedParticle {
            data: self.instrument as u32,
            pitch: self.pitch,
            ..encoded(dragonfly_plugin_sys::DF_PARTICLE_NOTE)
        }
    }
}

#[derive(Clone, Copy, Debug, Default, Eq, Hash, PartialEq)]
pub struct DragonEggTeleport {
    diff: BlockPos,
}

impl DragonEggTeleport {
    pub const fn new(diff: BlockPos) -> Self {
        Self { diff }
    }

    pub const fn diff(self) -> BlockPos {
        self.diff
    }
}

impl private::Sealed for DragonEggTeleport {}

impl Particle for DragonEggTeleport {
    fn encode(&self) -> EncodedParticle<'_> {
        EncodedParticle {
            diff: self.diff,
            ..encoded(dragonfly_plugin_sys::DF_PARTICLE_DRAGON_EGG_TELEPORT)
        }
    }
}

#[derive(Clone, Copy, Debug, Default, Eq, Hash, PartialEq)]
pub struct Splash {
    colour: Option<Colour>,
}

impl Splash {
    /// Creates a splash using Dragonfly's default potion colour.
    pub const fn new() -> Self {
        Self { colour: None }
    }

    pub const fn coloured(colour: Colour) -> Self {
        Self {
            colour: Some(colour),
        }
    }

    pub const fn colour(self) -> Option<Colour> {
        self.colour
    }
}

impl private::Sealed for Splash {}

impl Particle for Splash {
    fn encode(&self) -> EncodedParticle<'_> {
        EncodedParticle {
            colour: self.colour.unwrap_or_default(),
            ..encoded(dragonfly_plugin_sys::DF_PARTICLE_SPLASH)
        }
    }
}

#[derive(Clone, Copy, Debug, Default, Eq, Hash, PartialEq)]
pub struct Effect {
    colour: Colour,
}

impl Effect {
    pub const fn new(colour: Colour) -> Self {
        Self { colour }
    }

    pub const fn colour(self) -> Colour {
        self.colour
    }
}

impl private::Sealed for Effect {}

impl Particle for Effect {
    fn encode(&self) -> EncodedParticle<'_> {
        EncodedParticle {
            colour: self.colour,
            ..encoded(dragonfly_plugin_sys::DF_PARTICLE_EFFECT)
        }
    }
}

mod private {
    pub trait Sealed {}
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn coloured_particle_preserves_rgba() {
        let colour = Colour::new(1, 2, 3, 4);
        let flame = Flame::coloured(colour);
        let encoded = flame.encode();
        let result = encoded.with_raw(|raw| {
            assert_eq!(raw.kind, dragonfly_plugin_sys::DF_PARTICLE_FLAME);
            assert_eq!(
                (raw.colour.r, raw.colour.g, raw.colour.b, raw.colour.a),
                (1, 2, 3, 4)
            );
            assert!(raw.block.is_null());
        });
        assert!(result.is_some());
    }

    #[test]
    fn default_flame_uses_dragonfly_zero_colour_sentinel() {
        let result = Flame::new().encode().with_raw(|raw| {
            assert_eq!(
                (raw.colour.r, raw.colour.g, raw.colour.b, raw.colour.a),
                (0, 0, 0, 0)
            );
        });
        assert!(result.is_some());
    }

    #[test]
    fn block_particle_keeps_block_view_alive_during_call() {
        let particle = PunchBlock::new(
            block::WoodenDoor::new().with_open_bit(true),
            BlockFace::North,
        );
        let result = particle.encode().with_raw(|raw| {
            assert_eq!(raw.kind, dragonfly_plugin_sys::DF_PARTICLE_PUNCH_BLOCK);
            assert_eq!(raw.data, BlockFace::North as u32);
            assert!(!raw.block.is_null());
            let block = unsafe { &*raw.block };
            assert_eq!(
                unsafe { crate::string_view(block.identifier) },
                "minecraft:wooden_door"
            );
            assert_ne!(block.properties_nbt.len, 0);
        });
        assert!(result.is_some());
    }

    #[test]
    fn note_and_dragon_egg_fields_map_directly() {
        let note = Note::new(Instrument::Banjo, 24).encode().with_raw(|raw| {
            assert_eq!(raw.data, Instrument::Banjo as u32);
            assert_eq!(raw.pitch, 24);
        });
        assert!(note.is_some());

        let teleport = DragonEggTeleport::new(BlockPos { x: -3, y: 4, z: 5 })
            .encode()
            .with_raw(|raw| {
                assert_eq!((raw.diff.x, raw.diff.y, raw.diff.z), (-3, 4, 5));
            });
        assert!(teleport.is_some());
    }
}
