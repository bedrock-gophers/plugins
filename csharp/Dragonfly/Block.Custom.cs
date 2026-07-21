using System.Runtime.InteropServices;
using System.Text;
using Dragonfly.Native;

namespace Dragonfly;

public static partial class Block
{
    public enum CustomBlockRenderMethod { Opaque, AlphaTest, Blend, DoubleSided }

    public sealed record Custom : World.Block, World.Item
    {
        public Custom(string identifier) : this(identifier, new Dictionary<string, object?>()) { }
        internal Custom(string identifier, IReadOnlyDictionary<string, object?> properties)
        {
            Identifier = CustomBlockRegistry.ValidateIdentifier(identifier);
            Properties = CustomBlockRegistry.NormalizeProperties(Identifier, properties);
        }
        public string Identifier { get; }
        public IReadOnlyDictionary<string, object?> Properties { get; }

        public Custom WithState(string name, bool value) => WithStateValue(name, value);
        public Custom WithState(string name, byte value) => WithStateValue(name, value);
        public Custom WithState(string name, int value) => WithStateValue(name, value);
        public Custom WithState(string name, string value) => WithStateValue(name, value);

        private Custom WithStateValue(string name, object value)
        {
            ArgumentException.ThrowIfNullOrWhiteSpace(name);
            ArgumentNullException.ThrowIfNull(value);
            var properties = new Dictionary<string, object?>(Properties, StringComparer.Ordinal) { [name] = value };
            return new Custom(Identifier, properties);
        }
    }

    public static Custom RegisterCustom(
        string identifier,
        string name,
        byte[] texturePng,
        Item.CustomItemCategory category = Item.CustomItemCategory.Construction,
        int maxCount = 64,
        string group = "",
        byte[]? geometryJson = null,
        CustomBlockData? data = null) =>
        CustomBlockRegistry.Register(identifier, name, texturePng, category, maxCount, group, geometryJson, data);
}

internal static unsafe class CustomBlockRegistry
{
    private sealed class Entry : IDisposable
    {
        private readonly List<nint> _allocations = [];

        internal Entry(string identifier, string name, byte[] texture, Item.CustomItemCategory category,
            int maxCount, string group, byte[]? geometry, Block.CustomBlockData? data)
        {
            Descriptor = new CustomBlockDescriptor
            {
                Identifier = Allocate(Encoding.UTF8.GetBytes(identifier)),
                Name = Allocate(Encoding.UTF8.GetBytes(name)),
                TexturePng = Allocate(texture),
                GeometryJson = Allocate(geometry ?? []),
                Group = Allocate(Encoding.UTF8.GetBytes(group)),
                Category = (uint)category,
                MaxCount = maxCount,
                ComponentDataJson = Allocate(data?.ToJson() ?? []),
            };
        }

        internal CustomBlockDescriptor Descriptor { get; }

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
    private static readonly Dictionary<string, IReadOnlyDictionary<string, object?>> DefaultStates = new(StringComparer.Ordinal);

    internal static Block.Custom Register(string identifier, string name, byte[] texturePng,
        Item.CustomItemCategory category, int maxCount, string group, byte[]? geometryJson, Block.CustomBlockData? data)
    {
        identifier = ValidateIdentifier(identifier);
        ArgumentException.ThrowIfNullOrWhiteSpace(name);
        ArgumentNullException.ThrowIfNull(texturePng);
        if (texturePng.Length < 8 || !texturePng.AsSpan(0, 8).SequenceEqual(new byte[] { 137, 80, 78, 71, 13, 10, 26, 10 }))
            throw new ArgumentException("custom block texture must be a PNG image", nameof(texturePng));
        if (!Enum.IsDefined(category)) throw new ArgumentOutOfRangeException(nameof(category));
        if (maxCount is < 1 or > 64) throw new ArgumentOutOfRangeException(nameof(maxCount));
        ArgumentNullException.ThrowIfNull(group);
        if (!Identifiers.Add(identifier)) throw new InvalidOperationException($"custom block {identifier} is already registered");
        Entries.Add(new Entry(identifier, name, texturePng.ToArray(), category, maxCount, group, geometryJson?.ToArray(), data));
        var defaultState = data?.DefaultState() ?? new Dictionary<string, object?>();
        DefaultStates[identifier] = defaultState;
        return new Block.Custom(identifier, defaultState);
    }

    internal static string ValidateIdentifier(string identifier) => CustomItemRegistry.ValidateIdentifier(identifier);
    internal static bool Contains(string identifier) => Identifiers.Contains(identifier);
    internal static ulong Count => (ulong)Entries.Count;

    internal static IReadOnlyDictionary<string, object?> NormalizeProperties(string identifier, IReadOnlyDictionary<string, object?> properties)
    {
        var normalized = new Dictionary<string, object?>(properties, StringComparer.Ordinal);
        if (!DefaultStates.TryGetValue(identifier, out var defaults)) return normalized;
        foreach (var (name, defaultValue) in defaults)
        {
            if (!normalized.TryGetValue(name, out var value) || value is null || defaultValue is null) continue;
            normalized[name] = defaultValue switch
            {
                bool => Convert.ToInt64(value) != 0,
                byte => Convert.ToByte(value),
                int => Convert.ToInt32(value),
                string => Convert.ToString(value) ?? "",
                _ => value,
            };
        }
        return normalized;
    }

    internal static bool TryGet(ulong index, out CustomBlockDescriptor descriptor)
    {
        if (index >= (ulong)Entries.Count) { descriptor = default; return false; }
        descriptor = Entries[(int)index].Descriptor;
        return true;
    }

    internal static void Clear()
    {
        foreach (var entry in Entries) entry.Dispose();
        Entries.Clear();
        Identifiers.Clear();
        DefaultStates.Clear();
    }
}
