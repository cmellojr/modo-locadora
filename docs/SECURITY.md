# Política de Segurança

## Versões Suportadas

| Versão | Suportada          |
|--------|--------------------|
| main   | :white_check_mark: |

Apenas a versão mais recente em `main` recebe atualizações de segurança.

## Reportando Vulnerabilidades

**Não abra uma issue pública para vulnerabilidades de segurança.**

Reporte de forma privada por e-mail ao mantenedor do projeto ou use o [reporte privado de vulnerabilidades do GitHub](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing-information-about-vulnerabilities/privately-reporting-a-security-vulnerability).

Inclua: descrição, passos para reproduzir, impacto potencial e sugestão de correção.

## Autenticação e Sessões

- Senhas protegidas com **bcrypt** (custo padrão). Nunca logadas ou expostas em respostas da API.
- Sessões usam cookie assinado (`session_member`): `{member_uuid}.{hmac_sha256_hex}`
- Flags do cookie: `HttpOnly`, `SameSite=Strict`, `MaxAge=604800` (7 dias), `Path=/`
- `COOKIE_SECRET` deve ter pelo menos 32 caracteres.

## Autorização

| Escopo | Middleware | Verificação |
|--------|-----------|-------------|
| Rotas de sócio | `RequireAuth` | Cookie assinado válido |
| Rotas admin (`/admin/*`) | `RequireAdmin` | Cookie válido + e-mail bate com `ADMIN_EMAIL` |
| Rotas de turma (ações) | `RequireAuth` | Cookie assinado válido |
| Ações admin de turma | `RequireAuth` + verificação de cargo | Membro com role `admin` na turma |
| Exclusão de turma | `RequireAuth` + verificação de criador | `created_by` = sócio logado |

Requisições não autenticadas redirecionam para `/`. Usuários não-admin recebem `403 Forbidden`.

## Reputação do Sócio

- Aluguéis atrasados são auto-devolvidos por um job de background (intervalo de 5 minutos).
- Sócios infratores são marcados como `in_debt` com incremento permanente do `late_count`.
- Sócios em débito não podem alugar até se redimirem via `/membership/redeem`.
- O Painel da Vergonha na página de entrada exibe os maiores infratores.

## Integridade de Dados

- Operações de aluguel, devolução e aquisição usam **transações de banco de dados**.
- Todas as queries SQL usam placeholders parametrizados (sem interpolação de strings).

## Upload de Arquivos

- Uploads de capa e badges de turma restritos a arquivos de imagem (`accept="image/*"`).
- Tamanho máximo do formulário: 10 MB.
- Arquivos salvos com UUID como nome (previne path traversal).
- Capas de jogos: `web/static/covers/`. Badges de turmas: `web/static/clubs/`.

## Análise Estática

O projeto usa `golangci-lint` (`.golangci.yml`) com linters de segurança:
- **gosec** — Detecta problemas comuns de segurança em Go (credenciais hardcoded, criptografia fraca, padrões de SQL injection).
- **errcheck** — Garante que valores de retorno de erro são verificados.
- **staticcheck** — Captura construções suspeitas.

Execute `task check` ou `golangci-lint run ./...` antes de commits.

## Checklist de Deploy

- Defina um `COOKIE_SECRET` forte e aleatório (mínimo 32 caracteres).
- Defina `ADMIN_EMAIL` para restringir acesso admin.
- Use **HTTPS** em produção para proteger cookies e dados de formulário.
- Restrinja acesso ao banco apenas ao servidor da aplicação.
- Rotacione credenciais da API Twitch se comprometidas.
- Credenciais carregadas via `.env` (nunca commitado — listado no `.gitignore`).
