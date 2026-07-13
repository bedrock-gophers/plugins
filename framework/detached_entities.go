package framework

import (
	"math"
	"sync"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/world"
)

type detachedEntityEntry struct {
	handle  *world.EntityHandle
	cleanup func()
}

type detachedEntities struct {
	mu      sync.Mutex
	next    uint64
	entries map[native.DetachedEntityID]detachedEntityEntry
}

func newDetachedEntities() *detachedEntities {
	return &detachedEntities{entries: make(map[native.DetachedEntityID]detachedEntityEntry)}
}

func (d *detachedEntities) put(handle *world.EntityHandle, cleanup func()) native.DetachedEntityID {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.next == math.MaxUint64 {
		return native.DetachedEntityID{}
	}
	d.next++
	id := native.DetachedEntityID{Value: d.next, Generation: d.next}
	d.entries[id] = detachedEntityEntry{handle: handle, cleanup: cleanup}
	return id
}

func (d *detachedEntities) take(id native.DetachedEntityID) (*world.EntityHandle, bool) {
	if !id.Valid() {
		return nil, false
	}
	d.mu.Lock()
	entry, ok := d.entries[id]
	if ok {
		delete(d.entries, id)
	}
	d.mu.Unlock()
	return entry.handle, ok
}

func (d *detachedEntities) drop(id native.DetachedEntityID) {
	if !id.Valid() {
		return
	}
	d.mu.Lock()
	entry, ok := d.entries[id]
	if ok {
		delete(d.entries, id)
	}
	d.mu.Unlock()
	if ok && entry.cleanup != nil {
		entry.cleanup()
	}
}

func (d *detachedEntities) drain() {
	d.mu.Lock()
	entries := d.entries
	d.entries = make(map[native.DetachedEntityID]detachedEntityEntry)
	d.mu.Unlock()
	for _, entry := range entries {
		if entry.cleanup != nil {
			entry.cleanup()
		}
	}
}
