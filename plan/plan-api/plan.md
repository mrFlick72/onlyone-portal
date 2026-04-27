
## Tech Stack

| Concern | Choice                                                                                                                                                                                                                                                        |
|---------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Language | Go 1.25.1                                                                                                                                                                                                                                                     |
| Web framework | Gin (`github.com/gin-gonic/gin v1.11.0`)                                                                                                                                                                                                                      |
| Persistence | Postgres                                                                                                                                                                                                                                                      |
| ID generation | google/uuid (salt for DynamoDB range keys)                                                                                                                                                                                                                    |
| Auth | JWT validation via the shared `core-services/golang-web-framework` middleware; JWKS fetched from `http://local.api.vauthenticator.com:9090/oauth2/jwks`                                                                                                       |
| Shared framework | `github.com/mrflick72/onlyone-portal/core-services/golang-web-framework` — resolved via local `replace` directive in `go.mod` pointing to `../../core-services/golang-web-framework`                                                                          |
| Test assertions | testify + go-playground/assert                                                                                                                                                                                                                                |
| Build tag | `-tags test` is required for any test that imports shared fixtures. Helpers like `domain/tags/fixture.go` are guarded by `//go:build test` so they are never linked into the production binary. Running `go test ./...` without the tag will fail to compile. |

Config is read by the shared framework's `config.GetConfigurationManagerInstance()` (backed by Viper). The config file path is set via the `CONFIG_FILE_LOCATION` env var.
