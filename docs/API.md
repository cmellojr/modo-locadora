# API Reference

Modo Locadora uses server-rendered pages (HTML) and a few JSON endpoints. For authentication details, see [SECURITY.md](SECURITY.md).

## Pages (SSR)

### `GET /`

Landing page (Balcão) with login form and Wall of Shame (top late returners). Authenticated members are redirected to `/games`.

### `GET /games`

Without query parameters: 3-column layout with member mini-card + Wall of Shame (left), platform selection grid (center), activities feed + gaming almanac (right).

With `?platform=X`: simplified cartridge cards for that console — cover, title, copy count, availability status, golden star for completed games. Each card links to the game detail page.

### `GET /games/{id}`

Game detail page. Shows cover, title, platform, summary, source magazine, copy availability, total rentals, top renter, current renter, and acquired date. Logged-in members see the [ALUGAR] button if copies are available.

Query parameter: `error=em_debito` shows debt warning.

### `GET /carteirinha`

Digital membership card. Requires authentication. Shows membership number, progression title (Sócio Novato / Prata / Ouro / Dono da Calçada), profile, rental stats, status, password notebook, and active rentals with self-return (verdict selection).

Query parameter: `success` shows notification.

### `GET /admin/stock`

IGDB search and game acquisition page. Requires admin role. Query parameters: `q`, `magazine`, `selected`, `success`.

### `GET /admin/inventory`

Full catalog table with edit buttons and game health indicators (Cartucho Novo, Clássico Eterno, Precisa Soprar, Fita Gasta). Requires admin role. Query parameter: `success`.

### `GET /admin/edit/{id}`

Game edit form with cover upload (multipart) and cover display mode selector. Shows rental history (last 5 records). Requires admin role.

### `GET /admin/returns`

Active rentals dashboard with return buttons. Requires admin role. Query parameter: `success`.

---

## Form Endpoints

### `POST /login`

Authenticate and set session cookie. Content-Type: `application/x-www-form-urlencoded`.

| Field | Description |
|-------|-------------|
| `profile_name` | Member profile name |
| `password` | Member password |

**Success:** 303 redirect to `/games`, sets `session_member` cookie.

### `POST /rent`

Rent a game. Requires authentication.

| Field | Description |
|-------|-------------|
| `game_id` | Game UUID |

**Success:** 303 redirect to `/games/{id}`. Members in debt are redirected with `?error=em_debito`.

### `POST /carteirinha/notes`

Save password notebook. Requires authentication.

| Field | Description |
|-------|-------------|
| `notes` | Password notes text |

**Success:** 303 redirect to `/carteirinha?success=1`.

### `POST /carteirinha/return`

Self-return a rental with verdict. Requires authentication.

| Field | Description |
|-------|-------------|
| `rental_id` | Rental UUID |
| `verdict` | Play status: `zerei`, `joguei_um_pouco`, or `desisti` |

**Success:** 303 redirect to `/carteirinha?success=devolucao`. Fires an activity event based on the verdict.

### `POST /carteirinha/redeem`

Clear member's debt status. Requires authentication. No fields needed.

**Success:** 303 redirect to `/carteirinha?success=redencao`.

### `POST /admin/purchase`

Add a game from IGDB to the catalog. Requires admin role. Creates a `game_copy` atomically.

| Field | Description |
|-------|-------------|
| `title` | Game title |
| `igdb_id` | IGDB game ID |
| `platform` | Platform name (defaults to "N/A") |
| `summary` | Game description |
| `cover_url` | Cover image URL |
| `magazine` | Source magazine label |

**Success:** 303 redirect to `/admin/edit/{id}`.

### `POST /admin/update-game`

Update game details. Requires admin role. Content-Type: `multipart/form-data` (supports cover file upload).

| Field | Description |
|-------|-------------|
| `id` | Game UUID |
| `title` | Game title |
| `platform` | Platform name |
| `summary` | Description |
| `magazine` | Source magazine |
| `cover_url` | Existing cover URL (hidden, fallback) |
| `cover_display` | CSS object-fit mode: `cover` (default), `contain`, or `fill` |
| `cover_file` | Cover image file upload (optional) |

**Success:** 303 redirect to `/admin/inventory?success={title}`.

### `POST /admin/return-game`

Process a game return. Requires admin role.

| Field | Description |
|-------|-------------|
| `rental_id` | Rental UUID |

**Success:** 303 redirect to `/admin/returns?success=Fita+devolvida`.

---

## JSON API

### `POST /members`

Register a new member.

```json
{
  "profile_name": "Player1",
  "email": "player1@locadora.com",
  "password": "secret123",
  "favorite_console": "SNES"
}
```

**Response** `201 Created`: Member object with auto-assigned `MembershipNumber`. `PasswordHash` is always empty.

### `GET /search?q={query}`

Search IGDB database. Returns up to 10 results with name, summary, cover, and platforms.
