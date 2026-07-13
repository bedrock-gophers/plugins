package host

const (
	blockPropertyBool int32 = iota
	blockPropertyUint8
	blockPropertyInt32
	blockPropertyString
)

// EncodeBlockProperties preserves the concrete property types used by Dragonfly's block registry.
func EncodeBlockProperties(properties map[string]any) ([]byte, bool) {
	encoded := make(map[string]any, len(properties))
	for key, value := range properties {
		if key == "" {
			return nil, false
		}
		switch value := value.(type) {
		case bool:
			encoded[key] = map[string]any{"kind": blockPropertyBool, "value": value}
		case uint8:
			encoded[key] = map[string]any{"kind": blockPropertyUint8, "value": value}
		case int32:
			encoded[key] = map[string]any{"kind": blockPropertyInt32, "value": value}
		case string:
			encoded[key] = map[string]any{"kind": blockPropertyString, "value": value}
		default:
			return nil, false
		}
	}
	return MarshalNBT(encoded)
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
