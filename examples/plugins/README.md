# Plugins

Each source plugin owns its `.csproj`. `make run` regenerates the C# API, publishes every source
project as NativeAOT for the host platform, replaces staged example native libraries, then starts
the server. Windows, Linux, and macOS use `.dll`, `.so`, and `.dylib` files respectively.
It does not hot-reload an existing server process; stop that process and reconnect after rerunning
`make run` so Bedrock receives the new command metadata.

The loader accepts compatible precompiled native libraries, but `make run` is source-owned and
cleans its staged binaries. Stage a precompiled file after the build and start the server directly
when testing binary-only delivery. Set `DOTNET_RID` to override the inferred host RID.

`kitchen-sink` uses every exposed C# API and grows with each parity slice. Its `block` overload
demonstrates `World.Tx`, `Cube.Pos`, typed `Block.Sand`, and `World.SetOpts` without exposing block
identifiers or NBT. Its `crop` overload round-trips the promoted `Growth` field as a typed
`Block.WheatSeeds`. Its `item` overload round-trips typed finite-state items, books, fireworks,
and all 28 tiered armour states through player inventory while the private transport keeps item
identifiers and NBT out of plugin code, including dyed leather and armour trims.
Its `server` overload uses direct lazy server-wide player iteration, broadcasts inside each borrowed
player lifetime, and verifies stable UUID/name handle lookup.

`vanilla-commands` keeps each command in its own file and grows as more of Dragonfly's gameplay API is exposed.
