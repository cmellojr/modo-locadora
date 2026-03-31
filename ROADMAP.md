# Roadmap

> Plano de evolução do Modo Locadora.

O projeto segue [Versionamento Semântico](https://semver.org/). Enquanto estiver em `0.1.x`, o foco é desenvolvimento ativo e experimentação. A versão `0.2.0` marcará o primeiro milestone estável com roadmap definido até o `1.0`.

---

## Lançadas

### 0.1.0 — Estrutura Inicial

- Projeto Go + PostgreSQL + Docker Compose
- Cliente IGDB com OAuth2 Twitch
- Página de entrada e prateleira de jogos com SSR
- Migration `001_initial_schema.sql`

### 0.1.1 — Docker e Acervo

- Docker Compose com app + PostgreSQL
- Página de estoque admin com busca IGDB
- Tema NES.css dark com grid responsivo
- Migration `002_update_games_table.sql`

### 0.1.2 — Sócios e Aluguéis

- Sistema de carteirinha (`1991-XXX`)
- Aluguel e devolução com dashboard admin
- Autenticação bcrypt + cookies HMAC-SHA256
- Middleware `RequireAuth` e `RequireAdmin`
- Migration `003_membership_and_rental_support.sql`

### 0.1.3 — Feed e Vereditos

- Vereditos de devolução (Zerei / Joguei um pouco / Desisti)
- Feed "Aconteceu na Locadora" em tempo real
- Estrela Dourada para jogos zerados
- Auto-return com penalização e Painel da Vergonha
- Caderno de Passwords
- Almanaque do Tio (efemérides gaming)
- Seed com dados da Ação Games #1
- Migrations `004` a `007`

### 0.1.4 — Acervo e Progressão

- Títulos de progressão (Sócio Novato → Prata → Ouro → Dono da Calçada)
- Indicadores de saúde do acervo (Cartucho Novo / Clássico Eterno / Precisa Soprar / Fita Gasta)
- Histórico de aluguéis na ficha do jogo
- Campo `cover_display` para modo de exibição de capas
- Capas brasileiras TecToy no seed
- Taskfile para automação SRE
- golangci-lint configurado
- Migration `008_cover_display.sql`

### 0.1.5 — Turmas (atual)

- Turmas (comunidades gamers): podcasts, canais YouTube, grupos WhatsApp
- Criação, edição e exclusão de turmas com badge personalizado
- Relação M2M: sócios podem participar de múltiplas turmas
- Múltiplos admins por turma (promover/remover membros)
- Listagem pública, ações requerem autenticação
- Seção "MINHAS TURMAS" na carteirinha
- Feed de atividades registra criação e entrada em turmas
- Botão TURMAS em todas as páginas
- Migration `009_clubs.sql`

---

## Próximas

### 0.1.6+ — Interação Social

- **Verso da Capa** — Dicas públicas para os próximos jogadores ao devolver uma fita.
- **Regra da Sexta** — Alugou na sexta? Só precisa devolver na segunda!
- **Roleta do Tio** — Modo aleatório para quem não sabe o que alugar.
- **Menções na Mídia** — Registre em quais podcasts ou reportagens cada jogo foi mencionado.

### 0.2.0 — (a definir)

Primeiro milestone estável. Roadmap detalhado até o `1.0` será publicado nesta versão.

---

*Sugestões? Abra uma [issue](https://github.com/cmellojr/modo-locadora/issues) ou mande um PR.*
