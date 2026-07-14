package native

import (
	"math"
	"testing"
	"time"
)

func TestCopyPlayerEffectPreservesNanoseconds(t *testing.T) {
	view, ok := playerEffectSnapshotView(PlayerEffect{
		Type: EffectSpeed, Level: 2, Duration: time.Second + time.Nanosecond,
		Potency: 1, Tick: 3,
	})
	if !ok {
		t.Fatal("snapshot rejected valid effect")
	}
	effect, ok := copyPlayerEffect(PlayerEffectAdd, view)
	if !ok || effect.Duration != time.Second+time.Nanosecond || effect.Tick != 3 {
		t.Fatalf("effect = %+v, ok=%v", effect, ok)
	}
}

func TestCopyPlayerEffectPreservesNegativeNanoseconds(t *testing.T) {
	view, ok := playerEffectSnapshotView(PlayerEffect{
		Type: EffectSpeed, Level: 1, Duration: time.Nanosecond,
		Potency: 1,
	})
	if !ok {
		t.Fatal("snapshot rejected valid effect")
	}
	view.duration_nanoseconds = -1
	if effect, ok := copyPlayerEffect(PlayerEffectAdd, view); !ok || effect.Duration != -time.Nanosecond {
		t.Fatalf("effect = %+v, ok=%v", effect, ok)
	}
}

func TestCopyPlayerEffectPreservesMaximumDuration(t *testing.T) {
	view, ok := playerEffectSnapshotView(PlayerEffect{
		Type: EffectSpeed, Level: 1, Duration: time.Duration(math.MaxInt64),
		Potency: 1,
	})
	if !ok {
		t.Fatal("snapshot rejected maximum duration")
	}
	effect, ok := copyPlayerEffect(PlayerEffectAdd, view)
	if !ok || effect.Duration != time.Duration(math.MaxInt64) {
		t.Fatalf("effect duration = %v, ok=%v", effect.Duration, ok)
	}
}

func TestCopyPlayerEffectRejectsNonBooleanFlags(t *testing.T) {
	view, ok := playerEffectSnapshotView(PlayerEffect{Type: EffectSpeed, Level: 1, Potency: 1})
	if !ok {
		t.Fatal("snapshot rejected valid effect")
	}
	view.ambient = 2
	if effect, ok := copyPlayerEffect(PlayerEffectAdd, view); ok {
		t.Fatalf("non-boolean ambient accepted: %+v", effect)
	}
}
