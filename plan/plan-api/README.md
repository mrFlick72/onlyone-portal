# Plan API

Authenticated Go service for plan and todo management in OnlyOne Portal. It exposes `/api/plan` endpoints, validates JWTs through the shared framework, and persists plans/todos in Postgres.

## Features

- Create, list, read, and delete plans for the authenticated user.
- Add, edit, and delete todos inside a plan.
- Track todo status with validated transitions.
- Cascade-delete todos when a plan is deleted.
- Use shared Gin server setup, CORS, JWT middleware, config, logging, and health endpoint.

## Todo Status Workflow

| From | Allowed targets |
|------|-----------------|
| `TODO` | `IN_PROGRESS`, `ABORTED` |
| `IN_PROGRESS` | `TODO`, `DONE`, `ABORTED` |
| `DONE` | none |
| `ABORTED` | none |

Invalid transitions return `409 Conflict`. Unknown statuses return `400 Bad Request`.

## API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/plan` | List plans for the current user. |
| `GET` | `/api/plan/:id` | Get one plan with todos. |
| `POST` | `/api/plan` | Create a plan. Body: `{"title":"...","date":"YYYY-MM-DD"}`. |
| `DELETE` | `/api/plan/:id` | Delete a plan and its todos. |
| `PUT` | `/api/plan/:id/todo` | Add a todo. Body: `{"content":"...","date":"YYYY-MM-DD"}`. |
| `PUT` | `/api/plan/:id/todo/:todoId` | Update todo content/date. |
| `PUT` | `/api/plan/:id/todo/:todoId/status` | Change status. Body: `{"status":"IN_PROGRESS"}`. |
| `DELETE` | `/api/plan/:id/todo/:todoId` | Delete a todo. |

All API routes require a bearer token accepted by the shared JWT middleware.

## Local Development

Start Postgres:

```bash
cd test
docker compose up -d
```

Run the service:

```bash
cd ..
CONFIG_FILE_LOCATION=test/application.yml go run .
```

Build:

```bash
go build -o app .
```

## Tests

```bash
# Unit/domain/web tests.
go test -tags test ./domain/plan ./pkg/clock ./web/plan

# Postgres adapter tests.
CONFIG_FILE_LOCATION=test/application.yml go test -tags test ./adapter/plan/db

# Everything; requires Postgres.
CONFIG_FILE_LOCATION=test/application.yml go test -tags test ./...
```

## Configuration

The service reads configuration through the shared framework. Set `CONFIG_FILE_LOCATION` to a YAML file with these keys:

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

## Data Model

`scripts/init.sql` creates two tables:

- `plan`: `id`, `user_name`, `title`, `date`
- `todo`: `id`, `plan_id`, `user_name`, `date`, `content`, `status`

`todo.plan_id` references `plan.id` with `ON DELETE CASCADE`.

## License

See `LICENSE`.
