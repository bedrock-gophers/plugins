package framework

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bedrock-gophers/plugins/internal/host"
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type transferSessionConn struct {
	identity login.IdentityData
	closed   chan struct{}
	close    sync.Once
	mu       sync.Mutex
	writes   []packet.Packet
}

func newTransferSessionConn(id uuid.UUID) *transferSessionConn {
	return &transferSessionConn{
		identity: login.IdentityData{Identity: id.String(), DisplayName: "SessionTraveller"},
		closed:   make(chan struct{}),
	}
}

func (c *transferSessionConn) IdentityData() login.IdentityData { return c.identity }
func (*transferSessionConn) ClientData() login.ClientData {
	return login.ClientData{LanguageCode: "en_US"}
}
func (*transferSessionConn) ClientCacheEnabled() bool { return false }
func (*transferSessionConn) ChunkRadius() int         { return 1 }
func (*transferSessionConn) Latency() time.Duration   { return 0 }
func (*transferSessionConn) Flush() error             { return nil }
func (*transferSessionConn) RemoteAddr() net.Addr     { return &net.TCPAddr{} }
func (c *transferSessionConn) ReadPacket() (packet.Packet, error) {
	<-c.closed
	return nil, net.ErrClosed
}
func (c *transferSessionConn) WritePacket(value packet.Packet) error {
	c.mu.Lock()
	c.writes = append(c.writes, value)
	c.mu.Unlock()
	return nil
}
func (*transferSessionConn) StartGameContext(context.Context, minecraft.GameData) error { return nil }
func (c *transferSessionConn) Close() error {
	c.close.Do(func() { close(c.closed) })
	return nil
}
func (c *transferSessionConn) wroteChangeDimension() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, value := range c.writes {
		if _, ok := value.(*packet.ChangeDimension); ok {
			return true
		}
	}
	return false
}

type transferSessionQuitHandler struct {
	player.NopHandler
	players *host.Players
	quits   atomic.Int32
}

func (h *transferSessionQuitHandler) HandleQuit(connected *player.Player) {
	h.quits.Add(1)
	h.players.Unregister(connected)
}

func TestTransferPlayerSessionSwitchesLoaderAndClosesWhileWorldless(t *testing.T) {
	players := host.NewPlayers()
	manager := newWorldManager("", nil, players)
	source := world.Config{}.New()
	if err := manager.RegisterCore(OverworldID, source); err != nil {
		t.Fatal(err)
	}
	target, err := manager.Create("example:session-nether", world.Config{Dim: world.Nether})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = manager.CloseCustom()
		_ = source.Close()
	})
	sourceID, _ := manager.WorldByName(0, string(OverworldID))
	targetID, _ := manager.WorldByName(0, "example:session-nether")
	playerUUID := uuid.MustParse("68b74323-0008-4cc8-9f99-c31bcb487a63")
	connection := newTransferSessionConn(playerUUID)
	var stops atomic.Int32
	networkSession := session.Config{
		MaxChunkRadius: 1,
		HandleStop: func(*world.Tx, session.Controllable) {
			stops.Add(1)
		},
	}.New(connection)
	handle := world.EntitySpawnOpts{ID: playerUUID}.New(player.Type, player.Config{
		UUID: playerUUID, Name: "SessionTraveller", Session: networkSession,
	})
	networkSession.SetHandle(handle, skin.Skin{})
	quit := &transferSessionQuitHandler{players: players}
	var id native.PlayerID
	if err := source.Do(func(tx *world.Tx) {
		connected := tx.AddEntity(handle).(*player.Player)
		id = players.Register(connected, 96)
		connected.Handle(quit)
		networkSession.Spawn(connected, tx)
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !manager.TransferPlayer(0, id, targetID, native.Vec3{X: 3, Y: 70, Z: 4}) {
		t.Fatal("session transfer to nether was rejected")
	}
	if err := world.NewEntityRef[*player.Player](handle).Do(func(tx *world.Tx, connected *player.Player) {
		if tx.World() != target || connected.Position() != ([3]float64{3, 70, 4}) {
			t.Errorf("session destination world=%p position=%v", tx.World(), connected.Position())
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for !connection.wroteChangeDimension() {
		if time.Now().After(deadline) {
			t.Fatal("session loader did not emit a dimension switch")
		}
		time.Sleep(time.Millisecond)
	}

	entered := make(chan struct{})
	release := make(chan struct{})
	barrier := source.Do(func(*world.Tx) {
		close(entered)
		<-release
	})
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("source barrier did not start")
	}
	if !manager.TransferPlayer(0, id, sourceID, native.Vec3{Y: 65}) {
		t.Fatal("session return transfer was rejected")
	}
	if err := target.Do(func(tx *world.Tx) {
		if _, ok := handle.Entity(tx); ok {
			t.Error("session player remained in target while source add was blocked")
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
	close(release)
	if err := barrier.Wait(context.Background()); err != nil && !errors.Is(err, world.ErrWorldClosed) {
		t.Fatal(err)
	}
	deadline = time.Now().Add(2 * time.Second)
	for quit.quits.Load() != 1 || stops.Load() != 1 || !handle.Closed() {
		if time.Now().After(deadline) {
			t.Fatalf("session teardown quits=%d stops=%d closed=%v", quit.quits.Load(), stops.Load(), handle.Closed())
		}
		time.Sleep(time.Millisecond)
	}
	if _, ok := players.Handle(id); ok {
		t.Fatal("session player remained registered after worldless connection close")
	}
}
