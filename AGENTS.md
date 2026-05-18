# AGENTS.md

This file is the single source of truth for AI coding agents (Gemini CLI, Claude Code, Copilot, Cursor, etc.) working on the **Modo Locadora** project.

## Project Overview

**Modo Locadora** is a retro-gaming session manager and backlog tracker designed to emulate the experience of 1990s Brazilian video rental stores (*locadoras*). 

- **Core Concept**: "Scarcity by design" — each game has limited physical copies.
- **Reputation**: Members must manage their reputation to avoid the "Wall of Shame."
- **Stack**: Go 1.24+, PostgreSQL 15, Server-Side Rendering (SSR) with `html/template`, NES.css for an 8-bit UI. No JavaScript.
- **License**: GPL v3.

## Core Philosophy & Design Guidelines

1.  **Language Split**: 
    - **Code, database columns, routes, and comments**: English.
    - **UI/User-facing text**: Portuguese (BR).
2.  **No JavaScript**: Strictly zero-JS, fully SSR approach.
3.  **Retro Aesthetic**: Adhere to the NES.css dark theme (`is-dark`, `#0A0E1A`) and the "video rental store" metaphor (e.g., "Sopro" for redemptions).
4.  **Progression**: Members progress from `Sócio Novato` -> `Prata` -> `Ouro` -> `Dono da Calçada`.
5.  **Scarcity**: Never bypass the copy-limit logic; if copies are 0, the game is "Alugado."

## Build & Run

The project uses [Task](https://taskfile.dev/) (`Taskfile.yml`) for common commands.

```bash
task build     # Build the binary
task check     # static analysis (vet + lint + build) — run before every commit
task dev       # go run ./cmd/server
task seed      # apply migrations + seed data
task up        # docker compose up -d --build
task reset     # full reset (down -v + up + seed)
task logs      # docker compose logs -f app
task psql      # connect to DB container
```

## Database & Migrations

PostgreSQL 15 via Docker Compose. Migrations are in `internal/database/migrations/`.

Shortcut: `go run ./cmd/server --seed` applies all migrations (001-011) + seed data.
Default Credentials: `tio_da_locadora` / `sopre_a_fita` / `modo_locadora`.

Refer to [docs/setup.md](docs/setup.md) for a full list of migrations and test accounts.

## Code Architecture

Refer to [ARCHITECTURE.md](ARCHITECTURE.md) for a high-level system overview (in Portuguese).

### Package Structure

| Package | Purpose |
|---------|---------|
| `cmd/server/main.go` | Entrypoint: config, template parsing, pgx pool, route wiring |
| `internal/handlers/` | HTTP handlers encapsulated in a `Handler` struct |
| `internal/database/` | `Store` interface (store.go) and `PostgreSQL` implementation (postgres.go) |
| `internal/middleware/` | `RequireAuth` (session check) and `RequireAdmin` (auth + email check) |
| `internal/auth/` | HMAC-SHA256 cookie signing/verification |
| `internal/igdb/` | IGDB API client (Twitch OAuth2) |
| `internal/models/` | Domain structs: Member, Game, GameCopy, Rental, Club, etc. |
| `internal/storage/` | `StorageProvider` for file uploads (Local vs GCS) |
| `web/templates/` | Standalone HTML templates (Portuguese UI) |
| `web/static/` | Global CSS (`retro.css`), logos, and uploaded assets |

### Authentication
- **Login**: `POST /login` -> bcrypt verify -> HMAC-signed cookie (`{uuid}.{hmac_hex}`).
- **Admin**: Cookie verified + email checked against `ADMIN_EMAIL` env var.

### Navigation Flow
1. `GET /games`: Platform selection grid (`platforms.html`).
2. `GET /games?platform=X`: Filtered game shelf (`games.html`).
3. `GET /games/{id}`: Full game detail with rental stats (`game_detail.html`).

## Routing

### Public
- `GET /` — Landing page with login and "Wall of Shame"
- `POST /login` — Authentication
- `POST /members` — Registration (JSON API)
- `GET /search?q=` — IGDB search (JSON API)
- `GET /clubs` — Public club listing
- `GET /clubs/{id}` — Public club detail

### Authenticated (`RequireAuth`)
- `GET /membership` — Digital membership card & password notes
- `POST /rent` — Rent a game
- `POST /membership/return` — Self-return a rental (with verdict)
- `POST /membership/redeem` — Clear debt status ("Sopro")
- `POST /clubs/new`, `POST /clubs/{id}/join`, etc.

### Admin (`RequireAdmin`)
- `GET /admin/stock` — IGDB search & add games
- `GET /admin/inventory` — Full catalog listing
- `GET /admin/edit/{id}` — Edit game form & history
- `GET /admin/returns` — Active rentals dashboard

## Conventions & Rules

1.  **Commit Format**: Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`).
2.  **Branching**: `main` (stable), `develop` (active).
3.  **Routing**: Standard library only — `mux.HandleFunc("METHOD /path", handler)`.
4.  **Templates**: Standalone HTML files. Use NES.css classes and `retro.css` utilities.
5.  **Migrations**: New changes MUST be a new numbered `.sql` file and registered in `cmd/server/main.go`.
6.  **Security**: Parameterized SQL. Never store plaintext passwords. Cookie secrets min 32 chars.
7.  **SRE**: Always run `task check` before proposing a PR.

## Common Tasks

### Adding a new route
1. Add `Store` interface method in `store.go` and implement in `postgres.go`.
2. Add handler method to `internal/handlers/handler.go`.
3. Create/modify template in `web/templates/`.
4. Register route in `cmd/server/main.go` with middleware.
5. Run `task build` to verify template parsing.

### Adding a database migration
1. Create `internal/database/migrations/0XX_description.sql`.
2. Register the file in `sqlFiles` list in `cmd/server/main.go`.
3. Document in `docs/setup.md`.
