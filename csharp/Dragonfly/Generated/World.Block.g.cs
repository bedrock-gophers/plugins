// Code generated from Dragonfly server/world Go AST. DO NOT EDIT.
#nullable enable

namespace Dragonfly;

public sealed partial class World
{
    public interface Block { }

    public sealed class SetOpts
    {
        public bool DisableBlockUpdates;
        public bool DisableLiquidDisplacement;
        public bool DisableRedstoneUpdates;
    }

    public partial class Tx
    {
        public Block Block(Cube.Pos position) => PluginBridge.Host.WorldBlock(Invocation, position);

        public void SetBlock(Cube.Pos position, Block? block, SetOpts? options = null) =>
            PluginBridge.Host.SetWorldBlock(Invocation, position, block, options);
    }
}
