# ADR-002: Backend Language Choice

## Status
Accepted

## Date
2026-03-13

## Context
The Coffee Journal needs a separate backend API server to handle data logic, audio file storage, and search. The original plan chose Go (Chi router + database/sql + lib/pq).

However, the primary learning goal for this project is **backend fundamentals** — not a specific language. The skills that matter for backend development are:

- REST API design (resource modeling, status codes, error handling)
- Database design and SQL (schema, migrations, indexes, queries)
- Authentication and authorization patterns
- Testing strategies (unit, integration, e2e)
- Cloud operations (Docker, CI/CD, deployment, monitoring)

These skills transfer across any language. Learning a new language (Go) at the same time adds cognitive overhead that slows down learning the fundamentals.

## Decision
Use **Go with Chi router + `database/sql`**.

### Options Considered
- Option A: Go (Chi + database/sql)
- Option B: Node.js/TypeScript (Express, Hono, or Fastify)

## Option A: Go

### Pros
- Forces clear separation between frontend and backend (different language = no temptation to blur boundaries)
- Excellent for learning explicit error handling, strong typing, and concurrency patterns
- Compiles to a single binary — simple deployment
- Strong standard library for HTTP and JSON
- Highly valued in infrastructure and platform roles

### Cons
- New language to learn alongside backend fundamentals — splits focus
- Smaller ecosystem for web-specific tasks (ORMs, validation, middleware) compared to Node.js
- Different mental model from the TypeScript already used on the frontend
- Slower iteration speed while learning syntax, idioms, and tooling

## Option B: Node.js with TypeScript

### Pros
- **Same language as the frontend** — no context switching, focus stays on backend concepts
- Shared types between frontend and backend (can extract a shared `types/` package)
- Massive ecosystem for web backends (Express, Hono, Fastify, Prisma, Zod, etc.)
- Faster iteration — already familiar with TypeScript tooling (npm, tsconfig, etc.)
- One language to debug across the entire stack
- Strong job market demand for full-stack TypeScript

### Cons
- Easier to blur the frontend/backend boundary (temptation to skip the separate API)
- Runtime type safety requires extra tooling (Zod, io-ts) vs Go's compile-time checks
- Node.js single-threaded model requires understanding of async patterns
- Less exposure to a second language's idioms and trade-offs

## Analysis

The key question: **does learning Go serve the primary goal of this project?**

The primary goal is to build a working app while solidifying backend fundamentals (API design, SQL, testing, deployment). Go introduces a secondary learning goal (new language) that competes for attention.

| Factor | Go | Node.js/TS |
|--------|-----|------------|
| Time to first working endpoint | Slower (new syntax, tooling) | Faster (familiar language) |
| Focus on backend concepts | Split with language learning | Fully on backend concepts |
| Type sharing with frontend | Not possible | Possible |
| Deployment simplicity | Single binary | Needs Node runtime or bundler |
| Long-term career value | Strong for infra/platform | Strong for full-stack |
| Forces architectural discipline | Yes (different language) | Requires self-discipline |

### Framework Options (if Node.js/TS)
- **Express**: Most documented, largest ecosystem, but older patterns
- **Hono**: Lightweight, modern, works on edge runtimes, good DX
- **Fastify**: Fast, schema-based validation, good plugin system

## Rationale for Final Decision

The primary concern in the ADR — learning a new language alongside backend fundamentals — is reduced by existing API experience (Kotlin, Python). The concepts (REST handlers, migrations, connection pools) are already familiar; Go is a new syntax for known ideas, not an unknown domain.

Go better serves the stated goal of learning backend fundamentals *from infrastructure up*:
- Single binary deployment makes the infra layer transparent (no runtime, no process manager)
- `database/sql` with raw SQL keeps the database layer explicit — no ORM abstractions
- Explicit error handling (`if err != nil`) builds disciplined thinking about failure modes
- Goroutines expose a different concurrency model worth having alongside JVM/Node.js mental models
- Go's constraints (no exceptions, no inheritance) enforce clean separation between layers

TypeScript is the day-to-day language — Go provides genuine exposure to a second backend paradigm, which is more valuable given the learning goal than faster iteration speed.

## Consequences

- Expect some initial friction with Go syntax and tooling (offset by prior backend experience)
- Clean architectural separation between frontend and backend is enforced by the language boundary
- `database/sql` with raw SQL queries teaches database interaction without ORM abstraction
- Single binary + Dockerfile makes deployment straightforward to reason about
- Can leverage Go experience for infrastructure/platform roles in the future
