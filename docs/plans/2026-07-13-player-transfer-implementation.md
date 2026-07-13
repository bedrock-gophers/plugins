# Stable Player Transfer Implementation Plan

**Goal:** Add generation-safe, session-preserving player transfer between
managed Dragonfly worlds through host ABI v17 and a minimal Rust API.

**Architecture:** `WorldManager` validates a live invocation and stable player
handle, then defers cross-world removal until the plugin callback ends. It
moves that exact handle with `AddEntityAt`, protects both worlds with lifecycle
leases, and restores to source only on `ErrWorldClosed`. Same-world calls use
ordinary Dragonfly teleport. Tests live in separate focused files.

**Tech stack:** Go 1.26, Dragonfly v0.11, cgo/C ABI, Rust 2024 workspace.

## Constraints

- Host ABI v17 is 464 bytes; `world_open_spec` stays at 448 and
  `player_transfer` is appended at 456.
- Plugin ABI remains v3.
- Never retain `world.Tx` or transaction-scoped `*player.Player`.
- Retain the exact player `*world.EntityHandle` and existing player registry
  entry.
- Never block one world owner waiting for another.
- Invocation zero schedules through the stable handle; only stale nonzero
  invocations are rejected.
- Public Rust returns `()` and hides native status.
- Do not use or cherry-pick the transferable-entities experiment.
- Add production behavior test-first and commit coherent local milestones.

## Milestone 1: Documentation and host choreography

- [x] Audit Dragonfly `Tx.Defer`, `World.Do`, `RemoveEntity`, `AddEntityAt`,
  `Player.Teleport`, entity-handle scheduling, and natural change-world tick.
- [x] Record exact public, ABI, lease, rollback, identity, event, and failure
  contracts in the design document.
- [ ] Commit the documentation milestone.

## Milestone 2: Framework transfer core

Files:

- Create `framework/player_transfer.go`.
- Create `framework/player_transfer_test.go`.
- Modify `internal/host/players.go`.
- Modify `internal/host/players_test.go` or another focused existing host test
  file only for narrow handle lookup coverage.

Red-green sequence:

1. [x] Write focused tests for generation-checked player handle lookup.
2. [x] Add `Players.Handle(PlayerID)` and public departure cleanup needed only for
   rollback.
3. [x] Write same-world tests proving ordinary teleport semantics and invalid
   request rejection.
4. [x] Write cross-world tests proving exact handle, exact `PlayerID`, exact
   position, a source mutation performed after `TransferPlayer` in the same
   callback, and an invocation-zero transfer followed by a handle-based
   mutation in the destination.
5. [x] Implement non-blocking source/destination lifecycle leases and deferred handoff.
6. [x] Write and pass tests for destination `ErrWorldClosed` restoration, stale
   generation, player quit/closed handle, unload waiting, and idempotent lease
   release under `-race`.
7. [x] Write a first-destination-tick test proving exactly one natural
   `PlayerChangeWorld` event and no event on rollback.
8. [x] Add a public Dragonfly session fixture proving loader switch and
   connection teardown while the player handle is worldless.
9. [x] Add terminal double-world-close eviction and logging coverage.
10. Commit the framework milestone.

The integration fixture uses Dragonfly's real player entity type, world owners,
handlers, ticks, and public `session.Config.New`/`SetHandle`/`Spawn` APIs with a
fake network connection. No private session setup is copied.

## Milestone 3: Host ABI v17

Files:

- Modify `cmd/abi-gen/main.go` and regenerate outputs.
- Modify `rust/dragonfly-plugin-sys/src/lib.rs` layout tests.
- Modify `internal/native/bridge.c`.
- Create `internal/native/player_transfer_exports.go`.
- Create `internal/native/player_transfer_exports_test.go` or focused native
  integration coverage.
- Modify `internal/native/host.go`.
- Modify runtime and macro host pointer types.

Red-green sequence:

1. [x] Change layout tests first to require `DfHostApiV17`, size 464,
   `world_open_spec` offset 448, and `player_transfer` offset 456.
2. [x] Append `DfHostPlayerTransferFn` in the generator and regenerate C/Rust.
3. [x] Rename all generated/runtime/macro host API references from v16 to v17,
   retaining the exact v16 prefix.
4. [x] Add `Host.TransferPlayer` plus no-op implementation.
5. [x] Add a separate cgo exporter validating values and forwarding to
   the framework; it returns only private ABI status.
6. [x] Wire bridge extern, wrapper, initializer, and static assertions.
7. [x] Verify generator check, layout tests, native tests, runtime tests, and
   macros; commit the ABI milestone.

## Milestone 4: Rust API and integration

Files:

- Modify `rust/dragonfly-plugin/src/lib.rs`.
- Create `rust/dragonfly-plugin/src/player_transfer_test.rs`.
- Update an appropriate example plugin without combining unrelated commands.
- Update architecture/parity documentation and this plan's status.

Red-green sequence:

1. [x] Add a Rust test host recording player/world/position and require
   `Player::transfer(World, Vec3) -> ()`.
2. [x] Implement the method as one private host call, ignoring unavailable host or
   non-OK status.
3. [x] Add native/runtime integration coverage proving all values cross the ABI.
4. [x] Add a small transfer subcommand to the focused world-command example
   structure.
5. [x] Record the public session-backed fixture and its loader/teardown evidence.
6. Commit the SDK/example/documentation milestone.

## Final verification gate

Run from a clean worktree:

```bash
go run ./cmd/abi-gen -root . -check
go test -race ./...
cargo fmt --all --check
cargo clippy --workspace --all-targets -- -D warnings
cargo test --workspace
make native-integration
```

Also run focused transfer tests repeatedly under the race detector and verify
the example plugin build. Report all exact commands and results before handing
the local commits to the parent task.
