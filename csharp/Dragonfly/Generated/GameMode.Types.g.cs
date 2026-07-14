// Code generated from Dragonfly server/world/game_mode.go AST and live registry. DO NOT EDIT.
#nullable enable
using System;

namespace Dragonfly;

public sealed partial class World
{
    public interface GameMode
    {
        bool AllowsEditing();
        bool AllowsTakingDamage();
        bool CreativeInventory();
        bool HasCollision();
        bool AllowsFlying();
        bool AllowsInteraction();
        bool Visible();
        bool InstantPortalTravel();
    }

    public static readonly GameMode GameModeSurvival = new BuiltinGameMode(0, 0x6bUL);
    public static readonly GameMode GameModeCreative = new BuiltinGameMode(1, 0xfdUL);
    public static readonly GameMode GameModeAdventure = new BuiltinGameMode(2, 0x6aUL);
    public static readonly GameMode GameModeSpectator = new BuiltinGameMode(3, 0x10UL);

    public static (GameMode GameMode, bool Ok) GameModeByID(int id) => id switch
    {
        0 => (GameModeSurvival, true),
        1 => (GameModeCreative, true),
        2 => (GameModeAdventure, true),
        3 => (GameModeSpectator, true),
        _ => (GameModeSurvival, false),
    };

    public static (int ID, bool Ok) GameModeID(GameMode mode)
    {
        if (mode is BuiltinGameMode builtin) return (builtin.ID, true);
        return (0, false);
    }

    internal static long GameModeDescriptor(GameMode mode)
    {
        ArgumentNullException.ThrowIfNull(mode);
        if (mode is BuiltinGameMode builtin)
            return unchecked((long)(BuiltinGameModeFlag | (uint)builtin.ID));
        ulong capabilities = 0;
        if (mode.AllowsEditing()) capabilities |= 1UL << 0;
        if (mode.AllowsTakingDamage()) capabilities |= 1UL << 1;
        if (mode.CreativeInventory()) capabilities |= 1UL << 2;
        if (mode.HasCollision()) capabilities |= 1UL << 3;
        if (mode.AllowsFlying()) capabilities |= 1UL << 4;
        if (mode.AllowsInteraction()) capabilities |= 1UL << 5;
        if (mode.Visible()) capabilities |= 1UL << 6;
        if (mode.InstantPortalTravel()) capabilities |= 1UL << 7;
        return (long)capabilities;
    }

    internal static GameMode GameModeFromDescriptor(long descriptor)
    {
        var value = unchecked((ulong)descriptor);
        if ((value & BuiltinGameModeFlag) != 0)
        {
            var rawID = value & ~BuiltinGameModeFlag;
            if (rawID > int.MaxValue)
                throw new InvalidOperationException("invalid game mode descriptor");
            var (mode, ok) = GameModeByID((int)rawID);
            if (!ok) throw new InvalidOperationException("invalid game mode descriptor");
            return mode;
        }
        if ((value & ~CustomGameModeMask) != 0)
            throw new InvalidOperationException("invalid game mode descriptor");
        return new CapabilityGameMode(value);
    }

    private const ulong BuiltinGameModeFlag = 1UL << 63;
    private const ulong CustomGameModeMask = (1UL << 8) - 1;

    private class CapabilityGameMode(ulong capabilities) : GameMode
    {
        public bool AllowsEditing() => (capabilities & (1UL << 0)) != 0;
        public bool AllowsTakingDamage() => (capabilities & (1UL << 1)) != 0;
        public bool CreativeInventory() => (capabilities & (1UL << 2)) != 0;
        public bool HasCollision() => (capabilities & (1UL << 3)) != 0;
        public bool AllowsFlying() => (capabilities & (1UL << 4)) != 0;
        public bool AllowsInteraction() => (capabilities & (1UL << 5)) != 0;
        public bool Visible() => (capabilities & (1UL << 6)) != 0;
        public bool InstantPortalTravel() => (capabilities & (1UL << 7)) != 0;
    }

    private sealed class BuiltinGameMode(int id, ulong capabilities) : CapabilityGameMode(capabilities)
    {
        internal int ID { get; } = id;
    }
}
