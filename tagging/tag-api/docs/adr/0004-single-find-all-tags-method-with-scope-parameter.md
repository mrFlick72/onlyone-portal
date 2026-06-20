# FindAllTags takes a scope parameter; FindTagsByScope is merged into it, not kept as a sibling

When `Scope` was first added, `domain.TagRepository` got a new `FindTagsByScope(ctx, scope)` method deliberately kept separate from the existing `FindAllTags(ctx)`, specifically so the unscoped path stayed completely untouched and the change was as low-risk as possible (see the original issue discussion). At that point the two methods ran genuinely different queries: `FindAllTags` never touched `scope` at all, while `FindTagsByScope` always added a `FilterExpression`.

That distinction disappeared once [0003](./0003-scope-always-persisted-empty-string-default.md) landed: an empty `scope` now makes `FindTagsByScope` skip the `FilterExpression` and run the exact same `Query` as `FindAllTags`. From that point on, `FindAllTags` was just `FindTagsByScope(ctx, "")` — two names for what is, in every case, one query shaped by one optional parameter. Keeping both was pure duplication: two interface methods, two actions, two sets of tests, for one underlying behavior.

We collapsed them: `domain.TagRepository` now has a single method, `FindAllTags(ctx context.Context, scope string) ([]Tag, error)`, implemented exactly as `FindTagsByScope` was — `scope == ""` returns everything unfiltered, any other value filters on the normalized `scope` attribute. `domain.FindTagsByScope` (the action) is removed; `domain.FindAllTags` (the action) now takes the same `scope` argument and still appends the `UNKNOWN` sentinel regardless of it. The HTTP contract is unaffected: `GET /api/tags` calls the action with `scope = ""`, `GET /api/tags/scope/:scope` calls it with the path parameter — both routes now share one action and one repository method instead of two.

We kept the name `FindAllTags` rather than renaming everything to `FindTagsByScope`: it's the name already wired through `main.go` and both route handlers, and reads naturally for the call site that passes `""`. The optional parameter is the part that needed explaining, not the name.

## Considered Options

- Keep `FindAllTags(ctx)` and `FindTagsByScope(ctx, scope)` as two methods/actions (original design) — rejected once 0003 landed: the two query paths became identical except for an `if`, so two names for one behavior was just duplication to maintain.
- Drop `FindAllTags` and rename everything to `FindTagsByScope(ctx, scope)` — rejected: no functional difference from the chosen option, but it would have renamed every call site (route handlers, `main.go`, tests) for no behavioral gain.
- Single method, kept the name `FindAllTags(ctx, scope)` (chosen) — no duplicated query logic, minimal call-site churn, the existing name already reads fine when called with `""`.
