# Tag update is value-only; key and scope are immutable after creation

Adding update and delete to the tag catalog (previously create-only, #28) raises the question of what update may change. `Tag` is a `(Key, Value, Scope)` tuple: `Key` is server-generated and already immutable by construction — the client never mints or chooses it. `Value` is the user-authored label and the obvious target of "rename this tag." `Scope` is the open question: could update also re-classify a tag from `expense` to `revenue`?

Allowing `Scope` to change carries the same hazard as delete: a `Tag` referenced by real expense/revenue records, re-scoped, becomes unreachable from the scope those records still expect (`GET /api/tags/scope/expense` stops returning a tag now scoped to `revenue`), so budget-api's resolution of those records would need exactly the same non-destructive fallback this work is already adding for delete (see budget-api's ADR on the `UNKNOWN` fallback). Folding a second referential-integrity hazard into "update" — for a need (re-classifying a tag's domain) that's rare and already expressible as delete-old-create-new — buys nothing.

`PATCH /api/tags/scope/:scope/:key` therefore accepts only `{ "value": "..." }`. `Key` and `Scope` together are the resource's identity, carried in the path, not in the writable body; a request whose `:key` does not carry `:scope` is rejected with `404`, the same composite-identity check `DELETE` uses. Targeting the synthesized `UNKNOWN` sentinel with either verb is rejected with `400` — it is never persisted, so there is nothing to update or delete.

## Considered Options

- **Value + scope mutable** — rejected: re-scoping a referenced tag is exactly as destructive to budget-api reads as deleting it, so it would need the same non-destructive-fallback machinery delete requires, for a need that's rare and already achievable as delete + create.
- **Value-only mutable (chosen)** — keeps update free of the cross-service referential-integrity question entirely; `Key`/`Scope` stay the resource's fixed identity, `Value` is the only field a rename touches.
