// Code generated from Dragonfly server/world/tx.go Go AST. DO NOT EDIT.
package framework

import (
	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/world"
)

func runExactWorldDefer(tx *world.Tx, callback func(*world.Tx) error, kind native.WorldDeferKind) (*world.Task, bool) {
	switch kind {
	case native.WorldDeferDefer:
		return tx.Defer(func(tx *world.Tx) {
			if err := callback(tx); err != nil {
				panic(err)
			}
		}), true
	case native.WorldDeferDeferErr:
		return tx.DeferErr(callback), true
	default:
		return nil, false
	}
}
