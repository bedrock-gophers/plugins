package framework

import (
	"errors"
	"sync"
	"testing"

	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/world"
)

type scheduledWorldRuntime struct {
	host.WorldRuntime
	players *host.Players
	fail    bool

	mu         sync.Mutex
	executions int
	drops      int
	plugin     uint64
	callback   uint64
	invocation native.InvocationID
	live       bool
}

func (*scheduledWorldRuntime) Subscriptions() uint64 { return 0 }

func (r *scheduledWorldRuntime) HandleWorldScheduled(plugin, callback uint64, invocation native.InvocationID, execute bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugin, r.callback = plugin, callback
	if !execute {
		r.drops++
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
	if !manager.ScheduleWorld(id, 41, 73) {
		t.Fatal("ScheduleWorld rejected a live world")
	}
	manager.DrainScheduled()
	runtime.mu.Lock()
	executions, drops := runtime.executions, runtime.drops
	plugin, callback, invocation, live := runtime.plugin, runtime.callback, runtime.invocation, runtime.live
	runtime.mu.Unlock()
	if executions != 1 || drops != 0 || plugin != 41 || callback != 73 || invocation == 0 || !live {
		t.Fatalf("scheduled callback = executions %d drops %d plugin %d callback %d invocation %d live %v", executions, drops, plugin, callback, invocation, live)
	}
	if _, ok := players.InvocationTx(invocation); ok {
		t.Fatal("scheduled transaction escaped its callback")
	}
}

func TestWorldScheduleDropsFailedCallbackAndStopsAdmission(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	runtime := &scheduledWorldRuntime{players: players, fail: true}
	manager.attachRuntime(runtime)
	if _, err := manager.Create("example:scheduled", world.Config{Synchronous: true}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	id, _ := manager.WorldByName(0, "example:scheduled")
	if !manager.ScheduleWorld(id, 5, 9) {
		t.Fatal("ScheduleWorld rejected a live world")
	}
	manager.DrainScheduled()
	runtime.mu.Lock()
	executions, drops := runtime.executions, runtime.drops
	runtime.mu.Unlock()
	if executions != 1 || drops != 1 {
		t.Fatalf("failed callback = executions %d drops %d, want 1, 1", executions, drops)
	}
	manager.StopScheduling()
	if manager.ScheduleWorld(id, 5, 10) {
		t.Fatal("ScheduleWorld accepted after StopScheduling")
	}
	if manager.ScheduleWorld(native.WorldID(999), 5, 11) {
		t.Fatal("ScheduleWorld accepted an unknown world")
	}
}

func TestWorldScheduleDropsTaskRejectedByClosedWorld(t *testing.T) {
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
	if manager.ScheduleWorld(id, 7, 11) {
		t.Fatal("ScheduleWorld accepted a closed world")
	}
	manager.StopScheduling()
	manager.DrainScheduled()
	runtime.mu.Lock()
	executions, drops := runtime.executions, runtime.drops
	runtime.mu.Unlock()
	if executions != 0 || drops != 1 {
		t.Fatalf("closed-world callback = executions %d drops %d, want 0, 1", executions, drops)
	}
}
