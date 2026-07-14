// Code generated from Dragonfly particle, sound/instrument, and image/color Go AST. DO NOT EDIT.
#nullable enable

namespace Dragonfly
{
    public static partial class Color
    {
        public readonly record struct RGBA(byte R, byte G, byte B, byte A);
    }

    public static partial class Sound
    {
        public readonly struct Instrument
        {
            private readonly uint _id;
            internal Instrument(uint id) => _id = id;
            internal uint Id => _id;
        }

        public static Instrument Piano() => new(0u);
        public static Instrument BassDrum() => new(1u);
        public static Instrument Snare() => new(2u);
        public static Instrument ClicksAndSticks() => new(3u);
        public static Instrument Bass() => new(4u);
        public static Instrument Flute() => new(5u);
        public static Instrument Bell() => new(6u);
        public static Instrument Guitar() => new(7u);
        public static Instrument Chimes() => new(8u);
        public static Instrument Xylophone() => new(9u);
        public static Instrument IronXylophone() => new(10u);
        public static Instrument CowBell() => new(11u);
        public static Instrument Didgeridoo() => new(12u);
        public static Instrument Bit() => new(13u);
        public static Instrument Banjo() => new(14u);
        public static Instrument Pling() => new(15u);
    }

    public static partial class Particle
    {
        public readonly record struct Flame(Color.RGBA Colour) : World.Particle;
        public readonly record struct Dust(Color.RGBA Colour) : World.Particle;
        public readonly record struct BlockBreak(World.Block Block) : World.Particle;
        public readonly record struct PunchBlock(World.Block Block, Cube.Face Face) : World.Particle;
        public readonly record struct BlockForceField : World.Particle;
        public readonly record struct BoneMeal(bool Area) : World.Particle;
        public readonly record struct Note(Sound.Instrument Instrument, int Pitch) : World.Particle;
        public readonly record struct DragonEggTeleport(Cube.Pos Diff) : World.Particle;
        public readonly record struct Evaporate : World.Particle;
        public readonly record struct WaterDrip : World.Particle;
        public readonly record struct LavaDrip : World.Particle;
        public readonly record struct Lava : World.Particle;
        public readonly record struct DustPlume : World.Particle;
        public readonly record struct HugeExplosion : World.Particle;
        public readonly record struct EndermanTeleport : World.Particle;
        public readonly record struct SnowballPoof : World.Particle;
        public readonly record struct EggSmash : World.Particle;
        public readonly record struct Splash(Color.RGBA Colour) : World.Particle;
        public readonly record struct Effect(Color.RGBA Colour) : World.Particle;
        public readonly record struct EntityFlame : World.Particle;
    }

    internal readonly record struct EncodedParticle(
        uint Kind, uint Data, int Pitch, Color.RGBA Colour, Cube.Pos Diff, World.Block? Block);

    internal static class ParticleCodec
    {
        internal static bool TryEncode(World.Particle particle, out EncodedParticle encoded)
        {
            switch (particle)
            {
                case Particle.Flame value:
                    encoded = new(0u, 0u, 0, value.Colour, default, null); return true;
                case Particle.Dust value:
                    encoded = new(1u, 0u, 0, value.Colour, default, null); return true;
                case Particle.BlockBreak value:
                    encoded = new(2u, 0u, 0, default, default, value.Block); return true;
                case Particle.PunchBlock value:
                    encoded = new(3u, (uint)value.Face, 0, default, default, value.Block); return true;
                case Particle.BlockForceField _:
                    encoded = new(4u, 0u, 0, default, default, null); return true;
                case Particle.BoneMeal value:
                    encoded = new(5u, value.Area ? 1u : 0u, 0, default, default, null); return true;
                case Particle.Note value:
                    encoded = new(6u, value.Instrument.Id, value.Pitch, default, default, null); return true;
                case Particle.DragonEggTeleport value:
                    encoded = new(7u, 0u, 0, default, value.Diff, null); return true;
                case Particle.Evaporate _:
                    encoded = new(8u, 0u, 0, default, default, null); return true;
                case Particle.WaterDrip _:
                    encoded = new(9u, 0u, 0, default, default, null); return true;
                case Particle.LavaDrip _:
                    encoded = new(10u, 0u, 0, default, default, null); return true;
                case Particle.Lava _:
                    encoded = new(11u, 0u, 0, default, default, null); return true;
                case Particle.DustPlume _:
                    encoded = new(12u, 0u, 0, default, default, null); return true;
                case Particle.HugeExplosion _:
                    encoded = new(13u, 0u, 0, default, default, null); return true;
                case Particle.EndermanTeleport _:
                    encoded = new(14u, 0u, 0, default, default, null); return true;
                case Particle.SnowballPoof _:
                    encoded = new(15u, 0u, 0, default, default, null); return true;
                case Particle.EggSmash _:
                    encoded = new(16u, 0u, 0, default, default, null); return true;
                case Particle.Splash value:
                    encoded = new(17u, 0u, 0, value.Colour, default, null); return true;
                case Particle.Effect value:
                    encoded = new(18u, 0u, 0, value.Colour, default, null); return true;
                case Particle.EntityFlame _:
                    encoded = new(19u, 0u, 0, default, default, null); return true;
                default:
                    encoded = default; return false;
            }
        }
    }
}
