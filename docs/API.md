# Referência da API

O Modo Locadora usa páginas renderizadas no servidor (HTML) e alguns endpoints JSON. Para detalhes de autenticação, veja [Política de Segurança](SECURITY.md).

## Páginas (SSR)

### `GET /`

Página de entrada (Balcão) com formulário de login e Painel da Vergonha (maiores devedores). Sócios autenticados são redirecionados para `/games`.

### `GET /games`

Sem parâmetros: layout 3 colunas com mini-card do sócio + Painel da Vergonha (esquerda), grade de plataformas (centro), feed de atividades + almanaque (direita).

Com `?platform=X`: cards simplificados de cartucho para o console — capa, título, número de cópias, disponibilidade e estrela dourada para jogos completados. Cada card leva à página de detalhe.

### `GET /games/{id}`

Página de detalhe do jogo. Mostra capa, título, plataforma, resumo, revista de origem, disponibilidade de cópias, total de aluguéis, fã número 1, sócio atual e data de aquisição. Sócios logados veem o botão [ALUGAR] se houver cópias disponíveis.

Parâmetro: `error=em_debito` exibe aviso de débito.

### `GET /carteirinha`

Carteirinha digital de sócio. Requer autenticação. Mostra número de matrícula, título de progressão (Sócio Novato / Prata / Ouro / Dono da Calçada), perfil, stats de aluguel, status, caderno de passwords e aluguéis ativos com auto-devolução (seleção de veredito).

Parâmetro: `success` exibe notificação.

### `GET /admin/stock`

Busca IGDB e página de aquisição de jogos. Requer acesso de administrador. Parâmetros: `q`, `magazine`, `selected`, `success`.

### `GET /admin/inventory`

Tabela completa do acervo com botões de edição e indicadores de saúde (Cartucho Novo, Clássico Eterno, Precisa Soprar, Fita Gasta). Requer acesso de administrador. Parâmetro: `success`.

### `GET /admin/edit/{id}`

Formulário de edição do jogo com upload de capa (multipart) e seletor de modo de exibição. Mostra histórico de aluguéis (últimos 5 registros). Requer acesso de administrador.

### `GET /admin/returns`

Dashboard de aluguéis ativos com botões de devolução. Requer acesso de administrador. Parâmetro: `success`.

---

## Endpoints de Formulário

### `POST /login`

Autenticar e definir cookie de sessão. Content-Type: `application/x-www-form-urlencoded`.

| Campo | Descrição |
|-------|-----------|
| `profile_name` | Nome de perfil do sócio |
| `password` | Senha do sócio |

**Sucesso:** redireciona (303) para `/games`, define cookie `session_member`.

### `POST /rent`

Alugar um jogo. Requer autenticação.

| Campo | Descrição |
|-------|-----------|
| `game_id` | UUID do jogo |

**Sucesso:** redireciona (303) para `/games/{id}`. Sócios em débito são redirecionados com `?error=em_debito`.

### `POST /carteirinha/notes`

Salvar caderno de passwords. Requer autenticação.

| Campo | Descrição |
|-------|-----------|
| `notes` | Texto do caderno de passwords |

**Sucesso:** redireciona (303) para `/carteirinha?success=1`.

### `POST /carteirinha/return`

Auto-devolução de aluguel com veredito. Requer autenticação.

| Campo | Descrição |
|-------|-----------|
| `rental_id` | UUID do aluguel |
| `verdict` | Status de jogo: `zerei`, `joguei_um_pouco` ou `desisti` |

**Sucesso:** redireciona (303) para `/carteirinha?success=devolucao`. Dispara evento de atividade baseado no veredito.

### `POST /carteirinha/redeem`

Limpar status de débito do sócio. Requer autenticação. Sem campos.

**Sucesso:** redireciona (303) para `/carteirinha?success=redencao`.

### `POST /admin/purchase`

Adicionar um jogo do IGDB ao acervo. Requer acesso de administrador. Cria uma `game_copy` atomicamente.

| Campo | Descrição |
|-------|-----------|
| `title` | Título do jogo |
| `igdb_id` | ID do jogo no IGDB |
| `platform` | Nome da plataforma (padrão "N/A") |
| `summary` | Descrição do jogo |
| `cover_url` | URL da capa |
| `magazine` | Revista de origem |

**Sucesso:** redireciona (303) para `/admin/edit/{id}`.

### `POST /admin/update-game`

Atualizar dados do jogo. Requer acesso de administrador. Content-Type: `multipart/form-data` (suporta upload de capa).

| Campo | Descrição |
|-------|-----------|
| `id` | UUID do jogo |
| `title` | Título do jogo |
| `platform` | Nome da plataforma |
| `summary` | Descrição |
| `magazine` | Revista de origem |
| `cover_url` | URL da capa existente (hidden, fallback) |
| `cover_display` | Modo CSS object-fit: `cover` (padrão), `contain` ou `fill` |
| `cover_file` | Arquivo de imagem da capa (opcional) |

**Sucesso:** redireciona (303) para `/admin/inventory?success={title}`.

### `POST /admin/return-game`

Processar devolução de jogo. Requer acesso de administrador.

| Campo | Descrição |
|-------|-----------|
| `rental_id` | UUID do aluguel |

**Sucesso:** redireciona (303) para `/admin/returns?success=Fita+devolvida`.

---

## API JSON

### `POST /members`

Registrar um novo sócio.

```json
{
  "profile_name": "Player1",
  "email": "player1@locadora.com",
  "password": "secret123",
  "favorite_console": "SNES"
}
```

**Resposta** `201 Created`: objeto do sócio com `MembershipNumber` auto-atribuído. `PasswordHash` é sempre vazio.

### `GET /search?q={query}`

Buscar na base IGDB. Retorna até 10 resultados com nome, resumo, capa e plataformas.
