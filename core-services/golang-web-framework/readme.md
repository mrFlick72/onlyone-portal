# golang-web-framework

Shared Go library used by all Gin-based services (budget-api, tag-api, account-api).

`WebServerProvisioner` auto-configures: CORS, JWT validation, health endpoint at `GET /management/health`, and OpenTelemetry (traces, metrics, logs).

## Package overview

| Package | Import path | Purpose |
|---------|-------------|---------|
| `config` | `.../config` | Viper-based YAML config singleton — `GetConfigFor`, `GetConfigBoolFor` |
| `logging` | `.../logging` | Zap logger with file rotation — `GetLoggerInstanceForComponentByTypeName` |
| `middleware/security` | `.../middleware/security` | JWT validation — `SetUpOAuth2()`, `GetCurrentUser(ctx)` |
| `web/server` | `.../web/server` | `WebServerProvisioner` — `ConfigureEngine()`, `StartEngine()`, `Shutdown(ctx)` |
| `web/magangement` | `.../web/magangement` | Health endpoint registration |
| `otel` | `.../otel` | OTel provider setup — `Setup(ctx)`, `SetupTracerProvider(ctx)` |
| `httpclient` | `.../httpclient` | OTel-aware HTTP client — `NewHTTPClient()` |
| `awsclient` | `.../awsclient` | OTel-aware AWS SDK v2 config loader — `LoadDefaultConfig(ctx, opts...)` |
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
```

---

## WebServerProvisioner

`ConfigureEngine()` builds the Gin router with this middleware chain (in order):

1. `otelgin.Middleware` — server span wrapping the full request (skips `/management/*`)
2. `gin.Logger()` / `gin.Recovery()`
3. CORS
4. JWT/OAuth2

`StartEngine()` starts the server and registers a `defer` shutdown for OTel providers. Services that need to control shutdown timing can call `provisioner.Shutdown(ctx)` directly.

```go
provisioner := server.WebServerProvisioner{}
ginEngine := provisioner.ConfigureEngine()
// register routes…
provisioner.StartEngine()
```

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

**Required**: build requests with `http.NewRequestWithContext(ctx, method, url, body)` so the transport can read the active span from the context:

```go
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
// add headers…
resp, err := client.Do(req)
```

Using `http.NewRequest` (no context) silently disables propagation even with the instrumented transport.

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
