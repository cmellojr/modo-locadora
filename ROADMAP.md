# Roadmap

> Plano de evolução do Modo Locadora.

O projeto segue [Versionamento Semântico](https://semver.org/). Enquanto estiver em `0.1.x`, o foco é desenvolvimento ativo e experimentação.

## Lançadas

### 0.1.0 — Estrutura Inicial

- Projeto Go + PostgreSQL.
- Cliente IGDB com OAuth2 Twitch.
- Páginas SSR iniciais.
- Migration `001_initial_schema.sql`.

### 0.1.1 — Docker e Acervo

- Docker Compose com app + PostgreSQL.
- Fluxo admin de busca IGDB e aquisição de jogos.
- Tema NES.css dark.
- Migration `002_update_games_table.sql`.

### 0.1.2 — Sócios e Aluguéis

- Matrícula `1991-XXX`, autenticação bcrypt e cookies HMAC.
- Aluguel/devolução com dashboard admin.
- Middleware `RequireAuth` e `RequireAdmin`.
- Migration `003_membership_and_rental_support.sql`.

### 0.1.3 — Feed, Reputação e Passwords

- Caderno de passwords.
- Auto-devolução com penalização.
- Feed "Aconteceu na Locadora".
- Almanaque do Tio.
- Seed inicial.
- Migrations `004` a `007`.

### 0.1.4 — Acervo e Progressão

- Títulos de progressão do sócio.
- Campo `cover_display`.
- Histórico de aluguéis no admin.
- Taskfile e `golangci-lint`.
- Migration `008_cover_display.sql`.

### 0.1.5 — Turmas

- Criação, edição, exclusão, entrada e saída de turmas.
- Badges, URL, múltiplos admins e integração com carteirinha.
- Migration `009_clubs.sql`.

### 0.1.6 — Internacionalização Interna e Popularidade

- Rotas, query params, status internos e vereditos normalizados em inglês.
- Popularidade do acervo substitui a antiga saúde do cartucho.
- Cinco vereditos de devolução.
- Migrations `010_rename_status_english.sql` e `011_verdict_popularity.sql`.

### 0.1.7 — Storage e Documentação

- Abstração `StorageProvider` para uploads.
- `LocalStorage` para desenvolvimento e `GCSStorage` para produção.
- Consolidação de instruções de agentes em `AGENTS.md`.
- Documentação em português com referências cruzadas.

## Próximas

### 0.1.8+ — Interação Social

- **Verso da Capa**: dicas públicas para próximos jogadores ao devolver uma fita.
- **Regra da Sexta**: aluguel na sexta com devolução na segunda.
- **Roleta do Tio**: sugestão aleatória para quem não sabe o que alugar.
- **Menções na Mídia**: podcasts, vídeos e reportagens associados a jogos.

### 0.2.0 — Milestone Estável

- Congelamento do primeiro conjunto estável de regras.
- Revisão de UX e documentação de deploy.
- Roadmap detalhado até `1.0`.

Sugestões podem virar issue ou PR. Veja [docs/contributing.md](docs/contributing.md).
