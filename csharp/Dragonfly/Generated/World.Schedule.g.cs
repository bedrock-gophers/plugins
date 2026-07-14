// Code generated from Dragonfly server/world/task.go Go AST. DO NOT EDIT.
#nullable enable
using System;
using System.Threading;
using System.Threading.Tasks;

namespace Dragonfly;

public sealed partial class World
{
    public Task Do(Action<Tx> callback) =>
        PluginBridge.Host.ScheduleWorld(this, TimeSpan.Zero, callback);

    public Task DoAfter(TimeSpan delay, Action<Tx> callback) =>
        PluginBridge.Host.ScheduleWorld(this, delay, callback);

    public sealed class Task
    {
        private readonly ulong _callbackId;
        private readonly object _lock = new();
        private readonly TaskCompletionSource _done = new(TaskCreationOptions.RunContinuationsAsynchronously);
        private Exception? _callbackError;
        private Exception? _error;
        private bool _completed;

        internal Task(ulong callbackId) => _callbackId = callbackId;

        internal void CallbackFailed(Exception error)
        {
            ArgumentNullException.ThrowIfNull(error);
            lock (_lock)
            {
                if (!_completed) _callbackError ??= error;
            }
        }

        internal void Complete(uint result)
        {
            lock (_lock)
            {
                if (_completed) return;
                _error = _callbackError ?? ResultError(result);
                _completed = true;
                _done.TrySetResult();
            }
        }

        public System.Threading.Tasks.Task Done() => _done.Task;

        public Exception? Err()
        {
            lock (_lock) return _completed ? _error : null;
        }

        public Exception? Wait(CancellationToken cancellationToken = default)
        {
            try
            {
                _done.Task.WaitAsync(cancellationToken).GetAwaiter().GetResult();
                return Err();
            }
            catch (OperationCanceledException error)
            {
                return error;
            }
        }

        public void OnDone(Action<Exception?> callback)
        {
            ArgumentNullException.ThrowIfNull(callback);
            _ = _done.Task.ContinueWith(
                completed => { _ = System.Threading.Tasks.Task.Run(() => callback(Err())); },
                CancellationToken.None,
                TaskContinuationOptions.None,
                TaskScheduler.Default);
        }

        public bool Cancel()
        {
            lock (_lock)
            {
                if (_completed) return false;
            }
            if (!PluginBridge.Host.CancelWorldTask(_callbackId)) return false;
            Complete(1);
            return true;
        }

        private static Exception? ResultError(uint result) => result switch
        {
            0 => null,
            1 => new TaskCanceledException("world task was cancelled"),
            2 => new InvalidOperationException("world closed before task ran"),
            _ => new InvalidOperationException("world task failed"),
        };
    }
}
