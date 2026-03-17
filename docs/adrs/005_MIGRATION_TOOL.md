# ADR-005: Database Migration Tool

## Status
Accepted

## Date
2026-03-16

## Context
The Go backend needs a strategy for managing database schema changes over time. Migrations must be reproducible across local development, CI, and production (Fly Postgres).

Three options were considered:
- `golang-migrate` — CLI + Go library, plain SQL files, widely used
- `goose` — similar to golang-migrate, additionally supports Go migration files
- Plain SQL + psql — manual execution via Makefile, no dependency

## Decision
Use **`golang-migrate`** with plain SQL migration files.

## Rationale

### Plain SQL enforces fundamentals
Migration files are plain `.sql` files — no DSL, no abstraction. Writing `ALTER TABLE`, `CREATE INDEX`, and `DROP COLUMN` directly reinforces the SQL fundamentals that are a core learning goal of this project.

### Industry standard
`golang-migrate` is the most widely adopted migration tool in the Go ecosystem. Familiarity with it transfers directly to real projects.

### CLI integrates cleanly with Makefile
```bash
make migrate      # run pending up migrations
make migrate-down # roll back one step
```

### First-class CI/CD support
Migrations run automatically in GitHub Actions before tests — no manual step, no drift between environments.

### Go library for programmatic use
`golang-migrate` can also be embedded in the Go binary to run migrations on startup — useful for the Fly.io deployment where running a separate migration step is less convenient.

## File Convention

```
api/internal/database/migrations/
  001_init.sql          -- up migration
  001_init.down.sql     -- down migration (rollback)
  002_add_users.sql
  002_add_users.down.sql
```

Each migration has a corresponding `.down.sql` for rollback. Files are prefixed with a zero-padded sequence number.

## Makefile Integration

```makefile
MIGRATE=migrate -path api/internal/database/migrations -database "$(DATABASE_URL)"

migrate:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

migrate-status:
	$(MIGRATE) version
```

## CI Integration

```yaml
- name: Run migrations
  run: migrate -path api/internal/database/migrations -database "$DATABASE_URL" up
  env:
    DATABASE_URL: postgres://postgres:test@localhost:5432/coffee_test?sslmode=disable
```

## Alternatives Considered

### goose
- Equally capable for this project
- Go migration files are a useful feature but not needed here — all schema changes can be expressed in plain SQL
- Smaller community than golang-migrate

### Plain SQL + psql
- No dependency, maximum transparency
- No tracking of which migrations have run — must manage state manually
- Poor CI/CD integration
- Not suitable once the project has multiple environments

## Consequences

### Positive
- Schema changes are versioned, reproducible, and reviewable in git
- Rollback is explicit and deliberate (`.down.sql` files)
- Works identically across local Docker Postgres and Fly Postgres
- Plain SQL keeps the learning focus on SQL, not tooling

### Negative
- Requires `migrate` CLI to be installed locally and in CI
- Down migrations must be written manually — easy to forget or get wrong
- No automatic schema diffing — migrations are written by hand
