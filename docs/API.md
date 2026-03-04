# API Reference

Modo Locadora exposes a mix of server-rendered pages (HTML) and JSON API endpoints.

## Pages (SSR)

These routes return HTML rendered via Go's `html/template`.

### `GET /`

**Landing page (Balcao).** Displays the login form with profile name and password fields.

- No authentication required.

### `GET /games`

**Games shelf (Prateleira).** Displays the catalog with real-time availability status.

- Shows the authenticated member's name if logged in, otherwise "Visitante".
- Available games show [ALUGAR] button (logged in) or [DISPONIVEL] (not logged in).
- Rented games show [ALUGADO] with "Com o Socio: Nome" tag.
- Falls back to mock data if the database is unavailable or empty.
- Includes link to `/carteirinha` for logged-in members.

### `GET /carteirinha`

**Membership card (Carteirinha).** Displays the member's digital membership card.

- **Requires authentication** (`RequireAuth` middleware).
- Shows: membership number (`1991-XXX`), profile name, email, favorite console, join date.
- Redirects to `/` if not authenticated.

### `GET /admin/stock`

**Admin stock management page.** Search the IGDB database and add games to the catalog.

- **Requires authentication** and **admin role** (`ADMIN_EMAIL`).
- Accepts optional query parameters:
  - `q` — Search term for IGDB.
  - `magazine` — Magazine/edition label to associate with the purchase.
  - `selected` — IGDB game ID to show confirmation form.

### `GET /admin/inventory`

**Admin catalog listing.** Full table of all games in the database with edit buttons.

- **Requires authentication** and **admin role**.
- Shows: cover, title, platform, magazine for each game.
- Each row has an [Editar] link to `/admin/edit/{id}`.
- Accepts optional query parameter:
  - `success` — Game title to show in success balloon after an edit.

### `GET /admin/edit/{id}`

**Admin game edit form.** Pre-filled form for editing a game's details.

- **Requires authentication** and **admin role**.
- Editable fields: title, platform, source magazine, summary, cover URL.
- Path parameter: `id` (UUID of the game).

### `GET /admin/returns`

**Admin returns dashboard (Balcao de Devolucoes).** Lists all active (unreturned) rentals.

- **Requires authentication** and **admin role**.
- Shows: cover thumbnail, game title, member name, rental date for each active rental.
- Each row has a [Devolver] button (POST form).
- Accepts optional query parameter:
  - `success` — Message shown in success balloon after a return.

---

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
  "MembershipNumber": "1991-001",
  "Address": "",
  "Phone": "",
  "JoinedAt": "2026-03-03T12:00:00Z"
}
```

> Note: `PasswordHash` is always empty in the response for security. `MembershipNumber` is auto-assigned sequentially.

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

### `POST /rent`

**Rent a game.**

- **Auth:** Requires authentication (`RequireAuth` middleware).
- **Content-Type:** `application/x-www-form-urlencoded`

**Form fields:**

| Field     | Type   | Required | Description         |
|-----------|--------|----------|---------------------|
| `game_id` | UUID   | Yes      | The game's unique ID|

**Success:** Redirects to `/games` with `303 See Other`. The game will now show as rented by the member.

**Errors:**

| Status | Reason                          |
|--------|---------------------------------|
| 303    | Not authenticated (redirect to `/`) |
| 400    | Invalid game ID                 |
| 500    | No available copy or DB error   |
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
- Returns raw IGDB results as JSON (up to 10 results).

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
    },
    "platforms": [
      {"id": 19, "name": "Super Nintendo Entertainment System", "abbreviation": "SNES"}
    ]
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
- A `game_copy` record is automatically created in the same transaction.

**Form fields:**

| Field      | Type   | Required | Description                    |
|------------|--------|----------|--------------------------------|
| `title`    | string | Yes      | Game title                     |
| `igdb_id`  | string | Yes      | IGDB game ID                   |
| `platform` | string | No       | Platform name (default: "N/A") |
| `summary`  | string | No       | Game description               |
| `cover_url`| string | No       | Cover image URL from IGDB      |
| `magazine` | string | No       | Source magazine/edition label   |

**Success:** Redirects to `/admin/edit/{id}` with `303 See Other` for immediate editing.

**Errors:**

| Status | Reason                 |
|--------|------------------------|
| 403    | Not an admin           |
| 503    | Database not configured|

---

### `POST /admin/update-game`

**Update a game's details.**

- **Auth:** Requires admin role.
- **Content-Type:** `application/x-www-form-urlencoded`

**Form fields:**

| Field      | Type   | Required | Description                    |
|------------|--------|----------|--------------------------------|
| `id`       | UUID   | Yes      | Game's unique ID               |
| `title`    | string | Yes      | Game title                     |
| `platform` | string | Yes      | Platform name                  |
| `summary`  | string | No       | Game description               |
| `magazine` | string | No       | Source magazine/edition label   |
| `cover_url`| string | No       | Cover image URL (optional update)|

**Success:** Redirects to `/admin/inventory?success={title}` with `303 See Other`.

**Errors:**

| Status | Reason                 |
|--------|------------------------|
| 400    | Invalid UUID format    |
| 403    | Not an admin           |
| 404    | Game not found         |
| 503    | Database not configured|

---

### `POST /admin/return-game`

**Process a game return (check-in).**

- **Auth:** Requires admin role.
- **Content-Type:** `application/x-www-form-urlencoded`

**Form fields:**

| Field      | Type   | Required | Description              |
|------------|--------|----------|--------------------------|
| `rental_id`| UUID   | Yes      | The rental's unique ID   |

**Success:** Redirects to `/admin/returns?success=Fita+devolvida` with `303 See Other`. The game copy is marked as available.

**Errors:**

| Status | Reason                 |
|--------|------------------------|
| 400    | Invalid rental ID      |
| 403    | Not an admin           |
| 500    | Database error         |
| 503    | Database not configured|

---

## Authentication Flow

```
1. POST /members      -> Register (bcrypt hash stored, membership number assigned)
2. POST /login        -> Validate credentials -> Set signed cookie (member UUID)
3. GET  /games        -> Cookie verified -> Show rental buttons if logged in
4. GET  /carteirinha  -> Cookie verified -> Show membership card
5. POST /rent         -> Cookie verified -> Create rental record
6. GET  /admin/*      -> Cookie verified -> Email checked against ADMIN_EMAIL
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
