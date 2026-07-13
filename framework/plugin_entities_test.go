package framework

import (
	"context"
	"errors"
	"math"
	"path/filepath"
	"testing"
	"time"

	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block/cube"
	dfentity "github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"github.com/go-gl/mathgl/mgl64"
)

type recordingEntityRuntime struct {
	adoptTypeKey uint64
	adoptOpaque  uint64
	adoptResult  native.EntityInstanceID
	adoptErr     error
	tickInstance native.EntityInstanceID
	tickInput    native.EntityTickInput
	tickOutput   native.EntityTickOutput
	hurtInput    native.EntityHurtInput
	hurtOutput   native.EntityHurtOutput
	healInput    native.EntityHealInput
	healOutput   native.EntityHealOutput
	deathInput   native.EntityDeathInput
	deathOutput  native.EntityDeathOutput
	saveOutput   native.EntitySaveOutput
	loadInput    native.EntityLoadInput
	loadTypeKey  uint64
	loadResult   native.EntityInstanceID
	loadErr      error
	saveErr      error
	saved        []native.EntityInstanceID
	destroyed    []native.EntityInstanceID
	onTick       func(native.EntityTickInput)
	onInvocation func(native.InvocationID)
}

type recordingEntityViewer struct {
	world.NopViewer
	actions []world.EntityAction
}

func (v *recordingEntityViewer) ViewEntityAction(_ world.Entity, action world.EntityAction) {
	v.actions = append(v.actions, action)
}

func (r *recordingEntityRuntime) EntityAdopt(typeKey, opaque uint64) (native.EntityInstanceID, error) {
	r.adoptTypeKey, r.adoptOpaque = typeKey, opaque
	return r.adoptResult, r.adoptErr
}

func (r *recordingEntityRuntime) EntityLoad(typeKey uint64, input native.EntityLoadInput) (native.EntityInstanceID, error) {
	r.loadTypeKey, r.loadInput = typeKey, input
	return r.loadResult, r.loadErr
}
func (r *recordingEntityRuntime) EntitySave(instance native.EntityInstanceID) (native.EntitySaveOutput, error) {
	r.saved = append(r.saved, instance)
	return r.saveOutput, r.saveErr
}

func (r *recordingEntityRuntime) EntityTick(instance native.EntityInstanceID, input native.EntityTickInput) (native.EntityTickOutput, error) {
	r.tickInstance = instance
	r.tickInput = input
	if r.onTick != nil {
		r.onTick(input)
	}
	return r.tickOutput, nil
}
func (r *recordingEntityRuntime) EntityHurt(_ native.EntityInstanceID, input native.EntityHurtInput) (native.EntityHurtOutput, error) {
	r.hurtInput = input
	if r.onInvocation != nil {
		r.onInvocation(input.Invocation)
	}
	output := r.hurtOutput
	if output == (native.EntityHurtOutput{}) {
		output.Damage = input.Damage
	}
	return output, nil
}
func (r *recordingEntityRuntime) EntityHeal(_ native.EntityInstanceID, input native.EntityHealInput) (native.EntityHealOutput, error) {
	r.healInput = input
	if r.onInvocation != nil {
		r.onInvocation(input.Invocation)
	}
	output := r.healOutput
	if output == (native.EntityHealOutput{}) {
		output.Amount = input.Amount
	}
	return output, nil
}
func (r *recordingEntityRuntime) EntityDeath(_ native.EntityInstanceID, input native.EntityDeathInput) (native.EntityDeathOutput, error) {
	r.deathInput = input
	if r.onInvocation != nil {
		r.onInvocation(input.Invocation)
	}
	return r.deathOutput, nil
}
func (r *recordingEntityRuntime) EntityDestroy(instance native.EntityInstanceID) {
	r.destroyed = append(r.destroyed, instance)
}

func TestBuildEntityRegistryPreservesDragonflyAndAddsForeignTypes(t *testing.T) {
	definition := native.EntityTypeDefinition{
		SaveID: "example:marker", NetworkID: "minecraft:armor_stand",
		Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
	}
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{definition})
	if err != nil {
		t.Fatal(err)
	}
	if registry.Config().TNT == nil || registry.Config().Arrow == nil || registry.Config().Item == nil {
		t.Fatal("Dragonfly entity factories were not preserved")
	}
	if _, ok := registry.Lookup("minecraft:tnt"); !ok {
		t.Fatal("Dragonfly entity types were not preserved")
	}
	entityType, ok := registry.Lookup("example:marker")
	if !ok {
		t.Fatal("custom entity type is missing")
	}
	foreign, ok := entityType.(*foreignBaseEntityType)
	if !ok || foreign.NetworkEncodeEntity() != "minecraft:armor_stand" {
		t.Fatalf("custom entity type = %#v", entityType)
	}
	wantBounds := cube.Box(-0.25, 0, -0.25, 0.25, 1.975, 0.25)
	if got := foreign.BBox(nil); got != wantBounds {
		t.Fatalf("BBox() = %#v, want %#v", got, wantBounds)
	}
}

func TestBuildEntityRegistryRejectsInvalidAndDuplicateTypes(t *testing.T) {
	tests := []native.EntityTypeDefinition{
		{SaveID: "minecraft:tnt", NetworkID: "minecraft:tnt", Max: native.Vec3{X: 1, Y: 1, Z: 1}},
		{SaveID: "missing-namespace", NetworkID: "minecraft:pig", Max: native.Vec3{X: 1, Y: 1, Z: 1}},
		{SaveID: "example:nan", NetworkID: "minecraft:pig", Min: native.Vec3{X: math.NaN()}, Max: native.Vec3{X: 1, Y: 1, Z: 1}},
		{SaveID: "example:inverted", NetworkID: "minecraft:pig", Min: native.Vec3{X: 2}, Max: native.Vec3{X: 1, Y: 1, Z: 1}},
	}
	for _, definition := range tests {
		if _, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{definition}); err == nil {
			t.Fatalf("accepted invalid definition %#v", definition)
		}
	}
	duplicate := native.EntityTypeDefinition{
		SaveID: "example:duplicate", NetworkID: "minecraft:pig", Max: native.Vec3{X: 1, Y: 1, Z: 1},
	}
	if _, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{duplicate, duplicate}); err == nil {
		t.Fatal("accepted duplicate custom entity definitions")
	}
}

func TestWorldManagerSpawnsForeignBaseEntity(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{{
		SaveID: "example:marker", NetworkID: "minecraft:armor_stand",
		Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
	}})
	if err != nil {
		t.Fatal(err)
	}
	w, err := manager.Create("example:entities", world.Config{Synchronous: true, Entities: registry})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	worldID, _ := manager.WorldByName(0, "example:entities")
	if err := w.Do(func(tx *world.Tx) {
		invocation, end := players.BeginInvocation(tx)
		defer end()
		id, ok := manager.SpawnWorldEntity(invocation, worldID, native.EntitySpawn{
			Kind: native.EntityCustom, Type: "example:marker", Position: native.Vec3{X: 1, Y: 64, Z: 2}, NameTag: "Marker",
		})
		if !ok {
			t.Fatal("custom entity spawn failed")
		}
		state, ok := manager.EntityState(invocation, id)
		if !ok || state.Type != "example:marker" || state.NameTag != "Marker" || !state.CanTeleport || state.HasVelocity {
			t.Fatalf("state = %#v, %v", state, ok)
		}
		current, ok := manager.entityHandles.Resolve(id, tx)
		if !ok {
			t.Fatal("custom entity did not resolve")
		}
		if _, ok := current.(world.TickerEntity); ok {
			t.Fatal("base custom entity accidentally implements world.TickerEntity")
		}
		if _, ok := current.(dfentity.Living); ok {
			t.Fatal("base custom entity accidentally implements entity.Living")
		}
		if _, ok := current.(velocityEntity); ok {
			t.Fatal("base custom entity accidentally implements velocity")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestForeignEntityFamiliesExposeOnlyTheirDragonflyCapabilities(t *testing.T) {
	tests := []struct {
		name         string
		definition   native.EntityTypeDefinition
		wantTicker   bool
		wantLiving   bool
		wantVelocity bool
	}{
		{name: "base", definition: native.EntityTypeDefinition{
			SaveID: "example:base", NetworkID: "minecraft:armor_stand",
			Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
			Family: native.EntityFamilyBase,
		}},
		{name: "ticking", definition: native.EntityTypeDefinition{
			SaveID: "example:ticking", NetworkID: "minecraft:armor_stand",
			Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
			TypeKey: 1, Family: native.EntityFamilyTicking,
			CallbackFlags: native.EntityCallbackState | native.EntityCallbackTick,
		}, wantTicker: true, wantVelocity: true},
		{name: "living", definition: native.EntityTypeDefinition{
			SaveID: "example:living", NetworkID: "minecraft:iron_golem",
			Min: native.Vec3{X: -0.7, Z: -0.7}, Max: native.Vec3{X: 0.7, Y: 2.7, Z: 0.7},
			TypeKey: 2, Family: native.EntityFamilyLiving,
			CallbackFlags: native.EntityCallbackState | native.EntityCallbackTick,
			InitialHealth: 40, MaxHealth: 40, Speed: 0.1,
		}, wantTicker: true, wantLiving: true, wantVelocity: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{test.definition})
			if err != nil {
				t.Fatal(err)
			}
			w := world.Config{Synchronous: true, Entities: registry}.New()
			t.Cleanup(func() { _ = w.Close() })
			if err := w.Do(func(tx *world.Tx) {
				manager := NewWorldManager()
				handle, ok := manager.newEntityHandle(tx, native.EntitySpawn{
					Kind: native.EntityCustom, Type: test.definition.SaveID,
					Position: native.Vec3{Y: 64}, CustomInstance: 9,
				})
				if !ok {
					t.Fatal("custom entity family could not be spawned through WorldManager")
				}
				current := tx.AddEntity(handle)
				_, ticker := current.(world.TickerEntity)
				_, living := current.(dfentity.Living)
				_, velocity := current.(velocityEntity)
				if ticker != test.wantTicker || living != test.wantLiving || velocity != test.wantVelocity {
					t.Fatalf("capabilities ticker=%v living=%v velocity=%v", ticker, living, velocity)
				}
			}).Wait(context.Background()); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestForeignTickingEntityInvokesPluginInExactTransaction(t *testing.T) {
	players := host.NewPlayers()
	runtime := &recordingEntityRuntime{adoptResult: 99}
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{{
		SaveID: "example:ticking", NetworkID: "minecraft:armor_stand",
		Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
		TypeKey: 7, Family: native.EntityFamilyTicking,
		CallbackFlags: native.EntityCallbackState | native.EntityCallbackTick,
	}}, foreignEntityServices{runtime: runtime, players: players, entities: players.EntityRegistry()})
	if err != nil {
		t.Fatal(err)
	}
	w := world.Config{Synchronous: true, Entities: registry}.New()
	t.Cleanup(func() { _ = w.Close() })
	if err := w.Do(func(tx *world.Tx) {
		runtime.onTick = func(input native.EntityTickInput) {
			invocationTx, ok := players.InvocationTx(input.Invocation)
			if !ok || invocationTx != tx {
				t.Error("plugin tick did not receive the exact live transaction")
			}
		}
		manager := NewWorldManager()
		handle, ok := manager.newEntityHandle(tx, native.EntitySpawn{
			Kind: native.EntityCustom, Type: "example:ticking", Position: native.Vec3{Y: 64}, CustomInstance: 0xfeed,
		})
		if !ok || runtime.adoptTypeKey != 7 || runtime.adoptOpaque != 0xfeed {
			t.Fatalf("adopt type=%d opaque=%x ok=%v", runtime.adoptTypeKey, runtime.adoptOpaque, ok)
		}
		current := tx.AddEntity(handle)
		current.(world.TickerEntity).Tick(tx, 42)
		if runtime.tickInstance != 99 || runtime.tickInput.Current != 42 || runtime.tickInput.Entity.Generation == 0 {
			t.Fatalf("tick input = %#v", runtime.tickInput)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestForeignEntityDecodeInitializesRuntimeWithoutStateCodec(t *testing.T) {
	runtime := &recordingEntityRuntime{loadResult: 91}
	definition := native.EntityTypeDefinition{
		SaveID: "example:stateless", NetworkID: "minecraft:armor_stand", TypeKey: 77,
		Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
		Family: native.EntityFamilyTicking, CallbackFlags: native.EntityCallbackTick,
	}
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{definition}, foreignEntityServices{runtime: runtime})
	if err != nil {
		t.Fatal(err)
	}
	entityType, _ := registry.Lookup(definition.SaveID)
	data := &world.EntityData{}
	entityType.DecodeNBT(map[string]any{}, data)
	state := data.Data.(*foreignEntityState)
	if state.instance != 91 || runtime.loadTypeKey != 77 || len(runtime.loadInput.Data) != 0 || runtime.loadInput.Version != 0 {
		t.Fatalf("instance=%d type=%d input=%#v", state.instance, runtime.loadTypeKey, runtime.loadInput)
	}
}

func TestForeignEntityPluginTickRunsBeforePhysics(t *testing.T) {
	players := host.NewPlayers()
	runtime := &recordingEntityRuntime{}
	definition := native.EntityTypeDefinition{
		SaveID: "example:tick_order", NetworkID: "minecraft:armor_stand", TypeKey: 78,
		Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
		Family: native.EntityFamilyTicking, CallbackFlags: native.EntityCallbackTick,
		Physics: &native.EntityPhysics{Drag: 0},
	}
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{definition}, foreignEntityServices{
		runtime: runtime, players: players, entities: players.EntityRegistry(),
	})
	if err != nil {
		t.Fatal(err)
	}
	w := world.Config{Synchronous: true, Entities: registry}.New()
	t.Cleanup(func() { _ = w.Close() })
	if err := w.Do(func(tx *world.Tx) {
		entityType, _ := registry.Lookup(definition.SaveID)
		handle := (world.EntitySpawnOpts{Position: mgl64.Vec3{0, 64, 0}}).New(entityType, foreignEntityConfigFor(entityType, 55))
		current := tx.AddEntity(handle).(*foreignTickingEntity)
		runtime.onTick = func(native.EntityTickInput) { current.SetVelocity(mgl64.Vec3{1, 0, 0}) }
		current.Tick(tx, 1)
		if current.Position()[0] <= 0 {
			t.Fatalf("plugin velocity was not applied until after physics: position=%v velocity=%v", current.Position(), current.Velocity())
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestForeignEntityAdoptFailureLeavesOpaqueOwnedByPlugin(t *testing.T) {
	runtime := &recordingEntityRuntime{adoptErr: errors.New("adopt failed")}
	definition := native.EntityTypeDefinition{
		SaveID: "example:adopt_failure", NetworkID: "minecraft:armor_stand", TypeKey: 74,
		Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
		Family: native.EntityFamilyTicking, CallbackFlags: native.EntityCallbackState | native.EntityCallbackTick,
	}
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{definition}, foreignEntityServices{runtime: runtime})
	if err != nil {
		t.Fatal(err)
	}
	w := world.Config{Synchronous: true, Entities: registry}.New()
	t.Cleanup(func() { _ = w.Close() })
	if err := w.Do(func(tx *world.Tx) {
		manager := NewWorldManager()
		_, ok := manager.newEntityHandle(tx, native.EntitySpawn{
			Kind: native.EntityCustom, Type: definition.SaveID, Position: native.Vec3{Y: 64}, CustomInstance: 0xbeef,
		})
		if ok || runtime.adoptOpaque != 0xbeef || len(runtime.destroyed) != 0 {
			t.Fatalf("ok=%v opaque=%x destroyed=%v", ok, runtime.adoptOpaque, runtime.destroyed)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestForeignEntityPhysicsReturnsDragonflyMovement(t *testing.T) {
	definition := native.EntityTypeDefinition{
		SaveID: "example:physics", NetworkID: "minecraft:armor_stand", TypeKey: 73,
		Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
		Family: native.EntityFamilyTicking, CallbackFlags: native.EntityCallbackState | native.EntityCallbackTick,
		Physics: &native.EntityPhysics{Gravity: 0.08, Drag: 0.02, DragBeforeGravity: true},
	}
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{definition})
	if err != nil {
		t.Fatal(err)
	}
	w := world.Config{Synchronous: true, Entities: registry}.New()
	t.Cleanup(func() { _ = w.Close() })
	if err := w.Do(func(tx *world.Tx) {
		entityType, _ := registry.Lookup(definition.SaveID)
		handle := (world.EntitySpawnOpts{Position: mgl64.Vec3{0, 64, 0}, Velocity: mgl64.Vec3{0.1, 0, 0}}).New(entityType, foreignEntityConfigFor(entityType, 0))
		current := tx.AddEntity(handle).(*foreignTickingEntity)
		movement := current.state().Tick(current.Ent, tx)
		if movement == nil || current.Position() == (mgl64.Vec3{0, 64, 0}) {
			t.Fatalf("movement=%#v position=%v", movement, current.Position())
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestForeignLivingEntityPersistsPluginAndDragonflyState(t *testing.T) {
	runtime := &recordingEntityRuntime{
		saveOutput: native.EntitySaveOutput{Data: []byte{1, 2, 3}, Version: 5},
		loadResult: 23,
	}
	definition := native.EntityTypeDefinition{
		SaveID: "example:persistent_living", NetworkID: "minecraft:iron_golem",
		Min: native.Vec3{X: -0.7, Z: -0.7}, Max: native.Vec3{X: 0.7, Y: 2.7, Z: 0.7},
		TypeKey: 71, Family: native.EntityFamilyLiving,
		CallbackFlags: native.EntityCallbackState | native.EntityCallbackTick,
		InitialHealth: 40, MaxHealth: 40, Speed: 0.1, StateVersion: 5,
	}
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{definition}, foreignEntityServices{runtime: runtime})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "world")
	provider, err := (mcdb.Config{}).Open(path)
	if err != nil {
		t.Fatal(err)
	}
	first := world.Config{Synchronous: true, Provider: provider, Entities: registry}.New()
	if err := first.Do(func(tx *world.Tx) {
		entityType, _ := registry.Lookup(definition.SaveID)
		handle := (world.EntitySpawnOpts{Position: mgl64.Vec3{1, 64, 2}}).New(entityType, foreignEntityConfigFor(entityType, 17))
		living := tx.AddEntity(handle).(dfentity.Living)
		living.SetMaxHealth(60)
		living.Hurt(12, dfentity.FallDamageSource{})
		living.SetSpeed(0.35)
		living.AddEffect(effect.New(effect.Speed, 2, 5*time.Second))
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	if len(runtime.saved) != 1 || runtime.saved[0] != 17 {
		t.Fatalf("saved instances = %v", runtime.saved)
	}
	if len(runtime.destroyed) != 1 || runtime.destroyed[0] != 17 {
		t.Fatalf("destroyed instances = %v", runtime.destroyed)
	}

	provider, err = (mcdb.Config{}).Open(path)
	if err != nil {
		t.Fatal(err)
	}
	second := world.Config{Synchronous: true, Provider: provider, Entities: registry}.New()
	t.Cleanup(func() { _ = second.Close() })
	if err := second.Do(func(tx *world.Tx) {
		_ = tx.Block(cube.Pos{1, 64, 2})
		for current := range tx.Entities() {
			if current.H().Type().EncodeEntity() != definition.SaveID {
				continue
			}
			current.(world.TickerEntity).Tick(tx, 1)
			living := current.(dfentity.Living)
			if living.Health() != 28 || living.MaxHealth() != 60 || math.Abs(living.Speed()-0.49) > 0.000001 {
				t.Fatalf("loaded living state health=%v max=%v speed=%v", living.Health(), living.MaxHealth(), living.Speed())
			}
			effects := living.Effects()
			if len(effects) != 1 || effects[0].Type() != effect.Speed || effects[0].Level() != 2 {
				t.Fatalf("loaded effects = %#v", effects)
			}
			return
		}
		t.Fatal("persisted living entity was not loaded")
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if runtime.loadTypeKey != 71 || runtime.loadResult != 23 || runtime.loadInput.Version != 5 || string(runtime.loadInput.Data) != string([]byte{1, 2, 3}) {
		t.Fatalf("load type=%d input=%#v", runtime.loadTypeKey, runtime.loadInput)
	}
}

func TestForeignEntityPersistencePreservesRawStateAcrossRuntimeFailures(t *testing.T) {
	runtime := &recordingEntityRuntime{loadErr: errors.New("load failed")}
	definition := native.EntityTypeDefinition{
		SaveID: "example:failed_state", NetworkID: "minecraft:armor_stand", TypeKey: 79,
		Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
		Family: native.EntityFamilyTicking, CallbackFlags: native.EntityCallbackState | native.EntityCallbackTick,
		StateVersion: 4,
	}
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{definition}, foreignEntityServices{runtime: runtime})
	if err != nil {
		t.Fatal(err)
	}
	entityType, _ := registry.Lookup(definition.SaveID)
	data := &world.EntityData{}
	want := []byte{9, 8, 7}
	entityType.DecodeNBT(map[string]any{
		foreignEntityStateDataKey: want, foreignEntityStateVersionKey: int64(4),
	}, data)

	encoded := entityType.EncodeNBT(data)
	got, _ := encoded[foreignEntityStateDataKey].([]byte)
	if string(got) != string(want) || encoded[foreignEntityStateVersionKey] != int64(4) {
		t.Fatalf("encoded state=%v version=%v", got, encoded[foreignEntityStateVersionKey])
	}

	runtime.loadErr = nil
	runtime.loadResult = 92
	data = &world.EntityData{}
	entityType.DecodeNBT(map[string]any{
		foreignEntityStateDataKey: want, foreignEntityStateVersionKey: int64(4),
	}, data)
	runtime.saveErr = errors.New("save failed")
	encoded = entityType.EncodeNBT(data)
	got, _ = encoded[foreignEntityStateDataKey].([]byte)
	if string(got) != string(want) || encoded[foreignEntityStateVersionKey] != int64(4) || len(runtime.saved) != 1 {
		t.Fatalf("save fallback state=%v version=%v saves=%v", got, encoded[foreignEntityStateVersionKey], runtime.saved)
	}
}

func TestForeignLivingEntityRejectsPoisonedNBTState(t *testing.T) {
	definition := native.EntityTypeDefinition{
		SaveID: "example:poisoned_nbt", NetworkID: "minecraft:iron_golem", TypeKey: 75,
		Min: native.Vec3{X: -0.7, Z: -0.7}, Max: native.Vec3{X: 0.7, Y: 2.7, Z: 0.7},
		Family: native.EntityFamilyLiving, CallbackFlags: native.EntityCallbackTick,
		InitialHealth: 20, MaxHealth: 40, Speed: 0.1,
	}
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{definition})
	if err != nil {
		t.Fatal(err)
	}
	entityType, _ := registry.Lookup(definition.SaveID)
	data := &world.EntityData{}
	entityType.DecodeNBT(map[string]any{
		foreignEntityHealthKey: math.NaN(), foreignEntityMaxHealthKey: float64(-2), foreignEntitySpeedKey: math.Inf(1),
	}, data)
	state := data.Data.(*foreignEntityState)
	if state.health.Health() != 20 || state.health.MaxHealth() != 40 || state.speed != 0.1 || state.baseSpeed != 0.1 {
		t.Fatalf("health=%v max=%v speed=%v base=%v", state.health.Health(), state.health.MaxHealth(), state.speed, state.baseSpeed)
	}
}

func TestForeignLivingEntityCallbacksControlDamageHealingAndDeath(t *testing.T) {
	players := host.NewPlayers()
	runtime := &recordingEntityRuntime{
		hurtOutput:  native.EntityHurtOutput{Damage: 3},
		healOutput:  native.EntityHealOutput{Amount: 5},
		deathOutput: native.EntityDeathOutput{Cancelled: true},
	}
	definition := native.EntityTypeDefinition{
		SaveID: "example:callback_living", NetworkID: "minecraft:iron_golem",
		Min: native.Vec3{X: -0.7, Z: -0.7}, Max: native.Vec3{X: 0.7, Y: 2.7, Z: 0.7},
		TypeKey: 72, Family: native.EntityFamilyLiving,
		// No owner hurt/heal/death flags: Runtime must still receive these calls so it can fan out global events.
		CallbackFlags: native.EntityCallbackState | native.EntityCallbackTick,
		InitialHealth: 40, MaxHealth: 60, Speed: 0.1,
	}
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{definition}, foreignEntityServices{
		runtime: runtime, players: players, entities: players.EntityRegistry(),
	})
	if err != nil {
		t.Fatal(err)
	}
	w := world.Config{Synchronous: true, Entities: registry}.New()
	t.Cleanup(func() { _ = w.Close() })
	var livingHandle *world.EntityHandle
	if err := w.Do(func(tx *world.Tx) {
		runtime.onInvocation = func(invocation native.InvocationID) {
			invocationTx, ok := players.InvocationTx(invocation)
			if !ok || invocationTx != tx {
				t.Error("living callback did not receive the exact live transaction")
			}
		}
		entityType, _ := registry.Lookup(definition.SaveID)
		handle := (world.EntitySpawnOpts{Position: mgl64.Vec3{0, 64, 0}}).New(entityType, foreignEntityConfigFor(entityType, 31))
		livingHandle = handle
		living := tx.AddEntity(handle).(dfentity.Living)
		dealt, vulnerable := living.Hurt(4, dfentity.FallDamageSource{})
		if dealt != 3 || !vulnerable || living.Health() != 37 || runtime.hurtInput.Damage != 4 || runtime.hurtInput.Source.Kind != native.DamageSourceFall {
			t.Fatalf("hurt dealt=%v vulnerable=%v health=%v input=%#v", dealt, vulnerable, living.Health(), runtime.hurtInput)
		}
		healed := living.Heal(20, dfentity.FoodHealingSource{})
		if healed != 5 || living.Health() != 42 || runtime.healInput.Amount != 20 || runtime.healInput.Source.Kind != native.HealingSourceFood {
			t.Fatalf("heal amount=%v health=%v input=%#v", healed, living.Health(), runtime.healInput)
		}

		runtime.hurtOutput.Damage = 100
		dealt, vulnerable = living.Hurt(100, dfentity.VoidDamageSource{})
		if dealt != 0 || vulnerable || living.Health() != 42 || runtime.deathInput.Damage != 100 {
			t.Fatalf("cancelled death dealt=%v vulnerable=%v health=%v death=%#v", dealt, vulnerable, living.Health(), runtime.deathInput)
		}

		runtime.deathOutput.Cancelled = false
		dealt, vulnerable = living.Hurt(100, dfentity.VoidDamageSource{})
		if dealt != 42 || !vulnerable || !living.Dead() {
			t.Fatalf("fatal hurt dealt=%v vulnerable=%v health=%v", dealt, vulnerable, living.Health())
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := w.Do(func(tx *world.Tx) {
		for current := range tx.Entities() {
			if current.H() == livingHandle {
				t.Error("fatally damaged custom living entity remained in the world")
			}
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !livingHandle.Closed() || len(runtime.destroyed) != 1 || runtime.destroyed[0] != 31 {
		t.Fatalf("closed=%v destroyed=%v", livingHandle.Closed(), runtime.destroyed)
	}
}

func TestForeignLivingEntityBroadcastsHurtAndDeathUsingCurrentTransaction(t *testing.T) {
	definition := native.EntityTypeDefinition{
		SaveID: "example:actions", NetworkID: "minecraft:iron_golem", TypeKey: 80,
		Min: native.Vec3{X: -0.7, Z: -0.7}, Max: native.Vec3{X: 0.7, Y: 2.7, Z: 0.7},
		Family: native.EntityFamilyLiving, CallbackFlags: native.EntityCallbackTick,
		InitialHealth: 20, MaxHealth: 20, Speed: 0.1,
	}
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{definition})
	if err != nil {
		t.Fatal(err)
	}
	w := world.Config{Synchronous: true, Entities: registry}.New()
	t.Cleanup(func() { _ = w.Close() })
	viewer := &recordingEntityViewer{}
	loader := world.NewLoader(1, w, viewer)
	if err := w.Do(func(tx *world.Tx) {
		loader.Load(tx, 1)
		defer loader.Close(tx)
		if len(tx.Viewers(mgl64.Vec3{0, 64, 0})) == 0 {
			t.Fatal("loader did not register a viewer")
		}
		entityType, _ := registry.Lookup(definition.SaveID)
		handle := (world.EntitySpawnOpts{Position: mgl64.Vec3{0, 64, 0}}).New(entityType, foreignEntityConfigFor(entityType, 0))
		living := tx.AddEntity(handle).(dfentity.Living)
		living.Hurt(1, dfentity.FallDamageSource{})
		living.Hurt(19, dfentity.FallDamageSource{})
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(viewer.actions) != 3 {
		t.Fatalf("actions=%#v", viewer.actions)
	}
	if _, ok := viewer.actions[0].(dfentity.HurtAction); !ok {
		t.Fatalf("first action=%T", viewer.actions[0])
	}
	if _, ok := viewer.actions[1].(dfentity.HurtAction); !ok {
		t.Fatalf("second action=%T", viewer.actions[1])
	}
	if _, ok := viewer.actions[2].(dfentity.DeathAction); !ok {
		t.Fatalf("third action=%T", viewer.actions[2])
	}
}

func TestWorldManagerPropagatesPluginEntityRegistryToCustomWorlds(t *testing.T) {
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{{
		SaveID: "example:marker", NetworkID: "minecraft:armor_stand",
		Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
	}})
	if err != nil {
		t.Fatal(err)
	}
	manager := NewWorldManager()
	core := world.Config{Synchronous: true, Entities: registry}.New()
	t.Cleanup(func() { _ = core.Close() })
	if err := manager.RegisterCore(OverworldID, core); err != nil {
		t.Fatal(err)
	}
	custom, err := manager.Create("example:custom", world.Config{Synchronous: true})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })
	if _, ok := custom.EntityRegistry().Lookup("example:marker"); !ok {
		t.Fatal("custom world did not inherit plugin entity registry")
	}
}

func TestForeignBaseEntityPersistsThroughDragonflyProvider(t *testing.T) {
	registry, err := buildEntityRegistry(dfentity.DefaultRegistry, []native.EntityTypeDefinition{{
		SaveID: "example:marker", NetworkID: "minecraft:armor_stand",
		Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 1.975, Z: 0.25},
	}})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "world")
	provider, err := (mcdb.Config{}).Open(path)
	if err != nil {
		t.Fatal(err)
	}
	first := world.Config{Synchronous: true, Provider: provider, Entities: registry}.New()
	if err := first.Do(func(tx *world.Tx) {
		entityType, _ := registry.Lookup("example:marker")
		handle := (world.EntitySpawnOpts{
			Position: mgl64.Vec3{1, 64, 2}, NameTag: "Persistent marker",
		}).New(entityType, foreignBaseEntityConfig{})
		tx.AddEntity(handle)
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	provider, err = (mcdb.Config{}).Open(path)
	if err != nil {
		t.Fatal(err)
	}
	second := world.Config{Synchronous: true, Provider: provider, Entities: registry}.New()
	t.Cleanup(func() { _ = second.Close() })
	if err := second.Do(func(tx *world.Tx) {
		_ = tx.Block(cube.Pos{1, 64, 2})
		for current := range tx.Entities() {
			if current.H().Type().EncodeEntity() != "example:marker" {
				continue
			}
			named, ok := current.(nameTagEntity)
			if !ok {
				t.Fatalf("persisted entity type = %T", current)
			}
			if named.NameTag() != "Persistent marker" {
				t.Fatalf("persisted entity name = %q", named.NameTag())
			}
			return
		}
		t.Fatal("persisted custom entity was not loaded")
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}
