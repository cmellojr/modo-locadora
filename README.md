# Modo Locadora

> **"Sopre a fita, pegue o controle e respeite o tempo."**

O **Modo Locadora** e um sistema web open-source, desenvolvido em **Go**, projetado para gerenciar sessoes de jogos (jogatinas) simulando a experiencia das videolocadoras brasileiras dos anos 90.

Aqui, a escassez gera valor, o tempo e um compromisso e o conhecimento e compartilhado atraves do lendario "caderninho".

---

## A Experiencia

*   **O Balcao:** Um portal onde comunidades podem gerenciar suas proprias "unidades".
*   **Fitas Fisicas Virtuais:** O acervo e limitado. Se a copia de um jogo estiver alugada, voce tera que esperar alguem devolver.
*   **Prateleira:** Navegue pelo acervo de cartuchos em um layout que remete aos balcoes das locadoras.
*   **Carteirinha de Socio:** Cada membro recebe um numero de socio no formato `1991-XXX` e pode consultar sua carteirinha digital.
*   **O Caderninho:** Espaco pessoal para anotar passwords, codigos e mapas *(em breve)*.
*   **Verso da Capa:** Deixe dicas publicas para os proximos jogadores ao devolver uma fita *(em breve)*.
*   **Regra da Sexta:** Alugou na sexta? So precisa devolver na segunda!

---

## Funcionalidades

*   **Cadastro e Autenticacao de Socios:** Registro via API com senha bcrypt e login por nome + senha.
*   **Carteirinha Digital (`/carteirinha`):** Cartao de socio com numero de matricula `1991-XXX`, console favorito e data de ingresso.
*   **Prateleira de Jogos (`/games`):** Visualizacao do acervo com status de disponibilidade em tempo real. Socios autenticados podem alugar diretamente pela prateleira.
*   **Sistema de Aluguel:** Botao [ALUGAR] para jogos disponiveis. Jogos alugados exibem "Com o Socio: Nome".
*   **Balcao de Devolucoes (`/admin/returns`):** Area administrativa para dar baixa nas fitas alugadas.
*   **Abastecer Prateleiras (`/admin/stock`):** Busca de metadados via API do IGDB e adicao de novos jogos ao acervo.
*   **Acervo (`/admin/inventory`):** Listagem completa do catalogo com opcao de edicao para cada jogo.
*   **Edicao de Jogos (`/admin/edit/{id}`):** Formulario para traduzir e ajustar dados importados do IGDB.

---

## Stack Tecnica

Este projeto preza pela simplicidade, performance e legibilidade.

*   **Linguagem:** [Go](https://go.dev/) 1.24+ (seguindo o [Google Go Style Guide](https://google.github.io/styleguide/go/guide.html)).
*   **Banco de Dados:** [PostgreSQL](https://www.postgresql.org/) 15+ (utilizando `pgx/v5`).
*   **Interface:** Server-Side Rendering (SSR) com `html/template`.
*   **Estilizacao:** [NES.css](https://nostalgic-css.github.io/NES.css/) para uma estetica 8-bit e fonte "Press Start 2P".
*   **API de Dados:** [IGDB](https://api-docs.igdb.com/) para metadados e capas de jogos.
*   **Seguranca:** bcrypt para senhas, HMAC-SHA256 para cookies, middleware de autorizacao.

---

## Como Executar

### Pre-requisitos

*   **Go 1.24** ou superior.
*   **Docker & Docker Compose** (para o banco de dados).
*   **Credenciais da API do IGDB** (atraves do [Twitch Developer Portal](https://dev.twitch.tv/)).

### Configuracao

Copie o arquivo de exemplo e preencha com seus valores:

```bash
cp .env.example .env
```

Variaveis necessarias:

```env
TWITCH_CLIENT_ID=seu_client_id
TWITCH_CLIENT_SECRET=seu_client_secret
DATABASE_URL=postgres://tio_da_locadora:sopre_a_fita@localhost:5432/modo_locadora?sslmode=disable
COOKIE_SECRET=gere_uma_chave_secreta_aleatoria
ADMIN_EMAIL=seu_email_de_admin@example.com
```

### Passo a Passo

1.  **Subir o Banco de Dados:**
    ```bash
    docker compose up -d
    ```

2.  **Executar as Migracoes:**
    ```bash
    psql $DATABASE_URL -f internal/database/migrations/001_initial_schema.sql
    psql $DATABASE_URL -f internal/database/migrations/002_update_games_table.sql
    psql $DATABASE_URL -f internal/database/migrations/003_membership_and_rental_support.sql
    ```

3.  **Criar o Primeiro Socio (admin):**
    ```bash
    curl -X POST http://localhost:8080/members \
      -H "Content-Type: application/json" \
      -d '{
        "profile_name": "Tio da Locadora",
        "email": "seu_email_de_admin@example.com",
        "password": "sua_senha",
        "favorite_console": "Mega Drive"
      }'
    ```
    O email deve coincidir com `ADMIN_EMAIL` para acesso as rotas administrativas.

4.  **Iniciar o Servidor:**
    ```bash
    go run ./cmd/server
    ```

5.  **Acessar:**
    Abra `http://localhost:8080` no seu navegador.

---

## Estrutura do Projeto

```
modo-locadora/
├── cmd/server/main.go              # Ponto de entrada
├── internal/
│   ├── auth/                       # Assinatura de cookies (HMAC-SHA256)
│   ├── config/                     # Carregamento de .env
│   ├── database/
│   │   ├── store.go                # Interface Store
│   │   ├── postgres.go             # Implementacao PostgreSQL
│   │   └── migrations/             # Migracoes SQL (001-003)
│   ├── handlers/                   # Handlers HTTP
│   ├── igdb/                       # Cliente IGDB
│   ├── middleware/                 # Middleware de auth e admin
│   └── models/                     # Entidades (Member, Game, GameCopy, Rental)
├── web/
│   ├── static/css/                 # Tema NES retro
│   └── templates/                  # Templates HTML (PT-BR)
├── docs/                           # Documentacao do projeto
├── docker-compose.yml              # Container PostgreSQL
├── .env.example                    # Template de variaveis de ambiente
├── ARCHITECTURE.md                 # Arquitetura do sistema
└── go.mod                          # Modulo Go
```

---

## Licenca

Distribuido sob a licenca **GPL v3**. Veja `LICENSE` para mais informacoes.

---
*Desenvolvido pelo Tio da Locadora.*
