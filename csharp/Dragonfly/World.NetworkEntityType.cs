namespace Dragonfly;

public sealed partial class World
{
    private static readonly object RegisteredEntityTypesGate = new();
    private static readonly List<EntityType> RegisteredEntityTypes = [];

    /// <summary>
    /// Optionally gives a custom entity type a distinct identifier used when
    /// spawning it to clients. EntityType.EncodeEntity remains the save ID.
    /// </summary>
    public interface NetworkEntityType : EntityType
    {
        string NetworkEncodeEntity();
    }

    /// <summary>
    /// Registers an entity type created at runtime. Call this while the plugin
    /// is being constructed, before it is enabled.
    /// </summary>
    public static T RegisterEntityType<T>(T type) where T : EntityType
    {
        ArgumentNullException.ThrowIfNull(type);
        lock (RegisteredEntityTypesGate) RegisteredEntityTypes.Add(type);
        return type;
    }

    internal static EntityType[] SnapshotRegisteredEntityTypes()
    {
        lock (RegisteredEntityTypesGate) return RegisteredEntityTypes.ToArray();
    }

    internal static void ClearRegisteredEntityTypes()
    {
        lock (RegisteredEntityTypesGate) RegisteredEntityTypes.Clear();
    }
}
