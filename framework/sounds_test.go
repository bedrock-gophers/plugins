package framework

import (
	"context"
	"testing"

	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
)

func TestWorldSoundRejectsExpiredInvocation(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	w, err := manager.Create("example:sounds", world.Config{Synchronous: true, Entities: entity.DefaultRegistry})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	id, _ := manager.WorldByName(0, "example:sounds")
	var stale native.InvocationID
	if err := w.Do(func(tx *world.Tx) {
		var end func()
		stale, end = players.BeginInvocation(tx)
		end()
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if manager.PlayWorldSound(stale, id, native.Vec3{Y: 64}, native.WorldSound{Kind: native.SoundExplosion}) {
		t.Fatal("expired invocation played a sound")
	}
}

func TestWorldSoundReportsSameTransactionDecodeFailure(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	w, err := manager.Create("example:sound-decode", world.Config{Synchronous: true, Entities: entity.DefaultRegistry})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	id, _ := manager.WorldByName(0, "example:sound-decode")
	properties, ok := host.EncodeBlockProperties(map[string]any{})
	if !ok {
		t.Fatal("encode block properties")
	}
	if err := w.Do(func(tx *world.Tx) {
		invocation, end := players.BeginInvocation(tx)
		defer end()
		played := manager.PlayWorldSound(invocation, id, native.Vec3{Y: 64}, native.WorldSound{
			Kind:  native.SoundDoorOpen,
			Block: &native.WorldBlock{Identifier: "example:missing", PropertiesNBT: properties},
		})
		if played {
			t.Fatal("unregistered block sound reported success")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestWorldSoundUsesInvocationWorld(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	w, err := manager.Create("example:sound-invocation", world.Config{Synchronous: true, Entities: entity.DefaultRegistry})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	if err := w.Do(func(tx *world.Tx) {
		invocation, end := players.BeginInvocation(tx)
		defer end()
		if !manager.PlayWorldSound(invocation, 0, native.Vec3{Y: 64}, native.WorldSound{Kind: native.SoundExplosion}) {
			t.Fatal("current transaction world sound rejected")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}
