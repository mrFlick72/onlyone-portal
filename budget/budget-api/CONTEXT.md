# Budget Expense

Tracks a user's budget expenses and revenue, each optionally categorized by one or more tags from the [Tagging](../../tagging/tag-api/CONTEXT.md) context.

## Language

**Untagged Expense**:
A `BudgetExpense` submitted for create or update with an empty tag list. It is not a distinct stored state — at persist time it is given the `UNKNOWN` Sentinel Tag (see Tagging context) so it always has at least one tag once saved.
_Avoid_: Uncategorized expense, default-tagged expense
