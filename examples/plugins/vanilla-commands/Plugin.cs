using Dragonfly;

public sealed class VanillaCommands : Plugin
{
    public override void OnEnable()
    {
        HelpCommand.Register();
        GameModeCommand.Register();
        PingCommand.Register();
        PositionCommand.Register();
    }
}
