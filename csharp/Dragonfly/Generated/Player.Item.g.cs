// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.
#nullable enable
using Dragonfly.Native;

namespace Dragonfly;

public sealed partial class Player
{
    public Inventory.Value Inventory() => new(_invocation, Id, Abi.InventoryMain);
    public Inventory.Armour Armour() => new(_invocation, Id);
    public (Item.Stack MainHand, Item.Stack OffHand) HeldItems() =>
        PluginBridge.Host.HeldItems(_invocation, Id);
    public void SetHeldItems(Item.Stack mainHand, Item.Stack offHand) =>
        PluginBridge.Host.SetHeldItems(_invocation, Id, mainHand, offHand);
    public void SetHeldSlot(int to) => PluginBridge.Host.SetHeldSlot(_invocation, Id, to);
}
