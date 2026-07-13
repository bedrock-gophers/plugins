package framework

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/df-mc/dragonfly/server/world"
)

func TestDetachedRegistryConsumesExactlyOnce(t *testing.T) {
	var cleaned atomic.Int32
	registry := newDetachedEntities()
	handle := new(world.EntityHandle)
	id := registry.put(handle, func() { cleaned.Add(1) })
	var winners atomic.Int32
	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		defer wait.Done()
		if _, ok := registry.take(id); ok {
			winners.Add(1)
		}
	}()
	go func() { defer wait.Done(); registry.drop(id) }()
	wait.Wait()
	if winners.Load()+cleaned.Load() != 1 {
		t.Fatalf("take winners=%d cleanup=%d", winners.Load(), cleaned.Load())
	}
}

func TestDetachedRegistryAllocationIsMonotonicAndBounded(t *testing.T) {
	registry := newDetachedEntities()
	first := registry.put(new(world.EntityHandle), nil)
	second := registry.put(new(world.EntityHandle), nil)
	if !first.Valid() || second.Value <= first.Value || second.Generation <= first.Generation {
		t.Fatalf("detached IDs are not monotonic: %#v then %#v", first, second)
	}
	registry.next = math.MaxUint64
	if exhausted := registry.put(new(world.EntityHandle), nil); exhausted.Valid() || len(registry.entries) != 2 {
		t.Fatalf("exhausted allocation = %#v, entries=%d", exhausted, len(registry.entries))
	}
}

func TestDetachedRegistryDrainCleansEntriesExactlyOnce(t *testing.T) {
	var cleaned atomic.Int32
	registry := newDetachedEntities()
	first := registry.put(new(world.EntityHandle), func() { cleaned.Add(1) })
	second := registry.put(new(world.EntityHandle), func() { cleaned.Add(1) })
	registry.drain()
	registry.drop(first)
	registry.drop(second)
	registry.drain()
	if cleaned.Load() != 2 {
		t.Fatalf("cleanup calls = %d", cleaned.Load())
	}
	if _, ok := registry.take(first); ok {
		t.Fatal("drained entry remained available")
	}
}
