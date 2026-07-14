// Code generated from Dragonfly server/server.go Go AST. DO NOT EDIT.
#nullable enable
using System;
using System.Collections.Generic;

namespace Dragonfly;

public sealed partial class Server
{
    internal Server() { }

    public IEnumerable<Player> Players(World.Tx? tx = null) =>
        PluginBridge.Host.ServerPlayers(tx?.Invocation ?? 0);

    public (World.EntityHandle? Player, bool Ok) Player(Guid uuid) =>
        PluginBridge.Host.ServerPlayer(uuid);

    public (World.EntityHandle? Player, bool Ok) PlayerByName(string name) =>
        PluginBridge.Host.ServerPlayerByName(name);
}

public abstract partial class Plugin
{
    public Server Server() => new();
}
