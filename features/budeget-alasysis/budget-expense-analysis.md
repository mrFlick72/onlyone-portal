---
Title: Budget Expense Analytic
Projects:
    - budget/budget-api
    - budget/analytic-api
    - portal/application-shell
Scope:
    - Frontend in portal/application-shell
    - Backend in budget/budget-api and budget/analytic-api
Status: in progress
---

# Overview

A dedicated analytics page lets a user visualise their budget expenses through two
charts. The data is served by `budget/analytic-api`, which keeps a **projection** of
the user's expenses in its own Postgres database. The projection is kept up to date
by consuming domain events that `budget/budget-api` (the data owner) publishes to
Kafka whenever an expense is created, updated or deleted. Read requests are answered
purely from the projection — `analytic-api` never calls `budget-api` at request time.

# Architecture

```
budget-api  --(JSON CREATE/UPDATE/DELETE on topic "budget-api.expense")-->  Kafka
                                                                              |
analytic-api: confluent-kafka consumer (in-process, FastAPI lifespan)  <------+
        -> ExpenseProjectionRepository (Postgres): upsert by id / delete by id

analytic-api HTTP (consumed by the frontend, contract unchanged):
   PUT /api/analytic/budget/expense/total-by-tag    -> SQL GROUP BY over the projection
   PUT /api/analytic/budget/expense/total-by-year   -> SQL GROUP BY over the projection
```

`budget-api` is the source of truth. `analytic-api` owns only the read-optimised
projection and is eventually consistent with `budget-api`.

# Frontend

Developed in `portal/application-shell` as a dedicated page, reachable from the Home
page via a tile and from the global navigation menu. The UI is composed of two charts:

- **Total by tag** — budget expense for one year, optionally filtered by month and by
  one or more search tags. Bars are grouped/labelled by tag.
- **Total by year** — budget expense across a range of years, optionally filtered by a
  single search tag. Bars are grouped/labelled by year.

Tag contract (important, must be preserved by the backend): the frontend **filters by
tag _key_** (`SearchTag.key`) but **labels/groups results by tag _value_**
(`SearchTag.value`). The backend projection therefore stores both and must filter by
key while aggregating and labelling by value.

**Status: DONE.** `AnalyticsDashboardPage`, `AnalyticBarChart`, `AnalyticRepository`
(both endpoints), Home tile, nav item, and en/it message bundles are all implemented.
No frontend changes are required by this feature.

# Backend

`budget/analytic-api` is a Python 3.12+ FastAPI service. `budget/budget-api` is the
Go data owner.

## budget-api — event publishing (DONE)

When a budget expense is created, updated or deleted successfully, an event is fired on
the Kafka topic `budget-api.expense` via the `BudgetExpenseEventPublisher` abstraction,
implemented by `KafkaBudgetExpenseEventPublisher`. The publisher is already wired into
the create/update/delete actions.

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

Infra already present in `budget/docker-compose.yml`: Kafka, the `budget-api.expense`
topic (3 partitions), and a Confluent schema-registry (currently **unused** — events are
plain JSON; left in place for the future).

## analytic-api — projection & queries (TO IMPLEMENT)

The current `analytic-api` still answers the two endpoints by pulling every month/year
from `budget-api` over REST and aggregating in memory (`RestExpenseLoader`) — the exact
cost this feature exists to remove. It has **no Kafka consumer and no Postgres**. The
work below replaces that path.

### Confirmed design decisions

1. **Replace the REST-pull path entirely.** Remove `RestExpenseLoader`, the
   `ExpenseLoader` ABC, the `BUDGET_API_BASE_URL` dependency, and the legacy flat
   `PUT /api/analytic/budget/expense` endpoint (the frontend does not call it). Both
   chart endpoints serve from Postgres.
2. **Projection model: raw per-expense rows + SQL aggregation.** One row per expense
   (upsert on CREATE/UPDATE, delete by id on DELETE); totals computed with SQL
   `GROUP BY` at query time. Flexible for future charts.
3. **Consumer runtime: in-process via FastAPI lifespan, `confluent-kafka`.** Single
   deployable. The poll loop runs in a background thread (the client is synchronous)
   and is skipped under the test flag so `pytest` needs no broker.
4. **Tag aggregation: full amount to each tag.** A €100 expense tagged `[food, leisure]`
   contributes €100 to `food` and €100 to `leisure` (unchanged behaviour — falls out of
   the tag join naturally).

### HTTP contract (unchanged — keep frontend & `test_aggregation_end_point.py` green)

| Method | Path | Body | Response |
|---|---|---|---|
| PUT | `/api/analytic/budget/expense/total-by-tag` | `{ year, month?, tags? }` (`tags` = tag keys) | `[{ tag, total }]` (`tag` = tag value) |
| PUT | `/api/analytic/budget/expense/total-by-year` | `{ fromYear, toYear, tag? }` (`tag` = tag key) | `[{ year, total }]`, every year in range present (zero-filled) |

All queries are scoped to the authenticated `user_name` from the JWT.

### Work breakdown

1. **Dependencies & infra**
   - `pyproject.toml`: add `confluent-kafka`, `psycopg[binary]`; dev: a Postgres test
     fixture (testcontainers or the docker-compose db behind a marker).
   - `budget/docker-compose.yml`: add `postgres:16-alpine` (db `analytic`).
   - `local/.env` + AGENTS.md env table: add `DATABASE_URL`,
     `KAFKA_BOOTSTRAP_SERVERS`, `KAFKA_CONSUMER_GROUP`, `KAFKA_EXPENSE_TOPIC`.
2. **Postgres projection store** (`src/app/analytic/adapter/db/`)
   - Tables: `budget_expense_projection(id PK, user_name, expense_date DATE,
     amount NUMERIC(.,2), note)` and `budget_expense_tag(expense_id FK, tag_key,
     tag_value)`. Schema lives in `scripts/init.sql` (out of the app source, like
     plan-api); applied via docker-compose initdb locally, by the test fixture, and
     by deployment migrations — the adapter runs no DDL.
   - `PostgresExpenseProjectionRepository`: `upsert` (`INSERT ... ON CONFLICT (id) DO
     UPDATE`, replacing tag rows to handle re-tagging), `delete(id)`,
     `total_by_tag(user_name, year, month, tag_keys)` (filter by key, `GROUP BY`
     tag_value), `total_by_year(user_name, from_year, to_year, tag_key)` (year range,
     zero-fill missing years).
3. **Kafka consumer** (`src/app/analytic/adapter/kafka/`)
   - `BudgetExpenseConsumer`: decode `{action, payload}`, dispatch CREATE/UPDATE→upsert,
     DELETE→delete. Manual commit after the DB write (at-least-once + idempotent ops).
   - Started/stopped in a FastAPI `lifespan` handler in `server.py`.
4. **Service & DI rewiring**
   - Back `BudgetExpenseAnalysisService` with the Postgres repository instead of
     `ExpenseLoader`, keeping `total_by_tag` / `total_by_year` signatures identical.
   - Update `container.py`; delete `RestExpenseLoader`, `ExpenseLoader`, the flat
     endpoint + its representations, and the obsolete REST-pull tests.
5. **Tests**
   - Repository integration tests (upsert idempotency, delete, re-tagging, multi-tag
     full-amount counting, year-range zero-fill, user scoping) behind an infra marker.
   - Consumer unit tests (each action → repository calls).
   - `test_aggregation_end_point.py` passes unchanged. `mypy` clean.
6. **Docs**
   - Rewrite the `analytic-api/AGENTS.md` architecture section (Kafka+Postgres
     projection, new env vars, consumer lifecycle, tag key-vs-value note).
   - Add `analytic-api` and its `→ Kafka, Postgres` dependency to the root `CLAUDE.md`.

### Out of scope (follow-up)

**Historical backfill.** The projection only reflects expenses created/updated/deleted
after the consumer goes live; pre-existing expenses appear once touched. A one-off
replay/REST sweep to seed history is a separate task.

