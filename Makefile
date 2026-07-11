.PHONY: generate check-generated build-native build-server build test benchmark clean

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
RUNTIME_LIBRARY := libdragonfly_plugin_runtime.dylib
PLUGIN_LIBRARY := libmovement_guard_plugin.dylib
CHAT_PLUGIN_LIBRARY := libchat_filter_plugin.dylib
else
RUNTIME_LIBRARY := libdragonfly_plugin_runtime.so
PLUGIN_LIBRARY := libmovement_guard_plugin.so
CHAT_PLUGIN_LIBRARY := libchat_filter_plugin.so
endif

generate:
	go run ./tools/abi-gen -root .
	cargo fmt --all

check-generated:
	go run ./tools/abi-gen -root . -check
	cargo fmt --all -- --check

build-native: generate
	cargo build --release -p dragonfly-plugin-runtime -p movement-guard-plugin -p chat-filter-plugin
	mkdir -p build/lib build/plugins
	cp target/release/$(RUNTIME_LIBRARY) build/lib/
	cp target/release/$(PLUGIN_LIBRARY) build/plugins/
	cp target/release/$(CHAT_PLUGIN_LIBRARY) build/plugins/

build-server:
	mkdir -p build
	go build -o build/bedrock-gophers ./cmd/bedrock-gophers

build: build-native build-server

test: build-native check-generated
	cargo test --workspace
	go test ./...

benchmark: build-native
	go test ./internal/native -run '^$$' -bench . -benchmem

clean:
	cargo clean
	rm -rf build
