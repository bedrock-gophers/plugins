using System.Text;
using Dragonfly.Native;

namespace Dragonfly.Packet;

public interface Packet
{
    uint ID();
}

public sealed class Context
{
    private bool _cancelled;

    internal Context(string xuid, bool cancelled)
    {
        _xuid = xuid;
        _cancelled = cancelled;
    }

    private readonly string _xuid;
    public string XUID() => _xuid;
    public bool Cancelled() => _cancelled;
    public void Cancel() => _cancelled = true;
}

public readonly record struct Vector2(float X, float Y);
public readonly record struct Vector3(float X, float Y, float Z);

public sealed class Value
{
    private readonly ulong _packet;
    private readonly uint _field;

    internal Value(ulong packet, int field)
    {
        _packet = packet;
        _field = checked((uint)field);
    }

    // Json serialises complex protocol values lazily on the Go side. Typed
    // nested protocol proxies will replace this fallback as AST coverage grows.
    public string Json() => Encoding.UTF8.GetString(PacketBridge.Data(_packet, _field, 10));
}

public sealed class Unknown : Packet
{
    private readonly ulong _handle;
    private readonly uint _id;

    internal Unknown(ulong handle, uint id) => (_handle, _id) = (handle, id);
    public uint ID() => _id;
    public uint PacketID { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte[] Payload { get => PacketBridge.Bytes(_handle, 1); set => PacketBridge.SetBytes(_handle, 1, value); }
}

internal static unsafe class PacketBridge
{
    private const ulong MaxDataBytes = 16UL << 20;

    internal static bool Bool(ulong packet, int field) => Get(packet, field, 1).Unsigned != 0;
    internal static long Signed(ulong packet, int field) => Get(packet, field, 2).Signed;
    internal static ulong Unsigned(ulong packet, int field) => Get(packet, field, 3).Unsigned;
    internal static double Number(ulong packet, int field) => Get(packet, field, 4).Number;
    internal static string String(ulong packet, int field) => Encoding.UTF8.GetString(Data(packet, checked((uint)field), 5));
    internal static byte[] Bytes(ulong packet, int field) => Data(packet, checked((uint)field), 6);
    internal static Vector2 Vector2(ulong packet, int field)
    {
        var value = Get(packet, field, 7);
        return new((float)value.X, (float)value.Y);
    }
    internal static Vector3 Vector3(ulong packet, int field)
    {
        var value = Get(packet, field, 8);
        return new((float)value.X, (float)value.Y, (float)value.Z);
    }
    internal static Guid Guid(ulong packet, int field)
    {
        var value = Get(packet, field, 9);
        return new Guid(new ReadOnlySpan<byte>(value.Uuid.Bytes, 16), bigEndian: true);
    }

    internal static void SetBool(ulong packet, int field, bool value) =>
        Set(packet, field, new PacketFieldValue { Kind = 1, Unsigned = value ? 1UL : 0UL });
    internal static void SetSigned(ulong packet, int field, long value) =>
        Set(packet, field, new PacketFieldValue { Kind = 2, Signed = value });
    internal static void SetUnsigned(ulong packet, int field, ulong value) =>
        Set(packet, field, new PacketFieldValue { Kind = 3, Unsigned = value });
    internal static void SetNumber(ulong packet, int field, double value) =>
        Set(packet, field, new PacketFieldValue { Kind = 4, Number = value });
    internal static void SetString(ulong packet, int field, string value) =>
        SetData(packet, field, 5, Encoding.UTF8.GetBytes(value));
    internal static void SetBytes(ulong packet, int field, byte[] value) => SetData(packet, field, 6, value);
    internal static void SetVector2(ulong packet, int field, Vector2 value) =>
        Set(packet, field, new PacketFieldValue { Kind = 7, X = value.X, Y = value.Y });
    internal static void SetVector3(ulong packet, int field, Vector3 value) =>
        Set(packet, field, new PacketFieldValue { Kind = 8, X = value.X, Y = value.Y, Z = value.Z });
    internal static void SetGuid(ulong packet, int field, Guid value)
    {
        var native = new PacketFieldValue { Kind = 9 };
        value.TryWriteBytes(new Span<byte>(native.Uuid.Bytes, 16), bigEndian: true, out _);
        Set(packet, field, native);
    }

    private static PacketFieldValue Get(ulong packet, int field, uint kind)
    {
        var api = PluginBridge.Host.Api;
        if (api is null || api->PacketFieldGet == null) throw Expired();
        PacketFieldValue value = default;
        if (api->PacketFieldGet(api->Context, packet, checked((uint)field), &value) != Abi.Ok || value.Kind != kind)
            throw Expired();
        return value;
    }

    internal static byte[] Data(ulong packet, uint field, uint kind)
    {
        var api = PluginBridge.Host.Api;
        if (api is null || api->PacketFieldGet == null) throw Expired();
        PacketFieldValue value = default;
        if (api->PacketFieldGet(api->Context, packet, field, &value) != Abi.Ok || value.Kind != kind || value.Data.Length > MaxDataBytes)
            throw Expired();
        if (value.Data.Length == 0) return [];
        var data = new byte[checked((int)value.Data.Length)];
        fixed (byte* pointer = data)
        {
            value.Data = new StringBuffer { Data = pointer, Capacity = (ulong)data.Length };
            if (api->PacketFieldGet(api->Context, packet, field, &value) != Abi.Ok || value.Kind != kind || value.Data.Length != (ulong)data.Length)
                throw Expired();
        }
        return data;
    }

    private static void SetData(ulong packet, int field, uint kind, byte[] data)
    {
        ArgumentNullException.ThrowIfNull(data);
        if ((ulong)data.Length > MaxDataBytes) throw new ArgumentOutOfRangeException(nameof(data));
        fixed (byte* pointer = data)
        {
            var value = new PacketFieldValue
            {
                Kind = kind,
                Data = new StringBuffer { Data = pointer, Length = (ulong)data.Length, Capacity = (ulong)data.Length },
            };
            Set(packet, field, value);
        }
    }

    private static void Set(ulong packet, int field, PacketFieldValue value)
    {
        var api = PluginBridge.Host.Api;
        if (api is null || api->PacketFieldSet == null ||
            api->PacketFieldSet(api->Context, packet, checked((uint)field), &value) != Abi.Ok)
            throw Expired();
    }

    private static InvalidOperationException Expired() => new("packet is no longer available");
}
