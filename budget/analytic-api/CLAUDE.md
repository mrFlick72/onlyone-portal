# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

All commands run from `budget/analytic-api/`.

```bash
# Setup
python3 -m venv venv && source venv/bin/activate
pip install -r requirements.txt

# Run all tests
PYTHONPATH=src pytest

# Run a single test file
PYTHONPATH=src pytest tests/analytic/api/test_end_point.py

# Run a single test by name
PYTHONPATH=src pytest tests/analytic/api/test_end_point.py::test_hello_returns_200

# Start locally
ANALYTIC_API_CONFIG_FILE_LOCATION=local/.env PYTHONPATH=src python src/app/main.py
```

`PYTHONPATH=src` is always required because the package root is `src/app`, not the project root.

## Architecture

### Request lifecycle

`main.py` → uvicorn → `server.py` (FastAPI app) → `UserNameInjectorFilter` middleware → router handler

`server.py` is the wiring point: it loads `.env` via `python-dotenv` (path from `ANALYTIC_API_CONFIG_FILE_LOCATION`), builds the DI container, registers middleware, and mounts routers.

### Dependency injection

The project uses `dependency-injector` with a two-level container hierarchy:

```
ApplicationContainer          (src/app/container.py)
├── UserConfigContainer       (src/app/user/container.py)
└── AnalyticConfigContainer   (src/app/analytic/container.py)
    └── user_config_container (injected dependency)
```

`ApplicationContainer` is instantiated once in `server.py` and wired to endpoint modules via `application_container.wire(modules=["app.analytic.api.end_point"])`. Endpoint handlers receive injected services via `Depends(Provide[...])`.

When adding a new service, declare it in `AnalyticConfigContainer`, add its `Provide[...]` dependency to the endpoint, and wire the new module in `server.py`.

### Authentication

`UserNameInjectorFilter` (Starlette `BaseHTTPMiddleware`) intercepts every request except `GET /health` and `OPTIONS`. It:
1. Validates the `Authorization: Bearer <token>` header — returns `401` if missing or malformed
2. Verifies the JWT signature against JWKS fetched from `IDP_ISS/oauth2/jwks` at startup
3. Checks `issuer` against `IDP_ISS`; audience verification is intentionally disabled (`verify_aud: False`) because the IDP does not populate `aud` for this flow
4. Stores the resolved `UserName` in a thread-local via `LocalThreadUserNameResolver`

Retrieve the current user in a handler via `user_config_container.user_name_resolver().get_user_name()`.

Tests bypass the middleware entirely by setting `WITH_MIDDLEWARE=false` (configured in `pytest.ini`).

### Adding a new feature

1. Create `src/app/<feature>/domain/` — domain types and service interfaces
2. Create `src/app/<feature>/api/end_point.py` — `APIRouter` with handlers
3. Create `src/app/<feature>/container.py` — `DeclarativeContainer` wiring services
4. Register the container in `ApplicationContainer` and wire the module in `server.py`
5. Mirror the test structure under `tests/<feature>/`

## Environment variables

| Variable | Required | Description |
|---|---|---|
| `ANALYTIC_API_CONFIG_FILE_LOCATION` | yes | Path to `.env` file loaded by `python-dotenv` at startup |
| `IDP_ISS` | yes | Issuer URL for JWKS fetch and JWT `iss` validation (e.g. `http://local.api.vauthenticator.com:9090`) |
| `CORS_ALLOWED_ORIGINS` | yes | Comma-separated list of allowed origins — server refuses to start if empty or unset |
| `WITH_MIDDLEWARE` | no | Set to `false` to skip JWT auth (tests only) |
| `LOG_LEVEL` | no | Logging level, default `INFO` |
| `LOG_FILE_LOCATION` | no | Log file path, default `logs/app.log` |
