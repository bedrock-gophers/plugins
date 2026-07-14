using Dragonfly;

internal static class HelpCommand
{
    public static void Register() => Cmd.Register(Cmd.New(
        "help",
        "Shows the available commands.",
        ["?"],
        new Help()));

    internal sealed class Help : Cmd.Runnable
    {
        public Cmd.Optional<int> Page;

        public void Run(Cmd.Source source, Cmd.Output output, World.Tx? tx)
        {
            if (Page.LoadOr(1) != 1)
            {
                output.Error("That help page does not exist.");
                return;
            }

            output.Print("Available commands: /gamemode, /help, /ping, /position");
        }
    }
}
