# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Tech Stack

| Concern          | Choice                                                                                                                                                                                                                                                        |
|------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Language         | Go 1.25.1                                                                                                                                                                                                                                                     |
| Web framework    | Gin (`github.com/gin-gonic/gin v1.11.0`)                                                                                                                                                                                                                      |
| Persistence      | Postgres                                                                                                                                                                                                                                                      |
| ID generation    | google/uuid                                                                                                                                                                                                                                                   |
| Auth             | JWT validation via the shared `core-services/golang-web-framework` middleware; JWKS fetched from `http://local.api.vauthenticator.com:9090/oauth2/jwks`                                                                                                       |
| Shared framework | `github.com/mrflick72/onlyone-portal/core-services/golang-web-framework` — resolved via local `replace` directive in `go.mod` pointing to `../../core-services/golang-web-framework`                                                                          |
| Test assertions  | testify + go-playground/assert                                                                                                                                                                                                                                |
| Build tag        | `-tags test` is required for any test that imports shared fixtures. Helpers like `domain/tags/fixture.go` are guarded by `//go:build test` so they are never linked into the production binary. Running `go test ./...` without the tag will fail to compile. |

Config is read by the shared framework's `config.GetConfigurationManagerInstance()` (backed by Viper). The config file
path is set via the `CONFIG_FILE_LOCATION` env var.

## Commands

```bash
go build -o app .
go test -tags test ./...                        # all tests (require Postgres)
go test -tags test ./src/plan/...               # repository integration tests only
go test -tags test ./src/web/...                # endpoint integration tests only
go test -tags test -run TestGetAllTodo ./src/web/...  # single test
```

## Architecture

plan-service is a Go service using the **shared Gin framework** (`core-services/golang-web-framework`), consistent with
the other portal services. JWT validation is handled by the shared middleware.

**Request path:** `main.go` → `src/configuration/main.go` (wires Gin via `server.WebServerProvisioner` + repository) →
`src/web/todoEndpoints.go` (handlers) → `src/plan/` (Postgres repository).

**Domain packages:**

- `src/plan/` — domain model (`Todo`, `Plan`), repository interfaces, Postgres implementations
- `src/web/` — Gin handler struct `TodoEndpoints`; UUID assigned server-side on `POST`; date serialized as
  `"2006-01-02"` via `src/pkg/clock`
- `src/pkg/logging/` — Zap singleton (JSON to file + console)

**Shared framework behaviour (inherited):**

- `server.WebServerProvisioner` auto-configures CORS, JWT middleware, and `GET /management/health`
- JWT middleware skips `/management/*` and `OPTIONS`; authenticated user available via `security.GetCurrentUser(ctx)`
  after `server.GinContextToPlainContextFactory`
- Config values accessed via `config.GetConfigurationManagerInstance().GetConfigFor("key")`

**Notable unfinished work:** `PlanRepository` in `src/plan/planMysqlRepository.go` is a stub — the interface is defined
but methods are not implemented beyond no-op stubs.

## Configuration

Set `CONFIG_FILE_LOCATION` to a YAML config file path. Required keys:

```yaml
server:
  port: 5050

cors:
  allowed:
    origins: http://local.onlyone-portal.com:8070

idp:
  jwks-endpoint: http://local.api.vauthenticator.com:9090/oauth2/jwks

user:
  required-role: PLAN_READER

database:
  url: localhost
  name: todo
  user: root
  password: root

logger:
  level: info
  file-name: log.log
```

## Testing

All tests are integration tests hitting a real Postgres instance. Start one with:

```bash
cd src/test && docker compose up -d
```

Apply the schema from `src/scripts/init.sql` before running tests. Each test clears tables via `TRUNCATE` at start/end (
see `src/plan/databaseTestUtils.go`). The `-tags test` build tag is required.
