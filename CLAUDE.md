# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Status

This project is in **planning/design phase** — all architectural decisions are documented in ADRs and `PLAN.md`, but source code has not yet been implemented. Start from Phase 1 of `PLAN.md`.

## Commands

Once scaffolded, the Makefile will provide:

```bash
make dev          # Start PostgreSQL, Go API, and Next.js in parallel
make dev-db       # Start PostgreSQL via docker-compose
make dev-api      # Start Go API with hot reload (air)
make dev-web      # Start Next.js dev server
make migrate      # Run database migrations
make sqlc         # Regenerate sqlc types from SQL queries
make test         # Run Go tests + Next.js build check
```

**Go backend** (in `api/`):
```bash
go test ./...                        # Run all tests
go test ./internal/service/bean/...  # Run a specific package's tests
go run main.go serve                 # Start API server (port 8080)
go run main.go migrate up            # Run migrations
```

**Next.js frontend** (in `web/`):
```bash
npm run dev    # Dev server (port 3000)
npm run build  # Production build
npm run lint   # ESLint
```

## Architecture

### Stack
- **Frontend**: Next.js 16 (App Router), TypeScript, Tailwind CSS
- **Backend**: Go with Echo v4, PostgreSQL 16, pgx/v5, sqlc, golang-migrate, zerolog
- **Storage**: Cloudflare R2 for photos (S3-compatible)
- **Infrastructure**: Docker Compose (local DB), Vercel (frontend), Fly.io (API)

### Request Flow
```
Browser → Next.js (port 3000) → [API rewrite /api/*] → Go API (port 8080) → PostgreSQL (port 5433)
```

`next.config.ts` proxies all `/api/*` requests to the Go backend to avoid CORS.

### Go Backend: Clean Architecture (4 layers, inward dependencies only)

```
Handler → Service → Repository → Domain
```

- **`domain/`** — Pure Go entities (`Bean`, `Tasting`, `User`) and repository interfaces. No external imports.
- **`internal/service/`** — Business logic. Each domain has `usecase.go` (interface), `bean.go` (impl), `input.go`, `output.go`, `error.go`.
- **`internal/repository/`** — Implements domain interfaces using sqlc-generated queries. `common/` handles pgtype↔domain conversions.
- **`internal/http/`** — Echo handlers (thin adapters), middleware, and `converter.go` for request↔input and output↔response mapping.
- **`builder/`** — Manual dependency injection; wires all layers together at startup.
- **`internal/sqlc/`** — Generated code (do not edit manually). Source SQL lives in `internal/sqlc/query/*.sql`.
- **`migration/sql/`** — Plain SQL migration files, embedded into the binary via `golang-migrate`.

### Next.js Frontend: App Router

- Server Components for all list/detail pages (SSR for SEO and performance).
- Client Components only for interactive forms and search.
- `lib/api.ts` — All fetch calls to the Go API.

### Domain Model

- **Bean**: One per bag of coffee. Has roaster, origin, roast level, process, photo URL, `is_public` flag.
- **Tasting**: Evaluation session for a bean. Stores flavor tags (`TEXT[]`), brew method, grind size, scores (acidity/aroma/body/sweetness/overall 1–5), and free-form notes.
- **User**: Multi-user schema is in place. Auth (JWT + OAuth) is deferred to a later phase.

### Database Design Choices

- UUIDs as primary keys.
- Soft deletes via `deleted_at` (nullable timestamp).
- Full-text search via `tsvector` columns on beans (name, roaster, origin) and tastings (notes).
- GIN index on `flavor_tags` for array containment queries.

## Key Documentation

- `PLAN.md` — Master implementation plan with all 8 phases, DB schema, API endpoints, and project structure.
- `ADR/001–008` — Architecture Decision Records explaining every major technology choice.
