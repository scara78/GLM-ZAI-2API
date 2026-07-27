# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum* ./
RUN go mod download

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o token-collector .

# Runtime stage — needs Chromium
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache \
    chromium \
    ca-certificates \
    tzdata

# Tell chromedp where Chrome is
ENV CHROME_PATH=/usr/bin/chromium-browser
ENV CHROMEDP_NO_SANDBOX=1

COPY --from=builder /app/token-collector /app/token-collector

# COUNT  — how many tokens to collect (default 750)
# OUT    — output path (mount a volume here)
ENV COUNT=750
ENV OUT=/app/data/tokens.sqlite

# entrypoint.sh: skip if tokens.sqlite already has tokens, else collect
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

VOLUME ["/app/data"]

CMD ["/app/entrypoint.sh"]