package native

import (
	"sync"
	"sync/atomic"
)

const maxInventoryMenusPerHost = 128

var (
	inventoryMenuSequence  atomic.Uint64
	inventoryMenuMu        sync.Mutex
	inventoryMenuCond      = sync.NewCond(&inventoryMenuMu)
	inventoryMenus         = map[uint64]*inventoryMenuRegistration{}
	activeInventoryMenus   = map[inventoryMenuPlayerKey]*inventoryMenuRegistration{}
	inventoryMenuHostState = map[uint64]*inventoryMenuState{}
)

type inventoryMenuPlayerKey struct {
	host   uint64
	player PlayerID
}

type inventoryMenuRegistration struct {
	id     uint64
	host   uint64
	player PlayerID
	submit func(InvocationID, PlayerSnapshot, ItemStack) bool
	close  func(InvocationID, PlayerSnapshot) bool
	drop   func()

	inflight        int
	terminal        func()
	terminalStarted bool
}

type inventoryMenuState struct {
	closing, draining bool
	inflight, count   int
}

type inventoryMenuCallback struct {
	callback func()
	state    *inventoryMenuState
}

type inventoryMenuDiscarder interface {
	DiscardPlayerInventoryMenus()
}

func registerInventoryMenu(
	host uint64,
	player PlayerID,
	submit func(InvocationID, PlayerSnapshot, ItemStack) bool,
	close func(InvocationID, PlayerSnapshot) bool,
	drop func(),
) (uint64, bool) {
	inventoryMenuMu.Lock()
	state := inventoryMenuHostState[host]
	if state == nil {
		state = &inventoryMenuState{}
		inventoryMenuHostState[host] = state
	}
	effectiveCount := state.count
	if old := activeInventoryMenus[inventoryMenuPlayerKey{host: host, player: player}]; old != nil {
		effectiveCount--
	}
	if state.closing || state.draining || effectiveCount >= maxInventoryMenusPerHost {
		inventoryMenuMu.Unlock()
		return 0, false
	}
	id := inventoryMenuSequence.Add(1)
	registration := &inventoryMenuRegistration{
		id: id, host: host, player: player, submit: submit, close: close, drop: drop,
	}
	inventoryMenus[id] = registration
	state.count++
	inventoryMenuMu.Unlock()
	return id, true
}

// activateInventoryMenu replaces the previous player menu only after the host
// accepted the new menu. A failed send therefore leaves the old menu usable.
func activateInventoryMenu(id uint64) bool {
	inventoryMenuMu.Lock()
	registration := inventoryMenus[id]
	if registration == nil {
		inventoryMenuMu.Unlock()
		return false
	}
	state := inventoryMenuHostState[registration.host]
	if state == nil || state.draining {
		inventoryMenuMu.Unlock()
		return false
	}
	key := inventoryMenuPlayerKey{host: registration.host, player: registration.player}
	old := activeInventoryMenus[key]
	activeInventoryMenus[key] = registration
	var callback *inventoryMenuCallback
	if old != nil && old != registration {
		callback = detachInventoryMenuLocked(old, old.drop)
	}
	inventoryMenuMu.Unlock()
	runInventoryMenuCallback(callback)
	return true
}

// SubmitPlayerInventoryMenu invokes a non-terminal menu callback. It may be
// called repeatedly until the menu is closed, replaced, cancelled, or drained.
func SubmitPlayerInventoryMenu(id uint64, invocation InvocationID, submitter PlayerSnapshot, item ItemStack) bool {
	inventoryMenuMu.Lock()
	registration := inventoryMenus[id]
	if registration == nil {
		inventoryMenuMu.Unlock()
		return true
	}
	if submitter.Player != registration.player {
		inventoryMenuMu.Unlock()
		return false
	}
	state := inventoryMenuHostState[registration.host]
	if state == nil || state.draining {
		inventoryMenuMu.Unlock()
		return true
	}
	registration.inflight++
	state.inflight++
	inventoryMenuMu.Unlock()

	defer finishInventoryMenuSubmission(registration)
	return registration.submit(invocation, submitter, item)
}

// ClosePlayerInventoryMenu terminates a menu through its close callback.
func ClosePlayerInventoryMenu(id uint64, invocation InvocationID, player PlayerSnapshot) bool {
	inventoryMenuMu.Lock()
	registration := inventoryMenus[id]
	if registration == nil {
		inventoryMenuMu.Unlock()
		return true
	}
	if player.Player != registration.player {
		inventoryMenuMu.Unlock()
		return false
	}
	state := inventoryMenuHostState[registration.host]
	if state == nil || state.draining {
		inventoryMenuMu.Unlock()
		return true
	}
	callback := detachInventoryMenuLocked(registration, func() { registration.close(invocation, player) })
	inventoryMenuMu.Unlock()
	runInventoryMenuCallback(callback)
	return true
}

func CancelPlayerInventoryMenu(id uint64) {
	inventoryMenuMu.Lock()
	registration := inventoryMenus[id]
	var callback *inventoryMenuCallback
	if registration != nil {
		callback = detachInventoryMenuLocked(registration, registration.drop)
	}
	inventoryMenuMu.Unlock()
	runInventoryMenuCallback(callback)
}

func CancelPlayerInventoryMenus(host uint64, player PlayerID) {
	inventoryMenuMu.Lock()
	var callbacks []*inventoryMenuCallback
	for _, registration := range inventoryMenus {
		if registration.host == host && registration.player == player {
			if callback := detachInventoryMenuLocked(registration, registration.drop); callback != nil {
				callbacks = append(callbacks, callback)
			}
		}
	}
	inventoryMenuMu.Unlock()
	for _, callback := range callbacks {
		runInventoryMenuCallback(callback)
	}
}

func finishInventoryMenuSubmission(registration *inventoryMenuRegistration) {
	inventoryMenuMu.Lock()
	state := inventoryMenuHostState[registration.host]
	registration.inflight--
	state.inflight--
	callback := scheduleInventoryMenuTerminalLocked(registration)
	inventoryMenuCond.Broadcast()
	inventoryMenuMu.Unlock()
	runInventoryMenuCallback(callback)
}

func detachInventoryMenuLocked(registration *inventoryMenuRegistration, terminal func()) *inventoryMenuCallback {
	if inventoryMenus[registration.id] == registration {
		delete(inventoryMenus, registration.id)
		state := inventoryMenuHostState[registration.host]
		state.count--
	}
	key := inventoryMenuPlayerKey{host: registration.host, player: registration.player}
	if activeInventoryMenus[key] == registration {
		delete(activeInventoryMenus, key)
	}
	if registration.terminal == nil {
		registration.terminal = terminal
	}
	return scheduleInventoryMenuTerminalLocked(registration)
}

func scheduleInventoryMenuTerminalLocked(registration *inventoryMenuRegistration) *inventoryMenuCallback {
	if registration.terminal == nil || registration.terminalStarted || registration.inflight != 0 {
		return nil
	}
	registration.terminalStarted = true
	state := inventoryMenuHostState[registration.host]
	state.inflight++
	return &inventoryMenuCallback{callback: registration.terminal, state: state}
}

func runInventoryMenuCallback(callback *inventoryMenuCallback) {
	if callback == nil {
		return
	}
	defer func() {
		inventoryMenuMu.Lock()
		callback.state.inflight--
		inventoryMenuCond.Broadcast()
		inventoryMenuMu.Unlock()
	}()
	callback.callback()
}

func drainHostInventoryMenus(host uint64, closing bool) {
	inventoryMenuMu.Lock()
	state := inventoryMenuHostState[host]
	if state == nil {
		state = &inventoryMenuState{}
		inventoryMenuHostState[host] = state
	}
	state.draining = true
	if closing {
		state.closing = true
	}
	inventoryMenuMu.Unlock()

	discardHostInventoryMenus(host)

	inventoryMenuMu.Lock()
	callbacks := make([]*inventoryMenuCallback, 0, state.count)
	for _, registration := range inventoryMenus {
		if registration.host == host {
			if callback := detachInventoryMenuLocked(registration, registration.drop); callback != nil {
				callbacks = append(callbacks, callback)
			}
		}
	}
	inventoryMenuMu.Unlock()
	for _, callback := range callbacks {
		runInventoryMenuCallback(callback)
	}

	inventoryMenuMu.Lock()
	for state.inflight != 0 {
		inventoryMenuCond.Wait()
	}
	if !closing {
		state.draining = false
	}
	inventoryMenuMu.Unlock()
}

func discardHostInventoryMenus(host uint64) {
	value, ok := hosts.Load(host)
	if !ok {
		return
	}
	registered := value.(*registeredHost)
	if discarder, ok := registered.host.(inventoryMenuDiscarder); ok {
		discarder.DiscardPlayerInventoryMenus()
	}
}

func activateHostInventoryMenus(host uint64) {
	inventoryMenuMu.Lock()
	state := inventoryMenuHostState[host]
	if state == nil {
		inventoryMenuHostState[host] = &inventoryMenuState{}
	} else {
		state.closing = false
		state.draining = false
	}
	inventoryMenuMu.Unlock()
}

// CancelPlayerInventoryMenus terminates this runtime's menus for one player.
func (r *Runtime) CancelPlayerInventoryMenus(player PlayerID) {
	if r != nil {
		CancelPlayerInventoryMenus(r.hostContext, player)
	}
}
