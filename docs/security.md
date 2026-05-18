# Política de Segurança

## Versões Suportadas

Apenas a versão mais recente em `main` recebe correções de segurança.

## Reportando Vulnerabilidades

Não abra issues públicas para vulnerabilidades. Use reporte privado pelo GitHub ou contate o mantenedor em privado.

Inclua descrição, passos para reproduzir, impacto potencial e sugestão de correção.

## Autenticação e Sessões

- Senhas usam **bcrypt**.
- Sessões usam cookie assinado `session_member`: `{member_uuid}.{hmac_sha256_hex}`.
- Flags do cookie: `HttpOnly`, `SameSite=Strict`, `MaxAge=604800`, `Path=/`.
- `COOKIE_SECRET` deve ter 32+ caracteres em produção.
- `POST /logout` remove o cookie de sessão.

## Autorização

| Escopo | Proteção |
|--------|----------|
| Rotas de sócio | `RequireAuth` |
| Rotas admin (`/admin/*`) | `RequireAdmin` + e-mail igual a `ADMIN_EMAIL` |
| Ações de turma | `RequireAuth` |
| Administração de turma | Verificação de role `admin` em `club_members` |
| Exclusão de turma | Apenas `created_by` |

Requisições não autenticadas redirecionam para `/`. Usuários sem permissão recebem `403 Forbidden`.

## Reputação

- O job em `internal/jobs/overdue.go` verifica atrasos a cada 5 minutos.
- Aluguéis vencidos são auto-devolvidos com veredito `auto_return`.
- O sócio recebe `status = in_debt` e `late_count` é incrementado.
- `POST /membership/redeem` limpa o status de débito, mas mantém o histórico em `late_count`.

## Integridade de Dados

- Aluguel, devolução e aquisição usam transações no PostgreSQL.
- Queries usam placeholders parametrizados.
- Identificadores são UUIDs.
- `public_legacy` mantém os vereditos normalizados em inglês.

## Uploads

Uploads de capa e badge passam por `internal/storage.StorageProvider`.

| Ambiente | Provider | Destino |
|----------|----------|---------|
| Local/dev | `LocalStorage` | `web/static/covers/` e `web/static/clubs/` |
| Produção com `APP_ENV=production` | `GCSStorage` | Bucket `STORAGE_BUCKET_NAME` |

Regras atuais:

- Formulários aceitam `image/*`.
- Tamanho máximo do formulário: 10 MB.
- Nomes de arquivo usam UUID para reduzir risco de path traversal e colisão.
- Em GCS, objetos recebem ACL pública para exibição direta pela UI.

## Análise Estática

O projeto usa `.golangci.yml` com `errcheck`, `staticcheck`, `unused`, `gosec`, `govet`, `ineffassign` e `typecheck`.

```bash
task check
```

## Checklist de Deploy

- Defina `COOKIE_SECRET` forte.
- Defina `ADMIN_EMAIL`.
- Use HTTPS.
- Proteja o PostgreSQL fora da internet pública.
- Defina credenciais Twitch/IGDB de produção.
- Para GCS: configure `APP_ENV=production`, `STORAGE_BUCKET_NAME` e credenciais do Google Cloud no ambiente.
- Não commite `.env` ou credenciais.
