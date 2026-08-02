# 0004 — A tag-catalog fetch must validate HTTP status before an empty result can be trusted, and the durable write gets a kill switch

- Status: Accepted
- Date: 2026-08-02

## Context

ADR 0003 made `GetTagBy`'s not-found branch return `tags.UnknownSentinel()` instead of an error, specifically
because "key absent from the scoped catalog" needed to be distinguishable from every other failure mode, which
"are untouched and continue to propagate as errors." That distinction is the entire safety argument in ADR 0003:
read-time `UNKNOWN` fallback is safe to apply broadly (including, since #27, durably rewriting the stored `tag`
attribute) only because it is assumed to fire exclusively when tag-api's catalog was fetched successfully and
definitively does not contain the key.

`RestSearchTagRepository.GetAllTags` (`adapter/tags/rest/rest_tags_repository.go`) does not actually uphold that
assumption. It calls `Client.Do(req)`, reads the body, and unmarshals it into `[]tags.SearchTag` — but never
inspects `resp.StatusCode`. Any response whose body still unmarshals into a valid (typically empty) array is
treated as a successful fetch:

- A degraded tag-api instance responding `200 []` (e.g. a defensive empty-list fallback on an internal error).
- A gateway/load balancer serving a stub `200 []` health-check-style body instead of proxying the real request.
- A misrouted `404` whose body happens to be `[]`.

None of these are "key absent from the catalog" — they are "the catalog could not be confirmed" — but `GetTagBy`
cannot tell the difference from an empty slice with `err == nil`, so it applies the ADR 0003 fallback anyway. See
#33 for the full trace: this makes every `tagKey` looked up during the bad response window resolve to `UNKNOWN`,
which #27's durable reclassification then writes back to DynamoDB for every affected expense in the query range —
an entire month's expenses can be flattened to `UNKNOWN` from a single degraded response, without a Go `error`
ever existing to log.

Hard failures (connection refused, DNS failure, timeout) are not the gap — those already produce a Go `error` from
`Client.Do` and propagate correctly today. The gap is specifically a *reachable* response that is wrong.

`RistrettoCachedSearchTagRepository.GetAllTags` delegates to the same `RestSearchTagRepository` (or an equivalent)
for its cache-miss path and only skips caching when `err != nil`, so the same false-empty result also gets cached
under `search_tags_user_<user>_scope_<scope>` and served as ground truth for the Ristretto TTL — widening a single
bad response into every request for that user/scope during that window.

## Decision: validate the HTTP status before trusting an empty catalog

`RestSearchTagRepository.GetAllTags` must validate `resp.StatusCode` before unmarshaling: any non-2xx response
returns an error and never reaches the `[]tags.SearchTag` unmarshal path. A successful, empty `[]` result is only
ever produced by an actual 2xx response from tag-api.

This is the minimal change that restores the invariant ADR 0003 already assumed it had: `GetTagBy` returns
`UnknownSentinel()` if and only if the catalog fetch definitively succeeded and the key was definitively absent
from it. Everything else — unreachable, degraded, misrouted, malformed — must error, exactly as ADR 0003's own
"Errors that are not key-absent... are untouched" line already claimed.

No change is needed in `RistrettoCachedSearchTagRepository` itself: it already skips `setToTheCache` whenever the
delegate returns an error, so fixing the delegate's status handling closes the caching amplification for free.

## Decision: a config-gated kill switch for the durable write

Status validation closes the specific gap this ADR was written for, but it does not — and cannot — close the
residual risk below: an inference-from-absence design has an irreducible failure mode as long as tag-api can
return a *bona fide* 2xx empty/wrong catalog. Given #33 was found in production behavior, not in review, the
durable side of reclassification gets an explicit, config-driven kill switch rather than shipping the status-code
fix alone and hoping it's sufficient.

`DynamoDbBudgetExpenseRepository` (`adapter/budget/expense/dynamodb/dynamo_db_budget_expense_repository.go`) gains
a `ReclassificationEnabled bool` field, wired from a new config key `budget-api.reclassification.enabled` in
`config/configurations.go`. The gate sits in `publishReclassification` — the producer, at the point detection
happens — not in `UpdateBudgetExpense.reclassify` on the consumer side of the event bus:

- The read-time, non-destructive `UNKNOWN` fallback from ADR 0003 (`GetTagBy`, `fromDynamo`) is **untouched** and
  always active — it writes nothing and is the behavior the frontend depends on for a deleted tag to render as
  "Uncategorized" at all.
- `FindByDateRange` still detects `hasDeletedTagReference` and still calls `publishReclassification` for every
  detected record — detection itself does not change.
- When `ReclassificationEnabled` is `false`, `publishReclassification` logs "reclassification disabled by config:
  budget expense `<id>` would have been durably reclassified to tags `<tags>`" and returns **without publishing to
  `EventBus`**. `UpdateBudgetExpense.reclassify` never runs at all — no `Save`, no re-published Kafka event, and
  crucially no event sits on the bus either. When `true`, the event is published and behavior is exactly what #27
  shipped.

Gating at the producer rather than the consumer means the flag has exactly one decision point for the whole durable
path, instead of two components (`DynamoDbBudgetExpenseRepository` publishing and `UpdateBudgetExpense` separately
deciding whether to act) that would otherwise have to agree — publishing an event that the consumer then silently
drops is both a duplicated decision and a wasted round trip through the bus for something that was never going to
happen.

**Default is disabled.** `GetConfigBoolFor` returns Go's zero value (`false`) when the config key is absent, and
the `DynamoDbBudgetExpenseRepository{}` struct literal's zero value for a `bool` field is likewise `false` — the two
defaults agree, so a config file that hasn't been updated to mention this key is fail-safe by construction, not by
convention. Existing deployments must explicitly set `budget-api.reclassification.enabled: true` to keep the #27
behavior; until then, the durable-write blast radius described in this ADR's Context is fully inert while the log
line still gives full visibility into what *would* have happened, so behavior in production can be watched safely
before the write path is turned back on.

## Alternatives considered

- **Add a distinct "confirmed empty" vs. "fetch degraded" signal from tag-api** (e.g. a response header, or
  refusing to omit tags the caller doesn't have permission to see). Rejected for now: no such signal exists in
  tag-api today, and inventing one is materially larger scope than the bug being fixed. Worth revisiting if status
  validation alone proves insufficient in practice.
- **Move the reclassification decision out of the read path entirely** (e.g. require an explicit tag-deletion
  event from tag-api instead of inferring deletion from catalog absence). Rejected: this reopens the design ADR
  0003 and #27 already settled, and the inference itself is sound — it is only unsound when the "successful fetch"
  precondition is violated, which is what this ADR fixes.
- **Gate the flag at consume time instead of publish time** (leave `publishReclassification` unconditional and
  check the flag in `UpdateBudgetExpense.reclassify` before the listener acts). Rejected: it still publishes every
  detected event onto the bus regardless of the flag, so the event exists and is only discarded one hop later —
  a needless extra component in the decision, and it makes it possible for detection, publication, and the actual
  write decision to disagree if the wiring ever forgets to pass the same config value to both. Gating at the
  producer keeps the decision in one place and matches the actual requirement: don't publish an event for
  something that isn't going to happen.
- **Ship the status-code fix alone, no flag.** Rejected: the residual risk below (a true 2xx with a wrongly-empty
  body) is real and was not something this codebase had reasoned about before #33; a kill switch buys a window to
  observe production behavior via logs before trusting the write path again, at near-zero implementation cost.

## Consequences

- A tag-api outage or degraded response now correctly surfaces as an error on every affected lookup — expenses
  referencing tags looked up during that window are skipped from `FindByDateRange` results (the pre-existing,
  intentional behavior for genuine errors per ADR 0003) rather than being reclassified.
- The residual risk is unchanged for the case ADR 0003 always accepted: a tag-api response that is **both** a true
  2xx **and** wrongly empty (e.g. a request-routing or auth bug that legitimately returns 200 with the wrong
  user's/scope's empty catalog) is still indistinguishable from "this user has no tags in this scope" and will
  still drive reclassification. Status validation narrows the failure surface to that one case; it does not
  eliminate the inherent ambiguity of inferring deletion from absence.
- Every non-2xx branch of tag-api's contract (auth failure, 5xx, 404) must now be exercised in
  `RestSearchTagRepository`'s tests to confirm each one returns an error rather than falling through to the
  unmarshal path.
- Until `budget-api.reclassification.enabled` is explicitly set to `true` in a deployment's config, deleted-tag
  records stay read-time-only `UNKNOWN` (ADR 0003) indefinitely — the underlying DynamoDB `tag` attribute is never
  durably rewritten and analytics is never told about it, exactly as ADR 0003 already accepted as the pre-#27
  status quo. Re-enabling the flag is required to get #27's durable behavior back.
