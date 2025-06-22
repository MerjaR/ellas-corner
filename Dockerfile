# ───── Build stage ─────
FROM golang:1.24.4-alpine3.22 AS builder

# Install dependencies for building Go app
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy Go source code and build binary
COPY ./*.go ./
COPY ./internal/ ./internal/
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o forum-app .

# ───── Final stage ─────
FROM alpine:3.22

# Install SQLite CLI
RUN apk add --no-cache sqlite

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy built binary and required assets
COPY --from=builder /app/forum-app ./
COPY ./web/ ./web/
COPY ./data/ ./data/
COPY ./migrations/ ./migrations/

# Ensure proper ownership for runtime access
RUN chown -R appuser:appgroup /app

# Internal port the app listens on by default
EXPOSE 8080

# Switch to non-root user
USER appuser

# Start the app (uses env vars for config)
CMD ["./forum-app"]