package host

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBlockPropertiesEncodeCanonically(t *testing.T) {
	properties := map[string]any{
		"bool": true, "byte": uint8(7), "int": int32(-9), "string": "north",
	}
	want, ok := EncodeBlockProperties(properties)
	if !ok {
		t.Fatal("valid properties rejected")
	}
	for range 100 {
		got, ok := EncodeBlockProperties(properties)
		if !ok || !bytes.Equal(got, want) {
			t.Fatalf("encoding changed: %x != %x", got, want)
		}
	}
	decoded, ok := DecodeBlockProperties(want)
	if !ok || !reflect.DeepEqual(decoded, properties) {
		t.Fatalf("decoded = %#v, %v", decoded, ok)
	}
}

func TestBlockPropertiesRejectInvalidStringsAndTypes(t *testing.T) {
	for name, properties := range map[string]map[string]any{
		"empty key":     {"": true},
		"invalid key":   {string([]byte{0xff}): true},
		"invalid value": {"direction": string([]byte{0xff})},
		"unknown type":  {"age": int(1)},
	} {
		t.Run(name, func(t *testing.T) {
			if _, ok := EncodeBlockProperties(properties); ok {
				t.Fatal("invalid properties accepted")
			}
		})
	}
}
