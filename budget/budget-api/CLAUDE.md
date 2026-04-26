# CLAUDE.md — budget-api

---

## Tech Stack

| Concern            | Choice                                                                                                                                                                                                                                                        |
|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Language           | Go 1.25.1                                                                                                                                                                                                                                                     |
| Web framework      | Gin (`github.com/gin-gonic/gin v1.11.0`)                                                                                                                                                                                                                      |
| Persistence        | AWS DynamoDB via `aws-sdk-go-v2`                                                                                                                                                                                                                              |
| In-process cache   | Ristretto (`github.com/dgraph-io/ristretto v0.2.0`)                                                                                                                                                                                                           |
| Decimal arithmetic | shopspring/decimal                                                                                                                                                                                                                                            |
| ID generation      | google/uuid (salt for DynamoDB range keys)                                                                                                                                                                                                                    |
| Auth               | JWT validation via the shared `core-services/golang-web-framework` middleware; JWKS fetched from `http://local.api.vauthenticator.com:9090/oauth2/jwks`                                                                                                       |
| Shared framework   | `github.com/mrflick72/onlyone-portal/core-services/golang-web-framework` — resolved via local `replace` directive in `go.mod` pointing to `../../core-services/golang-web-framework`                                                                          |
| Test assertions    | testify + go-playground/assert                                                                                                                                                                                                                                |
| Build tag          | `-tags test` is required for any test that imports shared fixtures. Helpers like `domain/tags/fixture.go` are guarded by `//go:build test` so they are never linked into the production binary. Running `go test ./...` without the tag will fail to compile. |

Config is read by the shared framework's `config.GetConfigurationManagerInstance()` (backed by Viper). The config file
path is set via the `CONFIG_FILE_LOCATION` env var.

---

## Architecture Decisions

### Hexagonal layout

```
domain/
  budget/expense/   # BudgetExpense model, CreateBudgetExpense, UpdateBudgetExpense,
                    # FindSpentBudget, DeleteBudgetExpense, BudgetExpenseActionsFacade,
                    # BudgetExpenseRepository port
  budget/revenue/   # Revenue model, CreateRevenue, UpdateRevenue, FindRevenue,
                    # DeleteRevenue, RevenueActionsFacade, RevenueRepository port
  tags/             # SearchTagRepository port + SearchTag value object
  money/, time/     # value objects (Money, Date, Month, Year)
adapter/
  budget/expense/dynamodb/  # DynamoDbBudgetExpenseRepository + DynamoDbBudgetExpenseIdProvider
  budget/revenue/dynamodb/  # DynamoDbRevenueRepository + DynamoDbRevenueIdProvider
  tags/rest/                # RestSearchTagRepository + RistrettoCachedSearchTagRepository decorator
web/
  budget/expense/   # package expense — RegisterExpenseEndpoints, representations, converters
  budget/revenue/   # package revenue — RegisterRevenueEndpoints, representations, converters
config/             # composition root — NewBudgetExpenseActionsFacade, NewRevenueActionsFacade
main.go             # wires everything via WebServerProvisioner (shared framework)
```

Keep changes inside the right layer. Domain must not import adapter or web packages.

`web/budget/expense` imports domain expense as `domainexpense "...domain/budget/expense"` to avoid the package-name
clash with the web package itself. Same pattern in `web/budget/revenue` with `domainrevenue`.

### DynamoDB tables and config keys

| Table             | Config key                                       |
|-------------------|--------------------------------------------------|
| `BUDGET_EXPENSES` | `budget-api.dynamo-db.budget-expense.table-name` |
| `BUDGET_REVENUE`  | `budget-api.dynamo-db.revenue.table-name`        |

AWS region is hardcoded to `eu-central-1` in `config/configurations.go` (`NewBudgetExpenseRepository`,
`NewRevenueRepository`). LocalStack tests override the endpoint in their fixture, not in the constructor.

### DynamoDB key schemes — do not change without a data migration

**Expense** (`adapter/budget/expense/dynamodb/dynamo_db_budget_expense_id_provider.go`):

- PK: `base64("<year>_<month>_<userName>")` — one partition per (user, calendar-month)
- RK: `base64("<day>_<uuid-salt>")`
- Full id stored in `budget_id` DynamoDB attribute as `<pk>-<rk>`
- Consequence: `DynamoDbBudgetExpenseRepository.FindByDateRange` issues **one DynamoDB Query per calendar month** in the
  range. A multi-month search is N queries, not a single scan.

**Revenue** (`adapter/budget/revenue/dynamodb/dynamo_db_revenue_id_provider.go`):

- PK: `base64("<year>_<userName>")` — one partition per (user, year)
- RK: `base64("<month>_<day>_<uuid-salt>")`
- Full id stored in `budget_id` as `<pk>-<rk>`
- Revenue range queries are therefore one Query per year.
- This layout preserves the Python `revenue-api` composite key format so existing `BUDGET_REVENUE` records remain
  readable without migration.

### Ownership enforced in two layers — keep both

For expense update/delete (`domain/budget/expense/actions.go`):

1. **Domain action**: `UpdateBudgetExpense` / `DeleteBudgetExpense` call `FindFor` first and compare
   `existingBudgetExpense.UserName == currentUser.UserName`. Returns a clean "not authorized" error.
2. **DynamoDB repository**: `Save` (when `!isNew`) and `Delete` attach `ConditionExpression: user_name = :user_name`.
   Race-safe backstop if the domain check were ever bypassed.

Revenue follows the same pattern. Removing either layer silently widens the authorization boundary.

### Tag cache has no invalidation hook

`RistrettoCachedSearchTagRepository` caches `GetAllTags` per user under key `search_tags_user_<userName>`. There is no
invalidation on writes from `tag-api` — entries live until Ristretto evicts them or the process restarts. If you need
fresh tag data immediately after a tag mutation, depend on `rest.NewRestSearchTagRepository` directly instead of
`config.NewSearchTagRepository`.

### Composition root quirks

`config.NewBudgetExpenseActionsFacade()` constructs **two** independent `SearchTagRepository` instances — one injected
into `DynamoDbBudgetExpenseRepository` (for read-after-write tag resolution) and one injected into `FindSpentBudget`.
They do not share a Ristretto cache; each call to `NewRistrettoCachedSearchTagRepository` creates a fresh cache.
Consolidating them is a reasonable cleanup but would change caching behavior — do it deliberately.

---

## API Contract per Domain

### Expense — `web/budget/expense/endpoint.go`

| Method   | Path                      | Purpose                 | Request body                         | Response                          |
|----------|---------------------------|-------------------------|--------------------------------------|-----------------------------------|
| `PUT`    | `/api/budget/expense`     | **Search** (not update) | `BudgetSearchCriteriaRepresentation` | `SpentBudgetRepresentation` `200` |
| `POST`   | `/api/budget/expense`     | Create                  | `BudgetExpenseRepresentation`        | `201 No Content`                  |
| `PUT`    | `/api/budget/expense/:id` | Update                  | `BudgetExpenseRepresentation`        | `204 No Content`                  |
| `DELETE` | `/api/budget/expense/:id` | Delete                  | —                                    | `204 No Content`                  |

`PUT /api/budget/expense` is overloaded as search so the frontend can send criteria as a JSON body (plain `GET` won't
accept a body). Do not normalize this to `GET` or `POST` without coordinating with the `application-shell` budget
bundle — the wire contract is load-bearing.

**`BudgetSearchCriteriaRepresentation`** (search request):

```json
{
  "month": "01",
  "year": "2024",
  "searchTagList": [
    "tagKey1",
    "tagKey2"
  ]
}
```

**`BudgetExpenseRepresentation`** (create/update body):

```json
{
  "date": "DD/MM/YYYY",
  "amount": "100.00",
  "note": "string",
  "tagKey": "string",
  "tagValue": "string"
}
```

### Revenue — `web/budget/revenue/endpoint.go`

| Method   | Path                              | Purpose        | Request body            | Response                        |
|----------|-----------------------------------|----------------|-------------------------|---------------------------------|
| `GET`    | `/api/budget/revenue?q=year=YYYY` | Search by year | —                       | `[]RevenueRepresentation` `200` |
| `POST`   | `/api/budget/revenue`             | Create         | `RevenueRepresentation` | `201 No Content`                |
| `PUT`    | `/api/budget/revenue/:id`         | Update         | `RevenueRepresentation` | `204 No Content`                |
| `DELETE` | `/api/budget/revenue/:id`         | Delete         | —                       | `204 No Content`                |

The `?q=year=YYYY` query param format preserves the Python revenue-api wire format consumed by the frontend.

**`RevenueRepresentation`** (create/update body):

```json
{
  "date": "DD/MM/YYYY",
  "amount": "100.00",
  "note": "string"
}
```

---

## How to Test Locally

### Unit tests (no infrastructure needed)

```bash
go test -tags test ./domain/... ./web/...
```

### Full test suite (adapter tests require LocalStack)

Start LocalStack first:

```bash
cd test
docker compose up -d   # starts localstack/localstack:3.2 with DynamoDB on :4566
```

Export the required AWS env vars (LocalStack accepts any non-empty value):

```bash
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
```

Run all tests:

```bash
cd ..
go test -tags test ./...
```

### Run the service locally end-to-end

```bash
cd test
docker compose up -d   # LocalStack must be running
./start.sh             # sets CONFIG_FILE_LOCATION=application.yml and runs ../main.go
```

`test/application.yml` configures the service for local use:

- Server port: `3050`
- CORS allowed origin: `http://local.onlyone-portal.com:8070`
- JWKS endpoint: `http://local.api.vauthenticator.com:9090/oauth2/jwks`
- Required JWT role: `USER_ROLE`
- Tag API base URL: `http://local.tag-api.onlyone-portal.com:8000`

> `test/` is a local dev helper only — it is not a Go test package.

### Build

```bash
go build -o app .                          # local binary
CGO_ENABLED=0 GOOS=linux go build -o app . # cross-compile for Docker/Alpine
docker build -t mrflick72/budget/budget-api:1 -f ../../core-services/docker/ubuntu.Dockerfile .
```
