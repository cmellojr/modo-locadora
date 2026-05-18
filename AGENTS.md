# AGENTS.md

Este arquivo é a referência rápida para agentes de IA de programação (Gemini CLI, Jules, Claude Code, Copilot, Cursor, Codex, etc.) trabalhando no **Modo Locadora**.

Para evitar duplicação:

- Visão geral e início rápido: [README.md](README.md)
- Arquitetura e fluxo de domínio: [ARCHITECTURE.md](ARCHITECTURE.md)
- Setup, migrations, seed e storage: [docs/setup.md](docs/setup.md)
- Rotas e endpoints: [docs/api.md](docs/api.md)
- Segurança: [docs/security.md](docs/security.md)
- Contribuição: [docs/contributing.md](docs/contributing.md)

## Contexto Essencial

O Modo Locadora é uma aplicação Go SSR que simula uma locadora brasileira de videogames dos anos 1990. O domínio principal envolve sócios, jogos, cópias físicas, aluguéis, vereditos, reputação, turmas e uploads de capas/badges.

Stack atual:

- Go 1.24+
- PostgreSQL 15
- `net/http.ServeMux`
- `html/template`
- NES.css
- Docker Compose
- IGDB via Twitch OAuth2
- `internal/storage` com `LocalStorage` e `GCSStorage`

## Regras de Idioma

- Código, rotas, colunas de banco, env vars, logs, query params, tipos, funções e comentários: inglês.
- UI em `web/templates/`: português (BR).
- Documentação `.md`: português (BR), mantendo comandos, paths e identificadores técnicos em inglês.

## Convenções de Implementação

- Use a standard library para routing: `mux.HandleFunc("METHOD /path", handler)`.
- Não adicione router externo.
- Mantenha o frontend sem JavaScript.
- Siga os padrões existentes de `Handler`, `Store` e `PostgresStore`.
- Uploads devem passar por `storage.StorageProvider`; não grave diretamente em `web/static/covers` ou `web/static/clubs` em handlers novos.
- Mudanças de schema exigem nova migration em `internal/database/migrations/` e registro em `sqlFiles` em `cmd/server/main.go`.
- Texto visível ao usuário deve ficar em português (BR).

## Comandos Úteis

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

Use `task check` antes de propor PR. Se Task não estiver disponível:

```bash
go build ./...
go vet ./...
golangci-lint run ./...
```

## Fluxo de Branches

- `main`: estável.
- `develop`: desenvolvimento ativo e alvo de PRs.
- Branches sugeridas: `feature/*`, `fix/*`, `hotfix/*`, `docs/*`, `chore/*`.

Commits seguem Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:`, `chore:`).

## Antes de Editar

1. Confira `git status --short --branch`.
2. Não reverta mudanças locais que você não fez.
3. Leia o arquivo responsável pelo assunto antes de alterar documentação.
4. Se tocar comportamento do software, atualize docs relevantes e [docs/changelog.md](docs/changelog.md).
