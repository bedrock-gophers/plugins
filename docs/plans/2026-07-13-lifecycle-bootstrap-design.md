# Lifecycle and core-world bootstrap design

## Goal

Make plugin startup transactional, plugin shutdown joined, and Dragonfly core-world
policy configurable before the server starts. This unblocks practice-core owning its
application and FFA worlds without practice-specific Go code.

## Plugin lifecycle

`Plugin::on_enable` returns `PluginResult`, whose error is an owned, thread-safe
Rust error. Host transport status values stay private.

Plugin ABI v4 gives the enable callback a caller-owned bounded error buffer. The
macro catches panics and writes a fixed panic message; ordinary errors write their
display text. The Rust runtime prefixes the plugin package name, rolls back the
failing plugin first, then already-enabled plugins in reverse order, and propagates
the safe message through the C bridge to Go. Host ABI v18 is unchanged.

Dragonfly normally binds listener factories inside `server.Config.New`, before
plugin enable. The framework replaces them with gated wrappers, then opens every
real listener synchronously after enable and command registration but before
`Server.Listen`. A partial bind failure closes opened listeners in reverse and
returns from `Run`; enable failure never binds a socket.

## Task ownership

`TaskGroup` owns every thread it starts. A shared `CancellationToken` supports
polling and blocking waits. Shutdown atomically stops admission, cancels, joins all
workers, and is idempotent under concurrent callers. Worker panics are contained.
Drop performs the same joined shutdown; detached plugin work is not supported.

The native host rejects host calls before enable and after final disable. Runtime
admission separately rejects/drains commands and ordinary events before
`on_disable`. Host/world access and entity callbacks remain available through
custom-world close. Final disable rejects/drains entity callbacks, destroys leaked
entity instances, and only then deactivates the host. Runtime destruction performs
both phases first, so no plugin code can execute after its library unloads.

## Shutdown order

For a started server:

1. close Dragonfly, rejecting connections and synchronously draining player/world callbacks;
2. begin plugin disable: reject/drain ordinary callbacks, run `on_disable`, and
   cancel/join plugin tasks while custom worlds and host calls remain available;
3. close plugin-owned custom worlds while entity save/destroy callbacks remain admitted;
4. finish plugin disable: reject/drain entity callbacks and deactivate host access;
5. destroy the Rust runtime and unload plugin libraries.

For an unstarted server, begin-disable runs first, then custom and core worlds
close while entity callbacks remain admitted, then final disable runs. Plugin
disable may therefore inspect or unload its worlds. Repeated cleanup is harmless.

## Core-world policy

`[worlds.core]` contains typed policies for read-only state, random ticks, time,
and weather. Provider configuration remains Dragonfly's `[dragonfly.World]`; it is
the single source for the shared core provider path.

Defaults preserve ordinary Dragonfly behaviour: writable, three random ticks per
subchunk, existing time, and existing weather. Invalid enum or cross-field pairs
fail config loading before native runtime creation.

Before `server.Config.New`, the framework applies read-only/random-tick scalars and
locks the shared provider settings to apply time/weather. A missing provider is
replaced with a `world.NopProvider` carrying stable settings. The one shared core
provider remains Dragonfly-owned and is never registered as a custom-world
provider.

Practice config uses read-only worlds, disabled random ticks, fixed time 6000, and
clear weather. Its persistent spawn provider remains `data/spawn`.

## Required proof

- exact plugin error survives Rust, runtime, C, and Go boundaries;
- panicking or failing enable never starts listening and rolls back in reverse;
- disable is once-only and no callback/host call executes after deactivation;
- task cancellation, admission race, panic containment, concurrent shutdown, and Drop all join;
- plugin disable runs while custom worlds remain usable;
- custom worlds close before runtime unload;
- core policy parses, validates, preserves defaults, and mutates all shared settings before server creation;
- generated ABI files, Rust tests/clippy, Go race tests, native builds, and examples all pass.
