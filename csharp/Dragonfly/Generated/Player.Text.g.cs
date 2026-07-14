// Code generated from Dragonfly server/player/player.go. DO NOT EDIT.
#nullable enable
using Dragonfly.Native;

namespace Dragonfly;

public sealed partial class Player
{
    public void Message(params object?[] a) => SendText(Abi.PlayerTextMessage, FormatArguments(a));
    public void SendPopup(params object?[] a) => SendText(Abi.PlayerTextPopup, FormatArguments(a));
    public void SendTip(params object?[] a) => SendText(Abi.PlayerTextTip, FormatArguments(a));
    public void SendJukeboxPopup(params object?[] a) => SendText(Abi.PlayerTextJukeboxPopup, FormatArguments(a));
    public void SetNameTag(string name) => SendText(Abi.PlayerTextNameTag, name);
    public void Disconnect(params object?[] msg) => SendText(Abi.PlayerTextDisconnect, FormatArguments(msg));
}
