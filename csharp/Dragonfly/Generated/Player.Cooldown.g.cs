// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.
#nullable enable
using System;

namespace Dragonfly;

public sealed partial class Player
{
    public bool HasCooldown(World.Item? item) => PluginBridge.Host.HasPlayerCooldown(_invocation, Id, item);
    public void SetCooldown(World.Item? item, TimeSpan cooldown) => PluginBridge.Host.SetPlayerCooldown(_invocation, Id, item, cooldown);
}
