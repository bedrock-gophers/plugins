use crate::{ItemStack, block};

#[repr(u32)]
#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub enum Instrument {
    Piano = 0,
    BassDrum = 1,
    Snare = 2,
    ClicksAndSticks = 3,
    Bass = 4,
    Flute = 5,
    Bell = 6,
    Guitar = 7,
    Chimes = 8,
    Xylophone = 9,
    IronXylophone = 10,
    CowBell = 11,
    Didgeridoo = 12,
    Bit = 13,
    Banjo = 14,
    Pling = 15,
}

#[repr(u32)]
#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub enum DiscType {
    Disc13 = 0,
    DiscCat = 1,
    DiscBlocks = 2,
    DiscChirp = 3,
    DiscFar = 4,
    DiscMall = 5,
    DiscMellohi = 6,
    DiscStal = 7,
    DiscStrad = 8,
    DiscWard = 9,
    Disc11 = 10,
    DiscWait = 11,
    DiscOtherside = 12,
    DiscPigstep = 13,
    Disc5 = 14,
    DiscRelic = 15,
    DiscCreator = 16,
    DiscCreatorMusicBox = 17,
    DiscPrecipice = 18,
    DiscTears = 19,
    DiscLavaChicken = 20,
}

#[repr(u32)]
#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub enum Horn {
    Ponder = 0,
    Sing = 1,
    Seek = 2,
    Feel = 3,
    Admire = 4,
    Call = 5,
    Yearn = 6,
    Dream = 7,
}

#[repr(u32)]
#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub enum CrossbowStage {
    Start = 0,
    Middle = 1,
    End = 2,
}

#[repr(u32)]
#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub enum Liquid {
    Water = 0,
    Lava = 1,
}

/// A sound that Dragonfly can play in a world.
///
/// This trait is sealed because the host ABI can only represent Dragonfly's
/// concrete built-in sound types.
pub trait Sound: private::Sealed {
    #[doc(hidden)]
    fn encode(&self) -> EncodedSound<'_>;
}

/// The host-independent representation of a typed sound.
#[doc(hidden)]
#[derive(Clone, Copy, Debug)]
pub struct EncodedSound<'a> {
    kind: u32,
    data: u32,
    integer: i32,
    flags: u32,
    scalar: f64,
    block: Option<&'a block::Block>,
    item: Option<&'a ItemStack>,
}

impl EncodedSound<'_> {
    pub(crate) fn with_raw<R>(
        &self,
        callback: impl FnOnce(&dragonfly_plugin_sys::DfSoundViewV1) -> R,
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
        let block_pointer = raw_block
            .as_ref()
            .map_or(core::ptr::null(), core::ptr::from_ref);
        let make_view = |item| dragonfly_plugin_sys::DfSoundViewV1 {
            kind: self.kind,
            data: self.data,
            integer: self.integer,
            flags: self.flags,
            scalar: self.scalar,
            block: block_pointer,
            item,
        };
        match self.item {
            Some(item) => {
                item.with_raw(|raw_item| callback(&make_view(core::ptr::from_ref(raw_item))))
            }
            None => Some(callback(&make_view(core::ptr::null()))),
        }
    }
}

fn encoded(kind: u32) -> EncodedSound<'static> {
    EncodedSound {
        kind,
        data: 0,
        integer: 0,
        flags: 0,
        scalar: 0.0,
        block: None,
        item: None,
    }
}

// The constants are passed separately because Rust identifiers cannot be
// converted to screaming snake case by a declarative macro.
macro_rules! simple_sound_kinds {
    ($($name:ident => $kind:ident),+ $(,)?) => {$ (
        #[derive(Clone, Copy, Debug, Default, Eq, Hash, PartialEq)]
        pub struct $name;

        impl private::Sealed for $name {}

        impl Sound for $name {
            fn encode(&self) -> EncodedSound<'_> {
                encoded(dragonfly_plugin_sys::$kind)
            }
        }
    )+ };
}

simple_sound_kinds! {
    AnvilBreak => DF_SOUND_KIND_ANVIL_BREAK,
    AnvilLand => DF_SOUND_KIND_ANVIL_LAND,
    AnvilUse => DF_SOUND_KIND_ANVIL_USE,
    ArrowHit => DF_SOUND_KIND_ARROW_HIT,
    BarrelClose => DF_SOUND_KIND_BARREL_CLOSE,
    BarrelOpen => DF_SOUND_KIND_BARREL_OPEN,
    BlastFurnaceCrackle => DF_SOUND_KIND_BLAST_FURNACE_CRACKLE,
    BowShoot => DF_SOUND_KIND_BOW_SHOOT,
    Burning => DF_SOUND_KIND_BURNING,
    Burp => DF_SOUND_KIND_BURP,
    CampfireCrackle => DF_SOUND_KIND_CAMPFIRE_CRACKLE,
    ChestClose => DF_SOUND_KIND_CHEST_CLOSE,
    ChestOpen => DF_SOUND_KIND_CHEST_OPEN,
    Click => DF_SOUND_KIND_CLICK,
    ComposterEmpty => DF_SOUND_KIND_COMPOSTER_EMPTY,
    ComposterFill => DF_SOUND_KIND_COMPOSTER_FILL,
    ComposterFillLayer => DF_SOUND_KIND_COMPOSTER_FILL_LAYER,
    ComposterReady => DF_SOUND_KIND_COMPOSTER_READY,
    CopperScraped => DF_SOUND_KIND_COPPER_SCRAPED,
    CrossbowShoot => DF_SOUND_KIND_CROSSBOW_SHOOT,
    DecoratedPotInsertFailed => DF_SOUND_KIND_DECORATED_POT_INSERT_FAILED,
    Deny => DF_SOUND_KIND_DENY,
    DoorCrash => DF_SOUND_KIND_DOOR_CRASH,
    Drowning => DF_SOUND_KIND_DROWNING,
    EndPortalCreated => DF_SOUND_KIND_END_PORTAL_CREATED,
    EnderChestClose => DF_SOUND_KIND_ENDER_CHEST_CLOSE,
    EnderChestOpen => DF_SOUND_KIND_ENDER_CHEST_OPEN,
    EnderEyePlaced => DF_SOUND_KIND_ENDER_EYE_PLACED,
    Experience => DF_SOUND_KIND_EXPERIENCE,
    Explosion => DF_SOUND_KIND_EXPLOSION,
    FireCharge => DF_SOUND_KIND_FIRE_CHARGE,
    FireExtinguish => DF_SOUND_KIND_FIRE_EXTINGUISH,
    FireworkBlast => DF_SOUND_KIND_FIREWORK_BLAST,
    FireworkHugeBlast => DF_SOUND_KIND_FIREWORK_HUGE_BLAST,
    FireworkLaunch => DF_SOUND_KIND_FIREWORK_LAUNCH,
    FireworkTwinkle => DF_SOUND_KIND_FIREWORK_TWINKLE,
    Fizz => DF_SOUND_KIND_FIZZ,
    FurnaceCrackle => DF_SOUND_KIND_FURNACE_CRACKLE,
    GhastShoot => DF_SOUND_KIND_GHAST_SHOOT,
    GhastWarning => DF_SOUND_KIND_GHAST_WARNING,
    GlassBreak => DF_SOUND_KIND_GLASS_BREAK,
    Ignite => DF_SOUND_KIND_IGNITE,
    ItemAdd => DF_SOUND_KIND_ITEM_ADD,
    ItemBreak => DF_SOUND_KIND_ITEM_BREAK,
    ItemFrameRemove => DF_SOUND_KIND_ITEM_FRAME_REMOVE,
    ItemFrameRotate => DF_SOUND_KIND_ITEM_FRAME_ROTATE,
    ItemThrow => DF_SOUND_KIND_ITEM_THROW,
    LecternBookPlace => DF_SOUND_KIND_LECTERN_BOOK_PLACE,
    LevelUp => DF_SOUND_KIND_LEVEL_UP,
    LightningExplode => DF_SOUND_KIND_LIGHTNING_EXPLODE,
    LightningThunder => DF_SOUND_KIND_LIGHTNING_THUNDER,
    MusicDiscEnd => DF_SOUND_KIND_MUSIC_DISC_END,
    Pop => DF_SOUND_KIND_POP,
    PotionBrewed => DF_SOUND_KIND_POTION_BREWED,
    PowerOff => DF_SOUND_KIND_POWER_OFF,
    PowerOn => DF_SOUND_KIND_POWER_ON,
    ShulkerBoxClose => DF_SOUND_KIND_SHULKER_BOX_CLOSE,
    ShulkerBoxOpen => DF_SOUND_KIND_SHULKER_BOX_OPEN,
    SignWaxed => DF_SOUND_KIND_SIGN_WAXED,
    SmokerCrackle => DF_SOUND_KIND_SMOKER_CRACKLE,
    StopUsingSpyglass => DF_SOUND_KIND_STOP_USING_SPYGLASS,
    TNT => DF_SOUND_KIND_TNT,
    Teleport => DF_SOUND_KIND_TELEPORT,
    Thunder => DF_SOUND_KIND_THUNDER,
    Totem => DF_SOUND_KIND_TOTEM,
    UseSpyglass => DF_SOUND_KIND_USE_SPYGLASS,
    WaxRemoved => DF_SOUND_KIND_WAX_REMOVED,
    WaxedSignFailedInteraction => DF_SOUND_KIND_WAXED_SIGN_FAILED_INTERACTION,
}

macro_rules! block_sounds {
    ($($name:ident => $kind:ident),+ $(,)?) => {$ (
        #[derive(Clone, Debug, Eq, PartialEq)]
        pub struct $name {
            block: block::Block,
        }

        impl $name {
            pub const fn new(block: block::Block) -> Self {
                Self { block }
            }

            pub const fn block(&self) -> &block::Block {
                &self.block
            }
        }

        impl private::Sealed for $name {}

        impl Sound for $name {
            fn encode(&self) -> EncodedSound<'_> {
                EncodedSound {
                    block: Some(&self.block),
                    ..encoded(dragonfly_plugin_sys::$kind)
                }
            }
        }
    )+ };
}

block_sounds! {
    BlockPlace => DF_SOUND_KIND_BLOCK_PLACE,
    BlockBreaking => DF_SOUND_KIND_BLOCK_BREAKING,
    DoorOpen => DF_SOUND_KIND_DOOR_OPEN,
    DoorClose => DF_SOUND_KIND_DOOR_CLOSE,
    TrapdoorOpen => DF_SOUND_KIND_TRAPDOOR_OPEN,
    TrapdoorClose => DF_SOUND_KIND_TRAPDOOR_CLOSE,
    FenceGateOpen => DF_SOUND_KIND_FENCE_GATE_OPEN,
    FenceGateClose => DF_SOUND_KIND_FENCE_GATE_CLOSE,
    ItemUseOn => DF_SOUND_KIND_ITEM_USE_ON,
}

#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct Attack {
    damage: bool,
}

impl Attack {
    pub const fn new(damage: bool) -> Self {
        Self { damage }
    }

    pub const fn damage(self) -> bool {
        self.damage
    }
}

impl private::Sealed for Attack {}

impl Sound for Attack {
    fn encode(&self) -> EncodedSound<'_> {
        EncodedSound {
            flags: self.damage as u32,
            ..encoded(dragonfly_plugin_sys::DF_SOUND_KIND_ATTACK)
        }
    }
}

#[derive(Clone, Copy, Debug, PartialEq)]
pub struct Fall {
    distance: f64,
}

impl Fall {
    pub const fn new(distance: f64) -> Self {
        Self { distance }
    }

    pub const fn distance(self) -> f64 {
        self.distance
    }
}

impl private::Sealed for Fall {}

impl Sound for Fall {
    fn encode(&self) -> EncodedSound<'_> {
        EncodedSound {
            scalar: self.distance,
            ..encoded(dragonfly_plugin_sys::DF_SOUND_KIND_FALL)
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

impl Sound for Note {
    fn encode(&self) -> EncodedSound<'_> {
        EncodedSound {
            data: self.instrument as u32,
            integer: self.pitch,
            ..encoded(dragonfly_plugin_sys::DF_SOUND_KIND_NOTE)
        }
    }
}

#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct MusicDiscPlay {
    disc_type: DiscType,
}

impl MusicDiscPlay {
    pub const fn new(disc_type: DiscType) -> Self {
        Self { disc_type }
    }

    pub const fn disc_type(self) -> DiscType {
        self.disc_type
    }
}

impl private::Sealed for MusicDiscPlay {}

impl Sound for MusicDiscPlay {
    fn encode(&self) -> EncodedSound<'_> {
        EncodedSound {
            data: self.disc_type as u32,
            ..encoded(dragonfly_plugin_sys::DF_SOUND_KIND_MUSIC_DISC_PLAY)
        }
    }
}

#[derive(Clone, Copy, Debug, PartialEq)]
pub struct DecoratedPotInserted {
    progress: f64,
}

impl DecoratedPotInserted {
    pub const fn new(progress: f64) -> Self {
        Self { progress }
    }

    pub const fn progress(self) -> f64 {
        self.progress
    }
}

impl private::Sealed for DecoratedPotInserted {}

impl Sound for DecoratedPotInserted {
    fn encode(&self) -> EncodedSound<'_> {
        EncodedSound {
            scalar: self.progress,
            ..encoded(dragonfly_plugin_sys::DF_SOUND_KIND_DECORATED_POT_INSERTED)
        }
    }
}

#[derive(Clone, Debug, PartialEq)]
pub struct EquipItem {
    item: ItemStack,
}

impl EquipItem {
    pub fn new(item: impl crate::Item) -> Self {
        Self {
            item: crate::item::new(item, 1),
        }
    }

    pub const fn item(&self) -> &ItemStack {
        &self.item
    }
}

impl private::Sealed for EquipItem {}

impl Sound for EquipItem {
    fn encode(&self) -> EncodedSound<'_> {
        EncodedSound {
            item: Some(&self.item),
            ..encoded(dragonfly_plugin_sys::DF_SOUND_KIND_EQUIP_ITEM)
        }
    }
}

macro_rules! liquid_sounds {
    ($($name:ident => $kind:ident),+ $(,)?) => {$ (
        #[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
        pub struct $name {
            liquid: Liquid,
        }

        impl $name {
            pub const fn new(liquid: Liquid) -> Self {
                Self { liquid }
            }

            pub const fn liquid(self) -> Liquid {
                self.liquid
            }
        }

        impl private::Sealed for $name {}

        impl Sound for $name {
            fn encode(&self) -> EncodedSound<'_> {
                EncodedSound {
                    data: self.liquid as u32,
                    ..encoded(dragonfly_plugin_sys::$kind)
                }
            }
        }
    )+ };
}

liquid_sounds! {
    BucketFill => DF_SOUND_KIND_BUCKET_FILL,
    BucketEmpty => DF_SOUND_KIND_BUCKET_EMPTY,
}

#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct CrossbowLoad {
    stage: CrossbowStage,
    quick_charge: bool,
}

impl CrossbowLoad {
    pub const fn new(stage: CrossbowStage, quick_charge: bool) -> Self {
        Self {
            stage,
            quick_charge,
        }
    }

    pub const fn stage(self) -> CrossbowStage {
        self.stage
    }

    pub const fn quick_charge(self) -> bool {
        self.quick_charge
    }
}

impl private::Sealed for CrossbowLoad {}

impl Sound for CrossbowLoad {
    fn encode(&self) -> EncodedSound<'_> {
        EncodedSound {
            integer: self.stage as i32,
            flags: self.quick_charge as u32,
            ..encoded(dragonfly_plugin_sys::DF_SOUND_KIND_CROSSBOW_LOAD)
        }
    }
}

#[derive(Clone, Copy, Debug, Eq, Hash, PartialEq)]
pub struct GoatHorn {
    horn: Horn,
}

impl GoatHorn {
    pub const fn new(horn: Horn) -> Self {
        Self { horn }
    }

    pub const fn horn(self) -> Horn {
        self.horn
    }
}

impl private::Sealed for GoatHorn {}

impl Sound for GoatHorn {
    fn encode(&self) -> EncodedSound<'_> {
        EncodedSound {
            data: self.horn as u32,
            ..encoded(dragonfly_plugin_sys::DF_SOUND_KIND_GOAT_HORN)
        }
    }
}

mod private {
    pub trait Sealed {}
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::item;

    #[test]
    fn simple_and_parameter_fields_encode_directly() {
        let simple = Explosion.encode().with_raw(|raw| {
            assert_eq!(raw.kind, dragonfly_plugin_sys::DF_SOUND_KIND_EXPLOSION);
            assert!(raw.block.is_null());
            assert!(raw.item.is_null());
        });
        assert!(simple.is_some());

        let crossbow = CrossbowLoad::new(CrossbowStage::Middle, true)
            .encode()
            .with_raw(|raw| {
                assert_eq!(raw.integer, CrossbowStage::Middle as i32);
                assert_eq!(raw.flags, 1);
            });
        assert!(crossbow.is_some());

        let note = Note::new(Instrument::Banjo, 24).encode().with_raw(|raw| {
            assert_eq!(raw.data, Instrument::Banjo as u32);
            assert_eq!(raw.integer, 24);
        });
        assert!(note.is_some());
    }

    #[test]
    fn block_view_lives_for_callback() {
        let sound = DoorOpen::new(block::new("minecraft:oak_door").with_property("open_bit", true));
        let result = sound.encode().with_raw(|raw| {
            assert!(!raw.block.is_null());
            let block = unsafe { &*raw.block };
            assert_eq!(
                unsafe { crate::string_view(block.identifier) },
                "minecraft:oak_door"
            );
            assert_ne!(block.properties_nbt.len, 0);
        });
        assert!(result.is_some());
    }

    #[test]
    fn item_view_lives_for_callback() {
        let sound = EquipItem::new(item::Sword::new(item::ToolTier::Diamond));
        let result = sound.encode().with_raw(|raw| {
            assert!(raw.block.is_null());
            assert!(!raw.item.is_null());
            let item = unsafe { &*raw.item };
            assert_eq!(
                unsafe { crate::string_view(item.identifier) },
                "minecraft:diamond_sword"
            );
            assert_eq!(item.count, 1);
        });
        assert!(result.is_some());
    }

    #[test]
    fn every_enum_matches_dragonfly_discriminants() {
        assert_eq!(Instrument::Pling as u32, 15);
        assert_eq!(DiscType::DiscLavaChicken as u32, 20);
        assert_eq!(Horn::Dream as u32, 7);
        assert_eq!(CrossbowStage::End as u32, 2);
        assert_eq!(Liquid::Lava as u32, 1);
    }
}
