# Stable player transfer design

Status: Approved for implementation

## Goal

Add a dedicated player world-transfer operation that preserves Dragonfly's
player entity handle, session, framework player ID, and handler state. Rust
plugins use the same shape as Dragonfly world movement without seeing ABI
status values:

```rust
player.transfer(destination, Vec3::new(8.5, 72.0, -3.5));
```

This is deliberately separate from generic transferable entities. A player is
session-backed and must not pass through a detached-token API.

## Public contract

`Player::transfer(world: World, position: Vec3)` returns `()`.

- Unknown/stale player IDs, unknown/unloading world IDs, stale nonzero
  invocations, and non-finite coordinates are ignored. Invocation zero is the
  valid off-callback/scoped-task mode.
- A transfer to the player's current world calls ordinary
  `player.Player.Teleport`, including Dragonfly's cancellable teleport event.
- A transfer to another world moves the exact existing
  `*world.EntityHandle`; it does not create a player, replace its registry
  entry, or expose a detached state.
- Cross-world placement uses `Tx.AddEntityAt` with the exact requested
  position.
- The Rust API does not expose the private `DfStatus` used at the native
  boundary.

The operation is accepted synchronously but a cross-world handoff completes
after the current Dragonfly callback. Later operations in the same plugin
callback, such as healing or changing a kit, therefore still execute against
the live source player and source transaction. Off-callback calls schedule
through the stable entity handle; subsequent handle-based mutations wait while
it is worldless and follow it into the destination.

## Transaction choreography

The framework validates the exact player generation, destination world, finite
position, and current source player before changing anything. A nonzero
invocation must resolve to its exact live transaction. Invocation zero
schedules through `EntityHandle.Do` and discovers the current source owner in
that callback.

For a same-world destination it invokes `Player.Teleport` immediately.

For a cross-world destination it acquires lifecycle read leases for both
managed worlds and registers this work with the current source transaction:

1. `sourceTx.Defer` runs after all plugin code in the callback and receives a
   fresh source transaction.
2. The deferred callback resolves the same player handle in that transaction
   and calls `sourceTx.RemoveEntity`.
3. It submits `destination.World.Do` and calls
   `destinationTx.AddEntityAt(handle, position)` there.
4. Successful destination completion releases both lifecycle leases.
5. If and only if destination submission/completion reports
   `world.ErrWorldClosed`, the exact handle is re-added to the source world.
   A synchronous destination rejection restores through the still-valid
   deferred source transaction. An asynchronous rejection schedules a fresh
   source `World.Do` transaction.
6. Source restoration clears the recorded departure so a failed transfer
   cannot contaminate a later change-world event.

No transaction or transaction-scoped `*player.Player` is retained. The only
long-lived Dragonfly value is the stable `*world.EntityHandle`, which is the
documented cross-world ownership object.

The source and destination lifecycle leases remain held until the handle is
attached to one side or a terminal closed-handle/source-close path is reached.
Framework unload therefore waits for an accepted handoff rather than closing a
world between removal and addition. Calls never wait for another world owner.

## Identity and lifecycle

`host.Players` adds one narrow generation-checked handle lookup. It returns the
existing entity handle only when the complete `PlayerID` is registered and the
handle is open. Transfer orchestration never unregisters or re-registers the
player entry.

World despawn/spawn handlers temporarily deactivate and reactivate the generic
entity view of the same handle. They already preserve its entity generation.
The player registry entry and `PlayerID` remain unchanged for the entire
handoff.

Quit, close, and stale-generation races fail closed:

- a stale ID can never obtain the handle;
- a player absent from the fresh source/handle transaction is not removed;
- a closed handle is never added;
- destination `ErrWorldClosed` restores to source when source remains open;
- no rollback is attempted after destination code ran and returned another
  error, because the player may already be live in the destination;
- every lifecycle lease has a single idempotent release path.

## Change-world events

The transfer emits no synthetic plugin events. Dragonfly stores the player's
previous world on the same player data and naturally calls
`HandleChangeWorld` on the first destination tick. The source world handler
records the source framework handle at despawn, allowing the existing player
handler to report it even if unload begins afterward.

Therefore a successful transfer produces exactly one natural
`PlayerChangeWorld` event. A failed transfer restored to its source produces
none.

## ABI

Host ABI becomes v17. Plugin ABI remains v3. The complete v16 prefix is
unchanged:

- `world_open_spec` remains at offset 448;
- one `DfHostPlayerTransferFn player_transfer` pointer is appended at offset
  456;
- `DfHostApiV17` is 464 bytes on supported 64-bit targets.

The function signature is:

```c
typedef DfStatus (*DfHostPlayerTransferFn)(
    uint64_t context,
    DfInvocationId invocation,
    DfPlayerId player,
    DfWorldId world,
    DfVec3 position);
```

`DF_STATUS_OK` means the request was validated and applied or scheduled. Other
statuses remain private to generated/runtime SDK code.

## Structure

- `framework/player_transfer.go` owns transfer orchestration and lifecycle
  leases.
- `framework/player_transfer_test.go` contains focused behavioral/race tests.
- `internal/host/players.go` only supplies exact handle lookup and departure
  cleanup.
- `internal/native/player_transfer_exports.go` contains the cgo exporter.
- generated ABI files, runtime, macros, and Rust SDK consistently use v17.

## Verification evidence

Tests must cover same-world teleport, cross-world exact handle/ID preservation,
exact destination position, post-request source operations, stale player and
world handles, non-finite positions, quit/close races, destination-close
rollback, unload waiting, and one natural change-world event.

A true network session fixture is only added if Dragonfly exposes one without
duplicating its private session bootstrap. Otherwise the integration test uses
a real `player.Type` entity, real world owners/transactions, real player/world
handlers, and first-destination tick behavior, and this limitation is recorded
in the implementation plan and final verification notes.
