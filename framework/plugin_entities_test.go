package framework

import (
	"context"
	"errors"
	"testing"

	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type managedEntityRuntimeStub struct {
	adoptPlugin, adoptType, adoptOpaque uint64
	nextOpen                            native.EntityOpenID
	handles                             map[native.EntityOpenID]native.EntityHandleID
	openCommon                          native.EntityCommonData
	encodeErr                           error
	openErr                             error
	tickName                            string
	ticks, releases, destroys           int
}

func (r *managedEntityRuntimeStub) EntityAdoptLocal(plugin, entityType, opaque uint64) (native.EntityInstanceID, error) {
	r.adoptPlugin, r.adoptType, r.adoptOpaque = plugin, entityType, opaque
	return 11, nil
}

func (r *managedEntityRuntimeStub) EntityDecodeNBT(_ uint64, common native.EntityCommonData, _ []byte) (native.EntityInstanceID, native.EntityCommonData, error) {
	return 12, common, nil
}

func (r *managedEntityRuntimeStub) EntityEncodeNBT(_ native.EntityInstanceID, common native.EntityCommonData) ([]byte, native.EntityCommonData, error) {
	if r.encodeErr != nil {
		return nil, common, r.encodeErr
	}
	encoded, _ := host.MarshalNBT(map[string]any{"value": int32(7)})
	return encoded, common, nil
}

func TestManagedEntityEncodeNBTFallsBackToLastValidState(t *testing.T) {
	runtime := &managedEntityRuntimeStub{}
	entityType := &managedEntityType{
		definition: native.EntityTypeDefinition{SaveID: "example:marker"},
		services:   foreignEntityServices{runtime: runtime},
	}
	data := &world.EntityData{Data: &managedEntityState{runtime: runtime, instance: 11}}

	if value := entityType.EncodeNBT(data)["value"]; value != int32(7) {
		t.Fatalf("initial NBT value = %#v", value)
	}
	runtime.encodeErr = errors.New("plugin callback unavailable")
	if value := entityType.EncodeNBT(data)["value"]; value != int32(7) {
		t.Fatalf("fallback NBT value = %#v", value)
	}
}

func (r *managedEntityRuntimeStub) EntityOpen(_ native.EntityInstanceID, _ native.InvocationID, handle native.EntityHandleID, common native.EntityCommonData) (native.EntityOpenID, uint32, native.EntityCommonData, error) {
	if r.openErr != nil {
		return 0, 0, common, r.openErr
	}
	r.nextOpen++
	if r.handles == nil {
		r.handles = map[native.EntityOpenID]native.EntityHandleID{}
	}
	r.handles[r.nextOpen] = handle
	r.openCommon = common
	return r.nextOpen, native.EntityCapabilityTicker, common, nil
}

func (r *managedEntityRuntimeStub) EntityBBox(_ native.EntityOpenID, _ native.InvocationID, common native.EntityCommonData) (native.BBox, native.EntityCommonData, error) {
	return native.BBox{Min: native.Vec3{X: -0.25, Z: -0.25}, Max: native.Vec3{X: 0.25, Y: 0.75, Z: 0.25}}, common, nil
}

func (r *managedEntityRuntimeStub) EntityClose(_ native.EntityOpenID, _ native.InvocationID, common native.EntityCommonData) (native.EntityCommonData, error) {
	return common, nil
}

func (r *managedEntityRuntimeStub) EntityH(open native.EntityOpenID, _ native.InvocationID, common native.EntityCommonData) (native.EntityHandleID, native.EntityCommonData, error) {
	return r.handles[open], common, nil
}

func (r *managedEntityRuntimeStub) EntityPosition(_ native.EntityOpenID, _ native.InvocationID, common native.EntityCommonData) (native.Vec3, native.EntityCommonData, error) {
	return common.Position, common, nil
}

func (r *managedEntityRuntimeStub) EntityRotation(_ native.EntityOpenID, _ native.InvocationID, common native.EntityCommonData) (native.Rotation, native.EntityCommonData, error) {
	return common.Rotation, common, nil
}

func (r *managedEntityRuntimeStub) EntityTickExact(_ native.EntityOpenID, _ native.InvocationID, _ int64, common native.EntityCommonData) (native.EntityCommonData, error) {
	r.ticks++
	if r.tickName != "" {
		common.Name = r.tickName
	}
	return common, nil
}

func (r *managedEntityRuntimeStub) EntityReleaseOpen(native.EntityOpenID) { r.releases++ }
func (r *managedEntityRuntimeStub) EntityDestroy(native.EntityInstanceID) { r.destroys++ }

type managedEntityViewer struct {
	world.NopViewer
	stateUpdates int
}

func (v *managedEntityViewer) ViewEntityState(world.Entity) { v.stateUpdates++ }

func TestManagedEntityRegistryUsesEncodedIdentityOnly(t *testing.T) {
	registry, err := buildEntityRegistry(entity.DefaultRegistry, []native.EntityTypeDefinition{{
		SaveID: "example:marker", NetworkID: "example:marker", TypeKey: 9,
	}})
	if err != nil {
		t.Fatal(err)
	}
	entityType, ok := registry.Lookup("example:marker")
	managed, exact := managedEntityTypeInfo(entityType)
	if !ok || !exact || managed.definition.TypeKey != 9 {
		t.Fatalf("custom type = %#v, found=%v exact=%v", entityType, ok, exact)
	}
}

func TestManagedEntityRegistryRejectsDuplicate(t *testing.T) {
	definition := native.EntityTypeDefinition{SaveID: "example:marker", NetworkID: "example:marker", TypeKey: 1}
	if _, err := buildEntityRegistry(entity.DefaultRegistry, []native.EntityTypeDefinition{definition, definition}); err == nil {
		t.Fatal("expected duplicate entity type error")
	}
}

func TestManagedEntityWorldlessLifecycleUsesExactCallbacks(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	runtime := &managedEntityRuntimeStub{}
	registry, err := buildEntityRegistry(entity.DefaultRegistry, []native.EntityTypeDefinition{{
		SaveID: "example:marker", NetworkID: "example:marker", TypeKey: 9,
	}}, foreignEntityServices{runtime: runtime, players: players, entities: manager.entityHandles})
	if err != nil {
		t.Fatal(err)
	}
	manager.entityTypes, manager.registriesReady = registry, true
	openedWorld, err := manager.Create("example:entities", world.Config{Synchronous: true})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = manager.CloseCustom() })

	handleID, ok := manager.NewEntity(native.EntitySpawnOptions{
		Position: native.Vec3{X: 1, Y: 65, Z: 2}, Velocity: native.Vec3{X: 0.25},
		Rotation: native.Rotation{Yaw: 45, Pitch: -10}, NameTag: "Marker",
		Type: "example:marker", Plugin: 0xface, LocalType: 9, Opaque: 0xbeef,
	})
	if !ok || !handleID.Valid() {
		t.Fatalf("NewEntity() = %#v, %v", handleID, ok)
	}
	if runtime.adoptPlugin != 0xface || runtime.adoptType != 9 || runtime.adoptOpaque != 0xbeef {
		t.Fatalf("adopt = plugin %#x, type %d, opaque %#x", runtime.adoptPlugin, runtime.adoptType, runtime.adoptOpaque)
	}

	if err := openedWorld.Do(func(tx *world.Tx) {
		viewer := &managedEntityViewer{}
		loader := world.NewLoader(1, openedWorld, viewer)
		loader.Load(tx, 9)
		defer loader.Close(tx)
		invocation, end := players.BeginInvocation(tx)
		defer end()
		entityID, ok := manager.AddEntity(invocation, handleID, nil)
		if !ok {
			t.Fatal("worldless entity could not be added")
		}
		current, ok := manager.entityHandles.Resolve(entityID, tx)
		if !ok {
			t.Fatal("added entity could not be resolved")
		}
		if current.Position() != (mgl64.Vec3{1, 65, 2}) || current.Rotation() != (cube.Rotation{45, -10}) {
			t.Fatalf("transform = %v, %v", current.Position(), current.Rotation())
		}
		named, ok := current.(nameTagEntity)
		if !ok {
			t.Fatal("name tag capability was not preserved")
		}
		if named.NameTag() != "Marker" {
			t.Fatalf("name tag = %q", named.NameTag())
		}
		named.SetNameTag("Renamed marker")
		if named.NameTag() != "Renamed marker" {
			t.Fatalf("updated name tag = %q", named.NameTag())
		}
		moving, ok := current.(velocityEntity)
		if !ok {
			t.Fatal("velocity capability was not preserved")
		}
		if moving.Velocity() != (mgl64.Vec3{0.25, 0, 0}) {
			t.Fatalf("velocity = %v", moving.Velocity())
		}
		moving.SetVelocity(mgl64.Vec3{0, 0.5, 0})
		teleporting, ok := current.(teleportEntity)
		if !ok {
			t.Fatal("teleport capability was not preserved")
		}
		teleporting.Teleport(mgl64.Vec3{4, 70, 5})
		if current.Position() != (mgl64.Vec3{4, 70, 5}) {
			t.Fatalf("teleported position = %v", current.Position())
		}
		teleporting.Teleport(mgl64.Vec3{1, 65, 2})
		if _, ok := current.(world.TickerEntity); !ok {
			t.Fatal("C# ticker capability was not preserved")
		}
		viewer.stateUpdates = 0
		runtime.tickName = "Ticked marker"
		current.(world.TickerEntity).Tick(tx, 42)
		if named.NameTag() != "Ticked marker" || viewer.stateUpdates != 1 {
			t.Fatalf("ticked name tag = %q, state updates=%d", named.NameTag(), viewer.stateUpdates)
		}
		entityType := current.H().Type()
		if got := entityType.BBox(current); got != cube.Box(-0.25, 0, -0.25, 0.25, 0.75, 0.25) {
			t.Fatalf("BBox() = %v", got)
		}
		if _, ok := manager.RemoveEntity(invocation, entityID); !ok {
			t.Fatal("managed entity could not be removed")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if runtime.ticks != 1 || runtime.releases == 0 {
		t.Fatalf("ticks=%d releases=%d", runtime.ticks, runtime.releases)
	}
	if !manager.CloseEntityHandle(handleID) || runtime.destroys != 1 {
		t.Fatalf("close result/destroys = %d", runtime.destroys)
	}

	shutdownHandle, ok := manager.NewEntity(native.EntitySpawnOptions{
		Position: native.Vec3{X: 1, Y: 65, Z: 2}, NameTag: "Shutdown marker",
		Type: "example:marker", Plugin: 0xface, LocalType: 9, Opaque: 0xcafe,
	})
	if !ok {
		t.Fatal("shutdown entity could not be created")
	}
	if err := openedWorld.Do(func(tx *world.Tx) {
		invocation, end := players.BeginInvocation(tx)
		defer end()
		if _, ok := manager.AddEntity(invocation, shutdownHandle, nil); !ok {
			t.Fatal("shutdown entity could not be added")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	runtime.openErr = errors.New("plugin callback unavailable")
	if err := openedWorld.Close(); err != nil {
		t.Fatal(err)
	}
	if runtime.destroys != 2 {
		t.Fatalf("shutdown fallback destroys = %d", runtime.destroys)
	}
}
