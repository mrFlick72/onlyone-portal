# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

The frontend SPA for OnlyOne Portal — a React 19 + TypeScript multi-page application built with Vite and Material UI.

Use the local `typescript-expert`, `typescript-react-reviewer`, or `vercel-react-best-practices` skills when a change touches TypeScript/React behavior or performance. Use Context7/MUI documentation tools when current library details are needed.

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

Each entry point independently calls `authenticationChecker()` for token auto-refresh and mounts its own React tree into `<div id="app">`. Navigation between pages is full page navigation (HTML links), not client-side routing.

The three `budget*` entries all share `budget/index.tsx`, which mounts `SpentBudgetApp`. That component uses `createBrowserRouter` (react-router v7) to render the right sub-page (`BudgetExpensePage`, `BudgetRevenuePage`, `SearchTagsPage`) based on the URL path of the HTML that loaded it (`/budget/expense/index`, `/budget/revenue/index`, `/budget/search-tags`).

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
  components/       # Shared UI: Menu, form inputs, layout helpers
  messages/         # Hardcoded English message bundle (no i18n library)
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
- Base URLs resolved from `config/ConfigLoader.ts`, which reads `import.meta.env.*` values inlined by Vite at build time. `applicationConfigLoader()` returns the full `ApplicationConfig`; per-API helpers (`getBudgetApiBaseUrl`, `getAccountApiBaseUrl`, `getTagApiBaseUrl`, `getPlanApiBaseUrl`, …) wrap it.

### Messages and Page Configs

There is no i18n library. `messages/MessageRepository.ts` holds a single hardcoded English `MessageBundle`; `getMessageFor(bundle, key)` returns the value or the key itself as fallback.

`messages/OnlyonePortalPagesConfigMap.tsx` is the single place that turns the message bundle into per-page label objects (`menuMessages`, dialog titles, button labels). Pages call `new OnlyonePortalPagesConfigMap()` and pick the section they need (`configMap.plan(...)`, `configMap.budgetExpense(...)`, `configMap.budgetRevenue(...)`, etc.). Adding a new dialog or menu item normally means editing both `MessageRepository.ts` (the raw strings) and `OnlyonePortalPagesConfigMap.tsx` (the per-page mapping).

### Environment Config

`environments/.env.development` (and `.env.production`, when present) supply all backend URLs and OAuth2 settings. `vite.config.ts` calls `loadEnv(mode, '../environments', '')` and forwards each value into `define` so it's accessible as `import.meta.env.<NAME>` in source. The build mode (`development` / `production`) is selected by the `--mode` flag in the `build` / `production-build` scripts.

Key config values: `IDP_BASE_URL`, `CLIENT_APPLICATION_ID`, `REDIRECT_URI`, `BUDGET_API_BASE_URL`, `REVENUE_API_BASE_URL`, `ACCOUNT_API_BASE_URL`, `TAG_API_BASE_URL`, `PLAN_API_BASE_URL`, `AUTHENTICATION_CHECK_INTERVAL`.
