BEDROCK_GOPHERS_REV := 54dd9c54eb4fd273eedd3f0253061986afe0a920
BEDROCK_GOPHERS_SHORT_REV := $(shell printf '%.12s' $(BEDROCK_GOPHERS_REV))
GO_FRAMEWORK_REV := $(shell go list -m -f '{{.Version}}' github.com/bedrock-gophers/plugins | sed 's/.*-//')
CACHE := .cache/bedrock-gophers
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
RUNTIME_LIBRARY := libdragonfly_plugin_runtime.dylib
PLUGIN_LIBRARY := libmovement_guard.dylib
CHAT_LIBRARY := libchat_filter.dylib
LIFECYCLE_LIBRARY := liblifecycle_logger.dylib
COMMAND_LIBRARY := libhello_command.dylib
PING_LIBRARY := libping_command.dylib
else
RUNTIME_LIBRARY := libdragonfly_plugin_runtime.so
PLUGIN_LIBRARY := libmovement_guard.so
CHAT_LIBRARY := libchat_filter.so
LIFECYCLE_LIBRARY := liblifecycle_logger.so
COMMAND_LIBRARY := libhello_command.so
PING_LIBRARY := libping_command.so
endif

.PHONY: check-revision prepare build run clean

check-revision:
	@test "$(GO_FRAMEWORK_REV)" = "$(BEDROCK_GOPHERS_SHORT_REV)" || (printf 'Go framework is pinned to %s, expected %s\n' "$(GO_FRAMEWORK_REV)" "$(BEDROCK_GOPHERS_SHORT_REV)"; exit 1)

$(CACHE)/.git:
	mkdir -p .cache
	git clone --quiet https://github.com/bedrock-gophers/plugins $(CACHE)

prepare: $(CACHE)/.git
	git -C $(CACHE) fetch --quiet origin
	git -C $(CACHE) checkout --quiet $(BEDROCK_GOPHERS_REV)

build: check-revision prepare
	cargo build --release --manifest-path plugins/movement-guard/Cargo.toml
	cargo build --release --manifest-path plugins/chat-filter/Cargo.toml
	cargo build --release --manifest-path plugins/lifecycle-logger/Cargo.toml
	cargo build --release --manifest-path plugins/hello-command/Cargo.toml
	cargo build --release --manifest-path plugins/ping-command/Cargo.toml
	cargo build --release --manifest-path $(CACHE)/Cargo.toml -p dragonfly-plugin-runtime
	mkdir -p lib plugins
	cp $(CACHE)/target/release/$(RUNTIME_LIBRARY) lib/
	cp plugins/movement-guard/target/release/$(PLUGIN_LIBRARY) plugins/
	cp plugins/chat-filter/target/release/$(CHAT_LIBRARY) plugins/
	cp plugins/lifecycle-logger/target/release/$(LIFECYCLE_LIBRARY) plugins/
	cp plugins/hello-command/target/release/$(COMMAND_LIBRARY) plugins/
	cp plugins/ping-command/target/release/$(PING_LIBRARY) plugins/
	go mod download

run: build
	go run .

clean:
	rm -rf .cache .data lib plugins/movement-guard/target plugins/chat-filter/target plugins/lifecycle-logger/target plugins/hello-command/target plugins/ping-command/target
	rm -f plugins/$(PLUGIN_LIBRARY) plugins/$(CHAT_LIBRARY) plugins/$(LIFECYCLE_LIBRARY) plugins/$(COMMAND_LIBRARY) plugins/$(PING_LIBRARY)
