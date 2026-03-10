# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Login/logout flow**: Balcão is always the landing page; logged-in members see welcome message + navigation instead of login form. `POST /logout` clears session cookie.
- **Auth bar**: All pages display "Sócio: nome / [DESCONECTAR]" aligned top-right when logged in. CSS class `.auth-bar` in `retro.css`.
- **Console logos on platform grid**: Platform selection page shows SVG console logos (`web/static/img/logos/`) instead of game cover images. Auto-mapped via `platformLogoFile()` helper.
- **3-level game navigation**: `/games` shows platform selection grid, `?platform=X` filters by console, `/games/{id}` shows full game detail with rental stats (total rentals, top renter, current renter).
- **Brazilian cover upload**: Admin can upload local cover images (TecToy, Playtronic) via multipart form on the edit page. Uploaded covers stored in `web/static/covers/` (Docker volume).
- **Auto-return system**: Background job checks overdue rentals every 5 minutes, auto-returns them and penalizes members (`em_debito` status + `late_count` increment). Migration `005_auto_return_reputation.sql`.
- **Wall of Shame** (`Painel da Vergonha`): Landing page shows top members with late returns.
- **Member redemption**: `POST /carteirinha/redeem` clears debt status.
- **Dockerized application**: Multi-stage Dockerfile, Docker Compose runs app + PostgreSQL, `covers_data` volume for uploads.
- **Password notebook** (`Caderno de Passwords`): Members can save game codes on their membership card. Migration `004_password_notes.sql`.
- **`internal/jobs/` package**: Background goroutine for overdue rental processing.
- **CLAUDE.md** and **AGENTS.md**: AI agent guidance files.

### Changed
- **Landing page**: Removed authenticated redirect; Balcão always shown first with conditional login/welcome content.
- **Game shelf simplified**: Cards now show only cover, title, copy count, and availability status (no summary/magazine).
- **`GET /games/{id}`** changed from JSON API to server-rendered game detail page.
- **`POST /rent`** redirects to game detail page instead of shelf.
- **NES.css component expansion**: `nes-badge`, `nes-progress`, `nes-list`, `nes-avatar`, `nes-dialog`, `nes-balloon` patterns.
- **Font sizes and container widths** standardized across all pages.
- **Inline styles consolidated** into reusable CSS classes in `retro.css`.

## [0.3.0] - 2026-03-04

### Added
- **Membership system**: Sequential membership numbers in `1991-XXX` format.
- **Carteirinha page**: Digital membership card.
- **Rental system**: Members rent games via `POST /rent`.
- **Returns dashboard**: Admin page listing active rentals with return buttons.
- **Admin inventory and edit**: Catalog table with edit links and game edit form.
- **Migration `003`**: Membership fields, `membership_seq` sequence, auto-created copies.
- `RequireAuth` and `RequireAdmin` middleware.
- HMAC-SHA256 signed cookies and bcrypt password hashing.
- `.env.example` and `docs/` directory.

### Changed
- Game shelf displays real-time availability states.
- Login requires profile name + password with bcrypt verification.
- Session cookie stores signed UUID.

## [0.2.0] - 2026-03-03

### Added
- Docker Compose for PostgreSQL 15.
- Admin stock page with IGDB search.
- Purchase game flow, search and game detail JSON endpoints.
- Member registration endpoint.
- NES-style CSS theme.

### Changed
- UI migrated to dark navy NES theme with responsive grid.
- Games table extended with cover, magazine, acquired date (migration `002`).

## [0.1.0] - 2026-03-03

### Added
- Initial Go project structure with core models.
- PostgreSQL layer with `pgx/v5` and `Store` interface.
- Migration `001_initial_schema.sql`.
- IGDB API client with Twitch OAuth2.
- Landing page and games shelf with SSR.
- Graceful shutdown. GPL v3 license.
