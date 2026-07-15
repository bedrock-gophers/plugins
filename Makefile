BEDROCK_GOPHERS_REV := 8d621595973a99afaebfeaffa360aac3b5890461
BEDROCK_GOPHERS_SHORT_REV := $(shell printf '%.12s' $(BEDROCK_GOPHERS_REV))
GO_FRAMEWORK_REV := $(shell go list -m -f '{{.Version}}' github.com/bedrock-gophers/plugins | sed 's/.*-//')
CACHE := .cache/bedrock-gophers
BUILD := .build
PLUGIN_PROJECTS := $(sort $(wildcard plugins/*/*.csproj))
DOTNET_RID ?= linux-x64
DOCKER_IMAGE ?= bedrock-gophers-minimal
DOCKER_VOLUME ?= bedrock-gophers-minimal-data

.PHONY: check-revision prepare build run docker-build docker-run clean

check-revision:
	@test "$(GO_FRAMEWORK_REV)" = "$(BEDROCK_GOPHERS_SHORT_REV)" || (printf 'Go framework is pinned to %s, expected %s\n' "$(GO_FRAMEWORK_REV)" "$(BEDROCK_GOPHERS_SHORT_REV)"; exit 1)

$(CACHE)/.git:
	mkdir -p .cache
	git clone --quiet https://github.com/bedrock-gophers/plugins $(CACHE)

prepare: $(CACHE)/.git
	git -C $(CACHE) fetch --quiet origin
	git -C $(CACHE) checkout --quiet $(BEDROCK_GOPHERS_REV)

build: check-revision prepare
	dotnet publish $(CACHE)/csharp/Dragonfly.Runtime/Dragonfly.Runtime.csproj \
		-c Release -r $(DOTNET_RID) --self-contained true -o $(BUILD)/dotnet/runtime
	mkdir -p $(BUILD)/dotnet/plugins $(BUILD)/bin lib
	rm -f plugins/*.so
	@set -eu; for project in $(PLUGIN_PROJECTS); do \
		name=$$(basename "$$(dirname "$$project")"); \
		dotnet publish "$$project" -c Release -r $(DOTNET_RID) --self-contained true \
			-p:DragonflyFrameworkRoot=$(abspath $(CACHE)) -o "$(BUILD)/dotnet/plugins/$$name"; \
		find "$(BUILD)/dotnet/plugins/$$name" -maxdepth 1 -type f -name '*.so' -exec cp {} plugins/ \;; \
	done
	cp $(BUILD)/dotnet/runtime/Dragonfly.Runtime.so lib/libdragonfly_plugin_runtime.so
	go build -o $(BUILD)/bin/server .

run: build
	$(BUILD)/bin/server

docker-build:
	docker build --tag $(DOCKER_IMAGE) .

docker-run: docker-build
	docker run --rm --init \
		--publish 19132:19132/udp \
		--volume $(DOCKER_VOLUME):/app/.data \
		$(DOCKER_IMAGE)

clean:
	rm -rf .cache .data .build lib
	rm -f plugins/*.so
	find plugins -type d \( -name bin -o -name obj \) -prune -exec rm -rf {} +
