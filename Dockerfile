# syntax=docker/dockerfile:1

# ── Stage 1: Build ───────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# ── Stage 2: Runtime ─────────────────────────────────────────────
FROM alpine:3.21 AS runtime

RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY web/ /app/web/
COPY internal/database/migrations/ /app/migrations/

RUN chown -R appuser:appgroup /app

USER 1001:1001

EXPOSE 8080

CMD ["/app/server"]
