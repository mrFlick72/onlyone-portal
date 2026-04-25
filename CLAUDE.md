# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

OnlyOne Portal is a cloud-native microservices application. All services share a common OAuth2/JWT authentication model via the vauthenticator IDP.

## Architecture

```
portal/application-shell/   # React 19 + TypeScript SPA (Webpack) — see portal/application-shell/CLAUDE.md
account/account-api/        # Go (Gin) — user account management via vauthenticator REST
budget/budget-api/          # Go (Gin) — budget expense + revenue (DynamoDB) — see budget/CLAUDE.md
budget/revenue-api/         # Python FastAPI — legacy revenue service, pending decommission
budget/budget-exporter/     # Python — data export job
budget/analytic/            # Python — analytics/visualization scripts
tagging/tag-api/            # Go (Gin) — transaction tagging (DynamoDB)
plan/plan-service/          # Go (Echo) — todo/plan management (MySQL) — does NOT use the shared framework
core-services/golang-web-framework/  # Shared Go library (Gin, JWT, CORS, logging, caching)
idp/                        # OAuth2 token helper scripts for local dev
```

**Sub-CLAUDE.md files** contain deeper per-subtree guidance: `budget/CLAUDE.md`, `budget/budget-api/CLAUDE.md`, `portal/application-shell/CLAUDE.md`.

## Build & Run Commands

### Frontend (`portal/application-shell/src/`)
```bash
npm install
npm run build              # development build → ../dist
npm run production-build   # production build
npm run watch              # watch mode
```

### Go services using the shared framework (budget-api, tag-api, account-api)
Run from each service directory. All require `CONFIG_FILE_LOCATION` env var pointing to a YAML config file.
```bash
go build -o app .
go test -tags test ./...               # `-tags test` required; adapter tests need LocalStack
go test -tags test ./domain/... ./web/... # unit tests only, no infra needed
CGO_ENABLED=0 GOOS=linux go build -o app . # cross-compile for Docker/Alpine
```

Test fixtures (e.g. `fixture.go` files) are guarded by `//go:build test`. Running `go test ./...` without the tag will fail to compile.

### Plan service (`plan/plan-service/`) — different setup
Uses Echo (not the shared Gin framework) and reads its own env vars directly via Viper (no `CONFIG_FILE_LOCATION`):
```bash
go build -o app .
go test ./...              # no build tag needed
```
Required env vars: `DATABASE_URL`, `DATABASE_USER`, `DATABASE_PASSWORD`, `WEB_SERVER_PORT`, `LOGGING_FILE_NAME`.

### Python services
```bash
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt && pip install .
pytest
# revenue-api requires:
export BUDGET_API_CONFIG_FILE_LOCATION=<path-to-config>
```

## Testing with LocalStack (DynamoDB)

Go adapter tests for budget-api and tag-api run against LocalStack. Start it first:
```bash
cd <service>/test && docker compose up -d   # localstack/localstack:3.2 on :4566
```
Required env vars (any non-empty value works):
```
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
AWS_DEFAULT_REGION=us-east-1
```

## Authentication Flow

- IDP: vauthenticator at `http://local.api.vauthenticator.com:9090`
- JWKS endpoint: `http://local.api.vauthenticator.com:9090/oauth2/jwks`
- Frontend uses OAuth2 PKCE flow; entry points: `callback.html`, `logout.html`
- Local token generation helpers: `idp/get_access_token.py` / `idp/get_access_token.sh`

## Shared Go Framework (`core-services/golang-web-framework`)

All Gin-based services (budget-api, tag-api, account-api) depend on this via a local `replace` directive:
```
github.com/mrflick72/onlyone-portal/core-services/golang-web-framework => ../../core-services/golang-web-framework
```

**What `server.WebServerProvisioner` auto-configures:**
- CORS (origins from config key `cors.allowed.origins`)
- JWT validation middleware (`security.SetUpOAuth2()`) — reads `idp.jwks-endpoint` and `user.required-role`
- Health endpoint at `GET /management/health` → `{"status": "UP"}` (no auth required)

**JWT middleware behavior:**
- Skips `/management/*` paths and `OPTIONS` requests automatically
- Validated user is stored in Gin context under key `"user"` as `security.User{UserName, Authorities, AccessToken}`
- Retrieve in handlers via `security.GetCurrentUser(ctx)` after converting Gin context with `server.GinContextToPlainContextFactory`
- JWT claims used: `user_name` (string) and `authorities` ([]string)

**Config**: all Gin services load a YAML config file via Viper; path set in `CONFIG_FILE_LOCATION` env var. Access values via `config.GetConfigurationManagerInstance().GetConfigFor("key")`.

## Service Routes Summary

| Service | Method | Path | Notes |
|---------|--------|------|-------|
| tag-api | GET | `/api/tags` | Returns all tags for authenticated user |
| tag-api | PUT | `/api/tags` | Creates tag; UUID key is generated server-side |
| account-api | GET | `/api/account/user-account` | Proxies to vauthenticator REST at `idp.base-url` |
| account-api | PUT | `/api/account/user-account` | Proxies to vauthenticator REST |
| plan-service | GET | `/todo-service/todo` | MySQL-backed, no JWT auth |
| plan-service | POST | `/todo-service/todo` | |
| plan-service | GET/DELETE | `/todo-service/todo/:id` | |

Budget-api and revenue-api routes: see `budget/CLAUDE.md`.

## Key Cross-Service Dependencies

- `budget-api` → `tag-api` (REST, config key `tag-api.base-url`; cached with Ristretto)
- `account-api` → vauthenticator IDP (REST, config key `idp.base-url`)
- `budget-api`, `tag-api` → DynamoDB (region hardcoded to `eu-central-1`)
- `plan-service` → MySQL
- All Gin services → vauthenticator JWKS for JWT validation

## Deployment

- Docker images per service; Go services use Alpine base, Python services use `python:3.14.3-alpine3.23`
- All Go services share a common Dockerfile: `core-services/docker/ubuntu.Dockerfile`
- Kubernetes Helm charts: `account/helm/`, `budget/helm-charts/`
- CI/CD via GitHub Actions (`.github/workflows/`)
- Local frontend dev: nginx on port 8070 via `portal/application-shell/local/docker compose up`; app served at `http://local.onlyone-portal.com:8070`
