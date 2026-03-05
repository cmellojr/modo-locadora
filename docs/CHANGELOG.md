# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Password notebook** (`Caderno de Passwords`): Members can save game codes and passwords on their membership card. Migration `004_password_notes.sql`.
- **NES.css component expansion**: `nes-badge` for status indicators, `nes-progress` for copy availability bars, `nes-list` for statute rules, `nes-avatar` on membership card, `nes-dialog` for acquisition confirmation, `nes-text` for semantic colored text.
- **CLAUDE.md**: AI agent guidance file for Claude Code.
- **AGENTS.md**: AI agent guidance file following Jules/Google specification.

### Changed
- **Unified success notifications**: All pages now use `nes-balloon` + `nes-bcrikko` pattern instead of mixed approaches.
- **Font sizes increased globally**: Body 12px, labels 12px, titles 18px, buttons 12px for better readability.
- **Container width standardized**: All pages use 960px max-width via `.container-wrapper` class.
- **Cover thumbnails enlarged**: Inventory 90px, returns 70px, shelf 120x155px, edit 160px.
- **Game status indicators**: Replaced disabled buttons with semantic `nes-badge` elements (`is-primary` for available, `is-error` for rented).
- **Platform display**: Uses native `nes-badge is-primary` instead of custom styled spans.
- **Inline styles consolidated**: ~15+ repeated inline styles moved to reusable CSS classes in `retro.css` (`.btn-nav`, `.btn-sm`, `.title-main`, `.title-sub`, `.footer-copyright`, `.nav-bar`, `.form-actions`, `.empty-state`, `.success-balloon`).
- **Landing page**: Statute rules use `nes-list is-disc`, columns bottom-aligned with flexbox.
- **Membership card**: Added `nes-avatar`, increased to 750px max-width.
- **Documentation rewrite**: README.md rewritten with nostalgic tone (PT-BR). All docs/ files revised to eliminate redundancy.

## [0.3.0] - 2026-03-04

### Added
- **Membership system**: Sequential membership numbers in `1991-XXX` format.
- **Carteirinha page** (`GET /carteirinha`): Digital membership card with member number, name, email, favorite console, and join date.
- **Rental system**: Members rent games from the shelf via `POST /rent`. Rented games show the member's name.
- **Returns dashboard** (`GET /admin/returns`): Admin page listing active rentals with return buttons.
- **Admin inventory** (`GET /admin/inventory`): Full catalog table with edit links.
- **Game edit** (`GET /admin/edit/{id}`): Form for editing game details (title, platform, summary, magazine).
- **Migration `003`**: Membership fields, `membership_seq` sequence, auto-created copies for existing games.
- `RequireAuth` middleware for member-only routes.
- `internal/auth` package with HMAC-SHA256 signed cookies.
- `internal/middleware` package with `RequireAuth` and `RequireAdmin`.
- bcrypt password hashing on registration.
- `.env.example` with placeholder values.
- `docs/` directory with project documentation.

### Changed
- Game shelf displays real-time availability with [ALUGAR], [DISPONIVEL], and [ALUGADO] states.
- `POST /admin/purchase` redirects to edit page for immediate translation.
- Login requires profile name + password (was name-only).
- Session cookie stores signed UUID instead of plain-text name.
- Cookie includes `MaxAge`, `SameSite=Strict`, `HttpOnly` flags.

## [0.2.0] - 2026-03-03

### Added
- Docker Compose configuration for PostgreSQL 15.
- Admin stock page (`/admin/stock`) with IGDB search.
- Purchase game flow (`POST /admin/purchase`).
- `GET /search` and `GET /games/{id}` JSON endpoints.
- `POST /members` registration endpoint.
- NES-style CSS theme (`retro.css`).
- `DATABASE_URL` environment variable support.

### Changed
- UI migrated to dark navy NES theme.
- Games shelf uses responsive CSS grid.
- Games table extended with `cover_url`, `source_magazine`, `acquired_at` (migration `002`).

## [0.1.0] - 2026-03-03

### Added
- Initial Go project structure.
- Core models: `Member`, `Game`, `GameCopy`, `Rental`.
- PostgreSQL layer with `pgx/v5` and `Store` interface.
- Migration `001_initial_schema.sql`.
- IGDB API client with Twitch OAuth2.
- Environment config loader (godotenv).
- Landing page and games shelf with SSR.
- Graceful shutdown with OS signal handling.
- `ARCHITECTURE.md`.
- GPL v3 license.
