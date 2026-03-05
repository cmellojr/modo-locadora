# Modo Locadora

```
 ╔╦╗╔═╗╔╦╗╔═╗  ╦  ╔═╗╔═╗╔═╗╔╦╗╔═╗╦═╗╔═╗
 ║║║║ ║ ║║║ ║  ║  ║ ║║  ╠═╣ ║║║ ║╠╦╝╠═╣
 ╩ ╩╚═╝═╩╝╚═╝  ╩═╝╚═╝╚═╝╩ ╩═╩╝╚═╝╩╚═╩ ╩
```

> **"Sopre a fita, pegue o controle e respeite o tempo."**

---

Era sexta-feira. Você corria da escola, passava no balcão, mostrava a carteirinha e rezava para aquele cartucho ainda estar na prateleira. O tio conferia o caderno, pegava a fita, soprava o conector e dizia: *"Devolve na segunda, hein?"*

O **Modo Locadora** é uma homenagem a essa era dourada. Um sistema web que recria a experiência das videolocadoras brasileiras dos anos 90 — onde a escassez gerava valor, o tempo era um compromisso, e o conhecimento era compartilhado num caderninho de passwords grudento de guaraná.

Este projeto é para quem ainda lembra do cheiro de plástico dos cartuchos, da emoção de encontrar a última cópia de Mega Man 2 na prateleira, e da comunidade que se formava em torno daquele balcão.

---

## A Locadora

**O Balcão** — A tela de entrada. Você chega, mostra sua carteirinha e entra.

**A Prateleira** — O acervo de cartuchos, com capas pixeladas e status em tempo real. Se a fita estiver alugada, vai ter que esperar o próximo sócio devolver.

**A Carteirinha de Sócio** — Cada membro recebe um número no formato `1991-XXX`. É digital, mas carrega o espírito daquele cartão plastificado com foto 3x4.

**O Caderno de Passwords** — Espaço pessoal para anotar senhas, códigos e mapas. Porque ninguém merece perder o progresso do Metroid.

**O Tio da Locadora** — O administrador. Ele abastece as prateleiras, cuida do acervo e dá baixa nas devoluções. Tudo pelo painel admin, como se estivesse atrás do balcão.

---

## O Espírito da Coisa

Este não é um sistema genérico de catálogo. É uma experiência. Cada detalhe foi pensado para evocar aquela época:

- **Escassez real** — Cada jogo tem cópias limitadas. Se todas estiverem alugadas, o jogo fica indisponível. Assim como era na locadora.
- **Visual 8-bit** — Interface com [NES.css](https://nostalgic-css.github.io/NES.css/) e fonte Press Start 2P. Cada pixel no lugar.
- **Sem JavaScript** — Renderização no servidor, como os sites de 1996. Rápido, limpo, sem frescura.
- **Copyleft** — Licenciado sob GPL v3. O código é livre, como deveria ser.

---

## Stack

| Componente | Tecnologia |
|------------|-----------|
| Backend | [Go](https://go.dev/) 1.24+ |
| Banco de dados | [PostgreSQL](https://www.postgresql.org/) 15+ |
| Interface | Server-Side Rendering com `html/template` |
| Estilo | [NES.css](https://nostalgic-css.github.io/NES.css/) 2.3.0 + Press Start 2P |
| Dados de jogos | [IGDB](https://api-docs.igdb.com/) via Twitch OAuth2 |
| Segurança | bcrypt + HMAC-SHA256 + middleware de autorização |

---

## Começar a Jogar

### Pré-requisitos

- Go 1.24+
- Docker (para o PostgreSQL)
- Credenciais da [API do IGDB](https://dev.twitch.tv/)

### Início rápido

```bash
git clone https://github.com/cmellojr/modo-locadora.git
cd modo-locadora
cp .env.example .env        # preencha com seus valores
docker compose up -d         # sobe o banco
go run ./cmd/server          # abre a locadora em http://localhost:8080
```

Para o guia completo de configuração (migrações, primeiro sócio, solução de problemas), veja **[docs/SETUP.md](docs/SETUP.md)**.

---

## Documentação

| Documento | Conteúdo |
|-----------|---------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | Visão geral da arquitetura, entidades e fluxos |
| [docs/SETUP.md](docs/SETUP.md) | Guia completo de instalação e configuração |
| [docs/API.md](docs/API.md) | Referência de endpoints (SSR e JSON) |
| [docs/SECURITY.md](docs/SECURITY.md) | Política de segurança e práticas |
| [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) | Como contribuir com o projeto |
| [docs/CHANGELOG.md](docs/CHANGELOG.md) | Histórico de mudanças |

---

## Funcionalidades Futuras

- **Verso da Capa** — Deixe dicas públicas para os próximos jogadores ao devolver uma fita.
- **Regra da Sexta** — Alugou na sexta? Só precisa devolver na segunda!
- **Ranking de Sócios** — Quem mais alugou, quem mais devolveu no prazo.
- **Coleção Pessoal** — Marque os jogos que você já zerou.

---

## Licença

Distribuído sob a licença **GPL v3**. Veja [LICENSE](LICENSE) para mais informações.

---

*Desenvolvido com nostalgia pelo Tio da Locadora.*

*Em memória de todas as locadoras que fecharam, mas nunca foram esquecidas.*
