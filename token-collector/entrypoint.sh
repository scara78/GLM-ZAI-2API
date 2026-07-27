#!/bin/sh
set -e

OUT="${OUT:-/app/data/tokens.sqlite}"
COUNT="${COUNT:-750}"

# Check if tokens already exist in the DB
if [ -f "$OUT" ]; then
    # Use sqlite3 if available, otherwise just skip
    if command -v sqlite3 >/dev/null 2>&1; then
        EXISTING=$(sqlite3 "$OUT" "SELECT COUNT(*) FROM tokens;" 2>/dev/null || echo "0")
    else
        EXISTING=0
    fi

    if [ "$EXISTING" -gt "0" ] 2>/dev/null; then
        echo "[token-collector] $EXISTING tokens already in $OUT — skipping collection."
        echo "[token-collector] Delete $OUT and restart this service to re-collect."
        exit 0
    fi
fi

echo "[token-collector] Starting collection of $COUNT tokens → $OUT"

/app/token-collector \
    --count "$COUNT" \
    --out "$OUT" \
    --headless=true

echo "[token-collector] Done. Tokens saved to $OUT"