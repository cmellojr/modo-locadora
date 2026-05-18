# PRD: Modo Locadora

> **"A experiência definitiva do jogador honesto."**

## Visão Geral

O Modo Locadora é um simulador de locadora brasileira dos anos 1990 aplicado a backlog e sessões de jogos retrô. O produto combina catálogo, escassez, reputação e interação social para combater a "Síndrome do Labirinto": ter muitos jogos e terminar poucos.

Este PRD descreve intenção de produto. Implementação técnica e rotas ficam em [../ARCHITECTURE.md](../ARCHITECTURE.md) e [api.md](api.md).

## Pilares

- **Escassez real**: cada jogo tem cópias físicas limitadas.
- **Compromisso**: alugar cria prazo; atrasar gera consequência.
- **Memória retrogamer brasileira**: capas nacionais, revistas, locadora, cartucho e carteirinha.
- **Progressão social**: reputação, títulos, turmas e feed público.
- **Simplicidade técnica**: SSR, zero JavaScript e interface rápida.

## Requisitos Funcionais

### Acervo

- Navegação em 3 níveis: plataformas, prateleira filtrada e detalhe do jogo.
- Aquisição admin via IGDB.
- Upload de capa local ou via storage de produção.
- Edição de metadados, revista de origem e `cover_display`.
- Inventário admin com classificação de popularidade.

### Sócio

- Registro com matrícula `1991-XXX`.
- Login com perfil e senha.
- Carteirinha com estatísticas, status, título de progressão, caderno de passwords e turmas.
- Títulos calculados por histórico: `Sócio Novato`, `Sócio Prata`, `Sócio Ouro`, `Dono da Calçada`.

### Aluguel e Reputação

- Aluguel bloqueado para sócios em `in_debt`.
- Prazo padrão de 3 dias.
- Devolução com vereditos: `completed`, `enjoyed`, `quick_play`, `not_for_me`, `gave_up`.
- Auto-devolução para atrasos com `auto_return`.
- Redenção manual via "Sopro".
- Painel da Vergonha para maiores infratores.

### Social

- Feed "Aconteceu na Locadora".
- Almanaque do Tio com efemérides.
- Turmas públicas com badge, URL, admins e membros.
- Carteirinha lista as turmas do sócio.

## Requisitos Não Funcionais

- Go 1.24+ e PostgreSQL 15.
- SSR com `html/template`.
- Sem JavaScript no frontend.
- NES.css com tema escuro e estética 8-bit.
- Uploads independentes do ambiente via `StorageProvider`.
- Validação mínima antes de PR: `task check`.

## Fora do Escopo Atual

- SPA ou framework JavaScript.
- Sistema de pagamento.
- Chat em tempo real.
- Ranking competitivo global.
- API pública versionada.

## Plano

Versões e próximas ideias estão em [../ROADMAP.md](../ROADMAP.md). Mudanças notáveis ficam em [changelog.md](changelog.md).
