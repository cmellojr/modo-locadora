# Development Setup

Step-by-step guide for setting up Modo Locadora locally.

## Prerequisites

| Tool       | Version | Purpose                       |
|------------|---------|-------------------------------|
| Docker     | 20+     | App + PostgreSQL containers   |
| Git        | 2.x     | Version control               |

For local development without Docker, you also need **Go 1.24+** and a **PostgreSQL 15+** instance.

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

# Database (used by Docker Compose)
DB_USER=tio_da_locadora
DB_PASSWORD=sopre_a_fita
DB_NAME=modo_locadora

# Security
COOKIE_SECRET=generate-a-random-secret-here-min-32-chars
ADMIN_EMAIL=your_admin_email@example.com
```

### Getting IGDB Credentials

1. Create an account at [Twitch Developer Console](https://dev.twitch.tv/console).
2. Register a new application (any category).
3. Copy the **Client ID** and generate a **Client Secret**.

## 2. Start with Docker (recommended)

```bash
docker compose up -d --build
```

This starts both the Go app and PostgreSQL. The app auto-connects to the database. Access at `http://localhost:8080`.

Migrations are applied manually (see step 3).

## 3. Run Migrations

Migrations are in `internal/database/migrations/` and must be applied in order.

### Via Docker container:

```bash
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/001_initial_schema.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/002_update_games_table.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/003_membership_and_rental_support.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/004_password_notes.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/005_auto_return_reputation.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/006_activities_feed.sql
```

### Via psql directly:

```bash
psql $DATABASE_URL -f internal/database/migrations/001_initial_schema.sql
psql $DATABASE_URL -f internal/database/migrations/002_update_games_table.sql
psql $DATABASE_URL -f internal/database/migrations/003_membership_and_rental_support.sql
psql $DATABASE_URL -f internal/database/migrations/004_password_notes.sql
psql $DATABASE_URL -f internal/database/migrations/005_auto_return_reputation.sql
psql $DATABASE_URL -f internal/database/migrations/006_activities_feed.sql
```

### Quick setup with seed data:

```bash
go run ./cmd/server --seed
```

This applies all migrations (001-006) and populates the database with sample games, members, and rental history. Test credentials: `MegaDriveKid` / `sega1991`, `Devedor` / `atrasado123`, `Novato` / `novato2026`.

### Migration Summary

| Migration | Description |
|-----------|-------------|
| `001_initial_schema.sql` | Base tables: `members`, `games`, `game_copies`, `rentals` |
| `002_update_games_table.sql` | Adds `cover_url`, `source_magazine`, `acquired_at` to `games` |
| `003_membership_and_rental_support.sql` | Adds membership fields, `membership_seq` sequence, auto-creates copies |
| `004_password_notes.sql` | Adds `password_notes` field to `members` |
| `005_auto_return_reputation.sql` | Adds `status` and `late_count` fields to `members` |
| `006_activities_feed.sql` | Creates `activities` table for event feed |
| `007_seed_initial_data.sql` | Sample data (applied via `--seed` flag, not manually) |

## 4. Local Development (without Docker for the app)

If you prefer running the Go server locally while keeping PostgreSQL in Docker:

```bash
docker compose up -d db       # start only PostgreSQL
go run ./cmd/server            # run the Go server locally
```

Set `DATABASE_URL` in `.env` to point to `localhost:5432` (not `db:5432`).

```bash
# Build binary
go build -o modo-locadora ./cmd/server

# Static analysis
go vet ./...
```

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
| Login with member credentials | Redirects to `/games` (platform grid) |
| Click a platform | Shows cartridge cards |
| Click a cartridge | Shows game detail page |
| `/carteirinha` (logged in) | Membership card with `1991-XXX` |
| `/admin/stock` (as admin) | IGDB search page |

## Troubleshooting

### "No DATABASE_URL provided"
Ensure `DATABASE_URL` is set in `.env`. When using Docker Compose for the full stack, it's set automatically via `docker-compose.yml`.

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
