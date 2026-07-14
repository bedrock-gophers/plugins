// Code generated from Dragonfly server/server.go Go AST. DO NOT EDIT.
#nullable enable
using System;
using System.Collections.Generic;
using Dragonfly.Native;

namespace Dragonfly;

public sealed partial class Server
{
    internal Server() { }

    public World World() => PluginBridge.Host.ServerWorld(Abi.WorldDimensionOverworld);

    public World Nether() => PluginBridge.Host.ServerWorld(Abi.WorldDimensionNether);

    public World End() => PluginBridge.Host.ServerWorld(Abi.WorldDimensionEnd);

    public int MaxPlayerCount() => PluginBridge.Host.ServerMaxPlayerCount();

    public int PlayerCount() => PluginBridge.Host.ServerPlayerCount();

    public IEnumerable<Player> Players(World.Tx? tx = null) =>
        PluginBridge.Host.ServerPlayers(tx?.Invocation ?? 0);

    public (World.EntityHandle? Player, bool Ok) Player(Guid uuid) =>
        PluginBridge.Host.ServerPlayer(uuid);

    public (World.EntityHandle? Player, bool Ok) PlayerByName(string name) =>
        PluginBridge.Host.ServerPlayerByName(name);

    public (World.EntityHandle? Player, bool Ok) PlayerByXUID(string xuid) =>
        PluginBridge.Host.ServerPlayerByXUID(xuid);
}

public abstract partial class Plugin
{
    public Server Server() => new();
}
