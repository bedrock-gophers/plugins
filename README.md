# bedrock-gophers/plugins

Native multi-language plugin runtime for [df-mc/dragonfly](https://github.com/df-mc/dragonfly). Rust is the first supported plugin language.

Current status: architecture spike. One generated movement event travels through Go, the native Rust runtime, and a dynamically loaded Rust plugin.

## Build and test

Requirements:

- Go 1.26+
- Rust 1.96+
- C compiler and `dlopen` support

```shell
make test
make benchmark
```

Build and run owned server:

```shell
make build
./build/bedrock-gophers -config build/server.toml
```

First run creates configuration. Framework loads native runtime and plugins, creates Dragonfly, installs world/player handlers, owns accept loop, and closes cleanly on `SIGINT`/`SIGTERM`.

Regenerate ABI files after changing `schema/`:

```shell
make generate
```

## Rust plugin example

```rust
use dragonfly_plugin::{PlayerMoveEvent, Plugin, plugin};

#[derive(Default)]
struct MovementGuard;

#[plugin]
impl Plugin for MovementGuard {
    fn on_move(&self, event: &mut PlayerMoveEvent<'_>) {
        if event.new_position().y < 0.0 {
            event.cancel();
        }
    }
}
```

Events continue by default. Cancellation is monotonic; no `allow()` API exists.
Plugin identity defaults to Cargo's package name; handler code does not repeat it.

See [native plugin architecture](docs/plans/rust-plugin-architecture.md).

## Examples

- [Movement guard](examples/rust/movement-guard): cancels movement below Y=0.
- [Chat filter](examples/rust/chat-filter): replaces text and cancels a blocked message.

Both examples compile as native plugin libraries through `make build-native`.
