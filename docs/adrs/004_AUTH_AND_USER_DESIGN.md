# ADR-004: Auth and User Design

## Status
Accepted (schema now, implementation deferred to a later phase)

## Date
2026-03-14

## Context
The original plan was single-user with no auth. After review, multi-user support was chosen for scalability — the app is designed to go public in the future. Both email/password and OAuth (Google, GitHub) are required auth mechanisms.

Auth implementation is deferred to avoid blocking early progress on core backend fundamentals (API design, SQL, deployment). The schema is designed correctly from the start so that adding auth later does not require breaking migrations.

## Decision
- Add a `users` table and `oauth_accounts` table to the schema now
- Add `user_id` to both `beans` and `tastings` now
- Implement auth endpoints, middleware, and OAuth flow in a later phase

## Domain Model

### Ubiquitous Language

| Term | Meaning |
|------|---------|
| User | A registered person who owns beans and tastings |
| Credential | An email/password pair used to authenticate a user |
| OAuthAccount | A linked third-party identity (Google or GitHub) that authenticates a user |
| Session | A period of authenticated access — represented by a JWT |

### Entity

**User** — has identity, persists over time, owns beans and tastings.

### Aggregate

**User** is an aggregate root.
- Owns its credentials (password hash) and OAuth account links
- Does not own Bean or Tasting objects — those are separate aggregate roots that reference User by ID

### Ownership Model

Both `beans` and `tastings` carry a `user_id` directly. This is intentional:

- **Tasting is an independent aggregate root** (established in ADR-003) — it should own its user reference rather than inheriting it through Bean
- **Timeline queries** load tastings directly without joining through beans — a direct `user_id` on tastings keeps that query simple
- The denormalization is justified by the access pattern

A user's beans: `WHERE beans.user_id = $1`
A user's timeline: `WHERE tastings.user_id = $1`
A bean's tastings: `WHERE tastings.bean_id = $1` (no user join needed — bean_id implies ownership)

### Repository

**UserRepository** — create user, find by email, find by OAuth provider+ID, upsert OAuth account.

## Database Schema

```sql
CREATE TABLE users (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email          TEXT NOT NULL UNIQUE,
  password_hash  TEXT,          -- nullable: OAuth-only users have no password
  name           TEXT NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at     TIMESTAMPTZ    -- soft delete: set on deactivation, cascades to beans/tastings
);

CREATE TABLE oauth_accounts (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider         TEXT NOT NULL CHECK (provider IN ('google', 'github')),
  provider_user_id TEXT NOT NULL,
  UNIQUE (provider, provider_user_id)
);
```

`beans` and `tastings` both gain:
```sql
user_id BIGINT NOT NULL REFERENCES users(id)
```

### Indexes

```sql
CREATE INDEX users_email_idx ON users (email);
CREATE INDEX oauth_accounts_user_id_idx ON oauth_accounts (user_id);
```

## Decisions

### Separate oauth_accounts table
OAuth accounts are stored in a separate table rather than columns on `users`. This allows one user to link multiple OAuth providers without schema changes. A user can have both a password and OAuth accounts simultaneously.

### Nullable password_hash
OAuth-only users never set a password. `password_hash` is nullable to support this. Email/password users always have a non-null `password_hash` — enforced at the application layer, not the DB.

### JWT over server-side sessions
Auth implementation will use JWTs (stateless). This fits the architecture — Next.js on Vercel and Go API on Fly.io are separate services. A shared session store would require additional infrastructure (Redis). JWTs let the Go API verify identity without any shared state.

### Deferred implementation
Auth endpoints, password hashing (bcrypt), OAuth flow, and JWT middleware are deferred to a later phase. The schema is in place so no breaking migrations are needed when auth is implemented.

## Implementation Plan (deferred)

When auth is implemented, the following are needed:

**Go API endpoints:**
- `POST /api/auth/register` — email/password registration
- `POST /api/auth/login` — email/password login, returns JWT
- `GET /api/auth/oauth/{provider}` — redirect to OAuth provider
- `GET /api/auth/oauth/{provider}/callback` — handle callback, upsert user, return JWT
- `POST /api/auth/logout` — client-side JWT discard (stateless, no server action needed)

**Go middleware:**
- JWT validation middleware on all protected routes
- Inject `user_id` into request context

**Next.js:**
- Store JWT in `httpOnly` cookie (set by Go API on login)
- Middleware to redirect unauthenticated users to login page
- Login and register pages (client components)

**Packages (Go):**
- `golang-jwt/jwt` — JWT creation and validation
- `golang.org/x/crypto/bcrypt` — password hashing
- OAuth: `golang.org/x/oauth2` + provider packages

## Consequences

### Positive
- Schema is multi-user from day one — no breaking migrations when auth is added
- Independent `user_id` on both beans and tastings keeps queries simple
- Separate `oauth_accounts` table supports multiple providers and future providers
- Deferred implementation keeps early phases focused on core backend fundamentals

### Negative
- `user_id NOT NULL` on beans and tastings means seeding/testing requires a user record first
- JWT statelessness means tokens cannot be revoked until expiry — acceptable for MVP
- Auth implementation (OAuth flow especially) is non-trivial and should be treated as its own phase
