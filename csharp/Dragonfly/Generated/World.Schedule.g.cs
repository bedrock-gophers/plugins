// Code generated from Dragonfly server/world/task.go Go AST. DO NOT EDIT.
using System;

namespace Dragonfly;

public sealed partial class World
{
    public void Schedule(Action<World.Tx> callback) =>
        PluginBridge.Host.ScheduleWorld(this, callback);
}
