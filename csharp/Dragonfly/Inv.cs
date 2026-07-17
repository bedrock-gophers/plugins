#nullable enable

namespace Dragonfly;

/// <summary>Client-side inventory menus backed by fake block containers.</summary>
public static class Inv
{
    public enum Container
    {
        Chest,
        DoubleChest,
        Hopper,
        Dropper,
        Barrel,
        EnderChest,
    }

    public interface Submittable
    {
        void Submit(Player player, Item.Stack item, World.Tx tx);
    }

    public interface Closer
    {
        void Close(Player player, World.Tx tx);
    }

    public sealed class Menu
    {
        private readonly Item.Stack[] _stacks;

        internal Menu(Submittable submittable, string name, Container container, Item.Stack[] stacks)
        {
            Submittable = submittable;
            Name = name;
            Container = container;
            _stacks = stacks;
        }

        internal Submittable Submittable { get; }
        internal string Name { get; }
        internal Container Container { get; }
        internal Item.Stack[] Stacks => (Item.Stack[])_stacks.Clone();

        public int Size() => ContainerSize(Container);

        public Menu WithStacks(params Item.Stack[] stacks)
        {
            ArgumentNullException.ThrowIfNull(stacks);
            if (stacks.Length > Size())
                throw new ArgumentException("too many stacks for the container", nameof(stacks));
            var contents = new Item.Stack[Size()];
            stacks.CopyTo(contents, 0);
            return new Menu(Submittable, Name, Container, contents);
        }

        public Menu WithStack(int slot, Item.Stack stack)
        {
            if ((uint)slot >= (uint)Size()) throw new ArgumentOutOfRangeException(nameof(slot));
            var contents = (Item.Stack[])_stacks.Clone();
            contents[slot] = stack;
            return new Menu(Submittable, Name, Container, contents);
        }
    }

    public static Menu NewMenu(Submittable submittable, string name, Container container)
    {
        ArgumentNullException.ThrowIfNull(submittable);
        ArgumentNullException.ThrowIfNull(name);
        return new Menu(submittable, name, container, new Item.Stack[ContainerSize(container)]);
    }

    public static void SendMenu(Player player, Menu menu)
    {
        ArgumentNullException.ThrowIfNull(player);
        ArgumentNullException.ThrowIfNull(menu);
        PluginBridge.Host.SendPlayerInventoryMenu(player.Invocation, player.Id, menu, update: false);
    }

    public static void UpdateMenu(Player player, Menu menu)
    {
        ArgumentNullException.ThrowIfNull(player);
        ArgumentNullException.ThrowIfNull(menu);
        PluginBridge.Host.SendPlayerInventoryMenu(player.Invocation, player.Id, menu, update: true);
    }

    public static void CloseContainer(Player player)
    {
        ArgumentNullException.ThrowIfNull(player);
        PluginBridge.Host.ClosePlayerInventoryMenu(player.Invocation, player.Id);
    }

    internal static int ContainerSize(Container container) => container switch
    {
        Container.Chest or Container.Barrel or Container.EnderChest => 27,
        Container.DoubleChest => 54,
        Container.Hopper => 5,
        Container.Dropper => 9,
        _ => throw new ArgumentOutOfRangeException(nameof(container)),
    };
}
