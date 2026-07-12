package host

import (
	"slices"
	"strings"
	"sync"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/player"
)

// Players owns stable native IDs for the lifetime of connected Dragonfly players.
type Players struct {
	mu      sync.RWMutex
	entries map[*player.Player]native.PlayerID
}

func NewPlayers() *Players {
	return &Players{entries: map[*player.Player]native.PlayerID{}}
}

func (p *Players) Register(player *player.Player, generation uint64) native.PlayerID {
	id := native.PlayerID{Generation: generation}
	uuid := player.UUID()
	copy(id.UUID[:], uuid[:])
	p.mu.Lock()
	p.entries[player] = id
	p.mu.Unlock()
	return id
}

func (p *Players) Unregister(player *player.Player) {
	p.mu.Lock()
	delete(p.entries, player)
	p.mu.Unlock()
}

func (p *Players) ID(player *player.Player) (native.PlayerID, bool) {
	p.mu.RLock()
	id, ok := p.entries[player]
	p.mu.RUnlock()
	return id, ok
}

func (p *Players) Names() []string {
	p.mu.RLock()
	names := make([]string, 0, len(p.entries))
	for connected := range p.entries {
		names = append(names, connected.Name())
	}
	p.mu.RUnlock()
	slices.Sort(names)
	return names
}

func (p *Players) CommandSnapshots() []native.CommandPlayer {
	p.mu.RLock()
	snapshots := make([]native.CommandPlayer, 0, len(p.entries))
	for connected, id := range p.entries {
		snapshots = append(snapshots, native.CommandPlayer{
			Player:              id,
			Name:                connected.Name(),
			LatencyMilliseconds: uint64(connected.Latency().Milliseconds()),
		})
	}
	p.mu.RUnlock()
	slices.SortFunc(snapshots, func(left, right native.CommandPlayer) int {
		return strings.Compare(left.Name, right.Name)
	})
	return snapshots
}

func (p *Players) ResolveName(name string) (native.PlayerID, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for connected, id := range p.entries {
		if strings.EqualFold(connected.Name(), name) {
			return id, true
		}
	}
	return native.PlayerID{}, false
}

func (p *Players) ResolveUUID(uuid [16]byte) (native.PlayerID, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, id := range p.entries {
		if id.UUID == uuid {
			return id, true
		}
	}
	return native.PlayerID{}, false
}

func (p *Players) ResolveID(id native.PlayerID) (*player.Player, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for connected, candidate := range p.entries {
		if candidate == id {
			return connected, true
		}
	}
	return nil, false
}

func (p *Players) MessagePlayer(id native.PlayerID, message string) bool {
	connected, ok := p.ResolveID(id)
	if !ok {
		return false
	}
	connected.Message(message)
	return true
}
