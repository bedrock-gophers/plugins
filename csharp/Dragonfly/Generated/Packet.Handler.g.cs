// Code generated from bedrock-gophers/intercept Handler Go AST. DO NOT EDIT.
#nullable enable
namespace Dragonfly;

public abstract partial class Plugin
{
    [HandlerSubscription(18014398509481984UL)]
    public virtual void HandleClientPacket(Packet.Context ctx, Packet.Packet packet) { }

    [HandlerSubscription(36028797018963968UL)]
    public virtual void HandleServerPacket(Packet.Context ctx, Packet.Packet packet) { }
}
