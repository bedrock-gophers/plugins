# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.26.4
ARG DOTNET_TAG=10.0

FROM golang:${GO_VERSION}-bookworm AS go-toolchain

FROM mcr.microsoft.com/dotnet/sdk:${DOTNET_TAG} AS build
ARG TARGETARCH
ENV PATH="/usr/local/go/bin:${PATH}"
RUN apt-get update && apt-get install -y --no-install-recommends \
        clang \
        build-essential \
        git \
        make \
        zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*
COPY --from=go-toolchain /usr/local/go /usr/local/go
WORKDIR /src/minimal
COPY . .
RUN architecture="${TARGETARCH:-$(uname -m)}" \
    && case "${architecture}" in \
        amd64|x86_64) dotnet_arch=x64 ;; \
        arm64|aarch64) dotnet_arch=arm64 ;; \
        *) echo "unsupported architecture: ${architecture}" >&2; exit 1 ;; \
    esac \
    && make DOTNET_RID="linux-${dotnet_arch}" build \
    && mkdir -p /out/.data

FROM gcr.io/distroless/cc-debian12:nonroot
WORKDIR /app
COPY --from=build --chown=nonroot:nonroot /src/minimal/.build/bin/server /app/server
COPY --from=build --chown=nonroot:nonroot /out/.data /app/.data
COPY --chown=nonroot:nonroot server.toml /app/server.toml
COPY --from=build --chown=nonroot:nonroot /src/minimal/lib/libdragonfly_plugin_runtime.so /app/lib/libdragonfly_plugin_runtime.so
COPY --from=build --chown=nonroot:nonroot /src/minimal/plugins/*.so /app/plugins/
EXPOSE 19132/udp
VOLUME ["/app/.data"]
ENTRYPOINT ["/app/server"]
