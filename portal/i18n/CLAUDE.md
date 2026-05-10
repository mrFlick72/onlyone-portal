# i18n Service

REST API that serves internationalization message bundles to portal UI applications (sections). Written in Go using the shared `core-services/golang-web-framework` (Gin + JWKS-cached JWT middleware + OTel + AWS SDK v2), matching the conventions of `budget-api`, `tag-api`, `account-api`, and `plan-api`.

## Purpose

The portal UI is split into multiple application sections (budget, plan, account, tagging, etc.). Each section needs localized message bundles. This service is the single source of truth: it loads YAML bundles from S3 and exposes them over HTTP, picking the language from the caller's JWT.

## Storage Layout (S3)

A single S3 bucket holds every bundle, keyed by application section and page:

```
<APPLICATION_NAME>/<PAGE>/message_bundle_<LANG>.yaml
```

- `<APPLICATION_NAME>` — portal section identifier (e.g. `budget`, `plan`, `account`).
- `<PAGE>` — page or feature within that application.
- `<LANG>` — locale tag, lower-case `<lang>_<region>` (e.g. `en_en`, `it_it`).

Example:
```
budget/expense-list/message_bundle_en_en.yaml
budget/expense-list/message_bundle_it_it.yaml
plan/todo-board/message_bundle_en_en.yaml
```

Bucket name comes from a config key (e.g. `i18n.s3.bundle.bucket-name`), loaded via `config.GetConfigurationManagerInstance()`.

## Language Selection

1. Read the authenticated user from the Gin context (`security.GetCurrentUser`) — populated by the shared `OAuth2Configurer`.
2. Resolve the preferred language from a claim on the JWT access token (e.g. `preferred_language` / `locale`). The exact claim name is config-driven so vauthenticator can evolve.
3. If the claim is missing or the bundle for that language does not exist in S3, fall back to the **default language** (`en_en`). Default is config-driven (`i18n.default-language`).

The service must never 500 on a missing localized bundle when the default is available — it falls back transparently and logs the miss.

## Planned HTTP API

All routes are JWT-protected (the framework's `OAuth2Configurer` covers everything outside `/management/*`).

| Method | Path                                          | Purpose                                                       |
|--------|-----------------------------------------------|---------------------------------------------------------------|
| GET    | `/api/i18n/:application/:page`                | Returns the message bundle for the caller's preferred lang    |
| GET    | `/api/i18n/:application/:page?lang=<locale>`  | Optional explicit override (still validated against bundles)  |

Response: parsed YAML → JSON object (`map[string]any`), so the UI consumes a flat/nested JSON tree without parsing YAML.

`GET /management/health` is provided by the framework.

## Architecture (target package layout)

Mirror the other Go services:

```
portal/i18n/
  main.go                          # bootstraps server.WebServerProvisioner, registers routes
  application.yml                  # local config (gitignored if it carries secrets)
  domain/i18n/                     # MessageBundle, Locale, BundleRepository port
  adapter/i18n/s3/                 # S3-backed BundleRepository (uses awsclient.LoadDefaultConfig)
  web/i18n/                        # Gin handlers, request/response DTOs, language resolution
  pkg/                             # small utilities (locale parsing, etc.)
  test/                            # docker-compose for LocalStack S3; integration fixtures
```

Key framework usages:
- `awsclient.LoadDefaultConfig(ctx, ...)` for the S3 client — gives traced AWS spans when OTel is on.
- `httpclient.NewHTTPClient()` for any outbound HTTP (not expected initially).
- `security.GetCurrentUser(ctx)` after `server.GinContextToPlainContextFactory` to read the JWT-derived user and language claim.
- Build configurers in `main.go` in the standard order: OTel → standard middleware → OAuth2 → routes → management.

## Caching

S3 reads are expensive per-request. Plan to cache parsed bundles in-process (Ristretto, as `budget-api` does for `tag-api` lookups), keyed by `<application>/<page>/<lang>`, with a TTL or explicit invalidation endpoint. Cache stays out of the first cut if it complicates the MVP.

## Config Keys (planned)

| Key                                | Purpose                                                  |
|------------------------------------|----------------------------------------------------------|
| `server.port`                       | HTTP port                                                |
| `cors.allowed.origins`              | CORS origins                                             |
| `security.jwks-uri`                 | vauthenticator JWKS endpoint                             |
| `i18n.s3.bundle.bucket-name`        | Bucket holding `<app>/<page>/message_bundle_<lang>.yaml` |
| `i18n.default-language`             | Fallback locale, defaults to `en_en`                     |
| `i18n.jwt.language-claim`           | JWT claim name carrying the user's preferred language    |
| `otel.*`                            | Standard OTel keys (see root `CLAUDE.md`)                |

## Build & Test

Same as the other shared-framework services:

```bash
go build -o app .
go test -tags test ./domain/... ./web/...                        # unit
cd test && docker compose up -d                                  # LocalStack S3
CONFIG_FILE_LOCATION=test/application.yml go test -tags test ./adapter/i18n/s3
```

`-tags test` is required wherever fixtures (`//go:build test`) are imported.

## Cross-Service Position

- Consumed by: `portal/application-shell` (every section that renders localized text).
- Depends on: vauthenticator JWKS (JWT validation + language claim), S3 (bundle storage).
- No DB. No DynamoDB. Stateless beyond the in-process cache.
