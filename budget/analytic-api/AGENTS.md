# AGENTS.md

Guidance for AI coding agents (Claude Code, Codex, etc.) working in `budget/analytic-api`.

This is a Python 3.12+ FastAPI microservice that serves budget-expense analytics for the OnlyOne Portal. It keeps a **read-optimised projection** of the user's expenses in its own Postgres database, kept up to date by consuming `CREATE`/`UPDATE`/`DELETE` events that `budget-api` (the data owner) publishes to the Kafka topic `budget-api.expense`. Read requests are answered purely from the projection — `analytic-api` never calls `budget-api` at request time, and is eventually consistent with it. The parent `../CLAUDE.md` and `../../CLAUDE.md` cover the budget subtree and the monorepo (auth model, deployment, sibling services).

## Commands

All commands run from `budget/analytic-api/`.

```bash
# Setup (dev — includes test dependencies)
python3 -m venv venv && source venv/bin/activate
pip install -e ".[dev]"

# Type-check (configured in pyproject.toml, scans src/, strict-ish: disallow_untyped_defs)
mypy

# Run the fast suite (no infrastructure required)
pytest -m "not integration"

# Run the Postgres-backed projection tests (needs Docker — uses testcontainers)
pytest -m integration

# Run everything
pytest

# Start the Kafka broker + topic + Postgres projection store for local runs
(cd .. && docker compose up -d kafka init-kafka analytic-postgres)

# Start locally (listens on 0.0.0.0:8045)
ANALYTIC_API_CONFIG_FILE_LOCATION=local/.env analytic-api
```

Dependencies are declared entirely in `pyproject.toml` (hatchling build). There is no `requirements.txt`. The `[dev]` extra installs test/type tools; the base install (`pip install .`, used by the Dockerfile) is production-only. The package is installed in editable mode, so `app.*` imports resolve without `PYTHONPATH`.

Tests run with auth middleware disabled, a dummy CORS origin, and the Kafka consumer switched off via `[tool.pytest.ini_options].env` in `pyproject.toml` (`WITH_MIDDLEWARE=false`, `CORS_ALLOWED_ORIGINS=http://localhost`, `KAFKA_CONSUMER_ENABLED=false`) — so `pytest -m "not integration"` needs neither a broker nor a database. The `integration` marker covers the `PostgresExpenseProjectionRepository` tests, which spin up a throwaway Postgres via `testcontainers` and therefore require Docker.

## Architecture

### Request lifecycle

`main.py` (logging setup + uvicorn, port 8045) → `server.py` (FastAPI app) → CORS middleware → `SecurityContextInjectorFilter` → router handler

`server.py` is the wiring point: it loads the `.env` file via `python-dotenv` (path from `ANALYTIC_API_CONFIG_FILE_LOCATION`), builds the DI container, registers middleware, mounts the `health` and `analytic` routers, and owns the **`lifespan`** handler. When `KAFKA_CONSUMER_ENABLED` is true (the default) the lifespan opens the Postgres pool and starts the Kafka consumer on startup, then stops the consumer and closes the pool on shutdown. The app does **not** create the schema — see "Database schema" below. It refuses to start if `CORS_ALLOWED_ORIGINS` is empty — wildcard origins are deliberately not allowed because credentials are enabled.

### Routes

All analytics endpoints are scoped to the authenticated `user_name` from the JWT and read only from the local projection.

| Method | Path                                          | Auth | Notes                                                                                            |
|--------|-----------------------------------------------|------|--------------------------------------------------------------------------------------------------|
| GET    | `/health`                                     | no   | Liveness probe, empty 200                                                                        |
| PUT    | `/api/analytic/budget/expense/total-by-tag`   | yes  | Body `{year, month?, tags?}` (`tags` = tag **keys**) → `[{tag, total}]` (`tag` = tag **value**)  |
| PUT    | `/api/analytic/budget/expense/total-by-year`  | yes  | Body `{fromYear, toYear, tag?}` (`tag` = tag **key**) → `[{year, total}]`, every year zero-filled |

**Tag key-vs-value contract:** the frontend filters by `SearchTag.key` but labels/groups by `SearchTag.value`. The projection stores both; `total-by-tag` filters expenses by key (keeping all sibling tags of a matched expense) and aggregates/labels by value, while `total-by-year` counts each expense once.

### Dependency injection

`dependency-injector` with a two-level container hierarchy:

```
ApplicationContainer                  (src/app/container.py)
├── SecurityConfigContainer           (src/app/infrastructure/security/security_container.py)
│   └── security_context_resolver        → LocalThreadSecurityContextResolver (singleton)
└── AnalyticConfigContainer           (src/app/analytic/container.py)
    ├── expense_projection_repository    → PostgresExpenseProjectionRepository (singleton)
    ├── budget_expense_analysis_service  → BudgetExpenseAnalysisService (repo + resolver injected)
    ├── budget_expense_event_handler     → BudgetExpenseEventHandler (repo injected)
    └── budget_expense_consumer          → BudgetExpenseConsumer (handler injected)
```

`ApplicationContainer` is instantiated once in `server.py` and wired to endpoint modules via `application_container.wire(modules=["app.analytic.api.end_point"])`. Handlers receive services via `Annotated[..., Depends(Provide[ApplicationContainer....])]` plus the `@inject` decorator. The consumer and repository are resolved directly from the container inside the `lifespan` handler, not via `@inject`.

Note: `server.py`, `container.py`, and `analytic/container.py` each call `load_dotenv` at import time, and `AnalyticConfigContainer` reads `DATABASE_URL` and the `KAFKA_*` vars at class-definition time — env vars must be set before `app.*` modules are imported.

### Authentication

`SecurityContextInjectorFilter` (`src/app/infrastructure/middleware/`, a Starlette `BaseHTTPMiddleware`) intercepts every request except `GET /health` and `OPTIONS`. It:

1. Validates the `Authorization: Bearer <token>` header — `401` if missing or malformed
2. Verifies the JWT (RS256) against JWKS fetched once at startup from `IDP_ISS/oauth2/jwks` — `401` on unknown `kid`. Because the fetch happens in the middleware constructor, the IDP must be reachable when the server boots (unless `WITH_MIDDLEWARE=false`)
3. Checks `issuer` against `IDP_ISS`; audience verification is intentionally disabled (`verify_aud: False`) because the IDP does not populate `aud` for this flow
4. Stores a `SecurityContext(token, user_name)` (user name from the `user_name` claim, configured in `server.py`) in a thread-local via `LocalThreadSecurityContextResolver`

Retrieve the current context in downstream code via an injected `SecurityContextResolver` → `get_security_context()`. `BudgetExpenseAnalysisService` reads `user_name` from it to scope every projection query to the caller.

### The projection: Kafka consumer + Postgres

`budget-api` publishes plain-JSON events to the Kafka topic `budget-api.expense`:

```json
{ "action": "CREATE | UPDATE | DELETE",
  "payload": { "id", "userName", "date": "dd/MM/yyyy", "amount": "<decimal, scale 2>", "note", "tags": [{ "key", "value" }] } }
```

- **`BudgetExpenseConsumer`** (`src/app/analytic/adapter/kafka/consumer.py`) — an in-process `confluent-kafka` consumer that polls on a daemon thread. Offsets are committed only **after** the event is applied (at-least-once delivery; the projection ops are idempotent so redelivery is safe). It is started/stopped by the `server.py` lifespan.
- **`BudgetExpenseEventHandler`** (same module) — decodes a single event dict and applies it: `CREATE`/`UPDATE` → `repository.save`, `DELETE` → `repository.delete`. Kept free of any Kafka type so it is unit-testable with a plain dict.
- **`PostgresExpenseProjectionRepository`** (`src/app/analytic/adapter/db/repository.py`) — owns its `psycopg` connection pool (opened lazily from the lifespan). `save` upserts the expense row (`ON CONFLICT (id) DO UPDATE`) and **replaces** its tag rows (so re-tagging on `UPDATE` is handled); `delete` removes by id (tags cascade). `total_by_tag` / `total_by_year` are SQL `GROUP BY` aggregations over `budget_expense_projection` + `budget_expense_tag`. When the `budget-api.expense` event contract changes, `BudgetExpenseEventHandler._to_expense` is the impact point.

### Database schema

The schema lives in `scripts/init.sql` (service root), kept out of the application source — the adapter never runs DDL. It is applied out-of-band, mirroring `plan-api`:

- **Local dev**: `budget/docker-compose.yml` mounts it into the `analytic-postgres` container's `/docker-entrypoint-initdb.d/`, so it runs once on first init.
- **Tests**: the integration fixture reads and applies `scripts/init.sql` to the throwaway testcontainer.
- **Deployment**: applied by the environment's migration/provisioning step.

### Shared value objects

`src/app/time/domain/` (`Date`, `Month`, `Year`) and `src/app/money/domain/` (`Money`) mirror the Go value objects in `budget-api`'s `domain/time` and `domain/money`. `Date` round-trips two formats: `%d/%m/%Y` (`date_for`/`formatted_date` — the budget-api wire format) and ISO `%Y-%m-%d`. `Money` uses `Decimal` with scale 2 and `ROUND_HALF_DOWN`. Prefer these over raw `str`/`float` in domain code.

### Adding a new feature

1. `src/app/<feature>/domain/` — dataclasses and ABC service ports
2. `src/app/<feature>/adapter/` — implementations of the ports
3. `src/app/<feature>/api/end_point.py` — `APIRouter` + pydantic representations in `api/representation.py`
4. `src/app/<feature>/container.py` — `DeclarativeContainer` wiring the services
5. Register the container in `ApplicationContainer` and wire the module in `server.py`
6. Mirror the test structure under `tests/<feature>/`

## Environment variables

Loaded from the file pointed to by `ANALYTIC_API_CONFIG_FILE_LOCATION` (see `local/.env` for a working local example).

| Variable | Required | Description |
|---|---|---|
| `ANALYTIC_API_CONFIG_FILE_LOCATION` | yes | Path to `.env` file loaded by `python-dotenv` at startup |
| `IDP_ISS` | yes | Issuer URL for JWKS fetch and JWT `iss` validation (e.g. `http://local.api.vauthenticator.com:9090`) |
| `CORS_ALLOWED_ORIGINS` | yes | Comma-separated explicit origin list — server refuses to start if empty or unset |
| `DATABASE_URL` | yes | Postgres connection string for the projection store (e.g. `postgresql://analytic:analytic@localhost:5433/analytic`) |
| `KAFKA_BOOTSTRAP_SERVERS` | yes | Kafka bootstrap servers for the consumer (e.g. `localhost:9092`) |
| `KAFKA_EXPENSE_TOPIC` | no | Topic to consume, default `budget-api.expense` |
| `KAFKA_CONSUMER_GROUP` | no | Consumer group id, default `analytic-api` |
| `KAFKA_CONSUMER_ENABLED` | no | Set to `false` to skip opening the DB pool and starting the consumer (tests set this), default `true` |
| `WITH_MIDDLEWARE` | no | Set to `false` to skip JWT auth (tests only) |
| `LOG_LEVEL` | no | Logging level, default `INFO` |
| `LOG_FILE_LOCATION` | no | Log file path, default `logs/app.log` (rotating, 10 MB × 5) |
