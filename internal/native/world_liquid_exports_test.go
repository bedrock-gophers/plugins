package native

import (
	"bytes"
	"strings"
	"testing"
)

func TestWorldLiquidExportPayloadPreservesBlockData(t *testing.T) {
	block := WorldBlock{Identifier: "minecraft:water", PropertiesNBT: []byte{1, 2, 3}}
	identifier, properties, ok := worldBlockPayload(block, true)
	if !ok || !bytes.Equal(identifier, []byte(block.Identifier)) || !bytes.Equal(properties, block.PropertiesNBT) {
		t.Fatalf("payload = %q, %v, %v", identifier, properties, ok)
	}
}

func TestWorldLiquidExportPayloadReportsSmallBuffers(t *testing.T) {
	identifier, properties, ok := worldBlockPayload(WorldBlock{
		Identifier: "minecraft:flowing_water", PropertiesNBT: []byte{1, 2, 3, 4},
	}, true)
	if !ok {
		t.Fatal("valid payload rejected")
	}
	if worldBlockFits(identifier, properties, 2, 1) {
		t.Fatal("undersized buffers accepted")
	}
	if !worldBlockFits(identifier, properties, uint64(len(identifier)), uint64(len(properties))) {
		t.Fatal("exact buffers rejected")
	}
}

func TestWorldLiquidExportPayloadRejectsMissingAndMalformedValues(t *testing.T) {
	for name, test := range map[string]struct {
		block WorldBlock
		ok    bool
	}{
		"missing liquid":   {block: WorldBlock{Identifier: "minecraft:water"}},
		"empty identifier": {block: WorldBlock{}, ok: true},
		"invalid utf8":     {block: WorldBlock{Identifier: string([]byte{0xff})}, ok: true},
		"long identifier":  {block: WorldBlock{Identifier: strings.Repeat("x", maxBlockIdentifierBytes+1)}, ok: true},
		"long properties":  {block: WorldBlock{Identifier: "minecraft:water", PropertiesNBT: make([]byte, maxBlockPropertiesBytes+1)}, ok: true},
	} {
		t.Run(name, func(t *testing.T) {
			if _, _, ok := worldBlockPayload(test.block, test.ok); ok {
				t.Fatal("malformed payload accepted")
			}
		})
	}
}
