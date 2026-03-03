# API Reference

Modo Locadora exposes a mix of server-rendered pages (HTML) and JSON API endpoints.

## Pages (SSR)

These routes return HTML rendered via Go's `html/template`.

### `GET /`

**Landing page (Balcao).** Displays the login form with profile name and password fields.

- No authentication required.

### `GET /games`

**Games shelf (Prateleira).** Displays the catalog split into "Releases" and "Catalog" sections.

- Shows the authenticated member's name if logged in, otherwise "Visitante".
- Falls back to mock data if the database is unavailable or empty.

### `GET /admin/stock`

**Admin stock management page.** Search the IGDB database and add games to the catalog.

- **Requires authentication** and **admin role** (`ADMIN_EMAIL`).
- Accepts optional query parameters:
  - `q` — Search term for IGDB.
  - `magazine` — Magazine/edition label to associate with the purchase.

## JSON API

These routes return `application/json` responses.

### `POST /members`

**Register a new member.**

- **Auth:** None required.

**Request body:**

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
  "JoinedAt": "2026-03-03T12:00:00Z"
}
```

> Note: `PasswordHash` is always empty in the response for security.

**Errors:**

| Status | Reason                        |
|--------|-------------------------------|
| 400    | Missing or invalid JSON body  |
| 400    | Password is empty             |
| 503    | Database not configured       |

---

### `POST /login`

**Authenticate a member and set a session cookie.**

- **Content-Type:** `application/x-www-form-urlencoded`

**Form fields:**

| Field          | Type   | Required | Description         |
|----------------|--------|----------|---------------------|
| `profile_name` | string | Yes      | Member profile name |
| `password`     | string | Yes      | Member password     |

**Success:** Redirects to `/games` with a `303 See Other` and sets the `session_member` cookie (HMAC-signed, HttpOnly, SameSite=Strict, 7-day expiry).

**Errors:**

| Status | Reason                          |
|--------|---------------------------------|
| 400    | Missing profile name or password|
| 401    | Invalid credentials             |
| 503    | Database not configured         |

---

### `GET /games/{id}`

**Retrieve a single game by UUID.**

- **Auth:** None required.

**Path parameters:**

| Parameter | Type | Description            |
|-----------|------|------------------------|
| `id`      | UUID | The game's unique ID   |

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

**Errors:**

| Status | Reason               |
|--------|----------------------|
| 400    | Invalid UUID format  |
| 404    | Game not found       |
| 503    | Database not configured |

---

### `GET /search?q={query}`

**Search games on the IGDB database.**

- **Auth:** None required.
- Returns raw IGDB results as JSON.

**Query parameters:**

| Parameter | Type   | Required | Description          |
|-----------|--------|----------|----------------------|
| `q`       | string | Yes      | Search term          |

**Response** `200 OK`:

```json
[
  {
    "id": 1234,
    "name": "Chrono Trigger",
    "summary": "An epic RPG...",
    "first_release_date": 795830400,
    "cover": {
      "id": 5678,
      "url": "//images.igdb.com/igdb/image/upload/t_thumb/co1v9x.jpg"
    }
  }
]
```

**Errors:**

| Status | Reason                            |
|--------|-----------------------------------|
| 400    | Missing `q` parameter             |
| 503    | IGDB credentials not configured   |

---

### `POST /admin/purchase`

**Add a game from IGDB to the local catalog.**

- **Auth:** Requires admin role.
- **Content-Type:** `application/x-www-form-urlencoded`

**Form fields:**

| Field      | Type   | Required | Description                    |
|------------|--------|----------|--------------------------------|
| `title`    | string | Yes      | Game title                     |
| `igdb_id`  | string | Yes      | IGDB game ID                   |
| `summary`  | string | No       | Game description (PT-BR)       |
| `cover_url`| string | No       | Cover image URL from IGDB      |
| `magazine` | string | No       | Source magazine/edition label   |

**Success:** Redirects to `/games` with `303 See Other`.

**Errors:**

| Status | Reason                 |
|--------|------------------------|
| 403    | Not an admin           |
| 503    | Database not configured|

## Authentication Flow

```
1. POST /members          → Register (bcrypt-hashed password stored)
2. POST /login            → Validate credentials → Set signed cookie
3. GET  /games            → Cookie verified → Member name displayed
4. GET  /admin/stock      → Cookie verified → Email checked against ADMIN_EMAIL
```

## Cookie Details

| Property   | Value                           |
|------------|---------------------------------|
| Name       | `session_member`                |
| Value      | `{member_uuid}.{hmac_sha256}`   |
| HttpOnly   | `true`                          |
| SameSite   | `Strict`                        |
| MaxAge     | 604800 (7 days)                 |
| Path       | `/`                             |
