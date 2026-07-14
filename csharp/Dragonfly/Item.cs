#nullable enable
using System.Globalization;
using Dragonfly.Native;

namespace Dragonfly;

public static partial class Item
{
    public static Stack NewStack(World.Item item, int count)
    {
        ArgumentNullException.ThrowIfNull(item);
        if (count < 0) throw new ArgumentOutOfRangeException(nameof(count));
        return new Stack(item, count);
    }

    public readonly struct Stack
    {
        private readonly World.Item? _item;
        private readonly int _count;
        private readonly uint _damage;
        private readonly bool _unbreakable;
        private readonly int _anvilCost;
        private readonly string? _customName;
        private readonly string[]? _lore;
        private readonly byte[]? _itemNbt;
        private readonly byte[]? _valuesNbt;
        private readonly ItemEnchantment[]? _enchantments;

        internal Stack(
            World.Item? item,
            int count,
            uint damage = 0,
            bool unbreakable = false,
            int anvilCost = 0,
            string? customName = null,
            string[]? lore = null,
            byte[]? itemNbt = null,
            byte[]? valuesNbt = null,
            ItemEnchantment[]? enchantments = null)
        {
            _item = item;
            _count = Math.Max(count, 0);
            _damage = damage;
            _unbreakable = unbreakable;
            _anvilCost = anvilCost;
            _customName = customName;
            _lore = lore;
            _itemNbt = itemNbt;
            _valuesNbt = valuesNbt;
            _enchantments = enchantments;
        }

        public int Count() => _count;
        public bool Empty() => _count == 0 || _item is null || ItemCodec.IsAir(_item);
        public World.Item? Item() => Empty() ? null : _item;

        public Stack Grow(int count) => Copy(count: Math.Max(0, _count + count));

        public string CustomName() => _customName ?? string.Empty;

        public Stack WithCustomName(params object?[] values) => Copy(
            customName: string.Join(" ", values.Select(FormatValue)));

        public string[] Lore() => Empty() ? [] : (string[])(_lore?.Clone() ?? Array.Empty<string>());

        public Stack WithLore(params string[] lines)
        {
            ArgumentNullException.ThrowIfNull(lines);
            return Copy(lore: (string[])lines.Clone());
        }

        internal uint Damage => _damage;
        internal bool IsUnbreakable => _unbreakable;
        internal int AnvilCostValue => _anvilCost;
        internal byte[] ItemNbt => _itemNbt ?? [];
        internal byte[] ValuesNbt => _valuesNbt ?? [];
        internal ItemEnchantment[] Enchantments => _enchantments ?? [];

        internal bool TryEncode(out string identifier, out int metadata)
        {
            if (_item is not null) return ItemCodec.TryEncode(_item, out identifier, out metadata);
            identifier = string.Empty;
            metadata = 0;
            return false;
        }

        private Stack Copy(
            int? count = null,
            string? customName = null,
            string[]? lore = null) => new(
                _item,
                count ?? _count,
                _damage,
                _unbreakable,
                _anvilCost,
                customName ?? _customName,
                lore ?? _lore,
                _itemNbt,
                _valuesNbt,
                _enchantments);

        private static string FormatValue(object? value) => value switch
        {
            null => "<nil>",
            bool boolean => boolean ? "true" : "false",
            IFormattable formattable => formattable.ToString(null, CultureInfo.InvariantCulture) ?? string.Empty,
            _ => value.ToString() ?? string.Empty,
        };
    }
}
