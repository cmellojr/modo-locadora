# PRD: Modo Locadora

**"A experiência definitiva do jogador honesto"**

## 1. Visão Geral

O Modo Locadora é um simulador de ecossistema de locadora brasileira dos anos 90. Ele combina a robustez da engenharia moderna (Go, Docker, Postgres) com a alma dos projetos clássicos como Projeto Jogatina e NES Archive. O sistema não é apenas um catálogo, mas uma ferramenta de curadoria que combate a "Síndrome do Labirinto" através de mecânicas de escassez e reputação social.

## 2. Pilares do Produto

- **Escassez Real:** Cartuchos são itens finitos. Se estiver alugado, o sócio deve esperar.
- **Identidade Nacional:** Prioridade para capas TecToy/Playtronic e referências a revistas como Ação Games e VideoGame.
- **Consequência Social:** Ações do sócio (atrasos, zeramentos) impactam sua reputação pública na "comunidade".
- **Imersão Retro:** Interface 8-bit funcional, sonora e visual (NES.css).

## 3. Requisitos Funcionais

### 3.1. Gestão de Acervo (O Balcão)

- **Navegação em 3 Níveis:**
  1. Seleção de Plataforma (com logos SVG).
  2. Prateleira de Jogos (Grid com capas brasileiras).
  3. Detalhe do Título (Stats, curiosidades e botão de aluguel).
- **Upload de Capas:** Sistema para o Admin subir imagens locais (TecToy) via Multipart Form.
- **Sistema de Cópias:** Controle rígido de unidades disponíveis por título.

### 3.2. Experiência do Sócio (A Carteirinha)

- **Identidade:** Numeração sequencial `1991-XXX`.
- **Status de Progressão:** Títulos automáticos: Sócio Novato, Sócio Prata (10+ devoluções no prazo), Sócio Ouro (25+ devoluções no prazo) e Dono da Calçada (5+ jogos zerados). Devedores veem título esmaecido com indicador.
- **Caderno de Passwords:** Campo de texto persistente para anotações e códigos de jogos.
- **Estrela Dourada:** Badge visual para jogos marcados com o veredito "Zerei".

### 3.3. Dinâmica de Locação e Reputação

- **Auto-Return System:** Job de segundo plano que penaliza atrasos automaticamente.
- **Painel da Vergonha:** Exposição pública de sócios com status "em débito".
- **Fluxo de Redenção:** Botão "Soprar Cartucho" para limpar pendências e restaurar acesso.
- **Veredito de Devolução:** O sócio deve classificar sua experiência (Zerei, Joguei um pouco, Desisti).

### 3.4. Social e Conteúdo (Feed)

- **Aconteceu na Locadora:** Feed de atividades em tempo real (Novas fitas, Vereditos, Punições, Turmas).
- **Almanaque do Tio:** Notícias históricas baseadas em efemérides reais da indústria de games.

### 3.5. Turmas (Comunidades Gamers)

- **Criação de Turmas:** Sócios podem criar turmas representando podcasts, canais YouTube, grupos WhatsApp ou qualquer comunidade gamer. Cada turma tem nome, descrição, badge (upload de imagem) e URL.
- **Participação Múltipla:** Sócios podem participar de quantas turmas quiserem.
- **Administração Distribuída:** Turmas podem ter múltiplos admins. Admins podem promover membros e remover participantes. O criador é automaticamente o primeiro admin.
- **Listagem Pública:** Qualquer visitante pode ver as turmas. Ações (criar, entrar, sair, editar) requerem login.
- **Exclusão:** Apenas o criador original pode excluir uma turma.
- **Integração com Carteirinha:** Seção "MINHAS TURMAS" exibe turmas do sócio com badge, nome e cargo.

### 3.6. Auditoria do Acervo (Admin)

- **Saúde do Acervo:** Indicador visual no inventário: Cartucho Novo (0-1 aluguéis), Clássico Eterno (<25% ruins), Precisa Soprar (25-49%), Fita Gasta (50%+).
- **Histórico de Aluguéis:** Ficha do jogo mostra últimos 5 aluguéis com sócio, datas, veredito e indicador de atraso.
- **Modo de Exibição:** Campo `cover_display` controla CSS object-fit das capas (preencher, mostrar inteira, esticar).

## 4. Requisitos Não-Funcionais (SRE Stack)

- **Tecnologia:** Go 1.24+, PostgreSQL 15, Docker Compose.
- **Segurança:** Senhas protegidas com bcrypt; Cookies assinados com HMAC-SHA256.
- **Performance:** Interface estritamente SSR (Server-Side Rendering) para manter a leveza e velocidade.
- **Design:** CSS Grid responsivo com fidelidade aos componentes NES.css.
- **Áudio:** Feedback sonoro via Web Audio API (Ondas Quadradas/8-bit).

## 5. Arquitetura de Dados (Principais Tabelas)

| Tabela | Descrição |
|--------|-----------|
| `members` | Dados do sócio, reputação, status e contadores |
| `games` | Informações de catálogo, capas e metadados |
| `game_copies` | Instâncias físicas dos cartuchos e seu estado atual |
| `rentals` | Histórico de locação, prazos e vereditos |
| `activities` | Logs de eventos para o feed social |
| `clubs` | Turmas (comunidades gamers) com badge e URL |
| `club_members` | Relação M2M entre turmas e sócios (com cargo) |

## 6. Plano de Lançamento

Veja [ROADMAP.md](../ROADMAP.md) para o plano detalhado de versões e milestones.

---

*&copy; 1991-2026 Modo Locadora - Inspirado no Projeto Jogatina e Fórum NES Archive.*
