# Budget Expense

Tracks a user's budget expenses and revenue, each categorized by one or more tags from the [Tagging](../../tagging/tag-api/CONTEXT.md) context. Expense tags come from the `expense` Scope; revenue tags from the `revenue` Scope. Both aggregates store tag **keys** only and resolve each key to its current value from tag-api on read, so a renamed tag shows its new label without rewriting records.

## Language

**Untagged Expense**:
A `BudgetExpense` submitted for create or update with an empty tag list. It is not a distinct stored state — at persist time it is given the `UNKNOWN` Sentinel Tag (see Tagging context) so it always has at least one tag once saved.
_Avoid_: Uncategorized expense, default-tagged expense

**Untagged Revenue**:
A `Revenue` submitted for create or update with an empty tag list — the exact analogue of Untagged Expense. At persist time it is given the `UNKNOWN` Sentinel Tag so it always has at least one tag once saved. A `Revenue` row that predates tagging has no stored tag at all; on read its absent tag resolves to `UNKNOWN` identically, so legacy and explicitly-untagged revenue behave the same with no backfill. Revenue carries tags for categorization only — unlike expense it emits no events and has no by-tag totals aggregate yet (those arrive with revenue analytics; see `docs/adr/0002-revenue-tagging-mirrors-expense-without-events-or-totals.md`).
_Avoid_: Uncategorized revenue, default-tagged revenue

**UnknownSentinel**:
The single budget-api-side definition of the `UNKNOWN` Sentinel Tag (`tags.UnknownSentinel`), shared by both the expense and revenue default-if-missing helpers so the literal is written once per service. It mirrors the same convention tag-api synthesizes on read; the two services duplicate the string with no shared enforcement (see Tagging context and tag-api ADR 0001).
_Avoid_: Default tag, fallback tag

**Unresolvable Tag Reference**:
A stored tag `Key` on an `Untagged Expense`'s or `Untagged Revenue`'s record that no longer exists in tag-api's catalog — today, only because the tag was deleted (tag-api [ADR 0008](../../tagging/tag-api/docs/adr/0008-tag-update-is-value-only.md)). `GetTagBy` resolves it to the same `UnknownSentinel`, in place, at read time: no write touches the record's stored `Key`, and no analytics projection is updated (see [ADR 0003](./docs/adr/0003-deleted-tag-references-resolve-to-unknown-in-place.md)). Indistinguishable downstream from a record that was never tagged — both resolve to `UNKNOWN` the same way. Making the reclassification durable, and reflecting it in analytics, is deferred (#27).
_Avoid_: Orphaned tag, dangling tag reference

**Scheduled Expense**:
A per-user template stored separately from `BudgetExpense`, holding the owner's user name, an amount, notes, a tag list (tag **keys**, as on `BudgetExpense`), a **Day** and an optional **Month**. A daily evaluation compares today's day (and month, when set) to the template; on a match it creates a real `BudgetExpense` for the owner. The template itself is not an expense and is never counted in totals or analytics. Editing a Scheduled Expense affects only the expenses generated after the edit — expenses already generated are untouched.
With a Month set it recurs yearly (e.g. 15 March each year); with the Month empty it recurs monthly (e.g. the 5th of every month). A Day that doesn't exist in the current month (29–31) clamps to that month's last day, so a Day-31 definition still fires every month. A nullable **End Date** bounds it — null means it recurs forever, a set date means it stops once that date is reached. It can be paused and resumed, or deleted.

Every evaluation updates the definition's **Last Evaluated Date**. If the daily job did not run for one or more days (deploy, outage), the next run backfills: it walks each day from `Last Evaluated Date + 1` to today, capped at the End Date, generating a `BudgetExpense` for every day that matches, so downtime never silently drops a recurring cost.

A **paused** Scheduled Expense is skipped entirely by the daily job — not evaluated, not backfilled. On pause, `Last Evaluated Date` is stamped with the pause date; on resume, it is stamped with the resume date. No `BudgetExpense` is generated for the paused span.
_Avoid_: Recurring expense template, standing order, due date (ambiguous — the stored fields are Day and Month), recurrent expense (reserved for the recurrence behaviour, not the template)
