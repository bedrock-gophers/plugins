package framework

import "log/slog"

// runCleanup keeps teardown ordering explicit. A started Dragonfly server is
// closed first because Close synchronously rejects and drains player callbacks.
// Pending world tasks are cancelled and running tasks drain before on_disable.
// Plugin disable is then two-phase: on_disable runs with custom worlds and host
// calls available, custom worlds drain entity callbacks, then admission closes.
type runCleanup struct {
	log            *slog.Logger
	started        bool
	closeStarted   func() error
	stopScheduling func()
	beginPlugins   func()
	closeCustom    func() error
	drainDetached  func()
	finishPlugins  func()
	closeUnstarted func()
	drainScheduled func()
	closeRuntime   func()
}

func (cleanup *runCleanup) close() {
	if cleanup.started {
		if err := cleanup.closeStarted(); err != nil {
			cleanup.log.Error("close Dragonfly server", "error", err)
		}
	}
	cleanup.stopScheduling()
	cleanup.drainScheduled()
	cleanup.beginPlugins()
	if err := cleanup.closeCustom(); err != nil {
		cleanup.log.Error("close custom worlds", "error", err)
	}
	cleanup.drainDetached()
	if !cleanup.started {
		cleanup.closeUnstarted()
	}
	cleanup.finishPlugins()
	cleanup.closeRuntime()
}
