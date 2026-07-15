// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.
using Dragonfly.Native;

namespace Dragonfly;

public sealed partial class Player
{
    public (int Added, bool Ok) Collect(Item.Stack s) =>
        PluginBridge.Host.RunPlayerItemAction(_invocation, Id, s, Abi.PlayerItemActionCollect);
    public int Drop(Item.Stack s) =>
        PluginBridge.Host.RunPlayerItemAction(_invocation, Id, s, Abi.PlayerItemActionDrop).Count;
}
