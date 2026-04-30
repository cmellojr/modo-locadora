# Contribuindo com o Modo Locadora

Obrigado por contribuir! Toda contribuição importa, seja um relato de bug, sugestão de funcionalidade ou pull request.

## Reportando Bugs

Verifique as [Issues](https://github.com/cmellojr/modo-locadora/issues) existentes primeiro.

Inclua: título claro, passos para reproduzir, comportamento esperado vs. real, detalhes do ambiente (OS, versão do Go, navegador) e logs ou screenshots se aplicável.

## Sugerindo Funcionalidades

Abra uma issue com o label `enhancement`. Descreva o problema que resolve, como você imagina funcionando e alternativas consideradas.

## Pull Requests

1. Faça fork do repositório e clone seu fork.
2. Crie uma branch a partir de `develop`: `git checkout -b feature/nome-da-feature`
3. Siga as convenções abaixo.
4. Verifique com `task check` (ou `go build ./...`, `go vet ./...`, `golangci-lint run ./...`).
5. Faça commit com mensagem descritiva (veja convenções abaixo).
6. Faça push para seu fork e abra um PR contra `develop`.
7. Na descrição do PR, conte qual jogo você estava jogando enquanto programava (é tradição).

## Convenções

### Idiomas

- **Código, rotas, colunas de banco e variáveis**: inglês.
- **Documentação (.md)**: português (BR), exceto CLAUDE.md e AGENTS.md (guias para agentes de IA, em inglês).
- **Templates de UI** (`web/templates/`): português (BR).

### Estilo Go

- Siga o [Google Go Style Guide](https://google.github.io/styleguide/go/guide.html).
- Execute `task check` antes de commitar.
- Mantenha funções curtas e focadas.

### Mensagens de Commit

Use [Conventional Commits](https://www.conventionalcommits.org/): `feat:`, `fix:`, `docs:`, `style:`, `refactor:`.

### Branches

- `main` — releases estáveis
- `develop` — desenvolvimento ativo (alvo de PRs)
- Branches de feature: `feature/*`, `fix/*`, `hotfix/*`, `docs/*`

### CSS e Templates

- Classes NES.css com overrides de tema escuro em `retro.css`.
- Classes utilitárias compartilhadas vão no `retro.css`. Estilos específicos da página vão no `<style>` inline do template.

### Migrations de Banco de Dados

- Coloque novas migrations em `internal/database/migrations/`.
- Use numeração incremental: `009_description.sql`, `010_description.sql`.
- Adicione o arquivo à lista `sqlFiles` em `cmd/server/main.go` (para a flag `--seed`).
- Documente o que cada migration faz no cabeçalho do arquivo.

## Configuração do Ambiente

Veja [Configuração do Ambiente](setup.md) para instruções de instalação.

## Licença

Ao contribuir, você concorda que suas contribuições serão licenciadas sob a [GPL v3](../LICENSE).
