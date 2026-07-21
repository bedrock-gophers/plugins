package framework

import "testing"

func TestDecodeCustomBlockDataUsesNativeProperties(t *testing.T) {
	data, err := decodeCustomBlockData(`{
		"server":{"hardness":3,"render_method":1},
		"states":{"test:active":[false,true]},
		"permutations":[{"condition":"q.block_state('test:active')","properties":{"map_color":"#ff0000"}}]
	}`)
	if err != nil {
		t.Fatal(err)
	}
	if got := floatValue(data.server["hardness"]); got != 3 {
		t.Fatalf("hardness = %v, want 3", got)
	}
	if len(data.states["test:active"]) != 2 {
		t.Fatalf("state values = %v, want 2 values", data.states["test:active"])
	}
	if len(data.permutations) != 1 || stringValue(data.permutations[0].properties["map_color"]) != "#ff0000" {
		t.Fatalf("unexpected permutations: %#v", data.permutations)
	}
}

func TestDecodeCustomBlockDataRejectsArbitraryComponents(t *testing.T) {
	if _, err := decodeCustomBlockData(`{"components":{"minecraft:random_offset":{}}}`); err == nil {
		t.Fatal("expected arbitrary block components to be rejected")
	}
}

func TestCustomBlockStatesCartesianProduct(t *testing.T) {
	states, err := customBlockStates(map[string][]any{
		"test:active": {false, true},
		"test:level":  {0, 1, 2},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(states) != 6 {
		t.Fatalf("state count = %d, want 6", len(states))
	}
}

func TestCustomBlockWithoutGeometryReturnsNil(t *testing.T) {
	if geometry := customBlockGeometry(nil); geometry != nil {
		t.Fatalf("geometry = %v, want nil", geometry)
	}
}
