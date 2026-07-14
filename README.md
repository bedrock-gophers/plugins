# Minimal bedrock-gophers server

Strict consumer-only example. This branch contains no framework, runtime, SDK, or plugin source.

```shell
make run
```

The build fetches the exact framework revision pinned in `go.mod`, compiles its C# NativeAOT runtime and lifecycle example, stages both shared libraries, then starts the Go server.
