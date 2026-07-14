package host

import (
	"iter"
	"sync"

	"github.com/bedrock-gophers/plugins/internal/native"
	dragonflyserver "github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

type serverPlayerSource interface {
	Players(*world.Tx) iter.Seq[*player.Player]
	Player(uuid.UUID) (*world.EntityHandle, bool)
}

// Server exposes the owned Dragonfly server through the private plugin host.
type Server struct {
	mu      sync.RWMutex
	source  serverPlayerSource
	players *Players

	iteratorMu sync.Mutex
	iterators  map[native.PlayerIteratorID]*serverPlayerIterator
	next       native.PlayerIteratorID
}

type serverPlayerIterator struct {
	mu         sync.Mutex
	invocation native.InvocationID
	next       func() (*player.Player, bool)
	stop       func()
	endCurrent func()
	stopped    bool
}

func NewServer(players *Players) *Server {
	return &Server{players: players, iterators: map[native.PlayerIteratorID]*serverPlayerIterator{}}
}

// Attach publishes the constructed Dragonfly server before plugins are enabled.
func (s *Server) Attach(source *dragonflyserver.Server) {
	s.mu.Lock()
	s.source = source
	s.mu.Unlock()
}

func (s *Server) server() (serverPlayerSource, bool) {
	s.mu.RLock()
	source := s.source
	s.mu.RUnlock()
	return source, source != nil && s.players != nil
}

// OpenServerPlayerIterator maps directly to Dragonfly Server.Players. A zero
// invocation deliberately passes nil rather than guessing the current world.
func (s *Server) OpenServerPlayerIterator(invocation native.InvocationID) (native.PlayerIteratorID, bool) {
	source, ok := s.server()
	if !ok {
		return 0, false
	}
	var tx *world.Tx
	if invocation != 0 {
		tx, ok = s.players.InvocationTx(invocation)
		if !ok {
			return 0, false
		}
	}
	next, stop := iter.Pull(source.Players(tx))
	iterator := &serverPlayerIterator{invocation: invocation, next: next, stop: stop}
	s.iteratorMu.Lock()
	if s.next == native.PlayerIteratorID(^uint64(0)) {
		s.iteratorMu.Unlock()
		iterator.close()
		return 0, false
	}
	s.next++
	id := s.next
	s.iterators[id] = iterator
	s.iteratorMu.Unlock()
	if invocation != 0 && !s.players.OnInvocationEnd(invocation, func() { s.closeServerPlayers(invocation, id) }) {
		s.closeServerPlayers(invocation, id)
		return 0, false
	}
	return id, true
}

// NextServerPlayer advances one iterator. The yielded PlayerSnapshot is bound
// to a fresh nested invocation that remains live only until the next advance
// or close, matching the lifetime of Dragonfly's loop variable.
func (s *Server) NextServerPlayer(invocation native.InvocationID, id native.PlayerIteratorID) (native.InvocationID, native.PlayerSnapshot, bool, bool) {
	s.iteratorMu.Lock()
	iterator, ok := s.iterators[id]
	s.iteratorMu.Unlock()
	if !ok || iterator.invocation != invocation {
		return 0, native.PlayerSnapshot{}, false, false
	}
	nested, snapshot, found, valid := iterator.advance(s.players)
	if !valid || !found {
		s.closeServerPlayers(invocation, id)
	}
	return nested, snapshot, found, valid
}

func (s *Server) CloseServerPlayers(invocation native.InvocationID, id native.PlayerIteratorID) {
	s.closeServerPlayers(invocation, id)
}

func (s *Server) closeServerPlayers(invocation native.InvocationID, id native.PlayerIteratorID) {
	s.iteratorMu.Lock()
	iterator, ok := s.iterators[id]
	if ok && iterator.invocation == invocation {
		delete(s.iterators, id)
	} else {
		iterator = nil
	}
	s.iteratorMu.Unlock()
	if iterator != nil {
		iterator.close()
	}
}

// Close releases every iterator retained by a plugin during shutdown.
func (s *Server) Close() {
	s.iteratorMu.Lock()
	iterators := make([]*serverPlayerIterator, 0, len(s.iterators))
	for id, iterator := range s.iterators {
		delete(s.iterators, id)
		iterators = append(iterators, iterator)
	}
	s.iteratorMu.Unlock()
	for _, iterator := range iterators {
		iterator.close()
	}
}

func (i *serverPlayerIterator) advance(players *Players) (nested native.InvocationID, snapshot native.PlayerSnapshot, found bool, valid bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.stopped {
		return 0, native.PlayerSnapshot{}, false, false
	}
	defer func() {
		if recover() != nil {
			i.finishLocked()
			nested, snapshot, found, valid = 0, native.PlayerSnapshot{}, false, false
		}
	}()
	i.endCurrentLocked()
	connected, found := i.next()
	if !found {
		i.finishLocked()
		return 0, native.PlayerSnapshot{}, false, true
	}
	id, ok := players.ID(connected)
	if !ok {
		i.finishLocked()
		return 0, native.PlayerSnapshot{}, false, false
	}
	nested, end := players.BeginInvocation(connected.Tx())
	if nested == 0 {
		i.finishLocked()
		return 0, native.PlayerSnapshot{}, false, false
	}
	i.endCurrent = end
	position := connected.Position()
	return nested, native.PlayerSnapshot{
		Player: id, Name: connected.Name(),
		LatencyMilliseconds: uint64(max(connected.Latency().Milliseconds(), 0)),
		Position:            native.Vec3{X: position.X(), Y: position.Y(), Z: position.Z()},
	}, true, true
}

func (i *serverPlayerIterator) close() {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer func() { _ = recover() }()
	if i.stopped {
		return
	}
	i.finishLocked()
}

func (i *serverPlayerIterator) endCurrentLocked() {
	if i.endCurrent != nil {
		i.endCurrent()
		i.endCurrent = nil
	}
}

func (i *serverPlayerIterator) finishLocked() {
	if i.stopped {
		return
	}
	i.endCurrentLocked()
	i.stopped = true
	if i.stop != nil {
		i.stop()
	}
}

func (s *Server) ServerPlayer(id [16]byte) (native.EntityHandleID, bool, bool) {
	source, ok := s.server()
	if !ok {
		return native.EntityHandleID{}, false, false
	}
	handle, found := source.Player(uuid.UUID(id))
	if !found {
		return native.EntityHandleID{}, false, true
	}
	result, ok := s.playerHandleID(handle)
	return result, ok, ok
}

func (s *Server) ServerPlayerByName(name string) (native.EntityHandleID, bool, bool) {
	source, ok := s.server()
	if !ok {
		return native.EntityHandleID{}, false, false
	}
	id, found := s.playerUUIDByExactName(name)
	if !found {
		return native.EntityHandleID{}, false, true
	}
	handle, found := source.Player(id)
	if !found {
		return native.EntityHandleID{}, false, true
	}
	result, ok := s.playerHandleID(handle)
	return result, ok, ok
}

// playerUUIDByExactName preserves Dragonfly's case-sensitive lookup semantics
// without calling Server.PlayerByName, whose current implementation reads the
// online-player map without its mutex. The subsequent UUID lookup confirms the
// cached registry entry is still online through Dragonfly's locked path.
func (s *Server) playerUUIDByExactName(name string) (uuid.UUID, bool) {
	s.players.mu.RLock()
	defer s.players.mu.RUnlock()
	for _, entry := range s.players.entries {
		if entry.name == name && entry.handle != nil && !entry.handle.Closed() {
			return uuid.UUID(entry.id.UUID), true
		}
	}
	return uuid.Nil, false
}

func (s *Server) playerHandleID(handle *world.EntityHandle) (native.EntityHandleID, bool) {
	if handle == nil || s.players == nil {
		return native.EntityHandleID{}, false
	}
	s.players.mu.RLock()
	entry, ok := s.players.entries[handle]
	s.players.mu.RUnlock()
	if !ok {
		return native.EntityHandleID{}, false
	}
	return s.players.entities.EntityHandleID(native.EntityID{UUID: entry.id.UUID, Generation: entry.id.Generation})
}
