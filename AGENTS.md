# AGENTS.md

Este arquivo é a fonte única de referência para agentes de IA de programação (Gemini CLI, Jules, Claude Code, Copilot, Cursor, etc.) que trabalham no projeto **Modo Locadora**.

## Visão Geral do Projeto

**Modo Locadora** é um gerenciador de sessões e backlog de jogos retrô, projetado para emular a experiência das locadoras brasileiras de videogame dos anos 1990.

- **Conceito central**: "Scarcity by design" — cada jogo possui um número limitado de cópias físicas.
- **Reputação**: membros precisam cuidar da própria reputação para evitar o "Wall of Shame".
- **Stack**: Go 1.24+, PostgreSQL 15, Server-Side Rendering (SSR) com `html/template`, NES.css para uma UI 8-bit. Sem JavaScript.
- **Licença**: GPL v3.

## Filosofia Central e Diretrizes de Design

1. **Divisão de idioma**:
   - **Código, colunas de banco de dados, rotas e comentários**: inglês.
   - **Texto de UI e conteúdo visível para usuários**: português (BR).
2. **Sem JavaScript**: abordagem estritamente zero-JS, totalmente SSR.
3. **Estética retrô**: siga o tema escuro do NES.css (`is-dark`, `#0A0E1A`) e a metáfora de "locadora de videogame" (por exemplo, "Sopro" para redenções).
4. **Progressão**: membros evoluem de `Sócio Novato` -> `Prata` -> `Ouro` -> `Dono da Calçada`.
5. **Escassez**: nunca ignore a lógica de limite de cópias; se as cópias forem 0, o jogo está "Alugado".

## Build e Execução

O projeto usa [Task](https://taskfile.dev/) (`Taskfile.yml`) para comandos comuns.

```bash
task build     # Build the binary
task check     # static analysis (vet + lint + build) — run before every commit
task dev       # go run ./cmd/server
task seed      # apply migrations + seed data
task up        # docker compose up -d --build
task reset     # full reset (down -v + up + seed)
task logs      # docker compose logs -f app
task psql      # connect to DB container
```

## Banco de Dados e Migrations

PostgreSQL 15 via Docker Compose. As migrations ficam em `internal/database/migrations/`.

Atalho: `go run ./cmd/server --seed` aplica todas as migrations (001-011) + seed data.
Credenciais padrão: `tio_da_locadora` / `sopre_a_fita` / `modo_locadora`.

Consulte [docs/setup.md](docs/setup.md) para a lista completa de migrations e contas de teste.

## Arquitetura do Código

Consulte [ARCHITECTURE.md](ARCHITECTURE.md) para uma visão geral de alto nível do sistema.

### Estrutura de Pacotes

| Package | Propósito |
|---------|-----------|
| `cmd/server/main.go` | Entrypoint: config, template parsing, pgx pool, route wiring |
| `internal/handlers/` | HTTP handlers encapsulados em um `Handler` struct |
| `internal/database/` | Interface `Store` (`store.go`) e implementação `PostgreSQL` (`postgres.go`) |
| `internal/middleware/` | `RequireAuth` (session check) e `RequireAdmin` (auth + email check) |
| `internal/auth/` | HMAC-SHA256 cookie signing/verification |
| `internal/igdb/` | IGDB API client (Twitch OAuth2) |
| `internal/models/` | Domain structs: Member, Game, GameCopy, Rental, Club, etc. |
| `internal/storage/` | `StorageProvider` para file uploads (Local vs GCS) |
| `web/templates/` | Templates HTML standalone (UI em português) |
| `web/static/` | CSS global (`retro.css`), logos e assets enviados por upload |

### Authentication

- **Login**: `POST /login` -> bcrypt verify -> HMAC-signed cookie (`{uuid}.{hmac_hex}`).
- **Admin**: cookie verificado + email comparado com `ADMIN_EMAIL` env var.

### Fluxo de Navegação

1. `GET /games`: grid de seleção de plataforma (`platforms.html`).
2. `GET /games?platform=X`: prateleira de jogos filtrada (`games.html`).
3. `GET /games/{id}`: detalhe completo do jogo com estatísticas de aluguel (`game_detail.html`).

## Routing

### Public

- `GET /` — Landing page com login e "Wall of Shame"
- `POST /login` — Authentication
- `POST /members` — Registration (JSON API)
- `GET /search?q=` — IGDB search (JSON API)
- `GET /clubs` — Public club listing
- `GET /clubs/{id}` — Public club detail

### Authenticated (`RequireAuth`)

- `GET /membership` — Carteirinha digital e password notes
- `POST /rent` — Rent a game
- `POST /membership/return` — Self-return de aluguel (com verdict)
- `POST /membership/redeem` — Clear debt status ("Sopro")
- `POST /clubs/new`, `POST /clubs/{id}/join`, etc.

### Admin (`RequireAdmin`)

- `GET /admin/stock` — IGDB search e add games
- `GET /admin/inventory` — Full catalog listing
- `GET /admin/edit/{id}` — Edit game form e history
- `GET /admin/returns` — Active rentals dashboard

## Convenções e Regras

1. **Commit format**: Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`).
2. **Branching**: `main` (stable), `develop` (active).
3. **Routing**: somente standard library — `mux.HandleFunc("METHOD /path", handler)`.
4. **Templates**: arquivos HTML standalone. Use classes NES.css e utilities de `retro.css`.
5. **Migrations**: novas mudanças DEVEM ser adicionadas em um novo arquivo `.sql` numerado e registrado em `cmd/server/main.go`.
6. **Security**: parameterized SQL. Nunca armazene senhas em texto puro. Cookie secrets devem ter no mínimo 32 caracteres.
7. **SRE**: sempre rode `task check` antes de propor um PR.

## Tarefas Comuns

### Adicionar uma nova rota

1. Adicione um método na interface `Store` em `store.go` e implemente em `postgres.go`.
2. Adicione o handler em `internal/handlers/handler.go`.
3. Crie/modifique o template em `web/templates/`.
4. Registre a rota em `cmd/server/main.go` com o middleware apropriado.
5. Rode `task build` para verificar o template parsing.

### Adicionar uma database migration

1. Crie `internal/database/migrations/0XX_description.sql`.
2. Registre o arquivo na lista `sqlFiles` em `cmd/server/main.go`.
3. Documente em `docs/setup.md`.
