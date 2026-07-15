# Minimal bedrock-gophers server

Small consumer example. Plugin source lives in `plugins/`; framework and generated ABI code do not.

```shell
make run
```

Every `make run` fetches the exact framework revision pinned in `go.mod`, rebuilds the C#
NativeAOT runtime, rebuilds every local plugin to `.so`, replaces staged plugin binaries, then
starts the Go server. It does not reuse an older plugin build.

Every C# project under `plugins/*/` is discovered automatically. Adding another plugin only
requires a `.csproj` and its C# source.

The kitchen-sink plugin includes `/kitchen block`, demonstrating generated `World.Tx` range,
loaded-block, lazy block-search, height/light, typed block/liquid APIs, and registry-backed
`World.BlockByName` lookups for bool, byte, int, and string state properties. It probes an
empty liquid result, places/reads/removes typed water, then leaves matching water present and
schedules its update after 250 ms with `ScheduleBlockUpdate` and `TimeSpan`.
Its `/kitchen crop` overload round-trips promoted Dragonfly crop state through
`Block.WheatSeeds(Growth: 7)`.
Its `/kitchen biome` overload uses the generated vanilla biome types and exercises biome,
temperature, rain, snow, and thunder transaction queries before restoring the original biome.
`/kitchen tick` reads the transaction owner's current tick, not the world's day-time.
`/kitchen handle` exercises AST-generated `EntityHandle`, `RemoveEntity`, and `AddEntityAt`
through the command transaction while preserving handle identity.
`/kitchen particle` emits every typed Dragonfly particle and exercises all typed note instruments.
`/kitchen game-mode` exercises registered lookup, player reads, and a custom game mode without
leaving the player changed.
`/kitchen presentation` sends an AST-generated typed `Title` with subtitle, action text, and
durations.
`/kitchen cooldown` exercises typed `Player.HasCooldown` and `SetCooldown` with a diamond sword.
`/kitchen scoreboard` builds and sends an AST-generated typed scoreboard.
`/kitchen item` builds a typed diamond sword, round-trips inventory, armour, and held-item APIs,
then restores every changed slot.
`/kitchen form` covers reflected menu, custom, and modal forms with every typed element;
`/kitchen raw-form` demonstrates the open plugin-defined form interface.
Compatible precompiled `.so` plugins remain supported by the loader; because the source build
clears `plugins/*.so`, stage binary-only plugins after `make build` and start `.build/bin/server`
directly.

## Docker

Docker Compose:

```shell
docker compose up --build
```

The multi-stage image builds Go, the C# NativeAOT runtime, and every example plugin for the
machine's architecture, then bakes those `.so` files into the image. Plugin source changes require
`docker compose up --build`; there is no plugin hot reload or mounted plugin directory. Bedrock
listens on UDP port `19132`; server data is kept in the `bedrock-gophers-minimal-data` Docker
volume.

Plain Docker:

```shell
docker build -t bedrock-gophers-minimal .
docker run --rm --init -p 19132:19132/udp \
  -v bedrock-gophers-minimal-data:/app/.data bedrock-gophers-minimal
```

`make docker-run` is the short form when Make is installed.
