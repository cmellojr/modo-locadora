# Development Setup

This guide walks you through setting up the Modo Locadora project for local development.

## Prerequisites

| Tool       | Version | Purpose                        |
|------------|---------|--------------------------------|
| Go         | 1.24+   | Backend language               |
| PostgreSQL | 15+     | Database                       |
| Docker     | 20+     | Database container (optional)  |
| Git        | 2.x     | Version control                |

## 1. Clone the Repository

```bash
git clone https://github.com/cmellojr/modo-locadora.git
cd modo-locadora
```

## 2. Configure Environment Variables

Copy the example file and fill in your values:

```bash
cp .env.example .env
```

Edit `.env` with your settings:

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
COOKIE_SECRET=generate-a-random-secret-here
ADMIN_EMAIL=your_admin_email@example.com
```

### Getting IGDB Credentials

1. Create an account at [Twitch Developer Console](https://dev.twitch.tv/console).
2. Register a new application (any category).
3. Copy the **Client ID** and generate a **Client Secret**.
4. Paste them into your `.env` file.

## 3. Start the Database

### Option A: Docker Compose (recommended)

```bash
docker compose up -d
```

This starts a PostgreSQL 15 Alpine container on the port defined in `DB_PORT`.

### Option B: Local PostgreSQL

If you have PostgreSQL installed locally, create the database:

```bash
createdb -U your_user modo_locadora
```

## 4. Run Database Migrations

Migrations are in `internal/database/migrations/` and must be applied manually in order:

```bash
psql $DATABASE_URL -f internal/database/migrations/001_initial_schema.sql
psql $DATABASE_URL -f internal/database/migrations/002_update_games_table.sql
psql $DATABASE_URL -f internal/database/migrations/003_membership_and_rental_support.sql
```

On Windows with Docker:

```bash
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/001_initial_schema.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/002_update_games_table.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/003_membership_and_rental_support.sql
```

### What Each Migration Does

| Migration | Description |
|-----------|-------------|
| `001_initial_schema.sql` | Creates base tables: `members`, `games`, `game_copies`, `rentals` |
| `002_update_games_table.sql` | Adds `cover_url`, `source_magazine`, `acquired_at` columns to `games` |
| `003_membership_and_rental_support.sql` | Adds `membership_number`, `address`, `phone` to `members`; creates `membership_seq` sequence; backfills existing data; creates `game_copies` for existing games |

## 5. Install Dependencies

```bash
go mod download
```

## 6. Build and Run

```bash
go build -o modo-locadora ./cmd/server
./modo-locadora
```

Or run directly without building:

```bash
go run ./cmd/server
```

The server starts on `http://localhost:8080` by default. Set the `PORT` environment variable to change it.

## 7. Create Your First Member

The application requires a registered member to log in. Use the API to create one:

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

The member will automatically receive a membership number (e.g., `1991-001`).

To access admin routes, set the member's email to match the `ADMIN_EMAIL` value in `.env`.

## 8. Verify the Setup

| Check                        | Expected Result                          |
|------------------------------|------------------------------------------|
| `http://localhost:8080`      | Login page (Balcao) loads                |
| Login with member + password | Redirects to `/games` shelf              |
| `http://localhost:8080/games`| Games shelf with availability status     |
| `/carteirinha` (logged in)   | Membership card with `1991-XXX` number   |
| `/admin/stock` (as admin)    | IGDB search page loads                   |
| `/admin/inventory` (as admin)| Full catalog table with edit buttons     |
| `/admin/returns` (as admin)  | Active rentals list                      |
| `/admin/stock` (no login)    | Redirected to `/`                        |

## Project Structure

```
modo-locadora/
├── cmd/server/main.go              # Application entrypoint
├── internal/
│   ├── auth/auth.go                # Cookie signing (HMAC-SHA256)
│   ├── config/config.go            # .env loader (godotenv)
│   ├── database/
│   │   ├── store.go                # Store interface
│   │   ├── postgres.go             # PostgreSQL implementation
│   │   └── migrations/             # SQL migration files (001-003)
│   ├── handlers/handler.go         # HTTP handlers
│   ├── igdb/igdb.go                # IGDB API client
│   ├── middleware/middleware.go     # Auth & admin middleware
│   └── models/                     # Domain entities
│       ├── member.go
│       ├── game.go
│       ├── game_copy.go
│       └── rental.go
├── web/
│   ├── static/css/retro.css        # NES-style theme
│   └── templates/                  # Go HTML templates (PT-BR)
│       ├── index.html              # Login page (Balcao)
│       ├── games.html              # Game shelf with rental status
│       ├── carteirinha.html        # Membership card
│       ├── admin_stock.html        # IGDB search & purchase
│       ├── admin_inventory.html    # Catalog listing
│       ├── admin_edit.html         # Game edit form
│       └── admin_returns.html      # Active rentals check-in
├── docs/                           # Project documentation
├── docker-compose.yml              # PostgreSQL container
├── .env.example                    # Environment template
├── ARCHITECTURE.md                 # System architecture
└── go.mod                          # Go module definition
```

## Common Issues

### "No DATABASE_URL provided"

The server can start without a database (it uses mock data), but login, rental, and admin features require it. Make sure `DATABASE_URL` is set in your `.env`.

### "COOKIE_SECRET not set"

The server falls back to an insecure default in development. For production, always set a strong random value (at least 32 characters).

### "ADMIN_EMAIL not set"

Admin routes (`/admin/*`) will be inaccessible if `ADMIN_EMAIL` is not configured. Set it to the email address of the member who should have admin access.

### IGDB search returns no results

Verify your Twitch credentials are valid. You can test token retrieval manually:

```bash
curl -X POST "https://id.twitch.tv/oauth2/token?client_id=YOUR_ID&client_secret=YOUR_SECRET&grant_type=client_credentials"
```

### Port 8080 already in use

Kill the existing process:

```bash
# Linux/Mac
lsof -ti:8080 | xargs kill -9

# Windows
netstat -ano | findstr :8080
taskkill /F /PID <pid>
```
