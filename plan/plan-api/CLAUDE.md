# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Tech Stack

| Concern          | Choice                                                                                                                                                                                               |
|------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Language         | Go 1.25.1                                                                                                                                                                                            |
| Web framework    | Gin (`github.com/gin-gonic/gin v1.11.0`)                                                                                                                                                             |
| Persistence      | Postgres                                                                                                                                                                                             |
| ID generation    | `github.com/google/uuid` — UUIDs are generated server-side inside the repository (not by the caller)                                                                                                 |
| Auth             | JWT validation via the shared `core-services/golang-web-framework` middleware; JWKS fetched from `http://local.api.vauthenticator.com:9090/oauth2/jwks`                                              |
| Shared framework | `github.com/mrflick72/onlyone-portal/core-services/golang-web-framework` — resolved via local `replace` directive in `go.mod` pointing to `../../core-services/golang-web-framework`                 |
| Test assertions  | `github.com/stretchr/testify/assert`                                                                                                                                                                 |
| Build tag        | `-tags test` required for integration tests. Files guarded by `//go:build test` are never linked into the production binary. Running `go test ./...` without the tag will fail to compile.           |

Config is read by the shared framework's `config.GetConfigurationManagerInstance()` (backed by Viper). The config file
path is set via the `CONFIG_FILE_LOCATION` env var.

## Commands

```bash
# ⚠ main.go and web/plan/endpoints_test.go currently do NOT compile — see "Known broken state" below
go test -tags test ./adapter/plan/db/...      # integration tests — the only tests that currently work
go test -tags test -run TestGetPlan ./adapter/plan/db/...  # single integration test
```

## Architecture

plan-api is a Go service using the **shared Gin framework** (`core-services/golang-web-framework`). JWT validation
is handled by the shared middleware.

**Request path:** `main.go` → `web/plan/` (Gin handlers) → `adapter/plan/db/` (Postgres repository) → `domain/plan/`
(pure domain types).

**Packages:**
Hexagonal layout — keep changes inside the right layer:

```
domain/
  plan/     # Plan, Todo structs + PlanRepository interface (both domain types live here)
adapter/
  plan/db/  # Postgres impl of PlanRepository — fully implemented
web/
  plan/     # Gin handlers — routes registered, handler bodies are stubs pending implementation
config/     # composition root — NewPostgresDSN
pkg/
  clock/    # date parsing/formatting (YYYY-MM-DD ↔ time.Time)
  database/ # sql.Open + CloseResources helpers
main.go     # WebServerProvisioner — needs updating to wire plan endpoints
```

**Shared framework behaviour (inherited):**

- `server.WebServerProvisioner` auto-configures CORS, JWT middleware, and `GET /management/health`
- JWT middleware skips `/management/*` and `OPTIONS`; authenticated user available via `security.GetCurrentUser(ctx)`
  after `server.GinContextToPlainContextFactory`
- Config values accessed via `config.GetConfigurationManagerInstance().GetConfigFor("key")`

## Domain Model

Both `Plan` and `Todo` live in the single `domain/plan` package:

```go
// domain/plan/model.go
type Plan struct {
    Id       string
    UserName string
    Title    string
    Date     time.Time
    Todos    []*Todo
}

type Todo struct {
    Id       string
    UserName string
    Date     time.Time
    Content  string
}
```

```go
// domain/plan/repository.go
type PlanRepository interface {
    GetAllPlanBy(userName string) ([]*Plan, error)
    GetPlan(idPlanId string, userName string) (*Plan, error)
    CreateNewPlan(p Plan) (string, error)   // returns generated plan ID
    AddTodo(idPlanId string, t Todo) error
    RemoveTodo(idPlanId string, todoId string) error
}
```

At the database level, one Plan may have zero or more Todos. The FK `todo.plan_id → plan.id` is set with `ON DELETE CASCADE`.

## Repository Implementation (`adapter/plan/db/repository.go`)

All methods implemented:

| Method | Behaviour |
|--------|-----------|
| `CreateNewPlan` | Generates UUID server-side; inserts into `plan` table; returns the new ID |
| `GetPlan` | Queries `plan` by `id AND user_name`; calls private `loadTodosFor` for eager loading; returns not-found error when missing |
| `GetAllPlanBy` | Queries all `plan` rows by `user_name`; calls `loadTodosFor` for each plan |
| `AddTodo` | Inserts into `todo` with `plan_id` FK set |
| `RemoveTodo` | Deletes from `todo` by `id AND plan_id` |

**Key internals:**
- `loadTodosFor(planId string)` — private method; `SELECT … FROM todo WHERE plan_id = $1`; always returns a non-nil slice; normalises Postgres TIMESTAMP → `time.UTC` via `.UTC()`
- `buildPlans(rows)` — row mapper; initialises with `make([]*Plan, 0)` (never nil); normalises date to UTC; sets `Todos` to `[]*Todo{}`
- All methods open a new connection per call via `database.GetDatabaseConnectionFor` and defer cleanup with `database.CloseResources`

## Web Routes (`web/plan/endpoints.go`)

Routes are registered but all handler bodies are empty stubs — logic is commented out:

| Method | Path | Handler |
|--------|------|---------|
| GET | `/api/plan` | `getAll` — stub |
| GET | `/api/plan/:id` | `getOne` — stub |
| POST | `/api/plan` | `save` — stub |
| PUT | `/api/plan/:id/todo` | `addTodo` — stub |
| PUT | `/api/plan/:id/todo/:todoId` | `updateTodo` — stub |
| DELETE | `/api/plan/:id/todo/:todoId` | `removeTodo` — stub |
| DELETE | `/api/plan/:id` | `delete` — stub |

The `todoRepresentation` struct and `TodoEndpoints` wiring are in place; only the handler bodies need implementing.

## Known Broken State

Two files do not compile and must be fixed before `go build` or `go test ./web/...` will work:

| File | Problem |
|------|---------|
| `main.go` | Imports `adapter/todo/db` and `web/todo` — both packages were removed. Needs to be rewritten to wire `adapter/plan/db` + `web/plan`. |
| `web/plan/endpoints_test.go` | Imports `domain/todo` (removed) and mocks `todo.TodoRepository` (old interface). Needs to be rewritten for `plan.PlanRepository`. |

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
  name: postgres
  user: postgres
  password: postgres

logger:
  level: info
  file-name: log.log
```

A ready-to-use local config is at `test/application.yml`.

## Testing

**Unit tests** (`web/plan/`) — currently broken; `endpoints_test.go` must be rewritten first (see "Known Broken State").

**Integration tests** (`adapter/plan/db/`) — require Postgres and `-tags test`:

```bash
cd test && docker compose up -d
CONFIG_FILE_LOCATION=test/application.yml go test -tags test ./adapter/...
```

Schema is applied automatically by `initDatabase()` in `test_utils.go` (reads `scripts/init.sql`). No manual
schema apply step is needed after the first run.

**Test lifecycle (`adapter/plan/db/`):**
- `TestMain` calls `clearDatabase()` once before all tests (`TRUNCATE TABLE todo, plan` — both in one statement to respect the FK constraint)
- Individual tests do **not** clean up between themselves — data accumulates across tests
- Tests that insert under `"user-name"` coexist; `TestGetAllPlanBy` uses the distinct username `"all-plans-user"` for isolation

**Test helpers (`adapter/plan/db/test_utils.go`, guarded by `//go:build test`):**
- `aNewPlan()` — returns a `plan.Plan` with empty `Todos: []*plan.Todo{}`; `Id` is not set (repo generates it)
- `aNewTodoWith(content)` — returns a `plan.Todo` with a random UUID `Id`
- `assertEqualPlan(t, expected, actual)` — compares `expected plan.Plan` against dereferenced `*actual`; set `expected.Id` from the value returned by `CreateNewPlan` before calling
- `clearDatabase()` — runs `TRUNCATE TABLE todo, plan` (single statement); safe to call from `TestMain`
