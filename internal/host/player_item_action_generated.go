// Code generated from Dragonfly server/player/player.go Go AST. DO NOT EDIT.
package host

import (
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
)

func runExactPlayerItemAction(connected *player.Player, stack item.Stack, kind native.PlayerItemActionKind) (int, bool, bool) {
	switch kind {
	case native.PlayerItemActionCollect:
		count, ok := connected.Collect(stack)
		return count, ok, true
	case native.PlayerItemActionDrop:
		return connected.Drop(stack), true, true
	default:
		return 0, false, false
	}
}
