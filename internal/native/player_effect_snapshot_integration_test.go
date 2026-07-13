package native

import (
	"strings"
	"testing"
	"time"
)

type playerEffectSnapshotRecordingHost struct {
	noopHost
	effectsInvocation InvocationID
	clearInvocation   InvocationID
	player            PlayerID
	effectCalls       int
	clearCalls        int
	messages          []string
}

func (h *playerEffectSnapshotRecordingHost) PlayerEffects(invocation InvocationID, player PlayerID) ([]PlayerEffect, bool) {
	h.effectsInvocation, h.player = invocation, player
	h.effectCalls++
	return []PlayerEffect{
		{Type: -7, Level: 2, Duration: 1500 * time.Millisecond, Potency: 1, Mode: PlayerEffectAmbient, ParticlesHidden: true},
		{Type: EffectDarkness, Level: 1, Potency: 1, Mode: PlayerEffectInfinite},
	}, true
}

func (h *playerEffectSnapshotRecordingHost) ClearPlayerEffects(invocation InvocationID, player PlayerID) bool {
	h.clearInvocation, h.player = invocation, player
	h.clearCalls++
	return true
}

func (h *playerEffectSnapshotRecordingHost) SendPlayerText(_ InvocationID, _ PlayerID, _ PlayerTextKind, message string) bool {
	h.messages = append(h.messages, message)
	return true
}

func TestPlayerEffectsCrossRustAndCHostBoundary(t *testing.T) {
	library, plugins := nativeArtifacts(t)
	host := new(playerEffectSnapshotRecordingHost)
	runtime, err := OpenWithHost(library, plugins, host)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(runtime.Close)
	if err := runtime.Enable(); err != nil {
		t.Fatal(err)
	}
	commands, err := runtime.Commands()
	if err != nil {
		t.Fatal(err)
	}
	player := PlayerID{UUID: [16]byte{2, 4, 6, 8}, Generation: 23}
	input := CommandInput{
		Invocation: 67, Source: "Affected", SourceKind: CommandSourcePlayer, SourcePlayer: &player,
		OnlinePlayers: []CommandPlayer{{Player: player, Name: "Affected"}},
		Arguments:     "effects",
	}
	output, err := runtime.HandleCommand(commandNamed(t, commands, "player").Index, input)
	if err != nil || output.Failed {
		t.Fatalf("effects output=%+v error=%v", output, err)
	}
	input.Invocation, input.Arguments = 71, "clear-effects"
	output, err = runtime.HandleCommand(commandNamed(t, commands, "player").Index, input)
	if err != nil || output.Failed {
		t.Fatalf("clear output=%+v error=%v", output, err)
	}
	if host.effectCalls != 2 || host.clearCalls != 1 || host.effectsInvocation != 67 || host.clearInvocation != 71 || host.player != player {
		t.Fatalf("effect calls=%d clear calls=%d effects invocation=%d clear invocation=%d player=%#v", host.effectCalls, host.clearCalls, host.effectsInvocation, host.clearInvocation, host.player)
	}
	joined := strings.Join(host.messages, "\n")
	for _, fragment := range []string{
		"Effect -7 level=2 duration_ms=1500 ambient=true infinite=false particles_hidden=true",
		"Effect 30 level=1 duration_ms=0 ambient=false infinite=true particles_hidden=false",
		"Effects cleared.",
	} {
		if !strings.Contains(joined, fragment) {
			t.Fatalf("messages %q do not contain %q", joined, fragment)
		}
	}
}
