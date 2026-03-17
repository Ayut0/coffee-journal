# ADR-007: HTTP Router and Logging

## Status
Accepted

## Date
2026-03-16

## Context
The Go API needs an HTTP router and a logging strategy. These are chosen together as they are both foundational infrastructure decisions with no dependency on the domain model.

## Decisions

### Echo v4 as the HTTP router
Use `github.com/labstack/echo/v4` instead of Chi.

Both Echo and Chi are lightweight, idiomatic Go routers. Echo was chosen because:
- **Built-in request binding** — `c.Bind(&req)` decodes JSON request bodies into structs automatically
- **Built-in validation hook** — integrates with a validator interface cleanly
- **Middleware ecosystem** — CORS, request ID, rate limiting, JWT middleware are all first-class
- **Consistent handler signature** — `func(c echo.Context) error` is uniform across all handlers, including error handling
- **Reference project alignment** — the existing clean architecture reference uses Echo, reducing context switching

Chi would have been equally valid — it is more minimal and closer to `net/http`. Echo's additional ergonomics justify the choice given the number of handlers this project will have.

### zerolog for structured logging
Use `github.com/rs/zerolog` for logging.

Reasons:
- **Zero allocation** — zerolog is designed to produce no heap allocations on the hot path
- **JSON output** — structured logs are machine-readable, compatible with log aggregation tools (Fly.io logs, Datadog, etc.)
- **Simple API** — `log.Info().Str("key", "value").Msg("message")` is readable and explicit
- **Production habit** — structured logging is standard in production Go services; learning it from the start avoids retrofitting later

```go
log.Info().
    Str("method", r.Method).
    Str("path", r.URL.Path).
    Int("status", 200).
    Dur("latency", duration).
    Msg("request")
```

## Alternatives Considered

### Chi
- More minimal, closer to `net/http`
- No built-in binding or validation
- Would have been a valid choice — rejected in favour of Echo's ergonomics

### Go standard `log` / `slog`
- `slog` (added in Go 1.21) is a solid structured logger in the standard library
- zerolog is faster and has a more ergonomic API for request-scoped logging
- zerolog is already used in the reference project

### logrus
- Widely used but slower than zerolog
- Effectively in maintenance mode

## Consequences

### Positive
- Echo's request binding eliminates manual `json.Decode` boilerplate in every handler
- zerolog produces structured logs from day one — no migration from printf-style logs later
- Both packages are actively maintained with strong communities

### Negative
- Echo adds a thin abstraction over `net/http` — slightly less transparent than Chi for learning raw HTTP handling
- zerolog's chained API has a learning curve compared to `fmt.Println` debugging
