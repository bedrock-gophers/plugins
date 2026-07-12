package native

import (
	"sync"
	"sync/atomic"
)

// Host executes synchronous actions requested by native plugins.
type Host interface {
	MessagePlayer(PlayerID, string) bool
}

type noopHost struct{}

func (noopHost) MessagePlayer(PlayerID, string) bool { return false }

var (
	hostSequence atomic.Uint64
	hosts        sync.Map
)

func registerHost(host Host) uint64 {
	if host == nil {
		host = noopHost{}
	}
	id := hostSequence.Add(1)
	hosts.Store(id, host)
	return id
}

func unregisterHost(id uint64) {
	if id != 0 {
		hosts.Delete(id)
	}
}

func resolveHost(id uint64) (Host, bool) {
	host, ok := hosts.Load(id)
	if !ok {
		return nil, false
	}
	return host.(Host), true
}
