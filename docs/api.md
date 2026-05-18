# Referência de Rotas e Endpoints

O Modo Locadora é majoritariamente SSR: as rotas retornam HTML renderizado no servidor, com alguns endpoints JSON. Para autenticação e autorização, veja [security.md](security.md). Para a arquitetura dos handlers, veja [../ARCHITECTURE.md](../ARCHITECTURE.md).

## Páginas SSR

| Rota | Acesso | Descrição |
|------|--------|-----------|
| `GET /` | Público | Balcão com login, boas-vindas para sócios logados e Painel da Vergonha |
| `GET /games` | Público | Grade de plataformas, feed e almanaque |
| `GET /games?platform=X` | Público | Prateleira filtrada por plataforma |
| `GET /games/{id}` | Público | Detalhe do jogo, disponibilidade, estatísticas e botão de aluguel |
| `GET /membership` | Sócio | Carteirinha, notas, aluguéis ativos, vereditos e turmas |
| `GET /admin/stock` | Admin | Busca IGDB e aquisição de jogos |
| `GET /admin/inventory` | Admin | Inventário com classificação de popularidade |
| `GET /admin/edit/{id}` | Admin | Edição do jogo, upload de capa e histórico de aluguéis |
| `GET /admin/returns` | Admin | Dashboard de devoluções |
| `GET /clubs` | Público | Listagem pública de turmas |
| `GET /clubs/new` | Sócio | Formulário de criação de turma |
| `GET /clubs/{id}` | Público | Detalhe da turma e membros |
| `GET /clubs/{id}/edit` | Admin da turma | Formulário de edição da turma |

## Endpoints de Formulário

Todos usam `application/x-www-form-urlencoded`, exceto uploads com `multipart/form-data`.

### Auth

| Método e rota | Acesso | Campos | Resultado |
|---------------|--------|--------|-----------|
| `POST /login` | Público | `profile_name`, `password` | Define cookie `session_member` e redireciona para `/games` |
| `POST /logout` | Público | Nenhum | Limpa cookie e redireciona para `/` |

### Sócio

| Método e rota | Acesso | Campos | Resultado |
|---------------|--------|--------|-----------|
| `POST /rent` | Sócio | `game_id` | Aluga uma cópia disponível e redireciona para `/games/{id}` |
| `POST /membership/notes` | Sócio | `notes` | Salva caderno de passwords |
| `POST /membership/return` | Sócio | `rental_id`, `verdict` | Devolve aluguel com veredito |
| `POST /membership/redeem` | Sócio | Nenhum | Limpa `in_debt` |

Vereditos aceitos por `POST /membership/return`:

| Valor | Texto de UI |
|-------|-------------|
| `completed` | Zerei |
| `enjoyed` | Curti |
| `quick_play` | Joguei rápido |
| `not_for_me` | Não é pra mim |
| `gave_up` | Desisti |

### Admin

| Método e rota | Acesso | Campos principais | Resultado |
|---------------|--------|-------------------|-----------|
| `POST /admin/purchase` | Admin | `title`, `igdb_id`, `platform`, `summary`, `cover_url`, `magazine` | Cria jogo e primeira cópia |
| `POST /admin/update-game` | Admin | `id`, `title`, `platform`, `summary`, `magazine`, `cover_url`, `cover_display`, `cover_file` | Atualiza metadados e capa |
| `POST /admin/return-game` | Admin | `rental_id` | Processa devolução admin |

`cover_display` aceita `cover`, `contain` ou `fill`.

### Turmas

| Método e rota | Acesso | Campos | Resultado |
|---------------|--------|--------|-----------|
| `POST /clubs` | Sócio | `name`, `description`, `website_url`, `badge_file` | Cria turma e torna o criador admin |
| `POST /clubs/{id}/edit` | Admin da turma | `name`, `description`, `website_url`, `badge_file` | Atualiza turma |
| `POST /clubs/{id}/join` | Sócio | Nenhum | Entra na turma |
| `POST /clubs/{id}/leave` | Sócio | Nenhum | Sai da turma; último admin não pode sair |
| `POST /clubs/{id}/promote` | Admin da turma | `member_id` | Promove membro a admin |
| `POST /clubs/{id}/remove` | Admin da turma | `member_id` | Remove membro |
| `POST /clubs/{id}/delete` | Criador | Nenhum | Exclui turma |

## API JSON

### `POST /members`

Registra um novo sócio.

```json
{
  "profile_name": "Player1",
  "email": "player1@locadora.com",
  "password": "secret123",
  "favorite_console": "SNES"
}
```

Resposta `201 Created`: objeto do sócio com `MembershipNumber` auto-atribuído. `PasswordHash` não é retornado.

### `GET /search?q={query}`

Busca jogos na IGDB. Retorna até 10 resultados com nome, resumo, capa e plataformas.

## Convenções de Resposta

- Sucesso em formulários geralmente usa `303 See Other`.
- Falhas de autorização retornam `403 Forbidden`.
- Sócios em débito tentando alugar são redirecionados para o detalhe do jogo com `?error=in_debt`.
- Alguns fluxos usam `?success=...` para exibir mensagens nos templates.
