// Code generated from bedrock-gophers/unsafe WritePacket Go AST. DO NOT EDIT.
#nullable enable
using System;

namespace Dragonfly;

public sealed partial class Player
{
    public void WritePacket(Packet.Packet pk)
    {
        ArgumentNullException.ThrowIfNull(pk);
        PluginBridge.Host.WritePlayerPacket(_invocation, Id, pk.HostHandle());
    }
}
