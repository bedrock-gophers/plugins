package host

import (
	"context"
	"reflect"
	"testing"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

func TestPlayersInventoryItemRoundTrip(t *testing.T) {
	withPlayer(t, func(player *player.Player) {
		players := NewPlayers()
		playerID := players.Register(player, 7)
		invocation, leave := players.BeginInvocation(player.Tx())
		defer leave()
		valuesNBT, ok := marshalItemNBT(map[string]any{
			"owner": "rust", "level": int32(12), "flags": []byte{1, 0, 1},
		})
		if !ok {
			t.Fatal("encode values")
		}
		want := native.ItemStack{
			Identifier: "minecraft:diamond_sword", Count: 1, Damage: 3,
			Unbreakable: true, AnvilCost: 7,
			CustomName: "Plugin Sword", Lore: []string{"line one", "line two"}, ValuesNBT: valuesNBT,
			Enchantments: []native.ItemEnchantment{{ID: 9, Level: 5}, {ID: 17, Level: 3}},
		}
		main := native.InventoryID{Player: playerID, Kind: native.InventoryMain}
		if size, ok := players.InventorySize(invocation, main); !ok || size != 36 {
			t.Fatalf("main inventory size=%d ok=%v", size, ok)
		}
		if !players.SetInventoryItem(invocation, main, 2, want) {
			t.Fatal("set inventory item")
		}
		got, ok := players.InventoryItem(invocation, main, 2)
		if !ok {
			t.Fatal("read inventory item")
		}
		wantValues, _ := unmarshalItemNBT(want.ValuesNBT)
		gotValues, valuesOK := unmarshalItemNBT(got.ValuesNBT)
		got.ValuesNBT, want.ValuesNBT = nil, nil
		if !valuesOK || !reflect.DeepEqual(gotValues, wantValues) || !reflect.DeepEqual(got, want) {
			t.Fatalf("item mismatch\ngot:  %#v values=%#v\nwant: %#v values=%#v", got, gotValues, want, wantValues)
		}
		if !players.SetHeldSlot(invocation, playerID, 2) {
			t.Fatal("set held slot")
		}
		held, ok := players.HeldItem(invocation, playerID, 0)
		if !ok || held.Identifier != want.Identifier || held.CustomName != want.CustomName {
			t.Fatalf("held=%#v ok=%v", held, ok)
		}
	})
}

func TestPlayersInventoryAddClearAndOffhand(t *testing.T) {
	withPlayer(t, func(player *player.Player) {
		players := NewPlayers()
		playerID := players.Register(player, 9)
		invocation, leave := players.BeginInvocation(player.Tx())
		defer leave()
		main := native.InventoryID{Player: playerID, Kind: native.InventoryMain}
		added, ok := players.AddInventoryItem(invocation, main, native.ItemStack{Identifier: "minecraft:apple", Count: 70})
		if !ok || added != 70 {
			t.Fatalf("added=%d ok=%v", added, ok)
		}
		first, _ := players.InventoryItem(invocation, main, 0)
		second, _ := players.InventoryItem(invocation, main, 1)
		if first.Count != 64 || second.Count != 6 {
			t.Fatalf("apple counts=%d,%d", first.Count, second.Count)
		}
		offhand := native.InventoryID{Player: playerID, Kind: native.InventoryOffhand}
		if !players.SetInventoryItem(invocation, offhand, 0, native.ItemStack{Identifier: "minecraft:totem_of_undying", Count: 1}) {
			t.Fatal("set offhand")
		}
		item, ok := players.HeldItem(invocation, playerID, 1)
		if !ok || item.Identifier != "minecraft:totem_of_undying" {
			t.Fatalf("offhand=%#v ok=%v", item, ok)
		}
		if !players.ClearInventory(invocation, main) || !players.ClearInventory(invocation, offhand) {
			t.Fatal("clear inventory")
		}
		first, _ = players.InventoryItem(invocation, main, 0)
		item, _ = players.InventoryItem(invocation, offhand, 0)
		if first.Count != 0 || item.Count != 0 {
			t.Fatalf("inventories not cleared: main=%#v offhand=%#v", first, item)
		}
	})
}

func TestPlayersHeldItemsRoundTrip(t *testing.T) {
	withPlayer(t, func(player *player.Player) {
		players := NewPlayers()
		playerID := players.Register(player, 10)
		invocation, leave := players.BeginInvocation(player.Tx())
		defer leave()
		main := native.ItemStack{Identifier: "minecraft:apple", Count: 1, CustomName: "Main"}
		offhand := native.ItemStack{Identifier: "minecraft:totem_of_undying", Count: 1, CustomName: "Offhand"}
		if !players.SetHeldItems(invocation, playerID, main, offhand) {
			t.Fatal("set held items")
		}
		gotMain, gotOffhand, ok := players.HeldItems(invocation, playerID)
		if !ok || gotMain.Identifier != main.Identifier || gotMain.CustomName != main.CustomName ||
			gotOffhand.Identifier != offhand.Identifier || gotOffhand.CustomName != offhand.CustomName {
			t.Fatalf("held items = main %#v offhand %#v ok=%v", gotMain, gotOffhand, ok)
		}
	})
}

func TestPlayersRejectInvalidNativeItems(t *testing.T) {
	withPlayer(t, func(player *player.Player) {
		players := NewPlayers()
		id := native.InventoryID{Player: players.Register(player, 1), Kind: native.InventoryMain}
		invocation, leave := players.BeginInvocation(player.Tx())
		defer leave()
		invalid := []native.ItemStack{
			{Identifier: "missing:item", Count: 1},
			{Identifier: "minecraft:apple", Metadata: 1 << 16, Count: 1},
			{Identifier: "minecraft:diamond_sword", Count: 1, Enchantments: []native.ItemEnchantment{{ID: 9}}},
			{Identifier: "minecraft:diamond_sword", Count: 1, ValuesNBT: []byte{1, 2, 3}},
		}
		for index, value := range invalid {
			if players.SetInventoryItem(invocation, id, 0, value) {
				t.Fatalf("invalid item %d accepted", index)
			}
		}
	})
}

func TestPlayersSetInventoryItemRejectsInvalidSlots(t *testing.T) {
	withPlayer(t, func(player *player.Player) {
		players := NewPlayers()
		playerID := players.Register(player, 11)
		invocation, leave := players.BeginInvocation(player.Tx())
		defer leave()

		tests := []struct {
			name        string
			inventory   native.InventoryID
			validSlot   uint32
			invalidSlot uint32
			stack       native.ItemStack
		}{
			{
				name:        "main",
				inventory:   native.InventoryID{Player: playerID, Kind: native.InventoryMain},
				validSlot:   35,
				invalidSlot: 36,
				stack:       native.ItemStack{Identifier: "minecraft:apple", Count: 1},
			},
			{
				name:        "armour",
				inventory:   native.InventoryID{Player: playerID, Kind: native.InventoryArmour},
				validSlot:   3,
				invalidSlot: 4,
				stack:       native.ItemStack{Identifier: "minecraft:diamond_boots", Count: 1},
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				if players.SetInventoryItem(invocation, test.inventory, test.invalidSlot, test.stack) {
					t.Fatalf("invalid slot %d accepted", test.invalidSlot)
				}
				if !players.SetInventoryItem(invocation, test.inventory, test.validSlot, test.stack) {
					t.Fatalf("valid slot %d rejected", test.validSlot)
				}
				got, ok := players.InventoryItem(invocation, test.inventory, test.validSlot)
				if !ok || got.Identifier != test.stack.Identifier || got.Count != test.stack.Count {
					t.Fatalf("valid set did not persist: got=%#v ok=%v", got, ok)
				}
			})
		}
	})
}

func TestPlayersSetInventoryItemSchedulesAcrossWorldInvocation(t *testing.T) {
	source := world.Config{Synchronous: true}.New()
	destination := world.Config{Synchronous: true}.New()
	t.Cleanup(func() {
		_ = source.Close()
		_ = destination.Close()
	})
	players := NewPlayers()
	spawn := func(id uuid.UUID, name string) *world.EntityHandle {
		position := mgl64.Vec3{1, 64, 1}
		return world.EntitySpawnOpts{ID: id, Position: position}.New(
			player.Type,
			player.Config{UUID: id, Name: name, Position: position},
		)
	}
	sourceHandle := spawn(uuid.MustParse("ac1c3bc0-5f73-4561-93ec-f12d42f4ca41"), "Source")
	destinationHandle := spawn(uuid.MustParse("9743d9ae-bf63-47c3-9f38-9661d0313d03"), "Destination")
	var destinationID native.PlayerID
	if err := destination.Do(func(tx *world.Tx) {
		destinationID = players.Register(tx.AddEntity(destinationHandle).(*player.Player), 52)
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := source.Do(func(tx *world.Tx) {
		connected := tx.AddEntity(sourceHandle).(*player.Player)
		players.Register(connected, 51)
		invocation, leave := players.BeginInvocation(tx)
		defer leave()
		accepted := players.SetInventoryItem(invocation, native.InventoryID{
			Player: destinationID,
			Kind:   native.InventoryMain,
		}, 0, native.ItemStack{Identifier: "minecraft:apple", Count: 1})
		if !accepted {
			t.Fatal("cross-world inventory write was not scheduled")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	got, ok := players.InventoryItem(0, native.InventoryID{
		Player: destinationID,
		Kind:   native.InventoryMain,
	}, 0)
	if !ok || got.Identifier != "minecraft:apple" || got.Count != 1 {
		t.Fatalf("scheduled item = %#v, %v", got, ok)
	}
}
