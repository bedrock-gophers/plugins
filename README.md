# bedrock-gophers/plugins

Write [Dragonfly](https://github.com/df-mc/dragonfly) plugins in C# and compile them to native
shared libraries with .NET NativeAOT.

The public C# API is generated from Dragonfly's Go AST and follows Dragonfly's types and method
names. The native ABI is private implementation detail. Dragonfly parity is broad, but still in
progress.

## Quick start

Requirements:

- Go 1.26+
- .NET 10 SDK
- A C compiler and linker

Build every source plugin, stage its native library, and start the example server:

```shell
make run
```

`make run` always regenerates the C# API and republishes the runtime and every
`examples/plugins/*/*.csproj`. It does not reuse old plugin builds or hot-reload a running server.
Stop the old process and reconnect after restarting so Bedrock receives updated command metadata.

Consumer repositories place plugins in `plugins/`. See the [`minimal`](../../tree/minimal) branch
for the smallest local and Docker setup. The loader supports compatible precompiled `.dll`, `.so`,
and `.dylib` files too.

## Write a plugin

A plugin inherits `Dragonfly.Plugin`. It does not declare an ID, library path, or native entry
point: the project name becomes the plugin ID, and the source generator emits the entry point.

```csharp
using Dragonfly;

public sealed class Example : Plugin
{
    public override void OnEnable() => Console.WriteLine("enabled");
    public override void OnDisable() => Console.WriteLine("disabled");

    public override void OnJoin(Player.Context ctx)
    {
        Console.WriteLine($"{ctx.Player().Name()} joined");
    }

    public override void HandleQuit(Player player)
    {
        Console.WriteLine($"{player.Name()} quit");
    }
}
```

Use an example project as the starting point; it already contains the required NativeAOT and source
generator references. [`LifecycleLogger`](examples/plugins/lifecycle-logger) is the smallest one.

`OnJoin` is a host lifecycle extension fired after Dragonfly installs the player handler. All other
player and world callbacks mirror Dragonfly's handler interfaces. Cancellable callbacks allow the
action by default.

### Connection admission

Override `Allow` to reject a connection before a `Player` exists:

```csharp
public override (string Message, bool Allowed) Allow(
    Net.Addr addr,
    Login.IdentityData identity,
    Login.ClientData client)
{
    if (Banned(identity.XUID)) return ("You are banned.", false);
    return (string.Empty, true);
}
```

All plugins allow connections by default. Hooks run in deterministic plugin load order, and the
first denial wins. This callback may run concurrently. Identity data is verified; client data is
controlled by the connecting client.

### Packet interception

Packets are already decoded by gophertunnel:

```csharp
using Packet = Dragonfly.Packet;

public override void HandleClientPacket(Packet.Context ctx, Packet.Packet packet)
{
    if (packet is Packet.Text text) text.Message = text.Message.Trim();
    if (packet is Packet.CommandRequest command && command.CommandLine.Length == 0) ctx.Cancel();
}

public override void HandleServerPacket(Packet.Context ctx, Packet.Packet packet)
{
    Console.WriteLine($"send {packet.ID()} to {ctx.XUID()}");
}
```

All registered packet structs are generated. Simple top-level fields are typed and mutable; complex
protocol fields temporarily expose `Packet.Value.Json()`. Outgoing packets may be inspected or
cancelled, but not mutated because broadcasts may share packet objects between connections.

## API coverage

Current surface includes:

- Plugin lifecycle and connection admission
- Reflected commands
- Player and world handler callbacks
- Player state, movement, interactions, visibility, presentation, and world transfer
- Typed items, inventories, effects, forms, scoreboards, titles, sounds, and container menus
- Typed blocks, liquids, biomes, particles, dimensions, game modes, and world transactions
- World creation, MCDB persistence, scheduling, entities, and custom entities
- Server-wide player and world access
- Decoded packet interception

Types come from Dragonfly rather than Minecraft identifiers or public transport DTOs. Registries
and codecs stay behind the private ABI. [`KitchenSink`](examples/plugins/kitchen-sink) compiles
against the full exposed API and is the practical coverage reference.

## Lifetimes and transactions

Dragonfly ownership rules still apply across the C# boundary:

- `Player.Context` inherits `World.Context`, which inherits `World.Tx`.
- Callback transactions and values borrowed through them expire when the callback returns.
- Packet objects expire when their packet callback returns.
- Lazy player and entity iterators own transaction-scoped values. Do not collect or retain yielded
  objects; retain `Player.H()` or another stable entity handle instead.
- Keep world-owner callbacks short. Waiting on work scheduled to the same owner can deadlock.
- Use `Player.ChangeWorld` for players; generic entity removal does not complete player session
  transfer safely.

`World.Do` and `DoAfter` schedule work on a world's owner with a fresh transaction. Their returned
`World.Task` supports completion, errors, callbacks, cancellation, and waiting. `World.Tx.Defer` and
`DeferErr` run FIFO after the current callback and before its parent task completes.

## Examples

- [`lifecycle-logger`](examples/plugins/lifecycle-logger): minimal lifecycle plugin
- [`chat-filter`](examples/plugins/chat-filter): mutable chat callback
- [`movement-guard`](examples/plugins/movement-guard): cancellable movement callback
- [`vanilla-commands`](examples/plugins/vanilla-commands): reflected command layout
- [`kitchen-sink`](examples/plugins/kitchen-sink): complete API coverage and runtime exercises

See [examples/plugins/README.md](examples/plugins/README.md) for build and staging details.

## Development

```shell
make generate        # Regenerate C# API from Dragonfly's Go AST.
make check-generated # Verify generated files are current.
make build           # Build native plugins and server.
make test            # Build and run C# and Go checks.
make clean           # Remove generated build outputs.
```

Set `DOTNET_RID` to override the inferred host runtime identifier.
