# Phase 1 Makefile & Procfile Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create the `Makefile` and `Procfile` for the coffee-journal repo, and close the already-complete Issue #2.

**Architecture:** Two files at the repo root. `Procfile` defines the three foreground processes for `overmind`. `Makefile` exposes all developer targets, with `make dev` delegating to overmind and `make dev-api` auto-starting the DB first.

**Tech Stack:** GNU Make, overmind (tmux-based process manager), Docker Compose, air (Go hot reload), Next.js dev server.

---

## Files

| Action | Path | Purpose |
|--------|------|---------|
| Create | `Procfile` | Defines `db`, `api`, `web` processes for overmind |
| Create | `Makefile` | All developer targets with correct dependencies |

---

### Task 1: Close Issue #2

Issue #2 notes that `tasks/` and `tasks/lessons.md` already exist. Nothing to implement — just close it.

- [ ] **Step 1: Verify `tasks/` exists**

```bash
ls tasks/
```

Expected output includes `lessons.md` and `todo.md`.

- [ ] **Step 2: Close the issue**

```bash
gh issue close 2 --repo Ayut0/coffee-journal --comment "tasks/ and tasks/lessons.md already exist from the planning step. Nothing to implement."
```

---

### Task 2: Create Procfile

- [ ] **Step 1: Create `Procfile` at repo root**

```
db:  docker compose up
api: cd api && air
web: cd web && npm run dev
```

> Note: `docker compose up` runs without `-d` intentionally. overmind requires foreground processes to capture per-pane output.

- [ ] **Step 2: Verify file contents**

```bash
cat Procfile
```

Expected: three lines — `db`, `api`, `web`.

---

### Task 3: Create Makefile

- [ ] **Step 1: Create `Makefile` at repo root**

```makefile
.PHONY: dev dev-db dev-api dev-web migrate migrate-down sqlc test

dev:
	overmind start

dev-db:
	docker compose up -d

dev-api: dev-db
	cd api && air

dev-web:
	cd web && npm run dev

migrate:
	cd api && go run main.go migrate up

migrate-down:
	cd api && go run main.go migrate down

sqlc:
	cd api && sqlc generate

test:
	(cd api && go test ./...) && (cd web && npm run build)
```

> **Important:** Recipe lines (the commands under each target) must be indented with a **tab character**, not spaces. Most editors convert tabs to spaces — double-check if you paste this manually.

- [ ] **Step 2: Verify syntax with a dry-run**

```bash
make -n dev
make -n test
make -n dev-api
```

Expected for `make -n dev-api`:
```
docker compose up -d
cd api && air
```
This confirms `dev-db` runs first as a prerequisite.

Expected for `make -n test`:
```
(cd api && go test ./...) && (cd web && npm run build)
```

- [ ] **Step 3: Verify `.PHONY` is recognised**

```bash
make --print-data-base | grep "^\.PHONY"
```

Expected: `.PHONY` line lists all targets.

---

### Task 4: Commit and close Issue #3

- [ ] **Step 1: Stage and commit**

```bash
git add Makefile Procfile
git commit -m "feat: add Makefile and Procfile for local development (closes #3)"
```

- [ ] **Step 2: Close Issue #3**

```bash
gh issue close 3 --repo Ayut0/coffee-journal
```
