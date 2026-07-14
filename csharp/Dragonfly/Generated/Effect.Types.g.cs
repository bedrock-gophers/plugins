// Code generated from Dragonfly server/entity/effect Go AST and live registry. DO NOT EDIT.
#nullable enable
using System;
using System.Collections.Generic;

namespace Dragonfly;

public static partial class Effect
{
    public interface Type
    {
        Color.RGBA RGBA();
    }

    public interface LastingType : Type { }

    public readonly record struct Value
    {
        private Type? _type { get; init; }
        private TimeSpan _duration { get; init; }
        private int _level { get; init; }
        private double _potency { get; init; }
        private bool _ambient { get; init; }
        private bool _particlesHidden { get; init; }
        private bool _infinite { get; init; }
        private int _tick { get; init; }

        internal Value(Type type, TimeSpan duration, int level, double potency, bool ambient, bool particlesHidden, bool infinite, int tick)
        {
            _type = type;
            _duration = duration;
            _level = level;
            _potency = potency;
            _ambient = ambient;
            _particlesHidden = particlesHidden;
            _infinite = infinite;
            _tick = tick;
        }

        public Value WithoutParticles() => this with { _particlesHidden = true };
        public bool ParticlesHidden() => _particlesHidden;
        public int Level() => _level;
        public TimeSpan Duration() => _duration;
        public bool Ambient() => _ambient;
        public bool Infinite() => _infinite;
        public Type? Type() => _type;

        public Value TickDuration()
        {
            if (_type is not LastingType) return this;
            return this with
            {
                _duration = _infinite ? _duration : _duration - TimeSpan.FromMilliseconds(50),
                _tick = _tick + 1,
            };
        }

        public int Tick() => _tick;

        internal double Potency => _potency;
    }

    public static Value NewInstant(Type type, int level) => NewInstantWithPotency(type, level, 1d);

    public static Value NewInstantWithPotency(Type type, int level, double potency)
    {
        ArgumentNullException.ThrowIfNull(type);
        return new Value(type, default, level, potency, false, false, false, 0);
    }

    public static Value New(LastingType type, int level, TimeSpan duration)
    {
        ArgumentNullException.ThrowIfNull(type);
        return new Value(type, duration, level, 0d, false, false, false, 0);
    }

    public static Value NewAmbient(LastingType type, int level, TimeSpan duration)
    {
        ArgumentNullException.ThrowIfNull(type);
        return new Value(type, duration, level, 0d, true, false, false, 0);
    }

    public static Value NewInfinite(LastingType type, int level)
    {
        ArgumentNullException.ThrowIfNull(type);
        return new Value(type, default, level, 0d, false, false, true, 0);
    }

    public static (Color.RGBA Colour, bool Ambient) ResultingColour(IReadOnlyList<Value> effects)
    {
        ArgumentNullException.ThrowIfNull(effects);
        long red = 0;
        long green = 0;
        long blue = 0;
        long alpha = 0;
        var count = 0;
        var ambient = true;
        foreach (var effect in effects)
        {
            if (effect.ParticlesHidden()) continue;
            var type = effect.Type() ?? throw new InvalidOperationException("Effect has no type.");
            var colour = type.RGBA();
            red += colour.R;
            green += colour.G;
            blue += colour.B;
            alpha += colour.A;
            count++;
            if (!effect.Ambient()) ambient = false;
        }
        if (count == 0) return (new Color.RGBA(0x38, 0x5d, 0xc6, 0xff), false);
        return (new Color.RGBA((byte)(red / count), (byte)(green / count), (byte)(blue / count), (byte)(alpha / count)), ambient);
    }

    public static readonly LastingType Speed = new BuiltinLastingType(1, new Color.RGBA(51, 235, 255, 255));
    public static readonly LastingType Slowness = new BuiltinLastingType(2, new Color.RGBA(139, 175, 224, 255));
    public static readonly LastingType Haste = new BuiltinLastingType(3, new Color.RGBA(217, 192, 67, 255));
    public static readonly LastingType MiningFatigue = new BuiltinLastingType(4, new Color.RGBA(74, 66, 23, 255));
    public static readonly LastingType Strength = new BuiltinLastingType(5, new Color.RGBA(255, 199, 0, 255));
    public static readonly Type InstantHealth = new BuiltinInstantType(6, new Color.RGBA(248, 36, 35, 255));
    public static readonly Type InstantDamage = new BuiltinInstantType(7, new Color.RGBA(169, 101, 106, 255));
    public static readonly LastingType JumpBoost = new BuiltinLastingType(8, new Color.RGBA(253, 255, 132, 255));
    public static readonly LastingType Nausea = new BuiltinLastingType(9, new Color.RGBA(85, 29, 74, 255));
    public static readonly LastingType Regeneration = new BuiltinLastingType(10, new Color.RGBA(205, 92, 171, 255));
    public static readonly LastingType Resistance = new BuiltinLastingType(11, new Color.RGBA(145, 70, 240, 255));
    public static readonly LastingType FireResistance = new BuiltinLastingType(12, new Color.RGBA(255, 153, 0, 255));
    public static readonly LastingType WaterBreathing = new BuiltinLastingType(13, new Color.RGBA(152, 218, 192, 255));
    public static readonly LastingType Invisibility = new BuiltinLastingType(14, new Color.RGBA(246, 246, 246, 255));
    public static readonly LastingType Blindness = new BuiltinLastingType(15, new Color.RGBA(31, 31, 35, 255));
    public static readonly LastingType NightVision = new BuiltinLastingType(16, new Color.RGBA(194, 255, 102, 255));
    public static readonly LastingType Hunger = new BuiltinLastingType(17, new Color.RGBA(88, 118, 83, 255));
    public static readonly LastingType Weakness = new BuiltinLastingType(18, new Color.RGBA(72, 77, 72, 255));
    public static readonly LastingType Poison = new BuiltinLastingType(19, new Color.RGBA(135, 163, 99, 255));
    public static readonly LastingType Wither = new BuiltinLastingType(20, new Color.RGBA(115, 97, 86, 255));
    public static readonly LastingType HealthBoost = new BuiltinLastingType(21, new Color.RGBA(248, 125, 35, 255));
    public static readonly LastingType Absorption = new BuiltinLastingType(22, new Color.RGBA(37, 82, 165, 255));
    public static readonly LastingType Saturation = new BuiltinLastingType(23, new Color.RGBA(248, 36, 35, 255));
    public static readonly LastingType Levitation = new BuiltinLastingType(24, new Color.RGBA(206, 255, 255, 255));
    public static readonly LastingType FatalPoison = new BuiltinLastingType(25, new Color.RGBA(78, 147, 49, 255));
    public static readonly LastingType ConduitPower = new BuiltinLastingType(26, new Color.RGBA(29, 194, 209, 255));
    public static readonly LastingType SlowFalling = new BuiltinLastingType(27, new Color.RGBA(243, 207, 185, 255));
    public static readonly LastingType Darkness = new BuiltinLastingType(30, new Color.RGBA(41, 39, 33, 255));

    internal static bool TryID(Type? type, out int id)
    {
        switch (type)
        {
            case BuiltinInstantType instant:
                id = instant.ID;
                return true;
            case BuiltinLastingType lasting:
                id = lasting.ID;
                return true;
            default:
                id = 0;
                return false;
        }
    }

    public static (Type? Type, bool Ok) ByID(int id)
    {
        var type = TypeByID(id);
        return (type, type is not null);
    }

    public static (int ID, bool Ok) ID(Type type)
    {
        ArgumentNullException.ThrowIfNull(type);
        return TryID(type, out var id) ? (id, true) : (0, false);
    }

    internal static Type? TypeByID(int id) => id switch
    {
        1 => Speed,
        2 => Slowness,
        3 => Haste,
        4 => MiningFatigue,
        5 => Strength,
        6 => InstantHealth,
        7 => InstantDamage,
        8 => JumpBoost,
        9 => Nausea,
        10 => Regeneration,
        11 => Resistance,
        12 => FireResistance,
        13 => WaterBreathing,
        14 => Invisibility,
        15 => Blindness,
        16 => NightVision,
        17 => Hunger,
        18 => Weakness,
        19 => Poison,
        20 => Wither,
        21 => HealthBoost,
        22 => Absorption,
        23 => Saturation,
        24 => Levitation,
        25 => FatalPoison,
        26 => ConduitPower,
        27 => SlowFalling,
        30 => Darkness,
        _ => null,
    };

    private sealed record BuiltinInstantType(int ID, Color.RGBA Colour) : Type
    {
        public Color.RGBA RGBA() => Colour;
    }

    private sealed record BuiltinLastingType(int ID, Color.RGBA Colour) : LastingType
    {
        public Color.RGBA RGBA() => Colour;
    }
}
