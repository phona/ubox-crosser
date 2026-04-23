# ═══════════════════════════════════════════════════
# Multi-stage build for ubox-crosser binaries
# ═══════════════════════════════════════════════════

# --- Builder stage ---
FROM golang:1.23 AS builder

ARG BINARY=server
ARG GIT_SHA=unknown

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -X main.GitSHA=${GIT_SHA}" -o /app/crosser ./cmd/${BINARY}

# --- Runtime stage ---
FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/crosser /usr/local/bin/crosser

ENTRYPOINT ["crosser"]
