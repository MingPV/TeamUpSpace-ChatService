# Multi-stage Dockerfile (main at ./cmd/app/main.go)
# Examples:
#  - Local ARM build: docker build --platform linux/arm64 -t mingpv/tus-chat-service:latest .
#  - Multi-arch buildx & push: docker buildx build --platform linux/amd64,linux/arm64 -t <registry>/tus-chat-service:TAG --push .

ARG GO_VERSION=1.24
ARG BINARY=tus-chat-service
ARG BUILD_TARGET=cmd/app

FROM golang:${GO_VERSION}-alpine AS builder
# make ARG available inside this stage
ARG BINARY
ARG BUILD_TARGET
ARG TARGETARCH

WORKDIR /src
RUN apk add --no-cache git ca-certificates

# cache modules
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct && go mod download

# copy source
COPY . .

# build binary for target arch
ENV CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH}
WORKDIR /src/${BUILD_TARGET}
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o /usr/local/bin/${BINARY} -ldflags="-s -w" .

# verify binary exists (will show listing in build logs)
RUN ls -la /usr/local/bin || (echo "binary not found" && exit 1)

FROM alpine:3.18 AS runtime
ARG BINARY
RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app app

COPY --from=builder /usr/local/bin/${BINARY} /usr/local/bin/${BINARY}
RUN chmod +x /usr/local/bin/${BINARY}
USER app

ENV APP_PORT=8000
EXPOSE ${APP_PORT}

ENTRYPOINT ["/usr/local/bin/tus-chat-service"]