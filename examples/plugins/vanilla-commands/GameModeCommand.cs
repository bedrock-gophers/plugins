using Dragonfly;

internal static class GameModeCommand
{
    public static void Register() => Cmd.Register(Cmd.New(
        "gamemode",
        "Changes a player's game mode.",
        ["gm"],
        new GameMode()));

    internal enum GameModeValue
    {
        Survival,
        Creative,
        Adventure,
        Spectator,
    }

    internal sealed class GameMode : Cmd.Runnable
    {
        public GameModeValue Mode;
        public Cmd.Optional<Player> Target;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx)
        {
            var (target, _) = Target.Load();
            target ??= source as Player;
            if (target is null)
            {
                output.Error("Specify a player when using this command from the console.");
                return;
            }

            target.SetGameMode(Mode switch
            {
                GameModeValue.Survival => World.GameModeSurvival,
                GameModeValue.Creative => World.GameModeCreative,
                GameModeValue.Adventure => World.GameModeAdventure,
                GameModeValue.Spectator => World.GameModeSpectator,
                _ => World.GameModeSurvival,
            });
            output.Printf("Set {0}'s game mode to {1}.", target.Name(), Mode.ToString().ToLowerInvariant());
        }
    }
}
