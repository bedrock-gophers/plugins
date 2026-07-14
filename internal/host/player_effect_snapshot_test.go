package host

import (
	"context"
	"image/color"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

type snapshotCustomEffect struct{}

func (snapshotCustomEffect) RGBA() color.RGBA                  { return color.RGBA{} }
func (snapshotCustomEffect) Apply(world.Entity, effect.Effect) {}
func (snapshotCustomEffect) Start(world.Entity, int)           {}
func (snapshotCustomEffect) End(world.Entity, int)             {}

func TestPlayersSnapshotAndClearEffects(t *testing.T) {
	const customID = -31_997
	effect.Register(customID, snapshotCustomEffect{})
	withPlayer(t, func(connected *player.Player) {
		players := NewPlayers()
		id := players.Register(connected, 81)
		invocation, leave := players.BeginInvocation(connected.Tx())
		defer leave()

		connected.AddEffect(effect.New(effect.Speed, 2, 5*time.Second))
		connected.AddEffect(effect.NewAmbient(effect.Regeneration, 1, 2*time.Second).WithoutParticles())
		connected.AddEffect(effect.NewInfinite(effect.FireResistance, 3).WithoutParticles())
		connected.AddEffect(effect.New(snapshotCustomEffect{}, 4, time.Second))

		values, ok := players.PlayerEffects(invocation, id)
		if !ok || len(values) != 4 {
			t.Fatalf("effects = %#v ok=%v", values, ok)
		}
		byType := map[native.EffectType]native.PlayerEffect{}
		for _, value := range values {
			byType[value.Type] = value
		}
		if value := byType[native.EffectSpeed]; value.Level != 2 || value.Duration != 5*time.Second || value.Potency != 0 || value.Ambient || value.Infinite || value.ParticlesHidden || value.Tick != 0 {
			t.Fatalf("speed = %#v", value)
		}
		if value := byType[native.EffectRegeneration]; !value.Ambient || value.Infinite || !value.ParticlesHidden {
			t.Fatalf("ambient = %#v", value)
		}
		if value := byType[native.EffectFireResistance]; value.Ambient || !value.Infinite || value.Duration != 0 || !value.ParticlesHidden {
			t.Fatalf("infinite = %#v", value)
		}
		if value := byType[native.EffectType(customID)]; value.Level != 4 || value.Ambient || value.Infinite {
			t.Fatalf("custom = %#v", value)
		}

		if !players.ClearPlayerEffects(invocation, id) || len(connected.Effects()) != 0 {
			t.Fatalf("effects after clear = %#v", connected.Effects())
		}

		players.Unregister(connected)
		if values, ok := players.PlayerEffects(invocation, id); ok || values != nil {
			t.Fatalf("stale effects = %#v ok=%v", values, ok)
		}
		if players.ClearPlayerEffects(invocation, id) {
			t.Fatal("cleared effects for stale player")
		}
	})
}

func TestPlayersSnapshotEmptyEffects(t *testing.T) {
	withPlayer(t, func(connected *player.Player) {
		players := NewPlayers()
		id := players.Register(connected, 82)
		values, ok := players.PlayerEffects(0, id)
		if !ok || len(values) != 0 {
			t.Fatalf("effects = %#v ok=%v", values, ok)
		}
	})
}

func TestPlayersSnapshotIncludesInitialInstantEffects(t *testing.T) {
	w := world.Config{Synchronous: true}.New()
	t.Cleanup(func() { _ = w.Close() })
	id := uuid.MustParse("56b1b0c8-d899-43ae-a683-ef8714980289")
	handle := world.EntitySpawnOpts{ID: id}.New(player.Type, player.Config{
		UUID: id, Name: "InitialEffects", Effects: []effect.Effect{
			effect.New(effect.Speed, 2, time.Second),
			effect.NewInstant(effect.InstantHealth, 1),
		},
	})
	if err := w.Do(func(tx *world.Tx) {
		connected := tx.AddEntity(handle).(*player.Player)
		players := NewPlayers()
		playerID := players.Register(connected, 83)
		invocation, leave := players.BeginInvocation(tx)
		defer leave()
		values, ok := players.PlayerEffects(invocation, playerID)
		if !ok || len(values) != 2 {
			t.Fatalf("effects = %#v ok=%v", values, ok)
		}
		byType := map[native.EffectType]native.PlayerEffect{}
		for _, value := range values {
			byType[value.Type] = value
		}
		if byType[native.EffectSpeed].Level != 2 || byType[native.EffectInstantHealth].Level != 1 ||
			byType[native.EffectInstantHealth].Potency != 1 {
			t.Fatalf("effects = %#v", values)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestSnapshotPlayerEffectPreservesNegativeRemainingDurationAndTick(t *testing.T) {
	current := effect.New(effect.Speed, 1, time.Millisecond).TickDuration()
	value, ok := snapshotPlayerEffect(current)
	if !ok || current.Duration() >= 0 || value.Duration != current.Duration() || value.Tick != 1 {
		t.Fatalf("current duration=%v snapshot=%#v ok=%v", current.Duration(), value, ok)
	}
}

func TestPlayerEffectFromNativePreservesTickAndNanoseconds(t *testing.T) {
	value := native.PlayerEffect{
		Type: native.EffectSpeed, Level: 2, Duration: time.Second + time.Nanosecond,
		Potency: 1, Ambient: true, ParticlesHidden: true, Tick: 3,
	}
	reconstructed, ok := playerEffectFromNative(effect.Speed, value)
	if !ok || reconstructed.Duration() != value.Duration || reconstructed.Tick() != 3 ||
		!reconstructed.Ambient() || !reconstructed.ParticlesHidden() {
		t.Fatalf("reconstructed=%#v ok=%v", reconstructed, ok)
	}
}

func TestPlayerEffectFromNativePreservesInfiniteTick(t *testing.T) {
	value := native.PlayerEffect{
		Type: native.EffectFireResistance, Level: 1, Potency: 1, Infinite: true, Tick: 2,
	}
	reconstructed, ok := playerEffectFromNative(effect.FireResistance, value)
	if !ok || !reconstructed.Infinite() || reconstructed.Duration() != 0 || reconstructed.Tick() != 2 {
		t.Fatalf("reconstructed=%#v ok=%v", reconstructed, ok)
	}
}

func TestPlayerEffectFromNativeRejectsUnboundedTickWork(t *testing.T) {
	value := native.PlayerEffect{
		Type: native.EffectSpeed, Level: 1, Duration: time.Second,
		Potency: 1, Tick: maximumPlayerEffectTick + 1,
	}
	if reconstructed, ok := playerEffectFromNative(effect.Speed, value); ok {
		t.Fatalf("unbounded tick accepted: %#v", reconstructed)
	}
}

type snapshotLargeIDEffect struct{}

func (snapshotLargeIDEffect) RGBA() color.RGBA                  { return color.RGBA{} }
func (snapshotLargeIDEffect) Apply(world.Entity, effect.Effect) {}
func (snapshotLargeIDEffect) Start(world.Entity, int)           {}
func (snapshotLargeIDEffect) End(world.Entity, int)             {}

func TestSnapshotPlayerEffectRejectsOutOfRangeID(t *testing.T) {
	if strconv.IntSize < 64 {
		t.Skip("requires a 64-bit effect registry ID")
	}
	effect.Register(1<<40, snapshotLargeIDEffect{})
	if _, ok := snapshotPlayerEffect(effect.New(snapshotLargeIDEffect{}, 1, time.Second)); ok {
		t.Fatal("snapshot accepted an effect ID outside i32")
	}
}

func TestSnapshotPlayerEffectRejectsOutOfRangeLevel(t *testing.T) {
	if strconv.IntSize < 64 {
		t.Skip("requires a 64-bit effect level")
	}
	if _, ok := snapshotPlayerEffect(effect.New(effect.Speed, 1<<40, time.Second)); ok {
		t.Fatal("snapshot accepted an effect level outside i32")
	}
}

type snapshotUnregisteredEffect struct{}

func (snapshotUnregisteredEffect) RGBA() color.RGBA                  { return color.RGBA{} }
func (snapshotUnregisteredEffect) Apply(world.Entity, effect.Effect) {}
func (snapshotUnregisteredEffect) Start(world.Entity, int)           {}
func (snapshotUnregisteredEffect) End(world.Entity, int)             {}

func TestPlayersClearEffectsRejectsUnregisteredBeforeRemoval(t *testing.T) {
	w := world.Config{Synchronous: true}.New()
	t.Cleanup(func() { _ = w.Close() })
	id := uuid.MustParse("83e59098-d2ce-4421-8a51-9d62838fb026")
	handle := world.EntitySpawnOpts{ID: id}.New(player.Type, player.Config{
		UUID: id, Name: "UnregisteredEffect", Effects: []effect.Effect{
			effect.New(effect.Speed, 1, time.Second),
			effect.New(snapshotUnregisteredEffect{}, 1, time.Second),
		},
	})
	if err := w.Do(func(tx *world.Tx) {
		connected := tx.AddEntity(handle).(*player.Player)
		players := NewPlayers()
		playerID := players.Register(connected, 84)
		invocation, leave := players.BeginInvocation(tx)
		defer leave()
		if players.ClearPlayerEffects(invocation, playerID) {
			t.Fatal("clear accepted an unregistered effect type")
		}
		if effects := connected.Effects(); len(effects) != 2 {
			t.Fatalf("effects were partially cleared: %#v", effects)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

type snapshotInitialInstantEffect struct{}

var snapshotInitialInstantApplications atomic.Int64

func (snapshotInitialInstantEffect) RGBA() color.RGBA { return color.RGBA{} }
func (snapshotInitialInstantEffect) Apply(world.Entity, effect.Effect) {
	snapshotInitialInstantApplications.Add(1)
}

func TestPlayersClearEffectsMatchesDragonflyInitialEffectFlush(t *testing.T) {
	w := world.Config{Synchronous: true}.New()
	t.Cleanup(func() { _ = w.Close() })
	id := uuid.MustParse("a77ae85d-9493-4194-a8ec-83a30352a291")
	handle := world.EntitySpawnOpts{ID: id}.New(player.Type, player.Config{
		UUID: id, Name: "InitialInstantClear", Effects: []effect.Effect{
			effect.New(effect.Speed, 1, time.Second),
			effect.NewInstant(snapshotInitialInstantEffect{}, 1),
		},
	})
	snapshotInitialInstantApplications.Store(0)
	if err := w.Do(func(tx *world.Tx) {
		connected := tx.AddEntity(handle).(*player.Player)
		players := NewPlayers()
		playerID := players.Register(connected, 85)
		invocation, leave := players.BeginInvocation(tx)
		defer leave()
		if !players.ClearPlayerEffects(invocation, playerID) {
			t.Fatal("clear rejected valid initial effects")
		}
		if effects := connected.Effects(); len(effects) != 0 {
			t.Fatalf("effects after clear = %#v", effects)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if applications := snapshotInitialInstantApplications.Load(); applications != 1 {
		t.Fatalf("initial instant applications = %d", applications)
	}
}

func TestPlayersClearEffectsRejectsMalformedInitialInstantBeforeFlush(t *testing.T) {
	w := world.Config{Synchronous: true}.New()
	t.Cleanup(func() { _ = w.Close() })
	id := uuid.MustParse("c0c83213-c43f-46d5-82d7-c3bf27ee83c7")
	handle := world.EntitySpawnOpts{ID: id}.New(player.Type, player.Config{
		UUID: id, Name: "MalformedInitialClear", Effects: []effect.Effect{
			effect.New(effect.Speed, 1, time.Second),
			effect.NewInstant(snapshotInitialInstantEffect{}, 0),
		},
	})
	snapshotInitialInstantApplications.Store(0)
	if err := w.Do(func(tx *world.Tx) {
		connected := tx.AddEntity(handle).(*player.Player)
		players := NewPlayers()
		playerID := players.Register(connected, 86)
		invocation, leave := players.BeginInvocation(tx)
		defer leave()
		if players.ClearPlayerEffects(invocation, playerID) {
			t.Fatal("clear accepted a malformed initial instant effect")
		}
		if effects := connected.Effects(); len(effects) != 2 {
			t.Fatalf("initial effects changed after rejection: %#v", effects)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if applications := snapshotInitialInstantApplications.Load(); applications != 0 {
		t.Fatalf("malformed initial instant applications = %d", applications)
	}
}
