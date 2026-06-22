# 0002 — Revenue tagging mirrors expense storage/resolution, without events or the totals aggregate

- Status: Accepted
- Date: 2026-06-21

## Context

tag-api now treats `Scope` as authoritative (its ADR 0007): scope is mandatory on write and reads are strict, and a
fully functional `revenue` tag catalog exists, maintained through the frontend's per-scope tag-management UI (#20/#21).
Nothing consumes it yet — revenue *records* cannot carry tags. #22 wires revenue records to revenue-scoped tags so they
can be categorized the way expenses already are.

This ADR is where two forward-looking notes from the approved [ADR 0001](./0001-expense-scoped-tag-lookup-hardcoded-at-wiring.md)
are resolved (rather than editing that record):

- ADR 0001 anticipated a "future second-scope caller" sharing budget-api's tag-lookup code path — **that caller is
  revenue**, introduced here.
- ADR 0001's consequence that tag-api's unscoped `GET /api/tags` is "not removed" is now **historical**: it was removed
  when the `SearchTagsPage` redesign it depended on shipped (tag-api ADR 0007 / #20–#21). That change is unrelated to
  this ADR's decision but is noted so a reader of ADR 0001 finds the current state here.

Expense tagging is built from several parts:

1. **Storage**: a `BudgetExpense` stores its tags as a comma-joined string of tag **keys** in the DynamoDB `tag`
   attribute. Values are not stored.
2. **Resolution on read**: `DynamoDbBudgetExpenseRepository.fromDynamo` resolves each stored key to its current value
   via an injected expense-scoped `SearchTagRepository` (`GetTagBy`, backed by `GET /api/tags/scope/expense`, Ristretto
   cached). The live catalog value wins, so a renamed tag shows its new label without rewriting records.
3. **Default-if-missing**: `applyDefaultTagIfMissing` stamps the `UNKNOWN` sentinel on any expense saved with no tags, so
   every expense has ≥1 tag (tag-api ADR 0001 guarantees `UNKNOWN` is always resolvable without per-user onboarding).
4. **Totals aggregate**: `SpentBudget.TotalForSearchTags()` groups spend by tag for the expense page's by-tag view.
5. **Events**: create/update/delete publish to the `budget-api.expense` Kafka topic, consumed by analytic-api.

Parts 4 and 5 exist **only** to serve analytics/reporting. #22 explicitly defers revenue analytics.

## Decision

Revenue tagging replicates parts 1–3 and deliberately omits parts 4–5.

- **Storage (mirror)**: `Revenue` gains `Tags []tags.SearchTag`. The DynamoDB revenue item gains a `tag` attribute
  holding a comma-joined string of tag **keys**, exactly like expense. The revenue key scheme (PK/RK/`budget_id`) is
  unchanged — this is a new non-key attribute only, so no data migration of the key layout.
- **Resolution on read (mirror)**: `DynamoDbRevenueRepository` gains an injected revenue-scoped `SearchTagRepository`
  and resolves stored keys to values in `fromDynamo`, via `GET /api/tags/scope/revenue` (Ristretto cached under
  `search_tags_user_<userName>_scope_revenue`).
- **Read-tolerant default for legacy rows**: unlike expense, real revenue rows already exist with **no** `tag`
  attribute. Revenue's `fromDynamo` treats a missing/empty `tag` attribute as the `UNKNOWN` sentinel rather than
  asserting the attribute is present (expense asserts presence and would panic on a tagless row). Combined with
  `applyDefaultTagIfMissing` on write, every revenue — legacy or new — resolves to at least `UNKNOWN`. **No backfill.**
- **Shared sentinel**: the `UNKNOWN` sentinel is extracted to a single `tags.UnknownSentinel` in budget-api's domain and
  used by both expense's and revenue's default-if-missing helpers, so the literal is defined once per service. (The
  cross-service duplication — budget-api ↔ tag-api ↔ frontend — is inherent to separately deployed services and remains,
  as tag-api ADR 0001 already notes.)
- **Wiring**: two named constructors over a shared private helper —
  `NewExpenseSearchTagRepository()` / `NewRevenueSearchTagRepository()` — preserving ADR 0001's wiring-layer-literal
  pattern across two scopes. Revenue uses exactly **one** instance (injected into its repository); there is no second
  instance because there is no `FindSpentRevenue`.
- **No totals aggregate**: there is no `SpentRevenue` / `TotalForSearchTags` for revenue. `FindRevenue` returns
  `[]Revenue`, each now carrying resolved tags, and nothing more.
- **No events**: revenue create/update/delete publish nothing. Revenue actions remain `EventPublisher`-free.

## Why omit events and totals now

A future reader will reasonably ask "why does revenue carry tags but, unlike expense, emit no events and expose no by-tag
totals?" Because both exist solely to feed analytics, and revenue analytics is a separate, not-yet-built follow-up:

- A revenue Kafka stream with no consumer is speculative scaffolding — a new topic and event wire type maintained for no
  reader. It is cheaper and safer to introduce it together with the consumer that needs it.
- A `SpentRevenue` totals aggregate with no UI/endpoint to surface it is dead domain code.

Tagging a revenue record (store, resolve, display) needs neither. They are added with the revenue-analytics work, not
here.

## Consequences

- Revenue records can be created/updated with multiple revenue-scoped tags; tags display on revenue rows with
  live-resolved values; untagged and legacy revenue resolve to `UNKNOWN`. No migration.
- budget-api now resolves tags under two scopes, each with its own cache; expense behavior is unchanged.
- Revenue search stays year-only (`GET /api/budget/revenue?q=year=YYYY`) — no tag filter; that is reporting, deferred
  with analytics.
- When revenue analytics is built, it must add the revenue event stream and (if needed) a totals aggregate; the storage
  and resolution this ADR introduces are the foundation it builds on.
- The asymmetry between expense (events + totals) and revenue (neither) is intentional and temporary; revisit when
  revenue analytics lands.
