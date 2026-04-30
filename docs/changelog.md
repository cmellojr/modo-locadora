# Histórico de Mudanças

Todas as mudanças notáveis deste projeto serão documentadas neste arquivo.

O formato segue o [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) e o projeto adere ao [Versionamento Semântico](https://semver.org/spec/v2.0.0.html).

## [Não Lançado]

### Adicionado
- **Banner de imagem 728x90**: Título do site substituído por imagem PNG no formato leaderboard clássico dos anos 2000. Renderização pixel art via `image-rendering: pixelated`, escala responsiva automática.
- **Layout global 3 colunas (anos 2000)**: Estrutura de site inspirada em GameFAQs/Backloggery — sidebar esquerda (navegação + mini-card), área de conteúdo central, sidebar direita (feed + vergonha + almanaque). Template base `layout.html` com composição via `{{define "content"}}`. Todos os 12 templates convertidos.
- **Header com banner + barra de navegação**: Header dividido em duas linhas — banner com gradiente e logo no topo, barra de links tabulados abaixo (Balcão, Prateleira, Turmas, Carteirinha, Admin).
- **Turmas (Clubs)**: Sistema de comunidades gamers — sócios podem criar turmas (representando podcasts, canais YouTube, grupos WhatsApp etc.), personalizar com badge e URL, e convidar outros sócios. Múltiplos admins por turma; listagem pública. Relação M2M (`clubs` ↔ `club_members` ↔ `members`). Carteirinha exibe "MINHAS TURMAS" com badge, nome e cargo. Feed de atividades registra criação e entrada em turmas. Upload de badges salvos em `web/static/clubs/` (volume Docker). Migration `009_clubs.sql`. 11 rotas novas, 3 templates (`clubs.html`, `club_detail.html`, `club_form.html`). Botão TURMAS adicionado a todas as páginas.
- **Títulos de progressão** (`Status de Veterano`): Carteirinha exibe badge de progressão — Sócio Novato (padrão), Sócio Prata (10+ devoluções no prazo), Sócio Ouro (25+ devoluções no prazo), Dono da Calçada (5+ jogos zerados). Devedores veem título esmaecido com "(EM DÉBITO)". Computação pura via `ComputeMemberTitle()` em `internal/models/member.go`.
- **Indicadores de saúde do acervo** (`Saúde do Acervo`): Inventário admin mostra coluna de saúde por jogo — Cartucho Novo (0-1 aluguéis), Clássico Eterno (<25% vereditos ruins), Precisa Soprar (25-49%), Fita Gasta (50%+). Calculado a partir de vereditos e atrasos via `ListGamesWithHealth()`.
- **Histórico de aluguéis na ficha do jogo**: Página de edição admin mostra os últimos 5 registros de aluguel (sócio, datas, veredito, indicador de atraso) via `ListGameRentalHistory()`.
- **Modo de exibição de capa**: Jogos agora têm campo `cover_display` (cover/contain/fill) controlando CSS `object-fit` das imagens de capa. Editável na página admin. Migration `008_cover_display.sql`.
- **Configuração golangci-lint** (`.golangci.yml`): Linters — errcheck, staticcheck, unused, gosec, govet, ineffassign, typecheck.
- **Taskfile** (`Taskfile.yml`): Task runner SRE com comandos para build, vet, lint, check, dev, seed, up, down, reset, logs, psql.
- **Auto-detecção de diretório de migrations**: Flag `--seed` agora detecta diretório automaticamente (`migrations/` no Docker, `internal/database/migrations/` localmente).
- **Feed de atividades** (`Aconteceu na Locadora`): Feed de eventos em tempo real na página de plataformas mostrando novos jogos, vereditos, penalidades e redenções. Migration `006_activities_feed.sql`.
- **Almanaque do Tio** (`internal/almanac/`): Efemérides estáticas de gaming por dia do ano, exibidas ao lado do feed.
- **Sistema de vereditos**: Na devolução, sócios escolhem um veredito (Zerei / Joguei um pouco / Desisti). Vereditos armazenados em `public_legacy` e geram eventos de atividade.
- **Estrela Dourada**: Jogos completados ("Zerei") pelo sócio logado mostram estrela dourada na prateleira.
- **Auto-devolução na carteirinha**: Sócios podem devolver aluguéis diretamente da carteirinha via `POST /membership/return`.
- **Layout 3 colunas**: Página de plataformas usa CSS Grid (`.locadora-grid`): sidebar esquerda (mini-card + Painel da Vergonha), centro (grade de plataformas), sidebar direita (feed + almanaque).
- **Sistema de seed SQL**: Flag `--seed` aplica todas as migrations + `007_seed_initial_data.sql` com 5 jogos da Ação Games #1, 3 sócios de teste, histórico de aluguéis e feed. Execução via `go run ./cmd/server --seed`.
- **Fluxo de login/logout**: Balcão é sempre a página de entrada; sócios logados veem mensagem de boas-vindas + navegação. `POST /logout` limpa o cookie de sessão.
- **Barra de autenticação**: Todas as páginas exibem "Sócio: nome / [DESCONECTAR]" alinhado no topo direito quando logado. Classe CSS `.auth-bar` em `retro.css`.
- **Logos de console na grade**: Página de seleção de plataforma mostra logos SVG (`web/static/img/logos/`) ao invés de capas de jogos. Mapeamento automático via helper `platformLogoFile()`.
- **Navegação em 3 níveis**: `/games` mostra grade de plataformas, `?platform=X` filtra por console, `/games/{id}` mostra detalhe completo com stats.
- **Upload de capas brasileiras**: Admin pode enviar imagens locais (TecToy, Playtronic) via formulário multipart na página de edição. Capas salvas em `web/static/covers/` (volume Docker).
- **Sistema de auto-devolução**: Job de background verifica aluguéis atrasados a cada 5 minutos, auto-devolve e penaliza sócios (status `in_debt` + incremento de `late_count`). Migration `005_auto_return_reputation.sql`.
- **Painel da Vergonha**: Página de entrada mostra maiores devedores.
- **Redenção de sócio**: `POST /membership/redeem` limpa status de débito.
- **Aplicação Dockerizada**: Dockerfile multi-stage, Docker Compose roda app + PostgreSQL, volume `covers_data` para uploads.
- **Caderno de Passwords**: Sócios podem salvar códigos de jogos na carteirinha. Migration `004_password_notes.sql`.
- **Pacote `internal/jobs/`**: Goroutine de background para processamento de aluguéis atrasados.
- **CLAUDE.md** e **AGENTS.md**: Arquivos de orientação para agentes de IA.

### Alterado
- **Convenção de idioma reforçada**: Rotas `/carteirinha` renomeadas para `/membership`; status `em_debito` renomeado para `in_debt` (migration `010_rename_status_english.sql`). Todas as mensagens `http.Error`, logs e query params (`?success=`) traduzidos para inglês. Português restrito exclusivamente ao texto da interface web.
- **Breakpoint responsivo**: Reduzido de 1100px para 768px — sidebars permanecem visíveis em telas médias.
- **Página de entrada**: Removido redirecionamento automático; Balcão sempre exibido primeiro com conteúdo condicional login/boas-vindas.
- **Prateleira simplificada**: Cards agora mostram apenas capa, título, número de cópias e disponibilidade (sem resumo/revista).
- **`GET /games/{id}`** mudou de API JSON para página renderizada no servidor.
- **`POST /rent`** redireciona para página de detalhe ao invés da prateleira.
- **Página de plataformas**: Reestruturada de flex 2 colunas para CSS Grid 3 colunas.
- **Expansão de componentes NES.css**: `nes-radio` para veredito, `nes-icon star` para estrela dourada, `nes-progress`, `nes-list`, `nes-avatar`, `nes-dialog`, `nes-balloon`.
- **Tamanhos de fonte e larguras** padronizados em todas as páginas.
- **Estilos inline consolidados** em classes CSS reutilizáveis no `retro.css`.

## [0.1.2] - 2026-03-04

### Adicionado
- **Sistema de matrícula**: Números sequenciais no formato `1991-XXX`.
- **Página da carteirinha**: Carteirinha digital de sócio.
- **Sistema de aluguel**: Sócios alugam jogos via `POST /rent`.
- **Dashboard de devoluções**: Página admin listando aluguéis ativos com botões de devolução.
- **Inventário e edição admin**: Tabela do acervo com links de edição e formulário de edição.
- **Migration `003`**: Campos de matrícula, sequência `membership_seq`, criação automática de cópias.
- Middleware `RequireAuth` e `RequireAdmin`.
- Cookies assinados com HMAC-SHA256 e hash de senhas com bcrypt.
- `.env.example` e diretório `docs/`.

### Alterado
- Prateleira de jogos exibe estados de disponibilidade em tempo real.
- Login requer nome de perfil + senha com verificação bcrypt.
- Cookie de sessão armazena UUID assinado.

## [0.1.1] - 2026-03-03

### Adicionado
- Docker Compose para PostgreSQL 15.
- Página de estoque admin com busca IGDB.
- Fluxo de aquisição de jogos, endpoints JSON de busca e detalhe.
- Endpoint de registro de sócios.
- Tema CSS estilo NES.

### Alterado
- UI migrada para tema escuro NES com grid responsivo.
- Tabela de jogos estendida com capa, revista e data de aquisição (migration `002`).

## [0.1.0] - 2026-03-03

### Adicionado
- Estrutura inicial do projeto Go com modelos base.
- Camada PostgreSQL com `pgx/v5` e interface `Store`.
- Migration `001_initial_schema.sql`.
- Cliente IGDB com Twitch OAuth2.
- Página de entrada e prateleira de jogos com SSR.
- Graceful shutdown. Licença GPL v3.
