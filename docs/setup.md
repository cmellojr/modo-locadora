# Configuração do Ambiente

Este guia cobre setup local, Docker, migrations, seed e variáveis de ambiente. Para a visão de arquitetura, veja [../ARCHITECTURE.md](../ARCHITECTURE.md).

## Pré-requisitos

| Ferramenta | Uso |
|------------|-----|
| Docker + Docker Compose | App + PostgreSQL local |
| Git | Controle de versão |
| Go 1.24+ | Desenvolvimento sem container da app |
| PostgreSQL 15+ | Necessário se não usar o banco do Compose |
| Task | Atalhos de build, check, seed e Docker |
| golangci-lint | Usado por `task lint` e `task check` |

## Variáveis de Ambiente

Copie `.env.example` para `.env` e ajuste os valores.

| Variável | Obrigatória | Uso |
|----------|-------------|-----|
| `TWITCH_CLIENT_ID` | Sim | Client ID da Twitch para IGDB |
| `TWITCH_CLIENT_SECRET` | Sim | Client Secret da Twitch para IGDB |
| `DB_HOST` | Local | Host do PostgreSQL fora do Compose |
| `DB_PORT` | Local | Porta do PostgreSQL |
| `DB_USER` | Sim | Usuário do banco |
| `DB_PASSWORD` | Sim | Senha do banco |
| `DB_NAME` | Sim | Nome do banco |
| `DATABASE_URL` | Sim fora do Compose | String PostgreSQL usada pela app |
| `COOKIE_SECRET` | Produção | Chave HMAC; use 32+ caracteres |
| `ADMIN_EMAIL` | Sim para admin | E-mail que libera `/admin/*` |
| `PORT` | Não | Porta HTTP; padrão `8080` |
| `APP_ENV` | Produção GCS | Use `production` para ativar `GCSStorage` |
| `STORAGE_BUCKET_NAME` | Produção GCS | Bucket do Google Cloud Storage |

O Docker Compose local define `DATABASE_URL` para a app automaticamente. `APP_ENV` e `STORAGE_BUCKET_NAME` não são necessários para desenvolvimento local, porque uploads usam `LocalStorage`.

## Início com Docker

```bash
git clone https://github.com/cmellojr/modo-locadora.git
cd modo-locadora
cp .env.example .env
docker compose up -d --build
docker exec modo_locadora_app /app/server --seed
```

Acesse `http://localhost:8080`.

## Desenvolvimento Local

Para rodar a app fora do Docker mantendo o PostgreSQL no Compose:

```bash
docker compose up -d db
go run ./cmd/server --seed
go run ./cmd/server
```

Garanta que `DATABASE_URL` aponte para `localhost:5432`.

## Task Runner

```bash
task build
task vet
task lint
task check
task dev
task seed
task up
task down
task reset
task logs
task psql
```

Use `task check` antes de PRs. O fluxo de contribuição está em [contributing.md](contributing.md).

## Migrations

A forma recomendada para ambiente novo é:

```bash
go run ./cmd/server --seed
```

ou, no container:

```bash
docker exec modo_locadora_app /app/server --seed
```

A flag `--seed` aplica todas as migrations e depois `007_seed_initial_data.sql`. O diretório é detectado automaticamente:

- `migrations/` dentro do container.
- `internal/database/migrations/` localmente.

### Ordem Atual

| Migration | Descrição |
|-----------|-----------|
| `001_initial_schema.sql` | Tabelas base: `members`, `games`, `game_copies`, `rentals` |
| `002_update_games_table.sql` | Metadados de capa, revista e aquisição |
| `003_membership_and_rental_support.sql` | Matrícula, `membership_seq` e suporte a aluguéis |
| `004_password_notes.sql` | Caderno de passwords |
| `005_auto_return_reputation.sql` | `status` e `late_count` |
| `006_activities_feed.sql` | Feed `activities` |
| `007_seed_initial_data.sql` | Seed de jogos, sócios, histórico e feed |
| `008_cover_display.sql` | Campo `cover_display` |
| `009_clubs.sql` | `clubs` e `club_members` |
| `010_rename_status_english.sql` | `em_debito` -> `in_debt` |
| `011_verdict_popularity.sql` | Vereditos em inglês e eventos de atividade atualizados |

## Contas de Seed

| Perfil | Senha | Observação |
|--------|-------|------------|
| `MegaDriveKid` | `sega1991` | Sócio com histórico |
| `Devedor` | `atrasado123` | Sócio em débito |
| `Novato` | `novato2026` | Sócio novo |
| `tio_da_locadora` | `sopre_a_fita` | Admin se o e-mail bater com `ADMIN_EMAIL` |

## Primeiro Admin Manual

```bash
curl -X POST http://localhost:8080/members \
  -H "Content-Type: application/json" \
  -d '{
    "profile_name": "Tio da Locadora",
    "email": "admin@locadora.com",
    "password": "sopre_a_fita",
    "favorite_console": "Mega Drive"
  }'
```

O e-mail precisa ser igual a `ADMIN_EMAIL`.

## Uploads

Uploads passam por `internal/storage`:

- Local: `web/static/covers/` e `web/static/clubs/`.
- Docker local: volumes `covers_data` e `clubs_data`.
- Produção: `APP_ENV=production` ativa GCS; `STORAGE_BUCKET_NAME` deve estar definido.

Detalhes de segurança de upload estão em [security.md](security.md).

## Troubleshooting

| Sintoma | Verificação |
|---------|-------------|
| `No DATABASE_URL provided` | Defina `DATABASE_URL` ou use Docker Compose completo |
| `COOKIE_SECRET not set` | Em dev há fallback inseguro; em produção defina 32+ caracteres |
| Admin inacessível | Confira `ADMIN_EMAIL` e o e-mail do sócio logado |
| IGDB sem resultados | Valide `TWITCH_CLIENT_ID` e `TWITCH_CLIENT_SECRET` |
| Porta ocupada | Ajuste `PORT` ou libere a porta local |
