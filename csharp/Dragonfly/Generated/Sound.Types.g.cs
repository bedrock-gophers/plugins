// Code generated from Dragonfly server/world/sound Go AST. DO NOT EDIT.
#nullable enable

namespace Dragonfly;

public static partial class Sound
{
    public readonly record struct AnvilBreak : World.Sound;
    public readonly record struct AnvilLand : World.Sound;
    public readonly record struct AnvilUse : World.Sound;
    public readonly record struct ArrowHit : World.Sound;
    public readonly record struct BarrelClose : World.Sound;
    public readonly record struct BarrelOpen : World.Sound;
    public readonly record struct BlastFurnaceCrackle : World.Sound;
    public readonly record struct BowShoot : World.Sound;
    public readonly record struct Burning : World.Sound;
    public readonly record struct Burp : World.Sound;
    public readonly record struct CampfireCrackle : World.Sound;
    public readonly record struct ChestClose : World.Sound;
    public readonly record struct ChestOpen : World.Sound;
    public readonly record struct Click : World.Sound;
    public readonly record struct ComposterEmpty : World.Sound;
    public readonly record struct ComposterFill : World.Sound;
    public readonly record struct ComposterFillLayer : World.Sound;
    public readonly record struct ComposterReady : World.Sound;
    public readonly record struct CopperScraped : World.Sound;
    public readonly record struct CrossbowShoot : World.Sound;
    public readonly record struct DecoratedPotInsertFailed : World.Sound;
    public readonly record struct Deny : World.Sound;
    public readonly record struct DoorCrash : World.Sound;
    public readonly record struct Drowning : World.Sound;
    public readonly record struct EnderChestClose : World.Sound;
    public readonly record struct EnderChestOpen : World.Sound;
    public readonly record struct Experience : World.Sound;
    public readonly record struct Explosion : World.Sound;
    public readonly record struct FireCharge : World.Sound;
    public readonly record struct FireExtinguish : World.Sound;
    public readonly record struct FireworkBlast : World.Sound;
    public readonly record struct FireworkHugeBlast : World.Sound;
    public readonly record struct FireworkLaunch : World.Sound;
    public readonly record struct FireworkTwinkle : World.Sound;
    public readonly record struct Fizz : World.Sound;
    public readonly record struct FurnaceCrackle : World.Sound;
    public readonly record struct GhastShoot : World.Sound;
    public readonly record struct GhastWarning : World.Sound;
    public readonly record struct GlassBreak : World.Sound;
    public readonly record struct Ignite : World.Sound;
    public readonly record struct ItemAdd : World.Sound;
    public readonly record struct ItemBreak : World.Sound;
    public readonly record struct ItemFrameRemove : World.Sound;
    public readonly record struct ItemFrameRotate : World.Sound;
    public readonly record struct ItemThrow : World.Sound;
    public readonly record struct LecternBookPlace : World.Sound;
    public readonly record struct LevelUp : World.Sound;
    public readonly record struct LightningExplode : World.Sound;
    public readonly record struct LightningThunder : World.Sound;
    public readonly record struct MusicDiscEnd : World.Sound;
    public readonly record struct Pop : World.Sound;
    public readonly record struct PotionBrewed : World.Sound;
    public readonly record struct PowerOff : World.Sound;
    public readonly record struct PowerOn : World.Sound;
    public readonly record struct SignWaxed : World.Sound;
    public readonly record struct SmokerCrackle : World.Sound;
    public readonly record struct StopUsingSpyglass : World.Sound;
    public readonly record struct TNT : World.Sound;
    public readonly record struct Teleport : World.Sound;
    public readonly record struct Thunder : World.Sound;
    public readonly record struct Totem : World.Sound;
    public readonly record struct UseSpyglass : World.Sound;
    public readonly record struct WaxRemoved : World.Sound;
    public readonly record struct WaxedSignFailedInteraction : World.Sound;
    public readonly record struct ShulkerBoxOpen : World.Sound;
    public readonly record struct ShulkerBoxClose : World.Sound;
    public readonly record struct EnderEyePlaced : World.Sound;
    public readonly record struct EndPortalCreated : World.Sound;
    public readonly record struct Attack(bool Damage) : World.Sound;
    public readonly record struct Fall(double Distance) : World.Sound;
    public readonly record struct BlockPlace(World.Block Block) : World.Sound;
    public readonly record struct BlockBreaking(World.Block Block) : World.Sound;
    public readonly record struct DoorOpen(World.Block Block) : World.Sound;
    public readonly record struct DoorClose(World.Block Block) : World.Sound;
    public readonly record struct TrapdoorOpen(World.Block Block) : World.Sound;
    public readonly record struct TrapdoorClose(World.Block Block) : World.Sound;
    public readonly record struct FenceGateOpen(World.Block Block) : World.Sound;
    public readonly record struct FenceGateClose(World.Block Block) : World.Sound;
    public readonly record struct Note(Instrument Instrument, int Pitch) : World.Sound;
    public readonly record struct MusicDiscPlay(DiscType DiscType) : World.Sound;
    public readonly record struct DecoratedPotInserted(double Progress) : World.Sound;
    public readonly record struct ItemUseOn(World.Block Block) : World.Sound;
    public readonly record struct EquipItem(World.Item Item) : World.Sound;
    public readonly record struct BucketFill(World.Liquid Liquid) : World.Sound;
    public readonly record struct BucketEmpty(World.Liquid Liquid) : World.Sound;
    public readonly record struct CrossbowLoad(int Stage, bool QuickCharge) : World.Sound;
    public readonly record struct GoatHorn(Horn Horn) : World.Sound;

    internal static World.Sound DecodeEvent(
        uint kind, uint data, int integer, uint flags, double scalar, World.Block? block, World.Item? item) =>
        kind switch
        {
            0 => new AnvilBreak(),
            1 => new AnvilLand(),
            2 => new AnvilUse(),
            3 => new ArrowHit(),
            4 => new BarrelClose(),
            5 => new BarrelOpen(),
            6 => new BlastFurnaceCrackle(),
            7 => new BowShoot(),
            8 => new Burning(),
            9 => new Burp(),
            10 => new CampfireCrackle(),
            11 => new ChestClose(),
            12 => new ChestOpen(),
            13 => new Click(),
            14 => new ComposterEmpty(),
            15 => new ComposterFill(),
            16 => new ComposterFillLayer(),
            17 => new ComposterReady(),
            18 => new CopperScraped(),
            19 => new CrossbowShoot(),
            20 => new DecoratedPotInsertFailed(),
            21 => new Deny(),
            22 => new DoorCrash(),
            23 => new Drowning(),
            24 => new EnderChestClose(),
            25 => new EnderChestOpen(),
            26 => new Experience(),
            27 => new Explosion(),
            28 => new FireCharge(),
            29 => new FireExtinguish(),
            30 => new FireworkBlast(),
            31 => new FireworkHugeBlast(),
            32 => new FireworkLaunch(),
            33 => new FireworkTwinkle(),
            34 => new Fizz(),
            35 => new FurnaceCrackle(),
            36 => new GhastShoot(),
            37 => new GhastWarning(),
            38 => new GlassBreak(),
            39 => new Ignite(),
            40 => new ItemAdd(),
            41 => new ItemBreak(),
            42 => new ItemFrameRemove(),
            43 => new ItemFrameRotate(),
            44 => new ItemThrow(),
            45 => new LecternBookPlace(),
            46 => new LevelUp(),
            47 => new LightningExplode(),
            48 => new LightningThunder(),
            49 => new MusicDiscEnd(),
            50 => new Pop(),
            51 => new PotionBrewed(),
            52 => new PowerOff(),
            53 => new PowerOn(),
            54 => new SignWaxed(),
            55 => new SmokerCrackle(),
            56 => new StopUsingSpyglass(),
            57 => new TNT(),
            58 => new Teleport(),
            59 => new Thunder(),
            60 => new Totem(),
            61 => new UseSpyglass(),
            62 => new WaxRemoved(),
            63 => new WaxedSignFailedInteraction(),
            64 => new ShulkerBoxOpen(),
            65 => new ShulkerBoxClose(),
            66 => new EnderEyePlaced(),
            67 => new EndPortalCreated(),
            68 => new Attack(flags != 0),
            69 => new Fall(scalar),
            70 => new BlockPlace(block ?? throw new InvalidOperationException("Sound requires a block.")),
            71 => new BlockBreaking(block ?? throw new InvalidOperationException("Sound requires a block.")),
            72 => new DoorOpen(block ?? throw new InvalidOperationException("Sound requires a block.")),
            73 => new DoorClose(block ?? throw new InvalidOperationException("Sound requires a block.")),
            74 => new TrapdoorOpen(block ?? throw new InvalidOperationException("Sound requires a block.")),
            75 => new TrapdoorClose(block ?? throw new InvalidOperationException("Sound requires a block.")),
            76 => new FenceGateOpen(block ?? throw new InvalidOperationException("Sound requires a block.")),
            77 => new FenceGateClose(block ?? throw new InvalidOperationException("Sound requires a block.")),
            78 => new Note(new Instrument(data), integer),
            79 => new MusicDiscPlay(new DiscType(checked((int)data))),
            80 => new DecoratedPotInserted(scalar),
            81 => new ItemUseOn(block ?? throw new InvalidOperationException("Sound requires a block.")),
            82 => new EquipItem(item ?? throw new InvalidOperationException("Sound requires an item.")),
            83 => new BucketFill(block is World.Liquid liquid ? liquid : throw new InvalidOperationException("Sound requires a liquid.")),
            84 => new BucketEmpty(block is World.Liquid liquid ? liquid : throw new InvalidOperationException("Sound requires a liquid.")),
            85 => new CrossbowLoad(integer, flags != 0),
            86 => new GoatHorn(new Horn(checked((int)data))),
            _ => throw new InvalidOperationException("Invalid sound kind."),
        };
}
