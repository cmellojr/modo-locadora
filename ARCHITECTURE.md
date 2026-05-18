# Modo Locadora — Arquitetura do Sistema

## Visão Geral

O Modo Locadora é uma aplicação Go SSR que emula uma locadora brasileira de videogames. O domínio gira em torno de sócios, jogos, cópias físicas, aluguéis, vereditos, reputação e turmas.

Para setup, variáveis e comandos, veja [docs/setup.md](docs/setup.md). Para rotas detalhadas, veja [docs/api.md](docs/api.md).

## Componentes

| Área | Implementação |
|------|---------------|
| Entrypoint | `cmd/server/main.go` |
| HTTP | `net/http.ServeMux` com padrões `METHOD /path` |
| Handlers | `internal/handlers/handler.go` |
| Persistência | `internal/database.Store` + `internal/database/postgres.go` |
| Auth | `internal/auth` e `internal/middleware` |
| Jobs | `internal/jobs/overdue.go` |
| IGDB | `internal/igdb` |
| Almanaque | `internal/almanac` |
| Uploads | `internal/storage` (`LocalStorage` ou `GCSStorage`) |
| Templates | `web/templates` |
| Assets | `web/static` |

## Modelo de Domínio

```text
Member
  ├── membership_number: 1991-XXX
  ├── status: active | in_debt
  ├── late_count
  ├── password_notes
  └── MemberTitle calculado por histórico

Game
  ├── platform, summary, cover_url, cover_display, source_magazine
  └── GameCopy (1:N)
        ├── status: available | rented
        └── Rental (1:N)
              ├── rented_at, due_at, returned_at
              └── public_legacy: completed | enjoyed | quick_play | not_for_me | gave_up | auto_return

Activity
  └── event_type, member_name, game_title, created_at

Club
  ├── name, description, badge_url, website_url, created_by
  └── ClubMember (M:N com Member)
        └── role: admin | member
```

As tabelas e migrations ficam em `internal/database/migrations/`; o resumo operacional está em [docs/setup.md](docs/setup.md).

## Fluxo de Aluguel

1. Sócio navega por `/games`, escolhe plataforma e abre `/games/{id}`.
2. `POST /rent` seleciona uma cópia disponível em transação e cria o aluguel com prazo de 3 dias.
3. A cópia fica `rented`, reduzindo a disponibilidade do jogo.
4. A devolução pode acontecer por `POST /membership/return` com veredito ou por `POST /admin/return-game`.
5. O job `StartOverdueChecker` roda a cada 5 minutos; aluguéis vencidos são auto-devolvidos, recebem `auto_return` e penalizam o sócio com `in_debt`.
6. `POST /membership/redeem` limpa o débito, mantendo `late_count` como histórico.

## Templates

Todos os templates usam `web/templates/layout.html` como base.

| Template | Rota principal | Função |
|----------|----------------|--------|
| `index.html` | `GET /` | Balcão, login, boas-vindas e Painel da Vergonha |
| `platforms.html` | `GET /games` | Seleção de plataformas, feed e almanaque |
| `games.html` | `GET /games?platform=X` | Prateleira filtrada por console |
| `game_detail.html` | `GET /games/{id}` | Detalhe, disponibilidade e estatísticas |
| `membership.html` | `GET /membership` | Carteirinha, notas, aluguéis ativos e turmas |
| `admin_stock.html` | `GET /admin/stock` | Busca IGDB e aquisição |
| `admin_inventory.html` | `GET /admin/inventory` | Inventário com popularidade |
| `admin_edit.html` | `GET /admin/edit/{id}` | Edição, upload de capa e histórico |
| `admin_returns.html` | `GET /admin/returns` | Devoluções admin |
| `clubs.html` | `GET /clubs` | Listagem pública de turmas |
| `club_detail.html` | `GET /clubs/{id}` | Detalhe e membros da turma |
| `club_form.html` | `GET /clubs/new`, `GET /clubs/{id}/edit` | Criação e edição de turma |

## Uploads e Storage

Handlers de upload não gravam diretamente no filesystem. Eles usam `storage.StorageProvider`, criado em `cmd/server/main.go`:

- `APP_ENV=production`: usa `GCSStorage` e exige `STORAGE_BUCKET_NAME`.
- Qualquer outro valor: usa `LocalStorage` em `web/static/covers` e `web/static/clubs`.

O Docker Compose local mantém volumes `covers_data` e `clubs_data`. Para detalhes de configuração, veja [docs/setup.md](docs/setup.md).

## Deploy

O Dockerfile usa build multi-stage (`golang:1.24-alpine` -> `alpine:3.21`). O Compose local orquestra app + PostgreSQL com healthcheck e volumes persistentes.

Checklist de segurança e produção: [docs/security.md](docs/security.md).
