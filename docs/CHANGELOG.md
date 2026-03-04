# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Membership system:** Members now receive a sequential membership number in `1991-XXX` format upon registration.
- **Carteirinha page (`GET /carteirinha`):** Digital membership card showing member number, name, email, favorite console, and join date. Protected by `RequireAuth`.
- **Rental system:** Members can rent available games directly from the shelf via `POST /rent`. Games are marked as rented with the member's name displayed.
- **Returns dashboard (`GET /admin/returns`):** Admin page listing all active rentals with a [Devolver] button for each. Protected by `RequireAdmin`.
- **Return game flow (`POST /admin/return-game`):** Processes a game return, marking the copy as available again.
- **Admin inventory page (`GET /admin/inventory`):** Full catalog listing in a NES-style table with [Editar] buttons for each game. Protected by `RequireAdmin`.
- **Game edit page (`GET /admin/edit/{id}`):** Form to edit game title, platform, summary, and source magazine. Useful for translating IGDB data to Portuguese.
- **Game update flow (`POST /admin/update-game`):** Saves edited game fields to the database.
- **Migration `003_membership_and_rental_support.sql`:** Adds `membership_number`, `address`, and `phone` columns to members; creates `membership_seq` sequence; backfills existing members with membership numbers; creates `game_copies` for existing games.
- `NextMembershipNumber`, `ListGamesWithAvailability`, `RentGame`, `ReturnGame`, and `ListActiveRentals` methods to the `Store` interface.
- `GameAvailability` and `ActiveRental` structs for rental display data.
- `UpdateGame` and `GetGameByID` methods to the `Store` interface.
- IGDB search now returns platform data (`platforms.name`, `platforms.abbreviation`).
- IGDB search limit increased from 5 to 10 results.
- `RequireAuth` middleware for member-only routes.
- `AddGame` now atomically creates a `game_copy` record in the same transaction.
- Navigation links between admin pages (ACERVO, ABASTECER, PRATELEIRA, DEVOLVER).
- `internal/auth` package with HMAC-SHA256 signed cookies for session management.
- `internal/middleware` package with `RequireAuth` and `RequireAdmin` middleware.
- Password hashing with bcrypt (`golang.org/x/crypto/bcrypt`) on member registration.
- Real login validation: profile name + password verified against the database.
- `GetMemberByID` and `GetMemberByProfileName` methods to the `Store` interface.
- `COOKIE_SECRET` and `ADMIN_EMAIL` environment variables for security configuration.
- `.env.example` file with placeholder values for all required environment variables.
- `docs/` directory with project documentation.

### Changed
- Games shelf (`/games`) now displays real-time availability: [ALUGAR] for available games (logged in), [DISPONIVEL] for available games (not logged in), and [ALUGADO - Com o Socio: Nome] for rented games.
- Games shelf uses a single unified catalog view instead of separate Releases/Catalog sections.
- `POST /admin/purchase` now redirects to the edit page (`/admin/edit/{id}`) instead of back to stock, allowing immediate translation of IGDB data.
- Login form now requires both "Nome do Socio" (profile name) and password.
- Session cookie now stores a signed member UUID instead of a plain-text name.
- Session cookie now includes `MaxAge`, `SameSite=Strict`, and `HttpOnly` flags.
- `NewHandler` constructor now accepts a `cookieSecret` parameter.
- `CreateMember` now auto-assigns a membership number and response no longer exposes the password hash.

### Security
- **Passwords are now hashed with bcrypt** before being stored in the database.
- **Cookies are now HMAC-signed** — forged cookies are rejected automatically.
- **Admin routes (`/admin/*`) are protected** by `RequireAdmin` middleware that checks authentication and verifies the member's email against `ADMIN_EMAIL`.
- **Member routes (`/carteirinha`, `/rent`) are protected** by `RequireAuth` middleware.

## [0.2.0] - 2026-03-03

### Added
- Docker Compose configuration for PostgreSQL 15 (Alpine).
- Admin stock management page (`/admin/stock`) with IGDB search integration.
- Purchase game flow (`POST /admin/purchase`) to add games to the catalog.
- `GET /search` endpoint returning raw JSON from the IGDB API.
- `GET /games/{id}` endpoint returning a single game as JSON.
- `POST /members` endpoint for member registration (JSON API).
- Retro NES-style CSS theme (`web/static/css/retro.css`) inspired by forum aesthetics.
- `DATABASE_URL` environment variable support for PostgreSQL connection.

### Changed
- Migrated UI to "Nes Archive Forum V3" visual style with dark navy background.
- Games shelf uses a responsive CSS grid layout.
- Games table schema extended with `cover_url`, `source_magazine`, and `acquired_at` columns (migration `002`).

## [0.1.0] - 2026-03-03

### Added
- Initial Go project structure (`cmd/server`, `internal/` packages).
- Core data models: `Member`, `Game`, `GameCopy`, `Rental`.
- PostgreSQL database layer with `pgx/v5` connection pool.
- `Store` interface for decoupled data access.
- Database migration `001_initial_schema.sql` with `members`, `games`, `game_copies`, and `rentals` tables.
- IGDB API client (`internal/igdb`) with Twitch OAuth2 token retrieval and game search.
- Environment configuration loader using `godotenv`.
- Landing page (Balcao) with login form.
- Games shelf page (`/games`) with mock data fallback.
- Server-Side Rendering with `html/template` and NES.css + Press Start 2P font.
- Graceful shutdown with OS signal handling.
- `ARCHITECTURE.md` defining the project's vision, tech stack, and design principles.
- GPL v3 license.
