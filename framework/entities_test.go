package framework

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

func TestWorldManagerEntityLifecycleUsesExactTransaction(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	w, err := manager.Create("example:entities", world.Config{Synchronous: true, Entities: entity.DefaultRegistry})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	worldID, _ := manager.WorldByName(0, "example:entities")

	var staleInvocation native.InvocationID
	var spawned native.EntityID
	if err := w.Do(func(tx *world.Tx) {
		invocation, end := players.BeginInvocation(tx)
		defer end()
		staleInvocation = invocation
		id, ok := manager.SpawnWorldEntity(invocation, worldID, native.EntitySpawn{
			Kind: native.EntityText, Position: native.Vec3{X: 1, Y: 64, Z: 2}, Text: "hello",
		})
		if !ok {
			t.Fatal("spawn text failed")
		}
		spawned = id
		state, ok := manager.EntityState(invocation, id)
		if !ok || state.World != worldID || state.Type != "dragonfly:text" || state.Position != (native.Vec3{X: 1, Y: 64, Z: 2}) || !state.CanTeleport {
			t.Fatalf("state = %#v, %v", state, ok)
		}
		if !manager.TeleportEntity(invocation, id, native.Vec3{X: 4, Y: 70, Z: 5}) ||
			!manager.SetEntityVelocity(invocation, id, native.Vec3{Y: 0.5}) ||
			!manager.SetEntityNameTag(invocation, id, "renamed") {
			t.Fatal("entity mutation failed")
		}
		state, ok = manager.EntityState(invocation, id)
		if !ok || state.Position != (native.Vec3{X: 4, Y: 70, Z: 5}) || state.Velocity.Y != 0.5 || state.NameTag != "renamed" {
			t.Fatalf("mutated state = %#v, %v", state, ok)
		}
		ids, ok := manager.WorldEntities(invocation, worldID)
		if !ok || !slices.Contains(ids, id) {
			t.Fatalf("entities = %#v, %v", ids, ok)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if manager.SetEntityNameTag(staleInvocation, spawned, "late") {
		t.Fatal("expired invocation mutated an entity")
	}
	if !manager.DespawnEntity(0, spawned) {
		t.Fatal("off-callback despawn failed")
	}
}

func TestWorldManagerRejectsWorldlessEntityWithoutBlocking(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	w, err := manager.Create("example:worldless", world.Config{Synchronous: true, Entities: entity.DefaultRegistry})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	worldID, _ := manager.WorldByName(0, "example:worldless")
	var id native.EntityID
	if err := w.Do(func(tx *world.Tx) {
		invocation, end := players.BeginInvocation(tx)
		defer end()
		var ok bool
		id, ok = manager.SpawnWorldEntity(invocation, worldID, native.EntitySpawn{
			Kind: native.EntityText, Position: native.Vec3{Y: 64}, Text: "temporary",
		})
		if !ok {
			t.Fatal("spawn text failed")
		}
		current, ok := manager.entityHandles.Resolve(id, tx)
		if !ok || tx.RemoveEntity(current) == nil {
			t.Fatal("remove entity failed")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	done := make(chan bool, 1)
	go func() {
		_, ok := manager.EntityState(0, id)
		done <- ok
	}()
	select {
	case ok := <-done:
		if ok {
			t.Fatal("worldless entity returned state")
		}
	case <-time.After(time.Second):
		t.Fatal("worldless entity state blocked")
	}
}

func TestWorldManagerRejectsOverflowingEntityDuration(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	w, err := manager.Create("example:duration", world.Config{Synchronous: true, Entities: entity.DefaultRegistry})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	worldID, _ := manager.WorldByName(0, "example:duration")
	if err := w.Do(func(tx *world.Tx) {
		invocation, end := players.BeginInvocation(tx)
		defer end()
		if _, ok := manager.SpawnWorldEntity(invocation, worldID, native.EntitySpawn{
			Kind: native.EntityTNT, Position: native.Vec3{Y: 64}, FuseMilliseconds: ^uint64(0),
		}); ok {
			t.Fatal("overflowing duration was accepted")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestWorldManagerSpawnsTypedProjectiles(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	w, err := manager.Create("example:projectiles", world.Config{Synchronous: true, Entities: entity.DefaultRegistry})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	worldID, _ := manager.WorldByName(0, "example:projectiles")

	if err := w.Do(func(tx *world.Tx) {
		playerUUID := uuid.MustParse("f898b46d-5ad3-40d6-b48b-40345b9622be")
		handle := world.EntitySpawnOpts{ID: playerUUID}.New(player.Type, player.Config{UUID: playerUUID, Name: "Owner"})
		connected := tx.AddEntity(handle).(*player.Player)
		players.Register(connected, 80)
		owner, ok := players.EntityRegistry().ID(connected)
		if !ok {
			t.Fatal("owner ID missing")
		}
		invocation, end := players.BeginInvocation(tx)
		defer end()
		for _, spawn := range []native.EntitySpawn{
			{Kind: native.EntitySnowball, Owner: owner, Position: native.Vec3{Y: 65}, Velocity: native.Vec3{Z: 1}},
			{Kind: native.EntityArrow, Owner: owner, Position: native.Vec3{Y: 65}, Velocity: native.Vec3{Z: 1}, Damage: 3, Flags: native.EntityArrowCritical, Potion: 25},
			{Kind: native.EntityTNT, Position: native.Vec3{Y: 64}, FuseMilliseconds: uint64((10 * time.Second).Milliseconds())},
			{Kind: native.EntityExperienceOrb, Position: native.Vec3{Y: 64}, Experience: 7},
		} {
			id, ok := manager.SpawnWorldEntity(invocation, worldID, spawn)
			if !ok || id.Generation == 0 {
				t.Fatalf("spawn kind %d failed", spawn.Kind)
			}
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}
