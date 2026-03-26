# PRD: Modo Locadora

**"A experiencia definitiva do jogador honesto"**

## 1. Visao Geral

O Modo Locadora e um simulador de ecossistema de locadora brasileira dos anos 90. Ele combina a robustez da engenharia moderna (Go, Docker, Postgres) com a alma dos projetos classicos como Projeto Jogatina e NES Archive. O sistema nao e apenas um catalogo, mas uma ferramenta de curadoria que combate a "Sindrome do Labirinto" atraves de mecanicas de escassez e reputacao social.

## 2. Pilares do Produto

- **Escassez Real:** Cartuchos sao itens finitos. Se estiver alugado, o socio deve esperar.
- **Identidade Nacional:** Prioridade para capas TecToy/Playtronic e referencias a revistas como Acao Games e VideoGame.
- **Consequencia Social:** Acoes do socio (atrasos, zeramentos) impactam sua reputacao publica na "comunidade".
- **Imersao Retro:** Interface 8-bit funcional, sonora e visual (NES.css).

## 3. Requisitos Funcionais

### 3.1. Gestao de Acervo (O Balcao)

- **Navegacao em 3 Niveis:**
  1. Selecao de Plataforma (com logos SVG).
  2. Prateleira de Jogos (Grid com capas brasileiras).
  3. Detalhe do Titulo (Stats, curiosidades e botao de aluguel).
- **Upload de Capas:** Sistema para o Admin subir imagens locais (TecToy) via Multipart Form.
- **Sistema de Copias:** Controle rigido de unidades disponiveis por titulo.

### 3.2. Experiencia do Socio (A Carteirinha)

- **Identidade:** Numeracao sequencial `1991-XXX`.
- **Status de Progressao:** Titulos automaticos: Novato, Socio Prata (10+ devolucoes) e Dono da Calcada (5+ jogos zerados).
- **Caderno de Passwords:** Campo de texto persistente para anotacoes e codigos de jogos.
- **Estrela Dourada:** Badge visual para jogos marcados com o veredito "Zerei".

### 3.3. Dinamica de Locacao e Reputacao

- **Auto-Return System:** Job de segundo plano que penaliza atrasos automaticamente.
- **Painel da Vergonha:** Exposicao publica de socios com status "Em Debito".
- **Fluxo de Redencao:** Botao "Soprar Cartucho" para limpar pendencias e restaurar acesso.
- **Veredito de Devolucao:** O socio deve classificar sua experiencia (Zerei, Joguei um pouco, Desisti).

### 3.4. Social e Conteudo (Feed)

- **Aconteceu na Locadora:** Feed de atividades em tempo real (Novas fitas, Vereditos, Punicoes).
- **Almanaque do Tio:** Noticias historicas baseadas em efemerides reais da industria de games.

## 4. Requisitos Nao-Funcionais (SRE Stack)

- **Tecnologia:** Go 1.22+, PostgreSQL 15, Docker Compose.
- **Seguranca:** Senhas protegidas com bcrypt; Cookies assinados com HMAC-SHA256.
- **Performance:** Interface estritamente SSR (Server-Side Rendering) para manter a leveza e velocidade.
- **Design:** CSS Grid responsivo com fidelidade aos componentes NES.css.
- **Audio:** Feedback sonoro via Web Audio API (Ondas Quadradas/8-bit).

## 5. Arquitetura de Dados (Principais Tabelas)

| Tabela | Descricao |
|--------|-----------|
| `members` | Dados do socio, reputacao, status e contadores |
| `games` | Informacoes de catalogo, capas e metadados |
| `game_copies` | Instancias fisicas dos cartuchos e seu estado atual |
| `rentals` | Historico de locacao, prazos e vereditos |
| `activities` | Logs de eventos para o feed social |

## 6. Plano de Lancamento (Roadmap)

| Versao | Status | Escopo |
|--------|--------|--------|
| V0.1-0.3 | Concluido | Base tecnica, Auth, Docker e Prateleira Basica |
| V0.4 | Atual | Implementacao de Vereditos, Feed, Estrela Dourada e Seed Acao Games #1 |
| V0.5 | Proximo | Implementacao da "Roleta do Tio", Mural de Placares e Titulos de Status |

---

*&copy; 1991-2026 Modo Locadora - Inspirado no Projeto Jogatina e Forum NES Archive.*
