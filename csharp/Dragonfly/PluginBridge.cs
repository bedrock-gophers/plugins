using System.Runtime.CompilerServices;
using System.Runtime.InteropServices;
using System.Text;
using Dragonfly.Native;

namespace Dragonfly;

internal static unsafe class PluginBridge
{
    private static Func<Plugin>? Factory;
    private static PluginApi* Descriptor;

    internal static class Host
    {
        internal static HostApi* Api;

        internal static void SendPlayerText(ulong invocation, PlayerId player, uint kind, string message)
        {
            var api = Api;
            if (api is null || api->PlayerText == null) return;
            var bytes = Encoding.UTF8.GetBytes(message);
            fixed (byte* data = bytes)
            {
                _ = api->PlayerText(
                    api->Context,
                    invocation,
                    player,
                    kind,
                    new StringView { Data = data, Length = (ulong)bytes.Length });
            }
        }

        internal static void SetPlayerState(ulong invocation, PlayerId player, uint kind, PlayerStateValue value)
        {
            var api = Api;
            if (api is null || api->PlayerStateSet == null) return;
            _ = api->PlayerStateSet(api->Context, invocation, player, kind, value);
        }
    }

    internal static PluginApi* Initialize(Func<Plugin> factory, string id, ulong subscriptions)
    {
        if (Descriptor is not null) return Descriptor;
        Factory = factory;
        var bytes = Encoding.UTF8.GetBytes(id);
        var idPointer = (byte*)NativeMemory.Alloc((nuint)bytes.Length);
        bytes.CopyTo(new Span<byte>(idPointer, bytes.Length));
        Descriptor = (PluginApi*)NativeMemory.AllocZeroed((nuint)sizeof(PluginApi));
        *Descriptor = new PluginApi
        {
            Header = new AbiHeader
            {
                Version = Abi.PluginVersion,
                Size = (uint)sizeof(PluginApi),
                Subscriptions = subscriptions,
            },
            Id = new StringView { Data = idPointer, Length = (ulong)bytes.Length },
            Create = &Create,
            Enable = &Enable,
            Disable = &Disable,
            Commands = &Commands,
            HandleCommand = &HandleCommand,
            CommandEnumOptions = &CommandEnumOptions,
            SetHost = &SetHost,
            Destroy = &Destroy,
            HandleEvent = &HandleEvent,
        };
        return Descriptor;
    }

    [UnmanagedCallersOnly(CallConvs = [typeof(CallConvCdecl)])]
    private static void* Create()
    {
        try
        {
            CommandRegistry.Clear();
            return (void*)GCHandle.ToIntPtr(GCHandle.Alloc(Factory!()));
        }
        catch { return null; }
    }

    [UnmanagedCallersOnly(CallConvs = [typeof(CallConvCdecl)])]
    private static int SetHost(void* instance, void* host)
    {
        if (host is null) return Abi.Error;
        var header = (HostHeader*)host;
        if (header->Version != Abi.HostVersion || header->Size < 496) return Abi.Error;
        Host.Api = (HostApi*)host;
        return Abi.Ok;
    }

    [UnmanagedCallersOnly(CallConvs = [typeof(CallConvCdecl)])]
    private static int Enable(void* instance, StringBuffer* error)
    {
        try
        {
            Get(instance).OnEnable();
            if (error is not null) error->Length = 0;
            return Abi.Ok;
        }
        catch (Exception exception)
        {
            Write(error, exception.Message);
            return Abi.Error;
        }
    }

    [UnmanagedCallersOnly(CallConvs = [typeof(CallConvCdecl)])]
    private static int Disable(void* instance)
    {
        try { Get(instance).OnDisable(); return Abi.Ok; }
        catch { return Abi.Error; }
    }

    [UnmanagedCallersOnly(CallConvs = [typeof(CallConvCdecl)])]
    private static void Destroy(void* instance)
    {
        if (instance is not null) GCHandle.FromIntPtr((nint)instance).Free();
        Host.Api = null;
        CommandRegistry.Clear();
    }

    [UnmanagedCallersOnly(CallConvs = [typeof(CallConvCdecl)])]
    private static CommandDescriptor* Commands(void* instance, ulong* count)
    {
        try
        {
            if (instance is null) return null;
            return CommandRegistry.Native(count);
        }
        catch
        {
            if (count is not null) *count = 0;
            return null;
        }
    }

    [UnmanagedCallersOnly(CallConvs = [typeof(CallConvCdecl)])]
    private static int HandleCommand(void* instance, ulong index, CommandInput* input, CommandState* state)
    {
        try { return instance is null ? Abi.Error : CommandRegistry.Execute(index, input, state); }
        catch { return Abi.Error; }
    }

    [UnmanagedCallersOnly(CallConvs = [typeof(CallConvCdecl)])]
    private static int CommandEnumOptions(
        void* instance,
        ulong index,
        ulong overload,
        ulong parameter,
        CommandEnumContext* context,
        StringBuffer* output)
    {
        try
        {
            return instance is null
                ? Abi.Error
                : CommandRegistry.EnumOptions(index, overload, parameter, context, output);
        }
        catch { return Abi.Error; }
    }

    [UnmanagedCallersOnly(CallConvs = [typeof(CallConvCdecl)])]
    private static int HandleEvent(void* instance, uint eventId, void* input, void* state)
    {
        try
        {
            var plugin = Get(instance);
            switch (eventId)
            {
                case Abi.PlayerChatEvent:
                {
                    var value = (PlayerChatInput*)input;
                    var result = (PlayerChatState*)state;
                    var original = Utf8(value->Message);
                    var message = original;
                    var context = Event(value->Player, result->Cancelled, value->Invocation);
                    plugin.HandleChat(context, ref message);
                    ApplyCancellation(context, &result->Cancelled);
                    if (message != original)
                    {
                        if (!WriteExact(&result->Replacement, message)) return Abi.Error;
                        result->HasReplacement = 1;
                    }
                    return Abi.Ok;
                }
                case Abi.PlayerFoodLossEvent:
                {
                    var value = (PlayerFoodLossInput*)input;
                    var result = (PlayerFoodLossState*)state;
                    var context = Event(value->Player, result->Cancelled, value->Invocation);
                    var to = result->To;
                    plugin.HandleFoodLoss(context, value->From, ref to);
                    result->To = to;
                    ApplyCancellation(context, &result->Cancelled);
                    return Abi.Ok;
                }
                case Abi.PlayerJumpEvent:
                {
                    var value = (PlayerEventInput*)input;
                    plugin.HandleJump(new Player(value->Player, invocation: value->Invocation));
                    return Abi.Ok;
                }
                case Abi.PlayerMoveEvent:
                {
                    var value = (PlayerMoveInput*)input;
                    var result = (PlayerMoveState*)state;
                    var context = Event(value->Player, result->Cancelled, value->Invocation);
                    plugin.HandleMove(
                        context,
                        new Vector3(value->NewPosition.X, value->NewPosition.Y, value->NewPosition.Z),
                        new Rotation(value->Rotation.Yaw, value->Rotation.Pitch));
                    ApplyCancellation(context, &result->Cancelled);
                    return Abi.Ok;
                }
                case Abi.PlayerPunchAirEvent:
                {
                    var value = (PlayerEventInput*)input;
                    var result = (CancellableState*)state;
                    var context = Event(value->Player, result->Cancelled, value->Invocation);
                    plugin.HandlePunchAir(context);
                    ApplyCancellation(context, &result->Cancelled);
                    return Abi.Ok;
                }
                case Abi.PlayerQuitEvent:
                {
                    var value = (PlayerQuitInput*)input;
                    plugin.HandleQuit(new Player(value->Player, Utf8(value->Name), invocation: value->Invocation));
                    return Abi.Ok;
                }
                case Abi.PlayerTeleportEvent:
                {
                    var value = (PlayerTeleportInput*)input;
                    var result = (CancellableState*)state;
                    var context = Event(value->Player, result->Cancelled, value->Invocation);
                    plugin.HandleTeleport(context, new Vector3(value->Position.X, value->Position.Y, value->Position.Z));
                    ApplyCancellation(context, &result->Cancelled);
                    return Abi.Ok;
                }
                case Abi.PlayerToggleSneakEvent:
                case Abi.PlayerToggleSprintEvent:
                {
                    var value = (PlayerToggleInput*)input;
                    var result = (CancellableState*)state;
                    var context = Event(value->Player, result->Cancelled, value->Invocation);
                    if (eventId == Abi.PlayerToggleSneakEvent)
                        plugin.HandleToggleSneak(context, value->After != 0);
                    else
                        plugin.HandleToggleSprint(context, value->After != 0);
                    ApplyCancellation(context, &result->Cancelled);
                    return Abi.Ok;
                }
                default:
                    return Abi.Ok;
            }
        }
        catch { return Abi.Error; }
    }

    private static Plugin Get(void* instance) => (Plugin)GCHandle.FromIntPtr((nint)instance).Target!;

    private static Player.Context Event(PlayerId player, byte cancelled, ulong invocation = 0) =>
        new(new Player(player, invocation: invocation), cancelled != 0);

    private static void ApplyCancellation(Player.Context context, byte* cancelled)
    {
        if (context.Cancelled()) *cancelled = 1;
    }

    private static string Utf8(StringView value) => value.Length == 0
        ? string.Empty
        : Encoding.UTF8.GetString(new ReadOnlySpan<byte>(value.Data, checked((int)value.Length)));

    private static void Write(StringBuffer* output, string message)
    {
        if (output is null || output->Data is null || output->Capacity == 0) return;
        var bytes = Encoding.UTF8.GetBytes(message);
        var length = Math.Min(bytes.Length, checked((int)output->Capacity));
        bytes.AsSpan(0, length).CopyTo(new Span<byte>(output->Data, length));
        output->Length = (ulong)length;
    }

    private static bool WriteExact(StringBuffer* output, string message)
    {
        if (output is null) return false;
        var bytes = Encoding.UTF8.GetBytes(message);
        if ((ulong)bytes.Length > output->Capacity || bytes.Length != 0 && output->Data is null) return false;
        bytes.CopyTo(new Span<byte>(output->Data, bytes.Length));
        output->Length = (ulong)bytes.Length;
        return true;
    }
}
