---
status: accepted
---

# Silence k8s probe log noise, on two separate knobs

Every k8s liveness/readiness probe to `GET /management/health` produced two log lines: the `gin.Logger()` access log line, and an Info-level "skipping oauth2 evaluation for path" line from `NewOAuth2Middleware` (`middleware/security/jwt.go`), since `/management/*` is on the OAuth2 middleware's ignored-path list. Probes fire every few seconds per pod, so both drowned out real request/auth activity in the logs.

We addressed the two lines separately, since they come from different middleware layers with different audiences.

The gin access log line is now gated by a new boolean config key, `server.access-log.health-check-logging-enabled`. Viper's `GetBool` returns `false` for an absent key, so the default (no config change required in any consumer) is "skip `/management/health` from the access log." `StandardMiddlewareConfigurer.Configure()` builds the `gin.LoggerConfig.SkipPaths` list via a small pure function so the decision is unit-testable without needing to toggle `config.ConfigurationManager`'s `sync.Once` singleton mid-test.

The jwt.go "skipping oauth2 evaluation" lines (for both the ignored-path match and the OPTIONS-request match) moved from `LogInfofFor` to `LogDebugfFor`. These aren't specific to health checks — they fire for any OAuth2-ignored path or any OPTIONS preflight — so gating them behind the new access-log flag would have conflated two different concerns. The existing `logger.level` config (default `info`) already exists for exactly this: routine internal detail, off by default, available on demand.

## Considered Options

- Gate both log lines behind the same new boolean — rejected: the jwt.go lines aren't health-check-specific (they also fire on every CORS preflight and any future OAuth2-ignored path), so scoping them to a "health check logging" flag would be misleading and wouldn't help someone debugging, say, an OPTIONS request.
- Add a config-driven list of access-log skip paths instead of a boolean — rejected as unjustified flexibility; there is exactly one management endpoint today (`/management/health`) and no consumer has ever needed to skip a different path. A boolean is trivially extendable to a list later if a second endpoint appears.
- Remove the jwt.go lines outright instead of downgrading to Debug — rejected; they're genuinely useful when debugging why a request unexpectedly skipped (or didn't skip) auth, so Debug-level (not deleted) is the right tradeoff.

## Consequences

Re-enabling full probe visibility takes **two separate knobs**, not one: `server.access-log.health-check-logging-enabled: true` brings back only the gin access line; `logger.level: debug` is separately required to bring back the jwt.go skip line (and will also surface the four other Info-level lines jwt.go already emits per authenticated request). Anyone debugging a probe-related issue needs to know both exist.

This does not address the deeper reason probe-log investigation is hard: `gin.Logger()` writes to `gin.DefaultWriter` directly, bypassing the zap/OTel logging pipeline, so access log lines have no trace correlation even when `otel.enabled: true`. That is a larger, separate change and is out of scope here.
