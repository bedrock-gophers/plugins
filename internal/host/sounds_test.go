package host

import (
	"context"
	"math"
	"reflect"
	"testing"

	"github.com/bedrock-gophers/plugins/internal/native"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
)

func TestSoundFromNativeCoversEveryDragonflySound(t *testing.T) {
	w := world.Config{Synchronous: true, Entities: entity.DefaultRegistry}.New()
	t.Cleanup(func() { _ = w.Close() })
	properties, ok := EncodeBlockProperties(map[string]any{})
	if !ok {
		t.Fatal("encode block properties")
	}
	block := &native.WorldBlock{Identifier: "minecraft:stone", PropertiesNBT: properties}
	item := &native.ItemStack{Identifier: "minecraft:diamond_sword", Count: 1}
	if err := w.Do(func(tx *world.Tx) {
		for kind := native.SoundAnvilBreak; kind <= native.SoundGoatHorn; kind++ {
			value := testSoundDescriptor(kind, block, item)
			decoded, ok := SoundFromNative(tx, value)
			if !ok || decoded == nil {
				t.Fatalf("kind %d decoded as %T, ok=%v", kind, decoded, ok)
			}
			assertParameterizedSound(t, kind, decoded)
		}
		state, ok := EncodeBlockProperties(map[string]any{"pillar_axis": "x"})
		if !ok {
			t.Fatal("encode stateful block properties")
		}
		decoded, ok := SoundFromNative(tx, native.WorldSound{
			Kind:  native.SoundDoorOpen,
			Block: &native.WorldBlock{Identifier: "minecraft:oak_log", PropertiesNBT: state},
		})
		got, typed := decoded.(sound.DoorOpen)
		if !ok || !typed {
			t.Fatalf("stateful door sound decoded as %T, ok=%v", decoded, ok)
		}
		name, properties := got.Block.EncodeBlock()
		if name != "minecraft:oak_log" || properties["pillar_axis"] != "x" {
			t.Fatalf("stateful sound block = %s %#v", name, properties)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestSoundToNativeCoversEverySoundFromNative(t *testing.T) {
	w := world.Config{Synchronous: true, Entities: entity.DefaultRegistry}.New()
	t.Cleanup(func() { _ = w.Close() })
	properties, ok := EncodeBlockProperties(map[string]any{})
	if !ok {
		t.Fatal("encode block properties")
	}
	stone := &native.WorldBlock{Identifier: "minecraft:stone", PropertiesNBT: properties}
	sword := &native.ItemStack{Identifier: "minecraft:diamond_sword", Count: 1}
	if err := w.Do(func(tx *world.Tx) {
		for kind := native.SoundAnvilBreak; kind <= native.SoundGoatHorn; kind++ {
			want := testSoundDescriptor(kind, stone, sword)
			decoded, ok := SoundFromNative(tx, want)
			if !ok {
				t.Fatalf("kind %d did not decode", kind)
			}
			got, ok := SoundToNative(tx, decoded)
			if !ok {
				t.Fatalf("kind %d (%T) did not encode", kind, decoded)
			}
			pointer := reflect.New(reflect.TypeOf(decoded))
			pointer.Elem().Set(reflect.ValueOf(decoded))
			pointerSound, ok := pointer.Interface().(world.Sound)
			if !ok {
				t.Fatalf("kind %d pointer %T does not implement world.Sound", kind, pointer.Interface())
			}
			pointerGot, ok := SoundToNative(tx, pointerSound)
			if !ok || !reflect.DeepEqual(pointerGot, got) {
				t.Fatalf("kind %d pointer mismatch\nvalue:   %#v\npointer: %#v, ok=%v", kind, got, pointerGot, ok)
			}
			if kind == native.SoundBucketFill || kind == native.SoundBucketEmpty {
				if got.Kind != want.Kind || got.Data != want.Data || got.Block == nil {
					t.Fatalf("kind %d (%T) bucket descriptor = %#v", kind, decoded, got)
				}
				continue
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("kind %d (%T) mismatch\ngot:  %#v\nwant: %#v", kind, decoded, got, want)
			}
		}

		state, ok := EncodeBlockProperties(map[string]any{"pillar_axis": "x"})
		if !ok {
			t.Fatal("encode stateful block properties")
		}
		want := native.WorldSound{
			Kind:  native.SoundDoorOpen,
			Block: &native.WorldBlock{Identifier: "minecraft:oak_log", PropertiesNBT: state},
		}
		decoded, ok := SoundFromNative(tx, want)
		if !ok {
			t.Fatal("decode stateful block sound")
		}
		got, ok := SoundToNative(tx, decoded)
		if !ok || got.Kind != want.Kind || got.Block == nil || got.Block.Identifier != want.Block.Identifier {
			t.Fatalf("stateful block sound = %#v, ok=%v", got, ok)
		}
		gotProperties, gotOK := DecodeBlockProperties(got.Block.PropertiesNBT)
		wantProperties, wantOK := DecodeBlockProperties(want.Block.PropertiesNBT)
		if !gotOK || !wantOK || !reflect.DeepEqual(gotProperties, wantProperties) {
			t.Fatalf("stateful block properties = %#v, want %#v", gotProperties, wantProperties)
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestSoundToNativePreservesBucketLiquidState(t *testing.T) {
	w := world.Config{Synchronous: true, Entities: entity.DefaultRegistry}.New()
	t.Cleanup(func() { _ = w.Close() })
	tests := []struct {
		name  string
		sound world.Sound
		kind  native.SoundKind
		data  uint32
		block world.Block
	}{
		{name: "flowing water", sound: sound.BucketFill{Liquid: block.Water{Depth: 3, Falling: true}}, kind: native.SoundBucketFill, block: block.Water{Depth: 3, Falling: true}},
		{name: "still lava", sound: sound.BucketEmpty{Liquid: block.Lava{Depth: 6, Still: true}}, kind: native.SoundBucketEmpty, data: 1, block: block.Lava{Depth: 6, Still: true}},
	}
	if err := w.Do(func(tx *world.Tx) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got, ok := SoundToNative(tx, test.sound)
				if !ok || got.Kind != test.kind || got.Data != test.data || got.Block == nil {
					t.Fatalf("encoded bucket sound = %#v, ok=%v", got, ok)
				}
				identifier, properties := test.block.EncodeBlock()
				gotProperties, propertiesOK := DecodeBlockProperties(got.Block.PropertiesNBT)
				if !propertiesOK || got.Block.Identifier != identifier || !reflect.DeepEqual(gotProperties, properties) {
					t.Fatalf("encoded liquid = %s %#v, want %s %#v", got.Block.Identifier, gotProperties, identifier, properties)
				}
			})
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestSoundToNativeRejectsUnsupportedValues(t *testing.T) {
	w := world.Config{Synchronous: true, Entities: entity.DefaultRegistry}.New()
	t.Cleanup(func() { _ = w.Close() })
	if _, ok := SoundToNative(nil, sound.LevelUp{}); ok {
		t.Fatal("accepted nil transaction")
	}
	if err := w.Do(func(tx *world.Tx) {
		for _, value := range []world.Sound{
			nil,
			(*sound.LevelUp)(nil),
			sound.Fall{Distance: math.NaN()},
			sound.DecoratedPotInserted{Progress: math.Inf(1)},
			sound.BlockPlace{},
			sound.EquipItem{},
			sound.BucketFill{},
			sound.BucketEmpty{},
			sound.CrossbowLoad{Stage: sound.CrossbowLoadingEnd + 1},
		} {
			if encoded, ok := SoundToNative(tx, value); ok {
				t.Fatalf("encoded unsupported %T as %#v", value, encoded)
			}
		}
	}).Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func testSoundDescriptor(kind native.SoundKind, block *native.WorldBlock, item *native.ItemStack) native.WorldSound {
	value := native.WorldSound{Kind: kind}
	switch kind {
	case native.SoundAttack:
		value.Flags = 1
	case native.SoundFall:
		value.Scalar = 5.5
	case native.SoundBlockPlace, native.SoundBlockBreaking, native.SoundDoorOpen, native.SoundDoorClose,
		native.SoundTrapdoorOpen, native.SoundTrapdoorClose, native.SoundFenceGateOpen, native.SoundFenceGateClose,
		native.SoundItemUseOn:
		value.Block = block
	case native.SoundNote:
		value.Data, value.Integer = 15, 30
	case native.SoundMusicDiscPlay:
		value.Data = 20
	case native.SoundDecoratedPotInserted:
		value.Scalar = 0.75
	case native.SoundEquipItem:
		value.Item = item
	case native.SoundBucketFill, native.SoundBucketEmpty:
		value.Data = 1
	case native.SoundCrossbowLoad:
		value.Integer, value.Flags = sound.CrossbowLoadingEnd, 1
	case native.SoundGoatHorn:
		value.Data = 7
	}
	return value
}

func assertParameterizedSound(t *testing.T, kind native.SoundKind, decoded world.Sound) {
	t.Helper()
	switch kind {
	case native.SoundAttack:
		value, ok := decoded.(sound.Attack)
		if !ok || !value.Damage {
			t.Fatalf("attack = %#v", decoded)
		}
	case native.SoundFall:
		value, ok := decoded.(sound.Fall)
		if !ok || value.Distance != 5.5 {
			t.Fatalf("fall = %#v", decoded)
		}
	case native.SoundBlockPlace:
		_, ok := decoded.(sound.BlockPlace)
		if !ok {
			t.Fatalf("block place = %T", decoded)
		}
	case native.SoundBlockBreaking:
		_, ok := decoded.(sound.BlockBreaking)
		if !ok {
			t.Fatalf("block breaking = %T", decoded)
		}
	case native.SoundDoorOpen:
		_, ok := decoded.(sound.DoorOpen)
		if !ok {
			t.Fatalf("door open = %T", decoded)
		}
	case native.SoundDoorClose:
		_, ok := decoded.(sound.DoorClose)
		if !ok {
			t.Fatalf("door close = %T", decoded)
		}
	case native.SoundTrapdoorOpen:
		_, ok := decoded.(sound.TrapdoorOpen)
		if !ok {
			t.Fatalf("trapdoor open = %T", decoded)
		}
	case native.SoundTrapdoorClose:
		_, ok := decoded.(sound.TrapdoorClose)
		if !ok {
			t.Fatalf("trapdoor close = %T", decoded)
		}
	case native.SoundFenceGateOpen:
		_, ok := decoded.(sound.FenceGateOpen)
		if !ok {
			t.Fatalf("fence gate open = %T", decoded)
		}
	case native.SoundFenceGateClose:
		_, ok := decoded.(sound.FenceGateClose)
		if !ok {
			t.Fatalf("fence gate close = %T", decoded)
		}
	case native.SoundNote:
		value, ok := decoded.(sound.Note)
		if !ok || value.Instrument != sound.Pling() || value.Pitch != 30 {
			t.Fatalf("note = %#v", decoded)
		}
	case native.SoundMusicDiscPlay:
		value, ok := decoded.(sound.MusicDiscPlay)
		if !ok || value.DiscType != sound.DiscLavaChicken() {
			t.Fatalf("music disc = %#v", decoded)
		}
	case native.SoundDecoratedPotInserted:
		value, ok := decoded.(sound.DecoratedPotInserted)
		if !ok || value.Progress != 0.75 {
			t.Fatalf("decorated pot = %#v", decoded)
		}
	case native.SoundItemUseOn:
		_, ok := decoded.(sound.ItemUseOn)
		if !ok {
			t.Fatalf("item use on = %T", decoded)
		}
	case native.SoundEquipItem:
		value, ok := decoded.(sound.EquipItem)
		if !ok {
			t.Fatalf("equip item = %T", decoded)
		}
		name, _ := value.Item.EncodeItem()
		if name != "minecraft:diamond_sword" {
			t.Fatalf("equip item = %#v", decoded)
		}
	case native.SoundBucketFill:
		value, ok := decoded.(sound.BucketFill)
		_, lava := value.Liquid.(block.Lava)
		if !ok || !lava {
			t.Fatalf("bucket fill = %#v", decoded)
		}
	case native.SoundBucketEmpty:
		value, ok := decoded.(sound.BucketEmpty)
		_, lava := value.Liquid.(block.Lava)
		if !ok || !lava {
			t.Fatalf("bucket empty = %#v", decoded)
		}
	case native.SoundCrossbowLoad:
		value, ok := decoded.(sound.CrossbowLoad)
		if !ok || value.Stage != sound.CrossbowLoadingEnd || !value.QuickCharge {
			t.Fatalf("crossbow load = %#v", decoded)
		}
	case native.SoundGoatHorn:
		value, ok := decoded.(sound.GoatHorn)
		if !ok || value.Horn != sound.Dream() {
			t.Fatalf("goat horn = %#v", decoded)
		}
	}
}

func TestValidSoundRejectsMalformedPayloads(t *testing.T) {
	if !ValidSound(native.WorldSound{Callback: &native.WorldSoundCallback{Function: 1, Context: 2}}) {
		t.Fatal("rejected valid custom sound")
	}
	for _, value := range []native.WorldSound{
		{Kind: native.SoundGoatHorn + 1},
		{Kind: native.SoundBlockPlace},
		{Kind: native.SoundNote, Data: 16},
		{Kind: native.SoundMusicDiscPlay, Data: 21},
		{Kind: native.SoundEquipItem},
		{Kind: native.SoundBucketFill, Data: 2},
		{Kind: native.SoundCrossbowLoad, Integer: 3},
		{Kind: native.SoundGoatHorn, Data: 8},
		{Callback: &native.WorldSoundCallback{Function: 1}},
		{Kind: native.SoundExplosion, Callback: &native.WorldSoundCallback{Function: 1, Context: 2}},
	} {
		if ValidSound(value) {
			t.Fatalf("accepted malformed sound %#v", value)
		}
	}
}
