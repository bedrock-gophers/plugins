package host

import (
	"sync"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/world"
)

// Entities owns stable native IDs for Dragonfly entity handles. An ID remains
// valid until its handle is unregistered. Re-registering a handle receives a
// fresh generation, so an old plugin reference cannot resolve to a new entity.
type Entities struct {
	mu             sync.RWMutex
	byHandle       map[*world.EntityHandle]entityEntry
	byID           map[native.EntityID]*world.EntityHandle
	nextGeneration uint64
}

type entityEntry struct {
	id     native.EntityID
	active bool
}

func NewEntities() *Entities {
	return &Entities{
		byHandle: map[*world.EntityHandle]entityEntry{},
		byID:     map[native.EntityID]*world.EntityHandle{},
	}
}

// Register returns the stable ID of entity. Registering the same live handle
// more than once returns the same ID.
func (e *Entities) Register(entity world.Entity) native.EntityID {
	if entity == nil {
		return native.EntityID{}
	}
	return e.registerHandle(entity.H(), 0)
}

// Unregister expires the ID associated with entity.
func (e *Entities) Unregister(entity world.Entity) {
	if entity != nil {
		e.unregisterHandle(entity.H())
	}
}

// ID returns the live stable ID associated with entity.
func (e *Entities) ID(entity world.Entity) (native.EntityID, bool) {
	if entity == nil {
		return native.EntityID{}, false
	}
	e.mu.RLock()
	entry, ok := e.byHandle[entity.H()]
	e.mu.RUnlock()
	return entry.id, ok && entry.active
}

// Handle returns the live Dragonfly handle associated with id.
func (e *Entities) Handle(id native.EntityID) (*world.EntityHandle, bool) {
	if id.Generation == 0 {
		return nil, false
	}
	e.mu.RLock()
	handle, ok := e.byID[id]
	entry := e.byHandle[handle]
	e.mu.RUnlock()
	return handle, ok && entry.active && handle != nil && !handle.Closed()
}

// Resolve opens id in the exact transaction passed. It never scans worlds or
// schedules work on another owner.
func (e *Entities) Resolve(id native.EntityID, tx *world.Tx) (entity world.Entity, ok bool) {
	handle, ok := e.Handle(id)
	if !ok || tx == nil {
		return nil, false
	}
	defer func() {
		if recover() != nil {
			entity, ok = nil, false
		}
	}()
	return handle.Entity(tx)
}

func (e *Entities) registerHandle(handle *world.EntityHandle, generation uint64) native.EntityID {
	if handle == nil || handle.Closed() {
		return native.EntityID{}
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if entry, ok := e.byHandle[handle]; ok {
		entry.active = true
		e.byHandle[handle] = entry
		return entry.id
	}
	if generation == 0 {
		var ok bool
		generation, ok = e.nextAvailableGenerationLocked()
		if !ok {
			return native.EntityID{}
		}
	} else if generation > e.nextGeneration {
		e.nextGeneration = generation
	}
	id := native.EntityID{Generation: generation}
	uuid := handle.UUID()
	copy(id.UUID[:], uuid[:])
	for {
		if existing, found := e.byID[id]; !found || existing == handle {
			break
		}
		generation, ok := e.nextAvailableGenerationLocked()
		if !ok {
			return native.EntityID{}
		}
		id.Generation = generation
	}
	e.byHandle[handle] = entityEntry{id: id, active: true}
	e.byID[id] = handle
	return id
}

// deactivateHandle keeps identity across a Dragonfly world transfer while
// preventing operations from waiting on a handle that is not in any world.
func (e *Entities) deactivateHandle(handle *world.EntityHandle) {
	if handle == nil {
		return
	}
	e.mu.Lock()
	if entry, ok := e.byHandle[handle]; ok {
		entry.active = false
		e.byHandle[handle] = entry
	}
	e.mu.Unlock()
}

func (e *Entities) nextAvailableGenerationLocked() (uint64, bool) {
	if e.nextGeneration == ^uint64(0) {
		return 0, false
	}
	e.nextGeneration++
	return e.nextGeneration, true
}

func (e *Entities) unregisterHandle(handle *world.EntityHandle) {
	if handle == nil {
		return
	}
	e.mu.Lock()
	if entry, ok := e.byHandle[handle]; ok {
		delete(e.byHandle, handle)
		delete(e.byID, entry.id)
	}
	e.mu.Unlock()
}
