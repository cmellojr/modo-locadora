# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
# Build
go build -o modo-locadora ./cmd/server

# Run (requires .env or environment variables set)
go run ./cmd/server
# Server starts on port 8080 (override with PORT env var)

# Static analysis
go vet ./...
```

## Database

PostgreSQL 15 via Docker:
```bash
docker compose up -d
```

Migrations must be applied in order:
```bash
psql $DATABASE_URL -f internal/database/migrations/001_initial_schema.sql
psql $DATABASE_URL -f internal/database/migrations/002_update_games_table.sql
psql $DATABASE_URL -f internal/database/migrations/003_membership_and_rental_support.sql
psql $DATABASE_URL -f internal/database/migrations/004_password_notes.sql
```

Default credentials in docker-compose.yml: `tio_da_locadora` / `sopre_a_fita` / `modo_locadora`.

## Environment Variables

Required in `.env` (see `.env.example`):
- `TWITCH_CLIENT_ID`, `TWITCH_CLIENT_SECRET` — IGDB API via Twitch OAuth2
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DATABASE_URL`
- `COOKIE_SECRET` — HMAC-SHA256 key for session cookies (≥32 chars)
- `ADMIN_EMAIL` — email that grants admin panel access

## Architecture

Go 1.24.3, standard library `net/http.ServeMux` with method-pattern routing (`GET /path`, `POST /path`). Server-side rendered with `html/template`, no JavaScript. NES.css 2.3.0 + Press Start 2P font for retro 8-bit UI.

### Package Structure

- **`cmd/server/main.go`** — Entrypoint: loads config, parses templates, creates pgx pool, wires routes with middleware
- **`internal/handlers/handler.go`** — All HTTP handlers in a single `Handler` struct that holds Store, templates, and config
- **`internal/database/store.go`** — `Store` interface defining all DB operations
- **`internal/database/postgres.go`** — PostgreSQL implementation using `pgx/v5` connection pool with transactions for atomic operations (AddGame, RentGame)
- **`internal/middleware/middleware.go`** — `RequireAuth` (cookie verification) and `RequireAdmin` (auth + email check) middleware
- **`internal/auth/auth.go`** — HMAC-SHA256 cookie signing/verification
- **`internal/igdb/igdb.go`** — IGDB API client (Twitch OAuth2 → game search)
- **`internal/config/config.go`** — `.env` loader via godotenv
- **`internal/models/`** — Data structs: Member, Game, GameCopy, Rental
- **`web/templates/`** — Go HTML templates (Portuguese UI)
- **`web/static/css/retro.css`** — NES.css dark theme overrides and utility classes

### Request Flow

1. `main.go` registers routes on `http.ServeMux` with middleware wrappers
2. Middleware verifies `session_member` cookie via HMAC, loads member UUID into context
3. Handler reads context, calls `Store` interface methods, renders template with data struct
4. Store implementation executes parameterized SQL via pgx pool

### Database Schema

4 tables: `members`, `games`, `game_copies`, `rentals` + `membership_seq` sequence (generates `1991-XXX` membership numbers). Key relationship: Game → GameCopy → Rental ← Member. Copy status enum: `available` | `rented`.

### Auth Flow

POST `/login` → bcrypt verify → sign cookie as `{uuid}.{hmac_hex}` → middleware on subsequent requests splits and verifies signature.

## Conventions

- **Language split**: Code, routes, database columns in English. UI text (templates) in Portuguese (BR).
- **Commit format**: Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`)
- **Routing**: Standard library only — `mux.HandleFunc("METHOD /path", handler)`
- **No test framework**: Run `go build ./...` and `go vet ./...` before commits
- **CSS**: All template styling uses NES.css classes with dark theme overrides in `retro.css`. Shared utility classes: `.btn-nav`, `.btn-sm`, `.title-main`, `.title-sub`, `.footer-copyright`, `.nav-bar`, `.form-actions`, `.empty-state`, `.success-balloon`
- **Templates**: Each page is a standalone HTML file with inline `<style>` for page-specific CSS, shared classes from `retro.css`

## Dependencies

```
pgx/v5        — PostgreSQL driver + pool
godotenv      — .env loading
google/uuid   — UUID generation
golang.org/x/crypto — bcrypt
```
