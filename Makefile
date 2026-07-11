BEDROCK_GOPHERS_REV := 6e524b99f993edc278cc0f336ce4e53ee60d279a
CACHE := .cache/bedrock-gophers
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
RUNTIME_LIBRARY := libdragonfly_plugin_runtime.dylib
PLUGIN_LIBRARY := libmovement_guard.dylib
CHAT_LIBRARY := libchat_filter.dylib
LIFECYCLE_LIBRARY := liblifecycle_logger.dylib
else
RUNTIME_LIBRARY := libdragonfly_plugin_runtime.so
PLUGIN_LIBRARY := libmovement_guard.so
CHAT_LIBRARY := libchat_filter.so
LIFECYCLE_LIBRARY := liblifecycle_logger.so
endif

.PHONY: build run clean

$(CACHE)/Cargo.toml:
	mkdir -p .cache
	git clone --quiet https://github.com/bedrock-gophers/plugins $(CACHE)
	git -C $(CACHE) checkout --quiet $(BEDROCK_GOPHERS_REV)

build: $(CACHE)/Cargo.toml
	cargo build --release --manifest-path plugins/movement-guard/Cargo.toml
	cargo build --release --manifest-path plugins/chat-filter/Cargo.toml
	cargo build --release --manifest-path plugins/lifecycle-logger/Cargo.toml
	cargo build --release --manifest-path $(CACHE)/Cargo.toml -p dragonfly-plugin-runtime
	mkdir -p lib plugins
	cp $(CACHE)/target/release/$(RUNTIME_LIBRARY) lib/
	cp plugins/movement-guard/target/release/$(PLUGIN_LIBRARY) plugins/
	cp plugins/chat-filter/target/release/$(CHAT_LIBRARY) plugins/
	cp plugins/lifecycle-logger/target/release/$(LIFECYCLE_LIBRARY) plugins/
	go mod download

run: build
	go run .

clean:
	rm -rf .cache .data lib plugins/movement-guard/target plugins/chat-filter/target plugins/lifecycle-logger/target
	rm -f plugins/$(PLUGIN_LIBRARY) plugins/$(CHAT_LIBRARY) plugins/$(LIFECYCLE_LIBRARY)
