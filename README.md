# OnlyOne Portal

OnlyOne Portal is a cloud-native personal portal made of a React application shell and several backend services. The current feature set covers account management, budget expense/revenue tracking, attachment handling, tag search, and plan/todo management.

## Modules

| Path | Purpose |
|------|---------|
| `portal/application-shell` | React 19 + TypeScript Vite frontend. Builds separate HTML entries for home, auth callback/logout, budget pages, account, and plan pages. |
| `account/account-api` | Go service for user account management through vauthenticator REST APIs. |
| `budget/budget-api` | Go service for expenses, revenue, tags integration, and file attachments backed by DynamoDB/S3. |
| `tagging/tag-api` | Go service for user tag management backed by DynamoDB. |
| `plan/plan-api` | Go service for plans and todos backed by Postgres. |
| `core-services/golang-web-framework` | Shared Go framework for Gin server setup, JWT, CORS, config, logging, HTTP client, AWS client helpers, and OpenTelemetry. |
| `budget/revenue-api` | Legacy Python FastAPI revenue service. |

## Frontend

```bash
cd portal/application-shell/src
npm install
npm run dev
npm run type-check
npm run build
```

The local nginx wrapper serves the built app from `portal/application-shell/dist`:

```bash
cd portal/application-shell/local
docker compose up
```

Default local URL: `http://local.onlyone-portal.com:8070`.

## Go Services

Most Go services use the shared framework and load configuration from a YAML file selected by `CONFIG_FILE_LOCATION`.

```bash
go build -o app .
go test -tags test ./...
```

Adapter tests usually need local infrastructure:

- `budget/budget-api` and `tagging/tag-api`: LocalStack for DynamoDB/S3.
- `plan/plan-api`: Postgres from `plan/plan-api/test/docker-compose.yml`.

## Authentication

The portal uses OAuth2 Authorization Code with PKCE against vauthenticator. Backend services validate JWTs using the vauthenticator JWKS endpoint. The frontend stores `ACCESS_TOKEN` and `ID_TOKEN` in `sessionStorage` and sends bearer tokens to backend APIs.

## Main Routes

| Service | Routes |
|---------|--------|
| Portal | `/`, `/callback`, `/logout`, `/budget/expense/index`, `/budget/revenue/index`, `/budget/search-tags/index`, `/account/index`, `/plan/index`, `/plan/detail?id=...` |
| Plan API | `/api/plan`, `/api/plan/:id`, `/api/plan/:id/todo`, `/api/plan/:id/todo/:todoId`, `/api/plan/:id/todo/:todoId/status` |
| Budget API | Expense/revenue APIs plus `/api/attachment` metadata/content/delete endpoints. |
| Tag API | `/api/tags` |
| Account API | `/api/account/user-account` |

## More Documentation

- `CLAUDE.md`: repository-wide agent guidance.
- `portal/application-shell/README.md` and `portal/application-shell/CLAUDE.md`: frontend details.
- `plan/plan-api/README.md` and `plan/plan-api/CLAUDE.md`: plan service details.
- `budget/CLAUDE.md` and `budget/budget-api/CLAUDE.md`: budget service details.
