# ── Stage 1: build ────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/server ./cmd/main.go

# ── Stage 2: run ──────────────────────────────────────────────────────────────
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata && \
    adduser -D appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /app/bin/server ./server
COPY --from=builder --chown=appuser:appuser /app/migrations ./migrations

USER appuser

EXPOSE 8080

CMD ["./server"]

