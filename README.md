# Minimal bedrock-gophers server

Consumer-only example. This branch contains no bedrock-gophers framework/runtime source.

```shell
make run
```

Examples cover every current SDK area: lifecycle, join/quit, player messaging/tips/popups/titles, teleport/move/velocity/rotation actions, hurt/heal, food loss/death, block break/place/pick, sign edits, lectern page turns, item use/consume/release/damage/drop and item use on blocks, item-stack views, sprint/sneak toggles, jump/teleport events, experience gain, punch-air, held-slot changes, sleep, mutable chat, typed command contexts, subcommands, enum parameters, player latency, managed worlds, typed block state, stable entity handles and their worlds, typed built-in/projectile spawning, all typed Dragonfly particles and sounds, and cancellable attack-entity events. Each plugin owns its Cargo manifest, lockfile, source, and target. Build fetches the pinned bedrock-gophers revision, compiles runtime and plugins, then starts the server.
