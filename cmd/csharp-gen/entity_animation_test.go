package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEntityAnimationMatchesPinnedDragonfly(t *testing.T) {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/df-mc/dragonfly")
	directory, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	spec, err := inspectEntityAnimation(filepath.Join(string(bytes.TrimSpace(directory)), "server", "world", "entity_animation.go"))
	if err != nil {
		t.Fatal(err)
	}
	generated := string(generateEntityAnimation(spec))
	for _, expected := range []string{
		"public static EntityAnimation NewEntityAnimation(string name)",
		"public readonly struct EntityAnimation",
		"public string Name()",
		"public string Controller()",
		"public EntityAnimation WithController(string controller)",
		"public string NextState()",
		"public EntityAnimation WithNextState(string state)",
		"public string StopCondition()",
		"public EntityAnimation WithStopCondition(string condition)",
	} {
		if !strings.Contains(generated, expected) {
			t.Fatalf("generated entity animation missing %q:\n%s", expected, generated)
		}
	}
}

func TestEntityAnimationRejectsASTDrift(t *testing.T) {
	valid := `package world
type EntityAnimation struct { name string; nextState string; controller string; stopCondition string }
func NewEntityAnimation(name string) EntityAnimation { return EntityAnimation{name: name} }
func (a EntityAnimation) Name() string { return a.name }
func (a EntityAnimation) Controller() string { return a.controller }
func (a EntityAnimation) WithController(controller string) EntityAnimation { a.controller = controller; return a }
func (a EntityAnimation) NextState() string { return a.nextState }
func (a EntityAnimation) WithNextState(state string) EntityAnimation { a.nextState = state; return a }
func (a EntityAnimation) StopCondition() string { return a.stopCondition }
func (a EntityAnimation) WithStopCondition(condition string) EntityAnimation { a.stopCondition = condition; return a }`
	for name, source := range map[string]string{
		"field":       strings.Replace(valid, "nextState string", "next string", 1),
		"constructor": strings.Replace(valid, "EntityAnimation{name: name}", "EntityAnimation{controller: name}", 1),
		"getter":      strings.Replace(valid, "return a.controller", "return a.name", 1),
		"with method": strings.Replace(valid, "a.nextState = state", "a.controller = state", 1),
		"receiver":    strings.Replace(valid, "(a EntityAnimation) Name", "(a *EntityAnimation) Name", 1),
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "entity_animation.go")
			if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := inspectEntityAnimation(path); err == nil {
				t.Fatal("expected entity animation AST drift error")
			}
		})
	}
}
