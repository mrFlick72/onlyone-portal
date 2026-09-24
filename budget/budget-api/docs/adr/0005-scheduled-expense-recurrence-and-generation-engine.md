# 0005 — Scheduled Expense: day/month recurrence trigger, in-process daily generation engine

- Status: Accepted
- Date: 2026-09-22

## Context

Users want to define a recurring expense once (rent, a subscription, insurance) and have budget-api create a real
`BudgetExpense` for them automatically, instead of entering it by hand every period. This requires a new stored
template — the **Scheduled Expense** (see `CONTEXT.md`) — plus a mechanism that evaluates every user's templates and
generates expenses as they come due, without a request in flight to carry a JWT.

Several sub-decisions only make sense together, so they're recorded as one ADR rather than several.

## Decision

### Trigger shape: Day + optional Month, not a frequency enum

A Scheduled Expense stores a **Day** and an optional **Month**. The daily job compares today's day (and month, when
set) against the template. With Month set it fires yearly (15 March every year); with Month empty it fires monthly
(the 5th of every month). There is no `WEEKLY`/`DAILY`/custom-interval option.

This was chosen over a `frequency` enum (`WEEKLY`/`MONTHLY`/`YEARLY`) because it maps directly onto how users
describe these expenses ("the 5th of the month", "15 March every year") and needs no separate interval field or
next-occurrence-advance logic — the job just compares today's `(day, month)` to the template's, unconditionally,
every day. The trade-off: sub-monthly recurrence (weekly, biweekly) isn't representable. That's accepted as out of
scope for this iteration; today's typical scheduled expenses (rent, subscriptions, insurance) are monthly or yearly.

### Day 29–31 clamps to the month's last day

A Day that doesn't exist in the current month (e.g. 31 in April, or 29 in a non-leap February) **clamps** to that
month's last day rather than being skipped or rejected at creation. A Day-31 definition fires on 30 April, 28/29
February, etc. This was chosen over silently skipping those months (which would quietly under-generate a real
recurring cost) and over rejecting Days 29–31 at creation time (which would make "last day of the month" simply
unrepresentable). The check is `today.Day() == min(definition.Day, daysInMonth(today))`.

### End Date is nullable; pause is a full skip, not a backfilled gap

A nullable **End Date** bounds the recurrence — null means forever, a set date means the job stops generating once
that date has passed. This check runs unconditionally on every evaluation, including on a definition just resumed
from pause: there is no special "resume after End Date has passed" case, it simply becomes permanently inert.

A **paused** definition is skipped entirely by the daily job — not evaluated, not backfilled. Pausing and resuming
both stamp `LastEvaluatedDate` with the current date, so no gap accumulates while paused and no expense is
retroactively generated for the paused span when resumed. This is a deliberate difference from the downtime-backfill
behavior below: pausing is a deliberate user action to skip a span, not unplanned downtime, so it gets no catch-up.

The stamp happens only on an actual transition (#54): pausing an already-paused or resuming an already-active definition
is a no-op, with no write. Re-stamping an active definition (e.g. from a stale second tab) would otherwise silently wipe
a pending downtime backfill. "Today" is the UTC calendar day (`date.Today()` in `domain/time/date`) — the generation
engine (#55) must use the same definition, or stamps and evaluations drift by a day around midnight UTC vs local time.

### Downtime backfill via `LastEvaluatedDate`, generate-then-advance ordering

Each definition stores a `LastEvaluatedDate`, updated every time the job evaluates it. If the job didn't run for one
or more days (deploy, outage), the next run walks every day from `LastEvaluatedDate + 1` to today (capped at the End
Date) and generates an expense for each day that matches — so unplanned downtime never silently drops a recurring
cost, unlike a paused span which is deliberately never backfilled.

A definition that has never been evaluated (`LastEvaluatedDate` empty — i.e. just created) is evaluated for **today
only**: no backfill before its creation. Combined with the hourly schedule below, a definition created on its due day
still gets that day's expense within the hour (#55).

Per matching day, the job **generates the expense first, then advances `LastEvaluatedDate`**. A crash between the two
steps re-evaluates that day on the next run and produces a **duplicate** expense — accepted, since a duplicate is
visible in the expense list and trivially deletable. The alternative order (advance first, generate second) risks a
crash producing a **silent, permanent loss** of that day's expense, which is worse given these represent real costs
the user incurred.

### In-process `gocron` v2, single replica, no distributed lock yet

The evaluation runs as an in-process `github.com/go-co-op/gocron/v2` scheduler, wired as a `WebServerConfigurer`
(`Configure` starts it, `Dispose` calls `Shutdown()`), rather than a Kubernetes CronJob hitting an internal endpoint.
This keeps the whole feature inside budget-api's own deployable and lifecycle, with no new deployment artifact. It is
registered with the framework's `RegisterConfigurer` (golang-web-framework ADR 0006, #64), so `Dispose` runs within
the provisioner's `server.shutdown-timeout` budget and cancels an in-flight run.

It runs **once at startup and then hourly** (#55), not once a day: `LastEvaluatedDate` already makes each day evaluated
exactly once, so extra runs cost one Scan and do nothing, while they close the gap for a definition created on its due
day after that day's run, and make a restart a manual trigger. Singleton mode keeps a slow run from overlapping the next.

budget-api runs a single replica today, so no distributed lock or leader election is wired in. `gocron` v2 supports
`WithDistributedLocker`/`WithDistributedElector` for exactly this case; adopting one of those is the natural next step
if budget-api is ever scaled to multiple replicas, since two unelected replicas would otherwise double-generate every
matching day.

### No new field on `BudgetExpense`; traceability lives in `Notes`

The generated `BudgetExpense` gets no new schema field linking it back to its Scheduled Expense. Instead, its single
`Note` is the definition's `Notes` followed, on a new line, by a trace sentence naming the definition (#55):

```
Monthly rent, paid by bank transfer
Expense generated by the scheduled expense "Rent" with id: 1a2b3c
```

(just the trace line when the definition has no `Notes`). The `Description` in the trace is what tells a generated
expense apart from a manual one at a glance. This keeps `BudgetExpense` unchanged and generated expenses genuinely indistinguishable from manual ones in
storage — support/debugging traceability lives in free text, not a structured, queryable field. Consequence: nothing
programmatic (idempotency checks, cascade logic) can rely on this text; idempotency for a given day is `LastEvaluatedDate`
alone, not "does an expense already referencing this id exist."

### Generation reuses `CreateBudgetExpense` with a synthetic owner-only context

The job builds a `context.Context` carrying `security.User{UserName: &owner}` (no `AccessToken`, no `Authorities`) and
calls the existing `CreateBudgetExpense.Execute` unchanged. This works because that action only reads `UserName` from
the context — it does not call tag-api, and the Kafka publisher never touches `AccessToken` either. If
`CreateBudgetExpense.Execute` or the event publisher ever start requiring a real access token, this path breaks
silently unless re-verified.

The job never resolves anything through tag-api either: `rest_tags_repository` needs the user's access token. It reads
definitions with a cross-user `Scan` (`FindAllActive`) that skips tag resolution, and advances `LastEvaluatedDate` with a
conditional `UpdateItem` (row must still exist and be `ACTIVE`, so a definition deleted or paused mid-run is neither
upserted back as a stub nor advanced past its pause stamp).

### Tag display names are stored on the definition

The expense event on Kafka carries each tag's display name, and analytic-api's total-by-tag groups by that name — so a
generated expense with keys only would land in an unnamed bucket. The definition's `Save` therefore also stores each
tag's display name as of that save (`tag_names`, a key→name map), and the job copies keys and stored names onto the
generated expense (#55). Names are as fresh as the definition's last save — the same staleness a manually entered
expense already has, since it also carries the name current at entry time. User-facing reads of the definition keep
resolving live names from tag-api and ignore `tag_names`. Definitions saved before this change generate with empty
names until they are next saved. Alternatives rejected: a service credential for tag-api (IDP and tag-api authorization
changes for one job), and leaving names empty until analytic-api's per-user reindex.

### Delete is a hard delete of the definition only

Deleting a Scheduled Expense (#53) removes only its definition row; `BudgetExpense`s it already generated stay exactly
as they are. This follows from the decision above: generated expenses carry no structured link back to their
definition, so there is nothing to cascade to, and they are ordinary expenses the user can delete individually. No
soft-delete/tombstone either — pausing already covers "stop generating but keep the definition". Ownership is checked in
the domain (`FindFor` in the caller's own partition, as for Update), with a DynamoDB `attribute_exists(id)` condition on
`DeleteItem` as the backstop for a row that disappears between the check and the delete.

## Consequences

- Scheduled Expenses cannot express sub-monthly recurrence (weekly, biweekly); revisit the Day/Month model if that's
  ever needed.
- A crash at the exact instant between generate and advance produces a visible duplicate expense rather than a silent
  loss — acceptable, not invisible.
- Scaling budget-api beyond one replica requires adding `gocron`'s distributed locker/elector before it's safe to do
  so; deploying a second replica without that change will double-generate expenses.
  The replica count isn't enforced or even visible in this repo (there is no budget-api Helm chart here as of #55), so
  this is a deployment constraint to keep in mind, not one anything checks.
- The `Notes`-suffix trace is not machine-reliable (a user could edit or delete it); do not build future logic that
  depends on parsing it.
