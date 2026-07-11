BEDROCK_GOPHERS_REV := e5390b1f06507a3b36d9f947beb183e4c6001f6a
CACHE := .cache/bedrock-gophers
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
RUNTIME_LIBRARY := libdragonfly_plugin_runtime.dylib
PLUGIN_LIBRARY := libmovement_guard.dylib
else
RUNTIME_LIBRARY := libdragonfly_plugin_runtime.so
PLUGIN_LIBRARY := libmovement_guard.so
endif

.PHONY: build run clean

$(CACHE)/Cargo.toml:
	mkdir -p .cache
	git clone --quiet https://github.com/bedrock-gophers/plugins $(CACHE)
	git -C $(CACHE) checkout --quiet $(BEDROCK_GOPHERS_REV)

build: $(CACHE)/Cargo.toml
	cargo build --release --manifest-path plugins/movement-guard/Cargo.toml
	cargo build --release --manifest-path $(CACHE)/Cargo.toml -p dragonfly-plugin-runtime
	mkdir -p lib plugins
	cp $(CACHE)/target/release/$(RUNTIME_LIBRARY) lib/
	cp plugins/movement-guard/target/release/$(PLUGIN_LIBRARY) plugins/
	go mod download

run: build
	go run .

clean:
	rm -rf .cache .data lib plugins/movement-guard/target
	rm -f plugins/$(PLUGIN_LIBRARY)
