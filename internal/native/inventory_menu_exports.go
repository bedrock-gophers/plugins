package native

/*
#include "bridge.h"
*/
import "C"

import (
	"unicode/utf8"
	"unsafe"
)

const (
	maxInventoryMenuNameBytes = 4096
	maxInventoryMenuDataBytes = 64 << 20
)

//export bg_go_player_inventory_menu_send
func bg_go_player_inventory_menu_send(context C.uint64_t, invocation C.DfInvocationId, player C.DfPlayerId, update C.uint8_t, view *C.DfInventoryMenuView) C.DfStatus {
	if view == nil || update > 1 || view.reserved != 0 || view.submit == nil || view.close == nil || view.drop == nil {
		return C.DF_STATUS_ERROR
	}
	callbackContext, submitCallback, closeCallback, dropCallback := view.callback_context, view.submit, view.close, view.drop
	drop := func() { C.bg_call_inventory_menu_drop(dropCallback, callbackContext) }

	nameBytes, ok := copyNativeBytes(view.name, maxInventoryMenuNameBytes)
	container := InventoryMenuContainer(view.container)
	size, validContainer := inventoryMenuContainerSize(container)
	if !ok || !utf8.Valid(nameBytes) || !validContainer || uint64(view.item_count) != uint64(size) || size != 0 && view.items == nil {
		drop()
		return C.DF_STATUS_ERROR
	}
	items := make([]ItemStack, size)
	dataBytes := len(nameBytes)
	if size != 0 {
		for index := range items {
			item, valid := copyItemStackView(&unsafe.Slice(view.items, size)[index])
			if !valid {
				drop()
				return C.DF_STATUS_ERROR
			}
			dataBytes += inventoryMenuItemBytes(item)
			if dataBytes > maxInventoryMenuDataBytes {
				drop()
				return C.DF_STATUS_ERROR
			}
			items[index] = item
		}
	}

	ok = sendPlayerInventoryMenu(
		uint64(context), InvocationID(invocation), playerID(player),
		PlayerInventoryMenu{Name: string(nameBytes), Container: container, Items: items, Update: update != 0},
		func(callbackInvocation InvocationID, submitter PlayerSnapshot, item ItemStack) bool {
			arena := &nativeViewArena{}
			defer arena.release()
			snapshot, valid := nativeInventoryMenuPlayerSnapshot(submitter, arena)
			if !valid {
				return false
			}
			itemView, valid := nativeItemStackView(item, arena)
			if !valid {
				return false
			}
			return C.bg_call_inventory_menu_submit(
				submitCallback, callbackContext, C.DfInvocationId(callbackInvocation), &snapshot, &itemView,
			) == C.DF_STATUS_OK
		},
		func(callbackInvocation InvocationID, connected PlayerSnapshot) bool {
			arena := &nativeViewArena{}
			defer arena.release()
			snapshot, valid := nativeInventoryMenuPlayerSnapshot(connected, arena)
			if !valid {
				return false
			}
			return C.bg_call_inventory_menu_close(
				closeCallback, callbackContext, C.DfInvocationId(callbackInvocation), &snapshot,
			) == C.DF_STATUS_OK
		},
		drop,
	)
	if !ok {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}

func sendPlayerInventoryMenu(
	context uint64,
	invocation InvocationID,
	player PlayerID,
	menu PlayerInventoryMenu,
	submit func(InvocationID, PlayerSnapshot, ItemStack) bool,
	close func(InvocationID, PlayerSnapshot) bool,
	drop func(),
) bool {
	host, ok := resolveHost(context)
	if !ok {
		drop()
		return false
	}
	id, ok := registerInventoryMenu(context, player, submit, close, drop)
	if !ok {
		drop()
		return false
	}
	menu.ID = id
	if !host.SendPlayerInventoryMenu(invocation, player, menu) {
		host.DiscardPlayerInventoryMenu(player, id)
		CancelPlayerInventoryMenu(id)
		return false
	}
	if !activateInventoryMenu(id) {
		host.DiscardPlayerInventoryMenu(player, id)
		CancelPlayerInventoryMenu(id)
		return false
	}
	return true
}

//export bg_go_player_inventory_menu_close
func bg_go_player_inventory_menu_close(context C.uint64_t, invocation C.DfInvocationId, player C.DfPlayerId) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	if !ok || !host.ClosePlayerInventoryMenu(InvocationID(invocation), playerID(player)) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}

func inventoryMenuContainerSize(container InventoryMenuContainer) (int, bool) {
	switch container {
	case InventoryMenuChest, InventoryMenuBarrel, InventoryMenuEnderChest:
		return 27, true
	case InventoryMenuDoubleChest:
		return 54, true
	case InventoryMenuHopper:
		return 5, true
	case InventoryMenuDropper:
		return 9, true
	default:
		return 0, false
	}
}

func inventoryMenuItemBytes(item ItemStack) int {
	total := len(item.Identifier) + len(item.CustomName) + len(item.NBT) + len(item.ValuesNBT)
	for _, line := range item.Lore {
		total += len(line)
	}
	return total
}

func nativeInventoryMenuPlayerSnapshot(value PlayerSnapshot, arena *nativeViewArena) (C.DfPlayerSnapshot, bool) {
	name, ok := arena.stringView([]byte(value.Name))
	if !ok {
		return C.DfPlayerSnapshot{}, false
	}
	return C.DfPlayerSnapshot{
		player:               cPlayerID(value.Player),
		name:                 name,
		latency_milliseconds: C.uint64_t(value.LatencyMilliseconds),
		position: C.DfVec3{
			x: C.double(value.Position.X), y: C.double(value.Position.Y), z: C.double(value.Position.Z),
		},
	}, true
}
