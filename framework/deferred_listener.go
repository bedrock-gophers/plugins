package framework

import (
	"errors"
	"net"
	"sync"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/session"
)

// deferListenerCreation replaces eager Dragonfly listener factories with
// wrappers tracked by a gate. server.Config.New constructs the wrappers without
// binding sockets. The caller opens every real listener synchronously through
// openAll only after plugin enable succeeds.
func deferListenerCreation(config *server.Config) *deferredListenerGate {
	originals := append([]func(server.Config) (server.Listener, error){}, config.Listeners...)
	gate := newDeferredListenerGate()
	wrapped := make([]func(server.Config) (server.Listener, error), 0, len(originals))
	for _, factory := range originals {
		factory := factory
		wrapped = append(wrapped, func(current server.Config) (server.Listener, error) {
			current.Listeners = originals
			listener := newDeferredListener(func() (server.Listener, error) {
				return factory(current)
			})
			gate.add(listener)
			return listener, nil
		})
	}
	config.Listeners = wrapped
	return gate
}

type deferredListenerGateState uint8

const (
	deferredListenerGateCollecting deferredListenerGateState = iota
	deferredListenerGateOpening
	deferredListenerGateOpen
	deferredListenerGateClosing
	deferredListenerGateClosed
)

type deferredListenerGate struct {
	mu        sync.Mutex
	changed   *sync.Cond
	state     deferredListenerGateState
	listeners []*deferredListener
	openErr   error
	closeErr  error
}

func newDeferredListenerGate() *deferredListenerGate {
	gate := &deferredListenerGate{}
	gate.changed = sync.NewCond(&gate.mu)
	return gate
}

func (gate *deferredListenerGate) add(listener *deferredListener) {
	gate.mu.Lock()
	if gate.state == deferredListenerGateCollecting {
		gate.listeners = append(gate.listeners, listener)
		gate.mu.Unlock()
		return
	}
	gate.mu.Unlock()
	_ = listener.Close()
}

// openAll opens every listener in factory order. If any factory fails, every
// wrapper is closed and listeners opened so far are rolled back in reverse.
func (gate *deferredListenerGate) openAll() error {
	gate.mu.Lock()
	for {
		switch gate.state {
		case deferredListenerGateCollecting:
			gate.state = deferredListenerGateOpening
			listeners := append([]*deferredListener(nil), gate.listeners...)
			gate.mu.Unlock()
			return gate.open(listeners)
		case deferredListenerGateOpening, deferredListenerGateClosing:
			gate.changed.Wait()
		case deferredListenerGateOpen:
			gate.mu.Unlock()
			return nil
		case deferredListenerGateClosed:
			err := gate.openErr
			if err == nil {
				err = net.ErrClosed
			}
			gate.mu.Unlock()
			return err
		}
	}
}

func (gate *deferredListenerGate) open(listeners []*deferredListener) error {
	for _, listener := range listeners {
		if _, err := listener.ensureOpen(); err != nil {
			closeErr := closeDeferredListeners(listeners)
			openErr := errors.Join(err, closeErr)
			gate.mu.Lock()
			gate.state = deferredListenerGateClosed
			gate.openErr = openErr
			gate.closeErr = closeErr
			gate.changed.Broadcast()
			gate.mu.Unlock()
			return openErr
		}
	}

	gate.mu.Lock()
	gate.state = deferredListenerGateOpen
	gate.changed.Broadcast()
	gate.mu.Unlock()
	return nil
}

// closeAll closes all wrappers in reverse factory order. It is safe to call
// concurrently with openAll and is idempotent.
func (gate *deferredListenerGate) closeAll() error {
	gate.mu.Lock()
	for {
		switch gate.state {
		case deferredListenerGateOpening, deferredListenerGateClosing:
			gate.changed.Wait()
		case deferredListenerGateCollecting, deferredListenerGateOpen:
			gate.state = deferredListenerGateClosing
			listeners := append([]*deferredListener(nil), gate.listeners...)
			gate.mu.Unlock()
			closeErr := closeDeferredListeners(listeners)
			gate.mu.Lock()
			gate.state = deferredListenerGateClosed
			gate.closeErr = closeErr
			gate.changed.Broadcast()
			gate.mu.Unlock()
			return closeErr
		case deferredListenerGateClosed:
			err := gate.closeErr
			gate.mu.Unlock()
			return err
		}
	}
}

func closeDeferredListeners(listeners []*deferredListener) error {
	var closeErr error
	for index := len(listeners) - 1; index >= 0; index-- {
		closeErr = errors.Join(closeErr, listeners[index].Close())
	}
	return closeErr
}

type deferredListenerState uint8

const (
	deferredListenerUnopened deferredListenerState = iota
	deferredListenerOpening
	deferredListenerOpen
	deferredListenerFailed
	deferredListenerClosed
)

type deferredListener struct {
	mu       sync.Mutex
	changed  *sync.Cond
	open     func() (server.Listener, error)
	state    deferredListenerState
	listener server.Listener
	err      error
}

func newDeferredListener(open func() (server.Listener, error)) *deferredListener {
	listener := &deferredListener{open: open}
	listener.changed = sync.NewCond(&listener.mu)
	return listener
}

func (listener *deferredListener) Accept() (session.Conn, error) {
	inner, err := listener.awaitOpen()
	if err != nil {
		return nil, err
	}
	return inner.Accept()
}

func (listener *deferredListener) Disconnect(connection session.Conn, reason string) error {
	listener.mu.Lock()
	inner := listener.listener
	state := listener.state
	listener.mu.Unlock()
	if state != deferredListenerOpen || inner == nil {
		return net.ErrClosed
	}
	return inner.Disconnect(connection, reason)
}

func (listener *deferredListener) Close() error {
	listener.mu.Lock()
	for listener.state == deferredListenerOpening {
		listener.changed.Wait()
	}
	if listener.state == deferredListenerClosed {
		listener.mu.Unlock()
		return nil
	}
	inner := listener.listener
	listener.state = deferredListenerClosed
	listener.listener = nil
	listener.err = net.ErrClosed
	listener.changed.Broadcast()
	listener.mu.Unlock()
	if inner == nil {
		return nil
	}
	return inner.Close()
}

func (listener *deferredListener) awaitOpen() (server.Listener, error) {
	listener.mu.Lock()
	defer listener.mu.Unlock()
	for {
		switch listener.state {
		case deferredListenerUnopened, deferredListenerOpening:
			listener.changed.Wait()
		case deferredListenerOpen:
			return listener.listener, nil
		case deferredListenerFailed, deferredListenerClosed:
			return nil, listener.err
		}
	}
}

func (listener *deferredListener) ensureOpen() (server.Listener, error) {
	listener.mu.Lock()
	for {
		switch listener.state {
		case deferredListenerUnopened:
			listener.state = deferredListenerOpening
			open := listener.open
			listener.mu.Unlock()
			inner, err := open()
			if err != nil && inner != nil {
				err = errors.Join(err, inner.Close())
				inner = nil
			}
			listener.mu.Lock()
			if err != nil || inner == nil {
				if err == nil {
					err = errors.New("listener factory returned nil")
				}
				listener.state = deferredListenerFailed
				listener.err = err
			} else {
				listener.state = deferredListenerOpen
				listener.listener = inner
			}
			listener.changed.Broadcast()
		case deferredListenerOpening:
			listener.changed.Wait()
		case deferredListenerOpen:
			inner := listener.listener
			listener.mu.Unlock()
			return inner, nil
		case deferredListenerFailed, deferredListenerClosed:
			err := listener.err
			listener.mu.Unlock()
			return nil, err
		}
	}
}
