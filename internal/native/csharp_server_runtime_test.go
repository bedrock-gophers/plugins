package native

import (
	"slices"
	"testing"
)

type csharpServerHost struct {
	*recordingHost
	iterator               PlayerIteratorID
	snapshots              []PlayerSnapshot
	playerInvocations      []InvocationID
	index                  int
	openInvocations        []InvocationID
	nextInvocations        []InvocationID
	closeInvocations       []InvocationID
	closedIterators        []PlayerIteratorID
	lookupUUID             [16]byte
	lookupName             string
	lookupHandle           EntityHandleID
	entityHandleUUIDs      map[EntityHandleID][16]byte
	entityHandleCalls      []EntityID
	entityHandleInvocation []InvocationID
	textInvocations        []InvocationID
}

func (h *csharpServerHost) OpenServerPlayerIterator(invocation InvocationID) (PlayerIteratorID, bool) {
	h.openInvocations = append(h.openInvocations, invocation)
	h.index = 0
	return h.iterator, h.iterator != 0
}

func (h *csharpServerHost) NextServerPlayer(invocation InvocationID, iterator PlayerIteratorID) (InvocationID, PlayerSnapshot, bool, bool) {
	h.nextInvocations = append(h.nextInvocations, invocation)
	if iterator != h.iterator || h.index > len(h.snapshots) {
		return 0, PlayerSnapshot{}, false, false
	}
	if h.index == len(h.snapshots) {
		h.index++
		return 0, PlayerSnapshot{}, false, true
	}
	nested, snapshot := h.playerInvocations[h.index], h.snapshots[h.index]
	h.index++
	return nested, snapshot, true, true
}

func (h *csharpServerHost) CloseServerPlayers(invocation InvocationID, iterator PlayerIteratorID) {
	h.closeInvocations = append(h.closeInvocations, invocation)
	h.closedIterators = append(h.closedIterators, iterator)
}

func (h *csharpServerHost) ServerPlayer(uuid [16]byte) (EntityHandleID, bool, bool) {
	h.lookupUUID = uuid
	return h.lookupHandle, true, h.lookupHandle.Valid()
}

func (h *csharpServerHost) ServerPlayerByName(name string) (EntityHandleID, bool, bool) {
	h.lookupName = name
	return h.lookupHandle, true, h.lookupHandle.Valid()
}

func (h *csharpServerHost) EntityHandle(invocation InvocationID, entity EntityID) (EntityHandleID, bool) {
	h.entityHandleInvocation = append(h.entityHandleInvocation, invocation)
	h.entityHandleCalls = append(h.entityHandleCalls, entity)
	handle := EntityHandleID{Value: uint64(entity.UUID[0]) + 1, Generation: entity.Generation}
	if h.entityHandleUUIDs == nil {
		h.entityHandleUUIDs = map[EntityHandleID][16]byte{}
	}
	h.entityHandleUUIDs[handle] = entity.UUID
	return handle, handle.Valid()
}

func (h *csharpServerHost) EntityHandleUUID(handle EntityHandleID) ([16]byte, bool) {
	uuid, ok := h.entityHandleUUIDs[handle]
	return uuid, ok
}

func (h *csharpServerHost) SendPlayerText(invocation InvocationID, player PlayerID, kind PlayerTextKind, message string) bool {
	h.textInvocations = append(h.textInvocations, invocation)
	return h.recordingHost.SendPlayerText(invocation, player, kind, message)
}

func TestCSharpServerPlayersAndLookup(t *testing.T) {
	source := PlayerID{UUID: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, Generation: 7}
	first := PlayerID{UUID: [16]byte{2}, Generation: 8}
	second := PlayerID{UUID: [16]byte{3}, Generation: 9}
	host := &csharpServerHost{
		recordingHost:     &recordingHost{},
		iterator:          17,
		playerInvocations: []InvocationID{101, 102},
		snapshots: []PlayerSnapshot{
			{Player: first, Name: "Alpha", LatencyMilliseconds: 12, Position: Vec3{X: 1, Y: 64, Z: 2}},
			{Player: second, Name: "Bravo", LatencyMilliseconds: 34, Position: Vec3{X: 3, Y: 65, Z: 4}},
		},
		lookupHandle: EntityHandleID{Value: 2, Generation: source.Generation},
	}
	host.entityHandleUUIDs = map[EntityHandleID][16]byte{host.lookupHandle: source.UUID}

	pluginRuntime := openCSharpRuntimeWithHost(t, host)
	commands, err := pluginRuntime.Commands()
	if err != nil {
		t.Fatal(err)
	}
	kitchen := commandNamed(t, commands, "kitchen")
	var overload uint64
	found := false
	for index, candidate := range kitchen.Overloads {
		if len(candidate.Parameters) == 1 && candidate.Parameters[0].Name == "server" {
			overload, found = uint64(index), true
			break
		}
	}
	if !found {
		t.Fatalf("server overload missing: %#v", kitchen.Overloads)
	}
	if len(host.openInvocations) != 0 {
		t.Fatalf("server iteration opened before command execution: %v", host.openInvocations)
	}

	output, err := pluginRuntime.HandleCommand(kitchen.Index, CommandInput{
		Invocation: 42, Source: "Danick", SourceKind: CommandSourcePlayer, SourcePlayer: &source,
		Overload: overload, Arguments: []string{"server"},
		OnlinePlayers: []CommandPlayer{{Player: source, Name: "Danick"}},
	})
	if err != nil || output.Failed || output.Message != "players=2, first=02000000-0000-0000-0000-000000000000" {
		t.Fatalf("server output=%#v error=%v", output, err)
	}
	if !slices.Equal(host.openInvocations, []InvocationID{42}) ||
		!slices.Equal(host.nextInvocations, []InvocationID{42, 42, 42}) ||
		!slices.Equal(host.closeInvocations, []InvocationID{42}) ||
		!slices.Equal(host.closedIterators, []PlayerIteratorID{17}) {
		t.Fatalf("server iterator calls: open=%v next=%v close=%v iterators=%v",
			host.openInvocations, host.nextInvocations, host.closeInvocations, host.closedIterators)
	}
	if !slices.Equal(host.textInvocations, []InvocationID{101, 102}) ||
		!slices.Equal(host.textPlayers, []PlayerID{first, second}) ||
		!slices.Equal(host.texts, []string{"Kitchen server iteration is live.", "Kitchen server iteration is live."}) {
		t.Fatalf("iterated player calls: invocations=%v players=%v texts=%v",
			host.textInvocations, host.textPlayers, host.texts)
	}
	if host.lookupUUID != source.UUID || host.lookupName != "Danick" ||
		!slices.Equal(host.entityHandleInvocation, []InvocationID{101, 42}) ||
		!slices.Equal(host.entityHandleCalls, []EntityID{
			{UUID: first.UUID, Generation: first.Generation},
			{UUID: source.UUID, Generation: source.Generation},
		}) {
		t.Fatalf("server lookup calls: uuid=%x name=%q handle invocations=%v entities=%v",
			host.lookupUUID, host.lookupName, host.entityHandleInvocation, host.entityHandleCalls)
	}
}
