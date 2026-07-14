// Code generated from Dragonfly server/block Go AST and registry. DO NOT EDIT.
#nullable enable
using System;
using Dragonfly;

namespace Dragonfly
{
    public static partial class Block
    {
        public readonly record struct Air : World.Block;
        public readonly record struct Amethyst : World.Block;
        public readonly record struct AncientDebris : World.Block;
        public readonly record struct BambooMosaic : World.Block;
        public readonly record struct Barrier : World.Block;
        public readonly record struct BeetrootSeeds : World.Block;
        public readonly record struct BlueIce : World.Block;
        public readonly record struct Bookshelf : World.Block;
        public readonly record struct Bricks : World.Block;
        public readonly record struct Calcite : World.Block;
        public readonly record struct Carrot : World.Block;
        public readonly record struct ChiseledQuartz : World.Block;
        public readonly record struct Clay : World.Block;
        public readonly record struct Coal : World.Block;
        public readonly record struct Cobweb : World.Block;
        public readonly record struct CraftingTable : World.Block;
        public readonly record struct DeadBush : World.Block;
        public readonly record struct Diamond : World.Block;
        public readonly record struct DirtPath : World.Block;
        public readonly record struct DragonEgg : World.Block;
        public readonly record struct DriedKelp : World.Block;
        public readonly record struct Dripstone : World.Block;
        public readonly record struct Emerald : World.Block;
        public readonly record struct EnchantingTable : World.Block;
        public readonly record struct EndBricks : World.Block;
        public readonly record struct EndPortal : World.Block;
        public readonly record struct EndStone : World.Block;
        public readonly record struct Fern : World.Block;
        public readonly record struct FletchingTable : World.Block;
        public readonly record struct Glass : World.Block;
        public readonly record struct GlassPane : World.Block;
        public readonly record struct Glowstone : World.Block;
        public readonly record struct Gold : World.Block;
        public readonly record struct Grass : World.Block;
        public readonly record struct Gravel : World.Block;
        public readonly record struct Honeycomb : World.Block;
        public readonly record struct InfestedCobblestone : World.Block;
        public readonly record struct InfestedStone : World.Block;
        public readonly record struct InvisibleBedrock : World.Block;
        public readonly record struct Iron : World.Block;
        public readonly record struct IronBars : World.Block;
        public readonly record struct Lapis : World.Block;
        public readonly record struct LilyPad : World.Block;
        public readonly record struct Magma : World.Block;
        public readonly record struct Melon : World.Block;
        public readonly record struct MossCarpet : World.Block;
        public readonly record struct Mud : World.Block;
        public readonly record struct MudBricks : World.Block;
        public readonly record struct NetherBrickFence : World.Block;
        public readonly record struct NetherGoldOre : World.Block;
        public readonly record struct NetherQuartzOre : World.Block;
        public readonly record struct NetherSprouts : World.Block;
        public readonly record struct Netherite : World.Block;
        public readonly record struct Netherrack : World.Block;
        public readonly record struct PackedIce : World.Block;
        public readonly record struct PackedMud : World.Block;
        public readonly record struct Podzol : World.Block;
        public readonly record struct PolishedTuff : World.Block;
        public readonly record struct Potato : World.Block;
        public readonly record struct Purpur : World.Block;
        public readonly record struct QuartzBricks : World.Block;
        public readonly record struct RawCopper : World.Block;
        public readonly record struct RawGold : World.Block;
        public readonly record struct RawIron : World.Block;
        public readonly record struct RedstoneBlock : World.Block;
        public readonly record struct ReinforcedDeepslate : World.Block;
        public readonly record struct Resin : World.Block;
        public readonly record struct SeaLantern : World.Block;
        public readonly record struct Shroomlight : World.Block;
        public readonly record struct Slime : World.Block;
        public readonly record struct SmithingTable : World.Block;
        public readonly record struct SmoothBasalt : World.Block;
        public readonly record struct Snow : World.Block;
        public readonly record struct SoulSand : World.Block;
        public readonly record struct SoulSoil : World.Block;
        public readonly record struct SporeBlossom : World.Block;
        public readonly record struct TNT : World.Block;
        public readonly record struct Terracotta : World.Block;
        public readonly record struct WheatSeeds : World.Block;
        public readonly record struct Sand(bool Red = false) : World.Block;
        public readonly record struct Water(bool Still, int Depth, bool Falling) : World.Liquid
        {
            public string LiquidType() => "water";
        }
        public readonly record struct Lava(bool Still, int Depth, bool Falling) : World.Liquid
        {
            public string LiquidType() => "lava";
        }
    }

    internal static class BlockCodec
    {
        private static readonly byte[] State0 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State1 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State2 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State3 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State4 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State5 = [0x0a, 0x00, 0x00, 0x0a, 0x06, 0x00, 0x67, 0x72, 0x6f, 0x77, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State6 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State7 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State8 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State9 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State10 = [0x0a, 0x00, 0x00, 0x0a, 0x06, 0x00, 0x67, 0x72, 0x6f, 0x77, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State11 = [0x0a, 0x00, 0x00, 0x0a, 0x0b, 0x00, 0x70, 0x69, 0x6c, 0x6c, 0x61, 0x72, 0x5f, 0x61, 0x78, 0x69, 0x73, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x03, 0x00, 0x00, 0x00, 0x08, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x01, 0x00, 0x79, 0x00, 0x00];
        private static readonly byte[] State12 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State13 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State14 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State15 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State16 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State17 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State18 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State19 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State20 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State21 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State22 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State23 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State24 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State25 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State26 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State27 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State28 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State29 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State30 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State31 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State32 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State33 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State34 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State35 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State36 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State37 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State38 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State39 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State40 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State41 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State42 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State43 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State44 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State45 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State46 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State47 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State48 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State49 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State50 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State51 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State52 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State53 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State54 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State55 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State56 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State57 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State58 = [0x0a, 0x00, 0x00, 0x0a, 0x06, 0x00, 0x67, 0x72, 0x6f, 0x77, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State59 = [0x0a, 0x00, 0x00, 0x0a, 0x0b, 0x00, 0x70, 0x69, 0x6c, 0x6c, 0x61, 0x72, 0x5f, 0x61, 0x78, 0x69, 0x73, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x03, 0x00, 0x00, 0x00, 0x08, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x01, 0x00, 0x79, 0x00, 0x00];
        private static readonly byte[] State60 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State61 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State62 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State63 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State64 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State65 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State66 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State67 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State68 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State69 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State70 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State71 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State72 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State73 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State74 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State75 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State76 = [0x0a, 0x00, 0x00, 0x0a, 0x0b, 0x00, 0x65, 0x78, 0x70, 0x6c, 0x6f, 0x64, 0x65, 0x5f, 0x62, 0x69, 0x74, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x00, 0x00, 0x00, 0x00, 0x01, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x00, 0x00, 0x00];
        private static readonly byte[] State77 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State78 = [0x0a, 0x00, 0x00, 0x0a, 0x06, 0x00, 0x67, 0x72, 0x6f, 0x77, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State79 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State80 = [0x0a, 0x00, 0x00, 0x00];
        private static readonly byte[] State81 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State82 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State83 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State84 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State85 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State86 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State87 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State88 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State89 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State90 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State91 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State92 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0b, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State93 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State94 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State95 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State96 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State97 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State98 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State99 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State100 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State101 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State102 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State103 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State104 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State105 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State106 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State107 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State108 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0b, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State109 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State110 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State111 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State112 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State113 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State114 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State115 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State116 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State117 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State118 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State119 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State120 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State121 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State122 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State123 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State124 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0b, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State125 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State126 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State127 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State128 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State129 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State130 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State131 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State132 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State133 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State134 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State135 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State136 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State137 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State138 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State139 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State140 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0b, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State141 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State142 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State143 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00];
        private static readonly byte[] State144 = [0x0a, 0x00, 0x00, 0x0a, 0x0c, 0x00, 0x6c, 0x69, 0x71, 0x75, 0x69, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x74, 0x68, 0x03, 0x04, 0x00, 0x6b, 0x69, 0x6e, 0x64, 0x02, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00];

        internal static bool TryEncode(World.Block block, out string identifier, out byte[] properties)
        {
            switch (block)
            {
                case Block.Air _:
                    identifier = "minecraft:air"; properties = State0; return true;
                case Block.Amethyst _:
                    identifier = "minecraft:amethyst_block"; properties = State1; return true;
                case Block.AncientDebris _:
                    identifier = "minecraft:ancient_debris"; properties = State2; return true;
                case Block.BambooMosaic _:
                    identifier = "minecraft:bamboo_mosaic"; properties = State3; return true;
                case Block.Barrier _:
                    identifier = "minecraft:barrier"; properties = State4; return true;
                case Block.BeetrootSeeds _:
                    identifier = "minecraft:beetroot"; properties = State5; return true;
                case Block.BlueIce _:
                    identifier = "minecraft:blue_ice"; properties = State6; return true;
                case Block.Bookshelf _:
                    identifier = "minecraft:bookshelf"; properties = State7; return true;
                case Block.Bricks _:
                    identifier = "minecraft:brick_block"; properties = State8; return true;
                case Block.Calcite _:
                    identifier = "minecraft:calcite"; properties = State9; return true;
                case Block.Carrot _:
                    identifier = "minecraft:carrots"; properties = State10; return true;
                case Block.ChiseledQuartz _:
                    identifier = "minecraft:chiseled_quartz_block"; properties = State11; return true;
                case Block.Clay _:
                    identifier = "minecraft:clay"; properties = State12; return true;
                case Block.Coal _:
                    identifier = "minecraft:coal_block"; properties = State13; return true;
                case Block.Cobweb _:
                    identifier = "minecraft:web"; properties = State14; return true;
                case Block.CraftingTable _:
                    identifier = "minecraft:crafting_table"; properties = State15; return true;
                case Block.DeadBush _:
                    identifier = "minecraft:deadbush"; properties = State16; return true;
                case Block.Diamond _:
                    identifier = "minecraft:diamond_block"; properties = State17; return true;
                case Block.DirtPath _:
                    identifier = "minecraft:grass_path"; properties = State18; return true;
                case Block.DragonEgg _:
                    identifier = "minecraft:dragon_egg"; properties = State19; return true;
                case Block.DriedKelp _:
                    identifier = "minecraft:dried_kelp_block"; properties = State20; return true;
                case Block.Dripstone _:
                    identifier = "minecraft:dripstone_block"; properties = State21; return true;
                case Block.Emerald _:
                    identifier = "minecraft:emerald_block"; properties = State22; return true;
                case Block.EnchantingTable _:
                    identifier = "minecraft:enchanting_table"; properties = State23; return true;
                case Block.EndBricks _:
                    identifier = "minecraft:end_bricks"; properties = State24; return true;
                case Block.EndPortal _:
                    identifier = "minecraft:end_portal"; properties = State25; return true;
                case Block.EndStone _:
                    identifier = "minecraft:end_stone"; properties = State26; return true;
                case Block.Fern _:
                    identifier = "minecraft:fern"; properties = State27; return true;
                case Block.FletchingTable _:
                    identifier = "minecraft:fletching_table"; properties = State28; return true;
                case Block.Glass _:
                    identifier = "minecraft:glass"; properties = State29; return true;
                case Block.GlassPane _:
                    identifier = "minecraft:glass_pane"; properties = State30; return true;
                case Block.Glowstone _:
                    identifier = "minecraft:glowstone"; properties = State31; return true;
                case Block.Gold _:
                    identifier = "minecraft:gold_block"; properties = State32; return true;
                case Block.Grass _:
                    identifier = "minecraft:grass_block"; properties = State33; return true;
                case Block.Gravel _:
                    identifier = "minecraft:gravel"; properties = State34; return true;
                case Block.Honeycomb _:
                    identifier = "minecraft:honeycomb_block"; properties = State35; return true;
                case Block.InfestedCobblestone _:
                    identifier = "minecraft:infested_cobblestone"; properties = State36; return true;
                case Block.InfestedStone _:
                    identifier = "minecraft:infested_stone"; properties = State37; return true;
                case Block.InvisibleBedrock _:
                    identifier = "minecraft:invisible_bedrock"; properties = State38; return true;
                case Block.Iron _:
                    identifier = "minecraft:iron_block"; properties = State39; return true;
                case Block.IronBars _:
                    identifier = "minecraft:iron_bars"; properties = State40; return true;
                case Block.Lapis _:
                    identifier = "minecraft:lapis_block"; properties = State41; return true;
                case Block.LilyPad _:
                    identifier = "minecraft:waterlily"; properties = State42; return true;
                case Block.Magma _:
                    identifier = "minecraft:magma"; properties = State43; return true;
                case Block.Melon _:
                    identifier = "minecraft:melon_block"; properties = State44; return true;
                case Block.MossCarpet _:
                    identifier = "minecraft:moss_carpet"; properties = State45; return true;
                case Block.Mud _:
                    identifier = "minecraft:mud"; properties = State46; return true;
                case Block.MudBricks _:
                    identifier = "minecraft:mud_bricks"; properties = State47; return true;
                case Block.NetherBrickFence _:
                    identifier = "minecraft:nether_brick_fence"; properties = State48; return true;
                case Block.NetherGoldOre _:
                    identifier = "minecraft:nether_gold_ore"; properties = State49; return true;
                case Block.NetherQuartzOre _:
                    identifier = "minecraft:quartz_ore"; properties = State50; return true;
                case Block.NetherSprouts _:
                    identifier = "minecraft:nether_sprouts"; properties = State51; return true;
                case Block.Netherite _:
                    identifier = "minecraft:netherite_block"; properties = State52; return true;
                case Block.Netherrack _:
                    identifier = "minecraft:netherrack"; properties = State53; return true;
                case Block.PackedIce _:
                    identifier = "minecraft:packed_ice"; properties = State54; return true;
                case Block.PackedMud _:
                    identifier = "minecraft:packed_mud"; properties = State55; return true;
                case Block.Podzol _:
                    identifier = "minecraft:podzol"; properties = State56; return true;
                case Block.PolishedTuff _:
                    identifier = "minecraft:polished_tuff"; properties = State57; return true;
                case Block.Potato _:
                    identifier = "minecraft:potatoes"; properties = State58; return true;
                case Block.Purpur _:
                    identifier = "minecraft:purpur_block"; properties = State59; return true;
                case Block.QuartzBricks _:
                    identifier = "minecraft:quartz_bricks"; properties = State60; return true;
                case Block.RawCopper _:
                    identifier = "minecraft:raw_copper_block"; properties = State61; return true;
                case Block.RawGold _:
                    identifier = "minecraft:raw_gold_block"; properties = State62; return true;
                case Block.RawIron _:
                    identifier = "minecraft:raw_iron_block"; properties = State63; return true;
                case Block.RedstoneBlock _:
                    identifier = "minecraft:redstone_block"; properties = State64; return true;
                case Block.ReinforcedDeepslate _:
                    identifier = "minecraft:reinforced_deepslate"; properties = State65; return true;
                case Block.Resin _:
                    identifier = "minecraft:resin_block"; properties = State66; return true;
                case Block.SeaLantern _:
                    identifier = "minecraft:sea_lantern"; properties = State67; return true;
                case Block.Shroomlight _:
                    identifier = "minecraft:shroomlight"; properties = State68; return true;
                case Block.Slime _:
                    identifier = "minecraft:slime"; properties = State69; return true;
                case Block.SmithingTable _:
                    identifier = "minecraft:smithing_table"; properties = State70; return true;
                case Block.SmoothBasalt _:
                    identifier = "minecraft:smooth_basalt"; properties = State71; return true;
                case Block.Snow _:
                    identifier = "minecraft:snow"; properties = State72; return true;
                case Block.SoulSand _:
                    identifier = "minecraft:soul_sand"; properties = State73; return true;
                case Block.SoulSoil _:
                    identifier = "minecraft:soul_soil"; properties = State74; return true;
                case Block.SporeBlossom _:
                    identifier = "minecraft:spore_blossom"; properties = State75; return true;
                case Block.TNT _:
                    identifier = "minecraft:tnt"; properties = State76; return true;
                case Block.Terracotta _:
                    identifier = "minecraft:hardened_clay"; properties = State77; return true;
                case Block.WheatSeeds _:
                    identifier = "minecraft:wheat"; properties = State78; return true;
                case Block.Sand { Red: true }:
                    identifier = "minecraft:red_sand"; properties = State80; return true;
                case Block.Sand:
                    identifier = "minecraft:sand"; properties = State79; return true;
                case Block.Water { Still: false, Depth: 8, Falling: false }:
                    identifier = "minecraft:flowing_water"; properties = State81; return true;
                case Block.Water { Still: false, Depth: 7, Falling: false }:
                    identifier = "minecraft:flowing_water"; properties = State82; return true;
                case Block.Water { Still: false, Depth: 6, Falling: false }:
                    identifier = "minecraft:flowing_water"; properties = State83; return true;
                case Block.Water { Still: false, Depth: 5, Falling: false }:
                    identifier = "minecraft:flowing_water"; properties = State84; return true;
                case Block.Water { Still: false, Depth: 4, Falling: false }:
                    identifier = "minecraft:flowing_water"; properties = State85; return true;
                case Block.Water { Still: false, Depth: 3, Falling: false }:
                    identifier = "minecraft:flowing_water"; properties = State86; return true;
                case Block.Water { Still: false, Depth: 2, Falling: false }:
                    identifier = "minecraft:flowing_water"; properties = State87; return true;
                case Block.Water { Still: false, Depth: 1, Falling: false }:
                    identifier = "minecraft:flowing_water"; properties = State88; return true;
                case Block.Water { Still: false, Depth: 8, Falling: true }:
                    identifier = "minecraft:flowing_water"; properties = State89; return true;
                case Block.Water { Still: false, Depth: 7, Falling: true }:
                    identifier = "minecraft:flowing_water"; properties = State90; return true;
                case Block.Water { Still: false, Depth: 6, Falling: true }:
                    identifier = "minecraft:flowing_water"; properties = State91; return true;
                case Block.Water { Still: false, Depth: 5, Falling: true }:
                    identifier = "minecraft:flowing_water"; properties = State92; return true;
                case Block.Water { Still: false, Depth: 4, Falling: true }:
                    identifier = "minecraft:flowing_water"; properties = State93; return true;
                case Block.Water { Still: false, Depth: 3, Falling: true }:
                    identifier = "minecraft:flowing_water"; properties = State94; return true;
                case Block.Water { Still: false, Depth: 2, Falling: true }:
                    identifier = "minecraft:flowing_water"; properties = State95; return true;
                case Block.Water { Still: false, Depth: 1, Falling: true }:
                    identifier = "minecraft:flowing_water"; properties = State96; return true;
                case Block.Water { Still: true, Depth: 8, Falling: false }:
                    identifier = "minecraft:water"; properties = State97; return true;
                case Block.Water { Still: true, Depth: 7, Falling: false }:
                    identifier = "minecraft:water"; properties = State98; return true;
                case Block.Water { Still: true, Depth: 6, Falling: false }:
                    identifier = "minecraft:water"; properties = State99; return true;
                case Block.Water { Still: true, Depth: 5, Falling: false }:
                    identifier = "minecraft:water"; properties = State100; return true;
                case Block.Water { Still: true, Depth: 4, Falling: false }:
                    identifier = "minecraft:water"; properties = State101; return true;
                case Block.Water { Still: true, Depth: 3, Falling: false }:
                    identifier = "minecraft:water"; properties = State102; return true;
                case Block.Water { Still: true, Depth: 2, Falling: false }:
                    identifier = "minecraft:water"; properties = State103; return true;
                case Block.Water { Still: true, Depth: 1, Falling: false }:
                    identifier = "minecraft:water"; properties = State104; return true;
                case Block.Water { Still: true, Depth: 8, Falling: true }:
                    identifier = "minecraft:water"; properties = State105; return true;
                case Block.Water { Still: true, Depth: 7, Falling: true }:
                    identifier = "minecraft:water"; properties = State106; return true;
                case Block.Water { Still: true, Depth: 6, Falling: true }:
                    identifier = "minecraft:water"; properties = State107; return true;
                case Block.Water { Still: true, Depth: 5, Falling: true }:
                    identifier = "minecraft:water"; properties = State108; return true;
                case Block.Water { Still: true, Depth: 4, Falling: true }:
                    identifier = "minecraft:water"; properties = State109; return true;
                case Block.Water { Still: true, Depth: 3, Falling: true }:
                    identifier = "minecraft:water"; properties = State110; return true;
                case Block.Water { Still: true, Depth: 2, Falling: true }:
                    identifier = "minecraft:water"; properties = State111; return true;
                case Block.Water { Still: true, Depth: 1, Falling: true }:
                    identifier = "minecraft:water"; properties = State112; return true;
                case Block.Lava { Still: true, Depth: 8, Falling: false }:
                    identifier = "minecraft:lava"; properties = State113; return true;
                case Block.Lava { Still: true, Depth: 7, Falling: false }:
                    identifier = "minecraft:lava"; properties = State114; return true;
                case Block.Lava { Still: true, Depth: 6, Falling: false }:
                    identifier = "minecraft:lava"; properties = State115; return true;
                case Block.Lava { Still: true, Depth: 5, Falling: false }:
                    identifier = "minecraft:lava"; properties = State116; return true;
                case Block.Lava { Still: true, Depth: 4, Falling: false }:
                    identifier = "minecraft:lava"; properties = State117; return true;
                case Block.Lava { Still: true, Depth: 3, Falling: false }:
                    identifier = "minecraft:lava"; properties = State118; return true;
                case Block.Lava { Still: true, Depth: 2, Falling: false }:
                    identifier = "minecraft:lava"; properties = State119; return true;
                case Block.Lava { Still: true, Depth: 1, Falling: false }:
                    identifier = "minecraft:lava"; properties = State120; return true;
                case Block.Lava { Still: true, Depth: 8, Falling: true }:
                    identifier = "minecraft:lava"; properties = State121; return true;
                case Block.Lava { Still: true, Depth: 7, Falling: true }:
                    identifier = "minecraft:lava"; properties = State122; return true;
                case Block.Lava { Still: true, Depth: 6, Falling: true }:
                    identifier = "minecraft:lava"; properties = State123; return true;
                case Block.Lava { Still: true, Depth: 5, Falling: true }:
                    identifier = "minecraft:lava"; properties = State124; return true;
                case Block.Lava { Still: true, Depth: 4, Falling: true }:
                    identifier = "minecraft:lava"; properties = State125; return true;
                case Block.Lava { Still: true, Depth: 3, Falling: true }:
                    identifier = "minecraft:lava"; properties = State126; return true;
                case Block.Lava { Still: true, Depth: 2, Falling: true }:
                    identifier = "minecraft:lava"; properties = State127; return true;
                case Block.Lava { Still: true, Depth: 1, Falling: true }:
                    identifier = "minecraft:lava"; properties = State128; return true;
                case Block.Lava { Still: false, Depth: 8, Falling: false }:
                    identifier = "minecraft:flowing_lava"; properties = State129; return true;
                case Block.Lava { Still: false, Depth: 7, Falling: false }:
                    identifier = "minecraft:flowing_lava"; properties = State130; return true;
                case Block.Lava { Still: false, Depth: 6, Falling: false }:
                    identifier = "minecraft:flowing_lava"; properties = State131; return true;
                case Block.Lava { Still: false, Depth: 5, Falling: false }:
                    identifier = "minecraft:flowing_lava"; properties = State132; return true;
                case Block.Lava { Still: false, Depth: 4, Falling: false }:
                    identifier = "minecraft:flowing_lava"; properties = State133; return true;
                case Block.Lava { Still: false, Depth: 3, Falling: false }:
                    identifier = "minecraft:flowing_lava"; properties = State134; return true;
                case Block.Lava { Still: false, Depth: 2, Falling: false }:
                    identifier = "minecraft:flowing_lava"; properties = State135; return true;
                case Block.Lava { Still: false, Depth: 1, Falling: false }:
                    identifier = "minecraft:flowing_lava"; properties = State136; return true;
                case Block.Lava { Still: false, Depth: 8, Falling: true }:
                    identifier = "minecraft:flowing_lava"; properties = State137; return true;
                case Block.Lava { Still: false, Depth: 7, Falling: true }:
                    identifier = "minecraft:flowing_lava"; properties = State138; return true;
                case Block.Lava { Still: false, Depth: 6, Falling: true }:
                    identifier = "minecraft:flowing_lava"; properties = State139; return true;
                case Block.Lava { Still: false, Depth: 5, Falling: true }:
                    identifier = "minecraft:flowing_lava"; properties = State140; return true;
                case Block.Lava { Still: false, Depth: 4, Falling: true }:
                    identifier = "minecraft:flowing_lava"; properties = State141; return true;
                case Block.Lava { Still: false, Depth: 3, Falling: true }:
                    identifier = "minecraft:flowing_lava"; properties = State142; return true;
                case Block.Lava { Still: false, Depth: 2, Falling: true }:
                    identifier = "minecraft:flowing_lava"; properties = State143; return true;
                case Block.Lava { Still: false, Depth: 1, Falling: true }:
                    identifier = "minecraft:flowing_lava"; properties = State144; return true;
                case EncodedLiquid liquidEncoded:
                    identifier = liquidEncoded.Identifier; properties = liquidEncoded.Properties; return true;
                case EncodedBlock encoded:
                    identifier = encoded.Identifier; properties = encoded.Properties; return true;
                default:
                    identifier = string.Empty; properties = Array.Empty<byte>(); return false;
            }
        }

        internal static World.Block Decode(string identifier, ReadOnlySpan<byte> properties)
        {
            if (identifier == "minecraft:air" && properties.SequenceEqual(State0)) return new Block.Air();
            if (identifier == "minecraft:amethyst_block" && properties.SequenceEqual(State1)) return new Block.Amethyst();
            if (identifier == "minecraft:ancient_debris" && properties.SequenceEqual(State2)) return new Block.AncientDebris();
            if (identifier == "minecraft:bamboo_mosaic" && properties.SequenceEqual(State3)) return new Block.BambooMosaic();
            if (identifier == "minecraft:barrier" && properties.SequenceEqual(State4)) return new Block.Barrier();
            if (identifier == "minecraft:beetroot" && properties.SequenceEqual(State5)) return new Block.BeetrootSeeds();
            if (identifier == "minecraft:blue_ice" && properties.SequenceEqual(State6)) return new Block.BlueIce();
            if (identifier == "minecraft:bookshelf" && properties.SequenceEqual(State7)) return new Block.Bookshelf();
            if (identifier == "minecraft:brick_block" && properties.SequenceEqual(State8)) return new Block.Bricks();
            if (identifier == "minecraft:calcite" && properties.SequenceEqual(State9)) return new Block.Calcite();
            if (identifier == "minecraft:carrots" && properties.SequenceEqual(State10)) return new Block.Carrot();
            if (identifier == "minecraft:chiseled_quartz_block" && properties.SequenceEqual(State11)) return new Block.ChiseledQuartz();
            if (identifier == "minecraft:clay" && properties.SequenceEqual(State12)) return new Block.Clay();
            if (identifier == "minecraft:coal_block" && properties.SequenceEqual(State13)) return new Block.Coal();
            if (identifier == "minecraft:web" && properties.SequenceEqual(State14)) return new Block.Cobweb();
            if (identifier == "minecraft:crafting_table" && properties.SequenceEqual(State15)) return new Block.CraftingTable();
            if (identifier == "minecraft:deadbush" && properties.SequenceEqual(State16)) return new Block.DeadBush();
            if (identifier == "minecraft:diamond_block" && properties.SequenceEqual(State17)) return new Block.Diamond();
            if (identifier == "minecraft:grass_path" && properties.SequenceEqual(State18)) return new Block.DirtPath();
            if (identifier == "minecraft:dragon_egg" && properties.SequenceEqual(State19)) return new Block.DragonEgg();
            if (identifier == "minecraft:dried_kelp_block" && properties.SequenceEqual(State20)) return new Block.DriedKelp();
            if (identifier == "minecraft:dripstone_block" && properties.SequenceEqual(State21)) return new Block.Dripstone();
            if (identifier == "minecraft:emerald_block" && properties.SequenceEqual(State22)) return new Block.Emerald();
            if (identifier == "minecraft:enchanting_table" && properties.SequenceEqual(State23)) return new Block.EnchantingTable();
            if (identifier == "minecraft:end_bricks" && properties.SequenceEqual(State24)) return new Block.EndBricks();
            if (identifier == "minecraft:end_portal" && properties.SequenceEqual(State25)) return new Block.EndPortal();
            if (identifier == "minecraft:end_stone" && properties.SequenceEqual(State26)) return new Block.EndStone();
            if (identifier == "minecraft:fern" && properties.SequenceEqual(State27)) return new Block.Fern();
            if (identifier == "minecraft:fletching_table" && properties.SequenceEqual(State28)) return new Block.FletchingTable();
            if (identifier == "minecraft:glass" && properties.SequenceEqual(State29)) return new Block.Glass();
            if (identifier == "minecraft:glass_pane" && properties.SequenceEqual(State30)) return new Block.GlassPane();
            if (identifier == "minecraft:glowstone" && properties.SequenceEqual(State31)) return new Block.Glowstone();
            if (identifier == "minecraft:gold_block" && properties.SequenceEqual(State32)) return new Block.Gold();
            if (identifier == "minecraft:grass_block" && properties.SequenceEqual(State33)) return new Block.Grass();
            if (identifier == "minecraft:gravel" && properties.SequenceEqual(State34)) return new Block.Gravel();
            if (identifier == "minecraft:honeycomb_block" && properties.SequenceEqual(State35)) return new Block.Honeycomb();
            if (identifier == "minecraft:infested_cobblestone" && properties.SequenceEqual(State36)) return new Block.InfestedCobblestone();
            if (identifier == "minecraft:infested_stone" && properties.SequenceEqual(State37)) return new Block.InfestedStone();
            if (identifier == "minecraft:invisible_bedrock" && properties.SequenceEqual(State38)) return new Block.InvisibleBedrock();
            if (identifier == "minecraft:iron_block" && properties.SequenceEqual(State39)) return new Block.Iron();
            if (identifier == "minecraft:iron_bars" && properties.SequenceEqual(State40)) return new Block.IronBars();
            if (identifier == "minecraft:lapis_block" && properties.SequenceEqual(State41)) return new Block.Lapis();
            if (identifier == "minecraft:waterlily" && properties.SequenceEqual(State42)) return new Block.LilyPad();
            if (identifier == "minecraft:magma" && properties.SequenceEqual(State43)) return new Block.Magma();
            if (identifier == "minecraft:melon_block" && properties.SequenceEqual(State44)) return new Block.Melon();
            if (identifier == "minecraft:moss_carpet" && properties.SequenceEqual(State45)) return new Block.MossCarpet();
            if (identifier == "minecraft:mud" && properties.SequenceEqual(State46)) return new Block.Mud();
            if (identifier == "minecraft:mud_bricks" && properties.SequenceEqual(State47)) return new Block.MudBricks();
            if (identifier == "minecraft:nether_brick_fence" && properties.SequenceEqual(State48)) return new Block.NetherBrickFence();
            if (identifier == "minecraft:nether_gold_ore" && properties.SequenceEqual(State49)) return new Block.NetherGoldOre();
            if (identifier == "minecraft:quartz_ore" && properties.SequenceEqual(State50)) return new Block.NetherQuartzOre();
            if (identifier == "minecraft:nether_sprouts" && properties.SequenceEqual(State51)) return new Block.NetherSprouts();
            if (identifier == "minecraft:netherite_block" && properties.SequenceEqual(State52)) return new Block.Netherite();
            if (identifier == "minecraft:netherrack" && properties.SequenceEqual(State53)) return new Block.Netherrack();
            if (identifier == "minecraft:packed_ice" && properties.SequenceEqual(State54)) return new Block.PackedIce();
            if (identifier == "minecraft:packed_mud" && properties.SequenceEqual(State55)) return new Block.PackedMud();
            if (identifier == "minecraft:podzol" && properties.SequenceEqual(State56)) return new Block.Podzol();
            if (identifier == "minecraft:polished_tuff" && properties.SequenceEqual(State57)) return new Block.PolishedTuff();
            if (identifier == "minecraft:potatoes" && properties.SequenceEqual(State58)) return new Block.Potato();
            if (identifier == "minecraft:purpur_block" && properties.SequenceEqual(State59)) return new Block.Purpur();
            if (identifier == "minecraft:quartz_bricks" && properties.SequenceEqual(State60)) return new Block.QuartzBricks();
            if (identifier == "minecraft:raw_copper_block" && properties.SequenceEqual(State61)) return new Block.RawCopper();
            if (identifier == "minecraft:raw_gold_block" && properties.SequenceEqual(State62)) return new Block.RawGold();
            if (identifier == "minecraft:raw_iron_block" && properties.SequenceEqual(State63)) return new Block.RawIron();
            if (identifier == "minecraft:redstone_block" && properties.SequenceEqual(State64)) return new Block.RedstoneBlock();
            if (identifier == "minecraft:reinforced_deepslate" && properties.SequenceEqual(State65)) return new Block.ReinforcedDeepslate();
            if (identifier == "minecraft:resin_block" && properties.SequenceEqual(State66)) return new Block.Resin();
            if (identifier == "minecraft:sea_lantern" && properties.SequenceEqual(State67)) return new Block.SeaLantern();
            if (identifier == "minecraft:shroomlight" && properties.SequenceEqual(State68)) return new Block.Shroomlight();
            if (identifier == "minecraft:slime" && properties.SequenceEqual(State69)) return new Block.Slime();
            if (identifier == "minecraft:smithing_table" && properties.SequenceEqual(State70)) return new Block.SmithingTable();
            if (identifier == "minecraft:smooth_basalt" && properties.SequenceEqual(State71)) return new Block.SmoothBasalt();
            if (identifier == "minecraft:snow" && properties.SequenceEqual(State72)) return new Block.Snow();
            if (identifier == "minecraft:soul_sand" && properties.SequenceEqual(State73)) return new Block.SoulSand();
            if (identifier == "minecraft:soul_soil" && properties.SequenceEqual(State74)) return new Block.SoulSoil();
            if (identifier == "minecraft:spore_blossom" && properties.SequenceEqual(State75)) return new Block.SporeBlossom();
            if (identifier == "minecraft:tnt" && properties.SequenceEqual(State76)) return new Block.TNT();
            if (identifier == "minecraft:hardened_clay" && properties.SequenceEqual(State77)) return new Block.Terracotta();
            if (identifier == "minecraft:wheat" && properties.SequenceEqual(State78)) return new Block.WheatSeeds();
            if (identifier == "minecraft:sand" && properties.SequenceEqual(State79)) return new Block.Sand();
            if (identifier == "minecraft:red_sand" && properties.SequenceEqual(State80)) return new Block.Sand(true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State81)) return new Block.Water(false, 8, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State82)) return new Block.Water(false, 7, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State83)) return new Block.Water(false, 6, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State84)) return new Block.Water(false, 5, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State85)) return new Block.Water(false, 4, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State86)) return new Block.Water(false, 3, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State87)) return new Block.Water(false, 2, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State88)) return new Block.Water(false, 1, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State89)) return new Block.Water(false, 8, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State90)) return new Block.Water(false, 7, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State91)) return new Block.Water(false, 6, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State92)) return new Block.Water(false, 5, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State93)) return new Block.Water(false, 4, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State94)) return new Block.Water(false, 3, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State95)) return new Block.Water(false, 2, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State96)) return new Block.Water(false, 1, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State97)) return new Block.Water(true, 8, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State98)) return new Block.Water(true, 7, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State99)) return new Block.Water(true, 6, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State100)) return new Block.Water(true, 5, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State101)) return new Block.Water(true, 4, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State102)) return new Block.Water(true, 3, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State103)) return new Block.Water(true, 2, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State104)) return new Block.Water(true, 1, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State105)) return new Block.Water(true, 8, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State106)) return new Block.Water(true, 7, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State107)) return new Block.Water(true, 6, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State108)) return new Block.Water(true, 5, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State109)) return new Block.Water(true, 4, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State110)) return new Block.Water(true, 3, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State111)) return new Block.Water(true, 2, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State112)) return new Block.Water(true, 1, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State113)) return new Block.Lava(true, 8, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State114)) return new Block.Lava(true, 7, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State115)) return new Block.Lava(true, 6, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State116)) return new Block.Lava(true, 5, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State117)) return new Block.Lava(true, 4, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State118)) return new Block.Lava(true, 3, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State119)) return new Block.Lava(true, 2, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State120)) return new Block.Lava(true, 1, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State121)) return new Block.Lava(true, 8, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State122)) return new Block.Lava(true, 7, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State123)) return new Block.Lava(true, 6, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State124)) return new Block.Lava(true, 5, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State125)) return new Block.Lava(true, 4, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State126)) return new Block.Lava(true, 3, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State127)) return new Block.Lava(true, 2, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State128)) return new Block.Lava(true, 1, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State129)) return new Block.Lava(false, 8, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State130)) return new Block.Lava(false, 7, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State131)) return new Block.Lava(false, 6, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State132)) return new Block.Lava(false, 5, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State133)) return new Block.Lava(false, 4, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State134)) return new Block.Lava(false, 3, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State135)) return new Block.Lava(false, 2, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State136)) return new Block.Lava(false, 1, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State137)) return new Block.Lava(false, 8, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State138)) return new Block.Lava(false, 7, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State139)) return new Block.Lava(false, 6, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State140)) return new Block.Lava(false, 5, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State141)) return new Block.Lava(false, 4, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State142)) return new Block.Lava(false, 3, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State143)) return new Block.Lava(false, 2, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State144)) return new Block.Lava(false, 1, true);
            return new EncodedBlock(identifier, properties.ToArray());
        }

        internal static World.Liquid DecodeLiquid(string identifier, ReadOnlySpan<byte> properties)
        {
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State81)) return new Block.Water(false, 8, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State82)) return new Block.Water(false, 7, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State83)) return new Block.Water(false, 6, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State84)) return new Block.Water(false, 5, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State85)) return new Block.Water(false, 4, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State86)) return new Block.Water(false, 3, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State87)) return new Block.Water(false, 2, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State88)) return new Block.Water(false, 1, false);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State89)) return new Block.Water(false, 8, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State90)) return new Block.Water(false, 7, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State91)) return new Block.Water(false, 6, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State92)) return new Block.Water(false, 5, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State93)) return new Block.Water(false, 4, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State94)) return new Block.Water(false, 3, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State95)) return new Block.Water(false, 2, true);
            if (identifier == "minecraft:flowing_water" && properties.SequenceEqual(State96)) return new Block.Water(false, 1, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State97)) return new Block.Water(true, 8, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State98)) return new Block.Water(true, 7, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State99)) return new Block.Water(true, 6, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State100)) return new Block.Water(true, 5, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State101)) return new Block.Water(true, 4, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State102)) return new Block.Water(true, 3, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State103)) return new Block.Water(true, 2, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State104)) return new Block.Water(true, 1, false);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State105)) return new Block.Water(true, 8, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State106)) return new Block.Water(true, 7, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State107)) return new Block.Water(true, 6, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State108)) return new Block.Water(true, 5, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State109)) return new Block.Water(true, 4, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State110)) return new Block.Water(true, 3, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State111)) return new Block.Water(true, 2, true);
            if (identifier == "minecraft:water" && properties.SequenceEqual(State112)) return new Block.Water(true, 1, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State113)) return new Block.Lava(true, 8, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State114)) return new Block.Lava(true, 7, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State115)) return new Block.Lava(true, 6, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State116)) return new Block.Lava(true, 5, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State117)) return new Block.Lava(true, 4, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State118)) return new Block.Lava(true, 3, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State119)) return new Block.Lava(true, 2, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State120)) return new Block.Lava(true, 1, false);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State121)) return new Block.Lava(true, 8, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State122)) return new Block.Lava(true, 7, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State123)) return new Block.Lava(true, 6, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State124)) return new Block.Lava(true, 5, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State125)) return new Block.Lava(true, 4, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State126)) return new Block.Lava(true, 3, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State127)) return new Block.Lava(true, 2, true);
            if (identifier == "minecraft:lava" && properties.SequenceEqual(State128)) return new Block.Lava(true, 1, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State129)) return new Block.Lava(false, 8, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State130)) return new Block.Lava(false, 7, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State131)) return new Block.Lava(false, 6, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State132)) return new Block.Lava(false, 5, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State133)) return new Block.Lava(false, 4, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State134)) return new Block.Lava(false, 3, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State135)) return new Block.Lava(false, 2, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State136)) return new Block.Lava(false, 1, false);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State137)) return new Block.Lava(false, 8, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State138)) return new Block.Lava(false, 7, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State139)) return new Block.Lava(false, 6, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State140)) return new Block.Lava(false, 5, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State141)) return new Block.Lava(false, 4, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State142)) return new Block.Lava(false, 3, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State143)) return new Block.Lava(false, 2, true);
            if (identifier == "minecraft:flowing_lava" && properties.SequenceEqual(State144)) return new Block.Lava(false, 1, true);
            return new EncodedLiquid(identifier, properties.ToArray());
        }

        private sealed record EncodedBlock(string Identifier, byte[] Properties) : World.Block;
        private sealed record EncodedLiquid(string Identifier, byte[] Properties) : World.Liquid
        {
            public string LiquidType() => throw new InvalidOperationException("Opaque liquid type was not transported by the host.");
        }
    }
}
