# 0002 — Resolve UI language once at login, cache client-side

- Status: Accepted
- Date: 2026-09-15

## Context

The app has no runtime i18n library: `messages/MessageRepository.ts#getAllMessageRegistry(language = "it_it")` loads message bundles that are already compiled in at build time, and every one of the ten independent Vite page entries (`account`, `budget*`, `plan*`, `analytics`, `home`, …) calls it separately in its own `useEffect`, with nobody overriding the hardcoded `it_it` default. There is no client-side routing — navigation between pages is full-page (HTML links), so there is no shared in-memory app state that could carry a resolved language from one page to the next.

Account now gains a `locale` field (see `account/account-api/CONTEXT.md`), proxied from vauthenticator's standard OIDC `locale` claim, editable on `AccountDetailsPage`. We want every page — not just the Account page — to render in the account's chosen language.

Two things ruled out the obvious-looking alternatives:

- **Reading `locale` off a token already in hand** would be free (no network call) if it worked, but it doesn't: vauthenticator's `IdTokenEnhancer` only stamps `email` onto the ID token. Profile claims, `locale` included, are only available via an explicit `GET /userinfo` call — the same one `account-api`'s `FindAnAccount` already makes.
- **Calling account-api from every page** on load would work, but adds a new network dependency and a new failure mode to pages that have nothing to do with the Account feature today (e.g. `budget/expense`, `plan`), for a value that changes rarely.

## Decision

Resolve the language once per login, not once per page:

- `auth/Callback.tsx` — which every login already flows through, after token exchange and before the redirect into the app — fetches the account from account-api once and writes the resolved `locale` into `localStorage`.
- `getAllMessageRegistry()` reads that cached value instead of the hardcoded `it_it` default. No other page calls account-api for this purpose.
- `AccountDetailsPage` write-throughs the cache immediately on save, so a change is visible on the next page load without waiting for another login.
- If the cache is cold (cleared storage, or a page reached without a fresh login in this browser) **and** vauthenticator has no `locale` for the account (a new or legacy user), the app falls back to a browser-detected language: walk `navigator.languages` in order, match each entry's primary subtag against the supported set (`it`, `en`), and fall back further to the existing hardcoded `it_it` if nothing matches. This detected value is never written back to vauthenticator — it's recomputed on demand, not treated as a real user choice.

## Considered Options

- **Fetch account-api on every page load** — rejected: couples every page to account-api for a value that changes rarely, adds latency and a new failure mode app-wide.
- **Read `locale` from the ID token** — rejected: not available; vauthenticator's `IdTokenEnhancer` doesn't put profile claims on the token, only `/userinfo` does.
- **Scope to the Account page only, no app-wide reach** — rejected: the point of the feature is that the whole app renders in the chosen language, not just the page where it's set.
- **Auto-persist the browser-detected guess to vauthenticator on first detection** — rejected: would silently overwrite an account's genuinely-unset preference with a guess (borrowed device, VPN, mislabeled OS locale) as though the user had chosen it.

## Consequences

`localStorage` is per-browser, not per-account: a user logging in on a new browser/device sees the fallback chain (not their saved preference) until `Callback.tsx` runs there. If a user changes their language on another device mid-session here, this session's cache doesn't pick it up until its next login.
