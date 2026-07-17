.PHONY: generate check-generated build-native build-server build stage-examples run test clean

CSHARP_BUILD := go run ./cmd/csharp-build -root .
DOTNET_RID ?=
GOEXE := $(shell go env GOEXE)

generate:
	go run ./cmd/csharp-gen -root .

check-generated:
	go run ./cmd/csharp-gen -root . -check

build-native: generate
	$(CSHARP_BUILD) -action publish -rid "$(DOTNET_RID)"

build-server:
	$(CSHARP_BUILD) -action prepare
	go build -o build/bedrock-gophers$(GOEXE) ./cmd/bedrock-gophers

build: build-native build-server

stage-examples: build-native
	$(CSHARP_BUILD) -action stage

run: stage-examples
	go run ./cmd/bedrock-gophers -config examples/server.toml

test: build-native check-generated
	dotnet build csharp/Dragonfly.Generator/Dragonfly.Generator.csproj -c Release
	go test ./...

clean:
	$(CSHARP_BUILD) -action clean
