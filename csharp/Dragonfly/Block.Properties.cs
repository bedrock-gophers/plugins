#nullable enable
namespace Dragonfly;

internal static class BlockPropertyCodec
{
    internal static byte[] Encode(IReadOnlyDictionary<string, object?>? properties)
    {
        var root = new Nbt.Compound();
        if (properties is null) return Nbt.Encode(root);
        foreach (var (name, value) in properties)
        {
            ArgumentException.ThrowIfNullOrEmpty(name);
            var encoded = value switch
            {
                bool boolean => Entry(0, Nbt.Value.Byte(boolean ? (byte)1 : (byte)0)),
                byte number => Entry(1, Nbt.Value.Byte(number)),
                int number => Entry(2, Nbt.Value.Int(number)),
                string text => Entry(3, Nbt.Value.String(text)),
                null => throw new ArgumentException("block properties cannot contain null values", nameof(properties)),
                _ => throw new ArgumentException(
                    $"block property type {value.GetType()} is not supported by Dragonfly",
                    nameof(properties)),
            };
            root.Add(name, Nbt.Value.Compound(encoded));
        }
        return Nbt.Encode(root);
    }

    internal static IReadOnlyDictionary<string, object?> Decode(ReadOnlySpan<byte> data)
    {
        var result = new Dictionary<string, object?>(StringComparer.Ordinal);
        foreach (var (name, encoded) in Nbt.Decode(data))
        {
            if (encoded.Type != Nbt.TagType.Compound) continue;
            var entry = encoded.AsCompound();
            if (!entry.TryGetValue("kind", out var kind) || kind.Type != Nbt.TagType.Int ||
                !entry.TryGetValue("value", out var value)) continue;
            result[name] = kind.AsInt() switch
            {
                0 when value.Type == Nbt.TagType.Byte => value.AsByte() != 0,
                1 when value.Type == Nbt.TagType.Byte => value.AsByte(),
                2 when value.Type == Nbt.TagType.Int => value.AsInt(),
                3 when value.Type == Nbt.TagType.String => value.AsString(),
                _ => null,
            };
        }
        return result;
    }

    private static Nbt.Compound Entry(int kind, Nbt.Value value) => new()
    {
        ["kind"] = Nbt.Value.Int(kind),
        ["value"] = value,
    };
}
