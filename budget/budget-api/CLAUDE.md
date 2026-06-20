# CLAUDE.md — budget-api

---

## Tech Stack

| Concern            | Choice                                                                                                                                                                                                                                                        |
|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Language           | Go 1.25.1                                                                                                                                                                                                                                                     |
| Web framework      | Gin (`github.com/gin-gonic/gin v1.11.0`)                                                                                                                                                                                                                      |
| Persistence        | AWS DynamoDB via `aws-sdk-go-v2` (expense, revenue, attachment metadata)                                                                                                                                                                                      |
| Object storage     | AWS S3 via `aws-sdk-go-v2` — attachment file content                                                                                                                                                                                                          |
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
  budget/expense/    # BudgetExpense model, CreateBudgetExpense, UpdateBudgetExpense,
                     # FindSpentBudget, DeleteBudgetExpense, BudgetExpenseActionsFacade,
                     # BudgetExpenseRepository port
  budget/revenue/    # Revenue model, CreateRevenue, UpdateRevenue, FindRevenue,
                     # DeleteRevenue, RevenueActionsFacade, RevenueRepository port
  budget/attachment/ # Attachment + AttachmentMetadata models, SaveAttachment,
                     # GetAttachment, DeleteAttachment, AttachmentActionsFacade,
                     # AttachmentRepository port
  tags/              # SearchTagRepository port + SearchTag value object
  money/, time/      # value objects (Money, Date, Month, Year)
adapter/
  budget/expense/dynamodb/      # DynamoDbBudgetExpenseRepository + DynamoDbBudgetExpenseIdProvider
  budget/revenue/dynamodb/      # DynamoDbRevenueRepository + DynamoDbRevenueIdProvider
  budget/attachment/            # AwsCompositeAttachmentRepository — orchestrates dynamo + s3
  budget/attachment/dynamodb/   # DynamoDbAttachmentMetadataRepository + DynamoDbAttachmentIdProvider
  budget/attachment/s3/         # S3AttachmentContentRepository (file bytes)
  tags/rest/                    # RestSearchTagRepository + RistrettoCachedSearchTagRepository decorator
web/
  budget/expense/    # package expense    — RegisterExpenseEndpoints, representations, converters
  budget/revenue/    # package revenue    — RegisterRevenueEndpoints, representations, converters
  budget/attachment/ # package attachment — RegisterAttachmentEndpoints, representation, converter
config/              # composition root — NewBudgetExpenseActionsFacade,
                     # NewRevenueActionsFacade, NewAttachmentActionsFacade
main.go              # wires everything via WebServerProvisioner (shared framework)
```

Keep changes inside the right layer. Domain must not import adapter or web packages.

`web/budget/expense` imports domain expense as `domainexpense "...domain/budget/expense"` to avoid the package-name
clash with the web package itself. Same pattern in `web/budget/revenue` with `domainrevenue`.

### DynamoDB tables and config keys

| Table                          | Config key                                              |
|--------------------------------|---------------------------------------------------------|
| `BUDGET_EXPENSES`              | `budget-api.dynamo-db.budget-expense.table-name`        |
| `BUDGET_REVENUE`               | `budget-api.dynamo-db.revenue.table-name`               |
| `BUDGET_ATTACHMENT_METADATA`   | `budget-api.dynamo-db.attachment-metadata.table-name`   |

S3 bucket holding attachment file bytes:

| Bucket                | Config key                                |
|-----------------------|-------------------------------------------|
| `<attachment bucket>` | `budget-api.s3.attachment.bucket-name`    |

AWS region is hardcoded to `eu-central-1` in `config/configurations.go` (`NewBudgetExpenseRepository`,
`NewRevenueRepository`, `NewAttachmentRepository`). LocalStack tests override the endpoint in their fixture, not in the
constructor.

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

**Attachment metadata** (`adapter/budget/attachment/dynamodb/dynamo_db_attachment_id_provider.go`):

- PK: `<budgetId>_<UPPERCASE budgetType>` — one partition per parent expense/revenue (e.g. `budget-123_EXPENSE`)
- RK: `attachment_id` — UUID generated server-side on first save (re-used on update)
- GSI `<TableName>_GLOBAL_INDEX` keyed on `attachment_id` (HASH) with `ProjectionType: ALL`. Lookups by attachment id
  (read content, delete) hit the GSI, then resolve the base item by `(pk, attachment_id)`. Ownership is enforced by a
  `user_name = :user_name` filter on the query — never trust the id alone.
- Each item carries a `metadata` map attribute holding free-form key/value strings plus two reserved keys used by the
  composite repository:
  - `metadata_bucket` — S3 bucket name where the content lives
  - `metadata_object_key` — S3 object key (`<YYYY>/<MM>/<DD>/<pk>/<attachmentId>`)

**Attachment content** (`adapter/budget/attachment/s3/s3_attachment_content_repository.go`):

- Object key: `<YYYY>/<MM>/<DD>/<budgetId>_<UPPERCASE budgetType>/<attachmentId>` — date prefix is the attachment's
  business date, not the upload time.
- `file_location` stored alongside metadata is `<bucket>/<objectKey>` for traceability.

### Ownership enforced in two layers — keep both

For expense update/delete (`domain/budget/expense/actions.go`):

1. **Domain action**: `UpdateBudgetExpense` / `DeleteBudgetExpense` call `FindFor` first and compare
   `existingBudgetExpense.UserName == currentUser.UserName`. Returns a clean "not authorized" error.
2. **DynamoDB repository**: `Save` (when `!isNew`) and `Delete` attach `ConditionExpression: user_name = :user_name`.
   Race-safe backstop if the domain check were ever bypassed.

Revenue follows the same pattern. Removing either layer silently widens the authorization boundary.

For attachments, ownership is enforced inside the composite adapter and the metadata repository:

1. **`AwsCompositeAttachmentRepository`** resolves the current user from the request context before any read or
   delete (`GenAttachment`, `FindAllAttachment`, `DeleteAttachment`). The username is the partition lens for every
   subsequent call.
2. **`DynamoDbAttachmentMetadataRepository`** layers a `user_name = :user_name` filter expression on top of the GSI
   query in `GetAttachment` and `Delete`. A request with another user's attachment id therefore returns
   `attachment not found`, never the actual item.
3. **Save** sets `attachment.Owner = *user.UserName` in `SaveAttachment.Execute` (domain layer) before the repository
   sees the entity, so a client cannot impersonate another owner via the upload form fields.

### Tag cache has no invalidation hook

`RistrettoCachedSearchTagRepository` caches `GetAllTags` per user under key
`search_tags_user_<userName>_scope_<scope>` (scope is `expense`, the only scope budget-api looks up). There is no
invalidation on writes from `tag-api` — entries live until Ristretto evicts them or the process restarts. If you need
fresh tag data immediately after a tag mutation, depend on `rest.NewRestSearchTagRepository` directly instead of
`config.NewSearchTagRepository`.

`budget-api` only tags expenses, so `RestSearchTagRepository` fetches `GET /api/tags/scope/expense` rather than the
unscoped `GET /api/tags`. The `"expense"` scope literal is defined once at the wiring layer
(`config.NewSearchTagRepository`) and passed into both the REST repository (drives the request path) and the Ristretto
decorator (drives the cache key). Scope stays an adapter/wiring concern — `domain/tags.SearchTagRepository` and the
`SearchTag` value object are unchanged and carry no `Scope`. See
`docs/adr/0001-expense-scoped-tag-lookup-hardcoded-at-wiring.md`.

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

### Attachment — `web/budget/attachment/endpoint.go`

Attachments are file uploads attached to either an expense or a revenue (the same endpoint serves both — the
`budgetType` form field selects the parent aggregate).

| Method   | Path                                                | Purpose                              | Request                                     | Response                                              |
|----------|-----------------------------------------------------|--------------------------------------|---------------------------------------------|-------------------------------------------------------|
| `POST`   | `/api/attachment`                                   | Create or replace                    | `multipart/form-data`                       | `204 No Content`                                      |
| `GET`    | `/api/attachment/metadata/:budgetType/:budgetId`    | List metadata for a parent           | —                                           | `[]AttachmentMetadataRepresentation` `200`            |
| `GET`    | `/api/attachment/:attachmentId/content`             | Download file bytes                  | —                                           | Raw bytes with `Content-Type` + `Content-Disposition` |
| `DELETE` | `/api/attachment/:attachmentId`                     | Delete (metadata + S3 content)       | —                                           | `204 No Content`                                      |

**Upload form fields** (`POST /api/attachment`):

| Field          | Required | Notes                                                                                  |
|----------------|----------|----------------------------------------------------------------------------------------|
| `file`         | yes      | Multipart file part — used for filename and content type                               |
| `budgetId`     | yes      | Parent expense or revenue id                                                           |
| `budgetType`   | yes      | `expense` or `revenue` (case-insensitive — folded to upper case for the partition key) |
| `date`         | yes      | `DD/MM/YYYY` — drives the S3 date-scoped path                                          |
| `attachmentId` | no       | Provide to overwrite an existing attachment; omitted for create (server-generated)     |

**Download** sets `Content-Disposition: attachment; filename="<original-name>"` and falls back to
`application/octet-stream` if no content type was stored.

**Delete** removes the metadata row first, then the S3 object. If the metadata row is missing or owned by another user
the call returns an error and S3 is left untouched. If the S3 deletion fails after the metadata row is gone, the orphan
content can be reaped by an out-of-band sweeper — the metadata row is the source of truth.

**`AttachmentMetadataRepresentation`** (list response item):

```json
{
  "attachmentId": "uuid",
  "fileName": "receipt.pdf",
  "owner": "user-name",
  "budgetId": "budget-123",
  "budgetType": "expense"
}
```

The S3 bucket / object key are intentionally **not** exposed on the wire — clients reach the bytes via
`/api/attachment/:attachmentId/content`.

---

## How to Test Locally

### Unit tests (no infrastructure needed)

```bash
go test -tags test ./domain/... ./web/...
```

### Full test suite (adapter tests require LocalStack)

Start LocalStack first — DynamoDB and S3 are both required since attachment adapter tests round-trip metadata and
content against LocalStack:

```bash
cd test
docker compose up -d   # starts localstack/localstack:3.2 with DynamoDB + S3 on :4566
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
