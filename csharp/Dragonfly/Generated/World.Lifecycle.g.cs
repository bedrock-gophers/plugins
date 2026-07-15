// Code generated from Dragonfly server/world/world.go Go AST. DO NOT EDIT.
namespace Dragonfly;

public sealed partial class World
{
    public string Name() => PluginBridge.Host.WorldName(_invocation, Id) ?? string.Empty;
    public Cube.Range Range() => PluginBridge.Host.WorldRange(_invocation, Id);
    public int HighestLightBlocker(int x, int z) =>
        PluginBridge.Host.WorldHighestLightBlocker(_invocation, Id, x, z);
    public int Time() => PluginBridge.Host.WorldTime(_invocation, Id);
    public void SetTime(int @new) => PluginBridge.Host.SetWorldTime(_invocation, Id, @new);
    public void StopTime() => PluginBridge.Host.SetWorldTimeCycle(_invocation, Id, false);
    public void StartTime() => PluginBridge.Host.SetWorldTimeCycle(_invocation, Id, true);
    public bool TimeCycle() => PluginBridge.Host.WorldTimeCycle(_invocation, Id);
    public Cube.Pos Spawn() => PluginBridge.Host.WorldSpawn(_invocation, Id);
    public void SetSpawn(Cube.Pos pos) =>
        PluginBridge.Host.SetWorldSpawn(_invocation, Id, pos);
    public void SetRequiredSleepDuration(TimeSpan duration) =>
        PluginBridge.Host.SetWorldRequiredSleepDuration(_invocation, Id, duration);
    public GameMode DefaultGameMode() => PluginBridge.Host.WorldDefaultGameMode(_invocation, Id);
    public void SetTickRange(int v) => PluginBridge.Host.SetWorldTickRange(_invocation, Id, v);
    public void SetDefaultGameMode(GameMode mode) =>
        PluginBridge.Host.SetWorldDefaultGameMode(_invocation, Id, mode);
    public void SetDifficulty(Difficulty d) =>
        PluginBridge.Host.SetWorldDifficulty(_invocation, Id, d);
    public void Save() => PluginBridge.Host.SaveWorld(_invocation, Id);
    public void Close() => PluginBridge.Host.CloseWorld(_invocation, Id);
}

public static class WorldStateExtensions
{
    public static World.Dimension Dimension(this World world)
    {
        ArgumentNullException.ThrowIfNull(world);
        return PluginBridge.Host.WorldDimension(world.Invocation, world.Id);
    }

    public static World.Difficulty Difficulty(this World world)
    {
        ArgumentNullException.ThrowIfNull(world);
        return PluginBridge.Host.WorldDifficulty(world.Invocation, world.Id);
    }
}
