package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestItemSchemaNamesAreSafeRustIdentifiers(t *testing.T) {
	valid := []string{"beef", "cooked_beef", "tier2"}
	for _, value := range valid {
		if !validSchemaName(value) {
			t.Fatalf("expected %q to be valid", value)
		}
	}
	invalid := []string{"", "Beef", "2tier", "cooked_", "cooked__beef", "cooked-beef"}
	for _, value := range invalid {
		if validSchemaName(value) {
			t.Fatalf("expected %q to be invalid", value)
		}
	}
	if validRustValueName("type") {
		t.Fatal("Rust keyword passed value-name validation")
	}
}

func TestReadItemsRejectsGeneratedTypeCollisions(t *testing.T) {
	schema := itemTestSchema(t)
	schema = strings.Replace(schema, "name: apple", "name: helmet", 1)
	path := filepath.Join(t.TempDir(), "items.yaml")
	if err := os.WriteFile(path, []byte(schema), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readItems(path); err == nil {
		t.Fatal("expected generated Helmet collision to fail")
	}
}

func TestReadItemsRejectsReservedGeneratedTypes(t *testing.T) {
	schema := itemTestSchema(t)
	schema = strings.Replace(schema, "name: apple", "name: tool_tier", 1)
	path := filepath.Join(t.TempDir(), "items.yaml")
	if err := os.WriteFile(path, []byte(schema), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readItems(path); err == nil {
		t.Fatal("expected generated ToolTier collision to fail")
	}
}

func TestReadItemsRejectsRustKeywordMethods(t *testing.T) {
	schema := itemTestSchema(t)
	schema = strings.Replace(schema, "property: cooked", "property: type", 1)
	path := filepath.Join(t.TempDir(), "items.yaml")
	if err := os.WriteFile(path, []byte(schema), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readItems(path); err == nil {
		t.Fatal("expected Rust keyword method to fail")
	}
}

func itemTestSchema(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "schema", "items.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
