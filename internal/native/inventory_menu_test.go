package native

import (
	"sync/atomic"
	"testing"
)

func TestInventoryMenuCancelDuringClickDefersDrop(t *testing.T) {
	host := registerHost(noopHost{})
	defer unregisterHost(host)
	player := PlayerID{Generation: 1}
	snapshot := PlayerSnapshot{Player: player}
	var drops atomic.Int32
	var id uint64
	id, _ = registerInventoryMenu(host, player, func(InvocationID, PlayerSnapshot, uint32) bool {
		CancelPlayerInventoryMenu(id)
		if drops.Load() != 0 {
			t.Fatal("drop ran before the click callback returned")
		}
		return true
	}, func(InvocationID, PlayerSnapshot) bool { return true }, func() { drops.Add(1) })

	if !CompletePlayerInventoryMenuClick(id, 1, snapshot, 3) {
		t.Fatal("click was rejected")
	}
	if drops.Load() != 1 {
		t.Fatalf("drop count = %d, want 1", drops.Load())
	}
}

func TestInventoryMenuCloseDuringClickDefersClose(t *testing.T) {
	host := registerHost(noopHost{})
	defer unregisterHost(host)
	player := PlayerID{Generation: 1}
	snapshot := PlayerSnapshot{Player: player}
	var closes atomic.Int32
	var id uint64
	id, _ = registerInventoryMenu(host, player, func(invocation InvocationID, current PlayerSnapshot, _ uint32) bool {
		if !CompletePlayerInventoryMenuClose(id, invocation, current) {
			t.Fatal("close was rejected")
		}
		if closes.Load() != 0 {
			t.Fatal("close ran before the click callback returned")
		}
		return true
	}, func(InvocationID, PlayerSnapshot) bool {
		closes.Add(1)
		return true
	}, func() { t.Fatal("menu was dropped instead of closed") })

	if !CompletePlayerInventoryMenuClick(id, 1, snapshot, 3) {
		t.Fatal("click was rejected")
	}
	if closes.Load() != 1 {
		t.Fatalf("close count = %d, want 1", closes.Load())
	}
}
