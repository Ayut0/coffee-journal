# Coffee Bean & Flavor Journal — Implementation Plan

## Context
Build a daily-use Coffee Bean & Flavor Journal app. Log beans once per bag, record multiple tastings per bean with flavor profiles (multi-select tags, 1-5 scores, aftertaste), and a written note. Browse tastings in a timeline, view bean details with average scores. Multi-user with auth (deferred). Designed to go public in the future (SEO, social previews).

## Tech Stack
- **Frontend**: Next.js 16 (App Router) + TypeScript + Tailwind CSS — see [ADR-001](doc/adr/001_NEXTJS_OVER_REACT.md)
- **Backend**: Go + Echo v4 + `pgx/v5` + `sqlc` (Clean Architecture) — see [ADR-002](doc/adr/002_BACKEND_LANGUAGE_CHOICE.md), [ADR-006](doc/adr/006_DATABASE_DRIVER.md), [ADR-008](doc/adr/008_CLEAN_ARCHITECTURE.md)
- **Database**: PostgreSQL with UUID PKs, full-text search via `tsvector` — see [ADR-003](doc/adr/003_DATA_DESIGN.md)
- **Migrations**: `golang-migrate` with plain SQL files — see [ADR-005](doc/adr/005_MIGRATION_TOOL.md)
- **Logging**: `zerolog` (structured JSON) — see [ADR-007](doc/adr/007_ROUTER_AND_LOGGING.md)
- **Auth**: email/password + OAuth (Google, GitHub), JWT — deferred — see [ADR-004](doc/adr/004_AUTH_AND_USER_DESIGN.md)
- **Object storage**: Cloudflare R2 (package photos)
- **Location**: `~/projects/coffee-journal`

## Architecture Overview
```
┌──────────────────┐       ┌──────────────────┐
│   Next.js (web)  │──────▶│   Go API (api)   │
│   Port 3000      │ fetch │   Port 8080      │
│   SSR + React    │◀──────│   JSON REST      │
└──────────────────┘       └────────┬─────────┘
                                    │
                           ┌────────┴─────────┐
                           │   PostgreSQL      │
                           │   Port 5433       │
                           └──────────────────┘
```

- **Next.js** handles SSR (bean pages server-rendered for SEO/social previews), routing, and UI
- **Go** handles all data logic, photo uploads to R2, and search
- They communicate via HTTP (Next.js server components fetch from Go API)

## Go Backend — Clean Architecture

```
Handler → Service → Repository → Domain
```

Dependencies flow inward only. See [ADR-008](doc/adr/008_CLEAN_ARCHITECTURE.md) for full rationale.

## Project Structure
```
coffee-journal/
├── Makefile
├── docker-compose.yml
├── .env.example
│
├── api/                                    # Go backend
│   ├── go.mod
│   ├── main.go                             # CLI entrypoint (urfave/cli)
│   ├── sqlc.yaml                           # sqlc configuration
│   ├── builder/
│   │   ├── builder.go                      # Wire handlers, services, repositories
│   │   └── dependency.go                   # DB pool, config, shared deps
│   ├── cmd/
│   │   ├── serve.go                        # HTTP server command
│   │   └── migration.go                    # Migration command
│   ├── config/
│   │   └── config.go                       # Load env vars into Config struct
│   ├── domain/
│   │   ├── bean.go                         # Bean entity + BeanRepository interface
│   │   ├── tasting.go                      # Tasting entity + TastingRepository interface
│   │   └── user.go                         # User entity + UserRepository interface
│   ├── internal/
│   │   ├── http/
│   │   │   ├── server.go                   # Echo setup, middleware, route registration
│   │   │   ├── base/
│   │   │   │   └── base.go                 # ResponseRoot, ErrorResponse, HandleError
│   │   │   ├── handler/
│   │   │   │   ├── errors.go               # AppError types + HTTP mapping
│   │   │   │   ├── converter.go            # Request→Input, Output→Response helpers
│   │   │   │   ├── bean/
│   │   │   │   │   └── handler.go          # List, get, create, update, delete, photo
│   │   │   │   ├── tasting/
│   │   │   │   │   └── handler.go          # Timeline, get, create, update, delete
│   │   │   │   ├── user/
│   │   │   │   │   └── handler.go          # Register, login, OAuth (deferred)
│   │   │   │   └── search/
│   │   │   │       └── handler.go          # Full-text search
│   │   │   └── middleware/
│   │   │       └── middleware.go           # Recover, BodyDump, DefaultContentType
│   │   ├── repository/
│   │   │   ├── base.go                     # BaseRepository (pgxpool + GetQueries)
│   │   │   ├── transaction.go              # Transaction interface + implementation
│   │   │   ├── common/
│   │   │   │   └── common.go               # pgtype↔domain type conversion helpers
│   │   │   ├── bean/
│   │   │   │   └── bean.go                 # BeanRepository implementation
│   │   │   ├── tasting/
│   │   │   │   └── tasting.go              # TastingRepository implementation
│   │   │   └── user/
│   │   │       └── user.go                 # UserRepository implementation
│   │   ├── service/
│   │   │   ├── bean/
│   │   │   │   ├── usecase.go              # BeanUsecase interface
│   │   │   │   ├── bean.go                 # Implementation
│   │   │   │   ├── input.go                # CreateBeanInput, UpdateBeanInput
│   │   │   │   ├── output.go               # BeanOutput
│   │   │   │   └── error.go                # Domain errors
│   │   │   ├── tasting/
│   │   │   │   ├── usecase.go
│   │   │   │   ├── tasting.go
│   │   │   │   ├── input.go
│   │   │   │   ├── output.go
│   │   │   │   └── error.go
│   │   │   └── user/
│   │   │       ├── usecase.go
│   │   │       ├── user.go
│   │   │       ├── input.go
│   │   │       ├── output.go
│   │   │       └── error.go
│   │   └── sqlc/
│   │       ├── db.go                       # DBTX interface (pool or transaction)
│   │       ├── models.go                   # Generated DB structs
│   │       ├── querier.go                  # Generated Querier interface
│   │       ├── beans.sql.go                # Generated bean queries
│   │       ├── tastings.sql.go             # Generated tasting queries
│   │       ├── users.sql.go                # Generated user queries
│   │       └── query/
│   │           ├── beans.sql               # Bean SQL (sqlc annotated)
│   │           ├── tastings.sql            # Tasting SQL
│   │           └── users.sql               # User SQL
│   ├── migration/
│   │   ├── migrate.go                      # golang-migrate wrapper (//go:embed)
│   │   └── sql/
│   │       ├── 000000_init.up.sql          # update_updated_at trigger
│   │       ├── 000000_init.down.sql
│   │       ├── 000001_create_tables.up.sql
│   │       └── 000001_create_tables.down.sql
│   └── doc/
│       └── api.yaml                        # OpenAPI 3.0 spec (oapi-codegen source)
│
└── web/                                    # Next.js frontend
    ├── package.json
    ├── next.config.ts
    ├── tailwind.config.ts
    ├── app/
    │   ├── layout.tsx                      # Root layout with nav
    │   ├── page.tsx                        # Home: timeline (SSR)
    │   ├── beans/
    │   │   ├── page.tsx                    # Beans list (SSR)
    │   │   ├── new/page.tsx                # Create bean (client)
    │   │   └── [id]/
    │   │       ├── page.tsx                # Bean detail (SSR — SEO/social)
    │   │       ├── edit/page.tsx           # Edit bean (client)
    │   │       └── tastings/
    │   │           └── new/page.tsx        # New tasting (client)
    │   └── search/page.tsx                 # Search results (SSR)
    ├── components/
    │   ├── BeanCard.tsx
    │   ├── TastingCard.tsx
    │   ├── FlavorTagSelect.tsx             # Multi-select chip picker
    │   ├── ScoreSlider.tsx                 # Reusable 1-5 input
    │   └── SearchBar.tsx                   # "use client" — debounced input
    ├── lib/
    │   └── api.ts                          # Fetch helpers for Go API
    ├── hooks/
    │   └── useDebounce.ts
    └── types/
        └── index.ts
```

### Key Next.js Architecture Decisions
- **Server Components** for pages that need SEO: timeline, bean detail, beans list, search
- **Client Components** (`"use client"`) for interactive parts: forms, search bar
- **Server-side fetch**: Server components call Go API directly (e.g., `fetch('http://localhost:8080/api/beans')`)
- **`next.config.ts` rewrites**: Proxy `/api/*` to Go backend so client components can call `/api/...` without CORS issues

## Database Schema

See [ADR-003](doc/adr/003_DATA_DESIGN.md) for the full domain model, aggregate design, and rationale.
All primary keys are **UUID** (`gen_random_uuid()`). See [ADR-008](doc/adr/008_CLEAN_ARCHITECTURE.md).

**`users`** aggregate root — see [ADR-004](doc/adr/004_AUTH_AND_USER_DESIGN.md):
- `id` (UUID PK), `email` (UNIQUE), `password_hash` (nullable), `name`
- `created_at`, `updated_at`, `deleted_at` (soft delete)
- Linked `oauth_accounts` table: `user_id`, `provider` (google/github), `provider_user_id`

**`beans`** aggregate root:
- `id` (UUID PK), `user_id` (FK → users), `name`, `roaster`, `origin`, `roast_level` (Light/Medium/Dark), `process` (Washed/Natural/Honey, nullable)
- `altitude_min`, `altitude_max` (INT, nullable — MASL), `harvest_season` (TEXT, nullable — e.g. `"2023/24"`)
- `package_photo_url` (TEXT, nullable — public R2 URL), `is_public` (BOOLEAN, default false)
- `created_at`, `updated_at`, `deleted_at` (soft delete)

**`tastings`** aggregate root:
- `id` (UUID PK), `user_id` (FK → users), `bean_id` (FK → beans, CASCADE)
- `flavor_tags` (TEXT[]), `brew_method`, `grind_size`
- `acidity`, `aroma`, `body` (SMALLINT 1-5), `sweetness` (nullable), `overall` (SMALLINT 1-5)
- `aftertaste` (Short/Clean/Lingering), `note_text` (nullable)
- `created_at`, `updated_at`, `deleted_at` (soft delete)

### Full-Text Search Index (tsvector)
`beans` and `tastings` each carry a `search_vector TSVECTOR GENERATED ALWAYS AS (...) STORED` column, kept up-to-date automatically by PostgreSQL:
```sql
-- beans: index name + roaster + origin
ALTER TABLE beans ADD COLUMN search_vector tsvector
  GENERATED ALWAYS AS (
    to_tsvector('english', coalesce(name,'') || ' ' || coalesce(roaster,'') || ' ' || coalesce(origin,''))
  ) STORED;
CREATE INDEX beans_search_vector_idx ON beans USING GIN (search_vector);

-- tastings: index note_text
ALTER TABLE tastings ADD COLUMN search_vector tsvector
  GENERATED ALWAYS AS (
    to_tsvector('english', coalesce(note_text,''))
  ) STORED;
CREATE INDEX tastings_search_vector_idx ON tastings USING GIN (search_vector);
```
Search query: `WHERE search_vector @@ plainto_tsquery('english', $1)`.

## API Endpoints (Go — `/api`)

### Beans
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/beans` | List user's beans |
| POST | `/api/beans` | Create bean |
| GET | `/api/beans/{id}` | Bean detail + average scores + tastings |
| PUT | `/api/beans/{id}` | Update bean |
| DELETE | `/api/beans/{id}` | Soft delete bean |
| POST | `/api/beans/{id}/photo` | Upload package photo to R2 |

### Tastings
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/tastings` | Timeline: all tastings, newest first (paginated) |
| POST | `/api/beans/{id}/tastings` | Create tasting for a bean |
| GET | `/api/tastings/{id}` | Get single tasting |
| PUT | `/api/tastings/{id}` | Update tasting |
| DELETE | `/api/tastings/{id}` | Soft delete tasting |

### Search
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/search?q=` | Full-text search across beans + tasting notes |

### Health
| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check — returns 200 OK. Required by Fly.io for zero-downtime rolling deploys. |

## Frontend Pages

| Route | Rendering | Description |
|-------|-----------|-------------|
| `/` | SSR | Timeline — recent tastings with bean name, flavor chips, scores |
| `/beans` | SSR | Beans grid — BeanCards with name, roaster, roast level, tasting count |
| `/beans/new` | Client | Create bean form |
| `/beans/[id]` | SSR | Bean detail — info, average scores, tastings list. **Server-rendered for SEO + social sharing** |
| `/beans/[id]/edit` | Client | Edit bean form |
| `/beans/[id]/tastings/new` | Client | Tasting form — flavor tags, scores, brew method, grind size, aftertaste, note |
| `/search` | SSR | Search results across beans + notes |

### Social Preview (Open Graph)
```tsx
// app/beans/[id]/page.tsx
export async function generateMetadata({ params }) {
  const data = await fetch(`${API_URL}/api/beans/${params.id}`).then(r => r.json());
  return {
    title: `${data.bean.name} — Coffee Journal`,
    description: `${data.bean.roaster} · ${data.bean.origin} · ${data.averages.tasting_count} tastings`,
  };
}
```

## Implementation Phases

### Phase 1: Scaffolding
- Create `coffee-journal/` with `api/` and `web/` dirs, git init
- `api/`: `go mod init`, clean architecture directory structure, `sqlc.yaml`
- `web/`: `npx create-next-app@latest` with TypeScript + Tailwind + App Router
- `docker-compose.yml` (PostgreSQL), `Makefile`, `.env.example`, `.gitignore`
- Write OpenAPI spec (`doc/api.yaml`) for all endpoints

### Phase 2: Database + Go Foundation
- Migration files (users, beans, tastings, audit_log tables + indexes)
- `config/`, `builder/dependency.go` (DB pool setup)
- Domain entities (bean.go, tasting.go, user.go) with repository interfaces
- sqlc query files + `sqlc generate`

### Phase 3: Bean + Tasting API
- Repository implementations (bean, tasting)
- Service layer (bean, tasting) with input validation and domain errors
- Echo handlers (bean, tasting) with error mapping
- Wire everything in `builder/builder.go`
- Test all endpoints with curl

### Phase 4: Search + Photo Upload
- Search handler (full-text across beans + notes)
- R2 client setup, photo upload endpoint

### Phase 5: Next.js Foundation + Read Pages
- `next.config.ts` with API rewrites, `lib/api.ts` fetch helpers, TypeScript types
- Root layout with nav
- TimelinePage (SSR) + BeansListPage (SSR) + BeanDetailPage (SSR + Open Graph)
- BeanCard + TastingCard components

### Phase 6: Forms & CRUD (Client Components)
- BeanFormPage (create + edit)
- TastingFormPage (scores, tags, brew method, grind size, aftertaste, note)
- FlavorTagSelect, ScoreSlider components
- Delete confirmations

### Phase 7: Search & Polish
- SearchBar (client, debounced) + SearchPage (SSR results)
- Loading states (`loading.tsx`), error states (`error.tsx`)
- Responsive Tailwind CSS (mobile-friendly for cafe use)

### Phase 8: Auth (deferred)
- User registration + login (email/password)
- OAuth flow (Google, GitHub)
- JWT middleware on all protected routes
- Next.js auth middleware (redirect unauthenticated users)

## Key Packages

**Go**: `labstack/echo/v4`, `jackc/pgx/v5`, `sqlc-dev/sqlc`, `golang-migrate/migrate/v4`, `rs/zerolog`, `joho/godotenv`, `aws/aws-sdk-go-v2`, `golang.org/x/crypto`, `google/uuid`, `urfave/cli/v3`, `stretchr/testify`, `oapi-codegen/runtime`

**Next.js**: `tailwindcss`, `date-fns`, `react-hot-toast`

## Infrastructure

### Local Development
```
docker-compose.yml:
  - PostgreSQL 16 (port 5433)
  - Persistent volume for DB data

Makefile targets:
  make dev-db       → docker compose up -d
  make migrate      → ./api migrate
  make sqlc         → sqlc generate
  make dev-api      → cd api && air (Go hot reload on :8080)
  make dev-web      → cd web && npm run dev (Next.js on :3000)
  make dev          → run all three above in parallel
  make test         → go test ./... && cd web && npm test
```

- `air` for Go hot reload, Next.js built-in fast refresh for frontend
- `.env.local` in `web/` for `NEXT_PUBLIC_API_URL` and internal `API_URL`

### Deployment

```
┌──────────────────┐          ┌──────────────────┐
│   Vercel          │  fetch  │   Fly.io          │
│   (Next.js)       │────────▶│   (Go API)        │
│   Free tier       │         │   Free tier        │
│   Auto HTTPS      │         │   + Fly Postgres   │
│   Edge CDN        │         │                   │
└──────────────────┘          └──────────────────┘
```

- **Next.js on Vercel**: Free tier, auto-deploy from GitHub, edge CDN, built-in HTTPS
- **Go API on Fly.io**: Free tier, managed Postgres. No persistent volume needed — photos go to R2.

Files needed:
- **`api/Dockerfile`** — Go multi-stage build (build → alpine)
- **`api/fly.toml`** — Fly.io config (see release command below)
- **`web/vercel.json`** — (optional) Vercel config if needed

**Production migrations** — Fly.io runs a release command before starting the new instance. This ensures the DB is migrated before traffic shifts:
```toml
# api/fly.toml
[deploy]
  release_command = "./api migrate up"
```
The `migrate up` subcommand is already wired via `urfave/cli`. On each deploy Fly runs it in a temporary VM with access to `DATABASE_URL`, then starts the main process only if it exits 0.

### CI/CD (GitHub Actions)

```yaml
# .github/workflows/ci.yml
on:
  push: { branches: [main] }
  pull_request: { branches: [main] }

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env: { POSTGRES_DB: coffee_test, POSTGRES_PASSWORD: test }
        ports: [5432:5432]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.25' }
      - uses: actions/setup-node@v4
        with: { node-version: '20' }
      - run: cd api && go test ./...
      - run: cd web && npm ci && npm run build

  deploy-api:
    needs: test
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: cd api && flyctl deploy --remote-only
        env: { FLY_API_TOKEN: "${{ secrets.FLY_API_TOKEN }}" }

  # Vercel auto-deploys from GitHub — no action needed
```

### Infrastructure Summary

| Layer | Local Dev | Production |
|-------|-----------|------------|
| Frontend | Next.js dev server (:3000) | Vercel (auto-deploy, CDN, HTTPS) |
| Backend | Go + air (:8080) | Fly.io (Go binary) |
| Database | Docker PostgreSQL (:5433) | Fly Postgres (managed) |
| Photo storage | R2 bucket (dev credentials) | Cloudflare R2 (production bucket) |
| CI/CD | — | GitHub Actions (tests) + Vercel/Fly auto-deploy |
| HTTPS | — | Both Vercel and Fly.io provide auto-TLS |

## Verification
1. `make dev-db` → PostgreSQL running
2. `make migrate` → Tables + indexes created
3. `make dev-api` + `make dev-web` → Both servers running
4. Create a bean → appears in beans list (SSR)
5. Add a tasting with flavor tags + scores → appears in timeline
6. View bean detail → server-rendered, shows average scores
7. Share bean URL → Open Graph preview works (title + description)
8. Search → finds beans and tasting notes
9. Push to main → tests pass → API deploys to Fly.io, frontend deploys to Vercel

## Deferred / To Revisit

### Pagination — cursor-based (decided)
`GET /api/tastings` uses cursor-based pagination (`?cursor=<last_id>&limit=20`), not offset.
- Stable results when new tastings are inserted mid-browse
- Natural fit for infinite scroll / timeline UX
- sqlc query: `WHERE id < $cursor ORDER BY id DESC LIMIT $limit`
- Response includes `next_cursor` (last ID in the page), empty string when no more pages

### Secrets management
Before first production deploy, set all secrets explicitly:
```bash
# Fly.io
fly secrets set DATABASE_URL=... JWT_SECRET=... R2_ACCESS_KEY_ID=... R2_SECRET_ACCESS_KEY=... R2_BUCKET=... R2_ENDPOINT=...

# Vercel
# Set via dashboard or CLI: NEXT_PUBLIC_API_URL, API_URL
```
Document required vars in `.env.example`. Revisit during Phase 1 scaffolding.

### CORS
Currently not needed — Next.js rewrites proxy all `/api/*` calls server-side, avoiding browser CORS. Revisit in Phase 8 (auth) when OAuth redirects are implemented; may need Echo's CORS middleware depending on redirect flow.

### R2 CORS config
Photo uploads route through Go API (`POST /api/beans/{id}/photo`), so the browser never talks to R2 directly. No R2 CORS config needed unless a future decision switches to client-side presigned URL uploads. Revisit if upload strategy changes.
