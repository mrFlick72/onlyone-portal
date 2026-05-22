# OnlyOne Portal Application Shell

React 19 + TypeScript frontend for OnlyOne Portal. The app is built with Vite as a multi-page application: each main area has its own generated HTML file and bundle, while shared UI and repository patterns stay in `src/`.

## Features

- OAuth2 Authorization Code with PKCE against vauthenticator.
- Home page with navigation to budget, revenue, plan, and account areas.
- Budget expense tracking with tag search, category totals, and attachment upload/download/delete.
- Budget revenue tracking with year search and attachments.
- User account detail management.
- Plan management with todo add/edit/delete and status transitions.

## Pages

| URL | Entry | Purpose |
|-----|-------|---------|
| `/` | `home/index.tsx` | Home page and global feature navigation. |
| `/callback` | `auth/Callback.tsx` | OAuth2 callback handler. |
| `/logout` | `auth/Logout.tsx` | Session termination. |
| `/budget/expense/index` | `budget/index.tsx` | Budget expenses. |
| `/budget/revenue/index` | `budget/index.tsx` | Budget revenue. |
| `/budget/search-tags/index` | `budget/index.tsx` | Search tag management. |
| `/account/index` | `account/index.tsx` | Account details. |
| `/plan/index` | `plan/index.tsx` | Plan list. |
| `/plan/detail?id=<planId>` | `plan/index.tsx` | Todos for one plan. |

## Build and Run

All npm commands run from `src/`.

```bash
cd portal/application-shell/src
npm install
npm run dev
npm run type-check
npm run build
npm run production-build
```

The Vite dev server listens on port `3000`. Production-like local serving uses nginx and the built `dist/` directory:

```bash
cd portal/application-shell/local
docker compose up
```

The local nginx URL is `http://local.onlyone-portal.com:8070`.

## Environment

Environment files live in `portal/application-shell/environments/`. `vite.config.ts` loads the file for the selected mode and exposes values through `import.meta.env`.

Required values:

```text
REDIRECT_URI
CLIENT_APPLICATION_ID
IDP_BASE_URL
AUTHENTICATION_CHECK_INTERVAL
BUDGET_API_BASE_URL
REVENUE_API_BASE_URL
TAG_API_BASE_URL
ACCOUNT_API_BASE_URL
PLAN_API_BASE_URL
```

## Source Layout

```text
src/
  auth/        OAuth2 PKCE, token storage, JWT validation, callback/logout
  account/     account detail page and repository
  budget/      expense, revenue, search tags, attachments
  plan/        plan list/detail pages, todo dialogs, status transitions
  components/  shared menu, form, and layout helpers
  config/      environment-backed ApplicationConfig
  messages/    hardcoded English message bundle
  theme/       Material UI theme
  time/        month/date helpers
```

## API Pattern

Repository files call backend APIs with browser `fetch`, bearer token from `sessionStorage.ACCESS_TOKEN`, and `credentials: include`. Base URLs come from `ConfigLoader.ts`; plan calls use `PLAN_API_BASE_URL`.
