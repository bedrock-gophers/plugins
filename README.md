# Minimal bedrock-gophers server

Consumer-only example. This branch contains no bedrock-gophers framework/runtime source.

```shell
make run
```

Examples cover every current SDK area: lifecycle, join/quit, hurt/heal, food loss/death, block break/place/pick, sign edits, lectern page turns, item use, sprint/sneak toggles, jump/teleport, experience gain, punch-air, held-slot changes, sleep, movement, mutable chat, typed command contexts, subcommands, enum parameters, and player latency. Each plugin owns its Cargo manifest, lockfile, source, and target. Build fetches pinned bedrock-gophers revision, compiles runtime and plugins, then starts server.
