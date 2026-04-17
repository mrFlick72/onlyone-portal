# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

The frontend SPA for OnlyOne Portal — a React 19 + TypeScript multi-page application built with Webpack 5 and Material UI v6.

## Build Commands

All npm commands run from `src/` (that's where `package.json` lives):

```bash
cd src/
npm install
npm run build              # development build → ../dist
npm run production-build   # production build
npm run watch              # watch mode
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

Webpack builds **separate bundles** per page, each with its own HTML file. Entry points in `webpack.config.js`:

| Entry      | Source file               | Output HTML              | Purpose                        |
|------------|---------------------------|--------------------------|--------------------------------|
| `callback` | `auth/Callback.tsx`       | `callback.html`          | OAuth2 PKCE callback handler   |
| `logout`   | `auth/Logout.tsx`         | `logout.html`            | Session termination            |
| `home`     | `home/index.tsx`          | `index.html`             | Home/landing page              |
| `budget`   | `budget/index.tsx`        | `budget/index.html`      | Budget expenses + revenue      |
| `account`  | `account/index.tsx`       | `account/index.html`     | User profile management        |

Each entry point independently calls `authenticationChecker()` for token auto-refresh and mounts its own React tree into `<div id="app">`. Navigation between pages is full page navigation (HTML links), not client-side routing.

The **budget page only** uses React Router (HashRouter) internally for sub-routes: `/` (expenses), `/budget-revenue`, `/search-tags`.

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
- Base URLs resolved from `config/ConfigLoader.ts` (which reads `process.env.*` injected by dotenv-webpack)

### Environment Config

`environments/.env.development` and `.env.production` supply all backend URLs and OAuth2 settings. The `ENV` variable in npm scripts selects which file dotenv-webpack loads. Key config values: `IDP_BASE_URL`, `CLIENT_APPLICATION_ID`, `REDIRECT_URI`, `BUDGET_API_BASE_URL`, `REVENUE_API_BASE_URL`, `ACCOUNT_API_BASE_URL`, `TAG_API_BASE_URL`, `AUTHENTICATION_CHECK_INTERVAL`.
