package host

import (
	"testing"
	"time"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
)

func TestBucketItemsFromNative(t *testing.T) {
	tests := []struct {
		identifier string
		content    string
		liquid     bool
	}{
		{identifier: "minecraft:bucket"},
		{identifier: "minecraft:water_bucket", content: "water", liquid: true},
		{identifier: "minecraft:lava_bucket", content: "lava", liquid: true},
		{identifier: "minecraft:milk_bucket", content: "milk"},
	}
	for _, test := range tests {
		stack, ok := itemStackFromNative(native.ItemStack{Identifier: test.identifier, Count: 1})
		if !ok {
			t.Fatalf("decode %s", test.identifier)
		}
		bucket, ok := stack.Item().(item.Bucket)
		if !ok || bucket.Content.String() != test.content {
			t.Fatalf("%s decoded as %#v", test.identifier, stack.Item())
		}
		liquid, liquidOK := bucket.Content.Liquid()
		if liquidOK != test.liquid {
			t.Fatalf("%s liquid found=%v", test.identifier, liquidOK)
		}
		if liquidOK && liquid.LiquidType() != test.content {
			t.Fatalf("%s liquid type=%q", test.identifier, liquid.LiquidType())
		}
	}
}

func TestBucketItemsToNativeAndCapabilities(t *testing.T) {
	tests := []struct {
		bucket     item.Bucket
		identifier string
		maxCount   int
		fuel       time.Duration
		residue    bool
	}{
		{bucket: item.Bucket{}, identifier: "minecraft:bucket", maxCount: 16},
		{bucket: item.Bucket{Content: item.LiquidBucketContent(block.Water{})}, identifier: "minecraft:water_bucket", maxCount: 1},
		{bucket: item.Bucket{Content: item.LiquidBucketContent(block.Lava{})}, identifier: "minecraft:lava_bucket", maxCount: 1, fuel: 1000 * time.Second, residue: true},
		{bucket: item.Bucket{Content: item.MilkBucketContent()}, identifier: "minecraft:milk_bucket", maxCount: 1},
	}
	for _, test := range tests {
		encoded, ok := itemStackToNative(item.NewStack(test.bucket, 1))
		if !ok || encoded.Identifier != test.identifier || encoded.Metadata != 0 || encoded.Count != 1 {
			t.Fatalf("encode %s = %#v, %v", test.identifier, encoded, ok)
		}
		fuel := test.bucket.FuelInfo()
		if test.bucket.MaxCount() != test.maxCount || fuel.Duration != test.fuel || !fuel.Residue.Empty() != test.residue {
			t.Fatalf("%s capabilities: max=%d fuel=%s residue=%#v", test.identifier, test.bucket.MaxCount(), fuel.Duration, fuel.Residue)
		}
		if test.residue {
			residue, ok := fuel.Residue.Item().(item.Bucket)
			if !ok || !residue.Empty() || fuel.Residue.Count() != 1 {
				t.Fatalf("%s residue=%#v", test.identifier, fuel.Residue)
			}
		}
	}
}
