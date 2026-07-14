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

    private static Nbt.Compound Entry(int kind, Nbt.Value value) => new()
    {
        ["kind"] = Nbt.Value.Int(kind),
        ["value"] = value,
    };
}
