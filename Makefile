BEDROCK_GOPHERS_REV := 68f5e89fc5e5bcea544f9527ef39fa8c975d0a87
CACHE := .cache/bedrock-gophers
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
RUNTIME_LIBRARY := libdragonfly_plugin_runtime.dylib
PLUGIN_LIBRARY := libmovement_guard.dylib
CHAT_LIBRARY := libchat_filter.dylib
LIFECYCLE_LIBRARY := liblifecycle_logger.dylib
COMMAND_LIBRARY := libhello_command.dylib
else
RUNTIME_LIBRARY := libdragonfly_plugin_runtime.so
PLUGIN_LIBRARY := libmovement_guard.so
CHAT_LIBRARY := libchat_filter.so
LIFECYCLE_LIBRARY := liblifecycle_logger.so
COMMAND_LIBRARY := libhello_command.so
endif

.PHONY: prepare build run clean

$(CACHE)/.git:
	mkdir -p .cache
	git clone --quiet https://github.com/bedrock-gophers/plugins $(CACHE)

prepare: $(CACHE)/.git
	git -C $(CACHE) fetch --quiet origin
	git -C $(CACHE) checkout --quiet $(BEDROCK_GOPHERS_REV)

build: prepare
	cargo build --release --manifest-path plugins/movement-guard/Cargo.toml
	cargo build --release --manifest-path plugins/chat-filter/Cargo.toml
	cargo build --release --manifest-path plugins/lifecycle-logger/Cargo.toml
	cargo build --release --manifest-path plugins/hello-command/Cargo.toml
	cargo build --release --manifest-path $(CACHE)/Cargo.toml -p dragonfly-plugin-runtime
	mkdir -p lib plugins
	cp $(CACHE)/target/release/$(RUNTIME_LIBRARY) lib/
	cp plugins/movement-guard/target/release/$(PLUGIN_LIBRARY) plugins/
	cp plugins/chat-filter/target/release/$(CHAT_LIBRARY) plugins/
	cp plugins/lifecycle-logger/target/release/$(LIFECYCLE_LIBRARY) plugins/
	cp plugins/hello-command/target/release/$(COMMAND_LIBRARY) plugins/
	go mod download

run: build
	go run .

clean:
	rm -rf .cache .data lib plugins/movement-guard/target plugins/chat-filter/target plugins/lifecycle-logger/target plugins/hello-command/target
	rm -f plugins/$(PLUGIN_LIBRARY) plugins/$(CHAT_LIBRARY) plugins/$(LIFECYCLE_LIBRARY) plugins/$(COMMAND_LIBRARY)
