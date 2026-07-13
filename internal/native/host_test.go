package native

import "testing"

func TestDetachedEntityIDValidRequiresBothFields(t *testing.T) {
	for _, id := range []DetachedEntityID{
		{},
		{Value: 1},
		{Generation: 1},
	} {
		if id.Valid() {
			t.Fatalf("partial-zero ID is valid: %#v", id)
		}
	}
	if id := (DetachedEntityID{Value: 1, Generation: 1}); !id.Valid() {
		t.Fatalf("non-zero ID is invalid: %#v", id)
	}
}
