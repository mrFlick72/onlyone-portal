# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

The frontend SPA for OnlyOne Portal — a React 19 + TypeScript multi-page application built with Vite and Material UI v6.

Take in consideration $typescript-expert and $react-expert tags for this task.
Do not forget to use context7 and mui-mcp MCP if available for further details

## Build Commands

All npm commands run from `src/` (that's where `package.json` and `vite.config.ts` live):

```bash
cd src/
npm install
npm run dev                # Vite dev server on :3000
npm run build              # development-mode build → ../dist
npm run production-build   # production-mode build → ../dist
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

Each entry point independently calls `authenticationChecker()` for token auto-refresh and mounts its own React tree into `<div id="app">`. Navigation between pages is full page navigation (HTML links), not client-side routing.

The three `budget*` entries all share `budget/index.tsx`, which mounts `SpentBudgetApp`. That component uses `createBrowserRouter` (react-router v7) to render the right sub-page (`BudgetExpensePage`, `BudgetRevenuePage`, `SearchTagsPage`) based on the URL path of the HTML that loaded it (`/budget/expense/index`, `/budget/revenue/index`, `/budget/search-tags`).

### Feature-Based Source Organization

```
src/
  auth/           # OAuth2 PKCE flow + JWT validation (jose library)
  config/         # ApplicationConfig loaded from env vars via dotenv-webpack
  account/        # User profile feature
  budget/         # Budget feature (expenses, revenue, search tags)
    spent-budget/ # Expense tracking sub-feature
    budget-revenue/ # Revenue tracking sub-feature
    search-tags/  # Tag management sub-feature
  components/     # Shared UI: Menu, form inputs, layout helpers
  messages/       # Hardcoded English message bundle (no i18n library)
  theme/          # MUI theme configuration
  time/           # Month/date utilities
```

Each feature follows a `domain/` pattern: type definitions in `domain/*.ts`, API calls in `domain/*Repository.ts`, React components at the feature level.

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
- Base URLs resolved from `config/ConfigLoader.ts`, which reads `import.meta.env.*` values inlined by Vite at build time

### Environment Config

`environments/.env.development` (and `.env.production`, when present) supply all backend URLs and OAuth2 settings. `vite.config.ts` calls `loadEnv(mode, '../environments', '')` and forwards each value into `define` so it's accessible as `import.meta.env.<NAME>` in source. The build mode (`development` / `production`) is selected by the `--mode` flag in the `build` / `production-build` scripts. Key config values: `IDP_BASE_URL`, `CLIENT_APPLICATION_ID`, `REDIRECT_URI`, `BUDGET_API_BASE_URL`, `REVENUE_API_BASE_URL`, `ACCOUNT_API_BASE_URL`, `TAG_API_BASE_URL`, `AUTHENTICATION_CHECK_INTERVAL`.
