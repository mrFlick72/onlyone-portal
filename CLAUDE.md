# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

OnlyOne Portal is a cloud-native microservices application. All services share a common OAuth2/JWT authentication model via the vauthenticator IDP.

## Architecture

```
portal/application-shell/   # React 19 + TypeScript multi-page SPA (Vite 8, MUI 9) — see portal/application-shell/CLAUDE.md
portal/e2e/ai/              # Playwright-MCP-driven AI end-to-end tests (LOGIN, PLAN scenarios)
account/account-api/        # Go (Gin) — user account management via vauthenticator REST
budget/budget-api/          # Go (Gin) — budget expense + revenue + attachments (DynamoDB + S3); publishes expense events to Kafka — see budget/CLAUDE.md
budget/analytic-api/        # Python FastAPI — expense analytics; Kafka consumer projects into Postgres, served read-only — see budget/analytic-api/AGENTS.md
budget/revenue-api/         # Python FastAPI — legacy revenue service, pending decommission
budget/budget-exporter/     # Python — data export job
budget/analytic/            # Python — analytics/visualization scripts
tagging/tag-api/            # Go (Gin) — transaction tagging (DynamoDB)
plan/plan-api/              # Go (Gin) — plan/todo management (Postgres) — shared framework + JWT
core-services/golang-web-framework/  # Shared Go library (Gin, JWT, CORS, logging, caching, OTel, HTTP client)
idp/                        # OAuth2 token helper scripts for local dev
```

**Sub-CLAUDE.md files** contain deeper per-subtree guidance: `budget/CLAUDE.md`, `budget/budget-api/CLAUDE.md`, `budget/analytic-api/AGENTS.md`, `portal/application-shell/CLAUDE.md`, `plan/plan-api/CLAUDE.md`.

## Build & Run Commands

### Frontend (`portal/application-shell/src/`)
```bash
npm install
npm run dev                # Vite dev server on :3000
npm run build              # development build -> ../dist
npm run production-build   # production build
npm run type-check         # tsc --noEmit
```

### Go services using the shared framework (budget-api, tag-api, account-api, plan-api)
Run from each service directory. All require `CONFIG_FILE_LOCATION` env var pointing to a YAML config file.
```bash
go build -o app .
go test -tags test ./...               # `-tags test` required; adapter tests may need local infra
go test -tags test ./domain/... ./web/... # unit tests only, no infra needed
CGO_ENABLED=0 GOOS=linux go build -o app . # cross-compile for Docker/Alpine
```

Test fixtures (e.g. `fixture.go` files and `internal/test`) are guarded by `//go:build test`. Running `go test ./...` without the tag can fail to compile in packages that import those fixtures.

### Plan API (`plan/plan-api/`)
Uses the shared Gin framework, JWT middleware, and Postgres. Config is loaded through `CONFIG_FILE_LOCATION`.
```bash
cd plan/plan-api
go build -o app .
go test -tags test ./domain/plan ./pkg/clock ./web/plan

cd test && docker compose up -d
cd ..
CONFIG_FILE_LOCATION=test/application.yml go test -tags test ./adapter/plan/db
```

### Python services
```bash
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt && pip install .
pytest
# revenue-api requires:
export BUDGET_API_CONFIG_FILE_LOCATION=<path-to-config>
```

## Testing with LocalStack (DynamoDB + S3)

Go adapter tests for budget-api and tag-api run against LocalStack. budget-api also exercises S3 (attachment content).
Start it first:
```bash
cd <service>/test && docker compose up -d   # localstack/localstack:3.2 with DynamoDB + S3 on :4566
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

## End-to-End Testing

`portal/e2e/ai/` holds AI-agent-driven UI scenarios that run against the local stack with the Playwright MCP server. Each scenario is a Markdown brief (`LOGIN.md`, `PLAN.md`) listing acceptance criteria the agent verifies against PNG snapshots in `assertions/`. `prerequisite.md` describes how to launch Chrome in incognito with the insecure-origin flags needed to talk to the http-only local IDP, and `prompt.md` is the entry prompt the user feeds to the agent. Credentials come from `portal/e2e/ai/.env` (gitignored). There is no Jest/Vitest suite for the frontend — these scenarios and the Go test suites are the verification path.

## Shared Go Framework (`core-services/golang-web-framework`)

All Gin-based services (budget-api, tag-api, account-api, plan-api) depend on this via a local `replace` directive:
```
github.com/mrflick72/onlyone-portal/core-services/golang-web-framework => ../../core-services/golang-web-framework
```

**`server.WebServerProvisioner` is built around `WebServerConfigurer`s.** Each configurer owns one cross-cutting concern with a `Configure() error` step at boot and a `Dispose(ctx) error` step at shutdown. `ConfigureEngine()` registers them in this order (the order is the middleware order on the Gin engine):

1. `OTelConfigurer` — installs global OTel providers (trace/metric/log) via `otel.Setup(ctx)` and registers `otelgin.Middleware`; `/management/*` is filtered so health probes don't pollute traces. Failure is non-fatal: logs and falls back to a no-op shutdown so the service still boots without tracing.
2. `StandardMiddlewareConfigurer` — `gin.Logger()`, `gin.Recovery()`, CORS (origins from `cors.allowed.origins`).
3. `OAuth2Configurer` — calls `security.SetUpOAuth2(ctx)` which builds the JWKS-cached middleware. The lifetime `ctx` is owned by the configurer and cancelled on `Dispose`, stopping the JWKS refresh goroutine.

`management.RegisterEndpoints` then mounts `GET /management/health` → `{"status": "UP"}` (no auth).

**Lifecycle / graceful shutdown:**
- `StartEngine()` listens on `server.port` and blocks on SIGINT/SIGTERM. On signal it calls `srv.Shutdown(ctx)` to drain in-flight requests, then `provisioner.Shutdown(ctx)` which iterates `Dispose` on every configurer (cancels JWKS refresh, flushes OTel exporters).
- A `defer` in `StartEngine()` is the safety net for early returns (e.g. `ListenAndServe` failure) — same `Shutdown` path.
- `Shutdown` clears the configurer slice and `engine` so a second call is a no-op and a panicking `Configure` cannot leave the provisioner half-built.
- Tunable via Go duration strings: `server.read-timeout` (30s), `server.write-timeout` (30s), `server.idle-timeout` (120s), `server.read-header-timeout` (5s), `server.shutdown-timeout` (10s).
- `server.shutdown-timeout` is the single budget for the entire shutdown, including OTel exporter flushes.

**JWT middleware behavior:**
- `security.SetUpOAuth2(ctx)` takes a context — cancel it to stop the background JWKS refresh goroutine. (The `OAuth2Configurer` does this for you.)
- JWKS is fetched through `httpclient.NewHTTPClient()` so refresh requests appear as `JWKS refresh` client spans when OTel is enabled.
- Skips `/management/*` paths and `OPTIONS` requests automatically.
- Validated user is stored in Gin context under key `"user"` as `security.User{UserName, Authorities, AccessToken}`.
- Retrieve in handlers via `security.GetCurrentUser(ctx)` after converting Gin context with `server.GinContextToPlainContextFactory`.
- JWT claims used: `user_name` (string) and `authorities` ([]string).

**Crypto (`cypto` package):**
- `cypto.AesCbcCipher` — AES-CBC encrypt/decrypt with PKCS-style padding, base64-encoded output, random IV per encryption.
- `cypto.KeyRepository` — port for fetching `SymmetricKey` by id. `NewInMemoryKeyRepository()` reads a single key from config keys `key.in-memory.storage.key` (id) and `key.in-memory.storage.key-value` (raw key bytes).

**OpenTelemetry (`otel` package):**
- `otel.Setup(ctx)` is called automatically by `WebServerProvisioner.ConfigureEngine()`; services do not call it directly
- Initialises global `TracerProvider`, `MeterProvider`, `LoggerProvider` — all three signals to the same `otel.endpoint`
- `otel.enabled: false` (default) → no providers installed, zero overhead
- Config keys: `otel.enabled`, `otel.service-name`, `otel.protocol` (`http`/`grpc`), `otel.endpoint`, `otel.insecure`
- Resource includes `service.name`, `host.name`, `telemetry.sdk.*`, and `OTEL_RESOURCE_ATTRIBUTES` env var
- W3C TraceContext + Baggage propagation registered globally

**OTel-aware HTTP client (`httpclient` package):**
- Use `httpclient.NewHTTPClient()` for every outbound service-to-service HTTP call
- When `otel.enabled: true`: wraps transport with `otelhttp.NewTransport` — injects `traceparent`/`tracestate` and creates a client span
- When `otel.enabled: false`: returns a plain `&http.Client{}`
- **Requirement**: build requests with `http.NewRequestWithContext(ctx, ...)` — the transport reads the active span from the request context; `http.NewRequest` silently disables propagation

**OTel-aware AWS SDK v2 client (`awsclient` package):**
- Use `awsclient.LoadDefaultConfig(ctx, opts...)` instead of `aws_config.LoadDefaultConfig` when creating DynamoDB / other AWS clients
- When `otel.enabled: true`: appends `otelaws.AppendMiddlewares` to the SDK API middleware chain — each AWS API call becomes a traced child span
- When `otel.enabled: false`: returns a plain `aws.Config` with no overhead
- **Requirement**: pass the request `ctx` to every AWS API call (e.g. `client.PutItem(ctx, input)`, not `context.TODO()`) so the middleware can read the active span

**Config**: all Gin services load a YAML config file via Viper; path set in `CONFIG_FILE_LOCATION` env var. Access string values via `config.GetConfigurationManagerInstance().GetConfigFor("key")`, booleans via `GetConfigBoolFor("key")`, durations via `GetConfigDurationFor("key", default)` (returns `default` when missing or unparseable).

## Service Routes Summary

| Service      | Method     | Path                                             | Notes                                                                              |
|--------------|------------|--------------------------------------------------|------------------------------------------------------------------------------------|
| tag-api      | GET        | `/api/tags`                                      | Returns all tags for authenticated user                                            |
| tag-api      | PUT        | `/api/tags`                                      | Creates tag; UUID key is generated server-side                                     |
| account-api  | GET        | `/api/account/user-account`                      | Proxies to vauthenticator REST at `idp.base-url`                                   |
| account-api  | PUT        | `/api/account/user-account`                      | Proxies to vauthenticator REST                                                     |
| budget-api   | POST       | `/api/attachment`                                | Multipart upload (file + `budgetId`/`budgetType`/`date` + optional `attachmentId`) |
| budget-api   | GET        | `/api/attachment/metadata/:budgetType/:budgetId` | Lists attachments for a parent expense or revenue                                  |
| budget-api   | GET        | `/api/attachment/:attachmentId/content`          | Returns the raw file bytes with `Content-Disposition`                              |
| budget-api   | DELETE     | `/api/attachment/:attachmentId`                  | Deletes both the metadata row and the S3 object                                    |
| plan-api     | GET        | `/api/plan`                                      | Lists authenticated user's plans                                                   |
| plan-api     | POST       | `/api/plan`                                      | Creates a plan and returns `{ "id": "<uuid>" }`                                    |
| plan-api     | GET/DELETE | `/api/plan/:id`                                  | Gets or deletes a plan; todos are cascade-deleted in Postgres                      |
| plan-api     | PUT        | `/api/plan/:id/todo`                             | Adds a todo with default `TODO` status                                             |
| plan-api     | PUT        | `/api/plan/:id/todo/:todoId`                     | Updates todo content/date without changing status                                  |
| plan-api     | PUT        | `/api/plan/:id/todo/:todoId/status`              | Validates todo status transitions                                                  |
| plan-api     | DELETE     | `/api/plan/:id/todo/:todoId`                     | Deletes a todo                                                                     |
| analytic-api | PUT        | `/api/analytic/budget/expense/total-by-tag`      | Totals grouped by tag value for a year (optional month + tag-key filter)           |
| analytic-api | PUT        | `/api/analytic/budget/expense/total-by-year`     | Totals grouped by year over a range (optional tag-key filter); years zero-filled   |
| analytic-api | POST       | `/api/analytic/budget/expense/reindex`           | Rebuilds the caller's projection from budget-api over `{fromYear,toYear}` (recovery) |

Expense, revenue, and full attachment route details: see `budget/CLAUDE.md` and `budget/budget-api/CLAUDE.md`. Analytics route details: see `budget/analytic-api/AGENTS.md`.

## Key Cross-Service Dependencies

- `budget-api` → `tag-api` (REST, config key `tag-api.base-url`; cached with Ristretto; OTel trace propagation via `httpclient.NewHTTPClient()`)
- `account-api` → vauthenticator IDP (REST, config key `idp.base-url`)
- `budget-api`, `tag-api` → DynamoDB (region hardcoded to `eu-central-1`)
- `budget-api` → S3 (attachment content; bucket from config key `budget-api.s3.attachment.bucket-name`)
- `budget-api` → Kafka (publishes expense `CREATE`/`UPDATE`/`DELETE` events to topic `budget-api.expense`)
- `analytic-api` → Kafka (consumes `budget-api.expense`) → its own Postgres projection (the read path; no REST call to `budget-api`)
- `analytic-api` → `budget-api` (REST, config key `budget-api.base-url` via `BUDGET_API_BASE_URL`) — **only** for the per-user reindex/recovery endpoint, never on the read path
- `plan-api` → Postgres
- All Gin services, including `plan-api` → vauthenticator JWKS for JWT validation

## Deployment

- Docker images per service; Go services use Alpine base, Python services use `python:3.14.3-alpine3.23`
- Shared-framework Go services generally use `core-services/docker/ubuntu.Dockerfile`; `plan/plan-api` also has its own Dockerfile.
- Kubernetes Helm charts: `account/helm/`, `budget/helm-charts/`
- CI/CD via GitHub Actions (`.github/workflows/`)
- Local frontend dev: nginx on port 8070 via `portal/application-shell/local/docker compose up`; app served at `http://local.onlyone-portal.com:8070`
