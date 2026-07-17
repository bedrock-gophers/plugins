package native

import (
	"sync/atomic"
	"testing"
)

type inventoryMenuHost struct {
	noopHost
	send       func(InvocationID, PlayerID, PlayerInventoryMenu) bool
	discardOne func(PlayerID, uint64) bool
	discardAll func()
}

func (h inventoryMenuHost) SendPlayerInventoryMenu(invocation InvocationID, player PlayerID, menu PlayerInventoryMenu) bool {
	return h.send(invocation, player, menu)
}

func (h inventoryMenuHost) DiscardPlayerInventoryMenus() {
	if h.discardAll != nil {
		h.discardAll()
	}
}

func (h inventoryMenuHost) DiscardPlayerInventoryMenu(player PlayerID, id uint64) bool {
	return h.discardOne != nil && h.discardOne(player, id)
}

func TestInventoryMenuSupportsRepeatedSubmitAndTerminalClose(t *testing.T) {
	var sent PlayerInventoryMenu
	host := registerHost(inventoryMenuHost{send: func(_ InvocationID, _ PlayerID, menu PlayerInventoryMenu) bool {
		sent = menu
		return true
	}})
	t.Cleanup(func() { unregisterHost(host) })
	player := PlayerID{Generation: 7}
	var submits, closes, drops atomic.Int32
	if !sendPlayerInventoryMenu(host, 3, player, PlayerInventoryMenu{},
		func(InvocationID, PlayerSnapshot, ItemStack) bool { submits.Add(1); return true },
		func(InvocationID, PlayerSnapshot) bool { closes.Add(1); return true },
		func() { drops.Add(1) }) {
		t.Fatal("inventory menu send rejected")
	}
	snapshot := PlayerSnapshot{Player: player, Name: "MenuPlayer"}
	if !SubmitPlayerInventoryMenu(sent.ID, 4, snapshot, ItemStack{Identifier: "minecraft:apple", Count: 1}) ||
		!SubmitPlayerInventoryMenu(sent.ID, 5, snapshot, ItemStack{Identifier: "minecraft:diamond", Count: 1}) {
		t.Fatal("inventory menu submission rejected")
	}
	if !ClosePlayerInventoryMenu(sent.ID, 6, snapshot) {
		t.Fatal("inventory menu close rejected")
	}
	CancelPlayerInventoryMenu(sent.ID)
	if submits.Load() != 2 || closes.Load() != 1 || drops.Load() != 0 {
		t.Fatalf("callbacks = submits %d closes %d drops %d, want 2/1/0", submits.Load(), closes.Load(), drops.Load())
	}
}

func TestInventoryMenuReplacementDropsPreviousContext(t *testing.T) {
	var sent []PlayerInventoryMenu
	host := registerHost(inventoryMenuHost{send: func(_ InvocationID, _ PlayerID, menu PlayerInventoryMenu) bool {
		sent = append(sent, menu)
		return true
	}})
	t.Cleanup(func() { unregisterHost(host) })
	player := PlayerID{Generation: 9}
	var firstDrops, secondDrops atomic.Int32
	callbacks := func(drop *atomic.Int32) (func(InvocationID, PlayerSnapshot, ItemStack) bool, func(InvocationID, PlayerSnapshot) bool, func()) {
		return func(InvocationID, PlayerSnapshot, ItemStack) bool { return true },
			func(InvocationID, PlayerSnapshot) bool { return true },
			func() { drop.Add(1) }
	}
	firstSubmit, firstClose, firstDrop := callbacks(&firstDrops)
	if !sendPlayerInventoryMenu(host, 0, player, PlayerInventoryMenu{}, firstSubmit, firstClose, firstDrop) {
		t.Fatal("first inventory menu send rejected")
	}
	secondSubmit, secondClose, secondDrop := callbacks(&secondDrops)
	if !sendPlayerInventoryMenu(host, 0, player, PlayerInventoryMenu{}, secondSubmit, secondClose, secondDrop) {
		t.Fatal("second inventory menu send rejected")
	}
	if firstDrops.Load() != 1 || secondDrops.Load() != 0 {
		t.Fatalf("replacement drops = first %d second %d, want 1/0", firstDrops.Load(), secondDrops.Load())
	}
	CancelPlayerInventoryMenus(host, player)
	if secondDrops.Load() != 1 {
		t.Fatalf("active menu drops = %d, want 1", secondDrops.Load())
	}
	if len(sent) != 2 || sent[0].ID == sent[1].ID {
		t.Fatalf("sent menu IDs = %#v", sent)
	}
}

func TestInventoryMenuFailureDropsContextExactlyOnce(t *testing.T) {
	var sent PlayerInventoryMenu
	host := registerHost(inventoryMenuHost{send: func(_ InvocationID, _ PlayerID, menu PlayerInventoryMenu) bool {
		sent = menu
		return false
	}})
	t.Cleanup(func() { unregisterHost(host) })
	var drops atomic.Int32
	if sendPlayerInventoryMenu(host, 0, PlayerID{Generation: 1}, PlayerInventoryMenu{},
		func(InvocationID, PlayerSnapshot, ItemStack) bool { return true },
		func(InvocationID, PlayerSnapshot) bool { return true },
		func() { drops.Add(1) }) {
		t.Fatal("failed inventory menu send reported success")
	}
	CancelPlayerInventoryMenu(sent.ID)
	if drops.Load() != 1 {
		t.Fatalf("drops = %d, want 1", drops.Load())
	}
}

func TestInventoryMenuFailedReplacementKeepsPreviousContext(t *testing.T) {
	var sent []PlayerInventoryMenu
	host := registerHost(inventoryMenuHost{send: func(_ InvocationID, _ PlayerID, menu PlayerInventoryMenu) bool {
		sent = append(sent, menu)
		return len(sent) == 1
	}})
	t.Cleanup(func() { unregisterHost(host) })
	player := PlayerID{Generation: 11}
	var firstSubmits, firstDrops, secondDrops atomic.Int32
	if !sendPlayerInventoryMenu(host, 0, player, PlayerInventoryMenu{},
		func(InvocationID, PlayerSnapshot, ItemStack) bool { firstSubmits.Add(1); return true },
		func(InvocationID, PlayerSnapshot) bool { return true },
		func() { firstDrops.Add(1) }) {
		t.Fatal("first inventory menu send rejected")
	}
	if sendPlayerInventoryMenu(host, 0, player, PlayerInventoryMenu{},
		func(InvocationID, PlayerSnapshot, ItemStack) bool { return true },
		func(InvocationID, PlayerSnapshot) bool { return true },
		func() { secondDrops.Add(1) }) {
		t.Fatal("failed replacement reported success")
	}
	if len(sent) != 2 || firstDrops.Load() != 0 || secondDrops.Load() != 1 {
		t.Fatalf("failed replacement = sends %d first drops %d second drops %d, want 2/0/1",
			len(sent), firstDrops.Load(), secondDrops.Load())
	}
	if !SubmitPlayerInventoryMenu(sent[0].ID, 1, PlayerSnapshot{Player: player}, ItemStack{}) || firstSubmits.Load() != 1 {
		t.Fatal("previous menu was not usable after failed replacement")
	}
	CancelPlayerInventoryMenus(host, player)
	if firstDrops.Load() != 1 {
		t.Fatalf("first drops = %d, want 1", firstDrops.Load())
	}
}

func TestInventoryMenuTerminalWaitsForInflightSubmit(t *testing.T) {
	var sent PlayerInventoryMenu
	host := registerHost(inventoryMenuHost{send: func(_ InvocationID, _ PlayerID, menu PlayerInventoryMenu) bool {
		sent = menu
		return true
	}})
	t.Cleanup(func() { unregisterHost(host) })
	started, release, submitted := make(chan struct{}), make(chan struct{}), make(chan struct{})
	var drops atomic.Int32
	player := PlayerID{Generation: 13}
	if !sendPlayerInventoryMenu(host, 0, player, PlayerInventoryMenu{},
		func(InvocationID, PlayerSnapshot, ItemStack) bool {
			close(started)
			<-release
			return true
		},
		func(InvocationID, PlayerSnapshot) bool { return true },
		func() { drops.Add(1) }) {
		t.Fatal("inventory menu send rejected")
	}
	go func() {
		SubmitPlayerInventoryMenu(sent.ID, 1, PlayerSnapshot{Player: player}, ItemStack{})
		close(submitted)
	}()
	<-started
	CancelPlayerInventoryMenu(sent.ID)
	if drops.Load() != 0 {
		t.Fatal("menu context dropped while submit callback was running")
	}
	close(release)
	<-submitted
	if drops.Load() != 1 {
		t.Fatalf("drops = %d, want 1", drops.Load())
	}
}

func TestInventoryMenusWithSamePlayerAreIsolatedByHost(t *testing.T) {
	player := PlayerID{Generation: 14}
	var first, second PlayerInventoryMenu
	var firstDrops, secondDrops, secondSubmits atomic.Int32
	firstHost := registerHost(inventoryMenuHost{send: func(_ InvocationID, _ PlayerID, menu PlayerInventoryMenu) bool {
		first = menu
		return true
	}})
	defer unregisterHost(firstHost)
	secondHost := registerHost(inventoryMenuHost{send: func(_ InvocationID, _ PlayerID, menu PlayerInventoryMenu) bool {
		second = menu
		return true
	}})
	defer unregisterHost(secondHost)

	if !sendPlayerInventoryMenu(firstHost, 0, player, PlayerInventoryMenu{},
		func(InvocationID, PlayerSnapshot, ItemStack) bool { return true },
		func(InvocationID, PlayerSnapshot) bool { return true },
		func() { firstDrops.Add(1) }) ||
		!sendPlayerInventoryMenu(secondHost, 0, player, PlayerInventoryMenu{},
			func(InvocationID, PlayerSnapshot, ItemStack) bool { secondSubmits.Add(1); return true },
			func(InvocationID, PlayerSnapshot) bool { return true },
			func() { secondDrops.Add(1) }) {
		t.Fatal("inventory menu send rejected")
	}
	if first.ID == 0 || second.ID == 0 || firstDrops.Load() != 0 || secondDrops.Load() != 0 {
		t.Fatalf("activation = IDs %d/%d drops %d/%d, want non-zero IDs and no drops",
			first.ID, second.ID, firstDrops.Load(), secondDrops.Load())
	}

	CancelPlayerInventoryMenus(firstHost, player)
	if firstDrops.Load() != 1 || secondDrops.Load() != 0 {
		t.Fatalf("first-host cancel drops = %d/%d, want 1/0", firstDrops.Load(), secondDrops.Load())
	}
	if !SubmitPlayerInventoryMenu(second.ID, 1, PlayerSnapshot{Player: player}, ItemStack{}) || secondSubmits.Load() != 1 {
		t.Fatal("second-host menu was not usable after first-host cancellation")
	}
	CancelPlayerInventoryMenus(secondHost, player)
}

func TestInventoryMenuDrainDiscardsHostBeforeDroppingContext(t *testing.T) {
	var sent PlayerInventoryMenu
	var discarded, droppedBeforeDiscard atomic.Bool
	var drops atomic.Int32
	host := registerHost(inventoryMenuHost{
		send: func(_ InvocationID, _ PlayerID, menu PlayerInventoryMenu) bool {
			sent = menu
			return true
		},
		discardAll: func() { discarded.Store(true) },
	})
	defer unregisterHost(host)
	if !sendPlayerInventoryMenu(host, 0, PlayerID{Generation: 15}, PlayerInventoryMenu{},
		func(InvocationID, PlayerSnapshot, ItemStack) bool { return true },
		func(InvocationID, PlayerSnapshot) bool { return true },
		func() {
			if !discarded.Load() {
				droppedBeforeDiscard.Store(true)
			}
			drops.Add(1)
		}) {
		t.Fatal("inventory menu send rejected")
	}

	drainHostInventoryMenus(host, true)
	if sent.ID == 0 || !discarded.Load() || droppedBeforeDiscard.Load() || drops.Load() != 1 {
		t.Fatalf("drain = id %d discarded %v early drop %v drops %d, want non-zero/true/false/1",
			sent.ID, discarded.Load(), droppedBeforeDiscard.Load(), drops.Load())
	}
}

func TestInventoryMenuActivationFailureRollsBackExactHostMenu(t *testing.T) {
	started, release := make(chan struct{}), make(chan struct{})
	rolledBack := make(chan uint64, 1)
	var sent PlayerInventoryMenu
	var drops atomic.Int32
	host := registerHost(inventoryMenuHost{
		send: func(_ InvocationID, _ PlayerID, menu PlayerInventoryMenu) bool {
			sent = menu
			close(started)
			<-release
			return true
		},
		discardOne: func(_ PlayerID, id uint64) bool {
			rolledBack <- id
			return true
		},
	})
	defer unregisterHost(host)
	result := make(chan bool, 1)
	go func() {
		result <- sendPlayerInventoryMenu(host, 0, PlayerID{Generation: 17}, PlayerInventoryMenu{},
			func(InvocationID, PlayerSnapshot, ItemStack) bool { return true },
			func(InvocationID, PlayerSnapshot) bool { return true },
			func() { drops.Add(1) })
	}()
	<-started
	drainHostInventoryMenus(host, true)
	close(release)
	if <-result {
		t.Fatal("menu activated after host drain")
	}
	if id := <-rolledBack; id == 0 || id != sent.ID || drops.Load() != 1 {
		t.Fatalf("rollback = id %d sent %d drops %d, want matching non-zero ID and one drop", id, sent.ID, drops.Load())
	}
}
