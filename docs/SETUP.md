# Configuração do Ambiente

Guia passo a passo para configurar o Modo Locadora localmente.

## Pré-requisitos

| Ferramenta | Versão | Função |
|------------|--------|--------|
| Docker     | 20+    | Containers da app + PostgreSQL |
| Git        | 2.x    | Controle de versão |

Para desenvolvimento local sem Docker, você também precisa de **Go 1.24+** e uma instância **PostgreSQL 15+**.

## 1. Clone e Configuração

```bash
git clone https://github.com/cmellojr/modo-locadora.git
cd modo-locadora
cp .env.example .env
```

Edite o `.env` com seus valores:

```env
# IGDB API — Obtenha credenciais em https://dev.twitch.tv/console
TWITCH_CLIENT_ID=your_client_id
TWITCH_CLIENT_SECRET=your_client_secret

# Banco de dados (usado pelo Docker Compose)
DB_USER=tio_da_locadora
DB_PASSWORD=sopre_a_fita
DB_NAME=modo_locadora

# Segurança
COOKIE_SECRET=generate-a-random-secret-here-min-32-chars
ADMIN_EMAIL=your_admin_email@example.com
```

### Obtendo Credenciais da IGDB

1. Crie uma conta no [Twitch Developer Console](https://dev.twitch.tv/console).
2. Registre uma nova aplicação (qualquer categoria).
3. Copie o **Client ID** e gere um **Client Secret**.

## 2. Iniciar com Docker (recomendado)

```bash
docker compose up -d --build
```

Isso inicia a aplicação Go e o PostgreSQL. A app conecta ao banco automaticamente. Acesse em `http://localhost:8080`.

As migrations são aplicadas manualmente (veja passo 3).

## 3. Executar Migrations

As migrations ficam em `internal/database/migrations/` e devem ser aplicadas em ordem.

### Via container Docker:

```bash
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/001_initial_schema.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/002_update_games_table.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/003_membership_and_rental_support.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/004_password_notes.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/005_auto_return_reputation.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/006_activities_feed.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/008_cover_display.sql
docker exec -i modo_locadora_db psql -U tio_da_locadora -d modo_locadora < internal/database/migrations/009_clubs.sql
```

### Via psql direto:

```bash
psql $DATABASE_URL -f internal/database/migrations/001_initial_schema.sql
psql $DATABASE_URL -f internal/database/migrations/002_update_games_table.sql
psql $DATABASE_URL -f internal/database/migrations/003_membership_and_rental_support.sql
psql $DATABASE_URL -f internal/database/migrations/004_password_notes.sql
psql $DATABASE_URL -f internal/database/migrations/005_auto_return_reputation.sql
psql $DATABASE_URL -f internal/database/migrations/006_activities_feed.sql
psql $DATABASE_URL -f internal/database/migrations/008_cover_display.sql
psql $DATABASE_URL -f internal/database/migrations/009_clubs.sql
```

### Setup rápido com dados de teste:

```bash
# Desenvolvimento local:
go run ./cmd/server --seed

# Dentro do Docker (após docker compose up):
docker exec modo_locadora_app /app/server --seed
```

Isso aplica todas as migrations (001-009) e popula o banco com jogos, sócios, turmas e histórico de aluguéis. A flag `--seed` auto-detecta o diretório de migrations (`migrations/` no Docker, `internal/database/migrations/` localmente).

### Contas de teste

| Sócio | Senha | Perfil |
|-------|-------|--------|
| `MegaDriveKid` | `sega1991` | Sócio com histórico de aluguéis |
| `Devedor` | `atrasado123` | Sócio em débito |
| `Novato` | `novato2026` | Sócio sem histórico |

Admin: `tio_da_locadora` / `sopre_a_fita` (e-mail deve bater com `ADMIN_EMAIL`).

### Resumo das Migrations

| Migration | Descrição |
|-----------|-----------|
| `001_initial_schema.sql` | Tabelas base: `members`, `games`, `game_copies`, `rentals` |
| `002_update_games_table.sql` | Adiciona `cover_url`, `source_magazine`, `acquired_at` a `games` |
| `003_membership_and_rental_support.sql` | Campos de matrícula, sequência `membership_seq`, auto-criação de cópias |
| `004_password_notes.sql` | Campo `password_notes` em `members` |
| `005_auto_return_reputation.sql` | Campos `status` e `late_count` em `members` |
| `006_activities_feed.sql` | Tabela `activities` para feed de eventos |
| `007_seed_initial_data.sql` | Dados de teste (aplicado via flag `--seed`, não manualmente) |
| `008_cover_display.sql` | Campo `cover_display` em `games` (modo CSS object-fit) |
| `009_clubs.sql` | Tabelas `clubs` e `club_members` (turmas/comunidades gamers) + dados de seed |

## 4. Desenvolvimento Local (sem Docker para a app)

Se preferir rodar o servidor Go localmente mantendo o PostgreSQL no Docker:

```bash
docker compose up -d db       # inicia apenas o PostgreSQL
go run ./cmd/server            # roda o servidor Go localmente
```

Configure `DATABASE_URL` no `.env` apontando para `localhost:5432` (não `db:5432`).

```bash
# Build do binário
go build -o modo-locadora ./cmd/server

# Análise estática
go vet ./...

# Lint (requer golangci-lint)
golangci-lint run ./...

# Ou use o Task runner para todas as verificações
task check
```

## 5. Criando o Primeiro Sócio

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

O e-mail deve bater com `ADMIN_EMAIL` para acesso de administrador. Um número de matrícula (`1991-001`) é auto-atribuído.

## 6. Verificação

| Verificação | Esperado |
|-------------|----------|
| `http://localhost:8080` | Página de login carrega |
| Login com credenciais | Redireciona para `/games` (grade de plataformas) |
| Clicar numa plataforma | Mostra cards de cartucho |
| Clicar num cartucho | Mostra página de detalhe do jogo |
| `/carteirinha` (logado) | Carteirinha com `1991-XXX` + MINHAS TURMAS |
| `/clubs` | Listagem de turmas (com seed: "Turma da Acao Games") |
| `/admin/stock` (como admin) | Página de busca IGDB |

## Resolução de Problemas

### "No DATABASE_URL provided"
Verifique se `DATABASE_URL` está definido no `.env`. Ao usar Docker Compose para o stack completo, é definido automaticamente via `docker-compose.yml`.

### "COOKIE_SECRET not set"
O servidor usa um valor padrão inseguro em desenvolvimento. Para produção, defina um valor aleatório forte (pelo menos 32 caracteres).

### "ADMIN_EMAIL not set"
As rotas admin (`/admin/*`) ficam inacessíveis sem isso. Defina com o e-mail do sócio administrador.

### Busca IGDB não retorna resultados
Verifique as credenciais Twitch:
```bash
curl -X POST "https://id.twitch.tv/oauth2/token?client_id=YOUR_ID&client_secret=YOUR_SECRET&grant_type=client_credentials"
```

### Porta 8080 já em uso
```bash
# Linux/Mac
lsof -ti:8080 | xargs kill -9

# Windows
netstat -ano | findstr :8080
taskkill /F /PID <pid>
```
