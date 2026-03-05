# Modo Locadora

```
  ____  __  __  ___   ____    ___       _     ___    ____   ____  ____    ___   ____    ____
 |    \|  ||  |/   \ |    \  /   \     | |   /   \  /    | /    ||    \  /   \ |    \  /    |
 |  _  |  ||  |     ||  o  )|     |    | |  |     ||   __||  o  ||  D  )|     ||  D  )|  o  |
 |  |  |  ||  |  O  ||   _/ |  O  |    | |  |  O  ||  |   |     ||    / |  O  ||    / |     |
 |  |  |  ||  |     ||  |   |     |    | |  |     ||  |__ |  _  ||    \ |     ||    \ |  _  |
 |  |  |  ||  |     ||  |   |     |    | |  |     ||     ||  |  ||  .  \|     ||  .  \|  |  |
 |__|__|__||__|\_____/|__|    \___/     |_|   \___/ |_____||__|__||__|\_|\___/ |__|\_||__|__|
```

> **"Sopre a fita, pegue o controle e respeite o tempo."**

---

Era sexta-feira. Voce corria da escola, passava no balcao, mostrava a carteirinha e rezava para aquele cartucho ainda estar na prateleira. O tio conferia o caderno, pegava a fita, soprava o conector e dizia: *"Devolve na segunda, hein?"*

O **Modo Locadora** e uma homenagem a essa era dourada. Um sistema web que recria a experiencia das videolocadoras brasileiras dos anos 90 — onde a escassez gerava valor, o tempo era um compromisso, e o conhecimento era compartilhado num caderninho de passwords grudento de guarana.

Este projeto e para quem ainda lembra do cheiro de plástico dos cartuchos, da emoção de encontrar a última cópia de Mega Man 2 na prateleira, e da comunidade que se formava em torno daquele balcao.

---

## A Locadora

**O Balcao** — A tela de entrada. Voce chega, mostra sua carteirinha e entra.

**A Prateleira** — O acervo de cartuchos, com capas pixeladas e status em tempo real. Se a fita estiver alugada, vai ter que esperar o proximo socio devolver.

**A Carteirinha de Socio** — Cada membro recebe um numero no formato `1991-XXX`. E digital, mas carrega o espirito daquele cartao plastificado com foto 3x4.

**O Caderno de Passwords** — Espaco pessoal para anotar senhas, codigos e mapas. Porque ninguem merece perder o progresso do Metroid.

**O Tio da Locadora** — O administrador. Ele abastece as prateleiras, cuida do acervo e da baixa nas devolucoes. Tudo pelo painel admin, como se estivesse atras do balcao.

---

## O Espirito da Coisa

Este nao e um sistema generico de catalogo. E uma experiencia. Cada detalhe foi pensado para evocar aquela epoca:

- **Escassez real** — Cada jogo tem copias limitadas. Se todas estiverem alugadas, o jogo fica indisponivel. Assim como era na locadora.
- **Visual 8-bit** — Interface com [NES.css](https://nostalgic-css.github.io/NES.css/) e fonte Press Start 2P. Cada pixel no lugar.
- **Sem JavaScript** — Renderizacao no servidor, como os sites de 1996. Rapido, limpo, sem frescura.
- **Copyleft** — Licenciado sob GPL v3. O codigo e livre, como deveria ser.

---

## Stack

| Componente | Tecnologia |
|------------|-----------|
| Backend | [Go](https://go.dev/) 1.24+ |
| Banco de dados | [PostgreSQL](https://www.postgresql.org/) 15+ |
| Interface | Server-Side Rendering com `html/template` |
| Estilo | [NES.css](https://nostalgic-css.github.io/NES.css/) 2.3.0 + Press Start 2P |
| Dados de jogos | [IGDB](https://api-docs.igdb.com/) via Twitch OAuth2 |
| Seguranca | bcrypt + HMAC-SHA256 + middleware de autorizacao |

---

## Comecar a Jogar

### Pre-requisitos

- Go 1.24+
- Docker (para o PostgreSQL)
- Credenciais da [API do IGDB](https://dev.twitch.tv/)

### Inicio rapido

```bash
git clone https://github.com/cmellojr/modo-locadora.git
cd modo-locadora
cp .env.example .env        # preencha com seus valores
docker compose up -d         # sobe o banco
go run ./cmd/server          # abre a locadora em http://localhost:8080
```

Para o guia completo de configuracao (migracoes, primeiro socio, solucao de problemas), veja **[docs/SETUP.md](docs/SETUP.md)**.

---

## Documentacao

| Documento | Conteudo |
|-----------|---------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | Visao geral da arquitetura, entidades e fluxos |
| [docs/SETUP.md](docs/SETUP.md) | Guia completo de instalacao e configuracao |
| [docs/API.md](docs/API.md) | Referencia de endpoints (SSR e JSON) |
| [docs/SECURITY.md](docs/SECURITY.md) | Politica de seguranca e praticas |
| [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) | Como contribuir com o projeto |
| [docs/CHANGELOG.md](docs/CHANGELOG.md) | Historico de mudancas |

---

## Funcionalidades Futuras

- **Verso da Capa** — Deixe dicas publicas para os proximos jogadores ao devolver uma fita.
- **Regra da Sexta** — Alugou na sexta? So precisa devolver na segunda!
- **Ranking de Socios** — Quem mais alugou, quem mais devolveu no prazo.
- **Colecao Pessoal** — Marque os jogos que voce ja zerou.

---

## Licenca

Distribuido sob a licenca **GPL v3**. Veja [LICENSE](LICENSE) para mais informacoes.

---

*Desenvolvido com nostalgia pelo Tio da Locadora.*

*Em memoria de todas as locadoras que fecharam, mas nunca foram esquecidas.*
