// Code generated from Dragonfly server/world/tx.go Go AST. DO NOT EDIT.
#nullable enable
using System;
using Dragonfly.Native;

namespace Dragonfly;

public sealed partial class World
{
    public partial class Tx
    {
        public World.Task Defer(Action<Tx> f) =>
            PluginBridge.Host.DeferWorld(Invocation, f, Abi.WorldDeferDefer);
        public World.Task DeferErr(Func<Tx, Exception?> f) =>
            PluginBridge.Host.DeferWorld(Invocation, f, Abi.WorldDeferDeferErr);
    }
}
