namespace Dragonfly;

public static class ContainerMenu
{
    public enum Type : uint
    {
        Chest,
        DoubleChest,
        Hopper,
        Dropper,
        Barrel,
        EnderChest,
    }

    public interface Value
    {
        string Title();
        Type ContainerType();
        IReadOnlyList<Item.Stack> Items();
        void Submit(Player player, int slot, Item.Stack item, World.Tx transaction);
        void Close(Player player, World.Tx transaction) { }
    }

    public static int Size(Type type) => type switch
    {
        Type.Chest or Type.Barrel or Type.EnderChest => 27,
        Type.DoubleChest => 54,
        Type.Hopper => 5,
        Type.Dropper => 9,
        _ => throw new ArgumentOutOfRangeException(nameof(type), type, "Unknown inventory menu container."),
    };
}
