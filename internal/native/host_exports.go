package native

/*
#include "bridge.h"
*/
import "C"

//export bg_go_player_message
func bg_go_player_message(context C.uint64_t, player C.DfPlayerId, message C.DfStringView) C.DfStatus {
	host, ok := resolveHost(uint64(context))
	if !ok {
		return C.DF_STATUS_ERROR
	}
	var id PlayerID
	for index := range id.UUID {
		id.UUID[index] = byte(player.bytes[index])
	}
	id.Generation = uint64(player.generation)
	if !host.MessagePlayer(id, stringView(message)) {
		return C.DF_STATUS_ERROR
	}
	return C.DF_STATUS_OK
}
