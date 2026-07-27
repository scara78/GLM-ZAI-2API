# token-collector

Collects Aliyun captcha device tokens from chat.z.ai using headless Chromium.

## In Easypanel (Docker — no local Chrome needed)

Deploy the `token-collector` service from `easypanel.json`.  
It shares the `data` volume with `glm-zai-2api`.

**Steps:**
1. Deploy both services in Easypanel from this repo.
2. `token-collector` runs automatically on start, collects 750 tokens → `/app/data/tokens.sqlite`.
3. `glm-zai-2api` reads tokens from the same volume.
4. When tokens run out, restart the `token-collector` service in Easypanel to re-collect.

**Environment variables:**
| Variable | Default | Description |
|---|---|---|
| `COUNT` | `750` | How many tokens to collect (max 1250) |
| `OUT` | `/app/data/tokens.sqlite` | Output path (must match `DB_PATH` in main service) |

**Re-collect tokens:**
- Go to Easypanel → `token-collector` service → Restart.
- It detects existing tokens and skips if DB is not empty.
- To force re-collect: delete `/app/data/tokens.sqlite` via Easypanel file manager or shell, then restart.

## Local (Windows/Linux/Mac — needs Chrome installed)

```bash
cd token-collector
go mod tidy
go run . --count 750 --out ../tokens.sqlite
```