package host

import (
	"github.com/bedrock-gophers/plugins/internal/blockstate"
)

const (
	blockPropertyBool int32 = iota
	blockPropertyUint8
	blockPropertyInt32
	blockPropertyString
)

// EncodeBlockProperties preserves the concrete property types used by Dragonfly's block registry.
func EncodeBlockProperties(properties map[string]any) ([]byte, bool) {
	return blockstate.EncodeProperties(properties)
}

// DecodeBlockProperties restores the concrete property types used by Dragonfly's block registry.
func DecodeBlockProperties(data []byte) (map[string]any, bool) {
	encoded, ok := UnmarshalNBT(data)
	if !ok {
		return nil, false
	}
	properties := make(map[string]any, len(encoded))
	for key, raw := range encoded {
		entry, ok := raw.(map[string]any)
		if !ok || key == "" || len(entry) != 2 {
			return nil, false
		}
		kind, ok := entry["kind"].(int32)
		if !ok {
			return nil, false
		}
		switch kind {
		case blockPropertyBool:
			value, ok := entry["value"].(uint8)
			if !ok || value > 1 {
				return nil, false
			}
			properties[key] = value == 1
		case blockPropertyUint8:
			value, ok := entry["value"].(uint8)
			if !ok {
				return nil, false
			}
			properties[key] = value
		case blockPropertyInt32:
			value, ok := entry["value"].(int32)
			if !ok {
				return nil, false
			}
			properties[key] = value
		case blockPropertyString:
			value, ok := entry["value"].(string)
			if !ok {
				return nil, false
			}
			properties[key] = value
		default:
			return nil, false
		}
	}
	return properties, true
}
