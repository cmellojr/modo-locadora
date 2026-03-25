# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
# Local build
go build -o modo-locadora ./cmd/server
go run ./cmd/server          # starts on :8080 (override with PORT env var)
go vet ./...                 # static analysis — run before every commit

# Docker (full stack: app + PostgreSQL)
docker compose up -d --build
```

## Database

PostgreSQL 15 via Docker Compose. Migrations applied in order:

```bash
psql $DATABASE_URL -f internal/database/migrations/001_initial_schema.sql
psql $DATABASE_URL -f internal/database/migrations/002_update_games_table.sql
psql $DATABASE_URL -f internal/database/migrations/003_membership_and_rental_support.sql
psql $DATABASE_URL -f internal/database/migrations/004_password_notes.sql
psql $DATABASE_URL -f internal/database/migrations/005_auto_return_reputation.sql
psql $DATABASE_URL -f internal/database/migrations/006_activities_feed.sql
# 007 is seed data — applied via --seed flag, not manually
```

Shortcut: `go run ./cmd/server --seed` applies all migrations + seed data in one step.

Default DB credentials: `tio_da_locadora` / `sopre_a_fita` / `modo_locadora`.
Seed test members: `MegaDriveKid` / `sega1991`, `Devedor` / `atrasado123`, `Novato` / `novato2026`.

## Environment Variables

Required in `.env` (see `.env.example`):
- `TWITCH_CLIENT_ID`, `TWITCH_CLIENT_SECRET` — IGDB API via Twitch OAuth2
- `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DATABASE_URL`
- `COOKIE_SECRET` — HMAC-SHA256 key for session cookies (min 32 chars)
- `ADMIN_EMAIL` — email that grants admin panel access

## Architecture

Go 1.24, standard library `net/http.ServeMux` with method-pattern routing. Server-side rendered with `html/template`, no JavaScript. NES.css 2.3.0 + Press Start 2P for retro 8-bit UI. Multi-stage Dockerfile (golang:1.24-alpine builder + alpine:3.21 runtime).

### Package Structure

- **`cmd/server/main.go`** — Entrypoint: config, templates, pgx pool, routes, middleware
- **`internal/handlers/handler.go`** — All HTTP handlers in a `Handler` struct (Store + cookieSecret)
- **`internal/database/store.go`** — `Store` interface + view structs (GameAvailability, PlatformSummary, GameDetail, ActiveRental, ShameEntry, ActivityEntry, MemberRental)
- **`internal/database/postgres.go`** — PostgreSQL implementation (pgx/v5 pool, transactions)
- **`internal/middleware/middleware.go`** — `RequireAuth` and `RequireAdmin` middleware
- **`internal/auth/auth.go`** — HMAC-SHA256 cookie signing/verification
- **`internal/igdb/igdb.go`** — IGDB API client (Twitch OAuth2 + game search)
- **`internal/almanac/almanac.go`** — Static gaming ephemerides by day-of-year
- **`internal/jobs/overdue.go`** — Background goroutine: auto-returns overdue rentals every 5 min
- **`internal/config/config.go`** — `.env` loader via godotenv
- **`internal/models/`** — Domain structs: Member (with status/late_count), Game, GameCopy, Rental
- **`web/templates/`** — 9 standalone HTML templates (Portuguese UI)
- **`web/static/css/retro.css`** — NES.css dark theme overrides and utility classes
- **`web/static/covers/`** — Uploaded Brazilian game covers (Docker volume)

### Database Schema

6 tables + 1 sequence. Key relationship: `Game -> GameCopy -> Rental <- Member`.

- `members` — profile_name, email, password_hash, membership_number (`1991-XXX`), status (`active`|`em_debito`), late_count
- `games` — title, igdb_id, platform, summary, cover_url, source_magazine, acquired_at
- `game_copies` — game_id, status (`available`|`rented`)
- `rentals` — member_id, copy_id, rented_at, due_at (3 days), returned_at, public_legacy (verdict)
- `activities` — event_type, member_name, game_title, created_at (denormalized feed)

### Auth Flow

POST `/login` -> bcrypt verify -> HMAC-signed cookie `{uuid}.{hmac_hex}` -> middleware verifies on each request.

### 3-Level Game Navigation

- `GET /games` (no params) -> platform selection grid (platforms.html)
- `GET /games?platform=X` -> simplified cards for that console (games.html)
- `GET /games/{id}` -> full game detail with rental stats (game_detail.html)

## Conventions

- **Language split**: Code, routes, DB columns in English. UI text in Portuguese (BR).
- **Commit format**: Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`)
- **Branching**: `main` (stable) + `develop` (active). Feature branches: `feature/*`, `fix/*`, `hotfix/*`, `docs/*`
- **Routing**: Standard library only — `mux.HandleFunc("METHOD /path", handler)`
- **Validation**: `go build ./...` and `go vet ./...` before commits (no test framework)
- **CSS**: NES.css classes + dark theme overrides in `retro.css`. Shared utilities: `.btn-nav`, `.btn-sm`, `.title-main`, `.title-sub`, `.footer-copyright`, `.nav-bar`, `.form-actions`, `.empty-state`, `.success-balloon`
- **Templates**: Standalone HTML files. Page-specific CSS in inline `<style>`, shared CSS from `retro.css`
- **No JavaScript**: Fully static SSR

## Dependencies

```
pgx/v5           — PostgreSQL driver + pool
godotenv         — .env loading
google/uuid      — UUID generation
golang.org/x/crypto — bcrypt
```
