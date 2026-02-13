# ─── Build Stage ────────────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true

# Copy source and build
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o poll-api .

# ─── Runtime Stage ──────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/poll-api .

EXPOSE 5000

CMD ["./poll-api"]
