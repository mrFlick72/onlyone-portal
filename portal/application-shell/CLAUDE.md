# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

The frontend SPA for OnlyOne Portal — a React 19 + TypeScript multi-page application built with Vite and Material UI.

Use the local `typescript-expert`, `typescript-react-reviewer`, or `vercel-react-best-practices` skills when a change touches TypeScript/React behavior or performance. Use Context7/MUI documentation tools when current library details are needed.

## Conventions

- **Object shapes use `type`, not `interface`.** Declare component `Props`, domain models, and other object shapes as `type X = { … }`. Reserve `interface` for the two cases it earns: declaration merging, or extending a third-party `interface`. This is not yet lint-enforced (the frontend has no linter — only `tsc --noEmit`), so apply it by hand. See `docs/adr/0001-prefer-type-over-interface-for-object-shapes.md`.

## Tech Stack

| Concern | Choice |
|---------|--------|
| Language | TypeScript 6 (`tsc --noEmit` only — no Babel) |
| UI library | React 19.2.5 with `react-dom` and `react-router` v7 (used inside `SpentBudgetApp` and `PlanApp`) |
| Build tool | Vite 8 (`@vitejs/plugin-react` 6) |
| UI kit | Material UI 9 (`@mui/material`, `@mui/icons-material`, `@mui/x-date-pickers`) on Emotion |
| Auth | OAuth2 PKCE against vauthenticator; token introspection via `jose` 6 |
| Date handling | `moment` 2.30 with `@date-io/moment` for the MUI date pickers |
| Forms | `react-imask`, `react-number-format`, `react-select` |
| Messages | YAML bundles (`en_en` + `it_it`) compiled in at build time via `@modyfi/vite-plugin-yaml`; no runtime i18n library |
| Tests | None configured (no test runner, no test scripts, no test deps). Behaviour is verified through `portal/e2e/ai/` Playwright MCP scenarios. |

## Build Commands

All npm commands run from `src/` (that's where `package.json` and `vite.config.ts` live):

```bash
cd src/
npm install
npm run dev                # Vite dev server on :3000
npm run build              # development-mode build -> ../dist
npm run production-build   # production-mode build -> ../dist
npm run type-check         # tsc --noEmit
```

There is no test framework configured — no test runner, no test dependencies, no test scripts.

## Local Dev Server

The app is served locally via an nginx Docker container on port 8070:

```bash
cd local/
docker compose up
```

This mounts `dist/` into nginx. After building, the app is at `http://local.onlyone-portal.com:8070` (requires a `/etc/hosts` entry pointing to localhost).
The nginx config serves `/`, `/budget`, `/account`, and `/plan` by resolving extensionless paths to the generated HTML files.

## Architecture

### Multi-Page App, Not a Single SPA

Vite builds **separate bundles** per page, each with its own HTML file. Entries are declared in `src/vite.config.ts` (`pages` map), and HTML files are generated at build time from `src/template.html` by a small `htmlTemplatePlugin` that injects the matching entry script:

| Entry           | Source file               | Output HTML                    | Purpose                        |
|-----------------|---------------------------|--------------------------------|--------------------------------|
| `home`          | `home/index.tsx`          | `index.html`                   | Home/landing page              |
| `callback`      | `auth/Callback.tsx`       | `callback.html`                | OAuth2 PKCE callback handler   |
| `logout`        | `auth/Logout.tsx`         | `logout.html`                  | Session termination            |
| `budgetExpense` | `budget/index.tsx`        | `budget/expense/index.html`    | Budget expenses                |
| `budgetRevenue` | `budget/index.tsx`        | `budget/revenue/index.html`    | Budget revenue                 |
| `budgetTags`    | `budget/index.tsx`        | `budget/search-tags/index.html`| Budget tag search              |
| `account`       | `account/index.tsx`       | `account/index.html`           | User profile management        |
| `plan`          | `plan/index.tsx`          | `plan/index.html`              | Plan list                      |
| `planDetail`    | `plan/index.tsx`          | `plan/detail.html`             | Plan todo detail               |
| `analytics`     | `analytics/index.tsx`     | `analytics/index.html`         | Budget expense analytics       |

Each entry point independently calls `authenticationChecker()` for token auto-refresh and mounts its own React tree into `<div id="app">`. Navigation between pages is full page navigation (HTML links), not client-side routing.

The three `budget*` entries all share `budget/index.tsx`, which mounts `SpentBudgetApp`. That component uses `createBrowserRouter` (react-router v7) to render the right sub-page (`BudgetExpensePage`, `BudgetRevenuePage`, `SearchTagsPage`) based on the URL path of the HTML that loaded it (`/budget/expense/index`, `/budget/revenue/index`, `/budget/search-tags/index`).

The two `plan*` entries share `plan/index.tsx`, which mounts `PlanApp`. `PlanApp` routes `/plan/index` to the plan list and `/plan/detail?id=<planId>` to the todo detail page.

### Feature-Based Source Organization

```
src/
  auth/             # OAuth2 PKCE flow + JWT validation (jose library)
  config/           # ApplicationConfig loaded from Vite env values
  account/          # User profile feature
  budget/           # Budget feature (expenses, revenue, search tags, attachments)
    expense/        # Expense tracking sub-feature
    revenue/        # Revenue tracking sub-feature
    search-tags/    # Tag management sub-feature
    attachment/     # File attachments shared by expense + revenue
  plan/             # Plan and todo management, including todo status transitions
  analytics/        # Budget expense analytics dashboard (charts + reindex)
  components/       # Shared UI: Menu, form inputs, layout helpers
  messages/         # YAML message bundles (en_en + it_it) loaded at build time; no runtime i18n library
    bundle/         # Per-feature message_bundle_<lang>.yaml files
  theme/            # MUI theme configuration
  time/             # Month/date utilities
```

Each feature follows a `domain/` pattern: type definitions in `domain/*.ts`, API calls in `domain/*Repository.ts`, React components at the feature level.

### Attachments

`budget/attachment/UploadAttachmentPopUp.tsx` is a single dialog reused by both `BudgetExpensePage` and `BudgetRevenuePage`. The parent page passes an `AttachmentTarget` (`{ budgetId, budgetType, date, attachmentId? }`) and a label bundle; the dialog handles upload, listing existing attachments, per-row download, and per-row delete on its own. Backend calls live in `budget/attachment/domain/AttachmentRepository.ts`:

- `saveAttachment(target, file)` → `POST /api/attachment` (multipart). Pass `target.attachmentId` to overwrite an existing attachment instead of creating a new one.
- `getAttachmentsFor(target)` → `GET /api/attachment/metadata/:budgetType/:budgetId` — returns `[]` on non-2xx so callers don't have to branch.
- `downloadAttachment(attachmentId, fileName)` → `GET /api/attachment/:attachmentId/content` — streams to a synthesized `<a download>` and revokes the blob URL.
- `deleteAttachment(attachmentId)` → `DELETE /api/attachment/:attachmentId`.

All four go through `getBudgetApiBaseUrl()` and the standard bearer-token / `credentials: include` envelope used by every other repository. Tooltip / aria-label strings are sourced from `messages/MessageRepository.ts` under the `attachment.popup.*` keys and wired through `OnlyonePortalPagesConfigMap.tsx` for both the expense and revenue page configs.

### Plans and Todos

The plan feature lives in `src/plan/` and calls the plan API through `plan/domain/PlanRepository.ts` using `PLAN_API_BASE_URL`.

- `plan/index.tsx` is the shared entry for both `plan` and `planDetail` Vite entries. It calls `authenticationChecker()` and mounts `PlanApp`, which wires `createBrowserRouter` with `/plan/index` → `PlanListPage` and `/plan/detail` → `PlanDetailPage`.
- `PlanListPage` lists the authenticated user's plans (`getAllPlans`), opens the create-plan dialog (`createPlan`), deletes plans with confirmation (`deletePlan`), and navigates to `/plan/detail?id=<planId>` via `window.location.href` (full-page navigation, like every other section).
- `PlanDetailPage` reads `?id=<planId>` from `URLSearchParams`, loads the plan with `getPlan`, and manages its todos through `addTodo`, `updateTodo`, `removeTodo`, and `changeTodoStatus`. It treats `201` and `204` as success.
- `domain/Plan.ts` defines the wire types (`Plan`, `Todo`, `TodoStatus`, plus the `NewPlan` and `TodoPayload` request shapes used by `createPlan` and `addTodo`/`updateTodo`).
- `domain/TodoStatus.ts` is the frontend source of truth for allowed status buttons and labels. Allowed transitions are `TODO → IN_PROGRESS|ABORTED` and `IN_PROGRESS → TODO|DONE|ABORTED`; `DONE` and `ABORTED` are terminal. The same module exposes `labelFor` and `colorFor` for the MUI Chip variants used in `TodoRow` / `ChangeTodoStatusPopUp`.
- `ChangeTodoStatusPopUp` only renders transitions allowed by `TodoStatus.ts`; the backend still enforces the same rule and returns `409` for invalid transitions.
- The Plan area is reachable from every other page through `components/menu/PlanPageMenuItem.tsx` (linked via the `Plans` entry in `GlobalPageNavigation`) and from the home page tile.

### Analytics

The analytics feature lives in `src/analytics/` and calls `analytic-api` through `analytics/domain/AnalyticRepository.ts` using `ANALYTIC_API_BASE_URL`.

- `analytics/index.tsx` is the Vite entry; it calls `authenticationChecker()` and mounts `AnalyticsApp` → `AnalyticsDashboardPage`. Reachable via `components/menu/AnalyticsPageMenuItem.tsx` (the `Analytics` nav item) and the home page tile.
- `AnalyticsDashboardPage` renders two `AnalyticBarChart`s (`@mui/x-charts`): **total by tag** (`findTotalByTag` → `PUT …/total-by-tag`, filtered by year + optional month + tags) and **total by year** (`findTotalByYear` → `PUT …/total-by-year`, over a year range + optional tag). Tag filters send tag **keys** (`SearchTag.key`); bars are labelled by tag **value**. Each chart owns loading/error/empty state and a retry token.
- **Reindex/recovery**: a "Reindex data" section calls `reindexBudgetExpense` → `POST …/reindex` `{fromYear, toYear}` to rebuild the user's projection from budget-api when events were missed. It reuses the year-range validation (max 20 years), shows a spinner + disabled button while running, reports the result via a `Snackbar`, and on success bumps both charts' retry tokens to reload. This is the only mutating call on the page.
- Label strings come from `messages/bundle/analytics/message_bundle_{en_en,it_it}.yaml`, mapped in `OnlyonePortalPagesConfigMap.analytics()`; the `AnalyticsPageMessageBundle` type in `MessageBundles.ts` keeps producer and consumer in sync (includes the `reindex` group).

### Authentication Pattern

All authentication is in `auth/Authenticator.ts`. Key details:
- OAuth2 Authorization Code + PKCE flow against vauthenticator IDP
- Tokens stored in `sessionStorage` as `ACCESS_TOKEN` and `ID_TOKEN`
- JWT validation via `jose` library using JWKS from the IDP's `.well-known/openid-configuration`
- `authenticationChecker()` runs a periodic interval to silently re-authenticate before tokens expire
- If tokens are invalid/missing, the user is redirected to the IDP authorize endpoint

### API Communication Pattern

All repository files (`*Repository.ts`) use the browser `fetch` API with:
- `Authorization: Bearer ${window.sessionStorage.getItem("ACCESS_TOKEN")}`
- `credentials: 'include'`
- Base URLs resolved from `config/ConfigLoader.ts`, which reads `import.meta.env.*` values inlined by Vite at build time. `applicationConfigLoader()` returns the full `ApplicationConfig`; per-API helpers (`getBudgetApiBaseUrl`, `getAccountApiBaseUrl`, `getTagApiBaseUrl`, `getPlanApiBaseUrl`, `getAnalyticApiBaseUrl`, …) wrap it.

### Messages and Page Configs

There is no runtime i18n library. The message strings that used to be served by the now-removed `portal/i18n-api` service are vendored into the frontend as YAML bundles under `messages/bundle/<feature>/message_bundle_<lang>.yaml`. Feature folders are `account`, `budget-expense`, `budget-revenue`, `common`, `home`, `plan`, `plan-detail`, `search-tags`; each ships an `en_en` and an `it_it` file. The YAMLs are parsed at **build time** by `@modyfi/vite-plugin-yaml` (`yaml()` in `vite.config.ts`) — nothing is fetched at runtime.

`messages/MessageRepository.ts` is the registry loader:
- `import.meta.glob("./bundle/**/*.yaml", { eager: true })` pulls in every language bundle (Vite needs a static glob, so all locales are loaded and then filtered).
- `flattenMessages` turns the nested YAML into a flat `MessageBundle` with dot-separated keys (e.g. `menu.budgetPage.label`).
- `getAllMessageRegistry(language = "it_it")` filters to one locale and merges every feature bundle into a single flat registry. The default locale is `it_it`.
- `getMessageFor(bundle, key)` returns the value or `""` (empty string) on a miss.

Pages load the registry once in a `useEffect` (`setMessageRegistry(getAllMessageRegistry())`) into `MessageBundle` state, then hand it to the config map.

`messages/OnlyonePortalPagesConfigMap.tsx` turns that flat registry into per-page label objects (`menuMessages`, dialog titles, button labels). It is stateless: pages call `new OnlyonePortalPagesConfigMap()` and pass the registry into the section method they need (`configMap.plan(registry)`, `configMap.budgetExpense(registry)`, `configMap.common(registry)`, etc.). The structured label types those methods return live in `messages/MessageBundles.ts`, shared by both the config map (producer) and the components (consumers) so a mismatch is a compile error.

Adding a new dialog or menu item normally means: add the key to the relevant `bundle/<feature>/message_bundle_*.yaml` (both `en_en` and `it_it`), map it in `OnlyonePortalPagesConfigMap.tsx`, and — if it introduces a new label object — declare its type in `MessageBundles.ts`.

### Environment Config

`environments/.env.development` (and `.env.production`, when present) supply all backend URLs and OAuth2 settings. `vite.config.ts` calls `loadEnv(mode, '../environments', '')` and forwards each value into `define` so it's accessible as `import.meta.env.<NAME>` in source. The build mode (`development` / `production`) is selected by the `--mode` flag in the `build` / `production-build` scripts.

Key config values: `IDP_BASE_URL`, `CLIENT_APPLICATION_ID`, `REDIRECT_URI`, `BUDGET_API_BASE_URL`, `REVENUE_API_BASE_URL`, `ACCOUNT_API_BASE_URL`, `TAG_API_BASE_URL`, `PLAN_API_BASE_URL`, `ANALYTIC_API_BASE_URL`, `AUTHENTICATION_CHECK_INTERVAL`.
