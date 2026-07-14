package native

import (
	"math"
	"testing"
	"time"
	"unsafe"
)

func TestPlayerEffectSnapshotViewValidation(t *testing.T) {
	valid := PlayerEffect{Type: EffectSpeed, Level: 2, Duration: time.Second + time.Nanosecond, Potency: 0.75, Tick: 7}
	ambient := valid
	ambient.Ambient, ambient.ParticlesHidden = true, true
	infinite := valid
	infinite.Infinite, infinite.Duration = true, 0
	negative := valid
	negative.Duration = -time.Nanosecond
	nanPotency := valid
	nanPotency.Potency = math.NaN()
	odd := valid
	odd.Level, odd.Ambient, odd.Infinite, odd.Tick = 0, true, true, -1
	for name, value := range map[string]PlayerEffect{"timed": valid, "ambient": ambient, "infinite": infinite, "negative duration": negative, "nan potency": nanPotency, "raw fields": odd} {
		t.Run(name, func(t *testing.T) {
			view, ok := playerEffectSnapshotView(value)
			if !ok || int32(view.effect_type) != int32(value.Type) || int32(view.level) != value.Level ||
				time.Duration(view.duration_nanoseconds) != value.Duration || (view.ambient != 0) != value.Ambient ||
				(view.particles_hidden != 0) != value.ParticlesHidden || (view.infinite != 0) != value.Infinite ||
				int64(view.tick) != value.Tick || !(math.IsNaN(value.Potency) && math.IsNaN(float64(view.potency)) || float64(view.potency) == value.Potency) {
				t.Fatalf("view = %#v ok=%v", view, ok)
			}
		})
	}

}

func TestPlayerEffectSnapshotViewsHaveNoInventedCountLimit(t *testing.T) {
	value := PlayerEffect{Type: EffectSpeed, Level: 1, Duration: time.Second, Potency: 1}
	values := make([]PlayerEffect, 257)
	for index := range values {
		values[index] = value
	}
	if encoded, ok := playerEffectSnapshotViews(values); !ok || len(encoded) != len(values) {
		t.Fatalf("snapshot length=%d ok=%v", len(encoded), ok)
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
