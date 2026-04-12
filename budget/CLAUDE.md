# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

The parent `../CLAUDE.md` covers the monorepo as a whole (auth, framework, deployment). This file documents what is specific to the `budget/` subtree.

## Services in this directory

- `budget-api/` — Go (Gin). CRUD for both budget expenses and revenue. Persists to DynamoDB. Calls `tag-api` over REST for tag lookup.
- `revenue-api/` — Python FastAPI. Legacy revenue service; kept in place pending decommissioning once Go version is verified.
- `budget-exporter/` — Python. Data export job.
- `analytic/` — Python. Plot/visualization scripts plus a small `server.py`.
- `helm-charts/budget-api/` — Helm chart for `budget-api` deployment (the `old/` sibling is legacy).

## budget-api architecture (Go)

Hexagonal layout — keep changes inside the right layer:

```
domain/
  budget/expense/   # CreateBudgetExpense, UpdateBudgetExpense, FindSpentBudget, DeleteBudgetExpense + BudgetExpenseActionsFacade
  budget/revenue/   # CreateRevenue, UpdateRevenue, FindRevenue, DeleteRevenue + RevenueActionsFacade
  tags/             # SearchTagRepository port
  money/, time/     # value objects
adapter/
  budget/expense/dynamodb/  # DynamoDB impl of BudgetExpenseRepository + id provider
  budget/revenue/dynamodb/  # DynamoDB impl of RevenueRepository + id provider
  tags/rest/                # REST client for tag-api + Ristretto-cached decorator
web/
  budget/expense/   # package expense — endpoint, converter, representation for expense
  budget/revenue/   # package revenue — endpoint, converter, representation for revenue
config/             # composition root — NewBudgetExpenseActionsFacade, NewRevenueActionsFacade
main.go             # WebServerProvisioner + registers both expense and revenue endpoints
```

Key wiring facts:
- `config.NewBudgetExpenseActionsFacade()` / `config.NewRevenueActionsFacade()` are the entry points.
- DynamoDB table names come from config keys `budget-api.dynamo-db.budget-expense.table-name` and `budget-api.dynamo-db.revenue.table-name`.
- AWS region is hardcoded to `eu-central-1`.
- Tag repository is wrapped in `RistrettoCachedSearchTagRepository` — preserve cache invalidation semantics when modifying tag lookup.
- `web/budget/expense` imports domain expense as `domainexpense "...domain/budget/expense"` to avoid package-name clash. Same pattern in `web/budget/revenue` with `domainrevenue`.

## Revenue DynamoDB ID scheme

The Go revenue adapter preserves the Python `revenue-api` composite key layout so existing `BUDGET_REVENUE` records remain readable without migration:

- PK: `base64("<year>_<userName>")`
- RK: `base64("<month>_<day>_<salt>")`
- Full ID stored in `budget_id` attribute as `<pk>-<rk>`

Do not change this scheme without a data migration plan.

## Routes

| Aggregate | Method | Path |
|-----------|--------|------|
| Expense | GET | `/api/budget/expense?q=month=MM&year=YYYY` |
| Expense | POST/PUT/DELETE | `/api/budget/expense`, `/api/budget/expense/:id` |
| Revenue | GET | `/api/budget/revenue?q=year=YYYY` |
| Revenue | POST/PUT/DELETE | `/api/budget/revenue`, `/api/budget/revenue/:id` |

The revenue query param preserves Python's `?q=year=2023` wire format for frontend compatibility.

## Build & test (budget-api)

```bash
cd budget-api
go build -o app .
go test -tags test ./...                       # `test` build tag required; adapter tests need LocalStack
go test -tags test ./domain/... ./web/...      # unit tests only, no LocalStack needed
CGO_ENABLED=0 GOOS=linux go build -o app .    # cross-compile for Docker

# Docker image:
docker build -t mrflick72/budget/budget-api:1 -f ../../core-services/docker/ubuntu.Dockerfile .
```

DynamoDB adapter tests run against LocalStack on `http://localhost:4566` — see parent CLAUDE.md for required `AWS_*` env vars.

## Build & test (revenue-api, legacy Python)

```bash
cd revenue-api
python -m venv venv && source venv/bin/activate
pip install -r requirements.txt
PYTHONPATH=src pytest --cov=app --cov-report=term-missing --cov-report=html
export BUDGET_API_CONFIG_FILE_LOCATION=<path-to-config>   # required at runtime
```

## Cross-service dependencies

- `budget-api` → `tag-api` (REST, base URL from config key `tag-api.base-url`)
- `budget-api`, `revenue-api` → DynamoDB
- All services → vauthenticator JWTs (validation via shared Go framework / Python equivalent)

When changing the expense or tag contract, check `tag-api` and `portal/application-shell` budget bundle for downstream impact.
