// Code generated from Dragonfly server/world/difficulty.go AST and live registry. DO NOT EDIT.
#nullable enable
using System;
using Dragonfly.Native;

namespace Dragonfly;

public sealed partial class World
{
    public interface Difficulty
    {
        bool FoodRegenerates();
        double StarvationHealthLimit();
        int FireSpreadIncrease();
    }

    public static readonly Difficulty DifficultyPeaceful = new BuiltinDifficulty(0, true, 20d, 0);
    public static readonly Difficulty DifficultyEasy = new BuiltinDifficulty(1, false, 10d, 7);
    public static readonly Difficulty DifficultyNormal = new BuiltinDifficulty(2, false, 2d, 14);
    public static readonly Difficulty DifficultyHard = new BuiltinDifficulty(3, false, -1d, 21);

    public static (Difficulty Difficulty, bool Ok) DifficultyByID(int id) => id switch
    {
        0 => (DifficultyPeaceful, true),
        1 => (DifficultyEasy, true),
        2 => (DifficultyNormal, true),
        3 => (DifficultyHard, true),
        _ => (DifficultyNormal, false),
    };

    public static (int ID, bool Ok) DifficultyID(Difficulty diff)
    {
        if (diff is BuiltinDifficulty builtin) return (builtin.ID, true);
        return (0, false);
    }

    internal static DifficultyView DifficultyView(Difficulty difficulty)
    {
        ArgumentNullException.ThrowIfNull(difficulty);
        return new DifficultyView
        {
            ID = difficulty is BuiltinDifficulty builtin ? checked((uint)builtin.ID) : 0,
            Builtin = difficulty is BuiltinDifficulty ? (byte)1 : (byte)0,
            FoodRegenerates = difficulty.FoodRegenerates() ? (byte)1 : (byte)0,
            StarvationHealthLimit = difficulty.StarvationHealthLimit(),
            FireSpreadIncrease = difficulty.FireSpreadIncrease(),
        };
    }

    internal static Difficulty DifficultyFromView(DifficultyView view)
    {
        if (view.Builtin > 1 || view.FoodRegenerates > 1)
            throw new InvalidOperationException("invalid difficulty view");
        if (view.Builtin == 1)
        {
            if (view.ID > int.MaxValue)
                throw new InvalidOperationException("invalid difficulty view");
            var (difficulty, ok) = DifficultyByID((int)view.ID);
            if (!ok) throw new InvalidOperationException("invalid difficulty view");
            return difficulty;
        }
        if (view.ID != 0)
            throw new InvalidOperationException("invalid difficulty view");
        return new CapabilityDifficulty(
            view.FoodRegenerates != 0,
            view.StarvationHealthLimit,
            view.FireSpreadIncrease);
    }

    private class CapabilityDifficulty(
        bool foodRegenerates,
        double starvationHealthLimit,
        int fireSpreadIncrease) : Difficulty
    {
        public bool FoodRegenerates() => foodRegenerates;
        public double StarvationHealthLimit() => starvationHealthLimit;
        public int FireSpreadIncrease() => fireSpreadIncrease;
    }

    private sealed class BuiltinDifficulty(
        int id,
        bool foodRegenerates,
        double starvationHealthLimit,
        int fireSpreadIncrease) : CapabilityDifficulty(foodRegenerates, starvationHealthLimit, fireSpreadIncrease)
    {
        internal int ID { get; } = id;
    }
}
