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

    public static readonly GameMode GameModeSurvival = new BuiltinGameMode(0, true, true, false, true, false, true, true, false);
    public static readonly GameMode GameModeCreative = new BuiltinGameMode(1, true, false, true, true, true, true, true, true);
    public static readonly GameMode GameModeAdventure = new BuiltinGameMode(2, false, true, false, true, false, true, true, false);
    public static readonly GameMode GameModeSpectator = new BuiltinGameMode(3, false, false, false, false, true, false, false, true);

    internal static bool GameModeId(GameMode mode, out long id)
    {
        if (mode is BuiltinGameMode builtin)
        {
            id = builtin.Id;
            return true;
        }
        id = 0;
        return false;
    }

    public partial class Tx
    {
        internal Tx(ulong invocation) => Invocation = invocation;
        internal ulong Invocation { get; }
    }

    public class Context : Tx
    {
        private bool _cancelled;

        internal Context(ulong invocation, bool cancelled) : base(invocation) =>
            _cancelled = cancelled;

        public bool Cancelled() => _cancelled;
        public void Cancel() => _cancelled = true;
    }

    private sealed class BuiltinGameMode(
        long id,
        bool editing,
        bool damage,
        bool creativeInventory,
        bool collision,
        bool flying,
        bool interaction,
        bool visible,
        bool instantPortal) : GameMode
    {
        internal long Id { get; } = id;
        public bool AllowsEditing() => editing;
        public bool AllowsTakingDamage() => damage;
        public bool CreativeInventory() => creativeInventory;
        public bool HasCollision() => collision;
        public bool AllowsFlying() => flying;
        public bool AllowsInteraction() => interaction;
        public bool Visible() => visible;
        public bool InstantPortalTravel() => instantPortal;
    }
}
