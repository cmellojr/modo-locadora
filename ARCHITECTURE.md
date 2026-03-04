# Modo Locadora - System Architecture

## 1. Vision & Concept
Modo Locadora is a retro-gaming session manager focused on scarcity, community, and nostalgia.
It emulates the experience of 90s Brazilian video rental stores ("locadoras").

## 2. Tech Stack
- **Backend:** Go (Golang) 1.24+
- **Database:** PostgreSQL 15+ (using `pgx/v5` driver)
- **Frontend:** Server-Side Rendering (SSR) with `html/template`
- **Styling:** [NES.css](https://nostalgic-css.github.io/NES.css/) + Google Fonts ("Press Start 2P")
- **External API:** IGDB (Internet Game Database) for game metadata
- **Security:** bcrypt password hashing, HMAC-SHA256 cookie signing, role-based middleware

## 3. Design Principles (The "Nostalgic" Rulebook)
- **Clean & Static:** No heavy JavaScript, no ads, no trackers.
- **Visuals:** Use `image-rendering: pixelated` for game covers.
- **Language:** Code, routes, and database must be in **English**. The UI (templates) must be in **Portuguese (BR)**.
- **Copyleft:** Licensed under **GPL v3**.
- **Scarcity by Design:** Each game has a limited number of physical copies (cartridges). If all copies are rented, the game is unavailable.

## 4. Key Entities
- **Member:** A registered store member with a sequential membership number (`1991-XXX`), profile name, email, hashed password, and favorite console.
- **Game:** Metadata record fetched from IGDB, including title, platform, cover URL, summary, and source magazine.
- **GameCopy (Tape/Instance):** Each game has physical-like copies with a status (`available` or `rented`). A copy is automatically created when a game is added to the catalog.
- **Rental:** Tracks which member has which game copy, with rental and due dates. Active rentals block availability.

## 5. Routing Strategy

### Public Routes
- `GET /` — Landing page (Balcao/Login)
- `POST /login` — Authentication (profile name + password)
- `POST /members` — Member registration (JSON API)
- `GET /search?q=` — IGDB search (JSON API)
- `GET /games/{id}` — Single game details (JSON API)

### Authenticated Routes (RequireAuth)
- `GET /games` — The Shelf (Prateleira) with rental buttons
- `GET /carteirinha` — Member's digital membership card
- `POST /rent` — Rent a game

### Admin Routes (RequireAdmin)
- `GET /admin/stock` — IGDB search & add games to catalog
- `POST /admin/purchase` — Confirm game purchase from IGDB
- `GET /admin/inventory` — Full catalog listing with edit links
- `GET /admin/edit/{id}` — Edit game details form
- `POST /admin/update-game` — Save game edits
- `GET /admin/returns` — Active rentals dashboard (check-in counter)
- `POST /admin/return-game` — Process a game return

## 6. Database Schema

### Tables
- `members` — id, profile_name, email, password_hash, favorite_console, membership_number, address, phone, joined_at
- `games` — id, title, igdb_id, platform, summary, cover_url, source_magazine, acquired_at
- `game_copies` — id, game_id, status (`available`/`rented`)
- `rentals` — id, game_copy_id, member_id, rented_at, due_at, returned_at

### Sequences
- `membership_seq` — Generates sequential membership numbers (format: `1991-001`, `1991-002`, ...)

### Migrations
1. `001_initial_schema.sql` — Base tables (members, games, game_copies, rentals)
2. `002_update_games_table.sql` — Add cover_url, source_magazine, acquired_at
3. `003_membership_and_rental_support.sql` — Add membership fields, backfill data, auto-create copies

## 7. Authentication & Authorization

```
POST /members    -> Register (bcrypt hash stored, membership number assigned)
POST /login      -> Validate credentials -> Set HMAC-signed cookie (member UUID)
GET  /games      -> Cookie verified -> Show rental buttons if logged in
GET  /carteirinha -> Cookie verified -> Show membership card
POST /rent       -> Cookie verified -> Create rental record
GET  /admin/*    -> Cookie verified -> Email checked against ADMIN_EMAIL
```

## 8. Rental Flow

```
1. Member browses /games -> sees [ALUGAR] button on available games
2. Member clicks [ALUGAR] -> POST /rent creates rental, marks copy as rented
3. Game now shows "ALUGADO - Com o Socio: Nome" on the shelf
4. Admin visits /admin/returns -> sees all active rentals
5. Admin clicks [Devolver] -> POST /admin/return-game marks copy as available
6. Game becomes available again on the shelf
```

## 9. Templates

| Template | Route | Description |
|----------|-------|-------------|
| `index.html` | `GET /` | Login form (Balcao) |
| `games.html` | `GET /games` | Game shelf with rental status |
| `carteirinha.html` | `GET /carteirinha` | Member card with 1991-XXX number |
| `admin_stock.html` | `GET /admin/stock` | IGDB search & game purchase |
| `admin_inventory.html` | `GET /admin/inventory` | Full catalog table with edit links |
| `admin_edit.html` | `GET /admin/edit/{id}` | Game edit form |
| `admin_returns.html` | `GET /admin/returns` | Active rentals check-in counter |
