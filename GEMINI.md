# GEMINI.md — Project Context & Instructions

This file provides foundational context and instructions for Gemini CLI when working on the **Modo Locadora** project.

## Project Overview

**Modo Locadora** is a retro-gaming session manager and backlog tracker designed to emulate the experience of 1990s Brazilian video rental stores (*locadoras*). It focuses on "scarcity by design," where each game has limited physical copies, and members must manage their reputation to unlock titles and avoid the "Wall of Shame."

### Technical Stack
- **Backend:** Go 1.24+ (Standard Library `net/http` for routing).
- **Database:** PostgreSQL 15+ (accessed via `pgx/v5`).
- **Frontend:** Server-Side Rendering (SSR) with Go `html/template`.
- **Styling:** NES.css 2.3.0 (8-bit aesthetic) + "Press Start 2P" font.
- **External API:** IGDB (via Twitch OAuth2) for game metadata and covers.
- **Infrastructure:** Docker & Docker Compose (app + database).

### Core Philosophy
- **Portuguese for UI:** All user-facing text in templates must be in **Portuguese (BR)**.
- **English for Code:** All variable names, functions, database columns, routes, and comments must be in **English**.
- **No JavaScript:** The project strictly follows a zero-JS, fully SSR approach.
- **Retro Aesthetic:** Adhere to the NES.css dark theme and the "video rental store" metaphor in features.

---

## Architectural Map

### Directory Structure
- `cmd/server/main.go`: Application entrypoint, route registration, and template parsing.
- `internal/handlers/`: HTTP request handlers (encapsulated in the `Handler` struct).
- `internal/database/`: 
    - `store.go`: `Store` interface defining all DB operations.
    - `postgres.go`: PostgreSQL implementation of the `Store` interface.
    - `migrations/`: Numbered SQL migration files (001-011+).
- `internal/models/`: Domain entities (Member, Game, Rental, Club, etc.).
- `internal/middleware/`: Authentication and Authorization (`RequireAuth`, `RequireAdmin`).
- `internal/auth/`: Cookie-based session management (HMAC-SHA256 signing).
- `internal/igdb/`: Client for the IGDB API.
- `internal/almanac/`: Static gaming ephemerides and trivia logic.
- `web/templates/`: Standalone HTML templates with inline `<style>` for page-specific CSS.
- `web/static/`: Global CSS (`retro.css`), logos, and uploaded assets (covers, badges).

### Key Workflows
1. **Authentication:** Uses bcrypt for passwords and HMAC-signed cookies for sessions.
2. **Rental Cycle:** Member rents game (3-day deadline) -> Status changes to "Rented" -> Member returns with a "Verdict" (Zerou, Jogou, Desistiu) -> Copy becomes "Available."
3. **Wall of Shame:** A background job (`internal/jobs/overdue.go`) checks for overdue rentals every 5 minutes, auto-returns them, and penalizes the member's reputation.

---

## Development Standards

### Building & Running
- **Local Dev:** `go run ./cmd/server` (requires `DATABASE_URL` in `.env`).
- **Docker (Full Stack):** `docker compose up -d --build`.
- **Migrations/Seed:** `go run ./cmd/server --seed` (applies all migrations and test data).

### Task Runner (Taskfile.yml)
- `task build`: Compiles the server binary.
- `task check`: Runs `go vet`, `golangci-lint`, and `build`.
- `task seed`: Populates the database with initial data.
- `task reset`: Resets the database and re-seeds (Docker only).

### Contribution Rules
- **Migrations:** New database changes MUST be added as a new numbered `.sql` file in `internal/database/migrations/` and registered in the `sqlFiles` list in `cmd/server/main.go`.
- **Templates:** Use NES.css classes (`is-dark`, `nes-container`, etc.). Keep layouts consistent with `layout.html`.
- **Linting:** Ensure `golangci-lint run ./...` passes before proposing changes.

---

## Environment Variables
Ensure the following are defined in `.env`:
- `DATABASE_URL`: Postgres connection string.
- `COOKIE_SECRET`: Min 32-character string for session security.
- `ADMIN_EMAIL`: Email address granted admin privileges.
- `TWITCH_CLIENT_ID` & `TWITCH_CLIENT_SECRET`: For IGDB integration.

---

## Design Guidelines (Locadora Meta)
- **Status Indicators:** Use the "Sopro" (blowing on the cartridge) metaphor for redemptions.
- **Titles:** Members progress from `Sócio Novato` -> `Prata` -> `Ouro` -> `Dono da Calçada`.
- **Scarcity:** Never bypass the copy-limit logic; if copies are 0, the game is "Alugado."
