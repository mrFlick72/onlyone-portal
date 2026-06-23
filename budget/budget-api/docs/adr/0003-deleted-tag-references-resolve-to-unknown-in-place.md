# 0003 — Deleted tag references resolve to UNKNOWN in place, at read time; no migration this iteration

- Status: Accepted
- Date: 2026-06-23

## Context

tag-api's tag catalog is gaining update and delete (previously create-only; see #28 and tag-api's ADR on value-only
update). Delete is the hard case for budget-api: both expense and revenue store a tag by **key** only and resolve it
to its current value live from tag-api on every read (`GetTagBy`, via the Ristretto-cached
`RestSearchTagRepository` — see this service's ADR 0001/0002). `GetTagBy` today returns an error when a key is not
found in the scoped catalog.

That error propagates destructively once tags become deletable:

- `DynamoDbBudgetExpenseRepository.fromDynamo` (single read) propagates the error — a request for one expense that
  references a deleted tag returns `500`.
- `FindByDateRange`'s loop logs the error and **skips** the record — an expense referencing a deleted tag silently
  vanishes from list views, and its amount drops out of any total computed over the list.
- Revenue's equivalent resolution path has the same shape.

So the moment tag-api can delete a tag, every expense/revenue still referencing it breaks — immediately, and
silently in the list case. Some response to a missing key is required before delete can ship at all.

This service already has a sentinel for exactly this shape of problem: `tags.UnknownSentinel()`
(`{Key: "UNKNOWN", Value: "UNKNOWN"}`), used today to default a record saved with **no** tags so every record
resolves to at least one. A deleted-tag reference has a different cause (the record once had a real tag; the
catalog entry is now gone) but an identical resolution-time symptom: "this key doesn't resolve to a real catalog
entry right now."

## Decision

`GetTagBy`'s not-found branch returns `tags.UnknownSentinel()` instead of an error, in both
`RestSearchTagRepository` and `RistrettoCachedSearchTagRepository` (the two implementations duplicate the same
not-found logic and need the same change). A deleted tag's references therefore read exactly like an untagged
record — same sentinel, same key — with **no distinction preserved** between "always untagged" and "was tagged,
the tag was deleted."

This is read-time and in-place only:

- **No write happens.** The expense/revenue record's stored `tag` attribute (the deleted key) is untouched in
  DynamoDB. Every read re-resolves it to `UNKNOWN` via the same fallback, indefinitely.
- **No analytics update.** analytic-api's Postgres projection, built from the `budget-api.expense` Kafka stream at
  the time the expense was created/updated, still reflects whatever tag the expense carried back then. It does not
  learn about the deletion.
- Errors that are **not** "key absent from the scoped catalog" — tag-api unreachable, auth failure, malformed
  response — are untouched and continue to propagate as errors. Only the specific "looked through the whole list,
  key wasn't in it" branch changes.

A dedicated follow-up (#27) makes the `UNKNOWN` reclassification durable: rewriting the stored `tag` attribute on
affected records, and firing a Kafka reindex event so analytic-api's projection catches up for the affected month.

## Why not migrate now

- **Locating affected records requires a scan.** Neither the expense nor the revenue DynamoDB key scheme has a GSI
  on tag key; finding every record referencing a deleted key means a strategy this feature doesn't otherwise need
  (scan-by-key, or iterate by month/year). Designing that scan, plus the new Kafka reindex event and analytic-api
  handler, is real, separable scope — see #27.
- **Read-time fallback is sufficient for correctness today.** Once `GetTagBy` stops erroring, no record is lost
  from a list and no read 500s. The only thing the migration buys beyond that is no longer needing the
  fallback — a performance/cleanliness improvement, not a correctness one.
- **Decoupling ships delete sooner.** Update/delete is useful on its own; gating it on a Kafka reindex pipeline
  that doesn't exist yet would block a self-contained, low-risk change behind a much larger one.

## Consequences

- A user who deletes a tag sees every expense/revenue that used it immediately show "Uncategorized" (`UNKNOWN`) —
  non-destructive, but the original tag value is unrecoverable from budget-api once gone from tag-api (tag-api
  itself does not soft-delete either).
- Two distinct situations — "this record was never tagged" and "this record's tag was deleted" — are
  indistinguishable downstream until #27 ships. Acceptable: the UI and analytics never needed to tell them apart.
- Until #27 lands, the underlying DynamoDB `tag` attribute can reference a tag key that no longer exists in
  tag-api's catalog — a transient, intentional inconsistency the read-time fallback fully masks.
