package framework

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var errPacketInterceptActive = errors.New("packet intercept already active")

// packetInterceptEndpoint is the run-owned half of the process-wide intercept
// router. Implementations must finish every callback before their lease is
// released.
type packetInterceptEndpoint interface {
	HandleIncomingPacket(*intercept.Context, packet.Packet)
	HandleOutgoingPacket(*intercept.Context, packet.Packet)
}

type nativePacketEndpoint struct {
	runtime *native.Runtime
	packets *host.Packets
	menus   *host.InventoryMenus
	log     *slog.Logger
}

func (endpoint nativePacketEndpoint) HandleIncomingPacket(ctx *intercept.Context, value packet.Packet) {
	if endpoint.menus != nil {
		endpoint.menus.HandleIncomingPacket(ctx, value)
	}
	endpoint.handle(ctx, value, native.PacketClientEvent, native.PacketClientSubscription)
}

func (endpoint nativePacketEndpoint) HandleOutgoingPacket(ctx *intercept.Context, value packet.Packet) {
	endpoint.handle(ctx, value, native.PacketServerEvent, native.PacketServerSubscription)
}

func (endpoint nativePacketEndpoint) handle(ctx *intercept.Context, value packet.Packet, event uint32, subscription uint64) {
	if ctx == nil || value == nil || endpoint.runtime == nil || endpoint.packets == nil || endpoint.runtime.Subscriptions()&subscription == 0 {
		return
	}
	handle, release, ok := endpoint.packets.Borrow(value, event == native.PacketClientEvent)
	if !ok {
		return
	}
	defer release()
	xuid := ""
	if connection := ctx.Val(); connection != nil {
		xuid = connection.IdentityData().XUID
	}
	cancelled, err := endpoint.runtime.HandlePacket(event, handle, value.ID(), xuid, ctx.Cancelled())
	if cancelled {
		ctx.Cancel()
	}
	if err != nil && endpoint.log != nil {
		endpoint.log.Error("plugin packet handler failed", "packet", value.ID(), "xuid", xuid, "error", err)
	}
}

var (
	packetInterceptHookOnce sync.Once
	processPacketIntercept  packetInterceptRouter
)

// wrapPacketListeners keeps intercept outside the deferred listener. The gate
// therefore continues to own the exact Dragonfly factories and their bind
// lifecycle, while intercept only decorates the connections returned to the
// server.
func wrapPacketListeners(config *server.Config) {
	ensurePacketInterceptHook()
	config.Listeners = intercept.WrapListeners(config.Listeners)
}

func ensurePacketInterceptHook() {
	packetInterceptHookOnce.Do(func() {
		intercept.Hook(&processPacketIntercept)
	})
}

// activatePacketIntercept binds the process-wide intercept package to one
// server run. intercept v0.3.0 stores its server and handlers globally, so two
// active runs cannot be routed correctly. The returned lease drains callbacks
// before detaching the endpoint.
func activatePacketIntercept(srv *server.Server, endpoint packetInterceptEndpoint) (*packetInterceptLease, error) {
	if srv == nil || endpoint == nil {
		return nil, errors.New("packet intercept requires a server and endpoint")
	}
	ensurePacketInterceptHook()
	lease, err := processPacketIntercept.activate(endpoint)
	if err != nil {
		return nil, err
	}
	// The exclusive endpoint lease is acquired first. A rejected concurrent
	// run must never replace intercept's process-wide server pointer.
	intercept.Start(srv)
	return lease, nil
}

type packetInterceptRouter struct {
	mu       sync.RWMutex
	endpoint *packetInterceptSlot
}

type packetInterceptSlot struct {
	endpoint packetInterceptEndpoint
}

func (router *packetInterceptRouter) activate(endpoint packetInterceptEndpoint) (*packetInterceptLease, error) {
	router.mu.Lock()
	defer router.mu.Unlock()
	if router.endpoint != nil {
		return nil, errPacketInterceptActive
	}
	slot := &packetInterceptSlot{endpoint: endpoint}
	router.endpoint = slot
	return &packetInterceptLease{router: router, slot: slot}, nil
}

// HandleClientPacket implements intercept.Handler exactly. Holding the read
// lease through the callback makes Close a synchronous drain barrier.
func (router *packetInterceptRouter) HandleClientPacket(ctx *intercept.Context, value packet.Packet) {
	router.mu.RLock()
	defer router.mu.RUnlock()
	if router.endpoint != nil {
		router.endpoint.endpoint.HandleIncomingPacket(ctx, value)
	}
}

// HandleServerPacket implements intercept.Handler exactly.
func (router *packetInterceptRouter) HandleServerPacket(ctx *intercept.Context, value packet.Packet) {
	router.mu.RLock()
	defer router.mu.RUnlock()
	if router.endpoint != nil {
		router.endpoint.endpoint.HandleOutgoingPacket(ctx, value)
	}
}

type packetInterceptLease struct {
	router *packetInterceptRouter
	slot   *packetInterceptSlot
	once   sync.Once
}

// Close stops new callbacks and waits for callbacks already using the endpoint.
func (lease *packetInterceptLease) Close() {
	if lease == nil {
		return
	}
	lease.once.Do(func() {
		lease.router.mu.Lock()
		if lease.router.endpoint == lease.slot {
			lease.router.endpoint = nil
		}
		lease.router.mu.Unlock()
	})
}
