package framework

import (
	"errors"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/session"
)

func TestDeferredListenerGateOpensInFactoryOrder(t *testing.T) {
	recorder := &listenerRecorder{}
	config := server.Config{Listeners: []func(server.Config) (server.Listener, error){
		recordingListenerFactory(recorder, "one", nil),
		recordingListenerFactory(recorder, "two", nil),
		recordingListenerFactory(recorder, "three", nil),
	}}
	gate := deferListenerCreation(&config)
	listeners := constructDeferredListeners(t, config)
	if got := recorder.eventsCopy(); len(got) != 0 {
		t.Fatalf("events before openAll = %v", got)
	}
	if err := gate.openAll(); err != nil {
		t.Fatal(err)
	}
	if got, want := recorder.eventsCopy(), []string{"open one", "open two", "open three"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("events = %v, want %v", got, want)
	}

	if err := gate.closeAll(); err != nil {
		t.Fatal(err)
	}
	for _, listener := range listeners {
		if err := listener.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDeferredListenerGateRollsBackPartialOpenInReverse(t *testing.T) {
	recorder := &listenerRecorder{}
	wantErr := errors.New("bind three")
	config := server.Config{Listeners: []func(server.Config) (server.Listener, error){
		recordingListenerFactory(recorder, "one", nil),
		recordingListenerFactory(recorder, "two", nil),
		recordingListenerFactory(recorder, "three", wantErr),
		recordingListenerFactory(recorder, "four", nil),
	}}
	gate := deferListenerCreation(&config)
	constructDeferredListeners(t, config)

	if err := gate.openAll(); !errors.Is(err, wantErr) {
		t.Fatalf("openAll error = %v", err)
	}
	if got, want := recorder.eventsCopy(), []string{
		"open one", "open two", "open three", "close two", "close one",
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("events = %v, want %v", got, want)
	}
	if err := gate.closeAll(); err != nil {
		t.Fatal(err)
	}
	if got := recorder.eventsCopy(); !reflect.DeepEqual(got, []string{
		"open one", "open two", "open three", "close two", "close one",
	}) {
		t.Fatalf("idempotent close events = %v", got)
	}
}

func TestDeferredListenerGateCloseIsReverseAndIdempotent(t *testing.T) {
	recorder := &listenerRecorder{}
	config := server.Config{Listeners: []func(server.Config) (server.Listener, error){
		recordingListenerFactory(recorder, "one", nil),
		recordingListenerFactory(recorder, "two", nil),
		recordingListenerFactory(recorder, "three", nil),
	}}
	gate := deferListenerCreation(&config)
	constructDeferredListeners(t, config)
	if err := gate.openAll(); err != nil {
		t.Fatal(err)
	}
	if err := gate.closeAll(); err != nil {
		t.Fatal(err)
	}
	if err := gate.closeAll(); err != nil {
		t.Fatal(err)
	}
	if got, want := recorder.eventsCopy(), []string{
		"open one", "open two", "open three", "close three", "close two", "close one",
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("events = %v, want %v", got, want)
	}
}

func TestDeferredListenerAcceptWaitsForExplicitOpen(t *testing.T) {
	var opens atomic.Int32
	inner := newRecordingListener()
	config := server.Config{Listeners: []func(server.Config) (server.Listener, error){
		func(server.Config) (server.Listener, error) {
			opens.Add(1)
			return inner, nil
		},
	}}
	gate := deferListenerCreation(&config)
	listener := constructDeferredListeners(t, config)[0]
	done := make(chan error, 1)
	go func() {
		_, err := listener.Accept()
		done <- err
	}()
	if opens.Load() != 0 {
		t.Fatal("Accept opened the real listener")
	}
	select {
	case err := <-done:
		t.Fatalf("Accept returned before openAll: %v", err)
	case <-time.After(20 * time.Millisecond):
	}

	if err := gate.openAll(); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool { return opens.Load() == 1 })
	if err := gate.closeAll(); err != nil {
		t.Fatal(err)
	}
	if err := <-done; !errors.Is(err, net.ErrClosed) {
		t.Fatalf("Accept error = %v", err)
	}
	if inner.closes.Load() != 1 {
		t.Fatalf("inner close count = %d", inner.closes.Load())
	}
}

func TestDeferredListenerGateCloseWaitsForConcurrentOpen(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	inner := newRecordingListener()
	config := server.Config{Listeners: []func(server.Config) (server.Listener, error){
		func(server.Config) (server.Listener, error) {
			close(started)
			<-release
			return inner, nil
		},
	}}
	gate := deferListenerCreation(&config)
	constructDeferredListeners(t, config)
	openDone := make(chan error, 1)
	go func() { openDone <- gate.openAll() }()
	<-started
	closeDone := make(chan error, 1)
	go func() { closeDone <- gate.closeAll() }()
	select {
	case err := <-closeDone:
		t.Fatalf("closeAll returned during open: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	if err := <-openDone; err != nil {
		t.Fatal(err)
	}
	if err := <-closeDone; err != nil {
		t.Fatal(err)
	}
	if inner.closes.Load() != 1 {
		t.Fatalf("inner close count = %d", inner.closes.Load())
	}
}

func TestDeferredListenerCloseBeforeOpenNeverBinds(t *testing.T) {
	var opens atomic.Int32
	config := server.Config{Listeners: []func(server.Config) (server.Listener, error){
		func(server.Config) (server.Listener, error) {
			opens.Add(1)
			return newRecordingListener(), nil
		},
	}}
	gate := deferListenerCreation(&config)
	listener := constructDeferredListeners(t, config)[0]
	if err := gate.closeAll(); err != nil {
		t.Fatal(err)
	}
	if _, err := listener.Accept(); !errors.Is(err, net.ErrClosed) {
		t.Fatalf("Accept error = %v", err)
	}
	if opens.Load() != 0 {
		t.Fatalf("listener opened %d times", opens.Load())
	}
}

func constructDeferredListeners(t *testing.T, config server.Config) []server.Listener {
	t.Helper()
	listeners := make([]server.Listener, 0, len(config.Listeners))
	for _, factory := range config.Listeners {
		listener, err := factory(config)
		if err != nil {
			t.Fatal(err)
		}
		listeners = append(listeners, listener)
	}
	return listeners
}

type listenerRecorder struct {
	mu     sync.Mutex
	events []string
}

func (recorder *listenerRecorder) record(event string) {
	recorder.mu.Lock()
	recorder.events = append(recorder.events, event)
	recorder.mu.Unlock()
}

func (recorder *listenerRecorder) eventsCopy() []string {
	recorder.mu.Lock()
	defer recorder.mu.Unlock()
	return append([]string(nil), recorder.events...)
}

func recordingListenerFactory(recorder *listenerRecorder, name string, openErr error) func(server.Config) (server.Listener, error) {
	return func(server.Config) (server.Listener, error) {
		recorder.record("open " + name)
		if openErr != nil {
			return nil, openErr
		}
		return &orderedRecordingListener{name: name, recorder: recorder}, nil
	}
}

type orderedRecordingListener struct {
	name     string
	recorder *listenerRecorder
	closed   atomic.Bool
}

func (*orderedRecordingListener) Accept() (session.Conn, error)         { return nil, net.ErrClosed }
func (*orderedRecordingListener) Disconnect(session.Conn, string) error { return nil }
func (listener *orderedRecordingListener) Close() error {
	if listener.closed.CompareAndSwap(false, true) {
		listener.recorder.record("close " + listener.name)
	}
	return nil
}

func waitFor(t *testing.T, ready func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !ready() {
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for listener state")
		}
		time.Sleep(time.Millisecond)
	}
}

type recordingListener struct {
	closed chan struct{}
	closes atomic.Int32
}

func newRecordingListener() *recordingListener {
	return &recordingListener{closed: make(chan struct{})}
}

func (listener *recordingListener) Accept() (session.Conn, error) {
	<-listener.closed
	return nil, net.ErrClosed
}

func (*recordingListener) Disconnect(session.Conn, string) error { return nil }

func (listener *recordingListener) Close() error {
	if listener.closes.Add(1) == 1 {
		close(listener.closed)
	}
	return nil
}
