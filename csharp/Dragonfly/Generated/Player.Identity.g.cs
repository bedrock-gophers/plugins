// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.
using System;

namespace Dragonfly;

public sealed partial class Player
{
    public string Name() => PlayerName;
    public Guid UUID() => PluginBridge.Host.PlayerUUID(Id);
    public string XUID() => PluginBridge.Host.PlayerXUID(_invocation, Id);
}
