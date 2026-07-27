# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Runtime dependencies only (no chromium needed — pure HTTP mode)
RUN apk add --no-cache \
    ca-certificates \
    tzdata

# Set environment variables
ENV HOST=0.0.0.0
ENV PORT=5082
ENV DB_PATH=/app/data/tokens.sqlite

# Create data directory for persistent SQLite storage
RUN mkdir -p /app/data

# Copy compiled binary
COPY --from=builder /app/app /app/app

EXPOSE 5082

VOLUME ["/app/data"]

CMD ["/app/app"]