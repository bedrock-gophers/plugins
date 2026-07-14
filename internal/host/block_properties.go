package host

import (
	"bytes"
	"encoding/binary"
	"math"
	"sort"
	"unicode/utf8"
)

const (
	blockPropertyBool int32 = iota
	blockPropertyUint8
	blockPropertyInt32
	blockPropertyString
)

// EncodeBlockProperties preserves the concrete property types used by Dragonfly's block registry.
func EncodeBlockProperties(properties map[string]any) ([]byte, bool) {
	keys := make([]string, 0, len(properties))
	for key := range properties {
		if !validNBTString(key, false) {
			return nil, false
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var output bytes.Buffer
	output.Write([]byte{10, 0, 0}) // Root compound with an empty name.
	for _, key := range keys {
		if !writeBlockProperty(&output, key, properties[key]) {
			return nil, false
		}
	}
	output.WriteByte(0)
	return output.Bytes(), true
}

func writeBlockProperty(output *bytes.Buffer, key string, value any) bool {
	writeNBTTag(output, 10, key)
	writeNBTTag(output, 3, "kind")
	var kind [4]byte
	switch value := value.(type) {
	case bool:
		binary.LittleEndian.PutUint32(kind[:], uint32(blockPropertyBool))
		output.Write(kind[:])
		writeNBTTag(output, 1, "value")
		if value {
			output.WriteByte(1)
		} else {
			output.WriteByte(0)
		}
	case uint8:
		binary.LittleEndian.PutUint32(kind[:], uint32(blockPropertyUint8))
		output.Write(kind[:])
		writeNBTTag(output, 1, "value")
		output.WriteByte(value)
	case int32:
		binary.LittleEndian.PutUint32(kind[:], uint32(blockPropertyInt32))
		output.Write(kind[:])
		writeNBTTag(output, 3, "value")
		binary.LittleEndian.PutUint32(kind[:], uint32(value))
		output.Write(kind[:])
	case string:
		if !validNBTString(value, true) {
			return false
		}
		binary.LittleEndian.PutUint32(kind[:], uint32(blockPropertyString))
		output.Write(kind[:])
		writeNBTTag(output, 8, "value")
		writeNBTString(output, value)
	default:
		return false
	}
	output.WriteByte(0)
	return true
}

func writeNBTTag(output *bytes.Buffer, kind byte, name string) {
	output.WriteByte(kind)
	writeNBTString(output, name)
}

func writeNBTString(output *bytes.Buffer, value string) {
	var length [2]byte
	binary.LittleEndian.PutUint16(length[:], uint16(len(value)))
	output.Write(length[:])
	output.WriteString(value)
}

func validNBTString(value string, empty bool) bool {
	return (empty || value != "") && len(value) <= math.MaxUint16 && utf8.ValidString(value)
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
