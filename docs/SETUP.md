# Development Setup

Step-by-step guide for setting up Modo Locadora locally.

## Prerequisites

| Tool       | Version | Purpose                       |
|------------|---------|-------------------------------|
| Go         | 1.24+   | Backend language              |
| Docker     | 20+     | PostgreSQL container          |
| Git        | 2.x     | Version control               |

A local PostgreSQL 15+ installation works as an alternative to Docker.

## 1. Clone and Configure

```bash
git clone https://github.com/cmellojr/modo-locadora.git
cd modo-locadora
cp .env.example .env
```

Edit `.env` with your values:

```env
# IGDB API — Get credentials at https://dev.twitch.tv/console
TWITCH_CLIENT_ID=your_client_id
TWITCH_CLIENT_SECRET=your_client_secret

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=tio_da_locadora
DB_PASSWORD=sopre_a_fita
DB_NAME=modo_locadora
DATABASE_URL=postgres://tio_da_locadora:sopre_a_fita@localhost:5432/modo_locadora?sslmode=disable

# Security
COOKIE_SECRET=generate-a-random-secret-here-min-32-chars
ADMIN_EMAIL=your_admin_email@example.com
```

### Getting IGDB Credentials

1. Create an account at [Twitch Developer Console](https://dev.twitch.tv/console).
2. Register a new application (any category).
3. Copy the **Client ID** and generate a **Client Secret**.

## 2. Start the Database

### Option A: Docker Compose (recommended)

```bash
docker compose up -d
```

### Option B: Local PostgreSQL

```bash
createdb -U your_user modo_locadora
```

## 3. Run Migrations

Migrations are in `internal/database/migrations/` and must be applied in order:

```bash
psql $DATABASE_URL -f internal/database/migrations/001_initial_schema.sql
psql $DATABASE_URL -f internal/database/migrations/002_update_games_table.sql
psql $DATABASE_URL -f internal/database/migrations/003_membership_and_rental_support.sql
psql $DATABASE_URL -f internal/database/migrations/004_password_notes.sql
```

On Windows with Docker:

```bash
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/001_initial_schema.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/002_update_games_table.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/003_membership_and_rental_support.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/004_password_notes.sql
```

### Migration Summary

| Migration | Description |
|-----------|-------------|
| `001_initial_schema.sql` | Base tables: `members`, `games`, `game_copies`, `rentals` |
| `002_update_games_table.sql` | Adds `cover_url`, `source_magazine`, `acquired_at` to `games` |
| `003_membership_and_rental_support.sql` | Adds membership fields, `membership_seq` sequence, auto-creates copies |
| `004_password_notes.sql` | Adds `password_notes` field to `members` |

## 4. Build and Run

```bash
go build -o modo-locadora ./cmd/server
./modo-locadora
```

Or run directly:

```bash
go run ./cmd/server
```

The server starts on `http://localhost:8080`. Set the `PORT` environment variable to change it.

## 5. Create Your First Member

```bash
curl -X POST http://localhost:8080/members \
  -H "Content-Type: application/json" \
  -d '{
    "profile_name": "Tio da Locadora",
    "email": "your_admin_email@example.com",
    "password": "your_password",
    "favorite_console": "Mega Drive"
  }'
```

The email must match `ADMIN_EMAIL` for admin access. A membership number (`1991-001`) is auto-assigned.

## 6. Verify

| Check | Expected |
|-------|----------|
| `http://localhost:8080` | Login page loads |
| Login with member credentials | Redirects to `/games` |
| `/carteirinha` (logged in) | Membership card with `1991-XXX` |
| `/admin/stock` (as admin) | IGDB search page |
| `/admin/inventory` (as admin) | Catalog table |
| `/admin/returns` (as admin) | Active rentals list |

## Troubleshooting

### "No DATABASE_URL provided"

Login, rental, and admin features require the database. Ensure `DATABASE_URL` is set in `.env`.

### "COOKIE_SECRET not set"

The server falls back to an insecure default in development. For production, set a strong random value (at least 32 characters).

### "ADMIN_EMAIL not set"

Admin routes (`/admin/*`) are inaccessible without this. Set it to the admin member's email.

### IGDB search returns no results

Verify Twitch credentials:

```bash
curl -X POST "https://id.twitch.tv/oauth2/token?client_id=YOUR_ID&client_secret=YOUR_SECRET&grant_type=client_credentials"
```

### Port 8080 already in use

```bash
# Linux/Mac
lsof -ti:8080 | xargs kill -9

# Windows
netstat -ano | findstr :8080
taskkill /F /PID <pid>
```
