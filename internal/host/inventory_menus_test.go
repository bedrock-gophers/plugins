package host

import (
	"testing"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestInventoryMenuSpecsMatchPublicContainerSizes(t *testing.T) {
	want := map[native.InventoryMenuContainer]int{
		native.InventoryMenuChest:       27,
		native.InventoryMenuDoubleChest: 54,
		native.InventoryMenuHopper:      5,
		native.InventoryMenuDropper:     9,
		native.InventoryMenuBarrel:      27,
		native.InventoryMenuEnderChest:  27,
	}
	for container, size := range want {
		spec, ok := inventoryMenuSpecFor(container)
		if !ok || spec.block == nil || spec.size != size {
			t.Fatalf("container %d spec = %+v, %v; want size %d", container, spec, ok, size)
		}
	}
	if _, ok := inventoryMenuSpecFor(native.InventoryMenuEnderChest + 1); ok {
		t.Fatal("unknown inventory menu container accepted")
	}
}

func TestDoubleChestBlockActorDataPairsBothHalves(t *testing.T) {
	left := cube.Pos{10, 20, 30}
	right := left.Add(cube.Pos{1, 0, 0})
	chest := block.NewChest()
	leftData := inventoryMenuBlockActorData(chest, "Menu", "Chest", left, &right, true)
	rightData := inventoryMenuBlockActorData(chest, "Menu", "Chest", right, &left, false)

	if leftData["pairx"] != int32(right[0]) || leftData["pairz"] != int32(right[2]) {
		t.Fatalf("left pairing = %v/%v, want %d/%d", leftData["pairx"], leftData["pairz"], right[0], right[2])
	}
	if rightData["pairx"] != int32(left[0]) || rightData["pairz"] != int32(left[2]) {
		t.Fatalf("right pairing = %v/%v, want %d/%d", rightData["pairx"], rightData["pairz"], left[0], left[2])
	}
	if leftData["pairlead"] != byte(1) || rightData["pairlead"] != byte(0) {
		t.Fatalf("pair leads = %v/%v, want 1/0", leftData["pairlead"], rightData["pairlead"])
	}
	for label, data := range map[string]map[string]any{"left": leftData, "right": rightData} {
		if data["id"] != "Chest" || data["CustomName"] != "Menu" || data["Items"] == nil ||
			data["x"] == nil || data["y"] == nil || data["z"] == nil {
			t.Fatalf("%s block actor data is incomplete: %#v", label, data)
		}
	}
}

func TestDropperBlockActorDataIncludesEmptyInventory(t *testing.T) {
	spec, ok := inventoryMenuSpecFor(native.InventoryMenuDropper)
	if !ok {
		t.Fatal("dropper spec unavailable")
	}
	pos := cube.Pos{10, 20, 30}
	data := inventoryMenuBlockActorData(spec.block, "Menu", spec.blockActorID, pos, nil, false)
	if data["id"] != "Dropper" || data["CustomName"] != "Menu" || data["Items"] == nil ||
		data["x"] != int32(pos[0]) || data["y"] != int32(pos[1]) || data["z"] != int32(pos[2]) {
		t.Fatalf("dropper block actor data is incomplete: %#v", data)
	}
}

func TestInventoryMenuSessionPositionIsOutsideWorld(t *testing.T) {
	height := cube.Range{-64, 319}
	visible := cube.Pos{10, 20, 30}
	sessionPos := inventoryMenuSessionPosition(height, visible)
	if !sessionPos.OutOfBounds(height) {
		t.Fatalf("session position %v is inside height range %v", sessionPos, height)
	}
	if sessionPos[0] != visible[0] || sessionPos[2] != visible[2] {
		t.Fatalf("session position changed horizontal coordinates: got %v, visible %v", sessionPos, visible)
	}
}

func TestResetInventoryMenuWindowPreparesReplacement(t *testing.T) {
	s := &session.Session{}
	menu := &openInventoryMenu{
		inv:        inventory.New(27, nil),
		pos:        cube.Pos{10, 20, 30},
		sessionPos: cube.Pos{10, 320, 30},
		windowID:   7,
	}
	setInventoryMenuSession(s, menu, protocol.ContainerTypeContainer)

	containerType, ok := resetInventoryMenuWindow(s, menu.windowID)
	if !ok || containerType != protocol.ContainerTypeContainer {
		t.Fatalf("reset = %d, %v; want %d, true", containerType, ok, protocol.ContainerTypeContainer)
	}
	if inventoryMenuWindowOpen(s, menu.windowID) {
		t.Fatal("reset window remains open")
	}
	if next := nextInventoryMenuWindowID(s); next == menu.windowID {
		t.Fatalf("replacement reused closed window ID %d", next)
	}
}

func TestNormaliseInventoryMenuRequestsRewritesOnlyContainerSlots(t *testing.T) {
	take := &protocol.TakeStackRequestAction{}
	take.Source.Container.ContainerID = protocol.ContainerBarrel
	take.Destination.Container.ContainerID = protocol.ContainerCursor
	place := &protocol.PlaceStackRequestAction{}
	place.Source.Container.ContainerID = protocol.ContainerCursor
	place.Destination.Container.ContainerID = protocol.ContainerShulkerBox
	swap := &protocol.SwapStackRequestAction{}
	swap.Source.Container.ContainerID = protocol.ContainerHotBar
	swap.Destination.Container.ContainerID = protocol.ContainerBarrel
	drop := &protocol.DropStackRequestAction{}
	drop.Source.Container.ContainerID = protocol.ContainerDynamic
	destroyMenu := &protocol.DestroyStackRequestAction{}
	destroyMenu.Source.Container.ContainerID = protocol.ContainerLevelEntity
	destroyPlayer := &protocol.DestroyStackRequestAction{}
	destroyPlayer.Source.Container.ContainerID = protocol.ContainerHotBar

	normaliseInventoryMenuRequests([]protocol.ItemStackRequest{{Actions: []protocol.StackRequestAction{take, place, swap, drop, destroyMenu, destroyPlayer}}})
	if take.Source.Container.ContainerID != protocol.ContainerLevelEntity ||
		take.Destination.Container.ContainerID != protocol.ContainerCursor ||
		place.Source.Container.ContainerID != protocol.ContainerCursor ||
		place.Destination.Container.ContainerID != protocol.ContainerLevelEntity ||
		swap.Source.Container.ContainerID != protocol.ContainerHotBar ||
		swap.Destination.Container.ContainerID != protocol.ContainerLevelEntity ||
		drop.Source.Container.ContainerID != protocol.ContainerLevelEntity ||
		destroyMenu.Source.Container.ContainerID != inventoryMenuRejectedContainerID ||
		destroyPlayer.Source.Container.ContainerID != protocol.ContainerHotBar {
		t.Fatalf("normalised IDs = take %d/%d place %d/%d swap %d/%d drop %d destroy %d/%d",
			take.Source.Container.ContainerID, take.Destination.Container.ContainerID,
			place.Source.Container.ContainerID, place.Destination.Container.ContainerID,
			swap.Source.Container.ContainerID, swap.Destination.Container.ContainerID,
			drop.Source.Container.ContainerID,
			destroyMenu.Source.Container.ContainerID, destroyPlayer.Source.Container.ContainerID)
	}
}
