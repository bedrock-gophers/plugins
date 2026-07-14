package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func inspectWorldSchedule(path string) error {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return err
	}

	want := map[string]map[string]goSignature{
		"World": {
			"Do":      {Parameters: "func(*Tx)", Results: "*Task"},
			"DoAfter": {Parameters: "time.Duration, func(*Tx)", Results: "*Task"},
		},
		"Task": {
			"Done":   {Results: "<-chan struct{}"},
			"Err":    {Results: "error"},
			"Wait":   {Parameters: "context.Context", Results: "error"},
			"OnDone": {Parameters: "func(error)"},
			"Cancel": {Results: "bool"},
		},
	}
	found := map[string]map[string]goSignature{"World": {}, "Task": {}}
	hasTask := false
	for _, declaration := range file.Decls {
		switch declaration := declaration.(type) {
		case *ast.GenDecl:
			for _, specification := range declaration.Specs {
				typeSpec, ok := specification.(*ast.TypeSpec)
				if !ok || typeSpec.Name.Name != "Task" {
					continue
				}
				_, hasTask = typeSpec.Type.(*ast.StructType)
			}
		case *ast.FuncDecl:
			for receiver := range want {
				if pointerReceiver(declaration, receiver) {
					found[receiver][declaration.Name.Name] = worldScheduleSignature(declaration)
				}
			}
		}
	}
	if !hasTask {
		return fmt.Errorf("Dragonfly world has no Task struct")
	}
	for receiver, methods := range want {
		for name, signature := range methods {
			got, ok := found[receiver][name]
			if !ok {
				return fmt.Errorf("Dragonfly world.%s has no %s method", receiver, name)
			}
			if got != signature {
				return fmt.Errorf("Dragonfly world.%s.%s signature changed: %+v", receiver, name, got)
			}
		}
	}
	return nil
}

func worldScheduleSignature(function *ast.FuncDecl) goSignature {
	return goSignature{
		Parameters: formatWorldScheduleFields(function.Type.Params),
		Results:    formatWorldScheduleFields(function.Type.Results),
	}
}

func formatWorldScheduleFields(fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}
	var types []string
	for _, field := range fields.List {
		count := len(field.Names)
		if count == 0 {
			count = 1
		}
		for range count {
			types = append(types, formatWorldScheduleExpression(field.Type))
		}
	}
	return joinComma(types)
}

func formatWorldScheduleExpression(expression ast.Expr) string {
	switch value := expression.(type) {
	case *ast.ChanType:
		prefix := "chan "
		if value.Dir == ast.RECV {
			prefix = "<-chan "
		} else if value.Dir == ast.SEND {
			prefix = "chan<- "
		}
		return prefix + formatWorldScheduleExpression(value.Value)
	case *ast.StructType:
		if value.Fields == nil || len(value.Fields.List) == 0 {
			return "struct{}"
		}
	}
	return formatGoExpression(expression)
}

func generateWorldSchedule() []byte {
	var output bytes.Buffer
	output.WriteString("// Code generated from Dragonfly server/world/task.go Go AST. DO NOT EDIT.\n")
	output.WriteString("#nullable enable\nusing System;\nusing System.Threading;\nusing System.Threading.Tasks;\n\nnamespace Dragonfly;\n\n")
	output.WriteString(`public sealed partial class World
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
`)
	return output.Bytes()
}
