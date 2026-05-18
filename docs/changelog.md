# Histórico de Mudanças

Todas as mudanças notáveis deste projeto são documentadas aqui.

O formato segue [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) e o projeto adere ao [Versionamento Semântico](https://semver.org/spec/v2.0.0.html).

## [Não Lançado]

### Adicionado

- Abstração `StorageProvider` para uploads, com `LocalStorage` em desenvolvimento e `GCSStorage` em produção.
- Suporte a Google Cloud Storage para capas e badges quando `APP_ENV=production`.
- Documentação em português consolidada e com referências cruzadas entre README, arquitetura, setup, API, segurança, PRD, roadmap e contribuição.

### Alterado

- `AGENTS.md` agora é a fonte única de instruções para agentes de IA.
- Documentação removida/reestruturada para reduzir redundância entre arquivos.

### Removido

- Arquivos de instrução separados para agentes foram consolidados em `AGENTS.md`.
- Arquivo de descrição isolada do repositório foi removido.

## [0.1.6] - 2026-05-18

### Adicionado

- Cinco vereditos de devolução: `completed`, `enjoyed`, `quick_play`, `not_for_me`, `gave_up`.
- Veredito `auto_return` para devoluções automáticas.
- Classificação de popularidade no inventário admin.
- Migration `011_verdict_popularity.sql`.

### Alterado

- Vereditos antigos em português foram migrados para slugs em inglês.
- Eventos de atividade de veredito foram normalizados.
- Indicador de "saúde do acervo" foi substituído por popularidade.

## [0.1.5] - 2026-05-18

### Adicionado

- Sistema de turmas com criação, edição, exclusão, entrada, saída, promoção e remoção de membros.
- Badges de turma e integração com carteirinha.
- Migration `009_clubs.sql`.

### Alterado

- Navegação global passou a incluir Turmas.

## [0.1.4] - 2026-05-18

### Adicionado

- Títulos de progressão do sócio.
- `cover_display` para controlar exibição de capas.
- Histórico de aluguéis no admin.
- Taskfile e `golangci-lint`.
- Migration `008_cover_display.sql`.

## [0.1.3] - 2026-05-18

### Adicionado

- Caderno de passwords.
- Auto-devolução com penalização.
- Feed de atividades.
- Almanaque do Tio.
- Seed inicial.
- Migrations `004` a `007`.

## [0.1.2] - 2026-03-04

### Adicionado

- Matrículas `1991-XXX`.
- Carteirinha de sócio.
- Sistema de aluguel e devolução.
- Dashboard admin de devoluções.
- Inventário e edição admin.
- Middleware `RequireAuth` e `RequireAdmin`.
- Cookies assinados com HMAC-SHA256 e senhas com bcrypt.
- Migration `003_membership_and_rental_support.sql`.

## [0.1.1] - 2026-03-03

### Adicionado

- Docker Compose com PostgreSQL 15.
- Estoque admin com busca IGDB.
- Fluxo de aquisição de jogos.
- Registro de sócios.
- Tema NES.css.
- Migration `002_update_games_table.sql`.

## [0.1.0] - 2026-03-03

### Adicionado

- Estrutura inicial em Go.
- PostgreSQL com `pgx/v5`.
- Interface `Store`.
- Cliente IGDB com Twitch OAuth2.
- Primeiras páginas SSR.
- Migration `001_initial_schema.sql`.
- Licença GPL v3.
