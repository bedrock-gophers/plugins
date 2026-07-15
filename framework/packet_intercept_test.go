package framework

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type packetEndpointFuncs struct {
	incoming func(*intercept.Context, packet.Packet)
	outgoing func(*intercept.Context, packet.Packet)
}

func (endpoint packetEndpointFuncs) HandleIncomingPacket(ctx *intercept.Context, value packet.Packet) {
	if endpoint.incoming != nil {
		endpoint.incoming(ctx, value)
	}
}

func (endpoint packetEndpointFuncs) HandleOutgoingPacket(ctx *intercept.Context, value packet.Packet) {
	if endpoint.outgoing != nil {
		endpoint.outgoing(ctx, value)
	}
}

func TestPacketInterceptRouterRoutesExactDirections(t *testing.T) {
	router := &packetInterceptRouter{}
	value := intercept.NopPacket{}
	var incoming, outgoing packet.Packet
	lease, err := router.activate(packetEndpointFuncs{
		incoming: func(ctx *intercept.Context, got packet.Packet) {
			incoming = got
			ctx.Cancel()
		},
		outgoing: func(_ *intercept.Context, got packet.Packet) { outgoing = got },
	})
	if err != nil {
		t.Fatal(err)
	}
	defer lease.Close()

	clientContext, serverContext := new(intercept.Context), new(intercept.Context)
	router.HandleClientPacket(clientContext, value)
	router.HandleServerPacket(serverContext, value)
	if incoming != value || outgoing != value {
		t.Fatalf("packets = (%T, %T), want (%T, %T)", incoming, outgoing, value, value)
	}
	if !clientContext.Cancelled() || serverContext.Cancelled() {
		t.Fatalf("cancelled = (%v, %v), want (true, false)", clientContext.Cancelled(), serverContext.Cancelled())
	}
}

func TestPacketInterceptRouterLeaseIsExclusiveAndDrains(t *testing.T) {
	router := &packetInterceptRouter{}
	entered, release := make(chan struct{}), make(chan struct{})
	var calls atomic.Int32
	lease, err := router.activate(packetEndpointFuncs{incoming: func(*intercept.Context, packet.Packet) {
		calls.Add(1)
		close(entered)
		<-release
	}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := router.activate(packetEndpointFuncs{}); !errors.Is(err, errPacketInterceptActive) {
		t.Fatalf("second activation error = %v", err)
	}

	dispatched := make(chan struct{})
	go func() {
		router.HandleClientPacket(new(intercept.Context), intercept.NopPacket{})
		close(dispatched)
	}()
	<-entered
	closed := make(chan struct{})
	go func() {
		lease.Close()
		close(closed)
	}()
	select {
	case <-closed:
		t.Fatal("lease closed before callback drained")
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	<-dispatched
	<-closed

	router.HandleClientPacket(new(intercept.Context), intercept.NopPacket{})
	if calls.Load() != 1 {
		t.Fatalf("calls after close = %d, want 1", calls.Load())
	}
	lease.Close()
	replacement, err := router.activate(packetEndpointFuncs{})
	if err != nil {
		t.Fatalf("activation after close: %v", err)
	}
	replacement.Close()
}

func TestPacketInterceptWrapsOutsideDeferredListener(t *testing.T) {
	rawConnection := &packetSessionConn{identity: login.IdentityData{Identity: "33e52a01-35fc-4c83-a0f9-b5d3f34995ed"}}
	rawListener := &packetRecordingListener{connection: rawConnection}
	config := server.Config{Listeners: []func(server.Config) (server.Listener, error){
		func(server.Config) (server.Listener, error) {
			rawListener.opens.Add(1)
			return rawListener, nil
		},
	}}
	gate := deferListenerCreation(&config)
	wrapPacketListeners(&config)
	listener, err := config.Listeners[0](config)
	if err != nil {
		t.Fatal(err)
	}
	if rawListener.opens.Load() != 0 {
		t.Fatal("packet wrapper eagerly opened the raw listener")
	}
	if err := gate.openAll(); err != nil {
		t.Fatal(err)
	}
	connection, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}
	wrapped, ok := connection.(*intercept.Conn)
	if !ok || wrapped.Conn != rawConnection {
		t.Fatalf("connection = %#v, want intercept wrapper around raw connection", connection)
	}
	if err := listener.Disconnect(connection, "test"); err != nil {
		t.Fatal(err)
	}
	if rawListener.disconnected() != rawConnection {
		t.Fatal("disconnect did not unwrap the intercepted connection")
	}
	if err := gate.closeAll(); err != nil {
		t.Fatal(err)
	}
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	if rawListener.closes.Load() != 1 {
		t.Fatalf("raw listener close count = %d, want 1", rawListener.closes.Load())
	}
}

type packetSessionConn struct {
	identity login.IdentityData
}

func (connection *packetSessionConn) IdentityData() login.IdentityData { return connection.identity }
func (*packetSessionConn) ClientData() login.ClientData                { return login.ClientData{} }
func (*packetSessionConn) ClientCacheEnabled() bool                    { return false }
func (*packetSessionConn) ChunkRadius() int                            { return 1 }
func (*packetSessionConn) Latency() time.Duration                      { return 0 }
func (*packetSessionConn) Flush() error                                { return nil }
func (*packetSessionConn) RemoteAddr() net.Addr                        { return &net.TCPAddr{} }
func (*packetSessionConn) ReadPacket() (packet.Packet, error)          { return nil, net.ErrClosed }
func (*packetSessionConn) WritePacket(packet.Packet) error             { return nil }
func (*packetSessionConn) StartGameContext(context.Context, minecraft.GameData) error {
	return nil
}
func (*packetSessionConn) Close() error { return nil }

type packetRecordingListener struct {
	connection session.Conn
	opens      atomic.Int32
	closes     atomic.Int32
	mu         sync.Mutex
	last       session.Conn
}

func (listener *packetRecordingListener) Accept() (session.Conn, error) {
	if listener.connection == nil {
		return nil, net.ErrClosed
	}
	connection := listener.connection
	listener.connection = nil
	return connection, nil
}

func (listener *packetRecordingListener) Disconnect(connection session.Conn, _ string) error {
	listener.mu.Lock()
	listener.last = connection
	listener.mu.Unlock()
	return nil
}

func (listener *packetRecordingListener) disconnected() session.Conn {
	listener.mu.Lock()
	defer listener.mu.Unlock()
	return listener.last
}

func (listener *packetRecordingListener) Close() error {
	listener.closes.Add(1)
	return nil
}
