# golang-web-framework

Shared Go library used by all Gin-based services (budget-api, tag-api, account-api).

`WebServerProvisioner` auto-configures: CORS, JWT validation, health endpoint at `GET /management/health`, and OpenTelemetry (traces, metrics, logs).

## Package overview

| Package | Import path | Purpose |
|---------|-------------|---------|
| `config` | `.../config` | Viper-based YAML config singleton — `GetConfigFor`, `GetConfigBoolFor` |
| `logging` | `.../logging` | Zap logger with file rotation — `GetLoggerInstanceForComponentByTypeName` |
| `middleware/security` | `.../middleware/security` | JWT validation — `SetUpOAuth2()`, `GetCurrentUser(ctx)` |
| `web/server` | `.../web/server` | `WebServerProvisioner` + `WebServerConfigurer` lifecycle (OTel / standard middleware / OAuth2) — `ConfigureEngine()`, `StartEngine()`, `Shutdown(ctx)` |
| `web/management` | `.../web/management` | Health endpoint registration (`/management/health`) |
| `otel` | `.../otel` | OTel provider setup — `Setup(ctx)`, `SetupTracerProvider(ctx)` |
| `httpclient` | `.../httpclient` | OTel-aware HTTP client — `NewHTTPClient(opts...)` |
| `awsclient` | `.../awsclient` | OTel-aware AWS SDK v2 config loader — `LoadDefaultConfig(ctx, opts...)` |
| `cypto` | `.../cypto` | Symmetric crypto — `AesCbcCipher`, `KeyRepository`, `NewInMemoryKeyRepository()` |
| `cache` | `.../cache` | Cache provider interface |
| `cache/in_memory` | `.../cache/in_memory` | Ristretto in-memory cache implementation |
| `pkg/money` | `.../pkg/money` | Money/decimal helpers |
| `pkg/time/date` | `.../pkg/time/date` | Date helpers |

---

## Configuration

Set the `CONFIG_FILE_LOCATION` environment variable to the path of a YAML config file.

### Full configuration reference

```yaml
server:
  port: 3050
  # Optional HTTP timeouts (Go duration strings, e.g. "30s", "2m"). Defaults shown.
  read-header-timeout: 5s
  read-timeout: 30s
  write-timeout: 30s
  idle-timeout: 120s
  shutdown-timeout: 10s     # total budget for graceful shutdown (drain + configurer Dispose, including OTel flush)

idp:
  jwks-endpoint: http://local.api.vauthenticator.com:9090/oauth2/jwks

user:
  required-role: USER_ROLE

logger:
  level: info          # debug | info | warn | error
  file-name: logs.log

cors:
  allowed:
    origins: "http://local.onlyone-portal.com:8070"

otel:
  enabled: false             # true to enable OpenTelemetry (traces, metrics, logs)
  service-name: my-service   # reported as service.name in all telemetry signals
  protocol: http             # http (default, port 4318) | grpc (port 4317)
  endpoint: localhost:4318   # OTel Collector host:port, shared by all three signals
  insecure: true             # true for local/dev (no TLS); false for production

# Only required if you use cypto.NewInMemoryKeyRepository
key:
  in-memory:
    storage:
      key: my-key-id
      key-value: 0123456789abcdef0123456789abcdef   # 16/24/32 raw bytes for AES-128/192/256
```

---

## WebServerProvisioner

`WebServerProvisioner` is a thin host for `WebServerConfigurer`s. Each configurer owns one cross-cutting concern with two lifecycle hooks:

```go
type WebServerConfigurer interface {
    Name() string
    Configure() error              // runs at boot
    Dispose(ctx context.Context) error // runs at shutdown
}
```

`ConfigureEngine()` registers the production chain (the registration order is also the middleware order):

1. `OTelConfigurer` — `otel.Setup(ctx)` for global trace/metric/log providers; mounts `otelgin.Middleware` with a `/management/*` filter so health probes don't pollute traces. **Non-fatal**: if setup fails, logs and falls back to a no-op shutdown so the service still boots without tracing.
2. `StandardMiddlewareConfigurer` — `gin.Logger()`, `gin.Recovery()`, CORS.
3. `OAuth2Configurer` — calls `security.SetUpOAuth2(ctx)`. The lifetime `ctx` is owned by the configurer; `Dispose` cancels it, which stops the JWKS refresh goroutine.

After the configurers run, `management.RegisterEndpoints` mounts `GET /management/health` → `{"status": "UP"}` (no auth).

If any `Configure()` returns an error, `ConfigureEngine` invokes `Shutdown` (within the `server.shutdown-timeout` budget) before panicking — partial state is rolled back so a retry starts clean.

### Lifecycle and graceful shutdown

`StartEngine()` listens on `server.port` and blocks until the process receives SIGINT/SIGTERM or the listener fails. On signal it:

1. Calls `srv.Shutdown(ctx)` to drain in-flight requests.
2. Calls `provisioner.Shutdown(ctx)`, which iterates `Dispose(ctx)` on every configurer (cancels JWKS refresh, flushes OTel exporters) and joins every error so the process exit code reflects partial shutdown.
3. Resets the configurer slice and `engine` so a second call is a no-op and a panicking `Configure` cannot leave the provisioner half-built.

A `defer` in `StartEngine()` is the safety net for early returns (e.g. `ListenAndServe` failure) — same `Shutdown` path.

`server.shutdown-timeout` (default `10s`) is the **single budget** for the entire shutdown — drain + every `Dispose` (OTel flush included). Tune it once, in YAML.

```go
provisioner := server.WebServerProvisioner{}
ginEngine := provisioner.ConfigureEngine()
// register routes…
provisioner.StartEngine() // blocks until SIGINT/SIGTERM
```

Services that need to control shutdown timing themselves can call `provisioner.Shutdown(ctx)` directly with their own deadline.

### Adding a custom configurer

Append your own `WebServerConfigurer` after `ConfigureEngine()` returns and before `StartEngine()` if you need extra lifecycle-aware behaviour (background workers, custom middleware that owns goroutines, etc.). Each configurer should:

- Append itself to `wsp.configurers` in its constructor (see `NewOAuth2Configurer` for the pattern), so `Shutdown` will dispose it.
- Treat the `ctx` passed to `Dispose` as the shutdown deadline — **not** the lifetime ctx.

---

## OTel (`otel` package)

`otel.Setup(ctx)` is called automatically by `WebServerProvisioner`. It initialises the global `TracerProvider`, `MeterProvider`, and `LoggerProvider` from the YAML config.

- All three signals (traces, metrics, logs) share the same `otel.endpoint`.
- `otel.enabled: false` → nothing is installed, zero runtime overhead.
- Resource attributes: `service.name` (from config), `host.name`, `telemetry.sdk.*`, `OTEL_RESOURCE_ATTRIBUTES` env var.
- Propagation: W3C TraceContext + Baggage.

`SetupTracerProvider(ctx)` is available when only traces are needed (prefer `Setup` for new code).

---

## OTel-aware HTTP client (`httpclient` package)

Use `httpclient.NewHTTPClient()` for every outbound HTTP call that should participate in distributed traces.

```go
import "github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/httpclient"

client := httpclient.NewHTTPClient()
```

Behaviour driven by config:
- `otel.enabled: true` → transport wraps `http.DefaultTransport` with `otelhttp.NewTransport`: injects `traceparent`/`tracestate` headers and creates a client span.
- `otel.enabled: false` → returns a plain `&http.Client{}` with no overhead.
- Callers can pass `otelhttp.Option` values to customize span names or attributes for a specific outbound client.

**Required**: build requests with `http.NewRequestWithContext(ctx, method, url, body)` so the transport can read the active span from the context:

```go
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
// add headers…
resp, err := client.Do(req)
```

Using `http.NewRequest` (no context) silently disables propagation even with the instrumented transport.

The OAuth2 JWKS cache uses this client for auth-server key refreshes (refresh interval `15m`, boot fetch bounded by a `10s` timeout so an unreachable JWKS endpoint fails fast). With tracing enabled, refresh requests appear as `JWKS refresh` client spans.

---

## OTel-aware AWS SDK v2 client (`awsclient` package)

Use `awsclient.LoadDefaultConfig(ctx, opts...)` instead of `aws_config.LoadDefaultConfig` when creating any AWS service client (DynamoDB, S3, etc.).

```go
import (
    aws_config "github.com/aws/aws-sdk-go-v2/config"
    aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    awsclient "github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/awsclient"
)

cfg, err := awsclient.LoadDefaultConfig(ctx, aws_config.WithRegion("eu-central-1"))
client := aws_dynamodb.NewFromConfig(cfg)
```

Behaviour driven by config:
- `otel.enabled: true` → appends `otelaws.AppendMiddlewares` to the SDK API middleware chain; each AWS call becomes a traced child span with service, operation, and region attributes.
- `otel.enabled: false` → returns a plain `aws.Config` with no overhead.

**Required**: pass the request `ctx` to every AWS API call so the middleware can read the active span:

```go
result, err := client.PutItem(ctx, input)   // correct — span propagated
result, err := client.PutItem(context.TODO(), input)  // wrong — span lost
```

---

## Symmetric crypto (`cypto` package)

`cypto.Cipher` is a small interface for `Encrypt(plaintext) (ciphertext, error)` / `Decrypt(ciphertext) (plaintext, error)`. The default implementation is `AesCbcCipher` — AES-CBC with PKCS-style padding, base64 output, and a fresh random IV per encryption (prepended to the ciphertext).

Keys are looked up by id through a `KeyRepository` port:

```go
type KeyRepository interface {
    GetKeyFor(keyId string) (SymmetricKey, error)
}
```

`NewInMemoryKeyRepository()` is the bundled implementation — a single key loaded from config:

```yaml
key:
  in-memory:
    storage:
      key: my-key-id
      key-value: 0123456789abcdef0123456789abcdef   # raw key bytes — len must be 16/24/32 for AES-128/192/256
```

For production, write your own `KeyRepository` backed by a KMS / Secrets Manager and fetch keys lazily. The cipher only calls `GetKeyFor` on each Encrypt/Decrypt — caching is the repository's responsibility.

---

## Local dev with Docker (`grafana/otel-lgtm`)

`grafana/otel-lgtm` is an all-in-one image that bundles the complete LGTM stack — OTel Collector, Loki, Grafana, Tempo, and Mimir (Prometheus-compatible metrics) — into a single container. No separate config files needed.

**`local/docker-compose.yml`**
```yaml
services:
  lgtm:
    image: grafana/otel-lgtm:latest
    ports:
      - "3000:3000"   # Grafana UI
      - "4317:4317"   # OTLP gRPC
      - "4318:4318"   # OTLP HTTP
    environment:
      - GF_AUTH_ANONYMOUS_ENABLED=true
      - GF_AUTH_ANONYMOUS_ORG_ROLE=Admin
```

**Service `application.yml` for local dev**
```yaml
otel:
  enabled: true
  service-name: tag-api
  protocol: http
  endpoint: localhost:4318
  insecure: true
```

Start it:
```bash
docker compose -f local/docker-compose.yml up -d
```

Grafana UI is at `http://localhost:3000`. All datasources (Tempo, Mimir, Loki) are pre-configured and wired together — trace → logs correlation works out of the box via the `traceID` field.
