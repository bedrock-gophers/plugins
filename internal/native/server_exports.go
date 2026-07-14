package native

/*
#include "bridge.h"
*/
import "C"

import (
	"unicode/utf8"
	"unsafe"
)

const maxPlayerXUIDBytes = 64

//export bg_go_server_players_open
func bg_go_server_players_open(context C.uint64_t, invocation C.DfInvocationId, output *C.DfPlayerIteratorId) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	if !ok || output == nil {
		return C.DF_STATUS_ERROR
	}
	id, ok := host.OpenServerPlayerIterator(InvocationID(invocation))
	if !ok || id == 0 {
		return C.DF_STATUS_ERROR
	}
	*output = C.DfPlayerIteratorId(id)
	return C.DF_STATUS_OK
}

//export bg_go_server_players_next
func bg_go_server_players_next(context C.uint64_t, invocation C.DfInvocationId, iterator C.DfPlayerIteratorId, playerInvocation *C.DfInvocationId, output *C.DfPlayerSnapshotBuffer, found *C.uint8_t) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	if !ok || iterator == 0 || playerInvocation == nil || output == nil || found == nil {
		return C.DF_STATUS_ERROR
	}
	*playerInvocation = 0
	*found = 0
	nested, snapshot, hasValue, valid := host.NextServerPlayer(InvocationID(invocation), PlayerIteratorID(iterator))
	if !valid || hasValue && (nested == 0 || snapshot.Player.Generation == 0) {
		return C.DF_STATUS_ERROR
	}
	if !hasValue {
		return C.DF_STATUS_OK
	}
	if !writePlayerSnapshotBuffer(output, snapshot) {
		host.CloseServerPlayers(InvocationID(invocation), PlayerIteratorID(iterator))
		return C.DF_STATUS_ERROR
	}
	*playerInvocation = C.DfInvocationId(nested)
	*found = 1
	return C.DF_STATUS_OK
}

//export bg_go_server_players_close
func bg_go_server_players_close(context C.uint64_t, invocation C.DfInvocationId, iterator C.DfPlayerIteratorId) {
	host, ok := resolveHost(uint64(context))
	if ok && iterator != 0 {
		host.CloseServerPlayers(InvocationID(invocation), PlayerIteratorID(iterator))
	}
}

//export bg_go_server_player
func bg_go_server_player(context C.uint64_t, value C.DfUuid, output *C.DfEntityHandleId, found *C.uint8_t) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	if !ok || output == nil || found == nil {
		return C.DF_STATUS_ERROR
	}
	var id [16]byte
	for index := range id {
		id[index] = byte(value.bytes[index])
	}
	player, hasValue, valid := host.ServerPlayer(id)
	return writeServerPlayerLookup(player, hasValue, valid, output, found)
}

//export bg_go_server_player_by_name
func bg_go_server_player_by_name(context C.uint64_t, value C.DfStringView, output *C.DfEntityHandleId, found *C.uint8_t) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	name, valid := copyServerPlayerName(value)
	if !ok || !valid || output == nil || found == nil {
		return C.DF_STATUS_ERROR
	}
	player, hasValue, lookupValid := host.ServerPlayerByName(name)
	return writeServerPlayerLookup(player, hasValue, lookupValid, output, found)
}

//export bg_go_server_max_player_count
func bg_go_server_max_player_count(context C.uint64_t, output *C.int64_t) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	if !ok || output == nil {
		return C.DF_STATUS_ERROR
	}
	count, ok := host.ServerMaxPlayerCount()
	if !ok || count < 0 {
		return C.DF_STATUS_ERROR
	}
	*output = C.int64_t(count)
	return C.DF_STATUS_OK
}

//export bg_go_server_player_count
func bg_go_server_player_count(context C.uint64_t, output *C.int64_t) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	if !ok || output == nil {
		return C.DF_STATUS_ERROR
	}
	count, ok := host.ServerPlayerCount()
	if !ok || count < 0 {
		return C.DF_STATUS_ERROR
	}
	*output = C.int64_t(count)
	return C.DF_STATUS_OK
}

//export bg_go_server_player_by_xuid
func bg_go_server_player_by_xuid(context C.uint64_t, value C.DfStringView, output *C.DfEntityHandleId, found *C.uint8_t) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	xuid, valid := copyBoundedServerString(value, maxPlayerXUIDBytes)
	if !ok || !valid || output == nil || found == nil {
		return C.DF_STATUS_ERROR
	}
	player, hasValue, lookupValid := host.ServerPlayerByXUID(xuid)
	return writeServerPlayerLookup(player, hasValue, lookupValid, output, found)
}

//export bg_go_player_xuid
func bg_go_player_xuid(context C.uint64_t, invocation C.DfInvocationId, player C.DfPlayerId, output *C.DfStringBuffer) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	if !ok || output == nil {
		return C.DF_STATUS_ERROR
	}
	xuid, ok := host.PlayerXUID(InvocationID(invocation), playerID(player))
	if !ok || len(xuid) > maxPlayerXUIDBytes || !utf8.ValidString(xuid) || !writeSkinBuffer(output, []byte(xuid)) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}

func writeServerPlayerLookup(id EntityHandleID, hasValue, valid bool, output *C.DfEntityHandleId, found *C.uint8_t) C.DfStatus {
	*output = C.DfEntityHandleId{}
	*found = 0
	if !valid || hasValue && !id.Valid() {
		return C.DF_STATUS_ERROR
	}
	if hasValue {
		*output = cEntityHandleID(id)
		*found = 1
	}
	return C.DF_STATUS_OK
}

func copyServerPlayerName(view C.DfStringView) (string, bool) {
	return copyBoundedServerString(view, maxPlayerNameBytes)
}

func copyBoundedServerString(view C.DfStringView, maximum uint64) (string, bool) {
	if uint64(view.len) > maximum || view.len != 0 && view.data == nil {
		return "", false
	}
	value := C.GoBytes(unsafe.Pointer(view.data), C.int(view.len))
	return string(value), utf8.Valid(value)
}
