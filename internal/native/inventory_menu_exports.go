package native

/*
#include <stdlib.h>
#include "bridge.h"
*/
import "C"

import (
	"sync"
	"sync/atomic"
	"unicode/utf8"
	"unsafe"
)

const maxInventoryMenuTitleBytes = 4096

type inventoryMenuHost interface {
	SendPlayerInventoryMenu(InvocationID, PlayerID, PlayerInventoryMenu) bool
	ClosePlayerInventoryMenu(InvocationID, PlayerID) bool
}

type inventoryMenuRegistration struct {
	host   uint64
	player PlayerID
	click  func(InvocationID, PlayerSnapshot, uint32) bool
	close  func(InvocationID, PlayerSnapshot) bool
	drop   func()
	done   chan struct{}

	mu         sync.Mutex
	inflight   int
	terminated bool
	terminal   func()
}

var (
	inventoryMenuSequence atomic.Uint64
	inventoryMenuMu       sync.Mutex
	inventoryMenuRegistry = map[uint64]*inventoryMenuRegistration{}
)

//export bg_go_player_inventory_menu_send
func bg_go_player_inventory_menu_send(context C.uint64_t, invocation C.DfInvocationId, player C.DfPlayerId, view *C.DfInventoryMenuView) C.DfStatus {
	if view == nil || view.click == nil || view.close == nil || view.drop == nil {
		return C.DF_STATUS_ERROR
	}
	callbackContext, clickCallback, closeCallback, dropCallback := view.callback_context, view.click, view.close, view.drop
	drop := func() { C.bg_call_inventory_menu_drop(dropCallback, callbackContext) }
	titleBytes, ok := copySkinView(view.title)
	if !ok || len(titleBytes) == 0 || len(titleBytes) > maxInventoryMenuTitleBytes || !utf8.Valid(titleBytes) {
		drop()
		return C.DF_STATUS_ERROR
	}
	container := InventoryMenuContainer(view.container)
	size := inventoryMenuSize(container)
	if size == 0 || uint64(view.item_count) != uint64(size) || size != 0 && view.items == nil || view.update > 1 {
		drop()
		return C.DF_STATUS_ERROR
	}
	items := make([]ItemStack, size)
	for index := range items {
		item, valid := copyItemStackView(&unsafe.Slice(view.items, size)[index])
		if !valid {
			drop()
			return C.DF_STATUS_ERROR
		}
		items[index] = item
	}
	playerID := playerID(player)
	id, ok := registerInventoryMenu(uint64(context), playerID,
		func(callbackInvocation InvocationID, snapshot PlayerSnapshot, slot uint32) bool {
			return callInventoryMenuClick(clickCallback, callbackContext, callbackInvocation, snapshot, slot)
		},
		func(callbackInvocation InvocationID, snapshot PlayerSnapshot) bool {
			return callInventoryMenuClose(closeCallback, callbackContext, callbackInvocation, snapshot)
		}, drop)
	if !ok {
		drop()
		return C.DF_STATUS_ERROR
	}
	host, resolved := resolveHost(uint64(context))
	menuHost, available := host.(inventoryMenuHost)
	if !resolved || !available || !menuHost.SendPlayerInventoryMenu(InvocationID(invocation), playerID, PlayerInventoryMenu{
		ID: id, Title: string(titleBytes), Container: container, Items: items, Update: view.update != 0,
	}) {
		CancelPlayerInventoryMenu(id)
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}

//export bg_go_player_inventory_menu_close
func bg_go_player_inventory_menu_close(context C.uint64_t, invocation C.DfInvocationId, player C.DfPlayerId) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	menuHost, available := host.(inventoryMenuHost)
	if !ok || !available || !menuHost.ClosePlayerInventoryMenu(InvocationID(invocation), playerID(player)) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}

func inventoryMenuSize(container InventoryMenuContainer) int {
	switch container {
	case InventoryMenuChest, InventoryMenuBarrel, InventoryMenuEnderChest:
		return 27
	case InventoryMenuDoubleChest:
		return 54
	case InventoryMenuHopper:
		return 5
	case InventoryMenuDropper:
		return 9
	default:
		return 0
	}
}

func registerInventoryMenu(host uint64, player PlayerID, click func(InvocationID, PlayerSnapshot, uint32) bool, closeCallback func(InvocationID, PlayerSnapshot) bool, drop func()) (uint64, bool) {
	if _, ok := resolveHost(host); !ok {
		return 0, false
	}
	id := inventoryMenuSequence.Add(1)
	registration := &inventoryMenuRegistration{
		host: host, player: player, click: click, close: closeCallback, drop: drop, done: make(chan struct{}),
	}
	inventoryMenuMu.Lock()
	if len(inventoryMenuRegistry) >= maxFormsPerHost*4 {
		inventoryMenuMu.Unlock()
		return 0, false
	}
	inventoryMenuRegistry[id] = registration
	inventoryMenuMu.Unlock()
	return id, true
}

func CompletePlayerInventoryMenuClick(id uint64, invocation InvocationID, player PlayerSnapshot, slot uint32) bool {
	inventoryMenuMu.Lock()
	registration := inventoryMenuRegistry[id]
	inventoryMenuMu.Unlock()
	if registration == nil {
		return false
	}
	if player.Player != registration.player {
		CancelPlayerInventoryMenu(id)
		return false
	}
	registration.mu.Lock()
	if registration.terminated {
		registration.mu.Unlock()
		return false
	}
	registration.inflight++
	registration.mu.Unlock()
	accepted := registration.click(invocation, player, slot)
	finishInventoryMenuCallback(registration)
	return accepted
}

func CompletePlayerInventoryMenuClose(id uint64, invocation InvocationID, player PlayerSnapshot) bool {
	inventoryMenuMu.Lock()
	registration := inventoryMenuRegistry[id]
	if registration != nil {
		delete(inventoryMenuRegistry, id)
	}
	inventoryMenuMu.Unlock()
	if registration == nil {
		return false
	}
	if player.Player != registration.player {
		terminateInventoryMenu(registration, registration.drop)
		return false
	}
	return terminateInventoryMenu(registration, func() { registration.close(invocation, player) })
}

func CancelPlayerInventoryMenu(id uint64) {
	inventoryMenuMu.Lock()
	registration := inventoryMenuRegistry[id]
	if registration != nil {
		delete(inventoryMenuRegistry, id)
	}
	inventoryMenuMu.Unlock()
	if registration != nil {
		terminateInventoryMenu(registration, registration.drop)
	}
}

func CancelPlayerInventoryMenus(player PlayerID) {
	cancelInventoryMenus(func(registration *inventoryMenuRegistration) bool { return registration.player == player })
}

func drainHostInventoryMenus(host uint64) {
	registrations := cancelInventoryMenus(func(registration *inventoryMenuRegistration) bool { return registration.host == host })
	for _, registration := range registrations {
		<-registration.done
	}
}

func cancelInventoryMenus(match func(*inventoryMenuRegistration) bool) []*inventoryMenuRegistration {
	inventoryMenuMu.Lock()
	registrations := make([]*inventoryMenuRegistration, 0)
	for id, registration := range inventoryMenuRegistry {
		if match(registration) {
			delete(inventoryMenuRegistry, id)
			registrations = append(registrations, registration)
		}
	}
	inventoryMenuMu.Unlock()
	for _, registration := range registrations {
		terminateInventoryMenu(registration, registration.drop)
	}
	return registrations
}

func terminateInventoryMenu(registration *inventoryMenuRegistration, terminal func()) bool {
	registration.mu.Lock()
	if registration.terminated {
		registration.mu.Unlock()
		return false
	}
	registration.terminated = true
	if registration.inflight != 0 {
		registration.terminal = terminal
		registration.mu.Unlock()
		return true
	}
	registration.mu.Unlock()
	terminal()
	close(registration.done)
	return true
}

func finishInventoryMenuCallback(registration *inventoryMenuRegistration) {
	registration.mu.Lock()
	registration.inflight--
	terminal := registration.terminal
	if registration.inflight == 0 {
		registration.terminal = nil
	} else {
		terminal = nil
	}
	registration.mu.Unlock()
	if terminal != nil {
		terminal()
		close(registration.done)
	}
}

func callInventoryMenuClick(callback C.DfInventoryMenuClickFn, context unsafe.Pointer, invocation InvocationID, player PlayerSnapshot, slot uint32) bool {
	return withCPlayerSnapshot(player, func(snapshot *C.DfPlayerSnapshot) bool {
		return C.bg_call_inventory_menu_click(callback, context, C.DfInvocationId(invocation), snapshot, C.uint32_t(slot)) == C.DF_STATUS_OK
	})
}

func callInventoryMenuClose(callback C.DfInventoryMenuCloseFn, context unsafe.Pointer, invocation InvocationID, player PlayerSnapshot) bool {
	return withCPlayerSnapshot(player, func(snapshot *C.DfPlayerSnapshot) bool {
		return C.bg_call_inventory_menu_close(callback, context, C.DfInvocationId(invocation), snapshot) == C.DF_STATUS_OK
	})
}

func withCPlayerSnapshot(player PlayerSnapshot, callback func(*C.DfPlayerSnapshot) bool) bool {
	name := C.CBytes([]byte(player.Name))
	if name != nil {
		defer C.free(name)
	}
	snapshot := C.DfPlayerSnapshot{
		player:               cPlayerID(player.Player),
		name:                 C.DfStringView{data: (*C.uint8_t)(name), len: C.uint64_t(len(player.Name))},
		latency_milliseconds: C.uint64_t(player.LatencyMilliseconds),
		position:             C.DfVec3{x: C.double(player.Position.X), y: C.double(player.Position.Y), z: C.double(player.Position.Z)},
	}
	return callback(&snapshot)
}
