# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with the plan API.

## Current State

`plan-api` is an authenticated Go service for plans and todos. It is fully wired through the shared Gin framework, validates JWTs through vauthenticator JWKS, and stores data in Postgres.

Older plan docs describing a different service/storage shape are obsolete. The current HTTP surface is `/api/plan...`.

## Tech Stack

| Concern | Choice |
|---------|--------|
| Language | Go 1.25.1 |
| Web framework | Gin (`github.com/gin-gonic/gin v1.11.0`) through `core-services/golang-web-framework` |
| Persistence | Postgres via `database/sql` and `github.com/lib/pq` |
| ID generation | `github.com/google/uuid`; plan IDs are generated in the repository, todo IDs in the web handler |
| Auth | Shared JWT middleware; JWKS from config key `idp.jwks-endpoint` |
| Config | Shared framework config manager backed by Viper; config file selected with `CONFIG_FILE_LOCATION` |
| Test assertions | `github.com/stretchr/testify` |
| Build tag | `-tags test` is required for packages that import `internal/test` |

## Commands

Run from `plan/plan-api`.

```bash
go build -o app .

# Unit/domain tests that do not need Postgres.
go test -tags test ./domain/plan ./pkg/clock ./web/plan

# Adapter tests require local Postgres.
cd test && docker compose up -d
cd ..
CONFIG_FILE_LOCATION=test/application.yml go test -tags test ./adapter/plan/db

# Full package run; requires Postgres because adapter tests are included.
CONFIG_FILE_LOCATION=test/application.yml go test -tags test ./...
```

If the default Go build cache is read-only in the execution environment, set `GOCACHE=/tmp/go-cache` for test commands.

## Architecture

Request path:

```text
main.go
  -> server.WebServerProvisioner from core-services/golang-web-framework
  -> web/plan Gin handlers
  -> domain/plan PlanRepository interface
  -> adapter/plan/db Postgres repository
```

Package layout:

```text
domain/plan/      Plan, Todo, TodoStatus, transition rules, repository port
adapter/plan/db/  Postgres implementation of PlanRepository
web/plan/         Gin endpoint registration, request/response mapping, tests
config/           repository and Postgres DSN construction from shared config
pkg/clock/        YYYY-MM-DD date parsing/formatting helpers
pkg/database/     sql.Open and resource cleanup helpers
internal/test/    test fixtures guarded by //go:build test
scripts/init.sql  Postgres schema for plan and todo tables
```

## Domain Model

```go
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
    Status   TodoStatus
}
```

Todo statuses:

| Status | Meaning |
|--------|---------|
| `TODO` | Created and not started. |
| `IN_PROGRESS` | Work started. |
| `DONE` | Finished; terminal. |
| `ABORTED` | Cancelled; terminal. |

Allowed transitions:

| From | To |
|------|----|
| `TODO` | `IN_PROGRESS`, `ABORTED` |
| `IN_PROGRESS` | `TODO`, `DONE`, `ABORTED` |
| `DONE` | none |
| `ABORTED` | none |

The backend enforces these transitions in `TodoStatus.CanTransitionTo`. The frontend mirrors the same matrix in `portal/application-shell/src/plan/domain/TodoStatus.ts`.

## Repository Port

```go
type PlanRepository interface {
    GetAllPlanBy(userName string) ([]*Plan, error)
    GetPlan(idPlanId string, userName string) (*Plan, error)
    CreateNewPlan(p Plan) (string, error)
    DeletePlan(idPlanId string, userName string) error
    AddTodo(idPlanId string, t Todo) error
    UpdateTodo(idPlanId string, t Todo) error
    UpdateTodoStatus(idPlanId string, todoId string, status TodoStatus) error
    RemoveTodo(idPlanId string, todoId string) error
}
```

## Repository Implementation

`adapter/plan/db/repository.go` implements all repository methods.

| Method | Behavior |
|--------|----------|
| `CreateNewPlan` | Generates a UUID and inserts into `plan`; returns the generated ID. |
| `GetPlan` | Fetches a plan by `id AND user_name`; eagerly loads todos; returns not-found error when missing. |
| `GetAllPlanBy` | Fetches all plans for one user; initializes non-nil todo slices. |
| `DeletePlan` | Deletes by `id AND user_name`; Postgres cascade deletes todos. |
| `AddTodo` | Inserts a todo with status persisted from the domain value, normally `TODO`. |
| `UpdateTodo` | Updates todo content and date only; status is intentionally unchanged. |
| `UpdateTodoStatus` | Updates the status column for one todo. |
| `RemoveTodo` | Deletes a todo by `id AND plan_id`. |

Important details:

- Schema lives in `scripts/init.sql`. The script is idempotent: `CREATE TABLE IF NOT EXISTS` plus an `ALTER TABLE todo ADD COLUMN IF NOT EXISTS status` so it can run against both fresh and pre-status databases.
- `todo.plan_id -> plan.id` uses `ON DELETE CASCADE`; deleting a plan removes its todos in one statement.
- `loadTodosFor` returns a non-nil slice and normalizes timestamps to UTC.
- Each repository method opens a new connection via `database.GetDatabaseConnectionFor` (`sql.Open("postgres", dsn)`) and closes rows/stmt/db through `database.CloseResources`. There is no shared pool held by the repository struct.
- Only `GetAllPlanBy`, `GetPlan`, and `DeletePlan` filter by `user_name`. `UpdateTodo`, `UpdateTodoStatus`, and `RemoveTodo` filter by `id` (and `plan_id` for todo rows) only — ownership must be enforced by the caller. `changeTodoStatus` does this by loading the plan for the current user before issuing the update; `updateTodo` and `removeTodo` currently rely on todo IDs being unguessable.

## Web API

Routes are registered under `/api` in `web/plan/endpoints.go`.

| Method | Path | Handler | Success |
|--------|------|---------|---------|
| `GET` | `/api/plan` | `getAll` | `200` JSON array |
| `GET` | `/api/plan/:id` | `getOne` | `200` JSON object, `404` when missing |
| `POST` | `/api/plan` | `save` | `201 {"id":"<uuid>"}` |
| `DELETE` | `/api/plan/:id` | `delete` | `204` |
| `PUT` | `/api/plan/:id/todo` | `addTodo` | `201` |
| `PUT` | `/api/plan/:id/todo/:todoId` | `updateTodo` | `204` |
| `PUT` | `/api/plan/:id/todo/:todoId/status` | `changeTodoStatus` | `204`; `400` unknown status, `404` missing plan/todo, `409` invalid transition |
| `DELETE` | `/api/plan/:id/todo/:todoId` | `removeTodo` | `204` |

Representations in `web/plan/representations.go`:

```json
{
  "id": "plan-id",
  "user_name": "user-name",
  "title": "Plan title",
  "date": "2026-05-10",
  "todos": [
    {
      "id": "todo-id",
      "user_name": "user-name",
      "date": "2026-05-10",
      "content": "Do something",
      "status": "TODO"
    }
  ]
}
```

Handler conventions:

- The authenticated user comes from `security.GetCurrentUser` after `factory.CreateContextFromGin(c)`; request bodies must not be trusted for `user_name`. Every handler except `removeTodo` performs this lookup and returns `401` when it fails.
- Dates in JSON use `YYYY-MM-DD` through `pkg/clock`. `clock.ParseDateFor` returns the zero time on parse error rather than rejecting the request.
- Plan IDs are minted in the repository (`uuid.NewRandom`) and returned to the caller; todo IDs are minted in the `addTodo` handler (`uuid.NewString`) and persisted by `AddTodo`.
- New todos always start with `TODO`.
- `updateTodo` writes only `content` and `date`; the status column is intentionally left untouched.
- `changeTodoStatus` runs in this order: validate target with `TodoStatus.IsValid` (`400` on unknown), load the plan for the current user (`404` if missing), find the todo (`404` if missing), enforce `CanTransitionTo` (`409` on disallowed), then update. The frontend mirrors the same matrix in `portal/application-shell/src/plan/domain/TodoStatus.ts`.

## Configuration

Set `CONFIG_FILE_LOCATION` to a YAML file. Local config is `test/application.yml`.

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

## Testing

Tests are split one file per endpoint or behaviour rather than one big file. All web tests share `setupRouter` and `mockRepo` from `web/plan/endpoints_test.go`; assertions and mocking use `github.com/stretchr/testify`.

Web layer (`web/plan/`, no Postgres needed but `-tags test` is still required because the package imports `internal/test` fixtures):

- `endpoints_test.go` — `setupRouter` (injects a fake `security.User` into the Gin context) and the `mockRepo` shared by every other test in the package.
- `save_plan_test.go`, `get_plan_test.go`, `get_all_plans_test.go`, `delete_plan_test.go` — plan CRUD.
- `add_todo_test.go`, `update_todo_test.go`, `remove_todo_test.go` — todo CRUD.
- `change_todo_status_test.go` — every transition outcome (`204`, `400`, `404`, `409`).

Domain layer (`domain/plan/`, no infra, no build tag would also work but `-tags test` is harmless):

- `status_test.go` — `IsValid` and the full `CanTransitionTo` truth table.

Clock helpers (`pkg/clock/clock_test.go`) — round-trip parsing/formatting of `YYYY-MM-DD`.

Adapter layer (`adapter/plan/db/`, requires Postgres on `localhost:5432` with `postgres/postgres/postgres`):

```bash
cd test && docker compose up -d
cd ..
CONFIG_FILE_LOCATION=test/application.yml go test -tags test ./adapter/plan/db
```

- `repository_test.go` — `TestMain` runs `test.ClearDatabase()` once for the whole package.
- `create_plan_test.go`, `get_plan_test.go`, `get_all_plans_test.go`, `delete_plan_test.go`, `add_todo_test.go`, `update_todo_test.go`, `update_todo_status_test.go`, `remove_todo_test.go`.

`internal/test/test_utils.go` is guarded by `//go:build test` and exposes:

- `TestDSN` — the Postgres DSN used by adapter tests.
- `ANewPlan()` / `ANewTodoWith(content)` — fixtures with `UserName: "user-name"` and today's date in UTC.
- `ClearDatabase()` — re-runs `scripts/init.sql` and `TRUNCATE TABLE todo, plan`.
- `InitDatabase(conn)` — applies `scripts/init.sql` via the supplied connection.

Adapter tests do not reset state between each test, so use distinct users or fresh plans when isolation matters.
