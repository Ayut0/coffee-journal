# ADR-006: Database Engine, Driver, and Query Layer

## Status
Accepted

## Date
2026-03-16

## Context
The Go backend needs a PostgreSQL driver and a strategy for executing queries. Three decisions are made together here because they are tightly coupled: the database engine, how Go connects to it, and how queries are written and results scanned.

## Decisions

### PostgreSQL as the database engine
PostgreSQL is the only engine considered. The data model relies on Postgres-specific features:
- `TEXT[]` with GIN indexes for flavor tags
- `tsvector` / `pg_trgm` for full-text search
- `TIMESTAMPTZ` for timezone-aware timestamps
- `BIGSERIAL` for auto-incrementing primary keys
- `CHECK` constraints for enum enforcement

No other database engine supports this feature set with the same maturity.

### `pgx/v5` as the driver
Use `github.com/jackc/pgx/v5` instead of `github.com/lib/pq`.

`lib/pq` is in maintenance mode — it receives security fixes only and no new features. `pgx/v5` is the modern, actively maintained PostgreSQL driver for Go. Key advantages:
- Built-in connection pool via `pgxpool` — no separate pooling library needed
- Native support for PostgreSQL types (`pgtype`) including arrays, ranges, and `JSONB`
- Better performance — less reflection, more direct wire protocol handling
- `pgxpool.Pool` is safe for concurrent use and handles reconnection automatically

### sqlc for type-safe query generation
SQL queries are written as plain `.sql` files with sqlc comment annotations. `sqlc generate` produces type-safe Go functions and structs from those queries. No ORM (GORM, bun, ent) is used.

Reasons:
- SQL is written by hand — no DSL, no hidden queries, full control over what hits the database
- Generated Go code is compile-time type-safe — query parameter mismatches are caught at build time, not runtime
- No reflection overhead at runtime — generated code is plain Go function calls
- Explicit SQL transfers directly to other tools (psql, explain plans, Datadog)
- `pgx/v5` is used as the sqlc driver — native Postgres type support (`pgtype.UUID`, `pgtype.Timestamptz`) with no additional mapping layer

## Query Pattern

**SQL file** (`internal/sqlc/query/beans.sql`):
```sql
-- name: GetBean :one
SELECT * FROM beans
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: CreateBean :one
INSERT INTO beans (id, user_id, name, roaster, origin, roast_level)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
```

**Generated Go** (`internal/sqlc/beans.sql.go`):
```go
func (q *Queries) GetBean(ctx context.Context, id pgtype.UUID) (Bean, error)
func (q *Queries) CreateBean(ctx context.Context, arg CreateBeanParams) (Bean, error)
```

**sqlc.yaml**:
```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "migration/sql"
    queries: "internal/sqlc/query"
    gen:
      go:
        package: "sqlc"
        out: "internal/sqlc"
        sql_package: "pgx/v5"
        emit_interface: true
        emit_pointers_for_null_types: true
```

## Connection Pool Configuration

```go
config, _ := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
config.MaxConns = 25
config.MinConns = 5
config.MaxConnLifetime = 5 * time.Minute
pool, _ := pgxpool.NewWithConfig(ctx, config)
```

## Consequences

### Positive
- `pgxpool` replaces the manual pool configuration previously needed with `database/sql` + `lib/pq`
- Native Postgres type support — no manual serialisation for `TEXT[]`
- SQL queries are explicit and reviewable — no hidden queries from ORM methods
- Direct skill transfer — SQL written here works in any other tool (psql, Datadog, explain plans)

### Negative
- More boilerplate scanning code compared to an ORM — mitigated by consistent patterns across store files
- No auto-generated queries — all SQL must be written and maintained by hand
