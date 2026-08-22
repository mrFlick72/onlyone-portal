# Analytics

Serves budget-expense analytics from a read-optimised local copy of the user's expenses, kept current by consuming `budget-api`'s Kafka events.

## Language

**Projection**:
The Postgres-backed local copy of a user's expenses analytic-api answers every read from. Built and kept current by consuming `budget-api`'s `CREATE`/`UPDATE`/`DELETE` expense events off Kafka; analytic-api never calls `budget-api` on the read path, so it is eventually — not strictly — consistent with the source of truth.
_Avoid_: Cache, replica, read model

**Reindex**:
The one path where analytic-api does call `budget-api` — an explicit, per-caller recovery action that rebuilds their Projection over a year range by pulling expenses via REST. Upsert-only: it fills gaps from missed `CREATE`/`UPDATE` events but does not remove rows whose `DELETE` was missed.
_Avoid_: Resync, rebuild, backfill
