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
.PHONY: dev dev-db dev-api dev-web migrate migrate-down sqlc test

dev:     overmind start
dev-db:  docker compose up -d
dev-api: dev-db
	cd api && air
dev-web: cd web && npm run dev
migrate: cd api && go run main.go migrate up
migrate-down: cd api && go run main.go migrate down
sqlc:    cd api && sqlc generate
test:    (cd api && go test ./...) && (cd web && npm run build)
```

### Key decisions

- **`make dev` uses overmind**: gives each service a labeled tmux pane, eliminating interleaved output. Requires `overmind` + `tmux` (`brew install overmind`).
- **Individual `dev-*` targets**: still available for running services in isolation without overmind.
- **`dev-api` depends on `dev-db`**: `docker compose up -d` is idempotent — safe to run if Postgres is already up. Prevents a confusing startup failure if someone runs `make dev-api` directly. `dev-web` does not depend on `dev-db` since Next.js talks to the Go API, not the DB directly.
- **`.PHONY` on all targets**: none produce files, so without `.PHONY`, a file named `test` or `dev` would silently prevent the target from running.
- **`test` uses subshells**: `(cd api && go test ./...) && (cd web && npm run build)` — each `cd` is isolated so `cd web` always resolves from the repo root, not from inside `api/`.
- **`migrate-down` added**: useful during Phase 2 when iterating on migration files. Mirrors `migrate up` via `urfave/cli`.
- **`migrate` uses `go run main.go migrate up`**: avoids requiring a compiled binary during scaffolding. Production uses the compiled binary via Fly.io's release command.
- **Procfile omits `-d` for `db` intentionally**: overmind requires foreground processes to capture per-pane output. The `dev-db` target keeps `-d` for running the DB in isolation without overmind.
- **`test` uses `npm run build` not `npm test`**: acts as a TypeScript/import error check since no test suite exists yet in Phase 1.

## Out of Scope

All other Phase 1 tasks (directory structure, `go mod init`, `docker-compose.yml`, `.env.example`, Next.js init, OpenAPI spec) are separate issues.
