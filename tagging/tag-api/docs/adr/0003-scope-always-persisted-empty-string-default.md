# Scope is always persisted with an empty-string default, and an empty Scope query means "no filter"

The original plan for this iteration only wrote the `scope` attribute when the caller supplied a non-empty value, leaving it entirely absent otherwise — the natural way to model an "optional" field in DynamoDB, and consistent with how a sparse GSI would later treat it (see [0002](./0002-scope-filter-without-gsi.md)). We're changing that before this ships: starting with this iteration, every `SaveTag` call writes the `scope` attribute unconditionally, defaulting to `""` when the caller doesn't provide one. "Optional" now describes only what the caller must supply on write — not whether the attribute exists in storage.

The reason is the read-side contract: querying `FindAllTags(ctx, scope)` with an empty `scope` is meant to mean "no filter, give me everything," not "give me tags whose Scope is literally the empty string." At this stage, before any `Tag` carries a real, non-empty `Scope`, that's the only way nearly every tag stays reachable through the scope-filtered read path. So an empty (post-normalization) `scope` skips the `FilterExpression` entirely, rather than filtering on `scope = ""`.

A `Tag` saved before this attribute existed has no `scope` attribute at all. On read, that absence defaults to `""`, the same value an explicitly-saved-empty `Tag` would have — the two are indistinguishable once read, and both are reachable via an empty-`Scope` query and absent from any specific, non-empty `Scope` query until resaved with one. We accept losing the distinction between "never touched by this feature" and "explicitly has no Scope" — there is exactly one bucket for "no real Scope assigned," not two — in exchange for a uniform read/write contract.

## Considered Options

- Keep `scope` a sparse/omitted attribute when empty (original plan) — rejected: makes "give me everything" a special no-filter code path inconsistent with how every other `Scope` value is read, and conflates "predates this feature" with "explicitly unscoped" without actually distinguishing them anywhere that matters.
- Always persist `scope`, defaulting to `""`, with an empty `scope` query skipping the filter (chosen) — uniform schema for every `Tag` saved from this iteration forward; missing-attribute and empty-string read identically.
