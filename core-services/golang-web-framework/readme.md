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

### Local dev with Docker (`grafana/otel-lgtm`)

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
