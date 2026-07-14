package native

import "testing"

func TestItemSnapshotPairOwnsIndependentClones(t *testing.T) {
	host := registerHost(noopHost{})
	t.Cleanup(func() { unregisterHost(host) })
	main := ItemStack{Identifier: "minecraft:diamond_sword", Count: 1, Lore: []string{"main"}}
	offhand := ItemStack{Identifier: "minecraft:shield", Count: 1, Lore: []string{"offhand"}}
	mainID, offhandID, ok := registerItemSnapshotPair(host, main, offhand)
	if !ok || mainID == 0 || offhandID == 0 || mainID == offhandID {
		t.Fatalf("snapshot pair = %d, %d, %v", mainID, offhandID, ok)
	}
	main.Lore[0], offhand.Lore[0] = "changed", "changed"
	gotMain, mainOK := resolveItemSnapshot(host, mainID)
	gotOffhand, offhandOK := resolveItemSnapshot(host, offhandID)
	if !mainOK || !offhandOK || gotMain.Lore[0] != "main" || gotOffhand.Lore[0] != "offhand" {
		t.Fatalf("snapshot pair = main %#v/%v offhand %#v/%v", gotMain, mainOK, gotOffhand, offhandOK)
	}
	unregisterItemSnapshot(host, mainID)
	if _, ok := resolveItemSnapshot(host, mainID); ok {
		t.Fatal("main-hand snapshot remained after close")
	}
	if _, ok := resolveItemSnapshot(host, offhandID); !ok {
		t.Fatal("closing main-hand snapshot also closed off-hand snapshot")
	}
	unregisterItemSnapshot(host, offhandID)
}

func TestItemSnapshotPairRejectsWithoutTwoFreeSlots(t *testing.T) {
	host := registerHost(noopHost{})
	t.Cleanup(func() { unregisterHost(host) })
	for index := 0; index < maxItemSnapshotsPerHost-1; index++ {
		if _, ok := registerItemSnapshot(host, ItemStack{}); !ok {
			t.Fatalf("snapshot %d rejected", index)
		}
	}
	if first, second, ok := registerItemSnapshotPair(host, ItemStack{}, ItemStack{}); ok || first != 0 || second != 0 {
		t.Fatalf("snapshot pair accepted at capacity: %d, %d, %v", first, second, ok)
	}
	itemSnapshotMu.Lock()
	count := itemSnapshotCounts[host]
	itemSnapshotMu.Unlock()
	if count != maxItemSnapshotsPerHost-1 {
		t.Fatalf("snapshot count changed after rejected pair: %d", count)
	}
}
