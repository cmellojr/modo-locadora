# Modo Locadora — Arquitetura do Sistema

## Visão Geral

Gerenciador de sessões retro-gaming que emula as videolocadoras brasileiras dos anos 90. A escassez é princípio de design: cada jogo tem cópias físicas limitadas — se todas estiverem alugadas, o jogo fica indisponível. Aluguéis atrasados são devolvidos automaticamente com penalização de reputação.

## Modelo de Entidades

```
Sócio (1991-XXX)
  ├── status: active | in_debt
  ├── late_count: contador permanente de penalidades
  ├── password_notes: caderno pessoal de códigos de jogos
  └── MemberTitle: progressão calculada (Novato → Prata → Ouro → Dono da Calçada)

Jogo (metadados IGDB)
  ├── platform, summary, cover_url, cover_display, source_magazine
  └── GameCopy (1:N)
        ├── status: available | rented
        └── Rental (1:N)
              ├── member_id, rented_at, due_at (3 dias)
              ├── returned_at (NULL = ativo)
              └── public_legacy (veredito: zerei | joguei_um_pouco | desisti)

Atividade (feed desnormalizado)
  ├── event_type: penalty | redemption | new_game | verdict_complete | verdict_partial | verdict_quit | club_created | club_joined
  ├── member_name, game_title
  └── created_at

Turma (comunidade gamer)
  ├── name, description, badge_url, website_url
  ├── created_by → Sócio (criador)
  └── ClubMember (M2M com Sócio)
        ├── role: admin | member
        └── joined_at
```

## Fluxo de Locação

```
1. Sócio navega /games → seleciona console → seleciona jogo → /games/{id}
2. Clica [ALUGAR] → POST /rent → cópia marcada como rented, aluguel criado (3 dias de prazo)
3. Detalhe do jogo mostra "ALUGADO - Com o Sócio: Nome"
4a. Admin visita /admin/returns → clica [Devolver] → cópia disponível novamente
4b. Sócio visita /membership → escolhe veredito (Zerei/Joguei/Desisti) → POST /membership/return
5. Veredito salvo em public_legacy, evento de atividade dispara no feed
6. Se atrasado: job de background auto-devolve, sócio recebe in_debt + late_count++
7. Sócio pode se redimir via POST /membership/redeem
```

## Mapa de Navegação

```
GET /                     → Login (Balcão) — redireciona para /games se autenticado
GET /games                → Grade de seleção de plataformas (Mega Drive, SNES, ...)
GET /games?platform=X     → Cartuchos da plataforma selecionada
GET /games/{id}           → Detalhe do jogo (stats, botão de aluguel)
GET /membership           → Carteirinha de sócio + caderno de passwords
GET /admin/stock          → Busca IGDB e aquisição de jogos
GET /admin/inventory      → Tabela do acervo com links de edição
GET /admin/edit/{id}      → Edição do jogo (upload de capa, metadados)
GET /admin/returns        → Check-in de aluguéis ativos
GET /clubs                → Listagem pública de turmas
GET /clubs/new            → Formulário de criação de turma (auth)
GET /clubs/{id}           → Detalhe da turma (público)
GET /clubs/{id}/edit      → Formulário de edição (admin da turma)
```

## Templates

| Template | Rota | Página |
|----------|------|--------|
| `index.html` | `GET /` | Login + Painel da Vergonha |
| `platforms.html` | `GET /games` | Layout 3 colunas: mini-card + vergonha, plataformas, atividades + almanaque |
| `games.html` | `GET /games?platform=X` | Prateleira de cartuchos (cards simplificados) |
| `game_detail.html` | `GET /games/{id}` | Detalhe do jogo + stats de aluguel |
| `carteirinha.html` | `GET /membership` | Carteirinha + badge de título + caderno + aluguéis ativos com auto-devolução |
| `admin_stock.html` | `GET /admin/stock` | Busca IGDB e aquisição |
| `admin_inventory.html` | `GET /admin/inventory` | Tabela do acervo com indicadores de saúde |
| `admin_edit.html` | `GET /admin/edit/{id}` | Formulário de edição + histórico de aluguéis |
| `admin_returns.html` | `GET /admin/returns` | Balcão de devoluções |
| `clubs.html` | `GET /clubs` | Listagem de turmas (grid de cards) |
| `club_detail.html` | `GET /clubs/{id}` | Detalhe da turma + tabela de membros |
| `club_form.html` | `GET /clubs/new`, `GET /clubs/{id}/edit` | Formulário de criação/edição de turma |

## Deploy

Build Docker multi-stage (`golang:1.24-alpine` → `alpine:3.21`). Docker Compose orquestra app + PostgreSQL com health checks. Três volumes: `postgres_data` (banco), `covers_data` (capas enviadas) e `clubs_data` (badges de turmas).

Para configuração do ambiente, veja [SETUP.md](docs/SETUP.md). Para convenções de código, veja [CONTRIBUTING.md](docs/CONTRIBUTING.md).
