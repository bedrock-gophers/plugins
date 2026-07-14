package framework

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/world"
)

type scheduledWorldRuntime struct {
	host.WorldRuntime
	players *host.Players
	fail    bool

	mu          sync.Mutex
	executions  int
	completions int
	plugin      uint64
	callback    uint64
	invocation  native.InvocationID
	result      native.WorldTaskResult
	live        bool
}

func (*scheduledWorldRuntime) Subscriptions() uint64 { return 0 }

func (r *scheduledWorldRuntime) HandleWorldScheduled(plugin, callback uint64, invocation native.InvocationID, phase native.WorldTaskPhase, result native.WorldTaskResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugin, r.callback = plugin, callback
	if phase == native.WorldTaskComplete {
		r.completions++
		r.result = result
		return nil
	}
	r.executions++
	r.invocation = invocation
	_, r.live = r.players.InvocationTx(invocation)
	if r.fail {
		return errors.New("callback failed")
	}
	return nil
}

func TestWorldScheduleRunsOnceWithBorrowedTransaction(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	runtime := &scheduledWorldRuntime{players: players}
	manager.attachRuntime(runtime)
	if _, err := manager.Create("example:scheduled", world.Config{Synchronous: true}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	id, _ := manager.WorldByName(0, "example:scheduled")
	if !manager.ScheduleWorld(id, 41, 73, 0) {
		t.Fatal("ScheduleWorld rejected a live world")
	}
	manager.DrainScheduled()
	runtime.mu.Lock()
	executions, completions := runtime.executions, runtime.completions
	plugin, callback, invocation, result, live := runtime.plugin, runtime.callback, runtime.invocation, runtime.result, runtime.live
	runtime.mu.Unlock()
	if executions != 1 || completions != 1 || result != native.WorldTaskSuccess || plugin != 41 || callback != 73 || invocation == 0 || !live {
		t.Fatalf("scheduled callback = executions %d completions %d result %d plugin %d callback %d invocation %d live %v", executions, completions, result, plugin, callback, invocation, live)
	}
	if _, ok := players.InvocationTx(invocation); ok {
		t.Fatal("scheduled transaction escaped its callback")
	}
}

func TestWorldScheduleReportsPanickedCallbackAndStopsAdmission(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	runtime := &scheduledWorldRuntime{players: players, fail: true}
	manager.attachRuntime(runtime)
	if _, err := manager.Create("example:scheduled", world.Config{Synchronous: true}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	id, _ := manager.WorldByName(0, "example:scheduled")
	if !manager.ScheduleWorld(id, 5, 9, 0) {
		t.Fatal("ScheduleWorld rejected a live world")
	}
	manager.DrainScheduled()
	runtime.mu.Lock()
	executions, completions, result := runtime.executions, runtime.completions, runtime.result
	runtime.mu.Unlock()
	if executions != 1 || completions != 1 || result != native.WorldTaskPanicked {
		t.Fatalf("failed callback = executions %d completions %d result %d", executions, completions, result)
	}
	manager.StopScheduling()
	if manager.ScheduleWorld(id, 5, 10, 0) {
		t.Fatal("ScheduleWorld accepted after StopScheduling")
	}
	if manager.ScheduleWorld(native.WorldID(999), 5, 11, 0) {
		t.Fatal("ScheduleWorld accepted an unknown world")
	}
}

func TestWorldScheduleReturnsCompletedTaskForClosedWorld(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	runtime := &scheduledWorldRuntime{players: players}
	manager.attachRuntime(runtime)
	w, err := manager.Create("example:closed", world.Config{Synchronous: true})
	if err != nil {
		t.Fatal(err)
	}
	id, _ := manager.WorldByName(0, "example:closed")
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if !manager.ScheduleWorld(id, 7, 11, 0) {
		t.Fatal("ScheduleWorld rejected a known closed world")
	}
	manager.StopScheduling()
	manager.DrainScheduled()
	runtime.mu.Lock()
	executions, completions, result := runtime.executions, runtime.completions, runtime.result
	runtime.mu.Unlock()
	if executions != 0 || completions != 1 || result != native.WorldTaskWorldClosed {
		t.Fatalf("closed-world callback = executions %d completions %d result %d", executions, completions, result)
	}
}

func TestWorldScheduleCancelStopsDelayedTask(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	runtime := &scheduledWorldRuntime{players: players}
	manager.attachRuntime(runtime)
	if _, err := manager.Create("example:delayed", world.Config{Synchronous: true}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	id, _ := manager.WorldByName(0, "example:delayed")
	if !manager.ScheduleWorld(id, 8, 12, int64(time.Hour)) {
		t.Fatal("ScheduleWorld rejected delayed task")
	}
	if cancelled, found := manager.CancelWorldTask(8, 12); !found || !cancelled {
		t.Fatalf("CancelWorldTask() = %v, %v", cancelled, found)
	}
	manager.DrainScheduled()
	runtime.mu.Lock()
	executions, completions, result := runtime.executions, runtime.completions, runtime.result
	runtime.mu.Unlock()
	if executions != 0 || completions != 1 || result != native.WorldTaskCancelled {
		t.Fatalf("cancelled callback = executions %d completions %d result %d", executions, completions, result)
	}
	if cancelled, found := manager.CancelWorldTask(8, 12); found || cancelled {
		t.Fatalf("completed CancelWorldTask() = %v, %v", cancelled, found)
	}
}
