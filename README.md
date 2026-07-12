# Minimal bedrock-gophers server

Consumer-only example. This branch contains no bedrock-gophers framework/runtime source.

```shell
make run
```

Examples cover every current SDK area: lifecycle, join/quit, player messaging/tips/popups/titles, teleport/move/velocity/rotation actions, hurt/heal, food loss/death, block break/place/pick, sign edits, lectern page turns, item use/consume/release/damage/drop and item use on blocks, item-stack views, sprint/sneak toggles, jump/teleport events, experience gain, punch-air, held-slot changes, sleep, mutable chat, typed command contexts, subcommands, enum parameters, and player latency. Each plugin owns its Cargo manifest, lockfile, source, and target. Build fetches pinned bedrock-gophers revision, compiles runtime and plugins, then starts server.
