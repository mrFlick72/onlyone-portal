# AGENTS.md

Guidance for AI coding agents (Claude Code, Codex, etc.) working in `budget/analytic-api`.

This is a Python 3.12+ FastAPI microservice that computes budget-expense analytics for the OnlyOne Portal. It owns no data: it pulls expenses from `budget-api` over REST (propagating the caller's JWT) and transforms them. The parent `../CLAUDE.md` and `../../CLAUDE.md` cover the budget subtree and the monorepo (auth model, deployment, sibling services).

## Commands

All commands run from `budget/analytic-api/`.

```bash
# Setup (dev — includes test dependencies)
python3 -m venv venv && source venv/bin/activate
pip install -e ".[dev]"

# Type-check (configured in pyproject.toml, scans src/, strict-ish: disallow_untyped_defs)
mypy

# Run all tests
pytest

# Run a single test file
pytest tests/analytic/api/test_end_point.py

# Run a single test by name
pytest tests/analytic/api/test_end_point.py::test_hello_returns_200

# Start locally (listens on 0.0.0.0:8045)
ANALYTIC_API_CONFIG_FILE_LOCATION=local/.env analytic-api
```

Dependencies are declared entirely in `pyproject.toml` (hatchling build). There is no `requirements.txt`. The `[dev]` extra installs test/type tools; the base install (`pip install .`, used by the Dockerfile) is production-only. The package is installed in editable mode, so `app.*` imports resolve without `PYTHONPATH`.

Tests run with auth middleware disabled and a dummy CORS origin via `[tool.pytest.ini_options].env` in `pyproject.toml` (`WITH_MIDDLEWARE=false`, `CORS_ALLOWED_ORIGINS=http://localhost`) — no env setup is needed to run `pytest`.

## Architecture

### Request lifecycle

`main.py` (logging setup + uvicorn, port 8045) → `server.py` (FastAPI app) → CORS middleware → `SecurityContextInjectorFilter` → router handler

`server.py` is the wiring point: it loads the `.env` file via `python-dotenv` (path from `ANALYTIC_API_CONFIG_FILE_LOCATION`), builds the DI container, registers middleware, and mounts the `health` and `analytic` routers. It refuses to start if `CORS_ALLOWED_ORIGINS` is empty — wildcard origins are deliberately not allowed because credentials are enabled.

### Routes

| Method | Path                          | Auth | Notes                                                            |
|--------|-------------------------------|------|------------------------------------------------------------------|
| GET    | `/health`                     | no   | Liveness probe, empty 200                                        |
| GET    | `/api/analytic/hello`         | yes  | Smoke endpoint                                                   |
| PUT    | `/api/analytic/budget/expense`| yes  | Body `{year, month, tags}` → list of `{date, amount, tag_values}` |

Sample authenticated request: `local/sample_requests.md`.

### Dependency injection

`dependency-injector` with a two-level container hierarchy:

```
ApplicationContainer                  (src/app/container.py)
├── SecurityConfigContainer           (src/app/infrastructure/security/security_container.py)
│   └── security_context_resolver     → LocalThreadSecurityContextResolver (singleton)
└── AnalyticConfigContainer           (src/app/analytic/container.py)
    └── expense_loader                → RestExpenseLoader (singleton; gets the resolver injected)
```

`ApplicationContainer` is instantiated once in `server.py` and wired to endpoint modules via `application_container.wire(modules=["app.analytic.api.end_point"])`. Handlers receive services via `Annotated[..., Depends(Provide[ApplicationContainer....])]` plus the `@inject` decorator.

When adding a new service: declare a provider in the feature container, inject it in the endpoint with `Provide[...]`, and add the endpoint module to the `wire(modules=[...])` list in `server.py`.

Note: `server.py`, `container.py`, and `analytic/container.py` each call `load_dotenv` at import time, and `AnalyticConfigContainer` reads `BUDGET_API_BASE_URL` at class-definition time — env vars must be set before `app.*` modules are imported.

### Authentication

`SecurityContextInjectorFilter` (`src/app/infrastructure/middleware/`, a Starlette `BaseHTTPMiddleware`) intercepts every request except `GET /health` and `OPTIONS`. It:

1. Validates the `Authorization: Bearer <token>` header — `401` if missing or malformed
2. Verifies the JWT (RS256) against JWKS fetched once at startup from `IDP_ISS/oauth2/jwks` — `401` on unknown `kid`. Because the fetch happens in the middleware constructor, the IDP must be reachable when the server boots (unless `WITH_MIDDLEWARE=false`)
3. Checks `issuer` against `IDP_ISS`; audience verification is intentionally disabled (`verify_aud: False`) because the IDP does not populate `aud` for this flow
4. Stores a `SecurityContext(token, user_name)` (user name from the `user_name` claim, configured in `server.py`) in a thread-local via `LocalThreadSecurityContextResolver`

Retrieve the current context in downstream code via an injected `SecurityContextResolver` → `get_security_context()`. `RestExpenseLoader` uses the stored raw token to propagate the caller's identity on outbound calls to budget-api.

### Outbound dependency: budget-api

`RestExpenseLoader` (`src/app/analytic/adapter/service.py`) is the only adapter. It calls `PUT {BUDGET_API_BASE_URL}/api/budget/expense` with `{month, year, searchTagList}` and the propagated bearer token, then flattens budget-api's `dailyBudgetExpenseRepresentationList` response into domain `ExpenseRecord`s. When the budget-api expense contract changes, this parsing is the impact point.

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
| `BUDGET_API_BASE_URL` | yes | Base URL of budget-api for expense loading (e.g. `http://local.budget-api.onlyone-portal.com:3035`) |
| `WITH_MIDDLEWARE` | no | Set to `false` to skip JWT auth (tests only) |
| `LOG_LEVEL` | no | Logging level, default `INFO` |
| `LOG_FILE_LOCATION` | no | Log file path, default `logs/app.log` (rotating, 10 MB × 5) |
