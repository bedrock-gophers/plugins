using System.Runtime.InteropServices;
using System.Text;
using Dragonfly.Native;

namespace Dragonfly;

public static partial class Item
{
    public enum CustomItemCategory
    {
        Construction = 1,
        Nature = 2,
        Equipment = 3,
        Items = 4,
    }

    public sealed record Custom : World.Item
    {
        public Custom(string identifier) => Identifier = CustomItemRegistry.ValidateIdentifier(identifier);

        public string Identifier { get; }
    }

    public static Custom RegisterCustom(
        string identifier,
        string name,
        byte[] texturePng,
        CustomItemCategory category = CustomItemCategory.Items,
        int maxCount = 64,
        string group = "",
        CustomItemData? data = null) =>
        CustomItemRegistry.Register(identifier, name, texturePng, category, maxCount, group, data);
}

internal static unsafe class CustomItemRegistry
{
    private sealed class Entry : IDisposable
    {
        private readonly List<nint> _allocations = [];

        internal Entry(string identifier, string name, byte[] texture, Item.CustomItemCategory category, int maxCount, string group, Item.CustomItemData? data)
        {
            Descriptor = new CustomItemDescriptor
            {
                Identifier = Allocate(Encoding.UTF8.GetBytes(identifier)),
                Name = Allocate(Encoding.UTF8.GetBytes(name)),
                TexturePng = Allocate(texture),
                Group = Allocate(Encoding.UTF8.GetBytes(group)),
                Category = (uint)category,
                MaxCount = maxCount,
                ComponentDataJson = Allocate(data?.ToJson() ?? []),
            };
        }

        internal CustomItemDescriptor Descriptor { get; }

        private StringView Allocate(ReadOnlySpan<byte> value)
        {
            if (value.Length == 0) return default;
            var pointer = NativeMemory.Alloc((nuint)value.Length);
            value.CopyTo(new Span<byte>(pointer, value.Length));
            _allocations.Add((nint)pointer);
            return new StringView { Data = (byte*)pointer, Length = (ulong)value.Length };
        }

        public void Dispose()
        {
            foreach (var pointer in _allocations) NativeMemory.Free((void*)pointer);
            _allocations.Clear();
        }
    }

    private static readonly List<Entry> Entries = [];
    private static readonly HashSet<string> Identifiers = new(StringComparer.Ordinal);

    internal static Item.Custom Register(string identifier, string name, byte[] texturePng, Item.CustomItemCategory category, int maxCount, string group, Item.CustomItemData? data)
    {
        identifier = ValidateIdentifier(identifier);
        ArgumentException.ThrowIfNullOrWhiteSpace(name);
        ArgumentNullException.ThrowIfNull(texturePng);
        if (texturePng.Length < 8 || !texturePng.AsSpan(0, 8).SequenceEqual(new byte[] { 137, 80, 78, 71, 13, 10, 26, 10 }))
            throw new ArgumentException("custom item texture must be a PNG image", nameof(texturePng));
        if (!Enum.IsDefined(category)) throw new ArgumentOutOfRangeException(nameof(category));
        if (maxCount is < 1 or > 64) throw new ArgumentOutOfRangeException(nameof(maxCount));
        ArgumentNullException.ThrowIfNull(group);
        if (!Identifiers.Add(identifier)) throw new InvalidOperationException($"custom item {identifier} is already registered");
        Entries.Add(new Entry(identifier, name, texturePng.ToArray(), category, maxCount, group, data));
        return new Item.Custom(identifier);
    }

    internal static string ValidateIdentifier(string identifier)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(identifier);
        var parts = identifier.Split(':');
        if (identifier.Length > 256 || identifier.Count(character => character == ':') != 1 ||
            parts.Length != 2 || parts[0].Length == 0 || parts[1].Length == 0 ||
            string.Equals(parts[0], "minecraft", StringComparison.OrdinalIgnoreCase) ||
            identifier.Any(character => !(char.IsAsciiLetterOrDigit(character) || character is '_' or '-' or '.' or ':' or '/')))
            throw new ArgumentException("custom item identifier must use the namespace:name format", nameof(identifier));
        return identifier.ToLowerInvariant();
    }

    internal static bool Contains(string identifier) => Identifiers.Contains(identifier);

    internal static ulong Count => (ulong)Entries.Count;

    internal static bool TryGet(ulong index, out CustomItemDescriptor descriptor)
    {
        if (index >= (ulong)Entries.Count)
        {
            descriptor = default;
            return false;
        }
        descriptor = Entries[(int)index].Descriptor;
        return true;
    }

    internal static void Clear()
    {
        foreach (var entry in Entries) entry.Dispose();
        Entries.Clear();
        Identifiers.Clear();
    }
}
