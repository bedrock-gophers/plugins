// Code generated from Dragonfly server/item Go AST and live registry. DO NOT EDIT.
#nullable enable

namespace Dragonfly
{
    public static partial class Item
    {
        public readonly record struct ToolTier(int HarvestLevel, double BaseMiningEfficiency, double BaseAttackDamage, int EnchantmentValue, int Durability, string Name);

        public static readonly ToolTier ToolTierWood = new(1, 2d, 1d, 15, 59, "wooden");
        public static readonly ToolTier ToolTierGold = new(1, 12d, 1d, 22, 32, "golden");
        public static readonly ToolTier ToolTierStone = new(2, 4d, 2d, 5, 131, "stone");
        public static readonly ToolTier ToolTierCopper = new(2, 5d, 2d, 13, 190, "copper");
        public static readonly ToolTier ToolTierIron = new(3, 6d, 3d, 14, 250, "iron");
        public static readonly ToolTier ToolTierDiamond = new(4, 8d, 4d, 10, 1561, "diamond");
        public static readonly ToolTier ToolTierNetherite = new(4, 9d, 5d, 15, 2031, "netherite");

        public readonly record struct AmethystShard : World.Item;
        public readonly record struct Apple : World.Item;
        public readonly record struct Axe(ToolTier Tier) : World.Item;
        public readonly record struct BakedPotato : World.Item;
        public readonly record struct Beef(bool Cooked) : World.Item;
        public readonly record struct Beetroot : World.Item;
        public readonly record struct BeetrootSoup : World.Item;
        public readonly record struct BlazePowder : World.Item;
        public readonly record struct BlazeRod : World.Item;
        public readonly record struct Bone : World.Item;
        public readonly record struct BoneMeal : World.Item;
        public readonly record struct Book : World.Item;
        public readonly record struct BottleOfEnchanting : World.Item;
        public readonly record struct Bow : World.Item;
        public readonly record struct Bowl : World.Item;
        public readonly record struct Bread : World.Item;
        public readonly record struct Brick : World.Item;
        public readonly record struct CarrotOnAStick : World.Item;
        public readonly record struct Charcoal : World.Item;
        public readonly record struct Chicken(bool Cooked) : World.Item;
        public readonly record struct ClayBall : World.Item;
        public readonly record struct Clock : World.Item;
        public readonly record struct Coal : World.Item;
        public readonly record struct Cod(bool Cooked) : World.Item;
        public readonly record struct Compass : World.Item;
        public readonly record struct Cookie : World.Item;
        public readonly record struct CopperIngot : World.Item;
        public readonly record struct CopperNugget : World.Item;
        public readonly record struct Diamond : World.Item;
        public readonly record struct DiscFragment : World.Item;
        public readonly record struct DragonBreath : World.Item;
        public readonly record struct DriedKelp : World.Item;
        public readonly record struct EchoShard : World.Item;
        public readonly record struct Egg : World.Item;
        public readonly record struct Elytra : World.Item;
        public readonly record struct Emerald : World.Item;
        public readonly record struct EnchantedApple : World.Item;
        public readonly record struct EnchantedBook : World.Item;
        public readonly record struct EnderEye : World.Item;
        public readonly record struct EnderPearl : World.Item;
        public readonly record struct Feather : World.Item;
        public readonly record struct FermentedSpiderEye : World.Item;
        public readonly record struct FireCharge : World.Item;
        public readonly record struct Flint : World.Item;
        public readonly record struct FlintAndSteel : World.Item;
        public readonly record struct GhastTear : World.Item;
        public readonly record struct GlassBottle : World.Item;
        public readonly record struct GlisteringMelonSlice : World.Item;
        public readonly record struct GlowstoneDust : World.Item;
        public readonly record struct GoldIngot : World.Item;
        public readonly record struct GoldNugget : World.Item;
        public readonly record struct GoldenApple : World.Item;
        public readonly record struct GoldenCarrot : World.Item;
        public readonly record struct Gunpowder : World.Item;
        public readonly record struct HeartOfTheSea : World.Item;
        public readonly record struct Hoe(ToolTier Tier) : World.Item;
        public readonly record struct HoneyBottle : World.Item;
        public readonly record struct Honeycomb : World.Item;
        public readonly record struct InkSac(bool Glowing) : World.Item;
        public readonly record struct IronIngot : World.Item;
        public readonly record struct IronNugget : World.Item;
        public readonly record struct LapisLazuli : World.Item;
        public readonly record struct Leather : World.Item;
        public readonly record struct MagmaCream : World.Item;
        public readonly record struct MelonSlice : World.Item;
        public readonly record struct MushroomStew : World.Item;
        public readonly record struct Mutton(bool Cooked) : World.Item;
        public readonly record struct NautilusShell : World.Item;
        public readonly record struct NetherBrick : World.Item;
        public readonly record struct NetherQuartz : World.Item;
        public readonly record struct NetherStar : World.Item;
        public readonly record struct NetheriteIngot : World.Item;
        public readonly record struct NetheriteScrap : World.Item;
        public readonly record struct Paper : World.Item;
        public readonly record struct PhantomMembrane : World.Item;
        public readonly record struct Pickaxe(ToolTier Tier) : World.Item;
        public readonly record struct PoisonousPotato : World.Item;
        public readonly record struct PoppedChorusFruit : World.Item;
        public readonly record struct Porkchop(bool Cooked) : World.Item;
        public readonly record struct PrismarineCrystals : World.Item;
        public readonly record struct PrismarineShard : World.Item;
        public readonly record struct Pufferfish : World.Item;
        public readonly record struct PumpkinPie : World.Item;
        public readonly record struct Rabbit(bool Cooked) : World.Item;
        public readonly record struct RabbitFoot : World.Item;
        public readonly record struct RabbitHide : World.Item;
        public readonly record struct RabbitStew : World.Item;
        public readonly record struct RawCopper : World.Item;
        public readonly record struct RawGold : World.Item;
        public readonly record struct RawIron : World.Item;
        public readonly record struct RecoveryCompass : World.Item;
        public readonly record struct ResinBrick : World.Item;
        public readonly record struct RottenFlesh : World.Item;
        public readonly record struct Salmon(bool Cooked) : World.Item;
        public readonly record struct Scute : World.Item;
        public readonly record struct Shears : World.Item;
        public readonly record struct Shovel(ToolTier Tier) : World.Item;
        public readonly record struct ShulkerShell : World.Item;
        public readonly record struct Slimeball : World.Item;
        public readonly record struct Snowball : World.Item;
        public readonly record struct SpiderEye : World.Item;
        public readonly record struct Spyglass : World.Item;
        public readonly record struct Stick : World.Item;
        public readonly record struct Sugar : World.Item;
        public readonly record struct Sword(ToolTier Tier) : World.Item;
        public readonly record struct Totem : World.Item;
        public readonly record struct TropicalFish : World.Item;
        public readonly record struct TurtleShell : World.Item;
        public readonly record struct WarpedFungusOnAStick : World.Item;
        public readonly record struct Wheat : World.Item;
    }

    internal static class ItemCodec
    {
        internal static bool TryEncode(World.Item item, out string identifier, out int metadata)
        {
            switch (item)
            {
                case Item.AmethystShard _:
                    identifier = "minecraft:amethyst_shard"; metadata = 0; return true;
                case Item.Apple _:
                    identifier = "minecraft:apple"; metadata = 0; return true;
                case Item.Axe value when value.Tier == Item.ToolTierWood:
                    identifier = "minecraft:wooden_axe"; metadata = 0; return true;
                case Item.Axe value when value.Tier == Item.ToolTierGold:
                    identifier = "minecraft:golden_axe"; metadata = 0; return true;
                case Item.Axe value when value.Tier == Item.ToolTierStone:
                    identifier = "minecraft:stone_axe"; metadata = 0; return true;
                case Item.Axe value when value.Tier == Item.ToolTierCopper:
                    identifier = "minecraft:copper_axe"; metadata = 0; return true;
                case Item.Axe value when value.Tier == Item.ToolTierIron:
                    identifier = "minecraft:iron_axe"; metadata = 0; return true;
                case Item.Axe value when value.Tier == Item.ToolTierDiamond:
                    identifier = "minecraft:diamond_axe"; metadata = 0; return true;
                case Item.Axe value when value.Tier == Item.ToolTierNetherite:
                    identifier = "minecraft:netherite_axe"; metadata = 0; return true;
                case Item.BakedPotato _:
                    identifier = "minecraft:baked_potato"; metadata = 0; return true;
                case Item.Beef { Cooked: false }:
                    identifier = "minecraft:beef"; metadata = 0; return true;
                case Item.Beef { Cooked: true }:
                    identifier = "minecraft:cooked_beef"; metadata = 0; return true;
                case Item.Beetroot _:
                    identifier = "minecraft:beetroot"; metadata = 0; return true;
                case Item.BeetrootSoup _:
                    identifier = "minecraft:beetroot_soup"; metadata = 0; return true;
                case Item.BlazePowder _:
                    identifier = "minecraft:blaze_powder"; metadata = 0; return true;
                case Item.BlazeRod _:
                    identifier = "minecraft:blaze_rod"; metadata = 0; return true;
                case Item.Bone _:
                    identifier = "minecraft:bone"; metadata = 0; return true;
                case Item.BoneMeal _:
                    identifier = "minecraft:bone_meal"; metadata = 0; return true;
                case Item.Book _:
                    identifier = "minecraft:book"; metadata = 0; return true;
                case Item.BottleOfEnchanting _:
                    identifier = "minecraft:experience_bottle"; metadata = 0; return true;
                case Item.Bow _:
                    identifier = "minecraft:bow"; metadata = 0; return true;
                case Item.Bowl _:
                    identifier = "minecraft:bowl"; metadata = 0; return true;
                case Item.Bread _:
                    identifier = "minecraft:bread"; metadata = 0; return true;
                case Item.Brick _:
                    identifier = "minecraft:brick"; metadata = 0; return true;
                case Item.CarrotOnAStick _:
                    identifier = "minecraft:carrot_on_a_stick"; metadata = 0; return true;
                case Item.Charcoal _:
                    identifier = "minecraft:charcoal"; metadata = 0; return true;
                case Item.Chicken { Cooked: false }:
                    identifier = "minecraft:chicken"; metadata = 0; return true;
                case Item.Chicken { Cooked: true }:
                    identifier = "minecraft:cooked_chicken"; metadata = 0; return true;
                case Item.ClayBall _:
                    identifier = "minecraft:clay_ball"; metadata = 0; return true;
                case Item.Clock _:
                    identifier = "minecraft:clock"; metadata = 0; return true;
                case Item.Coal _:
                    identifier = "minecraft:coal"; metadata = 0; return true;
                case Item.Cod { Cooked: false }:
                    identifier = "minecraft:cod"; metadata = 0; return true;
                case Item.Cod { Cooked: true }:
                    identifier = "minecraft:cooked_cod"; metadata = 0; return true;
                case Item.Compass _:
                    identifier = "minecraft:compass"; metadata = 0; return true;
                case Item.Cookie _:
                    identifier = "minecraft:cookie"; metadata = 0; return true;
                case Item.CopperIngot _:
                    identifier = "minecraft:copper_ingot"; metadata = 0; return true;
                case Item.CopperNugget _:
                    identifier = "minecraft:copper_nugget"; metadata = 0; return true;
                case Item.Diamond _:
                    identifier = "minecraft:diamond"; metadata = 0; return true;
                case Item.DiscFragment _:
                    identifier = "minecraft:disc_fragment_5"; metadata = 0; return true;
                case Item.DragonBreath _:
                    identifier = "minecraft:dragon_breath"; metadata = 0; return true;
                case Item.DriedKelp _:
                    identifier = "minecraft:dried_kelp"; metadata = 0; return true;
                case Item.EchoShard _:
                    identifier = "minecraft:echo_shard"; metadata = 0; return true;
                case Item.Egg _:
                    identifier = "minecraft:egg"; metadata = 0; return true;
                case Item.Elytra _:
                    identifier = "minecraft:elytra"; metadata = 0; return true;
                case Item.Emerald _:
                    identifier = "minecraft:emerald"; metadata = 0; return true;
                case Item.EnchantedApple _:
                    identifier = "minecraft:enchanted_golden_apple"; metadata = 0; return true;
                case Item.EnchantedBook _:
                    identifier = "minecraft:enchanted_book"; metadata = 0; return true;
                case Item.EnderEye _:
                    identifier = "minecraft:ender_eye"; metadata = 0; return true;
                case Item.EnderPearl _:
                    identifier = "minecraft:ender_pearl"; metadata = 0; return true;
                case Item.Feather _:
                    identifier = "minecraft:feather"; metadata = 0; return true;
                case Item.FermentedSpiderEye _:
                    identifier = "minecraft:fermented_spider_eye"; metadata = 0; return true;
                case Item.FireCharge _:
                    identifier = "minecraft:fire_charge"; metadata = 0; return true;
                case Item.Flint _:
                    identifier = "minecraft:flint"; metadata = 0; return true;
                case Item.FlintAndSteel _:
                    identifier = "minecraft:flint_and_steel"; metadata = 0; return true;
                case Item.GhastTear _:
                    identifier = "minecraft:ghast_tear"; metadata = 0; return true;
                case Item.GlassBottle _:
                    identifier = "minecraft:glass_bottle"; metadata = 0; return true;
                case Item.GlisteringMelonSlice _:
                    identifier = "minecraft:glistering_melon_slice"; metadata = 0; return true;
                case Item.GlowstoneDust _:
                    identifier = "minecraft:glowstone_dust"; metadata = 0; return true;
                case Item.GoldIngot _:
                    identifier = "minecraft:gold_ingot"; metadata = 0; return true;
                case Item.GoldNugget _:
                    identifier = "minecraft:gold_nugget"; metadata = 0; return true;
                case Item.GoldenApple _:
                    identifier = "minecraft:golden_apple"; metadata = 0; return true;
                case Item.GoldenCarrot _:
                    identifier = "minecraft:golden_carrot"; metadata = 0; return true;
                case Item.Gunpowder _:
                    identifier = "minecraft:gunpowder"; metadata = 0; return true;
                case Item.HeartOfTheSea _:
                    identifier = "minecraft:heart_of_the_sea"; metadata = 0; return true;
                case Item.Hoe value when value.Tier == Item.ToolTierWood:
                    identifier = "minecraft:wooden_hoe"; metadata = 0; return true;
                case Item.Hoe value when value.Tier == Item.ToolTierGold:
                    identifier = "minecraft:golden_hoe"; metadata = 0; return true;
                case Item.Hoe value when value.Tier == Item.ToolTierStone:
                    identifier = "minecraft:stone_hoe"; metadata = 0; return true;
                case Item.Hoe value when value.Tier == Item.ToolTierCopper:
                    identifier = "minecraft:copper_hoe"; metadata = 0; return true;
                case Item.Hoe value when value.Tier == Item.ToolTierIron:
                    identifier = "minecraft:iron_hoe"; metadata = 0; return true;
                case Item.Hoe value when value.Tier == Item.ToolTierDiamond:
                    identifier = "minecraft:diamond_hoe"; metadata = 0; return true;
                case Item.Hoe value when value.Tier == Item.ToolTierNetherite:
                    identifier = "minecraft:netherite_hoe"; metadata = 0; return true;
                case Item.HoneyBottle _:
                    identifier = "minecraft:honey_bottle"; metadata = 0; return true;
                case Item.Honeycomb _:
                    identifier = "minecraft:honeycomb"; metadata = 0; return true;
                case Item.InkSac { Glowing: false }:
                    identifier = "minecraft:ink_sac"; metadata = 0; return true;
                case Item.InkSac { Glowing: true }:
                    identifier = "minecraft:glow_ink_sac"; metadata = 0; return true;
                case Item.IronIngot _:
                    identifier = "minecraft:iron_ingot"; metadata = 0; return true;
                case Item.IronNugget _:
                    identifier = "minecraft:iron_nugget"; metadata = 0; return true;
                case Item.LapisLazuli _:
                    identifier = "minecraft:lapis_lazuli"; metadata = 0; return true;
                case Item.Leather _:
                    identifier = "minecraft:leather"; metadata = 0; return true;
                case Item.MagmaCream _:
                    identifier = "minecraft:magma_cream"; metadata = 0; return true;
                case Item.MelonSlice _:
                    identifier = "minecraft:melon_slice"; metadata = 0; return true;
                case Item.MushroomStew _:
                    identifier = "minecraft:mushroom_stew"; metadata = 0; return true;
                case Item.Mutton { Cooked: false }:
                    identifier = "minecraft:mutton"; metadata = 0; return true;
                case Item.Mutton { Cooked: true }:
                    identifier = "minecraft:cooked_mutton"; metadata = 0; return true;
                case Item.NautilusShell _:
                    identifier = "minecraft:nautilus_shell"; metadata = 0; return true;
                case Item.NetherBrick _:
                    identifier = "minecraft:netherbrick"; metadata = 0; return true;
                case Item.NetherQuartz _:
                    identifier = "minecraft:quartz"; metadata = 0; return true;
                case Item.NetherStar _:
                    identifier = "minecraft:nether_star"; metadata = 0; return true;
                case Item.NetheriteIngot _:
                    identifier = "minecraft:netherite_ingot"; metadata = 0; return true;
                case Item.NetheriteScrap _:
                    identifier = "minecraft:netherite_scrap"; metadata = 0; return true;
                case Item.Paper _:
                    identifier = "minecraft:paper"; metadata = 0; return true;
                case Item.PhantomMembrane _:
                    identifier = "minecraft:phantom_membrane"; metadata = 0; return true;
                case Item.Pickaxe value when value.Tier == Item.ToolTierWood:
                    identifier = "minecraft:wooden_pickaxe"; metadata = 0; return true;
                case Item.Pickaxe value when value.Tier == Item.ToolTierGold:
                    identifier = "minecraft:golden_pickaxe"; metadata = 0; return true;
                case Item.Pickaxe value when value.Tier == Item.ToolTierStone:
                    identifier = "minecraft:stone_pickaxe"; metadata = 0; return true;
                case Item.Pickaxe value when value.Tier == Item.ToolTierCopper:
                    identifier = "minecraft:copper_pickaxe"; metadata = 0; return true;
                case Item.Pickaxe value when value.Tier == Item.ToolTierIron:
                    identifier = "minecraft:iron_pickaxe"; metadata = 0; return true;
                case Item.Pickaxe value when value.Tier == Item.ToolTierDiamond:
                    identifier = "minecraft:diamond_pickaxe"; metadata = 0; return true;
                case Item.Pickaxe value when value.Tier == Item.ToolTierNetherite:
                    identifier = "minecraft:netherite_pickaxe"; metadata = 0; return true;
                case Item.PoisonousPotato _:
                    identifier = "minecraft:poisonous_potato"; metadata = 0; return true;
                case Item.PoppedChorusFruit _:
                    identifier = "minecraft:popped_chorus_fruit"; metadata = 0; return true;
                case Item.Porkchop { Cooked: false }:
                    identifier = "minecraft:porkchop"; metadata = 0; return true;
                case Item.Porkchop { Cooked: true }:
                    identifier = "minecraft:cooked_porkchop"; metadata = 0; return true;
                case Item.PrismarineCrystals _:
                    identifier = "minecraft:prismarine_crystals"; metadata = 0; return true;
                case Item.PrismarineShard _:
                    identifier = "minecraft:prismarine_shard"; metadata = 0; return true;
                case Item.Pufferfish _:
                    identifier = "minecraft:pufferfish"; metadata = 0; return true;
                case Item.PumpkinPie _:
                    identifier = "minecraft:pumpkin_pie"; metadata = 0; return true;
                case Item.Rabbit { Cooked: false }:
                    identifier = "minecraft:rabbit"; metadata = 0; return true;
                case Item.Rabbit { Cooked: true }:
                    identifier = "minecraft:cooked_rabbit"; metadata = 0; return true;
                case Item.RabbitFoot _:
                    identifier = "minecraft:rabbit_foot"; metadata = 0; return true;
                case Item.RabbitHide _:
                    identifier = "minecraft:rabbit_hide"; metadata = 0; return true;
                case Item.RabbitStew _:
                    identifier = "minecraft:rabbit_stew"; metadata = 0; return true;
                case Item.RawCopper _:
                    identifier = "minecraft:raw_copper"; metadata = 0; return true;
                case Item.RawGold _:
                    identifier = "minecraft:raw_gold"; metadata = 0; return true;
                case Item.RawIron _:
                    identifier = "minecraft:raw_iron"; metadata = 0; return true;
                case Item.RecoveryCompass _:
                    identifier = "minecraft:recovery_compass"; metadata = 0; return true;
                case Item.ResinBrick _:
                    identifier = "minecraft:resin_brick"; metadata = 0; return true;
                case Item.RottenFlesh _:
                    identifier = "minecraft:rotten_flesh"; metadata = 0; return true;
                case Item.Salmon { Cooked: false }:
                    identifier = "minecraft:salmon"; metadata = 0; return true;
                case Item.Salmon { Cooked: true }:
                    identifier = "minecraft:cooked_salmon"; metadata = 0; return true;
                case Item.Scute _:
                    identifier = "minecraft:turtle_scute"; metadata = 0; return true;
                case Item.Shears _:
                    identifier = "minecraft:shears"; metadata = 0; return true;
                case Item.Shovel value when value.Tier == Item.ToolTierWood:
                    identifier = "minecraft:wooden_shovel"; metadata = 0; return true;
                case Item.Shovel value when value.Tier == Item.ToolTierGold:
                    identifier = "minecraft:golden_shovel"; metadata = 0; return true;
                case Item.Shovel value when value.Tier == Item.ToolTierStone:
                    identifier = "minecraft:stone_shovel"; metadata = 0; return true;
                case Item.Shovel value when value.Tier == Item.ToolTierCopper:
                    identifier = "minecraft:copper_shovel"; metadata = 0; return true;
                case Item.Shovel value when value.Tier == Item.ToolTierIron:
                    identifier = "minecraft:iron_shovel"; metadata = 0; return true;
                case Item.Shovel value when value.Tier == Item.ToolTierDiamond:
                    identifier = "minecraft:diamond_shovel"; metadata = 0; return true;
                case Item.Shovel value when value.Tier == Item.ToolTierNetherite:
                    identifier = "minecraft:netherite_shovel"; metadata = 0; return true;
                case Item.ShulkerShell _:
                    identifier = "minecraft:shulker_shell"; metadata = 0; return true;
                case Item.Slimeball _:
                    identifier = "minecraft:slime_ball"; metadata = 0; return true;
                case Item.Snowball _:
                    identifier = "minecraft:snowball"; metadata = 0; return true;
                case Item.SpiderEye _:
                    identifier = "minecraft:spider_eye"; metadata = 0; return true;
                case Item.Spyglass _:
                    identifier = "minecraft:spyglass"; metadata = 0; return true;
                case Item.Stick _:
                    identifier = "minecraft:stick"; metadata = 0; return true;
                case Item.Sugar _:
                    identifier = "minecraft:sugar"; metadata = 0; return true;
                case Item.Sword value when value.Tier == Item.ToolTierWood:
                    identifier = "minecraft:wooden_sword"; metadata = 0; return true;
                case Item.Sword value when value.Tier == Item.ToolTierGold:
                    identifier = "minecraft:golden_sword"; metadata = 0; return true;
                case Item.Sword value when value.Tier == Item.ToolTierStone:
                    identifier = "minecraft:stone_sword"; metadata = 0; return true;
                case Item.Sword value when value.Tier == Item.ToolTierCopper:
                    identifier = "minecraft:copper_sword"; metadata = 0; return true;
                case Item.Sword value when value.Tier == Item.ToolTierIron:
                    identifier = "minecraft:iron_sword"; metadata = 0; return true;
                case Item.Sword value when value.Tier == Item.ToolTierDiamond:
                    identifier = "minecraft:diamond_sword"; metadata = 0; return true;
                case Item.Sword value when value.Tier == Item.ToolTierNetherite:
                    identifier = "minecraft:netherite_sword"; metadata = 0; return true;
                case Item.Totem _:
                    identifier = "minecraft:totem_of_undying"; metadata = 0; return true;
                case Item.TropicalFish _:
                    identifier = "minecraft:tropical_fish"; metadata = 0; return true;
                case Item.TurtleShell _:
                    identifier = "minecraft:turtle_helmet"; metadata = 0; return true;
                case Item.WarpedFungusOnAStick _:
                    identifier = "minecraft:warped_fungus_on_a_stick"; metadata = 0; return true;
                case Item.Wheat _:
                    identifier = "minecraft:wheat"; metadata = 0; return true;
                case EncodedItem encoded:
                    identifier = encoded.Identifier; metadata = encoded.Metadata; return true;
                default:
                    identifier = string.Empty; metadata = 0; return false;
            }
        }

        internal static World.Item Decode(string identifier, int metadata)
        {
            if (identifier == "minecraft:amethyst_shard" && metadata == 0) return new Item.AmethystShard();
            if (identifier == "minecraft:apple" && metadata == 0) return new Item.Apple();
            if (identifier == "minecraft:wooden_axe" && metadata == 0) return new Item.Axe(Item.ToolTierWood);
            if (identifier == "minecraft:golden_axe" && metadata == 0) return new Item.Axe(Item.ToolTierGold);
            if (identifier == "minecraft:stone_axe" && metadata == 0) return new Item.Axe(Item.ToolTierStone);
            if (identifier == "minecraft:copper_axe" && metadata == 0) return new Item.Axe(Item.ToolTierCopper);
            if (identifier == "minecraft:iron_axe" && metadata == 0) return new Item.Axe(Item.ToolTierIron);
            if (identifier == "minecraft:diamond_axe" && metadata == 0) return new Item.Axe(Item.ToolTierDiamond);
            if (identifier == "minecraft:netherite_axe" && metadata == 0) return new Item.Axe(Item.ToolTierNetherite);
            if (identifier == "minecraft:baked_potato" && metadata == 0) return new Item.BakedPotato();
            if (identifier == "minecraft:beef" && metadata == 0) return new Item.Beef(false);
            if (identifier == "minecraft:cooked_beef" && metadata == 0) return new Item.Beef(true);
            if (identifier == "minecraft:beetroot" && metadata == 0) return new Item.Beetroot();
            if (identifier == "minecraft:beetroot_soup" && metadata == 0) return new Item.BeetrootSoup();
            if (identifier == "minecraft:blaze_powder" && metadata == 0) return new Item.BlazePowder();
            if (identifier == "minecraft:blaze_rod" && metadata == 0) return new Item.BlazeRod();
            if (identifier == "minecraft:bone" && metadata == 0) return new Item.Bone();
            if (identifier == "minecraft:bone_meal" && metadata == 0) return new Item.BoneMeal();
            if (identifier == "minecraft:book" && metadata == 0) return new Item.Book();
            if (identifier == "minecraft:experience_bottle" && metadata == 0) return new Item.BottleOfEnchanting();
            if (identifier == "minecraft:bow" && metadata == 0) return new Item.Bow();
            if (identifier == "minecraft:bowl" && metadata == 0) return new Item.Bowl();
            if (identifier == "minecraft:bread" && metadata == 0) return new Item.Bread();
            if (identifier == "minecraft:brick" && metadata == 0) return new Item.Brick();
            if (identifier == "minecraft:carrot_on_a_stick" && metadata == 0) return new Item.CarrotOnAStick();
            if (identifier == "minecraft:charcoal" && metadata == 0) return new Item.Charcoal();
            if (identifier == "minecraft:chicken" && metadata == 0) return new Item.Chicken(false);
            if (identifier == "minecraft:cooked_chicken" && metadata == 0) return new Item.Chicken(true);
            if (identifier == "minecraft:clay_ball" && metadata == 0) return new Item.ClayBall();
            if (identifier == "minecraft:clock" && metadata == 0) return new Item.Clock();
            if (identifier == "minecraft:coal" && metadata == 0) return new Item.Coal();
            if (identifier == "minecraft:cod" && metadata == 0) return new Item.Cod(false);
            if (identifier == "minecraft:cooked_cod" && metadata == 0) return new Item.Cod(true);
            if (identifier == "minecraft:compass" && metadata == 0) return new Item.Compass();
            if (identifier == "minecraft:cookie" && metadata == 0) return new Item.Cookie();
            if (identifier == "minecraft:copper_ingot" && metadata == 0) return new Item.CopperIngot();
            if (identifier == "minecraft:copper_nugget" && metadata == 0) return new Item.CopperNugget();
            if (identifier == "minecraft:diamond" && metadata == 0) return new Item.Diamond();
            if (identifier == "minecraft:disc_fragment_5" && metadata == 0) return new Item.DiscFragment();
            if (identifier == "minecraft:dragon_breath" && metadata == 0) return new Item.DragonBreath();
            if (identifier == "minecraft:dried_kelp" && metadata == 0) return new Item.DriedKelp();
            if (identifier == "minecraft:echo_shard" && metadata == 0) return new Item.EchoShard();
            if (identifier == "minecraft:egg" && metadata == 0) return new Item.Egg();
            if (identifier == "minecraft:elytra" && metadata == 0) return new Item.Elytra();
            if (identifier == "minecraft:emerald" && metadata == 0) return new Item.Emerald();
            if (identifier == "minecraft:enchanted_golden_apple" && metadata == 0) return new Item.EnchantedApple();
            if (identifier == "minecraft:enchanted_book" && metadata == 0) return new Item.EnchantedBook();
            if (identifier == "minecraft:ender_eye" && metadata == 0) return new Item.EnderEye();
            if (identifier == "minecraft:ender_pearl" && metadata == 0) return new Item.EnderPearl();
            if (identifier == "minecraft:feather" && metadata == 0) return new Item.Feather();
            if (identifier == "minecraft:fermented_spider_eye" && metadata == 0) return new Item.FermentedSpiderEye();
            if (identifier == "minecraft:fire_charge" && metadata == 0) return new Item.FireCharge();
            if (identifier == "minecraft:flint" && metadata == 0) return new Item.Flint();
            if (identifier == "minecraft:flint_and_steel" && metadata == 0) return new Item.FlintAndSteel();
            if (identifier == "minecraft:ghast_tear" && metadata == 0) return new Item.GhastTear();
            if (identifier == "minecraft:glass_bottle" && metadata == 0) return new Item.GlassBottle();
            if (identifier == "minecraft:glistering_melon_slice" && metadata == 0) return new Item.GlisteringMelonSlice();
            if (identifier == "minecraft:glowstone_dust" && metadata == 0) return new Item.GlowstoneDust();
            if (identifier == "minecraft:gold_ingot" && metadata == 0) return new Item.GoldIngot();
            if (identifier == "minecraft:gold_nugget" && metadata == 0) return new Item.GoldNugget();
            if (identifier == "minecraft:golden_apple" && metadata == 0) return new Item.GoldenApple();
            if (identifier == "minecraft:golden_carrot" && metadata == 0) return new Item.GoldenCarrot();
            if (identifier == "minecraft:gunpowder" && metadata == 0) return new Item.Gunpowder();
            if (identifier == "minecraft:heart_of_the_sea" && metadata == 0) return new Item.HeartOfTheSea();
            if (identifier == "minecraft:wooden_hoe" && metadata == 0) return new Item.Hoe(Item.ToolTierWood);
            if (identifier == "minecraft:golden_hoe" && metadata == 0) return new Item.Hoe(Item.ToolTierGold);
            if (identifier == "minecraft:stone_hoe" && metadata == 0) return new Item.Hoe(Item.ToolTierStone);
            if (identifier == "minecraft:copper_hoe" && metadata == 0) return new Item.Hoe(Item.ToolTierCopper);
            if (identifier == "minecraft:iron_hoe" && metadata == 0) return new Item.Hoe(Item.ToolTierIron);
            if (identifier == "minecraft:diamond_hoe" && metadata == 0) return new Item.Hoe(Item.ToolTierDiamond);
            if (identifier == "minecraft:netherite_hoe" && metadata == 0) return new Item.Hoe(Item.ToolTierNetherite);
            if (identifier == "minecraft:honey_bottle" && metadata == 0) return new Item.HoneyBottle();
            if (identifier == "minecraft:honeycomb" && metadata == 0) return new Item.Honeycomb();
            if (identifier == "minecraft:ink_sac" && metadata == 0) return new Item.InkSac(false);
            if (identifier == "minecraft:glow_ink_sac" && metadata == 0) return new Item.InkSac(true);
            if (identifier == "minecraft:iron_ingot" && metadata == 0) return new Item.IronIngot();
            if (identifier == "minecraft:iron_nugget" && metadata == 0) return new Item.IronNugget();
            if (identifier == "minecraft:lapis_lazuli" && metadata == 0) return new Item.LapisLazuli();
            if (identifier == "minecraft:leather" && metadata == 0) return new Item.Leather();
            if (identifier == "minecraft:magma_cream" && metadata == 0) return new Item.MagmaCream();
            if (identifier == "minecraft:melon_slice" && metadata == 0) return new Item.MelonSlice();
            if (identifier == "minecraft:mushroom_stew" && metadata == 0) return new Item.MushroomStew();
            if (identifier == "minecraft:mutton" && metadata == 0) return new Item.Mutton(false);
            if (identifier == "minecraft:cooked_mutton" && metadata == 0) return new Item.Mutton(true);
            if (identifier == "minecraft:nautilus_shell" && metadata == 0) return new Item.NautilusShell();
            if (identifier == "minecraft:netherbrick" && metadata == 0) return new Item.NetherBrick();
            if (identifier == "minecraft:quartz" && metadata == 0) return new Item.NetherQuartz();
            if (identifier == "minecraft:nether_star" && metadata == 0) return new Item.NetherStar();
            if (identifier == "minecraft:netherite_ingot" && metadata == 0) return new Item.NetheriteIngot();
            if (identifier == "minecraft:netherite_scrap" && metadata == 0) return new Item.NetheriteScrap();
            if (identifier == "minecraft:paper" && metadata == 0) return new Item.Paper();
            if (identifier == "minecraft:phantom_membrane" && metadata == 0) return new Item.PhantomMembrane();
            if (identifier == "minecraft:wooden_pickaxe" && metadata == 0) return new Item.Pickaxe(Item.ToolTierWood);
            if (identifier == "minecraft:golden_pickaxe" && metadata == 0) return new Item.Pickaxe(Item.ToolTierGold);
            if (identifier == "minecraft:stone_pickaxe" && metadata == 0) return new Item.Pickaxe(Item.ToolTierStone);
            if (identifier == "minecraft:copper_pickaxe" && metadata == 0) return new Item.Pickaxe(Item.ToolTierCopper);
            if (identifier == "minecraft:iron_pickaxe" && metadata == 0) return new Item.Pickaxe(Item.ToolTierIron);
            if (identifier == "minecraft:diamond_pickaxe" && metadata == 0) return new Item.Pickaxe(Item.ToolTierDiamond);
            if (identifier == "minecraft:netherite_pickaxe" && metadata == 0) return new Item.Pickaxe(Item.ToolTierNetherite);
            if (identifier == "minecraft:poisonous_potato" && metadata == 0) return new Item.PoisonousPotato();
            if (identifier == "minecraft:popped_chorus_fruit" && metadata == 0) return new Item.PoppedChorusFruit();
            if (identifier == "minecraft:porkchop" && metadata == 0) return new Item.Porkchop(false);
            if (identifier == "minecraft:cooked_porkchop" && metadata == 0) return new Item.Porkchop(true);
            if (identifier == "minecraft:prismarine_crystals" && metadata == 0) return new Item.PrismarineCrystals();
            if (identifier == "minecraft:prismarine_shard" && metadata == 0) return new Item.PrismarineShard();
            if (identifier == "minecraft:pufferfish" && metadata == 0) return new Item.Pufferfish();
            if (identifier == "minecraft:pumpkin_pie" && metadata == 0) return new Item.PumpkinPie();
            if (identifier == "minecraft:rabbit" && metadata == 0) return new Item.Rabbit(false);
            if (identifier == "minecraft:cooked_rabbit" && metadata == 0) return new Item.Rabbit(true);
            if (identifier == "minecraft:rabbit_foot" && metadata == 0) return new Item.RabbitFoot();
            if (identifier == "minecraft:rabbit_hide" && metadata == 0) return new Item.RabbitHide();
            if (identifier == "minecraft:rabbit_stew" && metadata == 0) return new Item.RabbitStew();
            if (identifier == "minecraft:raw_copper" && metadata == 0) return new Item.RawCopper();
            if (identifier == "minecraft:raw_gold" && metadata == 0) return new Item.RawGold();
            if (identifier == "minecraft:raw_iron" && metadata == 0) return new Item.RawIron();
            if (identifier == "minecraft:recovery_compass" && metadata == 0) return new Item.RecoveryCompass();
            if (identifier == "minecraft:resin_brick" && metadata == 0) return new Item.ResinBrick();
            if (identifier == "minecraft:rotten_flesh" && metadata == 0) return new Item.RottenFlesh();
            if (identifier == "minecraft:salmon" && metadata == 0) return new Item.Salmon(false);
            if (identifier == "minecraft:cooked_salmon" && metadata == 0) return new Item.Salmon(true);
            if (identifier == "minecraft:turtle_scute" && metadata == 0) return new Item.Scute();
            if (identifier == "minecraft:shears" && metadata == 0) return new Item.Shears();
            if (identifier == "minecraft:wooden_shovel" && metadata == 0) return new Item.Shovel(Item.ToolTierWood);
            if (identifier == "minecraft:golden_shovel" && metadata == 0) return new Item.Shovel(Item.ToolTierGold);
            if (identifier == "minecraft:stone_shovel" && metadata == 0) return new Item.Shovel(Item.ToolTierStone);
            if (identifier == "minecraft:copper_shovel" && metadata == 0) return new Item.Shovel(Item.ToolTierCopper);
            if (identifier == "minecraft:iron_shovel" && metadata == 0) return new Item.Shovel(Item.ToolTierIron);
            if (identifier == "minecraft:diamond_shovel" && metadata == 0) return new Item.Shovel(Item.ToolTierDiamond);
            if (identifier == "minecraft:netherite_shovel" && metadata == 0) return new Item.Shovel(Item.ToolTierNetherite);
            if (identifier == "minecraft:shulker_shell" && metadata == 0) return new Item.ShulkerShell();
            if (identifier == "minecraft:slime_ball" && metadata == 0) return new Item.Slimeball();
            if (identifier == "minecraft:snowball" && metadata == 0) return new Item.Snowball();
            if (identifier == "minecraft:spider_eye" && metadata == 0) return new Item.SpiderEye();
            if (identifier == "minecraft:spyglass" && metadata == 0) return new Item.Spyglass();
            if (identifier == "minecraft:stick" && metadata == 0) return new Item.Stick();
            if (identifier == "minecraft:sugar" && metadata == 0) return new Item.Sugar();
            if (identifier == "minecraft:wooden_sword" && metadata == 0) return new Item.Sword(Item.ToolTierWood);
            if (identifier == "minecraft:golden_sword" && metadata == 0) return new Item.Sword(Item.ToolTierGold);
            if (identifier == "minecraft:stone_sword" && metadata == 0) return new Item.Sword(Item.ToolTierStone);
            if (identifier == "minecraft:copper_sword" && metadata == 0) return new Item.Sword(Item.ToolTierCopper);
            if (identifier == "minecraft:iron_sword" && metadata == 0) return new Item.Sword(Item.ToolTierIron);
            if (identifier == "minecraft:diamond_sword" && metadata == 0) return new Item.Sword(Item.ToolTierDiamond);
            if (identifier == "minecraft:netherite_sword" && metadata == 0) return new Item.Sword(Item.ToolTierNetherite);
            if (identifier == "minecraft:totem_of_undying" && metadata == 0) return new Item.Totem();
            if (identifier == "minecraft:tropical_fish" && metadata == 0) return new Item.TropicalFish();
            if (identifier == "minecraft:turtle_helmet" && metadata == 0) return new Item.TurtleShell();
            if (identifier == "minecraft:warped_fungus_on_a_stick" && metadata == 0) return new Item.WarpedFungusOnAStick();
            if (identifier == "minecraft:wheat" && metadata == 0) return new Item.Wheat();
            return new EncodedItem(identifier, metadata);
        }

        internal static bool IsAir(World.Item item) =>
            TryEncode(item, out var identifier, out _) && identifier == "minecraft:air";

        private sealed record EncodedItem(string Identifier, int Metadata) : World.Item;
    }
}
