using Dragonfly;

internal static class PingCommand
{
    public static void Register() => Cmd.Register(Cmd.New(
        "ping",
        "Shows a player's network latency.",
        [],
        new Ping()));

    internal sealed class Ping : Cmd.Runnable
    {
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

            output.Printf("{0}'s ping is {1:0}ms.", target.Name(), target.Latency().TotalMilliseconds);
        }
    }
}
