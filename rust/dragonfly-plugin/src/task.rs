//! Scoped background work that is cancelled and joined during plugin shutdown.

use std::{
    error::Error,
    fmt,
    panic::{AssertUnwindSafe, catch_unwind},
    sync::{Arc, Condvar, Mutex, MutexGuard},
    thread::{self, JoinHandle},
    time::Duration,
};

#[cfg(test)]
#[path = "task_test.rs"]
mod tests;

/// Cooperative cancellation shared by every worker in a [`TaskGroup`].
#[derive(Clone, Debug)]
pub struct CancellationToken {
    inner: Arc<Cancellation>,
}

#[derive(Debug, Default)]
struct Cancellation {
    cancelled: Mutex<bool>,
    changed: Condvar,
}

impl CancellationToken {
    /// Reports whether shutdown cancellation has been requested.
    pub fn is_cancelled(&self) -> bool {
        *lock(&self.inner.cancelled)
    }

    /// Blocks until shutdown cancellation has been requested.
    pub fn wait(&self) {
        let mut cancelled = lock(&self.inner.cancelled);
        while !*cancelled {
            cancelled = self
                .inner
                .changed
                .wait(cancelled)
                .unwrap_or_else(|poison| poison.into_inner());
        }
    }

    /// Waits up to `timeout` for cancellation, returning whether it occurred.
    pub fn wait_timeout(&self, timeout: Duration) -> bool {
        let cancelled = lock(&self.inner.cancelled);
        if *cancelled {
            return true;
        }
        let (cancelled, _) = self
            .inner
            .changed
            .wait_timeout_while(cancelled, timeout, |cancelled| !*cancelled)
            .unwrap_or_else(|poison| poison.into_inner());
        *cancelled
    }

    fn cancel(&self) {
        let mut cancelled = lock(&self.inner.cancelled);
        if *cancelled {
            return;
        }
        *cancelled = true;
        self.inner.changed.notify_all();
    }
}

/// The externally observable lifecycle of a [`TaskGroup`].
#[derive(Clone, Copy, Debug, Default, Eq, PartialEq)]
pub enum TaskGroupState {
    #[default]
    Open,
    Closing,
    Closed,
}

/// A user-safe failure to start scoped work.
#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub enum SpawnError {
    /// Shutdown has begun, so the group no longer accepts work.
    ShuttingDown,
    /// The operating system could not create a worker thread.
    Unavailable,
}

impl fmt::Display for SpawnError {
    fn fmt(&self, formatter: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::ShuttingDown => formatter.write_str("task group is shutting down"),
            Self::Unavailable => formatter.write_str("background task is unavailable"),
        }
    }
}

impl Error for SpawnError {}

#[derive(Default)]
struct Lifecycle {
    state: TaskGroupState,
    workers: Vec<JoinHandle<()>>,
}

/// Owns plugin background workers until they have all stopped.
///
/// Workers should block on or periodically inspect their [`CancellationToken`].
/// Calling [`TaskGroup::shutdown`] requests cancellation and joins every accepted
/// worker. Shutdown is idempotent, and concurrent callers wait for the first
/// caller to finish joining.
pub struct TaskGroup {
    lifecycle: Mutex<Lifecycle>,
    changed: Condvar,
    cancellation: CancellationToken,
}

impl Default for TaskGroup {
    fn default() -> Self {
        Self::new()
    }
}

impl TaskGroup {
    /// Creates an open task group.
    pub fn new() -> Self {
        Self {
            lifecycle: Mutex::new(Lifecycle::default()),
            changed: Condvar::new(),
            cancellation: CancellationToken {
                inner: Arc::new(Cancellation::default()),
            },
        }
    }

    /// Returns the cancellation token associated with this group.
    pub fn cancellation_token(&self) -> CancellationToken {
        self.cancellation.clone()
    }

    /// Reports the current group lifecycle.
    pub fn state(&self) -> TaskGroupState {
        lock(&self.lifecycle).state
    }

    /// Starts a worker while the group is open.
    ///
    /// The worker receives the same cooperative cancellation token as every
    /// other worker. A worker panic is contained to that worker and does not
    /// escape shutdown or plugin unload.
    pub fn spawn<F>(&self, worker: F) -> Result<(), SpawnError>
    where
        F: FnOnce(CancellationToken) + Send + 'static,
    {
        let mut lifecycle = lock(&self.lifecycle);
        if lifecycle.state != TaskGroupState::Open {
            return Err(SpawnError::ShuttingDown);
        }

        let cancellation = self.cancellation.clone();
        let handle = thread::Builder::new()
            .name("dragonfly-plugin-task".to_owned())
            .spawn(move || {
                let _ = catch_unwind(AssertUnwindSafe(|| worker(cancellation)));
            })
            .map_err(|_| SpawnError::Unavailable)?;
        lifecycle.workers.push(handle);
        Ok(())
    }

    /// Rejects new work, requests cancellation, and joins every accepted worker.
    ///
    /// The owner of the group must call this from outside one of its workers.
    pub fn shutdown(&self) {
        let workers = {
            let mut lifecycle = lock(&self.lifecycle);
            loop {
                match lifecycle.state {
                    TaskGroupState::Open => {
                        lifecycle.state = TaskGroupState::Closing;
                        self.cancellation.cancel();
                        break std::mem::take(&mut lifecycle.workers);
                    }
                    TaskGroupState::Closing => {
                        lifecycle = self
                            .changed
                            .wait(lifecycle)
                            .unwrap_or_else(|poison| poison.into_inner());
                    }
                    TaskGroupState::Closed => return,
                }
            }
        };

        for worker in workers {
            let _ = worker.join();
        }

        let mut lifecycle = lock(&self.lifecycle);
        lifecycle.state = TaskGroupState::Closed;
        self.changed.notify_all();
    }
}

impl Drop for TaskGroup {
    fn drop(&mut self) {
        self.shutdown();
    }
}

fn lock<T>(mutex: &Mutex<T>) -> MutexGuard<'_, T> {
    mutex.lock().unwrap_or_else(|poison| poison.into_inner())
}
