package framework

import (
	"context"
	"math"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type customBlockClientListener struct {
	server.Listener
	blocks map[string]customBlockClientData
}

func (l *customBlockClientListener) Accept() (session.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return &customBlockClientConn{Conn: conn, blocks: l.blocks}, nil
}

func (l *customBlockClientListener) Disconnect(conn session.Conn, reason string) error {
	if wrapped, ok := conn.(*customBlockClientConn); ok {
		conn = wrapped.Conn
	}
	return l.Listener.Disconnect(conn, reason)
}

type customBlockClientConn struct {
	session.Conn
	blocks map[string]customBlockClientData
}

func (c *customBlockClientConn) StartGameContext(ctx context.Context, data minecraft.GameData) error {
	applyCustomBlockClientData(data.CustomBlocks, c.blocks)
	return c.Conn.StartGameContext(ctx, data)
}

func wrapCustomBlockClientListeners(config *server.Config, blocks map[string]customBlockClientData) {
	if len(blocks) == 0 {
		return
	}
	for index, factory := range config.Listeners {
		config.Listeners[index] = func(conf server.Config) (server.Listener, error) {
			listener, err := factory(conf)
			if err != nil {
				return nil, err
			}
			return &customBlockClientListener{Listener: listener, blocks: blocks}, nil
		}
	}
}

func applyCustomBlockClientData(entries []protocol.BlockEntry, blocks map[string]customBlockClientData) {
	for index := range entries {
		data, ok := blocks[entries[index].Name]
		if !ok {
			continue
		}
		components, ok := entries[index].Properties["components"].(map[string]any)
		if !ok {
			continue
		}

		delete(components, "minecraft:block_light_emission")
		delete(components, "minecraft:block_light_filter")
		if data.emission != 0 {
			components["minecraft:light_emission"] = map[string]any{"emission": data.emission}
		}
		if data.dampening != 15 {
			components["minecraft:light_dampening"] = map[string]any{"lightLevel": data.dampening}
		}

		// Dragonfly stores retained velocity (0.6 for normal blocks), while Bedrock's
		// component stores friction (0.4 for normal blocks).
		components["minecraft:friction"] = map[string]any{
			"value": float32(math.Min(0.9, math.Max(0, 1-data.friction))),
		}
	}
}
