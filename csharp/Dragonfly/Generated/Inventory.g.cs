// Code generated from Dragonfly server/item/inventory Go AST. DO NOT EDIT.
#nullable enable
using Dragonfly.Native;

namespace Dragonfly;

public static class Inventory
{
    public sealed class Value
    {
        private readonly ulong _invocation;
        private readonly InventoryId _id;

        internal Value(ulong invocation, PlayerId player, uint kind)
        {
            _invocation = invocation;
            _id = new InventoryId { Player = player, Kind = kind };
        }

        public int Size() => PluginBridge.Host.InventorySize(_invocation, _id);

        public Item.Stack Item(int slot)
        {
            CheckSlot(slot);
            return PluginBridge.Host.InventoryItem(_invocation, _id, slot);
        }

        public void SetItem(int slot, Item.Stack item)
        {
            CheckSlot(slot);
            PluginBridge.Host.SetInventoryItem(_invocation, _id, slot, item);
        }

        public int AddItem(Item.Stack item) =>
            PluginBridge.Host.AddInventoryItem(_invocation, _id, item);

        private void CheckSlot(int slot)
        {
            if (slot < 0 || slot >= Size()) throw new ArgumentOutOfRangeException(nameof(slot));
        }
    }

    public sealed class Armour
    {
        private readonly Value _inventory;

        internal Armour(ulong invocation, PlayerId player) =>
            _inventory = new Value(invocation, player, Abi.InventoryArmour);

        public Item.Stack Helmet() => _inventory.Item(0);
        public Item.Stack Chestplate() => _inventory.Item(1);
        public Item.Stack Leggings() => _inventory.Item(2);
        public Item.Stack Boots() => _inventory.Item(3);

        public void SetHelmet(Item.Stack helmet) => _inventory.SetItem(0, helmet);
        public void SetChestplate(Item.Stack chestplate) => _inventory.SetItem(1, chestplate);
        public void SetLeggings(Item.Stack leggings) => _inventory.SetItem(2, leggings);
        public void SetBoots(Item.Stack boots) => _inventory.SetItem(3, boots);

        public void Set(Item.Stack helmet, Item.Stack chestplate, Item.Stack leggings, Item.Stack boots)
        {
            SetHelmet(helmet);
            SetChestplate(chestplate);
            SetLeggings(leggings);
            SetBoots(boots);
        }

        public Value Inventory() => _inventory;
    }
}
