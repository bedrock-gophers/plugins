package host

import "github.com/df-mc/dragonfly/server/world"

// WorldHandler is installed on every framework-owned world before the server starts.
// Event callbacks will be added as their generated ABI definitions land.
type WorldHandler struct {
	world.NopHandler
	entities *Entities
}

var _ world.Handler = (*WorldHandler)(nil)

func NewWorldHandler(entities *Entities) *WorldHandler {
	return &WorldHandler{entities: entities}
}

func (h *WorldHandler) HandleEntitySpawn(_ *world.Tx, entity world.Entity) {
	if h.entities != nil {
		h.entities.Register(entity)
	}
}

func (h *WorldHandler) HandleEntityDespawn(tx *world.Tx, entity world.Entity) {
	if h.entities == nil || entity == nil {
		return
	}
	handle := entity.H()
	h.entities.deactivateHandle(handle)
	tx.Defer(func(*world.Tx) {
		if handle.Closed() {
			h.entities.unregisterHandle(handle)
		}
	})
}
