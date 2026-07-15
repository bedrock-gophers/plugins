// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.
namespace Dragonfly;

public sealed partial class Player
{
    public void PlaySound(World.Sound sound) =>
        PluginBridge.Host.PlayPlayerSound(_invocation, Id, sound);
}
