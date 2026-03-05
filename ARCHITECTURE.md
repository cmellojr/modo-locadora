# Modo Locadora - System Architecture

## Vision

Modo Locadora is a retro-gaming session manager that emulates the experience of 90s Brazilian video rental stores ("locadoras"). Scarcity is a core design principle: each game has limited physical copies, and if all are rented, the game is unavailable.

## Tech Stack

- **Backend:** Go 1.24+ with standard library `net/http.ServeMux` (method-pattern routing)
- **Database:** PostgreSQL 15+ via `pgx/v5` connection pool
- **Frontend:** Server-Side Rendering with `html/template`, no JavaScript
- **Styling:** [NES.css](https://nostalgic-css.github.io/NES.css/) 2.3.0 + Press Start 2P font
- **External API:** IGDB (via Twitch OAuth2) for game metadata and covers
- **Security:** bcrypt password hashing, HMAC-SHA256 cookie signing, role-based middleware

## Design Principles

- **Clean & Static:** No JavaScript, no ads, no trackers.
- **Pixelated:** `image-rendering: pixelated` for game covers.
- **Language Split:** Code, routes, and database in English. UI templates in Portuguese (BR).
- **Copyleft:** Licensed under GPL v3.
- **Scarcity by Design:** Limited copies per game. All rented = unavailable.

## Package Structure

```
cmd/server/main.go              Entrypoint: config, templates, pool, routes
internal/
  auth/auth.go                  HMAC-SHA256 cookie signing/verification
  config/config.go              .env loader (godotenv)
  database/
    store.go                    Store interface (all DB operations)
    postgres.go                 PostgreSQL implementation (pgx/v5 pool + transactions)
    migrations/                 SQL migrations (001-004)
  handlers/handler.go           All HTTP handlers (Handler struct: Store + templates + config)
  igdb/igdb.go                  IGDB API client (Twitch OAuth2 token + game search)
  middleware/middleware.go       RequireAuth and RequireAdmin middleware
  models/                       Domain entities: Member, Game, GameCopy, Rental
web/
  static/css/retro.css          NES.css dark theme overrides and utility classes
  templates/                    Go HTML templates (7 pages, all PT-BR)
```

## Key Entities

- **Member:** Profile name, email, bcrypt-hashed password, sequential membership number (`1991-XXX`), favorite console, password notebook.
- **Game:** Metadata from IGDB — title, platform, cover URL, summary, source magazine.
- **GameCopy:** Physical-like instance of a game with status (`available` | `rented`). Auto-created when a game is added.
- **Rental:** Links a member to a game copy with rental/due/return dates.

## Request Flow

```
1. main.go registers routes on http.ServeMux with middleware wrappers
2. Middleware verifies session_member cookie via HMAC, injects member UUID into context
3. Handler reads context, calls Store interface methods, renders template with data struct
4. Store executes parameterized SQL via pgx connection pool
```

## Database Schema

4 tables + 1 sequence:

```
members          -> id, profile_name, email, password_hash, favorite_console,
                    membership_number, address, phone, password_notes, joined_at

games            -> id, title, igdb_id, platform, summary, cover_url,
                    source_magazine, acquired_at

game_copies      -> id, game_id, status (available | rented)

rentals          -> id, game_copy_id, member_id, rented_at, due_at, returned_at

membership_seq   -> generates sequential numbers (1991-001, 1991-002, ...)
```

Key relationship: `Game -> GameCopy -> Rental <- Member`

## Routing

### Public
| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Landing page (Balcao) with login form |
| POST | `/login` | Authentication |
| POST | `/members` | Registration (JSON API) |
| GET | `/search?q=` | IGDB search (JSON API) |
| GET | `/games/{id}` | Single game details (JSON API) |

### Authenticated (RequireAuth)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/games` | Game shelf with rental buttons |
| GET | `/carteirinha` | Digital membership card |
| POST | `/carteirinha/notes` | Save password notebook |
| POST | `/rent` | Rent a game |

### Admin (RequireAdmin)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/admin/stock` | IGDB search and acquisition |
| POST | `/admin/purchase` | Confirm game acquisition |
| GET | `/admin/inventory` | Full catalog listing |
| GET | `/admin/edit/{id}` | Game edit form |
| POST | `/admin/update-game` | Save game edits |
| GET | `/admin/returns` | Active rentals dashboard |
| POST | `/admin/return-game` | Process game return |

## Rental Flow

```
1. Member browses /games -> sees [ALUGAR] on available games
2. Clicks [ALUGAR] -> POST /rent creates rental, marks copy as rented
3. Game shows "ALUGADO - Com o Socio: Nome" on the shelf
4. Admin visits /admin/returns -> sees active rentals
5. Admin clicks [Devolver] -> POST /admin/return-game marks copy as available
6. Game becomes available on the shelf again
```

## Templates

| Template | Route | Page |
|----------|-------|------|
| `index.html` | `GET /` | Login (Balcao) |
| `games.html` | `GET /games` | Game shelf with rental status |
| `carteirinha.html` | `GET /carteirinha` | Membership card + password notebook |
| `admin_stock.html` | `GET /admin/stock` | IGDB search and acquisition |
| `admin_inventory.html` | `GET /admin/inventory` | Catalog table with edit links |
| `admin_edit.html` | `GET /admin/edit/{id}` | Game edit form |
| `admin_returns.html` | `GET /admin/returns` | Returns check-in counter |

## Related Documentation

- [docs/SETUP.md](docs/SETUP.md) — Development environment setup
- [docs/API.md](docs/API.md) — Endpoint reference
- [docs/SECURITY.md](docs/SECURITY.md) — Security policy and practices
- [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) — Contribution guidelines and conventions
- [docs/CHANGELOG.md](docs/CHANGELOG.md) — Version history
