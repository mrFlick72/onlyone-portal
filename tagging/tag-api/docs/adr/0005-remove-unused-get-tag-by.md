# Remove TagRepository.GetTagBy instead of deprecating it in place

`GetTagBy(ctx, key)` has been on `domain.TagRepository` since the original scaffold and was never wired into a use case or HTTP route — `FindAllTags` is the only action that exists, and tag-api exposes no per-key route. Its only callers were tests and interface-satisfaction mocks. Two of those tests (`TestSaveTagPersistsNormalizedScope`, `TestSaveTagWithoutScopePersistsEmptyStringScopeAttribute`) used it as their read path to assert on Scope-persistence behavior from [0003](./0003-scope-always-persisted-empty-string-default.md); they were rewritten to assert via `FindAllTags` instead, which is the method actually wired to production.

We deleted the method rather than keeping it deprecated-in-place: it has zero callers, git history is the recovery path if a per-key lookup is ever needed, and budget-api's own `GetTagBy` (a different, used method on a different interface) already covers single-tag lookup by fetching `GET /api/tags` and filtering client-side — so there's no latent cross-service need for a per-key endpoint on tag-api today.

See https://github.com/mrFlick72/onlyone-portal/issues/15.

## Considered Options

- Keep `GetTagBy`, mark deprecated via doc comment — rejected: dead code with no callers and no scheduled consumer; a deprecation comment just defers the same deletion decision.
- Delete outright (chosen) — no production caller, no cross-service dependency, fully recoverable from git history if needed later.
