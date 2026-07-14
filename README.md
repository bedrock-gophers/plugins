# Minimal bedrock-gophers server

Small consumer example. Plugin source lives in `plugins/`; framework and generated ABI code do not.

```shell
make run
```

The build fetches the exact framework revision pinned in `go.mod`, compiles the C# NativeAOT
runtime, compiles every local plugin to `.so`, then starts the Go server.

Every C# project under `plugins/*/` is discovered automatically. Adding another plugin only
requires a `.csproj` and its C# source.

## Docker

Docker Compose:

```shell
docker compose up --build
```

The multi-stage image builds Go, the C# NativeAOT runtime, and every example plugin for the
machine's architecture. Bedrock listens on UDP port `19132`; server data is kept in the
`bedrock-gophers-minimal-data` Docker volume.

Plain Docker:

```shell
docker build -t bedrock-gophers-minimal .
docker run --rm --init -p 19132:19132/udp \
  -v bedrock-gophers-minimal-data:/app/.data bedrock-gophers-minimal
```

`make docker-run` is the short form when Make is installed.
