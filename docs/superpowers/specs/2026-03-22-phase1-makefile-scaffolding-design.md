# Phase 1 Scaffolding — Makefile & Procfile Design

**Date**: 2026-03-22
**Issues**: #2, #3
**Phase**: Phase 1 — Repo & Tooling

## Scope

Two issues from Phase 1:

- **Issue #2**: `tasks/` and `tasks/lessons.md` already exist — close as done, no implementation needed.
- **Issue #3**: Create `Makefile` with all targets and a `Procfile` for `make dev`.

## Design

### Procfile (repo root)

Used by `overmind` to run all three services in parallel with labeled tmux panes.

```
db:  docker compose up
api: cd api && air
web: cd web && npm run dev
```

### Makefile (repo root)

```makefile
dev:      overmind start
dev-db:   docker compose up -d
dev-api:  cd api && air
dev-web:  cd web && npm run dev
migrate:  cd api && go run main.go migrate up
sqlc:     cd api && sqlc generate
test:     cd api && go test ./... && cd web && npm run build
```

### Key decisions

- **`make dev` uses overmind**: gives each service a labeled tmux pane, eliminating interleaved output. Requires `overmind` + `tmux` (`brew install overmind`).
- **Individual `dev-*` targets**: still available for running services in isolation without overmind.
- **`migrate` uses `go run main.go migrate up`**: avoids requiring a compiled binary during scaffolding. The subcommand is `migrate up` (not bare `migrate`) per the `urfave/cli` wiring described in PLAN.md infrastructure section. Production uses the compiled binary via Fly.io's release command.
- **Procfile omits `-d` for `db` intentionally**: overmind requires processes to run in the foreground so it can capture output per pane. The `dev-db` target keeps `-d` for running the DB in isolation without overmind.
- **`test` uses `npm run build` not `npm test`**: acts as a TypeScript/import error check since no test suite exists yet in Phase 1. `npm ci` is omitted — it does a slow clean install and is inappropriate for local dev where `node_modules` already exists.

## Out of Scope

All other Phase 1 tasks (directory structure, `go mod init`, `docker-compose.yml`, `.env.example`, Next.js init, OpenAPI spec) are separate issues.
