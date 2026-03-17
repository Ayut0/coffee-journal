# ADR-003: Data Design

## Status
Accepted

## Date
2026-03-14

## Context
The Coffee Journal needs a data model that reflects the real domain — not just a schema that happens to work. We are applying Domain-Driven Design (DDD) principles to make the model explicit and intentional.

The key domain concepts are:
- A **Bean** represents a bag of coffee: who made it, where it's from, how it's roasted
- A **Tasting** represents one session of drinking and evaluating a bean
- A tasting captures flavor perception (tags, scores, aftertaste) and an optional written note

The primary access patterns are:
- Timeline: all tastings, newest first, across all beans
- Bean detail: one bean + all its tastings + average scores
- Search: beans and tasting notes by keyword

## Domain Model

### Ubiquitous Language

| Term | Meaning |
|------|---------|
| Bean | A specific coffee product from a roaster — identified by name, roaster, origin, and roast level |
| Tasting | A single evaluation session for a bean — one bean can have many tastings over time |
| FlavorTag | A label describing a perceived flavor (e.g. Chocolate, Fruity, Floral) — no identity, reusable across tastings |
| Score | An integer from 1–5 rating a sensory dimension (acidity, aroma, body, sweetness) |
| Aftertaste | The character of the finish after swallowing — one of: Short, Clean, Lingering |
| BrewMethod | The brewing technique used for a tasting — one of: Espresso, Pour Over, French Press, AeroPress, Moka Pot, Cold Brew, Drip |
| RoastLevel | The roast degree of a bean — one of: Light, Medium, Dark |
| Process | The post-harvest processing method — one of: Washed, Natural, Honey |
| PackagePhoto | A photo of the coffee bag — stored in Cloudflare R2, referenced by URL on the Bean — nullable |
| Altitude | The elevation range at which the coffee was grown, in meters above sea level (e.g. 1500–2000m) — nullable |
| HarvestSeason | The crop season the beans were harvested (e.g. "2023/24") — nullable |
| GrindSize | The grind coarseness used in a tasting — one of: Extra Fine, Fine, Medium Fine, Medium, Medium Coarse, Coarse, Extra Coarse |
| OverallScore | A single 1–5 integer summarising the overall impression of a tasting |
| Visibility | Whether a bean is visible to other users — private by default |
| SoftDelete | Marking a record as deleted without removing it from the database — tracked via `deleted_at` timestamp |
| AuditLog | An append-only record of deletion events for historical tracking |

### Entities

**Bean** — has identity, persists over time, can be updated independently.

**Tasting** — has identity, persists over time, has its own independent lifecycle.

### Value Objects

**FlavorTag** — a plain label (string). No identity. Two tastings with `Chocolate` share the same concept, not the same object. Stored as `TEXT[]` in the tastings table.

**Score** — an integer 1–5 with domain meaning. Four dimensions: `acidity`, `aroma`, `body`, `sweetness`. Stored as separate integer columns (not JSONB) so they can be aggregated with `AVG()` in SQL.

**Aftertaste** — an enum: `Short | Clean | Lingering`. Stored as `TEXT` with a CHECK constraint.

**BrewMethod** — an enum: `Espresso | Pour Over | French Press | AeroPress | Moka Pot | Cold Brew | Drip`. Stored as `TEXT` with a CHECK constraint. Belongs on Tasting, not Bean — the same bean can be evaluated across multiple brew methods, each producing a separate Tasting record.

**RoastLevel** — an enum: `Light | Medium | Dark`. Stored as `TEXT` with a CHECK constraint.

**Process** — an enum: `Washed | Natural | Honey`. Nullable (not all beans have a known process). Stored as `TEXT` with a CHECK constraint.

**PackagePhoto** — a URL string pointing to a photo of the coffee bag stored in Cloudflare R2 object storage. Nullable — not all beans will have a photo. The DB stores only the URL; the file lives in R2 and is served directly from its public URL.

**Altitude** — a range of integers (meters above sea level). Stored as two nullable columns: `altitude_min INT` and `altitude_max INT`. Two columns allow range queries (`WHERE altitude_min >= 1500`) and handle single-value entries (both columns set to the same value).

**HarvestSeason** — a free-form text string (e.g. `"2023/24"`). Nullable. Stored as `TEXT` — no numeric encoding since the season format is not uniformly queryable.

**GrindSize** — an enum: `Extra Fine | Fine | Medium Fine | Medium | Medium Coarse | Coarse | Extra Coarse`. Stored as `TEXT` with a CHECK constraint on Tasting.

**OverallScore** — a single integer 1–5 representing the overall impression of a tasting. Stored as `SMALLINT NOT NULL` — distinct from the individual dimension scores, gives a quick summary without requiring aggregation.

**Visibility** — a boolean on Bean (`is_public`). Defaults to `false` (private). Allows future public sharing without a schema change. Stored as `BOOLEAN NOT NULL DEFAULT false`.

### Aggregates

**Bean** is an aggregate root.
- Controls its own consistency boundary
- Does not own a collection of Tastings — Tastings reference Bean by ID

**Tasting** is an aggregate root.
- Has an independent lifecycle (create, edit, delete without touching Bean)
- References Bean by `bean_id` (foreign key) — this is a cross-aggregate reference by identity, not by object ownership
- If you want a bean's tastings, you query the TastingRepository filtered by `bean_id`

### Why Tasting is not inside the Bean aggregate

The timeline page loads all tastings across all beans — independently of which bean they belong to. If Tasting were inside the Bean aggregate, loading the timeline would require loading every Bean first, then extracting its tastings. This is an awkward fit for the access pattern.

Tasting has its own independent lifecycle: you create, edit, and delete a tasting without any change to the Bean. There is no invariant that requires Bean and Tasting to change together in a single transaction.

Cross-aggregate references are by ID only — the Tasting stores `bean_id`, not a Bean object.

### Repositories

One repository per aggregate root:

- **BeanRepository** — CRUD for Bean, average score queries
- **TastingRepository** — CRUD for Tasting, timeline queries, filter by `bean_id`

No repository for value objects — they are always accessed through their aggregate root.

## Database Schema

```sql
CREATE TABLE beans (
  id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id             UUID NOT NULL REFERENCES users(id),
  name                TEXT NOT NULL,
  roaster             TEXT NOT NULL,
  origin              TEXT NOT NULL,
  roast_level         TEXT NOT NULL CHECK (roast_level IN ('Light', 'Medium', 'Dark')),
  process             TEXT CHECK (process IN ('Washed', 'Natural', 'Honey')),
  altitude_min        INT,
  altitude_max        INT,
  harvest_season      TEXT,
  package_photo_url   TEXT,
  is_public           BOOLEAN NOT NULL DEFAULT false,
  created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at          TIMESTAMPTZ
);

CREATE TABLE tastings (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID NOT NULL REFERENCES users(id),
  bean_id     UUID NOT NULL REFERENCES beans(id) ON DELETE CASCADE,
  flavor_tags TEXT[] NOT NULL DEFAULT '{}',
  brew_method TEXT NOT NULL CHECK (brew_method IN ('Espresso', 'Pour Over', 'French Press', 'AeroPress', 'Moka Pot', 'Cold Brew', 'Drip')),
  grind_size  TEXT NOT NULL CHECK (grind_size IN ('Extra Fine', 'Fine', 'Medium Fine', 'Medium', 'Medium Coarse', 'Coarse', 'Extra Coarse')),
  acidity     SMALLINT NOT NULL CHECK (acidity BETWEEN 1 AND 5),
  aroma       SMALLINT NOT NULL CHECK (aroma BETWEEN 1 AND 5),
  body        SMALLINT NOT NULL CHECK (body BETWEEN 1 AND 5),
  sweetness   SMALLINT CHECK (sweetness BETWEEN 1 AND 5),
  overall     SMALLINT NOT NULL CHECK (overall BETWEEN 1 AND 5),
  aftertaste  TEXT NOT NULL CHECK (aftertaste IN ('Short', 'Clean', 'Lingering')),
  note_text   TEXT,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at  TIMESTAMPTZ
);
```

```sql
CREATE TABLE audit_log (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  entity_type TEXT NOT NULL,       -- 'user', 'bean', 'tasting'
  entity_id   UUID NOT NULL,
  action      TEXT NOT NULL,       -- 'deleted'
  actor_id    UUID,                -- user who performed the action (nullable for system actions)
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Indexes

```sql
-- Bean full-text search + user lookup
CREATE INDEX beans_fts_idx ON beans USING GIN (to_tsvector('english', name));
CREATE INDEX beans_created_at_idx ON beans (created_at DESC);
CREATE INDEX beans_user_id_idx ON beans (user_id);

-- Tasting timeline + bean + user lookup
CREATE INDEX tastings_created_at_idx ON tastings (created_at DESC);
CREATE INDEX tastings_bean_id_idx ON tastings (bean_id);
CREATE INDEX tastings_user_id_idx ON tastings (user_id);

-- Tasting full-text search on notes
CREATE INDEX tastings_fts_idx ON tastings USING GIN (to_tsvector('english', coalesce(note_text, '')));

-- Tasting flavor tag search
CREATE INDEX tastings_flavor_tags_idx ON tastings USING GIN (flavor_tags);

-- Audit log lookup by entity
CREATE INDEX audit_log_entity_idx ON audit_log (entity_type, entity_id);
```

## Decisions

### Scores as separate columns, not JSONB
Four score dimensions (`acidity`, `aroma`, `body`, `sweetness`) are stored as separate `SMALLINT` columns. This allows native `AVG(acidity)`, `AVG(aroma)` etc. in SQL without JSON extraction. Sweetness is nullable — it is considered optional in some evaluation frameworks.

### FlavorTags as TEXT[]
Flavor tags are value objects with no identity. A separate `tags` table with a join table would add two extra tables and joins for no benefit at this scale. `TEXT[]` with a GIN index supports containment queries (`@>`) efficiently.

### Enums as TEXT with CHECK constraints
Postgres `ENUM` types are hard to alter (adding a value requires a schema migration). `TEXT` with a `CHECK` constraint is equally safe at query time and easier to evolve.

### Package photo stored in Cloudflare R2, URL in DB
The coffee bag photo is stored in Cloudflare R2 (S3-compatible object storage). The `beans` table holds only the public URL. Reasons:
- Files served directly from R2's CDN — no load on the Go API for every image request
- R2 has no egress fees and a generous free tier
- S3-compatible API means Go uses `aws/aws-sdk-go-v2` with a custom endpoint — no proprietary SDK
- Eliminates the need for a persistent volume on Fly.io
- Upload flow: `POST /api/beans/{id}/photo` → Go uploads to R2 → stores returned URL on bean

### Altitude as two columns, not a range type
Altitude is stored as `altitude_min INT` and `altitude_max INT` rather than Postgres's native `int4range`. Two plain integer columns are simpler to work with in Go's `database/sql` (no custom scanner needed), straightforward to render in the UI, and sufficient for range queries. Both columns are nullable — altitude is not always listed on a bag.

### HarvestSeason as TEXT
Harvest season follows the specialty coffee convention of a slash-separated string (e.g. `"2023/24"`). Encoding this as two year integers would add complexity without enabling any useful queries. `TEXT` is the simplest correct representation.

### Soft delete with audit_log
Records are never hard deleted. Instead, `deleted_at TIMESTAMPTZ` is added to `users`, `beans`, and `tastings`. All queries filter `WHERE deleted_at IS NULL`. When a user is soft deleted, the application cascades `deleted_at` to their beans and tastings in a single transaction.

A separate `audit_log` table records each deletion event (`entity_type`, `entity_id`, `action`, `actor_id`). This provides a clean historical trail without polluting the main tables with tombstone rows that need to be queried around.

### Visibility defaults to private
`is_public BOOLEAN NOT NULL DEFAULT false` on `beans`. The app is designed to go public in the future — the column is in place now so no breaking migration is needed when public sharing is enabled. Private is the safe default.

### OverallScore as a required field
`overall SMALLINT NOT NULL` on tastings. An overall score gives a quick single-number summary without requiring score aggregation. It is intentionally separate from the dimension scores — a tasting can have high acidity and body but a low overall if the balance is off.

### GrindSize as a fixed enum
`grind_size TEXT NOT NULL CHECK (...)` with seven values from Extra Fine to Extra Coarse. A fixed list enforces consistency for filtering and comparison. Free-form text would make "Medium-Coarse" and "medium coarse" two different values in practice.

### No audio storage
Audio voice notes were removed from scope. The `tastings` table has no `audio_path` column. See PLAN.md for rationale.

## Consequences

### Positive
- Domain model is explicit and documented — ubiquitous language is established before writing code
- Repositories map 1:1 to aggregate roots — clean separation between domain and data layers
- Score columns enable simple SQL aggregation for average scores
- Schema constraints enforce domain rules at the database level (CHECK, NOT NULL, FK)

### Negative
- Two aggregate roots means the bean detail page requires two queries (fetch bean + fetch tastings by bean_id) — acceptable given the access pattern
- `TEXT[]` for flavor tags cannot enforce a controlled vocabulary at the DB level — validation must happen in the application layer
