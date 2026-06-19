# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this module is

Shared Go library providing cross-cutting concerns (HTTP server bootstrap, JWT auth, caching, OpenTelemetry, AWS/HTTP clients, symmetric crypto) consumed by every Gin-based service in OnlyOne Portal (`budget-api`, `tag-api`, `account-api`, `plan-api`). Those services pull this in via a local `replace` directive with no version pinning — **any breaking change to an exported method signature breaks every consumer's build immediately.** When extending an existing interface, widen it additively (add new methods) rather than changing existing signatures — see `docs/adr/0001-additive-cache-provider-interface.md` for the precedent and reasoning.

For full package-by-package API reference and config YAML reference, see `readme.md`. For caching domain vocabulary (Layered Cache, L1/L2, Backfill, Fail-open, Namespace), see `CONTEXT.md`. This file covers commands and architecture not already spelled out there.

## Commands

```bash
go build -o golang-web-framework .

# unit tests — most packages need the `test` build tag (fixture.go files are gated by `//go:build test`)
go test -tags test ./...

# single package / single test
go test -tags test ./cache/... -run TestLayeredCache
go test -tags test ./cypto/... -v

# integration tests (testcontainers-backed Redis) — separate tag, not run by default
go test -tags "test integration" ./cache/redis/...
```

CI (`.github/workflows/build-golang-web-framework.yml`) runs build → `go test -tags test ./...` → `go test -tags "test integration" ./cache/redis/...`, with a `localstack` service (KMS only) available on `4566-4599` for the `cypto.AwsKmsKayRepository` tests. To run those KMS-backed tests locally, start LocalStack yourself: `cd test && docker compose up -d` (the same `docker-compose.yml` CI implicitly relies on), with `AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY`/`AWS_DEFAULT_REGION` set to any non-empty value.

For local OTel trace/metric/log visualization (Grafana + Tempo + Mimir + Loki, all pre-wired): `docker compose -f hack/docker-compose.yml up -d`, then point a consuming service's `otel.endpoint` at `localhost:4318` — see `readme.md` for the full `local/docker-compose.yml` example and Grafana URL.

## Architecture notes beyond the readme

**Two `KeyRepository` implementations, same `cypto.Cipher` port.** `cypto.NewInMemoryKeyRepository()` reads one key straight from YAML config (dev/test only). `cypto.NewAwsKmsKeyRepository(ctx)` is the production path: it holds a KMS-encrypted ciphertext blob per key id (from `key.aws-kms.storage.key` / `key.aws-kms.storage.key-value`), calls KMS `Decrypt` on first use, and caches the plaintext key in memory under a `sync.RWMutex` thereafter — `kms.Client` is built via `awsclient.LoadDefaultConfig`, so it gets the same OTel span instrumentation as DynamoDB/S3 calls when `otel.enabled: true`. Both `KeyRepository.GetKeyFor` and `Cipher.Encrypt`/`Decrypt` take a `context.Context` (not mentioned in the parent repo's `CLAUDE.md`, which predates this change).

**`WebServerConfigurer` boot-failure handling.** If any configurer's `Configure()` returns an error during `ConfigureEngine()`, the provisioner calls `Shutdown` (bounded by `server.shutdown-timeout`) to unwind whatever already booted, *then* panics — so a process that fails to start never leaves dangling goroutines (JWKS refresh, OTel exporters) behind, and a supervisor-triggered restart begins from clean state.

**Cache provider error contract is enforced per-provider, not centrally.** `RistrettoCacheProvider` and `RedisCacheProvider` each independently guarantee fail-open behavior (`GetContext` → miss on error, `SetContext`/`EvictContext` → always `nil`) rather than relying on `LayeredCacheProvider` to swallow errors at the composite boundary. This was a deliberate choice (`docs/adr/0002-cache-providers-fail-open.md`) so each provider is safe to hand to a consumer standalone, not only when wrapped. Keep this invariant if you add a third provider (e.g. Memcached).

**`cache.CacheProvider` is context-only — no `Get`/`Set`/`Evict`.** Those non-context methods existed as additive-era backward-compat shims (`docs/adr/0001-additive-cache-provider-interface.md`); they were removed once an audit confirmed `budget-api` was the only real caller of `CacheProvider` anywhere in the monorepo and had migrated to `GetContext`/`SetContext` (`docs/adr/0004-remove-non-context-cache-provider-methods.md`). Use the `*Context` methods.

**`LayeredCacheProvider` does not invalidate other replicas.** A write/evict on one replica only updates that replica's L1 and the shared L2 — other replicas keep serving stale L1 reads until their own L1 TTL expires (`docs/adr/0003-cross-replica-l1-staleness-accepted.md`). This is an accepted bound, not a bug; don't "fix" it without raising a new ADR, since pub/sub invalidation was explicitly deferred as unjustified complexity.

**Test isolation by build tag**, three tiers:
- no tag: pure unit tests, always run, no external deps.
- `test`: unit tests that need shared fixtures (`fixture.go` under `//go:build test`) — required for `go test ./...` to even compile in some packages.
- `test integration`: `cache/redis/redis_cache_provider_integration_test.go` spins up a real Redis via `testcontainers-go/modules/redis`; needs Docker, not run in the default CI test step.
