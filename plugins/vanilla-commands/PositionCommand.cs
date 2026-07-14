using Dragonfly;

internal static class PositionCommand
{
    public static void Register() => Cmd.Register(Cmd.New(
        "position",
        "Shows your current position.",
        ["pos"],
        new Position()));

    internal sealed class Position : Cmd.Runnable, Cmd.Allower
    {
        public bool Allow(Cmd.Source source) => source is Player;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx)
        {
            var position = source.Position();
            output.Printf("Position: {0:0.##}, {1:0.##}, {2:0.##}", position.X, position.Y, position.Z);
        }
    }
}
