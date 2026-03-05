# AGENTS.md

This file provides guidance to AI coding agents (Jules, Claude Code, Copilot, etc.) when working with this repository.

## Project Overview

Modo Locadora is a retro-gaming session manager that emulates 90s Brazilian video rental stores. Built with Go 1.24+, PostgreSQL 15, server-side rendered HTML templates, and NES.css for an 8-bit UI aesthetic. Licensed under GPL v3.

## Build & Run

```bash
# Build the binary
go build -o modo-locadora ./cmd/server

# Run the server (requires .env configured)
go run ./cmd/server
# Starts on port 8080 (override with PORT env var)

# Static analysis (no test suite — use this before commits)
go vet ./...
```

## Database Setup

PostgreSQL 15 via Docker:
```bash
docker compose up -d
```

Apply migrations in order:
```bash
psql $DATABASE_URL -f internal/database/migrations/001_initial_schema.sql
psql $DATABASE_URL -f internal/database/migrations/002_update_games_table.sql
psql $DATABASE_URL -f internal/database/migrations/003_membership_and_rental_support.sql
psql $DATABASE_URL -f internal/database/migrations/004_password_notes.sql
```

Default credentials: `tio_da_locadora` / `sopre_a_fita` / database `modo_locadora`.

## Environment Variables

Required in `.env` (see `.env.example`):
- `TWITCH_CLIENT_ID`, `TWITCH_CLIENT_SECRET` — IGDB API via Twitch OAuth2
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DATABASE_URL`
- `COOKIE_SECRET` — HMAC-SHA256 key for session cookies (min 32 chars)
- `ADMIN_EMAIL` — email that grants admin panel access

## Code Architecture

### Package Structure

| Package | Purpose |
|---------|---------|
| `cmd/server/main.go` | Entrypoint: config, template parsing, pgx pool, route wiring |
| `internal/handlers/handler.go` | All HTTP handlers in a `Handler` struct (Store + templates + config) |
| `internal/database/store.go` | `Store` interface defining all DB operations |
| `internal/database/postgres.go` | PostgreSQL implementation using pgx/v5 pool with transactions |
| `internal/middleware/middleware.go` | `RequireAuth` (cookie check) and `RequireAdmin` (auth + email) |
| `internal/auth/auth.go` | HMAC-SHA256 cookie signing/verification |
| `internal/igdb/igdb.go` | IGDB API client (Twitch OAuth2 token flow + game search) |
| `internal/config/config.go` | `.env` loader via godotenv |
| `internal/models/` | Data structs: Member, Game, GameCopy, Rental |
| `web/templates/` | Go HTML templates (7 pages, all in Portuguese) |
| `web/static/css/retro.css` | NES.css dark theme overrides and shared utility classes |

### Request Flow

1. `main.go` registers routes on `http.ServeMux` with method-pattern routing (`GET /path`, `POST /path`)
2. Middleware verifies `session_member` cookie via HMAC, injects member UUID into request context
3. Handler reads context, calls `Store` interface methods, renders template with data struct
4. Store implementation executes parameterized SQL via pgx connection pool

### Database Schema

4 tables: `members`, `games`, `game_copies`, `rentals` + `membership_seq` sequence (generates `1991-XXX` numbers).
Key relationship: Game -> GameCopy -> Rental <- Member. Copy status enum: `available` | `rented`.

### Authentication

- Registration: bcrypt hash + sequential membership number
- Login: `POST /login` -> bcrypt verify -> HMAC-signed cookie (`{uuid}.{hmac_hex}`)
- Middleware splits cookie, verifies HMAC signature on each request
- Admin: cookie verified + email checked against `ADMIN_EMAIL` env var

## Routing

### Public
- `GET /` — Landing page with login form
- `POST /login` — Authentication
- `POST /members` — Registration (JSON API)
- `GET /search?q=` — IGDB search (JSON API)

### Authenticated (RequireAuth middleware)
- `GET /games` — Game shelf with rental buttons
- `GET /carteirinha` — Member card
- `POST /carteirinha/notes` — Save password notebook
- `POST /rent` — Rent a game

### Admin (RequireAdmin middleware)
- `GET /admin/stock` — IGDB search & add games
- `POST /admin/purchase` — Confirm game acquisition
- `GET /admin/inventory` — Full catalog listing
- `GET /admin/edit/{id}` — Edit game form
- `POST /admin/update-game` — Save game edits
- `GET /admin/returns` — Active rentals dashboard
- `POST /admin/return-game` — Process game return

## Conventions & Rules

1. **Language split**: Code, routes, database columns in **English**. UI text (templates) in **Portuguese (BR)**.
2. **Commit format**: Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`).
3. **Routing**: Standard library only — `mux.HandleFunc("METHOD /path", handler)`. No third-party routers.
4. **No test framework**: Validate with `go build ./...` and `go vet ./...` before commits.
5. **CSS**: NES.css 2.3.0 classes with dark theme overrides in `retro.css`. Shared utility classes: `.btn-nav`, `.btn-sm`, `.title-main`, `.title-sub`, `.footer-copyright`, `.nav-bar`, `.form-actions`, `.empty-state`, `.success-balloon`.
6. **Templates**: Each page is a standalone HTML file. Page-specific CSS in inline `<style>`, shared CSS from `retro.css`.
7. **No JavaScript**: The frontend is fully static SSR. No JS frameworks or inline scripts.
8. **Security**: Always use parameterized SQL queries. Never store plaintext passwords. Cookie secrets must be 32+ chars.
9. **Scarcity by design**: Each game has limited physical copies. All copies rented = game unavailable.

## Dependencies

```
pgx/v5           — PostgreSQL driver + connection pool
godotenv         — .env file loading
google/uuid      — UUID generation
golang.org/x/crypto — bcrypt password hashing
```

## Common Tasks for Agents

### Adding a new route
1. Add handler method to `internal/handlers/handler.go`
2. Register route in `cmd/server/main.go` with appropriate middleware wrapper
3. If needed, add Store interface method in `store.go` and implement in `postgres.go`
4. Create or modify template in `web/templates/`

### Adding a new database migration
1. Create numbered SQL file in `internal/database/migrations/` (e.g., `005_description.sql`)
2. Update Store interface and postgres implementation if schema changes affect queries
3. Document the migration in this file's Database Setup section

### Modifying UI/templates
1. Use NES.css classes from the framework — check existing usage in templates for patterns
2. Add shared styles to `web/static/css/retro.css`, page-specific styles in template `<style>` block
3. Follow the dark theme: dark backgrounds (#0A0E1A, #1A1A1A), light text, NES.css `is-dark` variants
4. Verify with `go build ./...` (template parsing happens at startup)
