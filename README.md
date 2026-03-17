# Coffee Journal

A personal app for logging coffee beans and tracking tasting notes over time. Record beans once per bag, then log multiple tasting sessions per bean with flavor profiles, brew method, scores, and free-form notes. Browse everything in a timeline, search across beans and notes, and view per-bean average scores.

## Stack

| Layer | Technology |
|---|---|
| Frontend | Next.js 16 (App Router), TypeScript, Tailwind CSS |
| Backend | Go, Echo v4, PostgreSQL 16 |
| DB access | pgx/v5 + sqlc (generated, type-safe queries) |
| Migrations | golang-migrate (plain SQL files) |
| Photo storage | Cloudflare R2 |
| Logging | zerolog (structured JSON) |
| Auth | Email/password + OAuth (Google, GitHub) — deferred |

## Local Development

**Prerequisites:** Docker, Go, Node.js

```bash
cp .env.example .env
make dev-db       # Start PostgreSQL (Docker)
make migrate      # Run database migrations
make dev          # Start Go API (port 8080) + Next.js (port 3000)
```

Next.js proxies all `/api/*` requests to the Go backend, so no CORS issues when calling the API from client components.

## API Endpoints

### Beans
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/beans` | List beans |
| POST | `/api/beans` | Create bean |
| GET | `/api/beans/{id}` | Bean detail with average scores and tastings |
| PUT | `/api/beans/{id}` | Update bean |
| DELETE | `/api/beans/{id}` | Soft delete |
| POST | `/api/beans/{id}/photo` | Upload package photo to R2 |

### Tastings
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/tastings` | Timeline (all tastings, newest first, paginated) |
| POST | `/api/beans/{id}/tastings` | Create tasting for a bean |
| GET | `/api/tastings/{id}` | Single tasting |
| PUT | `/api/tastings/{id}` | Update tasting |
| DELETE | `/api/tastings/{id}` | Soft delete |

### Other
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/search?q=` | Full-text search across beans and tasting notes |
| GET | `/health` | Health check |

## Architecture

The Go backend follows Clean Architecture — dependencies flow inward only:

```
Handler → Service → Repository → Domain
```

- **`domain/`** — Entities and repository interfaces. No external dependencies.
- **`internal/service/`** — Business logic and validation.
- **`internal/repository/`** — PostgreSQL implementations via sqlc.
- **`internal/http/`** — Echo HTTP handlers (thin adapters).
- **`builder/`** — Manual dependency injection.

The frontend uses Next.js App Router. Pages that benefit from SEO (timeline, bean detail, search results) are Server Components. Forms and the search bar are Client Components.

## Database

All primary keys are UUIDs. Soft deletes via `deleted_at`. Full-text search uses PostgreSQL `tsvector` generated columns with GIN indexes on bean name/roaster/origin and tasting notes.

See `PLAN.md` for the full schema, and `doc/adr/` for the decisions behind each technology choice.

## Deployment

| Service | Platform |
|---|---|
| Frontend | Vercel (auto-deploy from `main`) |
| API | Fly.io |
| Database | Fly.io Postgres |

Migrations run automatically on deploy via Fly.io `release_command`.
