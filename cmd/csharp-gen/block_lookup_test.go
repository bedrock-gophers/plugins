package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspectWorldBlockByName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "block.go")
	if err := os.WriteFile(path, []byte(`package world
type Block interface{}
func BlockByName(name string, properties map[string]any) (Block, bool) { return nil, false }
`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := inspectWorldBlockByName(path); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte(`package world
type Block interface{}
func BlockByName(identifier string, properties map[string]string) (Block, bool) { return nil, false }
`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := inspectWorldBlockByName(path); err == nil {
		t.Fatal("inspectWorldBlockByName accepted a changed signature")
	}
}
