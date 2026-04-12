# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

OnlyOne Portal is a cloud-native microservices application. All services share a common OAuth2/JWT authentication model via the vauthenticator IDP.

## Architecture

```
portal/application-shell/   # React 19 + TypeScript SPA (Webpack)
account/account-api/        # Go — user account management
budget/budget-api/          # Go — budget expense tracking
budget/revenue-api/         # Python FastAPI — revenue data
budget/budget-exporter/     # Python — data export
budget/analytic/            # Python — analytics/visualization
tagging/tag-api/            # Go — transaction tagging
plan/plan-service/          # Go (Echo) — todo/plan management
core-services/golang-web-framework/  # Shared Go library (Gin, JWT, CORS, logging, caching)
idp/                        # OAuth2 token helper scripts for local dev
```

All Go services depend on `core-services/golang-web-framework` via a local `replace` directive in their `go.mod`:
```
github.com/mrflick72/onlyone-portal/core-services/golang-web-framework => ../../core-services/golang-web-framework
```

## Build & Run Commands

### Frontend (`portal/application-shell/src/`)
```bash
npm install
npm run build              # development build → ../dist
npm run production-build   # production build
npm run watch              # watch mode
```

### Go Services (run from each service directory)
```bash
go build -o app .
go test -tags test ./...

# Cross-compile for Linux (used in Docker):
CGO_ENABLED=0 GOOS=linux go build -o app .
```

### Python Services (run from each service directory)
```bash
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt && pip install .
pytest
```

### Revenue API specifically
```bash
# Config file location set via env var:
export BUDGET_API_CONFIG_FILE_LOCATION=<path-to-config>
```

## Testing with LocalStack (DynamoDB)

Go and Python tests that touch DynamoDB use LocalStack. Required environment variables:
```
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
AWS_DEFAULT_REGION=us-east-1
```
Start LocalStack (`localstack/localstack:3.2`) before running tests locally.

## Plan Service (MySQL)

The plan service also uses MySQL. Schema is at `plan/plan-service/src/scripts/init.sql`. Connection configured via `.env` or environment variables including `WEB_SERVER_PORT`.

## Authentication Flow

- IDP: vauthenticator at `http://local.api.vauthenticator.com:9090`
- JWKS endpoint: `http://local.api.vauthenticator.com:9090/oauth2/jwks`
- Frontend uses OAuth2 PKCE flow; entry points: `callback.html`, `logout.html`
- All backend services validate JWT tokens using the shared framework middleware
- Local token generation helpers: `idp/get_access_token.py` / `idp/get_access_token.sh`

## Frontend Entry Points

Webpack builds separate bundles per page (configured in `webpack.config.js`):
- `callback`, `logout`, `home`, `budget`, `account`

Environment-specific config is handled via `dotenv-webpack` (development vs production).

## Core Go Framework

`core-services/golang-web-framework` provides reusable Gin-based infrastructure:
- CORS middleware
- JWT validation middleware
- Structured file-based logging (lumberjack)
- Caching (ristretto)
- Default server port: 3050

When modifying shared middleware, check all Go services for impact.

## Deployment

- Docker images per service; Go services use Alpine, Python services use `python:3.14.3-alpine3.23`
- Kubernetes Helm charts: `account/helm/`, `budget/helm-charts/`
- CI/CD via GitHub Actions (`.github/workflows/`)
