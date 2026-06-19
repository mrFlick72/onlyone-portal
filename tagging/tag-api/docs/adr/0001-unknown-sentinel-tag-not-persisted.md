# UNKNOWN sentinel tag is hardcoded, not persisted per-user

`budget-api` needs every expense to resolve to at least one valid tag, even when the user didn't pick one, so it can default untagged expenses without breaking `FindSpentBudget`'s per-key tag lookup. We considered seeding a real `UNKNOWN` row into every user's tag-api catalog (via `SaveTag`) so it would be an ordinary persisted tag, but rejected it: it requires onboarding every new (and existing) user with this row before they can safely create an untagged expense, and it permanently consumes one slot in a per-user keyspace that's otherwise entirely user-authored.

Instead, `FindAllTags` synthesizes and appends `Tag{Key: "UNKNOWN", Value: "UNKNOWN"}` to the catalog on every read — it never exists in DynamoDB, and `PUT /api/tags` cannot create, update, or collide with it. `budget-api` independently defaults an expense's empty tag list to the same literal `{Key: "UNKNOWN", Value: "UNKNOWN"}` at persist time. There is no shared code between the two services enforcing this string stays in sync — if one side ever changes it, the other must be updated by hand.

## Considered Options

- Real persisted per-user tag, seeded on user onboarding — rejected: requires onboarding work for every user, consumes catalog space.
- Hardcoded sentinel, synthesized only in `tag-api`'s read path (chosen) — no onboarding, no persisted row, but the literal string is duplicated across two services with no shared enforcement.
