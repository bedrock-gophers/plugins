namespace Dragonfly;

public sealed partial class Player
{
    public void SendContainerMenu(ContainerMenu.Value menu, bool update = false) =>
        PluginBridge.Host.SendPlayerInventoryMenu(Invocation, Id, menu, update);

    public void CloseContainerMenu() =>
        PluginBridge.Host.ClosePlayerInventoryMenu(Invocation, Id);
}
