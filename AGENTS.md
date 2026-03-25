# AGENTS.md

This file provides guidance to AI coding agents (Jules, Claude Code, Copilot, Cursor, etc.) when working with this repository.

## Project Overview

Modo Locadora is a retro-gaming session manager that emulates 90s Brazilian video rental stores. Built with Go 1.24+, PostgreSQL 15, server-side rendered HTML templates, and NES.css for an 8-bit UI aesthetic. Dockerized (app + DB). Licensed under GPL v3.

## Build & Run

```bash
# Build the binary
go build -o modo-locadora ./cmd/server

# Run the server (requires .env configured)
go run ./cmd/server
# Starts on port 8080 (override with PORT env var)

# Static analysis (no test suite — use this before commits)
go vet ./...

# Docker full stack (app + PostgreSQL)
docker compose up -d --build
```

## Database Setup

PostgreSQL 15 via Docker Compose. Apply migrations in order:

```bash
psql $DATABASE_URL -f internal/database/migrations/001_initial_schema.sql
psql $DATABASE_URL -f internal/database/migrations/002_update_games_table.sql
psql $DATABASE_URL -f internal/database/migrations/003_membership_and_rental_support.sql
psql $DATABASE_URL -f internal/database/migrations/004_password_notes.sql
psql $DATABASE_URL -f internal/database/migrations/005_auto_return_reputation.sql
psql $DATABASE_URL -f internal/database/migrations/006_activities_feed.sql
# 007 is seed data — applied via --seed flag
```

Shortcut: `go run ./cmd/server --seed` applies all migrations + seed data in one step.

Default DB credentials: `tio_da_locadora` / `sopre_a_fita` / database `modo_locadora`.
Seed test members: `MegaDriveKid` / `sega1991`, `Devedor` / `atrasado123`, `Novato` / `novato2026`.

## Environment Variables

Required in `.env` (see `.env.example`):
- `TWITCH_CLIENT_ID`, `TWITCH_CLIENT_SECRET` — IGDB API via Twitch OAuth2
- `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DATABASE_URL`
- `COOKIE_SECRET` — HMAC-SHA256 key for session cookies (min 32 chars)
- `ADMIN_EMAIL` — email that grants admin panel access

## Code Architecture

Go 1.24, standard library `net/http.ServeMux` with method-pattern routing. Server-side rendered with `html/template`, no JavaScript. NES.css 2.3.0 + Press Start 2P for retro 8-bit UI. Multi-stage Dockerfile (golang:1.24-alpine builder + alpine:3.21 runtime).

### Package Structure

| Package | Purpose |
|---------|---------|
| `cmd/server/main.go` | Entrypoint: config, template parsing, pgx pool, route wiring |
| `internal/handlers/handler.go` | All HTTP handlers in a `Handler` struct (Store + cookieSecret) |
| `internal/database/store.go` | `Store` interface + view structs (GameAvailability, PlatformSummary, GameDetail, ActiveRental, ShameEntry, ActivityEntry, MemberRental) |
| `internal/database/postgres.go` | PostgreSQL implementation (pgx/v5 pool, transactions) |
| `internal/middleware/middleware.go` | `RequireAuth` (cookie check) and `RequireAdmin` (auth + email) |
| `internal/auth/auth.go` | HMAC-SHA256 cookie signing/verification |
| `internal/igdb/igdb.go` | IGDB API client (Twitch OAuth2 token flow + game search) |
| `internal/almanac/almanac.go` | Static gaming ephemerides by day-of-year |
| `internal/jobs/overdue.go` | Background goroutine: auto-returns overdue rentals every 5 min |
| `internal/config/config.go` | `.env` loader via godotenv |
| `internal/models/` | Domain structs: Member (with status/late_count), Game, GameCopy, Rental |
| `web/templates/` | 9 standalone HTML templates (Portuguese UI) |
| `web/static/css/retro.css` | NES.css dark theme overrides and shared utility classes |
| `web/static/covers/` | Uploaded Brazilian game covers (Docker volume) |

### Request Flow

1. `main.go` registers routes on `http.ServeMux` with method-pattern routing (`GET /path`, `POST /path`)
2. Middleware verifies `session_member` cookie via HMAC, injects member UUID into request context
3. Handler reads context, calls `Store` interface methods, renders template with data struct
4. Store implementation executes parameterized SQL via pgx connection pool

### Database Schema

6 tables + 1 sequence. Key relationship: `Game -> GameCopy -> Rental <- Member`.

- `members` — profile_name, email, password_hash, membership_number (`1991-XXX`), status (`active`|`em_debito`), late_count
- `games` — title, igdb_id, platform, summary, cover_url, source_magazine, acquired_at
- `game_copies` — game_id, status (`available`|`rented`)
- `rentals` — member_id, copy_id, rented_at, due_at (3 days), returned_at, public_legacy (verdict)
- `activities` — event_type, member_name, game_title, created_at (denormalized feed)
- `membership_seq` — generates sequential numbers (1991-001, 1991-002, ...)

### Authentication

- Registration: bcrypt hash + sequential membership number
- Login: `POST /login` -> bcrypt verify -> HMAC-signed cookie (`{uuid}.{hmac_hex}`)
- Middleware splits cookie, verifies HMAC signature on each request
- Admin: cookie verified + email checked against `ADMIN_EMAIL` env var

### 3-Level Game Navigation

- `GET /games` (no params) -> platform selection grid (`platforms.html`)
- `GET /games?platform=X` -> simplified cartridge cards for that console (`games.html`)
- `GET /games/{id}` -> full game detail with rental stats (`game_detail.html`)

## Routing

### Public
- `GET /` — Landing page with login form (redirects to `/games` if authenticated)
- `POST /login` — Authentication
- `POST /members` — Registration (JSON API)
- `GET /search?q=` — IGDB search (JSON API)

### Authenticated (RequireAuth middleware)
- `GET /games` — Platform selection or filtered game shelf (`?platform=X`)
- `GET /games/{id}` — Game detail page with rental stats
- `GET /carteirinha` — Digital membership card
- `POST /carteirinha/notes` — Save password notebook
- `POST /carteirinha/redeem` — Clear debt status
- `POST /carteirinha/return` — Self-return a rental (with verdict)
- `POST /rent` — Rent a game

### Admin (RequireAdmin middleware)
- `GET /admin/stock` — IGDB search & add games
- `POST /admin/purchase` — Confirm game acquisition
- `GET /admin/inventory` — Full catalog listing
- `GET /admin/edit/{id}` — Edit game form (with cover upload)
- `POST /admin/update-game` — Save game edits (multipart/form-data)
- `GET /admin/returns` — Active rentals dashboard
- `POST /admin/return-game` — Process game return

## Conventions & Rules

1. **Language split**: Code, routes, database columns in **English**. UI text (templates) in **Portuguese (BR)**.
2. **Commit format**: Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`).
3. **Branching**: `main` (stable) + `develop` (active). Feature branches: `feature/*`, `fix/*`, `hotfix/*`, `docs/*`.
4. **Routing**: Standard library only — `mux.HandleFunc("METHOD /path", handler)`. No third-party routers.
5. **Validation**: `go build ./...` and `go vet ./...` before commits (no test framework).
6. **CSS**: NES.css 2.3.0 classes with dark theme overrides in `retro.css`. Shared utility classes: `.btn-nav`, `.btn-sm`, `.title-main`, `.title-sub`, `.footer-copyright`, `.nav-bar`, `.form-actions`, `.empty-state`, `.success-balloon`.
7. **Templates**: Each page is a standalone HTML file. Page-specific CSS in inline `<style>`, shared CSS from `retro.css`.
8. **No JavaScript**: The frontend is fully static SSR. No JS frameworks or inline scripts.
9. **Security**: Parameterized SQL queries. Never store plaintext passwords. Cookie secrets must be 32+ chars.
10. **Scarcity by design**: Each game has limited physical copies. All copies rented = game unavailable.

## Dependencies

```
pgx/v5           — PostgreSQL driver + connection pool
godotenv         — .env file loading
google/uuid      — UUID generation
golang.org/x/crypto — bcrypt password hashing
```

## Common Tasks for Agents

### Adding a new route
1. Add Store interface method in `store.go` and implement in `postgres.go`
2. Add handler method to `internal/handlers/handler.go`
3. Create or modify template in `web/templates/`
4. Register route in `cmd/server/main.go` with appropriate middleware wrapper
5. Run `go build ./...` and `go vet ./...`

### Adding a new database migration
1. Create numbered SQL file in `internal/database/migrations/` (e.g., `008_description.sql`)
2. Update Store interface and postgres implementation if schema changes affect queries
3. Document the migration in `docs/SETUP.md`

### Modifying UI/templates
1. Use NES.css classes — check existing templates for patterns
2. Shared styles: `web/static/css/retro.css`. Page-specific: inline `<style>` in template
3. Dark theme: backgrounds `#0A0E1A`, `#1A1A1A`; NES.css `is-dark` variants
4. Verify with `go build ./...` (template parsing happens at startup)
