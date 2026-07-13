BEDROCK_GOPHERS_REV := fe704649d3b741390aa290d0c916cad4cbc1684f
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
PLAYER_COMMAND_LIBRARY := libplayer_command_plugin.dylib
ITEMS_LIBRARY := libitems_command.dylib
SCOREBOARD_LIBRARY := libscoreboard.dylib
FORMS_LIBRARY := libforms.dylib
PING_LIBRARY := libping_command.dylib
WORLD_LIBRARY := libworld_command.dylib
ENTITY_LIBRARY := libentity_command_plugin.dylib
PARTICLE_LIBRARY := libparticle_command_plugin.dylib
SOUND_LIBRARY := libsound_command_plugin.dylib
else
RUNTIME_LIBRARY := libdragonfly_plugin_runtime.so
PLUGIN_LIBRARY := libmovement_guard.so
CHAT_LIBRARY := libchat_filter.so
LIFECYCLE_LIBRARY := liblifecycle_logger.so
COMMAND_LIBRARY := libhello_command.so
PLAYER_COMMAND_LIBRARY := libplayer_command_plugin.so
ITEMS_LIBRARY := libitems_command.so
SCOREBOARD_LIBRARY := libscoreboard.so
FORMS_LIBRARY := libforms.so
PING_LIBRARY := libping_command.so
WORLD_LIBRARY := libworld_command.so
ENTITY_LIBRARY := libentity_command_plugin.so
PARTICLE_LIBRARY := libparticle_command_plugin.so
SOUND_LIBRARY := libsound_command_plugin.so
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
	cargo build --release --manifest-path plugins/player-command/Cargo.toml
	cargo build --release --manifest-path plugins/items-command/Cargo.toml
	cargo build --release --manifest-path plugins/scoreboard/Cargo.toml
	cargo build --release --manifest-path plugins/forms/Cargo.toml
	cargo build --release --manifest-path plugins/ping-command/Cargo.toml
	cargo build --release --manifest-path plugins/world-command/Cargo.toml
	cargo build --release --manifest-path plugins/entity-command/Cargo.toml
	cargo build --release --manifest-path plugins/particle-command/Cargo.toml
	cargo build --release --manifest-path plugins/sound-command/Cargo.toml
	cargo build --release --manifest-path $(CACHE)/Cargo.toml -p dragonfly-plugin-runtime
	mkdir -p lib plugins
	cp $(CACHE)/target/release/$(RUNTIME_LIBRARY) lib/
	cp plugins/movement-guard/target/release/$(PLUGIN_LIBRARY) plugins/
	cp plugins/chat-filter/target/release/$(CHAT_LIBRARY) plugins/
	cp plugins/lifecycle-logger/target/release/$(LIFECYCLE_LIBRARY) plugins/
	cp plugins/hello-command/target/release/$(COMMAND_LIBRARY) plugins/
	cp plugins/player-command/target/release/$(PLAYER_COMMAND_LIBRARY) plugins/
	cp plugins/items-command/target/release/$(ITEMS_LIBRARY) plugins/
	cp plugins/scoreboard/target/release/$(SCOREBOARD_LIBRARY) plugins/
	cp plugins/forms/target/release/$(FORMS_LIBRARY) plugins/
	cp plugins/ping-command/target/release/$(PING_LIBRARY) plugins/
	cp plugins/world-command/target/release/$(WORLD_LIBRARY) plugins/
	cp plugins/entity-command/target/release/$(ENTITY_LIBRARY) plugins/
	cp plugins/particle-command/target/release/$(PARTICLE_LIBRARY) plugins/
	cp plugins/sound-command/target/release/$(SOUND_LIBRARY) plugins/
	go mod download

run: build
	go run .

clean:
	rm -rf .cache .data lib plugins/movement-guard/target plugins/chat-filter/target plugins/lifecycle-logger/target plugins/hello-command/target plugins/player-command/target plugins/items-command/target plugins/ping-command/target plugins/scoreboard/target plugins/forms/target plugins/world-command/target plugins/entity-command/target plugins/particle-command/target plugins/sound-command/target
	rm -f plugins/$(PLUGIN_LIBRARY) plugins/$(CHAT_LIBRARY) plugins/$(LIFECYCLE_LIBRARY) plugins/$(COMMAND_LIBRARY) plugins/$(PLAYER_COMMAND_LIBRARY) plugins/$(ITEMS_LIBRARY) plugins/$(PING_LIBRARY) plugins/$(SCOREBOARD_LIBRARY) plugins/$(FORMS_LIBRARY) plugins/$(WORLD_LIBRARY) plugins/$(ENTITY_LIBRARY) plugins/$(PARTICLE_LIBRARY) plugins/$(SOUND_LIBRARY)
