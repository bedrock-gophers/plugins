package native

import "testing"

type playerTransferRecordingHost struct {
	noopHost
	invocation InvocationID
	player     PlayerID
	world      WorldID
	position   Vec3
	calls      int
}

func (*playerTransferRecordingHost) WorldByName(InvocationID, string) (WorldID, bool) {
	return 72, true
}

func (h *playerTransferRecordingHost) TransferPlayer(invocation InvocationID, player PlayerID, world WorldID, position Vec3) bool {
	h.invocation, h.player, h.world, h.position = invocation, player, world, position
	h.calls++
	return true
}

func TestPlayerTransferCrossesRustAndCHostBoundary(t *testing.T) {
	library, plugins := nativeArtifacts(t)
	host := new(playerTransferRecordingHost)
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
	player := PlayerID{UUID: [16]byte{1, 3, 5, 7}, Generation: 19}
	output, err := runtime.HandleCommand(commandNamed(t, commands, "world").Index, CommandInput{
		Invocation: 61, Source: "Traveller", SourceKind: CommandSourcePlayer, SourcePlayer: &player,
		OnlinePlayers: []CommandPlayer{{Player: player, Name: "Traveller"}},
		Arguments:     "transfer example:target",
	})
	if err != nil || output.Failed {
		t.Fatalf("output=%+v error=%v", output, err)
	}
	if host.calls != 1 || host.invocation != 61 || host.player != player || host.world != 72 || host.position != (Vec3{X: 0.5, Y: 65, Z: 0.5}) {
		t.Fatalf("transfer calls=%d invocation=%d player=%#v world=%d position=%#v", host.calls, host.invocation, host.player, host.world, host.position)
	}
}
