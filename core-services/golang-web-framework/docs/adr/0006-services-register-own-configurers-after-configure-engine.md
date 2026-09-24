---
status: accepted
---

# Services register their own configurers, after ConfigureEngine

`WebServerProvisioner` already gives cross-cutting concerns a lifecycle: `Configure` at boot, `Dispose` on graceful shutdown, all within the single `server.shutdown-timeout` budget. Until now only the three built-ins (OTel, standard middleware, OAuth2) could use it — they append themselves to the unexported `configurers` slice, and `ConfigureEngine()` configures only its own local list. budget-api's Scheduled Expense generation engine (#55) is the first service-owned background component that needs the same lifecycle (start a `gocron` scheduler, stop it on shutdown), so we added `RegisterConfigurer(c WebServerConfigurer)` (#64).

It must be called **after** `ConfigureEngine()`, and panics if called before. That's the order services actually wire things in `main.go`: the engine first, then the facades and repositories a service configurer depends on. The alternative — letting services pass configurers into `ConfigureEngine` — would force those dependencies to exist before the engine, and would also have put service configurers *before* OAuth2 in the list, which only matters for middleware ordering but would have been surprising.

`Configure` runs immediately on registration. A failure is treated exactly like a built-in's at boot: the provisioner is shut down (disposing everything registered so far, the failing configurer included) and it panics. We kept the panic, rather than returning an error, so both boot paths fail the same way and a half-configured process never keeps serving; a service calling this from `main.go` has no meaningful recovery anyway.

`Dispose` runs in registration order, so service configurers are disposed after the built-ins — after `srv.Shutdown` has drained in-flight requests, and sharing what's left of the shutdown budget with the OTel exporter flush.

The change is additive: no existing signature or behavior moves.
