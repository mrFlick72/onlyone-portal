# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

This is a **multi-context** repo (a polyglot monorepo of independent services), so domain documentation is split per service.

## Before exploring, read these

- **`CONTEXT-MAP.md`** at the repo root — it points at one `CONTEXT.md` per context. Read each one relevant to the topic.
- The per-service **`CONTEXT.md`** for the service(s) you're about to work in.
- **`docs/adr/`** at the root for system-wide decisions, and the service-scoped `<service>/docs/adr/` for context-specific decisions.

If any of these files don't exist, **proceed silently**. Don't flag their absence; don't suggest creating them upfront. The `/domain-modeling` skill (reached via `/grill-with-docs` and `/improve-codebase-architecture`) creates them lazily when terms or decisions actually get resolved.

## File structure

Multi-context layout (presence of `CONTEXT-MAP.md` at the root). Services live under top-level directories rather than a single `src/`:

```
/
├── CONTEXT-MAP.md
├── docs/adr/                          ← system-wide decisions
├── budget/
│   ├── budget-api/
│   │   ├── CONTEXT.md
│   │   └── docs/adr/                  ← context-specific decisions
│   └── analytic-api/
│       └── CONTEXT.md
├── tagging/
│   └── tag-api/
│       ├── CONTEXT.md
│       └── docs/adr/
├── plan/
│   └── plan-api/
│       └── CONTEXT.md
├── account/
│   └── account-api/
│       └── CONTEXT.md
├── portal/
│   └── application-shell/
│       └── CONTEXT.md
└── core-services/
    └── golang-web-framework/
        └── CONTEXT.md
```

Not every context has its own `docs/adr/` — that directory is created lazily too, only once a context has a decision worth recording.

## Use the glossary's vocabulary

When your output names a domain concept (in an issue title, a refactor proposal, a hypothesis, a test name), use the term as defined in the relevant `CONTEXT.md`. Don't drift to synonyms the glossary explicitly avoids.

If the concept you need isn't in the glossary yet, that's a signal — either you're inventing language the project doesn't use (reconsider) or there's a real gap (note it for `/domain-modeling`).

## Flag ADR conflicts

If your output contradicts an existing ADR, surface it explicitly rather than silently overriding:

> _Contradicts ADR-0007 (event-sourced orders) — but worth reopening because…_
