# ── Stage 1: Build ───────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# ── Stage 2: Runtime ─────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server /app/server
COPY web/ /app/web/
COPY internal/database/migrations/ /app/migrations/

WORKDIR /app
EXPOSE 8080

CMD ["/app/server"]
