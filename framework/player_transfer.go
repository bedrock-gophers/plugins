package framework

import (
	"errors"
	"sync"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// worldTransferLeases keeps both managed worlds open until a player handle is
// attached to one of them. release is safe from all completion paths.
type worldTransferLeases struct {
	once        sync.Once
	source      *managedWorld
	destination *managedWorld
}

func (l *worldTransferLeases) release() {
	l.once.Do(func() {
		if l.destination != l.source {
			l.destination.lifecycle.RUnlock()
		}
		l.source.lifecycle.RUnlock()
	})
}

// TransferPlayer moves a registered player to a managed world. Cross-world
// movement is deferred until the current plugin callback has finished so all
// later player actions in that callback retain a valid source transaction.
func (m *WorldManager) TransferPlayer(invocation native.InvocationID, id native.PlayerID, destinationID native.WorldID, position native.Vec3) bool {
	if m.players == nil || !finiteVec3(position) {
		return false
	}
	handle, ok := m.players.Handle(id)
	if !ok {
		return false
	}
	if invocation == 0 {
		if _, ok := m.entryByHandle(destinationID); !ok {
			return false
		}
		task := handle.Do(func(tx *world.Tx, entity world.Entity) {
			connected, playerOK := entity.(*player.Player)
			if playerOK {
				m.transferPlayerFromTransaction(tx, connected, handle, destinationID, position)
			}
		})
		return task.Err() == nil
	}
	tx, ok := m.players.InvocationTx(invocation)
	if !ok {
		return false
	}
	connected, ok := transferPlayerInTransaction(handle, tx)
	if !ok {
		return false
	}
	return m.transferPlayerFromTransaction(tx, connected, handle, destinationID, position)
}

func (m *WorldManager) transferPlayerFromTransaction(tx *world.Tx, connected *player.Player, handle *world.EntityHandle, destinationID native.WorldID, position native.Vec3) bool {
	sourceWorld, ok := transactionWorld(tx)
	if !ok {
		return false
	}
	source, destination, leases, ok := m.acquireTransferLeases(sourceWorld, destinationID)
	if !ok {
		return false
	}
	target := mgl64.Vec3{position.X, position.Y, position.Z}
	if source == destination {
		defer leases.release()
		connected.Teleport(target)
		return true
	}

	tx.Defer(func(sourceTx *world.Tx) {
		m.runDeferredPlayerTransfer(sourceTx, handle, source, destination, target, leases)
	})
	return true
}

func (m *WorldManager) acquireTransferLeases(sourceWorld *world.World, destinationID native.WorldID) (*managedWorld, *managedWorld, *worldTransferLeases, bool) {
	m.mu.RLock()
	source := m.byWorld[sourceWorld]
	destination := m.handles[destinationID]
	if source == nil || destination == nil || source.unloading || destination.unloading {
		m.mu.RUnlock()
		return nil, nil, nil, false
	}
	m.mu.RUnlock()

	source.lifecycle.RLock()
	if destination != source {
		destination.lifecycle.RLock()
	}
	m.mu.RLock()
	sourceCurrent := m.byWorld[sourceWorld] == source && m.handles[source.id] == source && !source.unloading
	destinationCurrent := m.handles[destinationID] == destination && !destination.unloading
	m.mu.RUnlock()
	if !sourceCurrent || !destinationCurrent || source.closed || destination.closed {
		if destination != source {
			destination.lifecycle.RUnlock()
		}
		source.lifecycle.RUnlock()
		return nil, nil, nil, false
	}
	leases := &worldTransferLeases{source: source, destination: destination}
	return source, destination, leases, true
}

func (m *WorldManager) runDeferredPlayerTransfer(sourceTx *world.Tx, handle *world.EntityHandle, source, destination *managedWorld, position mgl64.Vec3, leases *worldTransferLeases) {
	deferredActive := true
	defer func() {
		if recover() != nil || deferredActive {
			leases.release()
		}
	}()
	connected, ok := transferPlayerInTransaction(handle, sourceTx)
	if !ok || sourceTx.World() != source.world || handle.Closed() {
		return
	}
	if removed := sourceTx.RemoveEntity(connected); removed != handle {
		return
	}

	task := destination.world.Do(func(destinationTx *world.Tx) {
		if !handle.Closed() {
			destinationTx.AddEntityAt(handle, position)
		}
	})
	deferredActive = false
	select {
	case <-task.Done():
		m.finishPlayerTransfer(task.Err(), sourceTx, handle, source, leases)
	default:
		task.OnDone(func(err error) {
			m.finishPlayerTransfer(err, nil, handle, source, leases)
		})
	}
}

func (m *WorldManager) finishPlayerTransfer(err error, sourceTx *world.Tx, handle *world.EntityHandle, source *managedWorld, leases *worldTransferLeases) {
	if !errors.Is(err, world.ErrWorldClosed) {
		leases.release()
		return
	}
	if sourceTx != nil {
		m.restoreTransferredPlayer(sourceTx, handle)
		leases.release()
		return
	}
	task := source.world.Do(func(tx *world.Tx) {
		m.restoreTransferredPlayer(tx, handle)
	})
	select {
	case <-task.Done():
		leases.release()
	default:
		task.OnDone(func(error) { leases.release() })
	}
}

func (m *WorldManager) restoreTransferredPlayer(tx *world.Tx, handle *world.EntityHandle) {
	if handle.Closed() {
		return
	}
	tx.AddEntity(handle)
	m.players.ForgetWorldDeparture(handle)
}

func transferPlayerInTransaction(handle *world.EntityHandle, tx *world.Tx) (connected *player.Player, ok bool) {
	if handle == nil || tx == nil {
		return nil, false
	}
	defer func() {
		if recover() != nil {
			connected, ok = nil, false
		}
	}()
	entity, ok := handle.Entity(tx)
	if !ok {
		return nil, false
	}
	connected, ok = entity.(*player.Player)
	return connected, ok
}
