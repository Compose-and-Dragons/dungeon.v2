#FROM golang:1.24.0-alpine AS builder
FROM --platform=$BUILDPLATFORM golang:1.25.2-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /workspace

# Copy go modules files from parent
COPY go.mod go.sum ./

# Copy compose-dragons directory from parent
COPY compose-dragons ./compose-dragons

# Copy dungeon-crawler-mcp-server files
COPY dungeon-crawler-mcp-server/go.mod dungeon-crawler-mcp-server/go.sum ./dungeon-crawler-mcp-server/
COPY dungeon-crawler-mcp-server/*.go ./dungeon-crawler-mcp-server/
COPY dungeon-crawler-mcp-server/tools ./dungeon-crawler-mcp-server/tools
COPY dungeon-crawler-mcp-server/types ./dungeon-crawler-mcp-server/types
COPY dungeon-crawler-mcp-server/data ./dungeon-crawler-mcp-server/data

WORKDIR /workspace/dungeon-crawler-mcp-server

RUN <<EOF
go mod tidy
#go build
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build
EOF

FROM alpine:latest
# for health check
RUN apk --no-cache add ca-certificates wget
WORKDIR /app
COPY --from=builder /workspace/dungeon-crawler-mcp-server/dungeon-mcp-server .
COPY --from=builder /workspace/dungeon-crawler-mcp-server/data ./data
ENTRYPOINT ["./dungeon-mcp-server"]

