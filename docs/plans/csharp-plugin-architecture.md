# C# plugin architecture

## Direction

- C# NativeAOT is the only plugin language.
- The public namespace mirrors Dragonfly's packages and exported types as closely as C# permits.
- Plugins subclass `Plugin`; generated build plumbing supplies the native entry point and project-name identity.
- The Go host owns Dragonfly and exposes a private flat C ABI. Plugins never use ABI types.
- Code generation reads the pinned Dragonfly Go source with `go/ast` and emits C#; there is no second public API schema.

## Shape

```text
Dragonfly Go API -> Go AST generator -> C# Dragonfly API
                                         |
                                         v
plugin source -> NativeAOT .so -> private C ABI -> Go host -> Dragonfly
```

The ABI is transport, not the API. C# names, interfaces, constructors, and behavior should come from Dragonfly. Hand-written code is limited to marshalling semantics that cannot be inferred from Go types.

## Order

1. NativeAOT loading and `OnEnable`/`OnDisable`.
2. `player.Handler` events. Movement, chat, food loss, jump, teleport, sprint/sneak toggles, punch-air, and quit are implemented.
3. Player methods and commands. Command interfaces and the implemented `Player` method surface are generated from Dragonfly's Go AST. C# uses `Cmd.New`/`Cmd.Register`, one `Cmd.Runnable` per overload, and reflected public fields as Dragonfly uses reflected Go struct fields. Supported command fields include subcommands, native enums, dynamic `Cmd.Enum` values, players, vectors, optional values, and `Cmd.Varargs`. The generator roots runnable fields and field types for NativeAOT; runnable types use `internal` visibility and require no linker annotations. Bedrock-facing subcommands and enum/player suggestions are always lowercase. `Player.SetGameMode(World.GameMode)` and Dragonfly's four built-in game-mode values are available. The text slice includes `Message`, `SendPopup`, `SendTip`, `SendJukeboxPopup`, `SetNameTag`, and `Disconnect`. `Messagef` remains absent until Go `fmt.Sprintf` semantics can be preserved honestly.
4. Worlds, items, blocks, forms, entities, particles, and sounds.
5. Convert practice-core and expand parity tests against Dragonfly.

Each slice removes the replaced legacy implementation. Unsupported API remains absent rather than gaining a parallel abstraction.

`examples/plugins/kitchen-sink` must use every exposed API. Its NativeAOT build is the compile-time parity canary for handlers, commands, worlds, items, blocks, forms, entities, particles, and sounds as those slices land.

`examples/plugins/vanilla-commands` keeps its plugin entry tiny and one command per file. It currently exercises `/gamemode`, `/help`, `/ping`, and `/position`, and expands with each gameplay parity slice.

Practice remains out of the parity loop until the framework API needed by practice exists. Feature work lands and is tested in this repository first.
