# Minimal bedrock-gophers server

Consumer-only example. This branch contains no bedrock-gophers framework/runtime source.

```shell
make run
```

Examples cover every current SDK area: lifecycle, movement, mutable chat, and commands with subcommands and enum parameters. Each plugin owns its Cargo manifest, lockfile, source, and target. Build fetches pinned bedrock-gophers revision, compiles runtime and plugins, then starts server.
