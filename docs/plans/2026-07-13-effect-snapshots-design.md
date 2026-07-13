# Player effect snapshots and clear-all

Status: Implemented; automated gates pass, Bedrock acceptance pending

## Goal

Expose Dragonfly's current lasting player effects to Rust plugins and add one
owner-transaction host operation for clearing them. The public API is deliberately small:

```rust
let effects = player.effects();
player.clear_effects();
```

Host transport failures remain private. `effects()` fails closed to an empty
vector and `clear_effects()` returns `()` like the existing player mutators.

## Dragonfly mapping

`player.Player.Effects()` is read inside the player's current or freshly
scheduled transaction. Every returned value must be a registered
`effect.LastingType`. `effect.ID` is the authoritative reverse mapping, so
custom registered lasting effects work without a generated Go switch or a new
schema adapter row. IDs must fit the ABI's signed 32-bit effect ID.

Snapshots preserve level, remaining duration, finite/ambient/infinite mode,
and hidden-particle state. Instant effects are applied immediately by
Dragonfly and are never valid snapshot entries. Finite duration zero remains
valid because an effect may have reached zero after the current tick and be
removed on the next expiry pass.

Clear-all snapshots `Player.Effects()` and calls `RemoveEffect` for every type
inside one player-owner mutation. It therefore clears built-in and registered
custom effects without one FFI round trip per effect. This is one owner pass,
not a rollback transaction: an effect `End` hook may add another effect.
Before the first player tick, Dragonfly's first `RemoveEffect` flushes its
pending initial list. Valid initial instant effects are therefore applied once
before the lasting effects are removed. The host rejects invalid initial level
or duration values before triggering that flush. A side-effect-free discard of
Dragonfly's private initial list is not available through its public API.

## Host ABI v18

Plugin ABI remains v3. Host ABI v18 retains the exact v17 prefix:

- `world_open_spec` remains at offset 448;
- `player_transfer` remains at offset 456;
- `player_effects` is appended at offset 464;
- `player_effects_clear` is appended at offset 472;
- `DfHostApiV18` is 480 bytes on supported 64-bit targets.

The fixed-size existing `DfEffectView` is reused for snapshot entries:

```c
typedef struct {
    DfEffectView *data;
    uint64_t len;
    uint64_t capacity;
} DfEffectBuffer;

typedef DfStatus (*DfHostPlayerEffectsFn)(
    uint64_t context,
    DfInvocationId invocation,
    DfPlayerId player,
    DfEffectBuffer *output);

typedef DfStatus (*DfHostPlayerEffectsClearFn)(
    uint64_t context,
    DfInvocationId invocation,
    DfPlayerId player);
```

`DfEffectBuffer` is 24 bytes with `data`, `len`, and `capacity` at offsets 0,
8, and 16. A snapshot contains at most 256 effects. Rust performs at most
three capacity attempts, beginning with a zero-capacity sizing probe.

The host converts the complete Go snapshot to a temporary C-view slice before
touching caller memory. It sets the required length and returns an error on
insufficient capacity without partial writes. A buffer with non-zero capacity
requires a non-null data pointer. Snapshot output accepts only positive levels,
potency exactly `1.0`, finite timed/ambient modes, or infinite mode with zero duration,
and a particle-hidden byte of zero or one. Durations use millisecond transport
granularity. Sub-millisecond positive values floor to zero; negative remaining
durations between Dragonfly expiry passes clamp to zero.

Each successful host call is one coherent transaction snapshot. No retained
snapshot registry is needed because entries contain no nested variable-length
data. Ordering is unspecified, matching Dragonfly's effect manager map.

## Safe Rust API

`Player::effects() -> Vec<effect::Effect>` privately validates the complete
buffer before returning it. Malformed length, excessive count, invalid mode,
invalid potency, invalid level, invalid particle byte, or host rejection makes
the public result empty. No raw `DfStatus` or host error type is exposed.

Initial instant effects are filtered because Dragonfly has not applied them
yet and does not retain instant effects as active status. `effect::Effect` adds
read-only accessors: `type_id`, `level`, `duration`,
`ambient`, `infinite`, and `particles_hidden`. A custom returned ID can be
used with `effect::RegisteredLasting::new(id)`. Tick age is not transported;
the snapshot is status inspection, not a byte-for-byte reapplication state.

## Verification

Tests cover the exact ABI prefix and tail, built-in and signed custom IDs,
finite/ambient/infinite and hidden-particle values, empty and stale players,
the 256-entry bound, capacity retry, malformed buffer values, insufficient
capacity without partial writes, and clear-all. A native integration fixture
must exercise Rust through the runtime and C bridge into a Go host.

Host ABI v18 is reserved exclusively for this milestone. The obsolete
transferable-entities experiment must rebase on the completed v18 prefix and
append its independent fields as host ABI v19.

Implementation milestones:

- `421a9cd` records the bounded snapshot and clear-all design.
- `fc1711a` appends the v18 ABI and implements the Go/C host path.
- `5be79e8` exposes the safe Rust snapshot, accessors, and clear operation.
- `e13a216` adds runnable command examples and Rust-to-C-to-Go integration.

Automated evidence completed on 2026-07-13:

- `go test -race ./...`
- generated ABI/block checks
- `cargo fmt --all -- --check`
- `cargo clippy --workspace --all-targets -- -D warnings`
- `cargo test --workspace`
- focused native cross-language integration
- full `make test`, including every runnable example plugin
