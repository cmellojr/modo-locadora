# Modo Locadora - System Architecture

## Vision

Retro-gaming session manager emulating 90s Brazilian video rental stores ("locadoras"). Scarcity is a core design principle: each game has limited physical copies, and if all are rented, the game is unavailable. Overdue rentals are auto-returned with reputation penalties.

## Entity Model

```
Member (1991-XXX)
  ├── status: active | em_debito
  ├── late_count: permanent penalty counter
  └── password_notes: personal game codes notebook

Game (IGDB metadata)
  ├── platform, summary, cover_url, source_magazine
  └── GameCopy (1:N)
        ├── status: available | rented
        └── Rental (1:N)
              ├── member_id, rented_at, due_at (3 days)
              └── returned_at (NULL = active)
```

## Rental Flow

```
1. Member browses /games → selects console → selects game → /games/{id}
2. Clicks [ALUGAR] → POST /rent → copy marked rented, rental created (3-day due)
3. Game detail shows "ALUGADO - Com o Sócio: Nome"
4. Admin visits /admin/returns → clicks [Devolver] → copy available again
5. If overdue: background job auto-returns, member gets em_debito + late_count++
6. Member can redeem via POST /carteirinha/redeem
```

## Navigation Map

```
GET /                     → Login (Balcão) — redirects to /games if authenticated
GET /games                → Platform selection grid (Mega Drive, SNES, ...)
GET /games?platform=X     → Cartridge cards for that console
GET /games/{id}           → Game detail (stats, rent button)
GET /carteirinha          → Membership card + password notebook
GET /admin/stock          → IGDB search & acquisition
GET /admin/inventory      → Catalog table with edit links
GET /admin/edit/{id}      → Edit game (cover upload, metadata)
GET /admin/returns        → Active rentals check-in
```

## Templates

| Template | Route | Page |
|----------|-------|------|
| `index.html` | `GET /` | Login + Wall of Shame |
| `platforms.html` | `GET /games` | Console selection grid |
| `games.html` | `GET /games?platform=X` | Cartridge shelf (simplified cards) |
| `game_detail.html` | `GET /games/{id}` | Game detail + rental stats |
| `carteirinha.html` | `GET /carteirinha` | Membership card + notebook |
| `admin_stock.html` | `GET /admin/stock` | IGDB search & acquisition |
| `admin_inventory.html` | `GET /admin/inventory` | Catalog table |
| `admin_edit.html` | `GET /admin/edit/{id}` | Game edit form |
| `admin_returns.html` | `GET /admin/returns` | Returns counter |

## Deployment

Multi-stage Docker build (`golang:1.24-alpine` → `alpine:3.21`). Docker Compose orchestrates app + PostgreSQL with health checks. Two volumes: `postgres_data` (DB) and `covers_data` (uploaded covers).

For build commands, package structure, DB schema, and conventions, see [CLAUDE.md](CLAUDE.md).
