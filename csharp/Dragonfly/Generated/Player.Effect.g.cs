// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.
#nullable enable
using System.Collections.Generic;

namespace Dragonfly;

public sealed partial class Player
{
    public void AddEffect(Effect.Value e) => PluginBridge.Host.AddPlayerEffect(_invocation, Id, e);
    public void RemoveEffect(Effect.Type e) => PluginBridge.Host.RemovePlayerEffect(_invocation, Id, e);
    public (Effect.Value Effect, bool Ok) Effect(Effect.Type e) => PluginBridge.Host.PlayerEffect(_invocation, Id, e);
    public IReadOnlyList<Effect.Value> Effects() => PluginBridge.Host.PlayerEffects(_invocation, Id);
}
