#FROM golang:1.24.0-alpine AS builder
FROM --platform=$BUILDPLATFORM golang:1.25.2-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
# COPY go.mod .
# COPY main.go .
COPY . .

RUN <<EOF
go mod tidy 
#go build
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build
EOF

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
WORKDIR /app
COPY --from=builder /app/mcp-toc-toc-server .
ENTRYPOINT ["./mcp-toc-toc-server"]
