# Minimal bedrock-gophers server

Consumer-only example. This branch contains no bedrock-gophers framework/runtime source.

```shell
make run
```

One Rust plugin lives under `plugins/movement-guard` and owns its Cargo manifest, lockfile, source, and target. Build fetches pinned bedrock-gophers revision, compiles runtime and plugin, then starts server.
