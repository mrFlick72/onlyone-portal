# golang-web-framework

Shared Go library used by all Gin-based services (budget-api, tag-api, account-api).

`WebServerProvisioner` auto-configures: CORS, JWT validation, health endpoint at `GET /management/health`, and OpenTelemetry (traces, metrics, logs).

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
  endpoint: localhost:4318   # OTel Collector host:port
  insecure: true             # true for local/dev (no TLS); false for production
```

### OpenTelemetry notes

- All three signals (traces, metrics, logs) are sent to the same collector endpoint.
- When `otel.enabled: false` (the default), no providers are installed and there is zero overhead.
- The `otelgin` middleware is position-first in the chain so the server span covers CORS and JWT auth time. Health probes (`/management/*`) are excluded from tracing.
- Call `provisioner.Shutdown(ctx)` on process exit to flush buffered spans/metrics/logs. `StartEngine()` registers a `defer` shutdown automatically.

### Local dev with Docker (OTel Collector + Jaeger)

```yaml
# docker-compose.yml excerpt
services:
  otel-collector:
    image: otel/opentelemetry-collector-contrib:latest
    ports:
      - "4317:4317"   # gRPC
      - "4318:4318"   # HTTP
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686" # UI
```

```yaml
# service application.yml for local dev
otel:
  enabled: true
  service-name: tag-api
  protocol: http
  endpoint: localhost:4318
  insecure: true
```
