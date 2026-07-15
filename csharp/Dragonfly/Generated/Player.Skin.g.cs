// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.
namespace Dragonfly;

public sealed partial class Player
{
    public Skin Skin() => PluginBridge.Host.PlayerSkin(_invocation, Id);
    public void SetSkin(Skin skin) => PluginBridge.Host.SetPlayerSkin(_invocation, Id, skin);
}
