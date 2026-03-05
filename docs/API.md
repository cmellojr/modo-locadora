# API Reference

Modo Locadora exposes server-rendered pages (HTML) and JSON API endpoints. For authentication and security details, see [SECURITY.md](SECURITY.md).

## Pages (SSR)

All page routes return HTML rendered via Go's `html/template`.

### `GET /`

Landing page (Balcao) with login form. No authentication required.

### `GET /games`

Game shelf (Prateleira) with real-time availability. Shows [ALUGAR] button for logged-in members, [DISPONIVEL] badge otherwise. Rented games display the renter's name.

### `GET /carteirinha`

Digital membership card. Requires authentication. Shows membership number, profile, and password notebook.

### `GET /admin/stock`

IGDB search and game acquisition page. Requires admin role.

Query parameters: `q` (search term), `magazine` (edition label), `selected` (IGDB game ID for confirmation).

### `GET /admin/inventory`

Full catalog table with edit buttons. Requires admin role.

Query parameter: `success` (game title for success notification).

### `GET /admin/edit/{id}`

Game edit form. Requires admin role. Path parameter: `id` (game UUID).

### `GET /admin/returns`

Active rentals dashboard with return buttons. Requires admin role.

Query parameter: `success` (notification message after return).

---

## JSON API

### `POST /members`

Register a new member. No authentication required.

**Request:**

```json
{
  "profile_name": "Player1",
  "email": "player1@locadora.com",
  "password": "secret123",
  "favorite_console": "SNES"
}
```

**Response** `201 Created`:

```json
{
  "ID": "a1b2c3d4-...",
  "ProfileName": "Player1",
  "Email": "player1@locadora.com",
  "PasswordHash": "",
  "FavoriteConsole": "SNES",
  "MembershipNumber": "1991-001",
  "JoinedAt": "2026-03-03T12:00:00Z"
}
```

`PasswordHash` is always empty in responses. `MembershipNumber` is auto-assigned.

| Status | Reason |
|--------|--------|
| 400 | Missing or invalid JSON body |
| 400 | Empty password |
| 503 | Database not configured |

---

### `POST /login`

Authenticate and set session cookie. Content-Type: `application/x-www-form-urlencoded`.

| Field | Required | Description |
|-------|----------|-------------|
| `profile_name` | Yes | Member profile name |
| `password` | Yes | Member password |

**Success:** 303 redirect to `/games`, sets `session_member` cookie.

| Status | Reason |
|--------|--------|
| 400 | Missing fields |
| 401 | Invalid credentials |
| 503 | Database not configured |

---

### `POST /rent`

Rent a game. Requires authentication. Content-Type: `application/x-www-form-urlencoded`.

| Field | Required | Description |
|-------|----------|-------------|
| `game_id` | Yes | Game UUID |

**Success:** 303 redirect to `/games`.

| Status | Reason |
|--------|--------|
| 303 | Not authenticated (redirect to `/`) |
| 400 | Invalid game ID |
| 500 | No available copy or DB error |
| 503 | Database not configured |

---

### `POST /carteirinha/notes`

Save password notebook. Requires authentication. Content-Type: `application/x-www-form-urlencoded`.

| Field | Required | Description |
|-------|----------|-------------|
| `notes` | Yes | Password notes text |

**Success:** 303 redirect to `/carteirinha`.

---

### `GET /games/{id}`

Retrieve a single game by UUID. No authentication required.

**Response** `200 OK`:

```json
{
  "ID": "a1b2c3d4-...",
  "Title": "Chrono Trigger",
  "IgdbID": "1234",
  "Platform": "SNES",
  "Summary": "An epic RPG...",
  "CoverURL": "https://images.igdb.com/...",
  "SourceMagazine": "Super Game Power #45",
  "AcquiredAt": "2026-03-03T12:00:00Z"
}
```

| Status | Reason |
|--------|--------|
| 400 | Invalid UUID |
| 404 | Not found |
| 503 | Database not configured |

---

### `GET /search?q={query}`

Search IGDB database. No authentication required. Returns up to 10 results.

**Response** `200 OK`:

```json
[
  {
    "id": 1234,
    "name": "Chrono Trigger",
    "summary": "An epic RPG...",
    "cover": { "id": 5678, "url": "//images.igdb.com/..." },
    "platforms": [{ "id": 19, "name": "Super Nintendo", "abbreviation": "SNES" }]
  }
]
```

| Status | Reason |
|--------|--------|
| 400 | Missing `q` parameter |
| 503 | IGDB credentials not configured |

---

### `POST /admin/purchase`

Add a game from IGDB to the catalog. Requires admin role. A `game_copy` is created atomically.

| Field | Required | Description |
|-------|----------|-------------|
| `title` | Yes | Game title |
| `igdb_id` | Yes | IGDB game ID |
| `platform` | No | Platform name |
| `summary` | No | Game description |
| `cover_url` | No | Cover image URL |
| `magazine` | No | Source magazine label |

**Success:** 303 redirect to `/admin/edit/{id}`.

---

### `POST /admin/update-game`

Update a game's details. Requires admin role.

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Game UUID |
| `title` | Yes | Game title |
| `platform` | Yes | Platform name |
| `summary` | No | Description |
| `magazine` | No | Source magazine |
| `cover_url` | No | Cover image URL |

**Success:** 303 redirect to `/admin/inventory?success={title}`.

---

### `POST /admin/return-game`

Process a game return. Requires admin role.

| Field | Required | Description |
|-------|----------|-------------|
| `rental_id` | Yes | Rental UUID |

**Success:** 303 redirect to `/admin/returns?success=Fita+devolvida`. Copy marked as available.
