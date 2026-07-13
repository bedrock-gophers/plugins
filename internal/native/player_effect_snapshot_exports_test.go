package native

import (
	"math"
	"testing"
	"time"
	"unsafe"
)

func TestPlayerEffectSnapshotViewValidation(t *testing.T) {
	valid := PlayerEffect{Type: EffectSpeed, Level: 2, Duration: time.Second, Potency: 1, Mode: PlayerEffectTimed}
	ambient := valid
	ambient.Mode, ambient.ParticlesHidden = PlayerEffectAmbient, true
	infinite := valid
	infinite.Mode, infinite.Duration = PlayerEffectInfinite, 0
	for name, value := range map[string]PlayerEffect{"timed": valid, "ambient": ambient, "infinite": infinite} {
		t.Run(name, func(t *testing.T) {
			view, ok := playerEffectSnapshotView(value)
			if !ok || int32(view.effect_type) != int32(value.Type) || int32(view.level) != value.Level || uint32(view.mode) != uint32(value.Mode) || (view.particles_hidden != 0) != value.ParticlesHidden {
				t.Fatalf("view = %#v ok=%v", view, ok)
			}
		})
	}

	invalid := map[string]PlayerEffect{
		"zero level":        func() PlayerEffect { value := valid; value.Level = 0; return value }(),
		"negative duration": func() PlayerEffect { value := valid; value.Duration = -1; return value }(),
		"zero potency":      func() PlayerEffect { value := valid; value.Potency = 0; return value }(),
		"nan potency":       func() PlayerEffect { value := valid; value.Potency = math.NaN(); return value }(),
		"instant":           func() PlayerEffect { value := valid; value.Mode = PlayerEffectInstant; return value }(),
		"unknown mode":      func() PlayerEffect { value := valid; value.Mode = 99; return value }(),
		"infinite duration": func() PlayerEffect { value := infinite; value.Duration = time.Second; return value }(),
	}
	for name, value := range invalid {
		t.Run(name, func(t *testing.T) {
			if _, ok := playerEffectSnapshotView(value); ok {
				t.Fatalf("invalid effect accepted: %#v", value)
			}
		})
	}
}

func TestPlayerEffectSnapshotViewsEnforceBound(t *testing.T) {
	value := PlayerEffect{Type: EffectSpeed, Level: 1, Duration: time.Second, Potency: 1, Mode: PlayerEffectTimed}
	values := make([]PlayerEffect, maxPlayerEffects)
	for index := range values {
		values[index] = value
	}
	if encoded, ok := playerEffectSnapshotViews(values); !ok || len(encoded) != maxPlayerEffects {
		t.Fatalf("maximum snapshot length=%d ok=%v", len(encoded), ok)
	}
	values = append(values, value)
	if encoded, ok := playerEffectSnapshotViews(values); ok || encoded != nil {
		t.Fatalf("oversized snapshot length=%d ok=%v", len(encoded), ok)
	}
}

func TestPlayerEffectBufferRequiresPointerForNonZeroCapacity(t *testing.T) {
	if validPlayerEffectBuffer(nil, 1) {
		t.Fatal("accepted null data with non-zero capacity")
	}
	if !validPlayerEffectBuffer(nil, 0) {
		t.Fatal("rejected zero-capacity sizing buffer")
	}
	value := byte(0)
	if !validPlayerEffectBuffer(unsafe.Pointer(&value), 1) {
		t.Fatal("rejected non-null data buffer")
	}
}

func TestWriteBoundedSnapshotDoesNotPartiallyWrite(t *testing.T) {
	sentinel := []int{41}
	writes := 0
	required, ok := writeBoundedSnapshot([]int{1, 2}, 1, func(values []int) {
		writes++
		copy(sentinel, values)
	})
	if ok || required != 2 || writes != 0 || sentinel[0] != 41 {
		t.Fatalf("required=%d ok=%v writes=%d sentinel=%v", required, ok, writes, sentinel)
	}
}
