// Code generated from gophertunnel minecraft/protocol/packet Go AST. DO NOT EDIT.
#nullable enable
using System;

namespace Dragonfly.Packet;

public sealed class Login : Packet
{
    private readonly ulong _handle;
    internal Login(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 1u;
    public int ClientProtocol { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public byte[] ConnectionRequest { get => PacketBridge.Bytes(_handle, 1); set => PacketBridge.SetBytes(_handle, 1, value); }
}

public sealed class PlayStatus : Packet
{
    private readonly ulong _handle;
    internal PlayStatus(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 2u;
    public int Status { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class ServerToClientHandshake : Packet
{
    private readonly ulong _handle;
    internal ServerToClientHandshake(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 3u;
    public byte[] JWT { get => PacketBridge.Bytes(_handle, 0); set => PacketBridge.SetBytes(_handle, 0, value); }
}

public sealed class ClientToServerHandshake : Packet
{
    private readonly ulong _handle;
    internal ClientToServerHandshake(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 4u;
}

public sealed class Disconnect : Packet
{
    private readonly ulong _handle;
    internal Disconnect(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 5u;
    public int Reason { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public bool HideDisconnectionScreen { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
    public string Message { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public string FilteredMessage { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
}

public sealed class ResourcePacksInfo : Packet
{
    private readonly ulong _handle;
    internal ResourcePacksInfo(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 6u;
    public bool TexturePackRequired { get => PacketBridge.Bool(_handle, 0); set => PacketBridge.SetBool(_handle, 0, value); }
    public bool HasAddons { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
    public bool HasScripts { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
    public bool ForceDisableVibrantVisuals { get => PacketBridge.Bool(_handle, 3); set => PacketBridge.SetBool(_handle, 3, value); }
    public Guid WorldTemplateUUID { get => PacketBridge.Guid(_handle, 4); set => PacketBridge.SetGuid(_handle, 4, value); }
    public string WorldTemplateVersion { get => PacketBridge.String(_handle, 5); set => PacketBridge.SetString(_handle, 5, value); }
    public Value TexturePacks => new(_handle, 6);
}

public sealed class ResourcePackStack : Packet
{
    private readonly ulong _handle;
    internal ResourcePackStack(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 7u;
    public bool TexturePackRequired { get => PacketBridge.Bool(_handle, 0); set => PacketBridge.SetBool(_handle, 0, value); }
    public Value TexturePacks => new(_handle, 1);
    public string BaseGameVersion { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public Value Experiments => new(_handle, 3);
    public bool ExperimentsPreviouslyToggled { get => PacketBridge.Bool(_handle, 4); set => PacketBridge.SetBool(_handle, 4, value); }
    public bool IncludeEditorPacks { get => PacketBridge.Bool(_handle, 5); set => PacketBridge.SetBool(_handle, 5, value); }
}

public sealed class ResourcePackClientResponse : Packet
{
    private readonly ulong _handle;
    internal ResourcePackClientResponse(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 8u;
    public byte Response { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value PacksToDownload => new(_handle, 1);
}

public sealed class Text : Packet
{
    private readonly ulong _handle;
    internal Text(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 9u;
    public byte TextType { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public bool NeedsTranslation { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
    public string SourceName { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public string Message { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public Value Parameters => new(_handle, 4);
    public string XUID { get => PacketBridge.String(_handle, 5); set => PacketBridge.SetString(_handle, 5, value); }
    public string PlatformChatID { get => PacketBridge.String(_handle, 6); set => PacketBridge.SetString(_handle, 6, value); }
    public Value FilteredMessage => new(_handle, 7);
}

public sealed class SetTime : Packet
{
    private readonly ulong _handle;
    internal SetTime(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 10u;
    public int Time { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class StartGame : Packet
{
    private readonly ulong _handle;
    internal StartGame(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 11u;
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public int PlayerGameMode { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public Vector3 PlayerPosition { get => PacketBridge.Vector3(_handle, 3); set => PacketBridge.SetVector3(_handle, 3, value); }
    public float Pitch { get => (float)PacketBridge.Number(_handle, 4); set => PacketBridge.SetNumber(_handle, 4, value); }
    public float Yaw { get => (float)PacketBridge.Number(_handle, 5); set => PacketBridge.SetNumber(_handle, 5, value); }
    public long WorldSeed { get => checked((long)PacketBridge.Signed(_handle, 6)); set => PacketBridge.SetSigned(_handle, 6, value); }
    public short SpawnBiomeType { get => checked((short)PacketBridge.Signed(_handle, 7)); set => PacketBridge.SetSigned(_handle, 7, value); }
    public string UserDefinedBiomeName { get => PacketBridge.String(_handle, 8); set => PacketBridge.SetString(_handle, 8, value); }
    public int Dimension { get => checked((int)PacketBridge.Signed(_handle, 9)); set => PacketBridge.SetSigned(_handle, 9, value); }
    public int Generator { get => checked((int)PacketBridge.Signed(_handle, 10)); set => PacketBridge.SetSigned(_handle, 10, value); }
    public int WorldGameMode { get => checked((int)PacketBridge.Signed(_handle, 11)); set => PacketBridge.SetSigned(_handle, 11, value); }
    public bool Hardcore { get => PacketBridge.Bool(_handle, 12); set => PacketBridge.SetBool(_handle, 12, value); }
    public int Difficulty { get => checked((int)PacketBridge.Signed(_handle, 13)); set => PacketBridge.SetSigned(_handle, 13, value); }
    public Value WorldSpawn => new(_handle, 14);
    public bool AchievementsDisabled { get => PacketBridge.Bool(_handle, 15); set => PacketBridge.SetBool(_handle, 15, value); }
    public int EditorWorldType { get => checked((int)PacketBridge.Signed(_handle, 16)); set => PacketBridge.SetSigned(_handle, 16, value); }
    public bool CreatedInEditor { get => PacketBridge.Bool(_handle, 17); set => PacketBridge.SetBool(_handle, 17, value); }
    public bool ExportedFromEditor { get => PacketBridge.Bool(_handle, 18); set => PacketBridge.SetBool(_handle, 18, value); }
    public int ServerEditorConnectionPolicy { get => checked((int)PacketBridge.Signed(_handle, 19)); set => PacketBridge.SetSigned(_handle, 19, value); }
    public bool AllowAnonymousBlockDropsInEditorWorlds { get => PacketBridge.Bool(_handle, 20); set => PacketBridge.SetBool(_handle, 20, value); }
    public int DayCycleLockTime { get => checked((int)PacketBridge.Signed(_handle, 21)); set => PacketBridge.SetSigned(_handle, 21, value); }
    public int EducationEditionOffer { get => checked((int)PacketBridge.Signed(_handle, 22)); set => PacketBridge.SetSigned(_handle, 22, value); }
    public bool EducationFeaturesEnabled { get => PacketBridge.Bool(_handle, 23); set => PacketBridge.SetBool(_handle, 23, value); }
    public string EducationProductID { get => PacketBridge.String(_handle, 24); set => PacketBridge.SetString(_handle, 24, value); }
    public float RainLevel { get => (float)PacketBridge.Number(_handle, 25); set => PacketBridge.SetNumber(_handle, 25, value); }
    public float LightningLevel { get => (float)PacketBridge.Number(_handle, 26); set => PacketBridge.SetNumber(_handle, 26, value); }
    public bool ConfirmedPlatformLockedContent { get => PacketBridge.Bool(_handle, 27); set => PacketBridge.SetBool(_handle, 27, value); }
    public bool MultiPlayerGame { get => PacketBridge.Bool(_handle, 28); set => PacketBridge.SetBool(_handle, 28, value); }
    public bool LANBroadcastEnabled { get => PacketBridge.Bool(_handle, 29); set => PacketBridge.SetBool(_handle, 29, value); }
    public int XBLBroadcastMode { get => checked((int)PacketBridge.Signed(_handle, 30)); set => PacketBridge.SetSigned(_handle, 30, value); }
    public int PlatformBroadcastMode { get => checked((int)PacketBridge.Signed(_handle, 31)); set => PacketBridge.SetSigned(_handle, 31, value); }
    public bool CommandsEnabled { get => PacketBridge.Bool(_handle, 32); set => PacketBridge.SetBool(_handle, 32, value); }
    public bool TexturePackRequired { get => PacketBridge.Bool(_handle, 33); set => PacketBridge.SetBool(_handle, 33, value); }
    public Value GameRules => new(_handle, 34);
    public Value Experiments => new(_handle, 35);
    public bool ExperimentsPreviouslyToggled { get => PacketBridge.Bool(_handle, 36); set => PacketBridge.SetBool(_handle, 36, value); }
    public bool BonusChestEnabled { get => PacketBridge.Bool(_handle, 37); set => PacketBridge.SetBool(_handle, 37, value); }
    public bool StartWithMapEnabled { get => PacketBridge.Bool(_handle, 38); set => PacketBridge.SetBool(_handle, 38, value); }
    public int PlayerPermissions { get => checked((int)PacketBridge.Signed(_handle, 39)); set => PacketBridge.SetSigned(_handle, 39, value); }
    public int ServerChunkTickRadius { get => checked((int)PacketBridge.Signed(_handle, 40)); set => PacketBridge.SetSigned(_handle, 40, value); }
    public bool HasLockedBehaviourPack { get => PacketBridge.Bool(_handle, 41); set => PacketBridge.SetBool(_handle, 41, value); }
    public bool HasLockedTexturePack { get => PacketBridge.Bool(_handle, 42); set => PacketBridge.SetBool(_handle, 42, value); }
    public bool FromLockedWorldTemplate { get => PacketBridge.Bool(_handle, 43); set => PacketBridge.SetBool(_handle, 43, value); }
    public bool MSAGamerTagsOnly { get => PacketBridge.Bool(_handle, 44); set => PacketBridge.SetBool(_handle, 44, value); }
    public bool FromWorldTemplate { get => PacketBridge.Bool(_handle, 45); set => PacketBridge.SetBool(_handle, 45, value); }
    public bool WorldTemplateSettingsLocked { get => PacketBridge.Bool(_handle, 46); set => PacketBridge.SetBool(_handle, 46, value); }
    public bool OnlySpawnV1Villagers { get => PacketBridge.Bool(_handle, 47); set => PacketBridge.SetBool(_handle, 47, value); }
    public bool PersonaDisabled { get => PacketBridge.Bool(_handle, 48); set => PacketBridge.SetBool(_handle, 48, value); }
    public bool CustomSkinsDisabled { get => PacketBridge.Bool(_handle, 49); set => PacketBridge.SetBool(_handle, 49, value); }
    public bool EmoteChatMuted { get => PacketBridge.Bool(_handle, 50); set => PacketBridge.SetBool(_handle, 50, value); }
    public string BaseGameVersion { get => PacketBridge.String(_handle, 51); set => PacketBridge.SetString(_handle, 51, value); }
    public int LimitedWorldWidth { get => checked((int)PacketBridge.Signed(_handle, 52)); set => PacketBridge.SetSigned(_handle, 52, value); }
    public int LimitedWorldDepth { get => checked((int)PacketBridge.Signed(_handle, 53)); set => PacketBridge.SetSigned(_handle, 53, value); }
    public bool NewNether { get => PacketBridge.Bool(_handle, 54); set => PacketBridge.SetBool(_handle, 54, value); }
    public Value EducationSharedResourceURI => new(_handle, 55);
    public Value ForceExperimentalGameplay => new(_handle, 56);
    public string LevelID { get => PacketBridge.String(_handle, 57); set => PacketBridge.SetString(_handle, 57, value); }
    public string WorldName { get => PacketBridge.String(_handle, 58); set => PacketBridge.SetString(_handle, 58, value); }
    public string TemplateContentIdentity { get => PacketBridge.String(_handle, 59); set => PacketBridge.SetString(_handle, 59, value); }
    public bool Trial { get => PacketBridge.Bool(_handle, 60); set => PacketBridge.SetBool(_handle, 60, value); }
    public Value PlayerMovementSettings => new(_handle, 61);
    public long Time { get => checked((long)PacketBridge.Signed(_handle, 62)); set => PacketBridge.SetSigned(_handle, 62, value); }
    public int EnchantmentSeed { get => checked((int)PacketBridge.Signed(_handle, 63)); set => PacketBridge.SetSigned(_handle, 63, value); }
    public Value Blocks => new(_handle, 64);
    public string MultiPlayerCorrelationID { get => PacketBridge.String(_handle, 65); set => PacketBridge.SetString(_handle, 65, value); }
    public bool ServerAuthoritativeInventory { get => PacketBridge.Bool(_handle, 66); set => PacketBridge.SetBool(_handle, 66, value); }
    public string GameVersion { get => PacketBridge.String(_handle, 67); set => PacketBridge.SetString(_handle, 67, value); }
    public Value PropertyData => new(_handle, 68);
    public ulong ServerBlockStateChecksum { get => checked((ulong)PacketBridge.Unsigned(_handle, 69)); set => PacketBridge.SetUnsigned(_handle, 69, value); }
    public bool ClientSideGeneration { get => PacketBridge.Bool(_handle, 70); set => PacketBridge.SetBool(_handle, 70, value); }
    public Guid WorldTemplateID { get => PacketBridge.Guid(_handle, 71); set => PacketBridge.SetGuid(_handle, 71, value); }
    public byte ChatRestrictionLevel { get => checked((byte)PacketBridge.Unsigned(_handle, 72)); set => PacketBridge.SetUnsigned(_handle, 72, value); }
    public bool DisablePlayerInteractions { get => PacketBridge.Bool(_handle, 73); set => PacketBridge.SetBool(_handle, 73, value); }
    public bool UseBlockNetworkIDHashes { get => PacketBridge.Bool(_handle, 74); set => PacketBridge.SetBool(_handle, 74, value); }
    public bool ServerAuthoritativeSound { get => PacketBridge.Bool(_handle, 75); set => PacketBridge.SetBool(_handle, 75, value); }
    public bool IsLoggingChat { get => PacketBridge.Bool(_handle, 76); set => PacketBridge.SetBool(_handle, 76, value); }
    public Value ServerJoinInformation => new(_handle, 77);
    public string ServerID { get => PacketBridge.String(_handle, 78); set => PacketBridge.SetString(_handle, 78, value); }
    public string ScenarioID { get => PacketBridge.String(_handle, 79); set => PacketBridge.SetString(_handle, 79, value); }
    public string WorldID { get => PacketBridge.String(_handle, 80); set => PacketBridge.SetString(_handle, 80, value); }
    public string OwnerID { get => PacketBridge.String(_handle, 81); set => PacketBridge.SetString(_handle, 81, value); }
}

public sealed class AddPlayer : Packet
{
    private readonly ulong _handle;
    internal AddPlayer(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 12u;
    public Guid UUID { get => PacketBridge.Guid(_handle, 0); set => PacketBridge.SetGuid(_handle, 0, value); }
    public string Username { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public string PlatformChatID { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 4); set => PacketBridge.SetVector3(_handle, 4, value); }
    public Vector3 Velocity { get => PacketBridge.Vector3(_handle, 5); set => PacketBridge.SetVector3(_handle, 5, value); }
    public float Pitch { get => (float)PacketBridge.Number(_handle, 6); set => PacketBridge.SetNumber(_handle, 6, value); }
    public float Yaw { get => (float)PacketBridge.Number(_handle, 7); set => PacketBridge.SetNumber(_handle, 7, value); }
    public float HeadYaw { get => (float)PacketBridge.Number(_handle, 8); set => PacketBridge.SetNumber(_handle, 8, value); }
    public Value HeldItem => new(_handle, 9);
    public int GameType { get => checked((int)PacketBridge.Signed(_handle, 10)); set => PacketBridge.SetSigned(_handle, 10, value); }
    public Value EntityMetadata => new(_handle, 11);
    public Value EntityProperties => new(_handle, 12);
    public Value AbilityData => new(_handle, 13);
    public Value EntityLinks => new(_handle, 14);
    public string DeviceID { get => PacketBridge.String(_handle, 15); set => PacketBridge.SetString(_handle, 15, value); }
    public int BuildPlatform { get => checked((int)PacketBridge.Signed(_handle, 16)); set => PacketBridge.SetSigned(_handle, 16, value); }
}

public sealed class AddActor : Packet
{
    private readonly ulong _handle;
    internal AddActor(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 13u;
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public string EntityType { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 3); set => PacketBridge.SetVector3(_handle, 3, value); }
    public Vector3 Velocity { get => PacketBridge.Vector3(_handle, 4); set => PacketBridge.SetVector3(_handle, 4, value); }
    public float Pitch { get => (float)PacketBridge.Number(_handle, 5); set => PacketBridge.SetNumber(_handle, 5, value); }
    public float Yaw { get => (float)PacketBridge.Number(_handle, 6); set => PacketBridge.SetNumber(_handle, 6, value); }
    public float HeadYaw { get => (float)PacketBridge.Number(_handle, 7); set => PacketBridge.SetNumber(_handle, 7, value); }
    public float BodyYaw { get => (float)PacketBridge.Number(_handle, 8); set => PacketBridge.SetNumber(_handle, 8, value); }
    public Value Attributes => new(_handle, 9);
    public Value EntityMetadata => new(_handle, 10);
    public Value EntityProperties => new(_handle, 11);
    public Value EntityLinks => new(_handle, 12);
}

public sealed class RemoveActor : Packet
{
    private readonly ulong _handle;
    internal RemoveActor(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 14u;
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class AddItemActor : Packet
{
    private readonly ulong _handle;
    internal AddItemActor(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 15u;
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public Value Item => new(_handle, 2);
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 3); set => PacketBridge.SetVector3(_handle, 3, value); }
    public Vector3 Velocity { get => PacketBridge.Vector3(_handle, 4); set => PacketBridge.SetVector3(_handle, 4, value); }
    public Value EntityMetadata => new(_handle, 5);
    public bool FromFishing { get => PacketBridge.Bool(_handle, 6); set => PacketBridge.SetBool(_handle, 6, value); }
}

public sealed class TakeItemActor : Packet
{
    private readonly ulong _handle;
    internal TakeItemActor(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 17u;
    public ulong ItemEntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public ulong TakerEntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
}

public sealed class MoveActorAbsolute : Packet
{
    private readonly ulong _handle;
    internal MoveActorAbsolute(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 18u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte Flags { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 2); set => PacketBridge.SetVector3(_handle, 2, value); }
    public Vector3 Rotation { get => PacketBridge.Vector3(_handle, 3); set => PacketBridge.SetVector3(_handle, 3, value); }
}

public sealed class MovePlayer : Packet
{
    private readonly ulong _handle;
    internal MovePlayer(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 19u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 1); set => PacketBridge.SetVector3(_handle, 1, value); }
    public float Pitch { get => (float)PacketBridge.Number(_handle, 2); set => PacketBridge.SetNumber(_handle, 2, value); }
    public float Yaw { get => (float)PacketBridge.Number(_handle, 3); set => PacketBridge.SetNumber(_handle, 3, value); }
    public float HeadYaw { get => (float)PacketBridge.Number(_handle, 4); set => PacketBridge.SetNumber(_handle, 4, value); }
    public byte Mode { get => checked((byte)PacketBridge.Unsigned(_handle, 5)); set => PacketBridge.SetUnsigned(_handle, 5, value); }
    public bool OnGround { get => PacketBridge.Bool(_handle, 6); set => PacketBridge.SetBool(_handle, 6, value); }
    public ulong RiddenEntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 7)); set => PacketBridge.SetUnsigned(_handle, 7, value); }
    public int TeleportCause { get => checked((int)PacketBridge.Signed(_handle, 8)); set => PacketBridge.SetSigned(_handle, 8, value); }
    public int TeleportSourceEntityType { get => checked((int)PacketBridge.Signed(_handle, 9)); set => PacketBridge.SetSigned(_handle, 9, value); }
    public ulong Tick { get => checked((ulong)PacketBridge.Unsigned(_handle, 10)); set => PacketBridge.SetUnsigned(_handle, 10, value); }
}

public sealed class UpdateBlock : Packet
{
    private readonly ulong _handle;
    internal UpdateBlock(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 21u;
    public Value Position => new(_handle, 0);
    public uint NewBlockRuntimeID { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public uint Flags { get => checked((uint)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public uint Layer { get => checked((uint)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
}

public sealed class AddPainting : Packet
{
    private readonly ulong _handle;
    internal AddPainting(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 22u;
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 2); set => PacketBridge.SetVector3(_handle, 2, value); }
    public int Direction { get => checked((int)PacketBridge.Signed(_handle, 3)); set => PacketBridge.SetSigned(_handle, 3, value); }
    public string Title { get => PacketBridge.String(_handle, 4); set => PacketBridge.SetString(_handle, 4, value); }
}

public sealed class LevelEvent : Packet
{
    private readonly ulong _handle;
    internal LevelEvent(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 25u;
    public int EventType { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 1); set => PacketBridge.SetVector3(_handle, 1, value); }
    public int EventData { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
}

public sealed class BlockEvent : Packet
{
    private readonly ulong _handle;
    internal BlockEvent(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 26u;
    public Value Position => new(_handle, 0);
    public int EventType { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public int EventData { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
}

public sealed class ActorEvent : Packet
{
    private readonly ulong _handle;
    internal ActorEvent(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 27u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte EventType { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public int EventData { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public Value FireAtPosition => new(_handle, 3);
}

public sealed class MobEffect : Packet
{
    private readonly ulong _handle;
    internal MobEffect(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 28u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte Operation { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public int EffectType { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public int Amplifier { get => checked((int)PacketBridge.Signed(_handle, 3)); set => PacketBridge.SetSigned(_handle, 3, value); }
    public bool Particles { get => PacketBridge.Bool(_handle, 4); set => PacketBridge.SetBool(_handle, 4, value); }
    public int Duration { get => checked((int)PacketBridge.Signed(_handle, 5)); set => PacketBridge.SetSigned(_handle, 5, value); }
    public ulong Tick { get => checked((ulong)PacketBridge.Unsigned(_handle, 6)); set => PacketBridge.SetUnsigned(_handle, 6, value); }
    public bool Ambient { get => PacketBridge.Bool(_handle, 7); set => PacketBridge.SetBool(_handle, 7, value); }
}

public sealed class UpdateAttributes : Packet
{
    private readonly ulong _handle;
    internal UpdateAttributes(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 29u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value Attributes => new(_handle, 1);
    public ulong Tick { get => checked((ulong)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
}

public sealed class InventoryTransaction : Packet
{
    private readonly ulong _handle;
    internal InventoryTransaction(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 30u;
    public int LegacyRequestID { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public Value LegacySetItemSlots => new(_handle, 1);
    public Value Actions => new(_handle, 2);
    public Value TransactionData => new(_handle, 3);
}

public sealed class MobEquipment : Packet
{
    private readonly ulong _handle;
    internal MobEquipment(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 31u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value NewItem => new(_handle, 1);
    public byte InventorySlot { get => checked((byte)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public byte HotBarSlot { get => checked((byte)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
    public byte WindowID { get => checked((byte)PacketBridge.Unsigned(_handle, 4)); set => PacketBridge.SetUnsigned(_handle, 4, value); }
}

public sealed class MobArmourEquipment : Packet
{
    private readonly ulong _handle;
    internal MobArmourEquipment(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 32u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value Helmet => new(_handle, 1);
    public Value Chestplate => new(_handle, 2);
    public Value Leggings => new(_handle, 3);
    public Value Boots => new(_handle, 4);
    public Value Body => new(_handle, 5);
}

public sealed class Interact : Packet
{
    private readonly ulong _handle;
    internal Interact(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 33u;
    public byte ActionType { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public ulong TargetEntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public Value Position => new(_handle, 2);
}

public sealed class BlockPickRequest : Packet
{
    private readonly ulong _handle;
    internal BlockPickRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 34u;
    public Value Position => new(_handle, 0);
    public bool AddBlockNBT { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
    public byte HotBarSlot { get => checked((byte)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
}

public sealed class ActorPickRequest : Packet
{
    private readonly ulong _handle;
    internal ActorPickRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 35u;
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public byte HotBarSlot { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public bool WithData { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
}

public sealed class PlayerAction : Packet
{
    private readonly ulong _handle;
    internal PlayerAction(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 36u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public int ActionType { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public Value BlockPosition => new(_handle, 2);
    public Value ResultPosition => new(_handle, 3);
    public int BlockFace { get => checked((int)PacketBridge.Signed(_handle, 4)); set => PacketBridge.SetSigned(_handle, 4, value); }
}

public sealed class HurtArmour : Packet
{
    private readonly ulong _handle;
    internal HurtArmour(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 38u;
    public int Cause { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public int Damage { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public long ArmourSlots { get => checked((long)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
}

public sealed class SetActorData : Packet
{
    private readonly ulong _handle;
    internal SetActorData(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 39u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value EntityMetadata => new(_handle, 1);
    public Value EntityProperties => new(_handle, 2);
    public ulong Tick { get => checked((ulong)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
}

public sealed class SetActorMotion : Packet
{
    private readonly ulong _handle;
    internal SetActorMotion(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 40u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Vector3 Velocity { get => PacketBridge.Vector3(_handle, 1); set => PacketBridge.SetVector3(_handle, 1, value); }
    public ulong Tick { get => checked((ulong)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
}

public sealed class SetActorLink : Packet
{
    private readonly ulong _handle;
    internal SetActorLink(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 41u;
    public Value EntityLink => new(_handle, 0);
}

public sealed class SetHealth : Packet
{
    private readonly ulong _handle;
    internal SetHealth(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 42u;
    public int Health { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class SetSpawnPosition : Packet
{
    private readonly ulong _handle;
    internal SetSpawnPosition(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 43u;
    public int SpawnType { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public Value Position => new(_handle, 1);
    public int Dimension { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public Value SpawnPosition => new(_handle, 3);
}

public sealed class Animate : Packet
{
    private readonly ulong _handle;
    internal Animate(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 44u;
    public byte ActionType { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public float Data { get => (float)PacketBridge.Number(_handle, 2); set => PacketBridge.SetNumber(_handle, 2, value); }
    public byte SwingSource { get => checked((byte)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
}

public sealed class Respawn : Packet
{
    private readonly ulong _handle;
    internal Respawn(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 45u;
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 0); set => PacketBridge.SetVector3(_handle, 0, value); }
    public byte State { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
}

public sealed class ContainerOpen : Packet
{
    private readonly ulong _handle;
    internal ContainerOpen(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 46u;
    public byte WindowID { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte ContainerType { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public Value ContainerPosition => new(_handle, 2);
    public long ContainerEntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 3)); set => PacketBridge.SetSigned(_handle, 3, value); }
}

public sealed class ContainerClose : Packet
{
    private readonly ulong _handle;
    internal ContainerClose(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 47u;
    public byte WindowID { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte ContainerType { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public bool ServerSide { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
}

public sealed class PlayerHotBar : Packet
{
    private readonly ulong _handle;
    internal PlayerHotBar(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 48u;
    public uint SelectedHotBarSlot { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte WindowID { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public bool SelectHotBarSlot { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
}

public sealed class InventoryContent : Packet
{
    private readonly ulong _handle;
    internal InventoryContent(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 49u;
    public uint WindowID { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value Content => new(_handle, 1);
    public Value Container => new(_handle, 2);
    public Value StorageItem => new(_handle, 3);
}

public sealed class InventorySlot : Packet
{
    private readonly ulong _handle;
    internal InventorySlot(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 50u;
    public uint WindowID { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public uint Slot { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public Value Container => new(_handle, 2);
    public Value StorageItem => new(_handle, 3);
    public Value NewItem => new(_handle, 4);
}

public sealed class ContainerSetData : Packet
{
    private readonly ulong _handle;
    internal ContainerSetData(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 51u;
    public byte WindowID { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public int Key { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public int Value { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
}

public sealed class CraftingData : Packet
{
    private readonly ulong _handle;
    internal CraftingData(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 52u;
    public Value Recipes => new(_handle, 0);
    public Value PotionRecipes => new(_handle, 1);
    public Value PotionContainerChangeRecipes => new(_handle, 2);
    public Value MaterialReducers => new(_handle, 3);
    public bool ClearRecipes { get => PacketBridge.Bool(_handle, 4); set => PacketBridge.SetBool(_handle, 4, value); }
}

public sealed class GUIDataPickItem : Packet
{
    private readonly ulong _handle;
    internal GUIDataPickItem(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 54u;
    public string ItemName { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public string ItemEffects { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public int HotBarSlot { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
}

public sealed class AdventureSettings : Packet
{
    private readonly ulong _handle;
    internal AdventureSettings(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 55u;
    public uint Flags { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public uint CommandPermissionLevel { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public uint ActionPermissions { get => checked((uint)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public uint PermissionLevel { get => checked((uint)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
    public uint CustomStoredPermissions { get => checked((uint)PacketBridge.Unsigned(_handle, 4)); set => PacketBridge.SetUnsigned(_handle, 4, value); }
    public long PlayerUniqueID { get => checked((long)PacketBridge.Signed(_handle, 5)); set => PacketBridge.SetSigned(_handle, 5, value); }
}

public sealed class BlockActorData : Packet
{
    private readonly ulong _handle;
    internal BlockActorData(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 56u;
    public Value Position => new(_handle, 0);
    public Value NBTData => new(_handle, 1);
}

public sealed class LevelChunk : Packet
{
    private readonly ulong _handle;
    internal LevelChunk(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 58u;
    public Value Position => new(_handle, 0);
    public int Dimension { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public ushort HighestSubChunk { get => checked((ushort)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public uint SubChunkCount { get => checked((uint)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
    public bool CacheEnabled { get => PacketBridge.Bool(_handle, 4); set => PacketBridge.SetBool(_handle, 4, value); }
    public Value BlobHashes => new(_handle, 5);
    public byte[] RawPayload { get => PacketBridge.Bytes(_handle, 6); set => PacketBridge.SetBytes(_handle, 6, value); }
}

public sealed class SetCommandsEnabled : Packet
{
    private readonly ulong _handle;
    internal SetCommandsEnabled(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 59u;
    public bool Enabled { get => PacketBridge.Bool(_handle, 0); set => PacketBridge.SetBool(_handle, 0, value); }
}

public sealed class SetDifficulty : Packet
{
    private readonly ulong _handle;
    internal SetDifficulty(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 60u;
    public uint Difficulty { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
}

public sealed class ChangeDimension : Packet
{
    private readonly ulong _handle;
    internal ChangeDimension(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 61u;
    public int Dimension { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 1); set => PacketBridge.SetVector3(_handle, 1, value); }
    public bool Respawn { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
    public Value LoadingScreenID => new(_handle, 3);
}

public sealed class SetPlayerGameType : Packet
{
    private readonly ulong _handle;
    internal SetPlayerGameType(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 62u;
    public int GameType { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class PlayerList : Packet
{
    private readonly ulong _handle;
    internal PlayerList(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 63u;
    public byte ActionType { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value Entries => new(_handle, 1);
}

public sealed class SimpleEvent : Packet
{
    private readonly ulong _handle;
    internal SimpleEvent(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 64u;
    public ushort EventType { get => checked((ushort)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
}

public sealed class Event : Packet
{
    private readonly ulong _handle;
    internal Event(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 65u;
    public long EntityRuntimeID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public bool UsePlayerID { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
    public Value EventValue => new(_handle, 2);
}

public sealed class SpawnExperienceOrb : Packet
{
    private readonly ulong _handle;
    internal SpawnExperienceOrb(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 66u;
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 0); set => PacketBridge.SetVector3(_handle, 0, value); }
    public int ExperienceAmount { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
}

public sealed class ClientBoundMapItemData : Packet
{
    private readonly ulong _handle;
    internal ClientBoundMapItemData(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 67u;
    public long MapID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public uint UpdateFlags { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public byte Dimension { get => checked((byte)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public bool LockedMap { get => PacketBridge.Bool(_handle, 3); set => PacketBridge.SetBool(_handle, 3, value); }
    public Value Origin => new(_handle, 4);
    public byte Scale { get => checked((byte)PacketBridge.Unsigned(_handle, 5)); set => PacketBridge.SetUnsigned(_handle, 5, value); }
    public Value MapsIncludedIn => new(_handle, 6);
    public Value TrackedObjects => new(_handle, 7);
    public Value Decorations => new(_handle, 8);
    public int Height { get => checked((int)PacketBridge.Signed(_handle, 9)); set => PacketBridge.SetSigned(_handle, 9, value); }
    public int Width { get => checked((int)PacketBridge.Signed(_handle, 10)); set => PacketBridge.SetSigned(_handle, 10, value); }
    public int XOffset { get => checked((int)PacketBridge.Signed(_handle, 11)); set => PacketBridge.SetSigned(_handle, 11, value); }
    public int YOffset { get => checked((int)PacketBridge.Signed(_handle, 12)); set => PacketBridge.SetSigned(_handle, 12, value); }
    public Value Pixels => new(_handle, 13);
}

public sealed class MapInfoRequest : Packet
{
    private readonly ulong _handle;
    internal MapInfoRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 68u;
    public long MapID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public Value ClientPixels => new(_handle, 1);
}

public sealed class RequestChunkRadius : Packet
{
    private readonly ulong _handle;
    internal RequestChunkRadius(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 69u;
    public int ChunkRadius { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public byte MaxChunkRadius { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
}

public sealed class ChunkRadiusUpdated : Packet
{
    private readonly ulong _handle;
    internal ChunkRadiusUpdated(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 70u;
    public int ChunkRadius { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class GameRulesChanged : Packet
{
    private readonly ulong _handle;
    internal GameRulesChanged(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 72u;
    public Value GameRules => new(_handle, 0);
}

public sealed class Camera : Packet
{
    private readonly ulong _handle;
    internal Camera(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 73u;
    public long CameraEntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public long TargetPlayerUniqueID { get => checked((long)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
}

public sealed class BossEvent : Packet
{
    private readonly ulong _handle;
    internal BossEvent(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 74u;
    public long BossEntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public long PlayerUniqueID { get => checked((long)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public byte EventType { get => checked((byte)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public string BossBarTitle { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public string FilteredBossBarTitle { get => PacketBridge.String(_handle, 4); set => PacketBridge.SetString(_handle, 4, value); }
    public float HealthPercentage { get => (float)PacketBridge.Number(_handle, 5); set => PacketBridge.SetNumber(_handle, 5, value); }
    public byte Colour { get => checked((byte)PacketBridge.Unsigned(_handle, 6)); set => PacketBridge.SetUnsigned(_handle, 6, value); }
    public byte Overlay { get => checked((byte)PacketBridge.Unsigned(_handle, 7)); set => PacketBridge.SetUnsigned(_handle, 7, value); }
}

public sealed class ShowCredits : Packet
{
    private readonly ulong _handle;
    internal ShowCredits(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 75u;
    public ulong PlayerRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public int StatusType { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
}

public sealed class AvailableCommands : Packet
{
    private readonly ulong _handle;
    internal AvailableCommands(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 76u;
    public Value EnumValues => new(_handle, 0);
    public Value ChainedSubcommandValues => new(_handle, 1);
    public Value Suffixes => new(_handle, 2);
    public Value Enums => new(_handle, 3);
    public Value ChainedSubcommands => new(_handle, 4);
    public Value Commands => new(_handle, 5);
    public Value DynamicEnums => new(_handle, 6);
    public Value Constraints => new(_handle, 7);
}

public sealed class CommandRequest : Packet
{
    private readonly ulong _handle;
    internal CommandRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 77u;
    public string CommandLine { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public Value CommandOrigin => new(_handle, 1);
    public bool Internal { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
    public string Version { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
}

public sealed class CommandBlockUpdate : Packet
{
    private readonly ulong _handle;
    internal CommandBlockUpdate(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 78u;
    public bool Block { get => PacketBridge.Bool(_handle, 0); set => PacketBridge.SetBool(_handle, 0, value); }
    public Value Position => new(_handle, 1);
    public uint Mode { get => checked((uint)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public bool NeedsRedstone { get => PacketBridge.Bool(_handle, 3); set => PacketBridge.SetBool(_handle, 3, value); }
    public bool Conditional { get => PacketBridge.Bool(_handle, 4); set => PacketBridge.SetBool(_handle, 4, value); }
    public ulong MinecartEntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 5)); set => PacketBridge.SetUnsigned(_handle, 5, value); }
    public string Command { get => PacketBridge.String(_handle, 6); set => PacketBridge.SetString(_handle, 6, value); }
    public string LastOutput { get => PacketBridge.String(_handle, 7); set => PacketBridge.SetString(_handle, 7, value); }
    public string Name { get => PacketBridge.String(_handle, 8); set => PacketBridge.SetString(_handle, 8, value); }
    public string FilteredName { get => PacketBridge.String(_handle, 9); set => PacketBridge.SetString(_handle, 9, value); }
    public bool ShouldTrackOutput { get => PacketBridge.Bool(_handle, 10); set => PacketBridge.SetBool(_handle, 10, value); }
    public uint TickDelay { get => checked((uint)PacketBridge.Unsigned(_handle, 11)); set => PacketBridge.SetUnsigned(_handle, 11, value); }
    public bool ExecuteOnFirstTick { get => PacketBridge.Bool(_handle, 12); set => PacketBridge.SetBool(_handle, 12, value); }
}

public sealed class CommandOutput : Packet
{
    private readonly ulong _handle;
    internal CommandOutput(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 79u;
    public Value CommandOrigin => new(_handle, 0);
    public byte OutputType { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public uint SuccessCount { get => checked((uint)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public Value OutputMessages => new(_handle, 3);
    public Value DataSet => new(_handle, 4);
}

public sealed class UpdateTrade : Packet
{
    private readonly ulong _handle;
    internal UpdateTrade(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 80u;
    public byte WindowID { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte WindowType { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public int Size { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public int TradeTier { get => checked((int)PacketBridge.Signed(_handle, 3)); set => PacketBridge.SetSigned(_handle, 3, value); }
    public long VillagerUniqueID { get => checked((long)PacketBridge.Signed(_handle, 4)); set => PacketBridge.SetSigned(_handle, 4, value); }
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 5)); set => PacketBridge.SetSigned(_handle, 5, value); }
    public string DisplayName { get => PacketBridge.String(_handle, 6); set => PacketBridge.SetString(_handle, 6, value); }
    public bool NewTradeUI { get => PacketBridge.Bool(_handle, 7); set => PacketBridge.SetBool(_handle, 7, value); }
    public bool DemandBasedPrices { get => PacketBridge.Bool(_handle, 8); set => PacketBridge.SetBool(_handle, 8, value); }
    public byte[] SerialisedOffers { get => PacketBridge.Bytes(_handle, 9); set => PacketBridge.SetBytes(_handle, 9, value); }
}

public sealed class UpdateEquip : Packet
{
    private readonly ulong _handle;
    internal UpdateEquip(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 81u;
    public byte WindowID { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte WindowType { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public int Size { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 3)); set => PacketBridge.SetSigned(_handle, 3, value); }
    public byte[] SerialisedInventoryData { get => PacketBridge.Bytes(_handle, 4); set => PacketBridge.SetBytes(_handle, 4, value); }
}

public sealed class ResourcePackDataInfo : Packet
{
    private readonly ulong _handle;
    internal ResourcePackDataInfo(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 82u;
    public string UUID { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public uint DataChunkSize { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public uint ChunkCount { get => checked((uint)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public ulong Size { get => checked((ulong)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
    public byte[] Hash { get => PacketBridge.Bytes(_handle, 4); set => PacketBridge.SetBytes(_handle, 4, value); }
    public bool Premium { get => PacketBridge.Bool(_handle, 5); set => PacketBridge.SetBool(_handle, 5, value); }
    public byte PackType { get => checked((byte)PacketBridge.Unsigned(_handle, 6)); set => PacketBridge.SetUnsigned(_handle, 6, value); }
}

public sealed class ResourcePackChunkData : Packet
{
    private readonly ulong _handle;
    internal ResourcePackChunkData(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 83u;
    public string UUID { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public uint ChunkIndex { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public ulong DataOffset { get => checked((ulong)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public byte[] Data { get => PacketBridge.Bytes(_handle, 3); set => PacketBridge.SetBytes(_handle, 3, value); }
}

public sealed class ResourcePackChunkRequest : Packet
{
    private readonly ulong _handle;
    internal ResourcePackChunkRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 84u;
    public string UUID { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public int ChunkIndex { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
}

public sealed class Transfer : Packet
{
    private readonly ulong _handle;
    internal Transfer(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 85u;
    public string Address { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public ushort Port { get => checked((ushort)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public bool ReloadWorld { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
}

public sealed class PlaySound : Packet
{
    private readonly ulong _handle;
    internal PlaySound(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 86u;
    public string SoundName { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 1); set => PacketBridge.SetVector3(_handle, 1, value); }
    public float Volume { get => (float)PacketBridge.Number(_handle, 2); set => PacketBridge.SetNumber(_handle, 2, value); }
    public float Pitch { get => (float)PacketBridge.Number(_handle, 3); set => PacketBridge.SetNumber(_handle, 3, value); }
    public Value Handle => new(_handle, 4);
}

public sealed class StopSound : Packet
{
    private readonly ulong _handle;
    internal StopSound(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 87u;
    public string SoundName { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public bool StopAll { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
    public bool StopMusicLegacy { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
}

public sealed class SetTitle : Packet
{
    private readonly ulong _handle;
    internal SetTitle(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 88u;
    public int ActionType { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public string Text { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public int FadeInDuration { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public int RemainDuration { get => checked((int)PacketBridge.Signed(_handle, 3)); set => PacketBridge.SetSigned(_handle, 3, value); }
    public int FadeOutDuration { get => checked((int)PacketBridge.Signed(_handle, 4)); set => PacketBridge.SetSigned(_handle, 4, value); }
    public string XUID { get => PacketBridge.String(_handle, 5); set => PacketBridge.SetString(_handle, 5, value); }
    public string PlatformOnlineID { get => PacketBridge.String(_handle, 6); set => PacketBridge.SetString(_handle, 6, value); }
    public string FilteredMessage { get => PacketBridge.String(_handle, 7); set => PacketBridge.SetString(_handle, 7, value); }
}

public sealed class AddBehaviourTree : Packet
{
    private readonly ulong _handle;
    internal AddBehaviourTree(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 89u;
    public string BehaviourTree { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
}

public sealed class StructureBlockUpdate : Packet
{
    private readonly ulong _handle;
    internal StructureBlockUpdate(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 90u;
    public Value Position => new(_handle, 0);
    public string StructureName { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public string FilteredStructureName { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public string DataField { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public bool IncludePlayers { get => PacketBridge.Bool(_handle, 4); set => PacketBridge.SetBool(_handle, 4, value); }
    public bool ShowBoundingBox { get => PacketBridge.Bool(_handle, 5); set => PacketBridge.SetBool(_handle, 5, value); }
    public int StructureBlockType { get => checked((int)PacketBridge.Signed(_handle, 6)); set => PacketBridge.SetSigned(_handle, 6, value); }
    public Value Settings => new(_handle, 7);
    public int RedstoneSaveMode { get => checked((int)PacketBridge.Signed(_handle, 8)); set => PacketBridge.SetSigned(_handle, 8, value); }
    public bool ShouldTrigger { get => PacketBridge.Bool(_handle, 9); set => PacketBridge.SetBool(_handle, 9, value); }
    public bool Waterlogged { get => PacketBridge.Bool(_handle, 10); set => PacketBridge.SetBool(_handle, 10, value); }
}

public sealed class ShowStoreOffer : Packet
{
    private readonly ulong _handle;
    internal ShowStoreOffer(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 91u;
    public Guid OfferID { get => PacketBridge.Guid(_handle, 0); set => PacketBridge.SetGuid(_handle, 0, value); }
    public byte Type { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
}

public sealed class PurchaseReceipt : Packet
{
    private readonly ulong _handle;
    internal PurchaseReceipt(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 92u;
    public Value Receipts => new(_handle, 0);
}

public sealed class PlayerSkin : Packet
{
    private readonly ulong _handle;
    internal PlayerSkin(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 93u;
    public Guid UUID { get => PacketBridge.Guid(_handle, 0); set => PacketBridge.SetGuid(_handle, 0, value); }
    public Value Skin => new(_handle, 1);
    public string NewSkinName { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public string OldSkinName { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
}

public sealed class SubClientLogin : Packet
{
    private readonly ulong _handle;
    internal SubClientLogin(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 94u;
    public byte[] ConnectionRequest { get => PacketBridge.Bytes(_handle, 0); set => PacketBridge.SetBytes(_handle, 0, value); }
}

public sealed class AutomationClientConnect : Packet
{
    private readonly ulong _handle;
    internal AutomationClientConnect(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 95u;
    public string ServerURI { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
}

public sealed class SetLastHurtBy : Packet
{
    private readonly ulong _handle;
    internal SetLastHurtBy(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 96u;
    public int EntityType { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class BookEdit : Packet
{
    private readonly ulong _handle;
    internal BookEdit(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 97u;
    public int InventorySlot { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public uint ActionType { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public int PageNumber { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public int SecondaryPageNumber { get => checked((int)PacketBridge.Signed(_handle, 3)); set => PacketBridge.SetSigned(_handle, 3, value); }
    public string Text { get => PacketBridge.String(_handle, 4); set => PacketBridge.SetString(_handle, 4, value); }
    public string PhotoName { get => PacketBridge.String(_handle, 5); set => PacketBridge.SetString(_handle, 5, value); }
    public string Title { get => PacketBridge.String(_handle, 6); set => PacketBridge.SetString(_handle, 6, value); }
    public string Author { get => PacketBridge.String(_handle, 7); set => PacketBridge.SetString(_handle, 7, value); }
    public string XUID { get => PacketBridge.String(_handle, 8); set => PacketBridge.SetString(_handle, 8, value); }
}

public sealed class NPCRequest : Packet
{
    private readonly ulong _handle;
    internal NPCRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 98u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte RequestType { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public string CommandString { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public byte ActionType { get => checked((byte)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
    public string SceneName { get => PacketBridge.String(_handle, 4); set => PacketBridge.SetString(_handle, 4, value); }
}

public sealed class PhotoTransfer : Packet
{
    private readonly ulong _handle;
    internal PhotoTransfer(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 99u;
    public string PhotoName { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public byte[] PhotoData { get => PacketBridge.Bytes(_handle, 1); set => PacketBridge.SetBytes(_handle, 1, value); }
    public string BookID { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public byte PhotoType { get => checked((byte)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
    public byte SourceType { get => checked((byte)PacketBridge.Unsigned(_handle, 4)); set => PacketBridge.SetUnsigned(_handle, 4, value); }
    public long OwnerEntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 5)); set => PacketBridge.SetSigned(_handle, 5, value); }
    public string NewPhotoName { get => PacketBridge.String(_handle, 6); set => PacketBridge.SetString(_handle, 6, value); }
}

public sealed class ModalFormRequest : Packet
{
    private readonly ulong _handle;
    internal ModalFormRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 100u;
    public uint FormID { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte[] FormData { get => PacketBridge.Bytes(_handle, 1); set => PacketBridge.SetBytes(_handle, 1, value); }
}

public sealed class ModalFormResponse : Packet
{
    private readonly ulong _handle;
    internal ModalFormResponse(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 101u;
    public uint FormID { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value ResponseData => new(_handle, 1);
    public Value CancelReason => new(_handle, 2);
}

public sealed class ServerSettingsRequest : Packet
{
    private readonly ulong _handle;
    internal ServerSettingsRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 102u;
}

public sealed class ServerSettingsResponse : Packet
{
    private readonly ulong _handle;
    internal ServerSettingsResponse(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 103u;
    public uint FormID { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte[] FormData { get => PacketBridge.Bytes(_handle, 1); set => PacketBridge.SetBytes(_handle, 1, value); }
}

public sealed class ShowProfile : Packet
{
    private readonly ulong _handle;
    internal ShowProfile(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 104u;
    public string XUID { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
}

public sealed class SetDefaultGameType : Packet
{
    private readonly ulong _handle;
    internal SetDefaultGameType(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 105u;
    public int GameType { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class RemoveObjective : Packet
{
    private readonly ulong _handle;
    internal RemoveObjective(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 106u;
    public string ObjectiveName { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
}

public sealed class SetDisplayObjective : Packet
{
    private readonly ulong _handle;
    internal SetDisplayObjective(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 107u;
    public string DisplaySlot { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public string ObjectiveName { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public string DisplayName { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public string CriteriaName { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public int SortOrder { get => checked((int)PacketBridge.Signed(_handle, 4)); set => PacketBridge.SetSigned(_handle, 4, value); }
}

public sealed class SetScore : Packet
{
    private readonly ulong _handle;
    internal SetScore(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 108u;
    public byte ActionType { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value Entries => new(_handle, 1);
}

public sealed class LabTable : Packet
{
    private readonly ulong _handle;
    internal LabTable(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 109u;
    public byte ActionType { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value Position => new(_handle, 1);
    public byte ReactionType { get => checked((byte)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
}

public sealed class UpdateBlockSynced : Packet
{
    private readonly ulong _handle;
    internal UpdateBlockSynced(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 110u;
    public Value Position => new(_handle, 0);
    public uint NewBlockRuntimeID { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public uint Flags { get => checked((uint)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public uint Layer { get => checked((uint)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
    public ulong EntityUniqueID { get => checked((ulong)PacketBridge.Unsigned(_handle, 4)); set => PacketBridge.SetUnsigned(_handle, 4, value); }
    public ulong TransitionType { get => checked((ulong)PacketBridge.Unsigned(_handle, 5)); set => PacketBridge.SetUnsigned(_handle, 5, value); }
}

public sealed class MoveActorDelta : Packet
{
    private readonly ulong _handle;
    internal MoveActorDelta(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 111u;
    public ushort Flags { get => checked((ushort)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 2); set => PacketBridge.SetVector3(_handle, 2, value); }
    public Vector3 Rotation { get => PacketBridge.Vector3(_handle, 3); set => PacketBridge.SetVector3(_handle, 3, value); }
}

public sealed class SetScoreboardIdentity : Packet
{
    private readonly ulong _handle;
    internal SetScoreboardIdentity(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 112u;
    public byte ActionType { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value Entries => new(_handle, 1);
}

public sealed class SetLocalPlayerAsInitialised : Packet
{
    private readonly ulong _handle;
    internal SetLocalPlayerAsInitialised(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 113u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
}

public sealed class UpdateSoftEnum : Packet
{
    private readonly ulong _handle;
    internal UpdateSoftEnum(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 114u;
    public string EnumType { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public Value Options => new(_handle, 1);
    public byte ActionType { get => checked((byte)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
}

public sealed class NetworkStackLatency : Packet
{
    private readonly ulong _handle;
    internal NetworkStackLatency(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 115u;
    public long Timestamp { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public bool NeedsResponse { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
}

public sealed class ScriptCustomEvent : Packet
{
    private readonly ulong _handle;
    internal ScriptCustomEvent(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 117u;
    public string EventName { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public byte[] EventData { get => PacketBridge.Bytes(_handle, 1); set => PacketBridge.SetBytes(_handle, 1, value); }
}

public sealed class SpawnParticleEffect : Packet
{
    private readonly ulong _handle;
    internal SpawnParticleEffect(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 118u;
    public byte Dimension { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 2); set => PacketBridge.SetVector3(_handle, 2, value); }
    public string ParticleName { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public Value MoLangVariables => new(_handle, 4);
}

public sealed class AvailableActorIdentifiers : Packet
{
    private readonly ulong _handle;
    internal AvailableActorIdentifiers(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 119u;
    public byte[] SerialisedEntityIdentifiers { get => PacketBridge.Bytes(_handle, 0); set => PacketBridge.SetBytes(_handle, 0, value); }
}

public sealed class NetworkChunkPublisherUpdate : Packet
{
    private readonly ulong _handle;
    internal NetworkChunkPublisherUpdate(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 121u;
    public Value Position => new(_handle, 0);
    public uint Radius { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public Value SavedChunks => new(_handle, 2);
}

public sealed class BiomeDefinitionList : Packet
{
    private readonly ulong _handle;
    internal BiomeDefinitionList(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 122u;
    public Value BiomeDefinitions => new(_handle, 0);
    public Value StringList => new(_handle, 1);
}

public sealed class LevelSoundEvent : Packet
{
    private readonly ulong _handle;
    internal LevelSoundEvent(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 123u;
    public string SoundType { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 1); set => PacketBridge.SetVector3(_handle, 1, value); }
    public int ExtraData { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public string EntityType { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public bool BabyMob { get => PacketBridge.Bool(_handle, 4); set => PacketBridge.SetBool(_handle, 4, value); }
    public bool DisableRelativeVolume { get => PacketBridge.Bool(_handle, 5); set => PacketBridge.SetBool(_handle, 5, value); }
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 6)); set => PacketBridge.SetSigned(_handle, 6, value); }
    public Value FireAtPosition => new(_handle, 7);
}

public sealed class LevelEventGeneric : Packet
{
    private readonly ulong _handle;
    internal LevelEventGeneric(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 124u;
    public int EventID { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public byte[] SerialisedEventData { get => PacketBridge.Bytes(_handle, 1); set => PacketBridge.SetBytes(_handle, 1, value); }
}

public sealed class LecternUpdate : Packet
{
    private readonly ulong _handle;
    internal LecternUpdate(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 125u;
    public byte Page { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte PageCount { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public Value Position => new(_handle, 2);
}

public sealed class ClientCacheStatus : Packet
{
    private readonly ulong _handle;
    internal ClientCacheStatus(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 129u;
    public bool Enabled { get => PacketBridge.Bool(_handle, 0); set => PacketBridge.SetBool(_handle, 0, value); }
}

public sealed class OnScreenTextureAnimation : Packet
{
    private readonly ulong _handle;
    internal OnScreenTextureAnimation(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 130u;
    public uint AnimationType { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
}

public sealed class MapCreateLockedCopy : Packet
{
    private readonly ulong _handle;
    internal MapCreateLockedCopy(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 131u;
    public long OriginalMapID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public long NewMapID { get => checked((long)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
}

public sealed class StructureTemplateDataRequest : Packet
{
    private readonly ulong _handle;
    internal StructureTemplateDataRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 132u;
    public string StructureName { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public Value Position => new(_handle, 1);
    public Value Settings => new(_handle, 2);
    public byte RequestType { get => checked((byte)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
}

public sealed class StructureTemplateDataResponse : Packet
{
    private readonly ulong _handle;
    internal StructureTemplateDataResponse(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 133u;
    public string StructureName { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public bool Success { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
    public byte ResponseType { get => checked((byte)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public Value StructureTemplate => new(_handle, 3);
}

public sealed class ClientCacheBlobStatus : Packet
{
    private readonly ulong _handle;
    internal ClientCacheBlobStatus(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 135u;
    public Value MissHashes => new(_handle, 0);
    public Value HitHashes => new(_handle, 1);
}

public sealed class ClientCacheMissResponse : Packet
{
    private readonly ulong _handle;
    internal ClientCacheMissResponse(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 136u;
    public Value Blobs => new(_handle, 0);
}

public sealed class EducationSettings : Packet
{
    private readonly ulong _handle;
    internal EducationSettings(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 137u;
    public string CodeBuilderDefaultURI { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public string CodeBuilderTitle { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public bool CanResizeCodeBuilder { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
    public bool DisableLegacyTitleBar { get => PacketBridge.Bool(_handle, 3); set => PacketBridge.SetBool(_handle, 3, value); }
    public string PostProcessFilter { get => PacketBridge.String(_handle, 4); set => PacketBridge.SetString(_handle, 4, value); }
    public string ScreenshotBorderPath { get => PacketBridge.String(_handle, 5); set => PacketBridge.SetString(_handle, 5, value); }
    public Value CanModifyBlocks => new(_handle, 6);
    public Value OverrideURI => new(_handle, 7);
    public bool HasQuiz { get => PacketBridge.Bool(_handle, 8); set => PacketBridge.SetBool(_handle, 8, value); }
    public Value ExternalLinkSettings => new(_handle, 9);
}

public sealed class Emote : Packet
{
    private readonly ulong _handle;
    internal Emote(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 138u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public uint EmoteLength { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public string EmoteID { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public string XUID { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public string PlatformID { get => PacketBridge.String(_handle, 4); set => PacketBridge.SetString(_handle, 4, value); }
    public byte Flags { get => checked((byte)PacketBridge.Unsigned(_handle, 5)); set => PacketBridge.SetUnsigned(_handle, 5, value); }
}

public sealed class MultiPlayerSettings : Packet
{
    private readonly ulong _handle;
    internal MultiPlayerSettings(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 139u;
    public int ActionType { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class SettingsCommand : Packet
{
    private readonly ulong _handle;
    internal SettingsCommand(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 140u;
    public string CommandLine { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public bool SuppressOutput { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
}

public sealed class AnvilDamage : Packet
{
    private readonly ulong _handle;
    internal AnvilDamage(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 141u;
    public byte Damage { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value AnvilPosition => new(_handle, 1);
}

public sealed class CompletedUsingItem : Packet
{
    private readonly ulong _handle;
    internal CompletedUsingItem(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 142u;
    public short UsedItemID { get => checked((short)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public int UseMethod { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
}

public sealed class NetworkSettings : Packet
{
    private readonly ulong _handle;
    internal NetworkSettings(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 143u;
    public ushort CompressionThreshold { get => checked((ushort)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public ushort CompressionAlgorithm { get => checked((ushort)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public bool ClientThrottle { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
    public byte ClientThrottleThreshold { get => checked((byte)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
    public float ClientThrottleScalar { get => (float)PacketBridge.Number(_handle, 4); set => PacketBridge.SetNumber(_handle, 4, value); }
}

public sealed class PlayerAuthInput : Packet
{
    private readonly ulong _handle;
    internal PlayerAuthInput(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 144u;
    public float Pitch { get => (float)PacketBridge.Number(_handle, 0); set => PacketBridge.SetNumber(_handle, 0, value); }
    public float Yaw { get => (float)PacketBridge.Number(_handle, 1); set => PacketBridge.SetNumber(_handle, 1, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 2); set => PacketBridge.SetVector3(_handle, 2, value); }
    public Vector2 MoveVector { get => PacketBridge.Vector2(_handle, 3); set => PacketBridge.SetVector2(_handle, 3, value); }
    public float HeadYaw { get => (float)PacketBridge.Number(_handle, 4); set => PacketBridge.SetNumber(_handle, 4, value); }
    public Value InputData => new(_handle, 5);
    public uint InputMode { get => checked((uint)PacketBridge.Unsigned(_handle, 6)); set => PacketBridge.SetUnsigned(_handle, 6, value); }
    public uint PlayMode { get => checked((uint)PacketBridge.Unsigned(_handle, 7)); set => PacketBridge.SetUnsigned(_handle, 7, value); }
    public uint InteractionModel { get => checked((uint)PacketBridge.Unsigned(_handle, 8)); set => PacketBridge.SetUnsigned(_handle, 8, value); }
    public float InteractPitch { get => (float)PacketBridge.Number(_handle, 9); set => PacketBridge.SetNumber(_handle, 9, value); }
    public float InteractYaw { get => (float)PacketBridge.Number(_handle, 10); set => PacketBridge.SetNumber(_handle, 10, value); }
    public ulong Tick { get => checked((ulong)PacketBridge.Unsigned(_handle, 11)); set => PacketBridge.SetUnsigned(_handle, 11, value); }
    public Vector3 Delta { get => PacketBridge.Vector3(_handle, 12); set => PacketBridge.SetVector3(_handle, 12, value); }
    public Value ItemInteractionData => new(_handle, 13);
    public Value ItemStackRequest => new(_handle, 14);
    public Value BlockActions => new(_handle, 15);
    public Vector2 VehicleRotation { get => PacketBridge.Vector2(_handle, 16); set => PacketBridge.SetVector2(_handle, 16, value); }
    public long ClientPredictedVehicle { get => checked((long)PacketBridge.Signed(_handle, 17)); set => PacketBridge.SetSigned(_handle, 17, value); }
    public Vector2 AnalogueMoveVector { get => PacketBridge.Vector2(_handle, 18); set => PacketBridge.SetVector2(_handle, 18, value); }
    public Vector3 CameraOrientation { get => PacketBridge.Vector3(_handle, 19); set => PacketBridge.SetVector3(_handle, 19, value); }
    public Vector2 RawMoveVector { get => PacketBridge.Vector2(_handle, 20); set => PacketBridge.SetVector2(_handle, 20, value); }
}

public sealed class CreativeContent : Packet
{
    private readonly ulong _handle;
    internal CreativeContent(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 145u;
    public Value Groups => new(_handle, 0);
    public Value Items => new(_handle, 1);
}

public sealed class PlayerEnchantOptions : Packet
{
    private readonly ulong _handle;
    internal PlayerEnchantOptions(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 146u;
    public Value Options => new(_handle, 0);
}

public sealed class ItemStackRequest : Packet
{
    private readonly ulong _handle;
    internal ItemStackRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 147u;
    public Value Requests => new(_handle, 0);
}

public sealed class ItemStackResponse : Packet
{
    private readonly ulong _handle;
    internal ItemStackResponse(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 148u;
    public Value Responses => new(_handle, 0);
}

public sealed class PlayerArmourDamage : Packet
{
    private readonly ulong _handle;
    internal PlayerArmourDamage(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 149u;
    public Value List => new(_handle, 0);
}

public sealed class CodeBuilder : Packet
{
    private readonly ulong _handle;
    internal CodeBuilder(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 150u;
    public string URL { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public bool ShouldOpenCodeBuilder { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
}

public sealed class UpdatePlayerGameType : Packet
{
    private readonly ulong _handle;
    internal UpdatePlayerGameType(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 151u;
    public int GameType { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public long PlayerUniqueID { get => checked((long)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public ulong Tick { get => checked((ulong)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
}

public sealed class EmoteList : Packet
{
    private readonly ulong _handle;
    internal EmoteList(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 152u;
    public ulong PlayerRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value EmotePieces => new(_handle, 1);
}

public sealed class PositionTrackingDBServerBroadcast : Packet
{
    private readonly ulong _handle;
    internal PositionTrackingDBServerBroadcast(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 153u;
    public byte BroadcastAction { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public int TrackingID { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public Value Payload => new(_handle, 2);
}

public sealed class PositionTrackingDBClientRequest : Packet
{
    private readonly ulong _handle;
    internal PositionTrackingDBClientRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 154u;
    public byte RequestAction { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public int TrackingID { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
}

public sealed class DebugInfo : Packet
{
    private readonly ulong _handle;
    internal DebugInfo(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 155u;
    public long PlayerUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public byte[] Data { get => PacketBridge.Bytes(_handle, 1); set => PacketBridge.SetBytes(_handle, 1, value); }
}

public sealed class PacketViolationWarning : Packet
{
    private readonly ulong _handle;
    internal PacketViolationWarning(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 156u;
    public int Type { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public int Severity { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public int PacketID { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public string ViolationContext { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
}

public sealed class MotionPredictionHints : Packet
{
    private readonly ulong _handle;
    internal MotionPredictionHints(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 157u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Vector3 Velocity { get => PacketBridge.Vector3(_handle, 1); set => PacketBridge.SetVector3(_handle, 1, value); }
    public bool OnGround { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
}

public sealed class AnimateEntity : Packet
{
    private readonly ulong _handle;
    internal AnimateEntity(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 158u;
    public string Animation { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public string NextState { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public string StopCondition { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public int StopConditionVersion { get => checked((int)PacketBridge.Signed(_handle, 3)); set => PacketBridge.SetSigned(_handle, 3, value); }
    public string Controller { get => PacketBridge.String(_handle, 4); set => PacketBridge.SetString(_handle, 4, value); }
    public float BlendOutTime { get => (float)PacketBridge.Number(_handle, 5); set => PacketBridge.SetNumber(_handle, 5, value); }
    public Value EntityRuntimeIDs => new(_handle, 6);
}

public sealed class CameraShake : Packet
{
    private readonly ulong _handle;
    internal CameraShake(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 159u;
    public float Intensity { get => (float)PacketBridge.Number(_handle, 0); set => PacketBridge.SetNumber(_handle, 0, value); }
    public float Duration { get => (float)PacketBridge.Number(_handle, 1); set => PacketBridge.SetNumber(_handle, 1, value); }
    public byte Type { get => checked((byte)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public byte Action { get => checked((byte)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
}

public sealed class PlayerFog : Packet
{
    private readonly ulong _handle;
    internal PlayerFog(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 160u;
    public Value Stack => new(_handle, 0);
}

public sealed class CorrectPlayerMovePrediction : Packet
{
    private readonly ulong _handle;
    internal CorrectPlayerMovePrediction(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 161u;
    public byte PredictionType { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 1); set => PacketBridge.SetVector3(_handle, 1, value); }
    public Vector3 Delta { get => PacketBridge.Vector3(_handle, 2); set => PacketBridge.SetVector3(_handle, 2, value); }
    public Vector2 Rotation { get => PacketBridge.Vector2(_handle, 3); set => PacketBridge.SetVector2(_handle, 3, value); }
    public Value VehicleAngularVelocity => new(_handle, 4);
    public bool OnGround { get => PacketBridge.Bool(_handle, 5); set => PacketBridge.SetBool(_handle, 5, value); }
    public ulong Tick { get => checked((ulong)PacketBridge.Unsigned(_handle, 6)); set => PacketBridge.SetUnsigned(_handle, 6, value); }
}

public sealed class ItemRegistry : Packet
{
    private readonly ulong _handle;
    internal ItemRegistry(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 162u;
    public Value Items => new(_handle, 0);
}

public sealed class FilterText : Packet
{
    private readonly ulong _handle;
    internal FilterText(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 163u;
    public string Text { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public bool FromServer { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
}

public sealed class ClientBoundDebugRenderer : Packet
{
    private readonly ulong _handle;
    internal ClientBoundDebugRenderer(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 164u;
    public uint Type { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public string Text { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 2); set => PacketBridge.SetVector3(_handle, 2, value); }
    public float Red { get => (float)PacketBridge.Number(_handle, 3); set => PacketBridge.SetNumber(_handle, 3, value); }
    public float Green { get => (float)PacketBridge.Number(_handle, 4); set => PacketBridge.SetNumber(_handle, 4, value); }
    public float Blue { get => (float)PacketBridge.Number(_handle, 5); set => PacketBridge.SetNumber(_handle, 5, value); }
    public float Alpha { get => (float)PacketBridge.Number(_handle, 6); set => PacketBridge.SetNumber(_handle, 6, value); }
    public ulong Duration { get => checked((ulong)PacketBridge.Unsigned(_handle, 7)); set => PacketBridge.SetUnsigned(_handle, 7, value); }
}

public sealed class SyncActorProperty : Packet
{
    private readonly ulong _handle;
    internal SyncActorProperty(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 165u;
    public Value PropertyData => new(_handle, 0);
}

public sealed class AddVolumeEntity : Packet
{
    private readonly ulong _handle;
    internal AddVolumeEntity(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 166u;
    public uint EntityRuntimeID { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value EntityMetadata => new(_handle, 1);
    public string EncodingIdentifier { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public string InstanceIdentifier { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public Value Bounds => new(_handle, 4);
    public int Dimension { get => checked((int)PacketBridge.Signed(_handle, 5)); set => PacketBridge.SetSigned(_handle, 5, value); }
    public string EngineVersion { get => PacketBridge.String(_handle, 6); set => PacketBridge.SetString(_handle, 6, value); }
}

public sealed class RemoveVolumeEntity : Packet
{
    private readonly ulong _handle;
    internal RemoveVolumeEntity(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 167u;
    public uint EntityRuntimeID { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public int Dimension { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
}

public sealed class SimulationType : Packet
{
    private readonly ulong _handle;
    internal SimulationType(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 168u;
    public byte SimulationTypeValue { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
}

public sealed class NPCDialogue : Packet
{
    private readonly ulong _handle;
    internal NPCDialogue(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 169u;
    public ulong EntityUniqueID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public int ActionType { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public string Dialogue { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public string SceneName { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public string NPCName { get => PacketBridge.String(_handle, 4); set => PacketBridge.SetString(_handle, 4, value); }
    public string ActionJSON { get => PacketBridge.String(_handle, 5); set => PacketBridge.SetString(_handle, 5, value); }
}

public sealed class EducationResourceURI : Packet
{
    private readonly ulong _handle;
    internal EducationResourceURI(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 170u;
    public Value Resource => new(_handle, 0);
}

public sealed class CreatePhoto : Packet
{
    private readonly ulong _handle;
    internal CreatePhoto(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 171u;
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public string PhotoName { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public string ItemName { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
}

public sealed class UpdateSubChunkBlocks : Packet
{
    private readonly ulong _handle;
    internal UpdateSubChunkBlocks(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 172u;
    public Value Position => new(_handle, 0);
    public Value Blocks => new(_handle, 1);
    public Value Extra => new(_handle, 2);
}

public sealed class PhotoInfoRequest : Packet
{
    private readonly ulong _handle;
    internal PhotoInfoRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 173u;
    public long PhotoID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class SubChunk : Packet
{
    private readonly ulong _handle;
    internal SubChunk(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 174u;
    public bool CacheEnabled { get => PacketBridge.Bool(_handle, 0); set => PacketBridge.SetBool(_handle, 0, value); }
    public int Dimension { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public Value Position => new(_handle, 2);
    public Value SubChunkEntries => new(_handle, 3);
}

public sealed class SubChunkRequest : Packet
{
    private readonly ulong _handle;
    internal SubChunkRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 175u;
    public int Dimension { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public Value Offsets => new(_handle, 1);
    public Value Position => new(_handle, 2);
}

public sealed class ClientStartItemCooldown : Packet
{
    private readonly ulong _handle;
    internal ClientStartItemCooldown(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 176u;
    public string Category { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public int Duration { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
}

public sealed class ScriptMessage : Packet
{
    private readonly ulong _handle;
    internal ScriptMessage(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 177u;
    public string Identifier { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public byte[] Data { get => PacketBridge.Bytes(_handle, 1); set => PacketBridge.SetBytes(_handle, 1, value); }
}

public sealed class CodeBuilderSource : Packet
{
    private readonly ulong _handle;
    internal CodeBuilderSource(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 178u;
    public byte Operation { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public byte Category { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public byte CodeStatus { get => checked((byte)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
}

public sealed class TickingAreasLoadStatus : Packet
{
    private readonly ulong _handle;
    internal TickingAreasLoadStatus(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 179u;
    public bool Preload { get => PacketBridge.Bool(_handle, 0); set => PacketBridge.SetBool(_handle, 0, value); }
}

public sealed class DimensionData : Packet
{
    private readonly ulong _handle;
    internal DimensionData(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 180u;
    public Value Definitions => new(_handle, 0);
}

public sealed class AgentAction : Packet
{
    private readonly ulong _handle;
    internal AgentAction(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 181u;
    public string Identifier { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public int Action { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public byte[] Response { get => PacketBridge.Bytes(_handle, 2); set => PacketBridge.SetBytes(_handle, 2, value); }
}

public sealed class ChangeMobProperty : Packet
{
    private readonly ulong _handle;
    internal ChangeMobProperty(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 182u;
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public string Property { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public bool BoolValue { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
    public string StringValue { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public int IntValue { get => checked((int)PacketBridge.Signed(_handle, 4)); set => PacketBridge.SetSigned(_handle, 4, value); }
    public float FloatValue { get => (float)PacketBridge.Number(_handle, 5); set => PacketBridge.SetNumber(_handle, 5, value); }
}

public sealed class LessonProgress : Packet
{
    private readonly ulong _handle;
    internal LessonProgress(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 183u;
    public string Identifier { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public int Action { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public int Score { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
}

public sealed class RequestAbility : Packet
{
    private readonly ulong _handle;
    internal RequestAbility(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 184u;
    public int Ability { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public Value Value => new(_handle, 1);
}

public sealed class RequestPermissions : Packet
{
    private readonly ulong _handle;
    internal RequestPermissions(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 185u;
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public int PermissionLevel { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public ushort RequestedPermissions { get => checked((ushort)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
}

public sealed class ToastRequest : Packet
{
    private readonly ulong _handle;
    internal ToastRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 186u;
    public string Title { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public string Message { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
}

public sealed class UpdateAbilities : Packet
{
    private readonly ulong _handle;
    internal UpdateAbilities(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 187u;
    public Value AbilityData => new(_handle, 0);
}

public sealed class UpdateAdventureSettings : Packet
{
    private readonly ulong _handle;
    internal UpdateAdventureSettings(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 188u;
    public bool NoPvM { get => PacketBridge.Bool(_handle, 0); set => PacketBridge.SetBool(_handle, 0, value); }
    public bool NoMvP { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
    public bool ImmutableWorld { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
    public bool ShowNameTags { get => PacketBridge.Bool(_handle, 3); set => PacketBridge.SetBool(_handle, 3, value); }
    public bool AutoJump { get => PacketBridge.Bool(_handle, 4); set => PacketBridge.SetBool(_handle, 4, value); }
}

public sealed class DeathInfo : Packet
{
    private readonly ulong _handle;
    internal DeathInfo(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 189u;
    public string Cause { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public Value Messages => new(_handle, 1);
}

public sealed class EditorNetwork : Packet
{
    private readonly ulong _handle;
    internal EditorNetwork(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 190u;
    public bool RouteToManager { get => PacketBridge.Bool(_handle, 0); set => PacketBridge.SetBool(_handle, 0, value); }
    public Value Payload => new(_handle, 1);
}

public sealed class FeatureRegistry : Packet
{
    private readonly ulong _handle;
    internal FeatureRegistry(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 191u;
    public Value Features => new(_handle, 0);
}

public sealed class ServerStats : Packet
{
    private readonly ulong _handle;
    internal ServerStats(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 192u;
    public float ServerTime { get => (float)PacketBridge.Number(_handle, 0); set => PacketBridge.SetNumber(_handle, 0, value); }
    public float NetworkTime { get => (float)PacketBridge.Number(_handle, 1); set => PacketBridge.SetNumber(_handle, 1, value); }
}

public sealed class RequestNetworkSettings : Packet
{
    private readonly ulong _handle;
    internal RequestNetworkSettings(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 193u;
    public int ClientProtocol { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class GameTestRequest : Packet
{
    private readonly ulong _handle;
    internal GameTestRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 194u;
    public string Name { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public byte Rotation { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public int Repetitions { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public Value Position => new(_handle, 3);
    public bool StopOnError { get => PacketBridge.Bool(_handle, 4); set => PacketBridge.SetBool(_handle, 4, value); }
    public int TestsPerRow { get => checked((int)PacketBridge.Signed(_handle, 5)); set => PacketBridge.SetSigned(_handle, 5, value); }
    public int MaxTestsPerBatch { get => checked((int)PacketBridge.Signed(_handle, 6)); set => PacketBridge.SetSigned(_handle, 6, value); }
}

public sealed class GameTestResults : Packet
{
    private readonly ulong _handle;
    internal GameTestResults(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 195u;
    public string Name { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public bool Succeeded { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
    public string Error { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
}

public sealed class UpdateClientInputLocks : Packet
{
    private readonly ulong _handle;
    internal UpdateClientInputLocks(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 196u;
    public uint Locks { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
}

public sealed class ClientCheatAbility : Packet
{
    private readonly ulong _handle;
    internal ClientCheatAbility(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 197u;
    public Value AbilityData => new(_handle, 0);
}

public sealed class CameraPresets : Packet
{
    private readonly ulong _handle;
    internal CameraPresets(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 198u;
    public Value Presets => new(_handle, 0);
}

public sealed class UnlockedRecipes : Packet
{
    private readonly ulong _handle;
    internal UnlockedRecipes(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 199u;
    public uint UnlockType { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value Recipes => new(_handle, 1);
}

public sealed class CameraInstruction : Packet
{
    private readonly ulong _handle;
    internal CameraInstruction(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 300u;
    public Value Set => new(_handle, 0);
    public Value Clear => new(_handle, 1);
    public Value Fade => new(_handle, 2);
    public Value Target => new(_handle, 3);
    public Value RemoveTarget => new(_handle, 4);
    public Value FieldOfView => new(_handle, 5);
    public Value Spline => new(_handle, 6);
    public Value AttachToEntity => new(_handle, 7);
    public Value DetachFromEntity => new(_handle, 8);
}

public sealed class TrimData : Packet
{
    private readonly ulong _handle;
    internal TrimData(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 302u;
    public Value Patterns => new(_handle, 0);
    public Value Materials => new(_handle, 1);
}

public sealed class OpenSign : Packet
{
    private readonly ulong _handle;
    internal OpenSign(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 303u;
    public Value Position => new(_handle, 0);
    public bool FrontSide { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
}

public sealed class AgentAnimation : Packet
{
    private readonly ulong _handle;
    internal AgentAnimation(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 304u;
    public byte Animation { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
}

public sealed class RefreshEntitlements : Packet
{
    private readonly ulong _handle;
    internal RefreshEntitlements(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 305u;
}

public sealed class PlayerToggleCrafterSlotRequest : Packet
{
    private readonly ulong _handle;
    internal PlayerToggleCrafterSlotRequest(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 306u;
    public int PosX { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public int PosY { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public int PosZ { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public byte Slot { get => checked((byte)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
    public bool Disabled { get => PacketBridge.Bool(_handle, 4); set => PacketBridge.SetBool(_handle, 4, value); }
}

public sealed class SetPlayerInventoryOptions : Packet
{
    private readonly ulong _handle;
    internal SetPlayerInventoryOptions(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 307u;
    public int LeftInventoryTab { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public int RightInventoryTab { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public bool Filtering { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
    public int InventoryLayout { get => checked((int)PacketBridge.Signed(_handle, 3)); set => PacketBridge.SetSigned(_handle, 3, value); }
    public int CraftingLayout { get => checked((int)PacketBridge.Signed(_handle, 4)); set => PacketBridge.SetSigned(_handle, 4, value); }
}

public sealed class SetHud : Packet
{
    private readonly ulong _handle;
    internal SetHud(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 308u;
    public Value Elements => new(_handle, 0);
    public int Visibility { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
}

public sealed class AwardAchievement : Packet
{
    private readonly ulong _handle;
    internal AwardAchievement(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 309u;
    public int AchievementID { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
}

public sealed class ClientBoundCloseForm : Packet
{
    private readonly ulong _handle;
    internal ClientBoundCloseForm(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 310u;
}

public sealed class ServerBoundLoadingScreen : Packet
{
    private readonly ulong _handle;
    internal ServerBoundLoadingScreen(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 312u;
    public int Type { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public Value LoadingScreenID => new(_handle, 1);
}

public sealed class JigsawStructureData : Packet
{
    private readonly ulong _handle;
    internal JigsawStructureData(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 313u;
    public Value StructureData => new(_handle, 0);
}

public sealed class CurrentStructureFeature : Packet
{
    private readonly ulong _handle;
    internal CurrentStructureFeature(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 314u;
    public string CurrentFeature { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
}

public sealed class ServerBoundDiagnostics : Packet
{
    private readonly ulong _handle;
    internal ServerBoundDiagnostics(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 315u;
    public float AverageFramesPerSecond { get => (float)PacketBridge.Number(_handle, 0); set => PacketBridge.SetNumber(_handle, 0, value); }
    public float AverageServerSimTickTime { get => (float)PacketBridge.Number(_handle, 1); set => PacketBridge.SetNumber(_handle, 1, value); }
    public float AverageClientSimTickTime { get => (float)PacketBridge.Number(_handle, 2); set => PacketBridge.SetNumber(_handle, 2, value); }
    public float AverageBeginFrameTime { get => (float)PacketBridge.Number(_handle, 3); set => PacketBridge.SetNumber(_handle, 3, value); }
    public float AverageInputTime { get => (float)PacketBridge.Number(_handle, 4); set => PacketBridge.SetNumber(_handle, 4, value); }
    public float AverageRenderTime { get => (float)PacketBridge.Number(_handle, 5); set => PacketBridge.SetNumber(_handle, 5, value); }
    public float AverageEndFrameTime { get => (float)PacketBridge.Number(_handle, 6); set => PacketBridge.SetNumber(_handle, 6, value); }
    public float AverageRemainderTimePercent { get => (float)PacketBridge.Number(_handle, 7); set => PacketBridge.SetNumber(_handle, 7, value); }
    public float AverageUnaccountedTimePercent { get => (float)PacketBridge.Number(_handle, 8); set => PacketBridge.SetNumber(_handle, 8, value); }
    public Value MemoryCategoryValues => new(_handle, 9);
    public Value EntityDiagnostics => new(_handle, 10);
    public Value SystemDiagnostics => new(_handle, 11);
    public Value WhiskerScopes => new(_handle, 12);
}

public sealed class CameraAimAssist : Packet
{
    private readonly ulong _handle;
    internal CameraAimAssist(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 316u;
    public string Preset { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public Vector2 Angle { get => PacketBridge.Vector2(_handle, 1); set => PacketBridge.SetVector2(_handle, 1, value); }
    public float Distance { get => (float)PacketBridge.Number(_handle, 2); set => PacketBridge.SetNumber(_handle, 2, value); }
    public byte TargetMode { get => checked((byte)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
    public byte Action { get => checked((byte)PacketBridge.Unsigned(_handle, 4)); set => PacketBridge.SetUnsigned(_handle, 4, value); }
    public bool ShowDebugRender { get => PacketBridge.Bool(_handle, 5); set => PacketBridge.SetBool(_handle, 5, value); }
}

public sealed class ContainerRegistryCleanup : Packet
{
    private readonly ulong _handle;
    internal ContainerRegistryCleanup(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 317u;
    public Value RemovedContainers => new(_handle, 0);
}

public sealed class MovementEffect : Packet
{
    private readonly ulong _handle;
    internal MovementEffect(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 318u;
    public ulong EntityRuntimeID { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public int Type { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public int Duration { get => checked((int)PacketBridge.Signed(_handle, 2)); set => PacketBridge.SetSigned(_handle, 2, value); }
    public ulong Tick { get => checked((ulong)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
}

public sealed class CameraAimAssistPresets : Packet
{
    private readonly ulong _handle;
    internal CameraAimAssistPresets(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 320u;
    public Value Categories => new(_handle, 0);
    public Value Presets => new(_handle, 1);
    public byte Operation { get => checked((byte)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
}

public sealed class ClientCameraAimAssist : Packet
{
    private readonly ulong _handle;
    internal ClientCameraAimAssist(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 321u;
    public string PresetID { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public byte Action { get => checked((byte)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public bool AllowAimAssist { get => PacketBridge.Bool(_handle, 2); set => PacketBridge.SetBool(_handle, 2, value); }
}

public sealed class ClientMovementPredictionSync : Packet
{
    private readonly ulong _handle;
    internal ClientMovementPredictionSync(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 322u;
    public Value ActorFlags => new(_handle, 0);
    public float BoundingBoxScale { get => (float)PacketBridge.Number(_handle, 1); set => PacketBridge.SetNumber(_handle, 1, value); }
    public float BoundingBoxWidth { get => (float)PacketBridge.Number(_handle, 2); set => PacketBridge.SetNumber(_handle, 2, value); }
    public float BoundingBoxHeight { get => (float)PacketBridge.Number(_handle, 3); set => PacketBridge.SetNumber(_handle, 3, value); }
    public float MovementSpeed { get => (float)PacketBridge.Number(_handle, 4); set => PacketBridge.SetNumber(_handle, 4, value); }
    public float UnderwaterMovementSpeed { get => (float)PacketBridge.Number(_handle, 5); set => PacketBridge.SetNumber(_handle, 5, value); }
    public float LavaMovementSpeed { get => (float)PacketBridge.Number(_handle, 6); set => PacketBridge.SetNumber(_handle, 6, value); }
    public float JumpStrength { get => (float)PacketBridge.Number(_handle, 7); set => PacketBridge.SetNumber(_handle, 7, value); }
    public float Health { get => (float)PacketBridge.Number(_handle, 8); set => PacketBridge.SetNumber(_handle, 8, value); }
    public float Hunger { get => (float)PacketBridge.Number(_handle, 9); set => PacketBridge.SetNumber(_handle, 9, value); }
    public float FrictionModifier { get => (float)PacketBridge.Number(_handle, 10); set => PacketBridge.SetNumber(_handle, 10, value); }
    public float Bounciness { get => (float)PacketBridge.Number(_handle, 11); set => PacketBridge.SetNumber(_handle, 11, value); }
    public float AirDragModifier { get => (float)PacketBridge.Number(_handle, 12); set => PacketBridge.SetNumber(_handle, 12, value); }
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 13)); set => PacketBridge.SetSigned(_handle, 13, value); }
    public bool Flying { get => PacketBridge.Bool(_handle, 14); set => PacketBridge.SetBool(_handle, 14, value); }
}

public sealed class UpdateClientOptions : Packet
{
    private readonly ulong _handle;
    internal UpdateClientOptions(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 323u;
    public Value GraphicsMode => new(_handle, 0);
    public Value FilterProfanity => new(_handle, 1);
}

public sealed class PlayerVideoCapture : Packet
{
    private readonly ulong _handle;
    internal PlayerVideoCapture(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 324u;
    public byte Action { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public int FrameRate { get => checked((int)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public string FilePrefix { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
}

public sealed class PlayerUpdateEntityOverrides : Packet
{
    private readonly ulong _handle;
    internal PlayerUpdateEntityOverrides(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 325u;
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public uint PropertyIndex { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public byte Type { get => checked((byte)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
    public int IntValue { get => checked((int)PacketBridge.Signed(_handle, 3)); set => PacketBridge.SetSigned(_handle, 3, value); }
    public float FloatValue { get => (float)PacketBridge.Number(_handle, 4); set => PacketBridge.SetNumber(_handle, 4, value); }
}

public sealed class PlayerLocation : Packet
{
    private readonly ulong _handle;
    internal PlayerLocation(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 326u;
    public int Type { get => checked((int)PacketBridge.Signed(_handle, 0)); set => PacketBridge.SetSigned(_handle, 0, value); }
    public long EntityUniqueID { get => checked((long)PacketBridge.Signed(_handle, 1)); set => PacketBridge.SetSigned(_handle, 1, value); }
    public Vector3 Position { get => PacketBridge.Vector3(_handle, 2); set => PacketBridge.SetVector3(_handle, 2, value); }
}

public sealed class ClientBoundControlSchemeSet : Packet
{
    private readonly ulong _handle;
    internal ClientBoundControlSchemeSet(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 327u;
    public byte ControlScheme { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
}

public sealed class PrimitiveShapes : Packet
{
    private readonly ulong _handle;
    internal PrimitiveShapes(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 328u;
    public Value Shapes => new(_handle, 0);
}

public sealed class ServerBoundPackSettingChange : Packet
{
    private readonly ulong _handle;
    internal ServerBoundPackSettingChange(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 329u;
    public Guid PackID { get => PacketBridge.Guid(_handle, 0); set => PacketBridge.SetGuid(_handle, 0, value); }
    public Value PackSetting => new(_handle, 1);
}

public sealed class ClientBoundDataStore : Packet
{
    private readonly ulong _handle;
    internal ClientBoundDataStore(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 330u;
    public Value Updates => new(_handle, 0);
}

public sealed class GraphicsOverrideParameter : Packet
{
    private readonly ulong _handle;
    internal GraphicsOverrideParameter(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 331u;
    public Value Values => new(_handle, 0);
    public Value FloatValue => new(_handle, 1);
    public Value Vec3Value => new(_handle, 2);
    public string BiomeIdentifier { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public Value PlayerIdentifier => new(_handle, 4);
    public byte ParameterType { get => checked((byte)PacketBridge.Unsigned(_handle, 5)); set => PacketBridge.SetUnsigned(_handle, 5, value); }
    public bool Reset { get => PacketBridge.Bool(_handle, 6); set => PacketBridge.SetBool(_handle, 6, value); }
}

public sealed class ServerBoundDataStore : Packet
{
    private readonly ulong _handle;
    internal ServerBoundDataStore(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 332u;
    public Value Update => new(_handle, 0);
}

public sealed class ClientBoundDataDrivenUIShowScreen : Packet
{
    private readonly ulong _handle;
    internal ClientBoundDataDrivenUIShowScreen(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 333u;
    public string ScreenID { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public uint FormID { get => checked((uint)PacketBridge.Unsigned(_handle, 1)); set => PacketBridge.SetUnsigned(_handle, 1, value); }
    public Value DataInstanceID => new(_handle, 2);
}

public sealed class ClientBoundDataDrivenUICloseScreen : Packet
{
    private readonly ulong _handle;
    internal ClientBoundDataDrivenUICloseScreen(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 334u;
    public Value FormID => new(_handle, 0);
}

public sealed class ClientBoundDataDrivenUIReload : Packet
{
    private readonly ulong _handle;
    internal ClientBoundDataDrivenUIReload(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 335u;
}

public sealed class ClientBoundTextureShift : Packet
{
    private readonly ulong _handle;
    internal ClientBoundTextureShift(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 336u;
    public byte ActionID { get => checked((byte)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public string CollectionName { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public string FromStep { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public string ToStep { get => PacketBridge.String(_handle, 3); set => PacketBridge.SetString(_handle, 3, value); }
    public Value AllSteps => new(_handle, 4);
    public ulong CurrentLengthTicks { get => checked((ulong)PacketBridge.Unsigned(_handle, 5)); set => PacketBridge.SetUnsigned(_handle, 5, value); }
    public ulong TotalLengthTicks { get => checked((ulong)PacketBridge.Unsigned(_handle, 6)); set => PacketBridge.SetUnsigned(_handle, 6, value); }
    public bool Enabled { get => PacketBridge.Bool(_handle, 7); set => PacketBridge.SetBool(_handle, 7, value); }
}

public sealed class VoxelShapes : Packet
{
    private readonly ulong _handle;
    internal VoxelShapes(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 337u;
    public Value Shapes => new(_handle, 0);
    public Value NameMap => new(_handle, 1);
    public ushort CustomShapeCount { get => checked((ushort)PacketBridge.Unsigned(_handle, 2)); set => PacketBridge.SetUnsigned(_handle, 2, value); }
}

public sealed class CameraSpline : Packet
{
    private readonly ulong _handle;
    internal CameraSpline(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 338u;
    public Value Splines => new(_handle, 0);
}

public sealed class CameraAimAssistActorPriority : Packet
{
    private readonly ulong _handle;
    internal CameraAimAssistActorPriority(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 339u;
    public Value PriorityData => new(_handle, 0);
}

public sealed class ResourcePacksReadyForValidation : Packet
{
    private readonly ulong _handle;
    internal ResourcePacksReadyForValidation(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 340u;
}

public sealed class LocatorBar : Packet
{
    private readonly ulong _handle;
    internal LocatorBar(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 341u;
    public Value Waypoints => new(_handle, 0);
}

public sealed class PartyChanged : Packet
{
    private readonly ulong _handle;
    internal PartyChanged(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 342u;
    public Value PartyInfo => new(_handle, 0);
}

public sealed class ServerBoundDataDrivenScreenClosed : Packet
{
    private readonly ulong _handle;
    internal ServerBoundDataDrivenScreenClosed(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 343u;
    public Value FormID => new(_handle, 0);
    public string CloseReason { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
}

public sealed class SyncWorldClocks : Packet
{
    private readonly ulong _handle;
    internal SyncWorldClocks(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 344u;
    public uint PayloadType { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value SyncStates => new(_handle, 1);
    public Value Clocks => new(_handle, 2);
    public ulong AddClockID { get => checked((ulong)PacketBridge.Unsigned(_handle, 3)); set => PacketBridge.SetUnsigned(_handle, 3, value); }
    public Value AddTimeMarkers => new(_handle, 4);
    public ulong RemoveClockID { get => checked((ulong)PacketBridge.Unsigned(_handle, 5)); set => PacketBridge.SetUnsigned(_handle, 5, value); }
    public Value RemoveTimeMarkerIDs => new(_handle, 6);
}

public sealed class ClientBoundAttributeLayerSync : Packet
{
    private readonly ulong _handle;
    internal ClientBoundAttributeLayerSync(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 345u;
    public uint PayloadType { get => checked((uint)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public Value Layers => new(_handle, 1);
    public string LayerName { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
    public int DimensionID { get => checked((int)PacketBridge.Signed(_handle, 3)); set => PacketBridge.SetSigned(_handle, 3, value); }
    public Value Settings => new(_handle, 4);
    public Value EnvironmentAttributes => new(_handle, 5);
    public Value RemoveAttributeNames => new(_handle, 6);
}

public sealed class ServerStoreInfo : Packet
{
    private readonly ulong _handle;
    internal ServerStoreInfo(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 346u;
    public Value StoreInfo => new(_handle, 0);
}

public sealed class ServerPresenceInfo : Packet
{
    private readonly ulong _handle;
    internal ServerPresenceInfo(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 347u;
    public Value PresenceInfo => new(_handle, 0);
}

public sealed class ClientboundUpdateSoundData : Packet
{
    private readonly ulong _handle;
    internal ClientboundUpdateSoundData(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 348u;
    public ulong ServerSoundHandle { get => checked((ulong)PacketBridge.Unsigned(_handle, 0)); set => PacketBridge.SetUnsigned(_handle, 0, value); }
    public string SoundEvent { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
}

public sealed class SendPartyDestinationCookie : Packet
{
    private readonly ulong _handle;
    internal SendPartyDestinationCookie(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 349u;
    public string Cookie { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public string Intent { get => PacketBridge.String(_handle, 1); set => PacketBridge.SetString(_handle, 1, value); }
    public string DestinationName { get => PacketBridge.String(_handle, 2); set => PacketBridge.SetString(_handle, 2, value); }
}

public sealed class PartyDestinationCookieResponse : Packet
{
    private readonly ulong _handle;
    internal PartyDestinationCookieResponse(ulong handle) => _handle = handle;
    ulong Packet.HostHandle() => _handle;
    public uint ID() => 350u;
    public string Cookie { get => PacketBridge.String(_handle, 0); set => PacketBridge.SetString(_handle, 0, value); }
    public bool Accepted { get => PacketBridge.Bool(_handle, 1); set => PacketBridge.SetBool(_handle, 1, value); }
}

internal static class PacketCodec
{
    internal static Packet Decode(ulong handle, uint id) => id switch
    {
        1u => new Login(handle),
        2u => new PlayStatus(handle),
        3u => new ServerToClientHandshake(handle),
        4u => new ClientToServerHandshake(handle),
        5u => new Disconnect(handle),
        6u => new ResourcePacksInfo(handle),
        7u => new ResourcePackStack(handle),
        8u => new ResourcePackClientResponse(handle),
        9u => new Text(handle),
        10u => new SetTime(handle),
        11u => new StartGame(handle),
        12u => new AddPlayer(handle),
        13u => new AddActor(handle),
        14u => new RemoveActor(handle),
        15u => new AddItemActor(handle),
        17u => new TakeItemActor(handle),
        18u => new MoveActorAbsolute(handle),
        19u => new MovePlayer(handle),
        21u => new UpdateBlock(handle),
        22u => new AddPainting(handle),
        25u => new LevelEvent(handle),
        26u => new BlockEvent(handle),
        27u => new ActorEvent(handle),
        28u => new MobEffect(handle),
        29u => new UpdateAttributes(handle),
        30u => new InventoryTransaction(handle),
        31u => new MobEquipment(handle),
        32u => new MobArmourEquipment(handle),
        33u => new Interact(handle),
        34u => new BlockPickRequest(handle),
        35u => new ActorPickRequest(handle),
        36u => new PlayerAction(handle),
        38u => new HurtArmour(handle),
        39u => new SetActorData(handle),
        40u => new SetActorMotion(handle),
        41u => new SetActorLink(handle),
        42u => new SetHealth(handle),
        43u => new SetSpawnPosition(handle),
        44u => new Animate(handle),
        45u => new Respawn(handle),
        46u => new ContainerOpen(handle),
        47u => new ContainerClose(handle),
        48u => new PlayerHotBar(handle),
        49u => new InventoryContent(handle),
        50u => new InventorySlot(handle),
        51u => new ContainerSetData(handle),
        52u => new CraftingData(handle),
        54u => new GUIDataPickItem(handle),
        55u => new AdventureSettings(handle),
        56u => new BlockActorData(handle),
        58u => new LevelChunk(handle),
        59u => new SetCommandsEnabled(handle),
        60u => new SetDifficulty(handle),
        61u => new ChangeDimension(handle),
        62u => new SetPlayerGameType(handle),
        63u => new PlayerList(handle),
        64u => new SimpleEvent(handle),
        65u => new Event(handle),
        66u => new SpawnExperienceOrb(handle),
        67u => new ClientBoundMapItemData(handle),
        68u => new MapInfoRequest(handle),
        69u => new RequestChunkRadius(handle),
        70u => new ChunkRadiusUpdated(handle),
        72u => new GameRulesChanged(handle),
        73u => new Camera(handle),
        74u => new BossEvent(handle),
        75u => new ShowCredits(handle),
        76u => new AvailableCommands(handle),
        77u => new CommandRequest(handle),
        78u => new CommandBlockUpdate(handle),
        79u => new CommandOutput(handle),
        80u => new UpdateTrade(handle),
        81u => new UpdateEquip(handle),
        82u => new ResourcePackDataInfo(handle),
        83u => new ResourcePackChunkData(handle),
        84u => new ResourcePackChunkRequest(handle),
        85u => new Transfer(handle),
        86u => new PlaySound(handle),
        87u => new StopSound(handle),
        88u => new SetTitle(handle),
        89u => new AddBehaviourTree(handle),
        90u => new StructureBlockUpdate(handle),
        91u => new ShowStoreOffer(handle),
        92u => new PurchaseReceipt(handle),
        93u => new PlayerSkin(handle),
        94u => new SubClientLogin(handle),
        95u => new AutomationClientConnect(handle),
        96u => new SetLastHurtBy(handle),
        97u => new BookEdit(handle),
        98u => new NPCRequest(handle),
        99u => new PhotoTransfer(handle),
        100u => new ModalFormRequest(handle),
        101u => new ModalFormResponse(handle),
        102u => new ServerSettingsRequest(handle),
        103u => new ServerSettingsResponse(handle),
        104u => new ShowProfile(handle),
        105u => new SetDefaultGameType(handle),
        106u => new RemoveObjective(handle),
        107u => new SetDisplayObjective(handle),
        108u => new SetScore(handle),
        109u => new LabTable(handle),
        110u => new UpdateBlockSynced(handle),
        111u => new MoveActorDelta(handle),
        112u => new SetScoreboardIdentity(handle),
        113u => new SetLocalPlayerAsInitialised(handle),
        114u => new UpdateSoftEnum(handle),
        115u => new NetworkStackLatency(handle),
        117u => new ScriptCustomEvent(handle),
        118u => new SpawnParticleEffect(handle),
        119u => new AvailableActorIdentifiers(handle),
        121u => new NetworkChunkPublisherUpdate(handle),
        122u => new BiomeDefinitionList(handle),
        123u => new LevelSoundEvent(handle),
        124u => new LevelEventGeneric(handle),
        125u => new LecternUpdate(handle),
        129u => new ClientCacheStatus(handle),
        130u => new OnScreenTextureAnimation(handle),
        131u => new MapCreateLockedCopy(handle),
        132u => new StructureTemplateDataRequest(handle),
        133u => new StructureTemplateDataResponse(handle),
        135u => new ClientCacheBlobStatus(handle),
        136u => new ClientCacheMissResponse(handle),
        137u => new EducationSettings(handle),
        138u => new Emote(handle),
        139u => new MultiPlayerSettings(handle),
        140u => new SettingsCommand(handle),
        141u => new AnvilDamage(handle),
        142u => new CompletedUsingItem(handle),
        143u => new NetworkSettings(handle),
        144u => new PlayerAuthInput(handle),
        145u => new CreativeContent(handle),
        146u => new PlayerEnchantOptions(handle),
        147u => new ItemStackRequest(handle),
        148u => new ItemStackResponse(handle),
        149u => new PlayerArmourDamage(handle),
        150u => new CodeBuilder(handle),
        151u => new UpdatePlayerGameType(handle),
        152u => new EmoteList(handle),
        153u => new PositionTrackingDBServerBroadcast(handle),
        154u => new PositionTrackingDBClientRequest(handle),
        155u => new DebugInfo(handle),
        156u => new PacketViolationWarning(handle),
        157u => new MotionPredictionHints(handle),
        158u => new AnimateEntity(handle),
        159u => new CameraShake(handle),
        160u => new PlayerFog(handle),
        161u => new CorrectPlayerMovePrediction(handle),
        162u => new ItemRegistry(handle),
        163u => new FilterText(handle),
        164u => new ClientBoundDebugRenderer(handle),
        165u => new SyncActorProperty(handle),
        166u => new AddVolumeEntity(handle),
        167u => new RemoveVolumeEntity(handle),
        168u => new SimulationType(handle),
        169u => new NPCDialogue(handle),
        170u => new EducationResourceURI(handle),
        171u => new CreatePhoto(handle),
        172u => new UpdateSubChunkBlocks(handle),
        173u => new PhotoInfoRequest(handle),
        174u => new SubChunk(handle),
        175u => new SubChunkRequest(handle),
        176u => new ClientStartItemCooldown(handle),
        177u => new ScriptMessage(handle),
        178u => new CodeBuilderSource(handle),
        179u => new TickingAreasLoadStatus(handle),
        180u => new DimensionData(handle),
        181u => new AgentAction(handle),
        182u => new ChangeMobProperty(handle),
        183u => new LessonProgress(handle),
        184u => new RequestAbility(handle),
        185u => new RequestPermissions(handle),
        186u => new ToastRequest(handle),
        187u => new UpdateAbilities(handle),
        188u => new UpdateAdventureSettings(handle),
        189u => new DeathInfo(handle),
        190u => new EditorNetwork(handle),
        191u => new FeatureRegistry(handle),
        192u => new ServerStats(handle),
        193u => new RequestNetworkSettings(handle),
        194u => new GameTestRequest(handle),
        195u => new GameTestResults(handle),
        196u => new UpdateClientInputLocks(handle),
        197u => new ClientCheatAbility(handle),
        198u => new CameraPresets(handle),
        199u => new UnlockedRecipes(handle),
        300u => new CameraInstruction(handle),
        302u => new TrimData(handle),
        303u => new OpenSign(handle),
        304u => new AgentAnimation(handle),
        305u => new RefreshEntitlements(handle),
        306u => new PlayerToggleCrafterSlotRequest(handle),
        307u => new SetPlayerInventoryOptions(handle),
        308u => new SetHud(handle),
        309u => new AwardAchievement(handle),
        310u => new ClientBoundCloseForm(handle),
        312u => new ServerBoundLoadingScreen(handle),
        313u => new JigsawStructureData(handle),
        314u => new CurrentStructureFeature(handle),
        315u => new ServerBoundDiagnostics(handle),
        316u => new CameraAimAssist(handle),
        317u => new ContainerRegistryCleanup(handle),
        318u => new MovementEffect(handle),
        320u => new CameraAimAssistPresets(handle),
        321u => new ClientCameraAimAssist(handle),
        322u => new ClientMovementPredictionSync(handle),
        323u => new UpdateClientOptions(handle),
        324u => new PlayerVideoCapture(handle),
        325u => new PlayerUpdateEntityOverrides(handle),
        326u => new PlayerLocation(handle),
        327u => new ClientBoundControlSchemeSet(handle),
        328u => new PrimitiveShapes(handle),
        329u => new ServerBoundPackSettingChange(handle),
        330u => new ClientBoundDataStore(handle),
        331u => new GraphicsOverrideParameter(handle),
        332u => new ServerBoundDataStore(handle),
        333u => new ClientBoundDataDrivenUIShowScreen(handle),
        334u => new ClientBoundDataDrivenUICloseScreen(handle),
        335u => new ClientBoundDataDrivenUIReload(handle),
        336u => new ClientBoundTextureShift(handle),
        337u => new VoxelShapes(handle),
        338u => new CameraSpline(handle),
        339u => new CameraAimAssistActorPriority(handle),
        340u => new ResourcePacksReadyForValidation(handle),
        341u => new LocatorBar(handle),
        342u => new PartyChanged(handle),
        343u => new ServerBoundDataDrivenScreenClosed(handle),
        344u => new SyncWorldClocks(handle),
        345u => new ClientBoundAttributeLayerSync(handle),
        346u => new ServerStoreInfo(handle),
        347u => new ServerPresenceInfo(handle),
        348u => new ClientboundUpdateSoundData(handle),
        349u => new SendPartyDestinationCookie(handle),
        350u => new PartyDestinationCookieResponse(handle),
        _ => new Unknown(handle, id),
    };
}
