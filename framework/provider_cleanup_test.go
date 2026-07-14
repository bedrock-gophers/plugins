package framework

import (
	"io"
	"log/slog"
	"testing"
)

type cleanupTestProvider struct {
	closed bool
}

func (provider *cleanupTestProvider) Close() error {
	provider.closed = true
	return nil
}

func TestCloseDragonflyProviderIgnoresTypedNil(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	var missing *cleanupTestProvider
	closeDragonflyProvider("test", missing, log)

	present := new(cleanupTestProvider)
	closeDragonflyProvider("test", present, log)
	if !present.closed {
		t.Fatal("non-nil provider was not closed")
	}
}
