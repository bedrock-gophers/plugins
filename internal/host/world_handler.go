package host

import "github.com/df-mc/dragonfly/server/world"

// WorldHandler is installed on every framework-owned world before the server starts.
// Event callbacks will be added as their generated ABI definitions land.
type WorldHandler struct {
	world.NopHandler
}

var _ world.Handler = (*WorldHandler)(nil)

func NewWorldHandler() *WorldHandler {
	return &WorldHandler{}
}
