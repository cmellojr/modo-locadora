# 🕹️ Modo Locadora

> **"Sopre a fita, pegue o controle e respeite o tempo."**

O **Modo Locadora** é um sistema web open-source, desenvolvido em **Go**, projetado para gerenciar sessões de jogos (jogatinas) simulando a experiência das videolocadoras brasileiras dos anos 90.

Aqui, a escassez gera valor, o tempo é um compromisso e o conhecimento é compartilhado através do lendário "caderninho".

---

## 📺 A Experiência

*   **O Balcão:** Um portal onde comunidades podem gerenciar suas próprias "unidades".
*   **Fitas Físicas Virtuais:** O acervo é limitado. Se todas as cópias de um jogo estiverem alugadas, você terá que esperar alguém devolver.
*   **Prateleiras Horizontais:** Navegue pelo acervo de Lançamentos e Catálogo em um layout que remete aos balcões das locadoras.
*   **O Caderninho:** Espaço pessoal para anotar passwords, códigos e mapas.
*   **Verso da Capa:** Deixe dicas públicas para os próximos jogadores ao devolver uma fita.
*   **Regra da Sexta:** Alugou na sexta? Só precisa devolver na segunda!

---

## ✨ Funcionalidades

*   **Autenticação de Sócios:** Sistema de login simples para identificação na locadora.
*   **Prateleira de Jogos (`/games`):** Visualização do acervo dividido entre `[NOVIDADES DA SEMANA]` e `[ARQUIVO HISTÓRICO]`.
*   **Abastecer Prateleiras (`/admin/stock`):** Interface administrativa para busca de metadados via API do IGDB e adição de novos jogos ao acervo.
*   **Gestão de Estoque:** Controle de cópias disponíveis para locação (escassas por design).

---

## 🛠️ Stack Técnica

Este projeto preza pela simplicidade, performance e legibilidade.

*   **Linguagem:** [Go](https://go.dev/) 1.24+ (seguindo o [Google Go Style Guide](https://google.github.io/styleguide/go/guide.html)).
*   **Banco de Dados:** [PostgreSQL](https://www.postgresql.org/) (utilizando `pgx`).
*   **Interface:** Server-Side Rendering (SSR) com `html/template`.
*   **Estilização:** [NES.css](https://nostalgic-css.github.io/NES.css/) para uma estética 8-bit e fonte "Press Start 2P".
*   **API de Dados:** [IGDB](https://api-docs.igdb.com/) para metadados e capas de jogos.

---

## 🚀 Como Executar

### Pré-requisitos

*   **Go 1.24** ou superior.
*   **Docker & Docker Compose** (para o banco de dados).
*   **Credenciais da API do IGDB** (através do [Twitch Developer Portal](https://dev.twitch.tv/)).

### Configuração

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/modo_locadora
TWITCH_CLIENT_ID=seu_client_id
TWITCH_CLIENT_SECRET=seu_client_secret
```

### Passo a Passo

1.  **Subir o Banco de Dados:**
    ```bash
    docker-compose up -d
    ```

2.  **Executar as Migrações:**
    (Atualmente as migrações são gerenciadas manualmente no diretório `internal/database/migrations`)

3.  **Iniciar o Servidor:**
    ```bash
    go run cmd/server/main.go
    ```

4.  **Acessar:**
    Abra `http://localhost:8080` no seu navegador.

---

## 📂 Estrutura do Projeto

*   `cmd/server/`: Ponto de entrada da aplicação.
*   `internal/`: Lógica de negócio, handlers, banco de dados e integrações (IGDB).
*   `web/templates/`: Arquivos HTML (SSR).
*   `web/static/`: Ativos estáticos (CSS, imagens).

---

## 📜 Licença

Distribuído sob a licença **GPL v3**. Veja `LICENSE` para mais informações.

---
*Desenvolvido com ❤️ por entusiastas da era de ouro dos videogames.*
