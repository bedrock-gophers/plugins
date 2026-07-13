use crate::{Item, item};

#[test]
fn armour_uses_typed_piece_and_tier_identifiers() {
    let helmet = item::Helmet::new(item::ArmourTier::Diamond);
    let chestplate = item::Chestplate::new(item::ArmourTier::Diamond);
    let leggings = item::Leggings::new(item::ArmourTier::Diamond);
    let boots = item::Boots::new(item::ArmourTier::Diamond);

    assert_eq!(helmet.tier(), item::ArmourTier::Diamond);
    assert_eq!(helmet.identifier(), "minecraft:diamond_helmet");
    assert_eq!(chestplate.identifier(), "minecraft:diamond_chestplate");
    assert_eq!(leggings.identifier(), "minecraft:diamond_leggings");
    assert_eq!(boots.identifier(), "minecraft:diamond_boots");

    macro_rules! assert_identifiers {
        ($(($actual:expr, $expected:literal $(,)?)),* $(,)?) => {
            $(assert_eq!($actual, $expected);)*
        };
    }
    assert_identifiers![
        (
            item::Helmet::new(item::ArmourTier::Leather).identifier(),
            "minecraft:leather_helmet",
        ),
        (
            item::Chestplate::new(item::ArmourTier::Leather).identifier(),
            "minecraft:leather_chestplate",
        ),
        (
            item::Leggings::new(item::ArmourTier::Leather).identifier(),
            "minecraft:leather_leggings",
        ),
        (
            item::Boots::new(item::ArmourTier::Leather).identifier(),
            "minecraft:leather_boots",
        ),
        (
            item::Helmet::new(item::ArmourTier::Copper).identifier(),
            "minecraft:copper_helmet",
        ),
        (
            item::Chestplate::new(item::ArmourTier::Copper).identifier(),
            "minecraft:copper_chestplate",
        ),
        (
            item::Leggings::new(item::ArmourTier::Copper).identifier(),
            "minecraft:copper_leggings",
        ),
        (
            item::Boots::new(item::ArmourTier::Copper).identifier(),
            "minecraft:copper_boots",
        ),
        (
            item::Helmet::new(item::ArmourTier::Gold).identifier(),
            "minecraft:golden_helmet",
        ),
        (
            item::Chestplate::new(item::ArmourTier::Gold).identifier(),
            "minecraft:golden_chestplate",
        ),
        (
            item::Leggings::new(item::ArmourTier::Gold).identifier(),
            "minecraft:golden_leggings",
        ),
        (
            item::Boots::new(item::ArmourTier::Gold).identifier(),
            "minecraft:golden_boots",
        ),
        (
            item::Helmet::new(item::ArmourTier::Chain).identifier(),
            "minecraft:chainmail_helmet",
        ),
        (
            item::Chestplate::new(item::ArmourTier::Chain).identifier(),
            "minecraft:chainmail_chestplate",
        ),
        (
            item::Leggings::new(item::ArmourTier::Chain).identifier(),
            "minecraft:chainmail_leggings",
        ),
        (
            item::Boots::new(item::ArmourTier::Chain).identifier(),
            "minecraft:chainmail_boots",
        ),
        (
            item::Helmet::new(item::ArmourTier::Iron).identifier(),
            "minecraft:iron_helmet",
        ),
        (
            item::Chestplate::new(item::ArmourTier::Iron).identifier(),
            "minecraft:iron_chestplate",
        ),
        (
            item::Leggings::new(item::ArmourTier::Iron).identifier(),
            "minecraft:iron_leggings",
        ),
        (
            item::Boots::new(item::ArmourTier::Iron).identifier(),
            "minecraft:iron_boots",
        ),
        (
            item::Helmet::new(item::ArmourTier::Diamond).identifier(),
            "minecraft:diamond_helmet",
        ),
        (
            item::Chestplate::new(item::ArmourTier::Diamond).identifier(),
            "minecraft:diamond_chestplate",
        ),
        (
            item::Leggings::new(item::ArmourTier::Diamond).identifier(),
            "minecraft:diamond_leggings",
        ),
        (
            item::Boots::new(item::ArmourTier::Diamond).identifier(),
            "minecraft:diamond_boots",
        ),
        (
            item::Helmet::new(item::ArmourTier::Netherite).identifier(),
            "minecraft:netherite_helmet",
        ),
        (
            item::Chestplate::new(item::ArmourTier::Netherite).identifier(),
            "minecraft:netherite_chestplate",
        ),
        (
            item::Leggings::new(item::ArmourTier::Netherite).identifier(),
            "minecraft:netherite_leggings",
        ),
        (
            item::Boots::new(item::ArmourTier::Netherite).identifier(),
            "minecraft:netherite_boots",
        ),
    ];
}

#[test]
fn beef_exposes_cooked_and_raw_typed_variants() {
    assert_eq!(item::Beef::raw().identifier(), "minecraft:beef");
    assert!(!item::Beef::raw().is_cooked());
    assert_eq!(item::Beef::cooked().identifier(), "minecraft:cooked_beef");
    assert!(item::Beef::cooked().is_cooked());
}
