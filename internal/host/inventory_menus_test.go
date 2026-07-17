package host

import (
	"testing"
	"time"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestClassifyInventoryMenuClose(t *testing.T) {
	tests := []struct {
		name                        string
		requestWindow, activeWindow byte
		want                        inventoryMenuCloseMode
	}{
		{name: "active window", requestWindow: 7, activeWindow: 7, want: inventoryMenuCloseAcknowledge},
		{name: "chat transition", requestWindow: 0xff, activeWindow: 7, want: inventoryMenuCloseIgnore},
		{name: "stale window", requestWindow: 6, activeWindow: 7, want: inventoryMenuCloseIgnore},
		{name: "player inventory", requestWindow: 0, activeWindow: 7, want: inventoryMenuCloseIgnore},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := classifyInventoryMenuClose(test.requestWindow, test.activeWindow, false); got != test.want {
				t.Fatalf("classifyInventoryMenuClose(%d, %d, false) = %d, want %d", test.requestWindow, test.activeWindow, got, test.want)
			}
		})
	}
	if got := classifyInventoryMenuClose(0xff, 7, true); got != inventoryMenuCloseAcknowledge {
		t.Fatalf("0xff close while awaiting acknowledgement = %d, want %d", got, inventoryMenuCloseAcknowledge)
	}
}

func TestPendingInventoryMenuID(t *testing.T) {
	if pendingInventoryMenuID(nil) != 0 {
		t.Fatal("nil pending menu has a non-zero ID")
	}
	pending := &pendingInventoryMenu{value: native.PlayerInventoryMenu{ID: 42}}
	if got := pendingInventoryMenuID(pending); got != 42 {
		t.Fatalf("pending menu ID = %d, want 42", got)
	}
}

func TestExpirePendingInventoryMenuClearsTransition(t *testing.T) {
	id := native.PlayerID{Generation: 1}
	menu := &activeInventoryMenu{pending: &pendingInventoryMenu{value: native.PlayerInventoryMenu{ID: 42}}}
	menus := NewInventoryMenus(nil)
	menus.active[id] = menu

	menus.expirePendingInventoryMenu(id, menu)
	if menus.active[id] != nil {
		t.Fatal("timed out transition remained active")
	}
	if menu.pending != nil {
		t.Fatal("timed out transition retained its pending menu")
	}
}

func TestInventoryMenuCloseCancelsPendingOpen(t *testing.T) {
	menu := &activeInventoryMenu{windowID: 7}
	menu.timer = time.AfterFunc(time.Hour, func() {})
	menus := NewInventoryMenus(nil)
	menus.closeLocked(session.Nop, nil, menu, false, false)
	if menu.timer != nil {
		t.Fatal("pending open timer was retained after close")
	}
	if menu.opened {
		t.Fatal("pending menu was marked open during close")
	}
}

func TestInventoryMenuReusesWindow(t *testing.T) {
	tests := []struct {
		name                         string
		update, previousOpened, same bool
		want                         bool
	}{
		{name: "open same container", previousOpened: true, same: true, want: true},
		{name: "update pending same container", update: true, same: true, want: true},
		{name: "open different container", previousOpened: true},
		{name: "open pending same container", same: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := inventoryMenuReusesWindow(test.update, test.previousOpened, test.same); got != test.want {
				t.Fatalf("inventoryMenuReusesWindow(%t, %t, %t) = %t, want %t", test.update, test.previousOpened, test.same, got, test.want)
			}
		})
	}
}

func TestInventoryMenuDefersReplacement(t *testing.T) {
	if !inventoryMenuDefersReplacement(true, true) {
		t.Fatal("mutating a menu from a click must be deferred until the invocation ends")
	}
	for _, test := range []struct {
		invocation, opened bool
	}{
		{opened: true},
		{invocation: true},
	} {
		if inventoryMenuDefersReplacement(test.invocation, test.opened) {
			t.Fatalf("unexpected deferral for invocation=%t opened=%t", test.invocation, test.opened)
		}
	}
}

func TestRewriteInventoryMenuRequest(t *testing.T) {
	menuSource := protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerBarrel},
		Slot:      10,
	}
	playerDestination := protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{ContainerID: protocol.ContainerInventory},
		Slot:      2,
	}
	request := &packet.ItemStackRequest{Requests: []protocol.ItemStackRequest{{
		Actions: []protocol.StackRequestAction{&protocol.SwapStackRequestAction{
			Source: menuSource, Destination: playerDestination,
		}},
	}}}

	rewriteInventoryMenuRequest(request, 27)
	action := request.Requests[0].Actions[0].(*protocol.SwapStackRequestAction)
	if action.Source.Container.ContainerID != protocol.ContainerLevelEntity {
		t.Fatalf("menu source container = %d, want %d", action.Source.Container.ContainerID, protocol.ContainerLevelEntity)
	}
	if action.Destination.Container.ContainerID != protocol.ContainerInventory {
		t.Fatalf("player destination container was rewritten to %d", action.Destination.Container.ContainerID)
	}
}

func TestRewriteInventoryMenuSlotLeavesUnrelatedSlots(t *testing.T) {
	for _, slot := range []protocol.StackRequestSlotInfo{
		{Container: protocol.FullContainerName{ContainerID: protocol.ContainerCursor}, Slot: 0},
		{Container: protocol.FullContainerName{ContainerID: protocol.ContainerBarrel}, Slot: 27},
	} {
		original := slot.Container.ContainerID
		rewriteInventoryMenuSlot(&slot, 27)
		if slot.Container.ContainerID != original {
			t.Fatalf("container %d slot %d was rewritten to %d", original, slot.Slot, slot.Container.ContainerID)
		}
	}
}

func TestInventoryMenuContainersUseFakeBlocks(t *testing.T) {
	tests := []struct {
		container native.InventoryMenuContainer
		window    byte
		actor     string
		size      int
		double    bool
	}{
		{native.InventoryMenuChest, protocol.ContainerTypeContainer, "Chest", 27, false},
		{native.InventoryMenuDoubleChest, protocol.ContainerTypeContainer, "Chest", 54, true},
		{native.InventoryMenuHopper, protocol.ContainerTypeHopper, "Hopper", 5, false},
		{native.InventoryMenuDropper, protocol.ContainerTypeDropper, "Dropper", 9, false},
		{native.InventoryMenuBarrel, protocol.ContainerTypeContainer, "Barrel", 27, false},
		{native.InventoryMenuEnderChest, protocol.ContainerTypeContainer, "EnderChest", 27, false},
	}
	for _, test := range tests {
		info, ok := resolveInventoryMenuContainer(test.container)
		if !ok || info.block == nil {
			t.Fatalf("container %d info=%#v ok=%t", test.container, info, ok)
		}
		if info.packetContainer != test.window || info.blockActorID != test.actor || info.size != test.size || info.double != test.double {
			t.Fatalf("container %d info=%#v", test.container, info)
		}
	}
}
