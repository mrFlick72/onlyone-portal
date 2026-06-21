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
