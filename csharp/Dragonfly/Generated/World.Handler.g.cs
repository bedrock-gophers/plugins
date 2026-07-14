// Code generated from Dragonfly server/world/handler.go. DO NOT EDIT.
#nullable enable
namespace Dragonfly;

public sealed partial class World
{
    public interface Sound { }

    public enum RedstoneUpdateCause
    {
        BlockUpdate = 0,
        ScheduledTick = 1,
        CompilerRebuild = 2,
    }

    public readonly record struct RedstoneUpdate(
        Cube.Pos Pos,
        Cube.Pos ChangedNeighbour,
        bool HasChangedNeighbour,
        bool ChangedRedstoneRelevant,
        Cube.Pos Source,
        bool HasSource,
        World.Block Before,
        World.Block? After,
        int OldPower,
        int NewPower,
        long CurrentTick,
        World.RedstoneUpdateCause Cause
    );

    public interface Handler
    {
        void HandleLiquidFlow(World.Context ctx, Cube.Pos from, Cube.Pos into, World.Liquid liquid, World.Block replaced);
        void HandleLiquidDecay(World.Context ctx, Cube.Pos pos, World.Liquid before, World.Liquid? after);
        void HandleLiquidHarden(World.Context ctx, Cube.Pos hardenedPos, World.Block liquidHardened, World.Block otherLiquid, World.Block newBlock);
        void HandleSound(World.Context ctx, World.Sound s, Vector3 pos);
        void HandleFireSpread(World.Context ctx, Cube.Pos from, Cube.Pos to);
        void HandleBlockBurn(World.Context ctx, Cube.Pos pos);
        void HandleCropTrample(World.Context ctx, Cube.Pos pos);
        void HandleLeavesDecay(World.Context ctx, Cube.Pos pos);
        void HandleEntitySpawn(World.Tx tx, World.Entity e);
        void HandleEntityDespawn(World.Tx tx, World.Entity e);
        void HandleExplosion(World.Context ctx, Vector3 position, ref World.Entity[] entities, ref Cube.Pos[] blocks, ref double itemDropChance, ref bool spawnFire);
        void HandleRedstoneUpdate(World.Context ctx, World.RedstoneUpdate update);
        void HandleClose(World.Tx tx);
    }
}

public abstract partial class Plugin : World.Handler
{
    [HandlerSubscription(2199023255552UL)]
    public virtual void HandleLiquidFlow(World.Context ctx, Cube.Pos from, Cube.Pos into, World.Liquid liquid, World.Block replaced) { }
    [HandlerSubscription(4398046511104UL)]
    public virtual void HandleLiquidDecay(World.Context ctx, Cube.Pos pos, World.Liquid before, World.Liquid? after) { }
    [HandlerSubscription(8796093022208UL)]
    public virtual void HandleLiquidHarden(World.Context ctx, Cube.Pos hardenedPos, World.Block liquidHardened, World.Block otherLiquid, World.Block newBlock) { }
    [HandlerSubscription(17592186044416UL)]
    public virtual void HandleSound(World.Context ctx, World.Sound s, Vector3 pos) { }
    [HandlerSubscription(35184372088832UL)]
    public virtual void HandleFireSpread(World.Context ctx, Cube.Pos from, Cube.Pos to) { }
    [HandlerSubscription(70368744177664UL)]
    public virtual void HandleBlockBurn(World.Context ctx, Cube.Pos pos) { }
    [HandlerSubscription(140737488355328UL)]
    public virtual void HandleCropTrample(World.Context ctx, Cube.Pos pos) { }
    [HandlerSubscription(281474976710656UL)]
    public virtual void HandleLeavesDecay(World.Context ctx, Cube.Pos pos) { }
    [HandlerSubscription(562949953421312UL)]
    public virtual void HandleEntitySpawn(World.Tx tx, World.Entity e) { }
    [HandlerSubscription(1125899906842624UL)]
    public virtual void HandleEntityDespawn(World.Tx tx, World.Entity e) { }
    [HandlerSubscription(2251799813685248UL)]
    public virtual void HandleExplosion(World.Context ctx, Vector3 position, ref World.Entity[] entities, ref Cube.Pos[] blocks, ref double itemDropChance, ref bool spawnFire) { }
    [HandlerSubscription(4503599627370496UL)]
    public virtual void HandleRedstoneUpdate(World.Context ctx, World.RedstoneUpdate update) { }
    [HandlerSubscription(9007199254740992UL)]
    public virtual void HandleClose(World.Tx tx) { }
}
