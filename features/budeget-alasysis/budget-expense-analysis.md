---
Title: Budget Expense Analytic
Projects:
    - budget/budget-api
    - budget/analytic-api
    - portal/application-shell
Scope:
    - Frontend in portal/application-shell
    - Backend in budget/budget-api and budget/analytic-api
Status: done (as-built — this document reconciles the spec to the shipped code)
---

# Overview

A dedicated analytics page lets a user visualise their budget expenses through two
charts plus a self-service data-rebuild action. The data is served by
`budget/analytic-api`, which keeps a **projection** of the user's expenses in its own
Postgres database. The projection is kept up to date by consuming domain events that
`budget/budget-api` (the data owner) publishes to Kafka whenever an expense is created,
updated or deleted. Read requests are answered purely from the projection —
`analytic-api` never calls `budget-api` at request time. The single exception is the
explicit **reindex** action, which pulls from `budget-api` over REST to rebuild the
caller's projection; it is off the read path.

**Roles.** There is one role: any authenticated portal user (a valid vauthenticator
JWT). The analytics endpoints and the reindex action are scoped to the caller's
`user_name` and perform no authority/role check — every authenticated user sees and
operates only on their own data. There is no admin or cross-user view.

# Architecture

```
budget-api  --(JSON CREATE/UPDATE/DELETE on topic "budget-api.expense")-->  Kafka
                                                                              |
analytic-api: confluent-kafka consumer (in-process, FastAPI lifespan)  <------+
        -> BudgetExpenseEventHandler -> ExpenseProjectionRepository (Postgres):
           save (upsert) by id / delete by id

analytic-api HTTP (consumed by the frontend):
   PUT  /api/analytic/budget/expense/total-by-tag    -> SQL GROUP BY over the projection
   PUT  /api/analytic/budget/expense/total-by-year   -> SQL GROUP BY over the projection
   POST /api/analytic/budget/expense/reindex         -> pull from budget-api REST, re-save

reindex only:
analytic-api --(PUT /api/budget/expense per (year,month), caller's JWT)--> budget-api
```

`budget-api` is the source of truth. `analytic-api` owns only the read-optimised
projection and is **eventually consistent** with `budget-api`: a just-created expense
appears in the charts once its Kafka event has been consumed, not synchronously.

# Cross-cutting invariants

These are enforced and must be preserved on both sides of the wire:

1. **Tag key-vs-value contract.** The frontend **filters by tag _key_**
   (`SearchTag.key`) but **labels/groups results by tag _value_** (`SearchTag.value`).
   The projection stores both. `total-by-tag` filters expenses by key (keeping all
   sibling tags of a matched expense) and aggregates/labels by value; `total-by-year`
   counts each expense once.
2. **Per-user scoping.** Every projection query and the reindex are scoped to the
   authenticated `user_name`. No endpoint can read or rebuild another user's data.
3. **Year-range cap = 20 (hard invariant).** `total-by-year` and `reindex` reject a
   range where `fromYear > toYear` or `toYear - fromYear + 1 > 20`. Enforced in the
   React page (`MAX_YEAR_RANGE_SIZE = 20`, inline validation) **and** in the FastAPI
   representations (`MAX_YEAR_RANGE_SIZE = 20`, pydantic `model_validator` → `422`).
   The two must stay in lockstep.
4. **Idempotent projection writes.** `CREATE`/`UPDATE` → upsert by id
   (`INSERT ... ON CONFLICT (id) DO UPDATE`, tag rows replaced to handle re-tagging);
   `DELETE` → delete by id (tags cascade). Re-applying any event is safe.
5. **At-least-once consumption.** Offsets are committed only after the DB write, so a
   crash mid-apply redelivers; invariant 4 makes redelivery harmless.

# Frontend (`portal/application-shell`) — DONE (as-built)

A dedicated multi-page-app entry (`analytics/index.tsx` → `AnalyticsApp` →
`AnalyticsDashboardPage`), reachable from the Home page tile and the global navigation
menu (`components/menu/AnalyticsPageMenuItem.tsx`, the `Analytics` nav item). Built on
`@mui/x-charts`. The page has **three** stacked panels:

1. **Total by tag** — budget expense for one year, optionally filtered by month and by
   one or more search tags. Bars are grouped/labelled by tag _value_; the filter sends
   tag _keys_. Re-queries on any filter change.
2. **Total by year** — budget expense across a range of years, optionally filtered by a
   single search tag (by key). Bars are grouped/labelled by year, every year in range
   present (zero-filled). Shows an inline warning and skips the request when the range
   is invalid or exceeds 20 years.
3. **Reindex data** — a self-service rebuild of the user's own projection over a
   from/to-year range. Spinner + disabled button while running; result reported via a
   `Snackbar` (success shows the imported count, failure shows a retry message). On
   success it bumps both charts' retry tokens so they reload against the freshly rebuilt
   projection. Same 20-year range validation as panel 2. This is the only mutating call
   on the page.

Each chart owns its own loading / error / empty state with a **Retry** affordance
(error path re-issues the query). Charts do not auto-refresh: new expenses surface on
the next filter change or after a reindex — the accepted freshness model is
**re-query-on-interaction + manual reindex** (no polling/SSE).

Message strings live in `messages/bundle/analytics/message_bundle_{en_en,it_it}.yaml`,
mapped in `OnlyonePortalPagesConfigMap.analytics()`, typed in `MessageBundles.ts`
(`AnalyticsPageMessageBundle`, including the `reindex` group).

# Backend

`budget/analytic-api` is a Python 3.12+ FastAPI service. `budget/budget-api` is the Go
data owner.

## budget-api — event publishing (DONE)

When a budget expense is created, updated or deleted successfully, an event is fired on
the Kafka topic `budget-api.expense` via the `BudgetExpenseEventPublisher` abstraction,
implemented by `KafkaBudgetExpenseEventPublisher`, wired into the create/update/delete
actions.

Wire contract (plain JSON, not schema-registry encoded):

```json
{
  "action": "CREATE | UPDATE | DELETE",
  "payload": {
    "id": "<uuid>",
    "userName": "<user>",
    "date": "dd/MM/yyyy",
    "amount": "<decimal string, scale 2>",
    "note": "<string>",
    "tags": [ { "key": "<tag key>", "value": "<tag value>" } ]
  }
}
```

Infra in `budget/docker-compose.yml`: Kafka, the `budget-api.expense` topic
(3 partitions), a Confluent schema-registry (present but **unused** — events are plain
JSON; left in place for the future), and `analytic-postgres`.

> Note the tag shape asymmetry: the **Kafka event** uses `{key, value}`, but
> `budget-api`'s **REST** expense response (used by reindex) serialises tags as
> `{tagKey, tagValue}`. Both are mapped to `ExpenseTag(key, value)` on the
> analytic-api side — see `RestBudgetExpenseSource._to_expenses`.

## analytic-api — projection & queries (DONE)

The legacy REST-pull read path (`RestExpenseLoader` / `ExpenseLoader` / the flat
`PUT /api/analytic/budget/expense` endpoint) has been **removed**. Both chart endpoints
now serve from Postgres; budget-api is touched only by reindex.

### Design decisions (as shipped)

1. **Read path is Postgres only.** `RestExpenseLoader`, the `ExpenseLoader` ABC and the
   flat endpoint are gone. Both chart endpoints aggregate the projection with SQL.
2. **Projection model: raw per-expense rows + SQL aggregation.** One row per expense in
   `budget_expense_projection`, its tags in `budget_expense_tag`; totals computed with
   `GROUP BY` at query time. Flexible for future charts.
3. **Consumer runtime: in-process via FastAPI lifespan, `confluent-kafka`.** Single
   deployable. The poll loop runs on a daemon thread (the client is synchronous) and is
   skipped under `KAFKA_CONSUMER_ENABLED=false` so `pytest` needs no broker.
4. **Tag aggregation: full amount to each tag.** A €100 expense tagged `[food, leisure]`
   contributes €100 to `food` and €100 to `leisure` — falls out of the tag join
   naturally.

### HTTP contract

| Method | Path | Body | Response |
|---|---|---|---|
| PUT | `/api/analytic/budget/expense/total-by-tag` | `{ year, month?, tags? }` (`tags` = tag keys; `month` 1–12) | `[{ tag, total }]` (`tag` = tag value) |
| PUT | `/api/analytic/budget/expense/total-by-year` | `{ fromYear, toYear, tag? }` (`tag` = tag key) | `[{ year, total }]`, every year in range present (zero-filled) |
| POST | `/api/analytic/budget/expense/reindex` | `{ fromYear, toYear }` | `{ imported }` — count rebuilt into the caller's projection |

All queries are scoped to the authenticated `user_name` from the JWT. `total-by-year`
and `reindex` enforce the 20-year cap (invariant 3) → `422` on violation.

### Components (as-built)

1. **Postgres projection store** (`src/app/analytic/adapter/db/repository.py`)
   - Tables: `budget_expense_projection(id PK, user_name, expense_date DATE,
     amount NUMERIC(.,2), note)` and `budget_expense_tag(expense_id FK, tag_key,
     tag_value)`. Schema in `scripts/init.sql`, applied **out-of-band** (docker-compose
     initdb locally, the testcontainer in integration tests, deployment migrations) —
     the adapter never runs DDL.
   - `PostgresExpenseProjectionRepository`: `save` (upsert + replace tag rows,
     de-duplicating tags by key), `delete(id)`, `total_by_tag(user, year, month,
     tag_keys)` (filter by key via `EXISTS`, `GROUP BY tag_value`),
     `total_by_year(user, from_year, to_year, tag_key)` (year range, zero-fill missing
     years in Python). Owns a lazily-opened `psycopg` connection pool.
2. **Kafka consumer** (`src/app/analytic/adapter/kafka/consumer.py`)
   - `BudgetExpenseEventHandler`: decode `{action, payload}`, dispatch
     CREATE/UPDATE→`save`, DELETE→`delete`. Kafka-type-free for unit testing. Raises
     `MalformedEventError` on bad structure / unknown action / DELETE-without-id /
     unparseable payload.
   - `BudgetExpenseConsumer`: daemon-thread poll loop. **Poison-message handling** —
     invalid JSON or `MalformedEventError` is logged and **committed/skipped** so it
     cannot wedge the partition; a transient failure (e.g. DB down) is left
     **uncommitted** for redelivery. Manual commit after a successful apply.
   - Started/stopped in the `server.py` FastAPI `lifespan` (which also opens/closes the
     DB pool) when `KAFKA_CONSUMER_ENABLED` is true.
3. **Service & DI** (`container.py`, `domain/service.py`)
   - `BudgetExpenseAnalysisService` is backed by the Postgres repository and reads
     `user_name` from the injected `SecurityContextResolver` to scope every query.
4. **Reindex** (`domain/reindex.py`, `adapter/rest/source.py`) — see below.

### Reindex / reimport (recovery)

`POST /api/analytic/budget/expense/reindex` `{fromYear, toYear}` rebuilds the **calling
user's** projection. Available to every user as a first-class self-service action.

- **Per-user, self-service.** Scoped to the JWT user. An all-users reindex would need
  new cross-user query support in budget-api (separate task) and is out of scope.
- **Pulls from budget-api over REST.** `RestBudgetExpenseSource` queries
  `PUT /api/budget/expense` per (year, month) across the range, fanning out over a
  `ThreadPoolExecutor` (`MAX_CONCURRENT_BUDGET_API_CALLS = 6`). The caller's token and
  user name are captured once on the request thread (the security context is a
  thread-local) before fan-out. The REST response carries full tag `key`+`value`
  (as `{tagKey, tagValue}`), so the projection is fully repopulated. This is the only
  place analytic-api calls budget-api and it is off the read path.
- **All-or-nothing fetch.** `fetch_expenses` materialises every (year, month) batch
  before saving; if any budget-api call fails (`raise_for_status`), the whole reindex
  aborts with an error and **no rows are written** — the frontend shows the failure
  Snackbar. On success, `ExpenseReindexService` `save`s each expense and returns the
  count as `imported`.
- **Upsert-only.** Fills gaps from missed `CREATE`/`UPDATE` events; it does **not**
  remove rows whose `DELETE` was missed (those linger until re-touched).
- Requires the `BUDGET_API_BASE_URL` env var (reindex-only).

### Tests (as-built)

- Repository integration tests (upsert idempotency, delete, re-tagging, multi-tag
  full-amount counting, year-range zero-fill, user scoping) behind the `integration`
  marker (testcontainers Postgres, needs Docker).
- Consumer unit tests (each action → repository calls; poison-message skip vs retry).
- `test_aggregation_end_point.py` runs against the fast suite (no broker/DB) with
  `KAFKA_CONSUMER_ENABLED=false`. `mypy` clean.

### Docs (DONE)

`analytic-api/AGENTS.md` documents the Kafka+Postgres projection, env vars, consumer
lifecycle, reindex, and the tag key-vs-value note. Root `CLAUDE.md`, `budget/CLAUDE.md`
and `portal/application-shell/CLAUDE.md` describe the service and its `→ Kafka, Postgres`
dependency and the frontend reindex panel.

# Known issues / follow-ups (not part of this reconciliation)

- **Message-bundle corruption.** `messages/bundle/analytics/message_bundle_en_en.yaml`
  line 25 reads `reinde  reindex:` (a corrupted key) — the `reindex.*` strings the page
  reads via `messages.reindex.*` will not resolve from the `en_en` bundle. Needs a
  one-line fix to `reindex:`; verify the `it_it` bundle too. Tracked as a bug, outside
  this spec-reconciliation task.
- **All-users reindex** is unsupported (needs budget-api cross-user query).
- **Missed-DELETE rows** linger after reindex until re-touched (upsert-only by design).

# User Stories

- **US-1**: As an authenticated user, I see an **Analytics** tile on the Home page and
  an **Analytics** item in the global navigation menu, so that I can reach the analytics
  dashboard from anywhere in the portal.
- **US-2**: As an authenticated user, I open the analytics dashboard and land on a page
  with three panels — Total by tag, Total by year, and Reindex data — so that I can both
  explore my spending and rebuild my data from one place.
- **US-3**: As an authenticated user, I view a **Total by tag** bar chart for a chosen
  year, optionally narrowed by month and by one or more tags, so that I can see how my
  spending breaks down across tags; bars are labelled by tag value while my filter
  selects tag keys.
- **US-4**: As an authenticated user, I view a **Total by year** bar chart across a
  from/to-year range, optionally narrowed by a single tag, so that I can see my spending
  trend year over year; every year in the range is shown, zero-filled when there is no
  spend.
- **US-5**: As an authenticated user, when I change any filter (year, month, tags,
  year range, tag), the corresponding chart re-queries and updates, so that I always see
  data matching my current selection.
- **US-6**: As an authenticated user, when a year range is inverted or wider than 20
  years, I see an inline warning and the chart does not query, and the reindex/year
  endpoints reject the same range with a `422`, so that I cannot request an unbounded
  span and client and server agree.
- **US-7**: As an authenticated user, when a chart fails to load, I see an error state
  with a **Retry** button that re-issues the request, so that a transient failure does
  not strand me.
- **US-8**: As an authenticated user, when my filter matches no data, I see an explicit
  empty-state message rather than a blank chart, so that I can tell "no spend" apart
  from "still loading."
- **US-9**: As an authenticated user, I rebuild my own analytics data over a from/to-year
  range via the **Reindex data** panel, so that I can recover after missed events or a
  reimport; the button shows a spinner and is disabled while running.
- **US-10**: As an authenticated user, when a reindex succeeds, I see a Snackbar with the
  imported count and both charts reload against the rebuilt projection, so that I
  immediately see the corrected data without leaving the page.
- **US-11**: As an authenticated user, when a reindex fails (e.g. budget-api is
  unreachable), I see a failure Snackbar and my existing projection is unchanged (the
  rebuild is all-or-nothing), so that a partial/corrupt rebuild can never occur.
- **US-12**: As an authenticated user, my analytics, and my reindex, only ever read or
  rebuild my own expenses (scoped to my JWT `user_name`), so that I can never see or
  affect another user's data.
- **US-13**: As the analytic-api service, I keep the Postgres projection current by
  consuming `budget-api.expense` events — upserting on CREATE/UPDATE and deleting on
  DELETE — so that the charts reflect budget changes without calling budget-api on the
  read path.
- **US-14**: As the analytic-api service, I commit a Kafka offset only after the DB write
  (at-least-once) and treat upsert/delete as idempotent, so that redelivery after a crash
  is harmless.
- **US-15**: As the analytic-api service, I skip and commit a permanently-unprocessable
  event (bad JSON, unknown action, malformed payload) while leaving a transient failure
  uncommitted for retry, so that a poison message cannot wedge a partition but a
  recoverable error is not silently dropped.
- **US-16**: As a budget user, when I create, update, or delete an expense in budget-api,
  an event is published to `budget-api.expense`, so that my analytics become consistent
  shortly afterwards (eventual consistency); if they look stale I can force consistency
  with Reindex.
