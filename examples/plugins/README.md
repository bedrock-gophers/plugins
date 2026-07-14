# Plugins

Each source plugin owns its `.csproj`. `make run` regenerates the C# API, publishes every source
project as NativeAOT, replaces staged example `.so` files, then starts the server.

The loader accepts compatible precompiled `.so` files, but `make run` is source-owned and cleans
its staged binaries. Stage a precompiled file after the build and start the server directly when
testing binary-only delivery.

`kitchen-sink` uses every exposed C# API and grows with each parity slice. Its `block` overload
demonstrates `World.Tx`, `Cube.Pos`, typed `Block.Sand`, and `World.SetOpts` without exposing block
identifiers or NBT.

`vanilla-commands` keeps each command in its own file and grows as more of Dragonfly's gameplay API is exposed.
