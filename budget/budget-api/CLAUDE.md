# CLAUDE.md — budget-api

---

## Tech Stack

| Concern            | Choice                                                                                                                                                                                                                                                        |
|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Language           | Go 1.25.1                                                                                                                                                                                                                                                     |
| Web framework      | Gin (`github.com/gin-gonic/gin v1.11.0`)                                                                                                                                                                                                                      |
| Persistence        | AWS DynamoDB via `aws-sdk-go-v2` (expense, revenue, attachment metadata)                                                                                                                                                                                      |
| Object storage     | AWS S3 via `aws-sdk-go-v2` — attachment file content                                                                                                                                                                                                          |
| Scheduling         | gocron v2 (`github.com/go-co-op/gocron/v2`) — the in-process Scheduled Expense job                                                                                                                                                              |
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
  budget/scheduledexpense/ # ScheduledExpense model, Create/Find(s)/Update/
                     # DeleteScheduledExpense, ScheduledExpenseActionsFacade,
                     # ScheduledExpenseRepository port (additive — Save+FindAll
                     # in #51, FindFor in #52, Delete in #53, UpdateStatus in #54);
                     # ScheduledExpenseJob (#55) — the port's FindAllActive
                     # (cross-user) + AdvanceLastEvaluatedDate are engine-only
  tags/              # SearchTagRepository port + SearchTag value object
  money/, time/      # value objects (Money, Date, Month, Year)
adapter/
  budget/expense/dynamodb/      # DynamoDbBudgetExpenseRepository + DynamoDbBudgetExpenseIdProvider
  budget/revenue/dynamodb/      # DynamoDbRevenueRepository + DynamoDbRevenueIdProvider
  budget/attachment/            # AwsCompositeAttachmentRepository — orchestrates dynamo + s3
  budget/attachment/dynamodb/   # DynamoDbAttachmentMetadataRepository + DynamoDbAttachmentIdProvider
  budget/attachment/s3/         # S3AttachmentContentRepository (file bytes)
  budget/scheduledexpense/dynamodb/ # DynamoDbScheduledExpenseRepository + DynamoDbScheduledExpenseIdProvider
  budget/scheduledexpense/scheduler/ # ScheduledExpenseJobConfigurer — runs the scheduled expense job (WebServerConfigurer)
  tags/rest/                    # RestSearchTagRepository + RistrettoCachedSearchTagRepository decorator
web/
  budget/expense/    # package expense    — RegisterExpenseEndpoints, representations, converters
  budget/revenue/    # package revenue    — RegisterRevenueEndpoints, representations, converters
  budget/attachment/ # package attachment — RegisterAttachmentEndpoints, representation, converter
  budget/scheduledexpense/ # package scheduledexpense — RegisterScheduledExpenseEndpoints, representation, converter
config/              # composition root — NewBudgetExpenseActionsFacade,
                     # NewRevenueActionsFacade, NewAttachmentActionsFacade,
                     # NewScheduledExpenseActionsFacade,
                     # NewScheduledExpenseJobConfigurer
main.go              # wires everything via WebServerProvisioner (shared framework);
                     # registers the scheduled expense job with RegisterConfigurer
```

Keep changes inside the right layer. Domain must not import adapter or web packages.

`web/budget/expense` imports domain expense as `domainexpense "...domain/budget/expense"` to avoid the package-name
clash with the web package itself. Same pattern in `web/budget/revenue` with `domainrevenue`, and in
`web/budget/scheduledexpense` with `domainscheduledexpense`.

### DynamoDB tables and config keys

| Table                          | Config key                                              |
|--------------------------------|---------------------------------------------------------|
| `BUDGET_EXPENSES`              | `budget-api.dynamo-db.budget-expense.table-name`        |
| `BUDGET_REVENUE`               | `budget-api.dynamo-db.revenue.table-name`               |
| `BUDGET_ATTACHMENT_METADATA`   | `budget-api.dynamo-db.attachment-metadata.table-name`   |
| `BUDGET_SCHEDULED_EXPENSE`     | `budget-api.dynamo-db.scheduled-expense.table-name`     |

Scheduled Expense job interval: `budget-api.scheduled-expense.job.interval` (Go duration, default `1h`).

S3 bucket holding attachment file bytes:

| Bucket                | Config key                                |
|-----------------------|-------------------------------------------|
| `<attachment bucket>` | `budget-api.s3.attachment.bucket-name`    |

AWS region is hardcoded to `eu-central-1` in `config/configurations.go` (`NewBudgetExpenseRepository`,
`NewRevenueRepository`, `NewAttachmentRepository`, `NewScheduledExpenseRepository`). LocalStack tests override the
endpoint in their fixture, not in the constructor.

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

**Scheduled Expense** (`adapter/budget/scheduledexpense/dynamodb/dynamo_db_scheduled_expense_id_provider.go`):

- PK: `user_name` — stored verbatim (not derived/encoded, unlike expense/revenue)
- SK: `id` — a plain UUID
- `FindAll`'s per-user scoping is therefore structural (the partition key itself), not a filter expression.
- The daily scheduled expense job (#55) does a plain table Scan across all users' definitions — no derived-key
  bookkeeping needed, at the accepted cost of a full scan (fine at this feature's expected low volume; see
  `docs/adr/0005-scheduled-expense-recurrence-and-generation-engine.md`).

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
`search_tags_user_<userName>_scope_<scope>`. budget-api now looks up **two** scopes — `expense` (expense tagging) and
`revenue` (revenue tagging) — each with its own cache instance and key suffix, so they never collide. There is no
invalidation on writes from `tag-api` — entries live until Ristretto evicts them or the process restarts. If you need
fresh tag data immediately after a tag mutation, depend on `rest.NewRestSearchTagRepository` directly instead of the
cached wiring constructors.

Both expense and revenue store tag **keys** only and resolve each key to its current value from tag-api on read, so
`RestSearchTagRepository` fetches `GET /api/tags/scope/<scope>` rather than the unscoped `GET /api/tags`. Each scope
literal is defined once at the wiring layer via a named constructor — `config.NewExpenseSearchTagRepository`
(`expense`) and `config.NewRevenueSearchTagRepository` (`revenue`), both over a private `newSearchTagRepository(scope)`
helper. The literal is passed into both the REST repository (drives the request path) and the Ristretto decorator
(drives the cache key). Scope stays an adapter/wiring concern — `domain/tags.SearchTagRepository` and the `SearchTag`
value object are unchanged and carry no `Scope`. See
`docs/adr/0001-expense-scoped-tag-lookup-hardcoded-at-wiring.md` and
`docs/adr/0002-revenue-tagging-mirrors-expense-without-events-or-totals.md`.

### Scheduled Expense job (`ScheduledExpenseJob`) — single replica only

`scheduledexpense.ScheduledExpenseJob` turns every `ACTIVE` Scheduled Expense into real `BudgetExpense`s as they
come due. Full rationale in `docs/adr/0005-scheduled-expense-recurrence-and-generation-engine.md`; the load-bearing facts:

- **Wiring:** `config.NewScheduledExpenseJobConfigurer(expenseFacade.CreateBudgetExpenseAction)` →
  `adapter/budget/scheduledexpense/scheduler.ScheduledExpenseJobConfigurer`, registered in `main.go` with the framework's
  `RegisterConfigurer`. It runs **once at startup, then every `budget-api.scheduled-expense.job.interval`**
  (default `1h`), singleton mode; `Dispose` cancels an in-flight run. A restart is the manual trigger.
- **Reuse the existing `CreateBudgetExpense` action.** The job creates expenses through the action held by the facade
  `main.go` already built (`NewBudgetExpenseActionsFacade()` returns the concrete facade for this). Calling
  `NewBudgetExpenseActionsFacade()` again would start a second reclassification listener and Kafka client.
- **Each day is evaluated once**, driven by `LastEvaluatedDate`: from `LastEvaluatedDate + 1` (today, when never
  evaluated) to today (`date.Today()`, UTC). A day is due when `day == min(Day, daysInMonth)`, in `Month` when set, and
  not past `EndDate`. Generate first, then advance — a crash in between yields a visible duplicate, never a loss.
- **No user token.** The job runs under an owner-only context (`security.User{UserName}`), so it must never reach
  tag-api: `FindAllActive` (cross-user paginated `Scan`) reads tag keys plus the display names stored at save time
  (`tag_names` attribute) instead of resolving them. User-facing reads still resolve live names and ignore `tag_names`.
- **Generated `Note`:** the definition's `Notes`, a newline, then
  `Expense generated by the scheduled expense "<Description>" with id: <id>`.
- **Single replica only** — there is no distributed lock; a second replica double-generates every due day.

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
| `POST`   | `/api/budget/expense`     | Create                  | `BudgetExpenseRepresentation`        | `201 No Content`, `500` on failure |
| `PUT`    | `/api/budget/expense/:id` | Update                  | `BudgetExpenseRepresentation`        | `204 No Content`, `500` on failure |
| `DELETE` | `/api/budget/expense/:id` | Delete                  | —                                    | `204 No Content`, `500` on failure |

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
| `POST`   | `/api/budget/revenue`             | Create         | `RevenueRepresentation` | `201 No Content`, `500` on failure |
| `PUT`    | `/api/budget/revenue/:id`         | Update         | `RevenueRepresentation` | `204 No Content`, `500` on failure |
| `DELETE` | `/api/budget/revenue/:id`         | Delete         | —                       | `204 No Content`, `500` on failure |

The `?q=year=YYYY` query param format preserves the Python revenue-api wire format consumed by the frontend.

**`RevenueRepresentation`** (create/update body):

```json
{
  "date": "DD/MM/YYYY",
  "amount": "100.00",
  "note": "string"
}
```

### Scheduled Expense — `web/budget/scheduledexpense/endpoint.go`

Create, list, get-by-id, update, pause/resume and delete (#51-#54) — see the parent issue and
`docs/adr/0005-scheduled-expense-recurrence-and-generation-engine.md`.

| Method | Path                                | Purpose                   | Request body                    | Response                                   |
|--------|--------------------------------------|----------------------------|----------------------------------|---------------------------------------------|
| `GET`  | `/api/budget/scheduled-expense`      | List (current user only)  | —                                | `ScheduledExpenseListRepresentation` `200` |
| `GET`  | `/api/budget/scheduled-expense/:id`  | Get one (current user only) | —                              | `ScheduledExpenseRepresentation` `200`, `404` if not found/not owned |
| `POST` | `/api/budget/scheduled-expense`      | Create                    | `ScheduledExpenseRepresentation` | `201 No Content`                           |
| `PUT`  | `/api/budget/scheduled-expense/:id`  | Update                    | `ScheduledExpenseRepresentation` | `204 No Content`, `404` if not found/not owned |
| `PATCH` | `/api/budget/scheduled-expense/:id` | Pause/resume              | `{"status":"ACTIVE"\|"PAUSED"}`   | `204 No Content`, `400` bad/missing status, `404` if not found/not owned |
| `DELETE` | `/api/budget/scheduled-expense/:id` | Delete (definition only)  | —                                | `204 No Content`, `404` if not found/not owned |

Update preserves `Status` and the scheduled expense job's internal `LastEvaluatedDate` from the existing record — neither
travels on the wire representation, and the DynamoDB adapter's `Save` replaces the whole item (see
`domain/budget/scheduledexpense/actions.go`'s `UpdateScheduledExpense.Execute`).

`PATCH /:id` is budget-api's first `PATCH` — a partial update of `Status` only, kept separate from `PUT /:id` (which
never carries `Status`). `PAUSED` → `PauseScheduledExpense`, `ACTIVE` → `ResumeScheduledExpense`; both stamp
`LastEvaluatedDate` with today (`date.Today()`, UTC) on an actual transition, and requesting the state a definition is
already in is a no-op `204` with no write (see ADR 0005). The adapter's `UpdateStatus` is a targeted `UpdateItem` on
`status` + `last_evaluated_date` only (not a full-item `Save`), so stored tag keys and concurrent edits are untouched.

Delete hard-deletes only the definition row; `BudgetExpense`s it already generated are left untouched (they carry no
reference back to it — see ADR 0005's "Delete" section).

Ownership on `GET /:id`, `PUT /:id`, `PATCH /:id` and `DELETE /:id` is enforced structurally, not by an explicit `UserName`
comparison: `FindFor` only ever looks inside the current user's own DynamoDB partition (`PK = user_name` from ctx), so
an `:id` belonging to another user is simply not found — `GET` returns `404`; the Update, Delete, Pause
and Resume actions return `ErrScheduledExpenseNotFound`, which the endpoints map to `404`. The adapter's `DeleteItem`
and `UpdateStatus` add an `attribute_exists(id)` condition as a backstop, so a row removed between the action's
`FindFor` and the write also surfaces as `ErrScheduledExpenseNotFound` rather than DynamoDB's silent no-op success. See the Scheduled Expense DynamoDB key scheme note above.

**`ScheduledExpenseRepresentation`** (create/update body; `id`/`status` are server-set and ignored on write — tag
shape is `{tagKey, tagValue}`, not `{key, value}`):

```json
{
  "description": "Rent",
  "amount": "1200.00",
  "notes": "string",
  "tags": [{"tagKey": "housing", "tagValue": "Housing"}],
  "day": 5,
  "month": 3,
  "endDate": "DD/MM/YYYY"
}
```

`month` and `endDate` are omitted (not `null`) when unset — a definition with no `month` recurs monthly, one with a
`month` recurs yearly on that day/month, and a definition with no `endDate` recurs forever. `status` on read is
`"ACTIVE"` or `"PAUSED"`.

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

Starting the service also starts the Scheduled Expense job, which runs once immediately — so restarting
is how to trigger it by hand (e.g. create a scheduled expense due today, restart, then check the expense list).

> `test/` is a local dev helper only — it is not a Go test package.

### Build

```bash
go build -o app .                          # local binary
CGO_ENABLED=0 GOOS=linux go build -o app . # cross-compile for Docker/Alpine
docker build -t mrflick72/budget/budget-api:1 -f ../../core-services/docker/ubuntu.Dockerfile .
```
