# ADR-008: Clean Architecture and Project Structure

## Status
Accepted

## Date
2026-03-16

## Context
The Go backend needs a clear, consistent architectural pattern that separates concerns, enforces dependency direction, and scales cleanly as more domain entities are added (Bean, Tasting, User). The reference project (go-clean-starter) provides a proven template for this.

## Decision
Use **Clean Architecture** with four concentric layers. Dependencies flow inward only — outer layers depend on inner layers, never the reverse.

Also adopt:
- **Manual dependency injection** via a `builder/` package
- **`urfave/cli/v3`** for CLI commands (serve, migrate)
- **`oapi-codegen`** for OpenAPI spec-first handler types

## Architecture Layers

```
┌─────────────────────────────────────────┐
│   HTTP Handler (Echo)                   │  ← Outermost: HTTP, request/response
├─────────────────────────────────────────┤
│   Service / Usecase                     │  ← Business logic, validation
├─────────────────────────────────────────┤
│   Repository (Data Access)              │  ← sqlc queries, pgx, type conversion
├─────────────────────────────────────────┤
│   Domain (Entities + Interfaces)        │  ← Innermost: pure Go, no dependencies
└─────────────────────────────────────────┘
```

**Dependency rule:** Domain knows nothing about the outside world. Services depend on domain interfaces. Repositories implement domain interfaces. Handlers depend on service interfaces.

## Layer Responsibilities

### Domain (`domain/`)
- Domain entities (Bean, Tasting, User) with private fields and getters
- Repository interfaces (BeanRepository, TastingRepository, UserRepository)
- Factory constructors that enforce invariants (e.g. `NewBean(...)`)
- No imports from any other internal package

### Service (`internal/service/`)
- One package per domain entity (bean/, tasting/, user/)
- Each contains: usecase interface, implementation, input/output types, domain errors
- Validates input before calling repositories
- Orchestrates cross-aggregate operations
- Depends on domain interfaces only

### Repository (`internal/repository/`)
- Implements domain repository interfaces
- Calls sqlc-generated query functions
- Converts between pgtype and domain types
- Handles transaction context propagation via `BaseRepository`

### Handler (`internal/http/handler/`)
- Thin HTTP adapter — parses request, calls service, writes response
- One package per domain entity (bean/, tasting/, user/, search/)
- Maps domain errors to HTTP status codes via `base.HandleError()`
- Converts request → service input, service output → response

## Project Structure

```
api/
├── main.go                          # CLI entrypoint
├── sqlc.yaml                        # sqlc configuration
├── builder/
│   ├── builder.go                   # Wire handlers, services, repositories
│   └── dependency.go                # DB pool, config, shared deps
├── cmd/
│   ├── serve.go                     # Start HTTP server (urfave/cli command)
│   └── migration.go                 # Run migrations (urfave/cli command)
├── config/
│   └── config.go                    # Load env vars into Config struct
├── domain/
│   ├── bean.go                      # Bean entity + BeanRepository interface
│   ├── tasting.go                   # Tasting entity + TastingRepository interface
│   └── user.go                      # User entity + UserRepository interface
├── internal/
│   ├── http/
│   │   ├── server.go                # Echo setup, middleware, route registration
│   │   ├── base/
│   │   │   └── base.go              # ResponseRoot, ErrorResponse, HandleError
│   │   ├── handler/
│   │   │   ├── errors.go            # AppError types + HTTP mapping
│   │   │   ├── converter.go         # Request→Input, Output→Response helpers
│   │   │   ├── bean/
│   │   │   │   └── handler.go       # BeanHandler (list, get, create, update, delete, photo)
│   │   │   ├── tasting/
│   │   │   │   └── handler.go       # TastingHandler (timeline, get, create, update, delete)
│   │   │   ├── user/
│   │   │   │   └── handler.go       # UserHandler (register, login, OAuth)
│   │   │   └── search/
│   │   │       └── handler.go       # SearchHandler
│   │   └── middleware/
│   │       └── middleware.go        # Recover, BodyDump, DefaultContentType
│   ├── repository/
│   │   ├── base.go                  # BaseRepository (pgxpool wrapper + GetQueries)
│   │   ├── transaction.go           # Transaction interface + implementation
│   │   ├── common/
│   │   │   └── common.go            # pgtype↔domain type conversion helpers
│   │   ├── bean/
│   │   │   └── bean.go              # BeanRepository implementation
│   │   ├── tasting/
│   │   │   └── tasting.go           # TastingRepository implementation
│   │   └── user/
│   │       └── user.go              # UserRepository implementation
│   ├── service/
│   │   ├── bean/
│   │   │   ├── usecase.go           # BeanUsecase interface
│   │   │   ├── bean.go              # Implementation
│   │   │   ├── input.go             # CreateBeanInput, UpdateBeanInput
│   │   │   ├── output.go            # BeanOutput
│   │   │   └── error.go             # Domain errors
│   │   ├── tasting/
│   │   │   ├── usecase.go
│   │   │   ├── tasting.go
│   │   │   ├── input.go
│   │   │   ├── output.go
│   │   │   └── error.go
│   │   └── user/
│   │       ├── usecase.go
│   │       ├── user.go
│   │       ├── input.go
│   │       ├── output.go
│   │       └── error.go
│   └── sqlc/
│       ├── db.go                    # DBTX interface (pool or transaction)
│       ├── models.go                # Generated DB structs
│       ├── querier.go               # Generated Querier interface
│       ├── beans.sql.go             # Generated bean queries
│       ├── tastings.sql.go          # Generated tasting queries
│       ├── users.sql.go             # Generated user queries
│       └── query/
│           ├── beans.sql            # Bean SQL queries (sqlc annotated)
│           ├── tastings.sql         # Tasting SQL queries
│           └── users.sql            # User SQL queries
├── migration/
│   ├── migrate.go                   # golang-migrate wrapper (embed sql/)
│   └── sql/
│       ├── 000000_init.up.sql       # update_updated_at trigger function
│       ├── 000000_init.down.sql
│       ├── 000001_create_tables.up.sql
│       └── 000001_create_tables.down.sql
└── pkg/                             # Shared utilities (if needed)
```

## Manual Dependency Injection

Dependencies are wired explicitly in `builder/` — no DI framework, no code generation.

```go
// builder/builder.go
func InitializeBeanHandler(d *Dependency) *bean.BeanHandler {
    beanRepo := beanRepo.NewBeanRepository(d.DB)
    beanService := beanService.NewBeanUsecase(beanRepo)
    return bean.NewBeanHandler(beanService)
}
```

**Why manual DI:**
- No magic — every dependency is traceable by reading the code
- No build-time code generation (Wire) or reflection (Fx)
- Idiomatic Go — explicit is better than implicit
- Easy to debug — the full wiring is visible in one file

## OpenAPI + oapi-codegen

The API contract is defined in `doc/api.yaml` (OpenAPI 3.0). `oapi-codegen` generates Go types for request/response structs. Handlers use these generated types for input binding and output serialisation.

**Why:**
- API contract is defined once and shared with the frontend
- Generated types are compile-time safe — no manual struct duplication
- Consistent with the reference template

## CLI (urfave/cli/v3)

`main.go` uses `urfave/cli/v3` to expose commands:

```
./api serve      # Start HTTP server
./api migrate    # Run database migrations
```

This keeps the binary self-contained — no separate migration binary needed.

## Consequences

### Positive
- Dependency rule enforced by package structure — impossible to accidentally import the wrong layer
- Each layer is independently testable — services can be tested with mock repositories, handlers with mock services
- Adding a new entity (e.g. Equipment) follows a clear pattern: domain → service → repository → handler → builder → routes
- Manual DI keeps all wiring explicit and readable

### Negative
- More files and packages than a flat structure — initial setup is slower
- Every new entity requires touching multiple layers (by design, not by accident)
- oapi-codegen adds a code generation step to the development workflow
