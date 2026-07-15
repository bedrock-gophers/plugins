package native

/*
#include "dragonfly_plugin.h"
*/
import "C"

import (
	"time"
	"unicode/utf8"
)

const maxItemCooldownIdentifierBytes = 256

//export bg_go_player_cooldown
func bg_go_player_cooldown(
	context C.uint64_t,
	invocation C.DfInvocationId,
	player C.DfPlayerId,
	operation C.uint32_t,
	identifier C.DfStringView,
	metadata C.int32_t,
	duration C.int64_t,
	active *C.uint8_t,
) C.DfStatus {
	if active == nil || uint32(operation) > uint32(PlayerCooldownSet) {
		return C.DF_STATUS_ERROR
	}
	*active = 0
	host, ok := resolveHost(uint64(context))
	identifierBytes, valid := copyNativeBytes(identifier, maxItemCooldownIdentifierBytes)
	if !ok || !valid || len(identifierBytes) == 0 || !utf8.Valid(identifierBytes) {
		return C.DF_STATUS_ERROR
	}
	value, ok := host.PlayerCooldown(
		InvocationID(invocation), playerID(player), PlayerCooldownOperation(operation),
		string(identifierBytes), int32(metadata), time.Duration(duration),
	)
	if !ok {
		return C.DF_STATUS_ERROR
	}
	*active = C.uint8_t(boolByte(value))
	return C.DF_STATUS_OK
}
