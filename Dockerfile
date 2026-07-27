# ── Stage 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Build deps (gcc needed for cgo sqlite fallback, but we use modernc pure-Go)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build main binary (pure Go — no CGO needed)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app .

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

# Chromium + runtime deps for headless browser (token collector)
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    chromium \
    chromium-chromedriver \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    font-noto \
    dbus \
    udev

# Tell collector where Chromium lives
ENV CHROME_PATH=/usr/bin/chromium-browser

# Server config defaults
ENV HOST=0.0.0.0
ENV PORT=5082
ENV DB_PATH=/app/data/tokens.sqlite

# Chromium sandbox workaround for Docker (no-sandbox flag set in code)
ENV CHROMIUM_FLAGS="--no-sandbox --disable-dev-shm-usage --disable-gpu"

RUN mkdir -p /app/data

COPY --from=builder /app/app /app/app

EXPOSE 5082

VOLUME ["/app/data"]

CMD ["/app/app"]