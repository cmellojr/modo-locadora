# Contribuindo com o Modo Locadora

Obrigado por contribuir. Issues, sugestões e pull requests são bem-vindos.

## Antes de Começar

- Leia [../README.md](../README.md) para entender o projeto.
- Use [setup.md](setup.md) para configurar o ambiente.
- Consulte [../AGENTS.md](../AGENTS.md) se estiver usando um agente de IA.
- Segurança vai em [security.md](security.md).

## Fluxo de Pull Request

1. Crie uma branch a partir de `develop`.
2. Faça mudanças pequenas e focadas.
3. Atualize documentação quando alterar comportamento, rotas, env vars, migrations ou fluxo de usuário.
4. Rode `task check`.
5. Use Conventional Commits.
6. Abra PR contra `develop`.

## Convenções

### Idiomas

- Código, rotas, colunas de banco, env vars, logs e query params: inglês.
- UI em `web/templates/`: português (BR).
- Documentação `.md`: português (BR), mantendo comandos, paths, nomes de tipos, funções e variáveis em inglês.

### Commits

Use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:`
- `fix:`
- `docs:`
- `refactor:`
- `style:`
- `test:`
- `chore:`

### Branches

- `main`: estável.
- `develop`: desenvolvimento ativo e alvo de PRs.
- Prefixos sugeridos: `feature/*`, `fix/*`, `hotfix/*`, `docs/*`, `chore/*`.

### Go

- Siga o estilo idiomático de Go.
- Prefira a standard library quando o projeto já usa standard library.
- Não adicione router externo; use `mux.HandleFunc("METHOD /path", handler)`.
- Atualize `Store` e `PostgresStore` juntos quando a mudança tocar persistência.

### Templates e CSS

- Sem JavaScript.
- Templates ficam em `web/templates/` e usam `layout.html`.
- Estilos compartilhados ficam em `web/static/css/retro.css`.
- Estilos muito específicos podem ficar no `<style>` do template.
- Texto visível ao usuário deve ser português (BR).

### Migrations

- Crie arquivo novo em `internal/database/migrations/`.
- Use numeração incremental.
- Registre o arquivo em `sqlFiles` em `cmd/server/main.go`.
- Atualize [setup.md](setup.md) e, se relevante, [changelog.md](changelog.md).

### Storage

- Uploads devem passar por `internal/storage.StorageProvider`.
- Não escreva diretamente em `web/static/covers` ou `web/static/clubs` dentro de handlers novos.

## Validação

```bash
task check
```

Sem Task:

```bash
go build ./...
go vet ./...
golangci-lint run ./...
```

## Licença

Ao contribuir, você concorda que suas contribuições serão licenciadas sob a [GPL v3](../LICENSE).
