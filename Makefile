BEDROCK_GOPHERS_REV := 7897b0b1e439786315be998eb61fdd44c924e6d6
BEDROCK_GOPHERS_SHORT_REV := $(shell printf '%.12s' $(BEDROCK_GOPHERS_REV))
GO_FRAMEWORK_REV := $(shell go list -m -f '{{.Version}}' github.com/bedrock-gophers/plugins | sed 's/.*-//')
CACHE := .cache/bedrock-gophers

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
	$(MAKE) -C $(CACHE) build-native
	mkdir -p lib plugins
	cp $(CACHE)/build/lib/libdragonfly_plugin_runtime.so lib/
	cp $(CACHE)/build/plugins/LifecycleLogger.so plugins/
	go mod download

run: build
	go run .

clean:
	rm -rf .cache .data lib plugins
