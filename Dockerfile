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
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install Chromium and runtime dependencies for chromedp
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    tzdata

# Set environment variables
ENV HOST=0.0.0.0
ENV PORT=5082
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV DB_PATH=/app/data/tokens.sqlite

# Create data directory for persistent SQLite storage
RUN mkdir -p /app/data

# Copy compiled binary
COPY --from=builder /app/app /app/app

EXPOSE 5082

VOLUME ["/app/data"]

CMD ["/app/app"]