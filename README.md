# Modo Locadora

> **"Sopre a fita, pegue o controle e respeite o tempo."**

---

Era sexta-feira. Você corria da escola, passava no balcão, mostrava a carteirinha e rezava para aquele cartucho ainda estar na prateleira. O tio conferia o caderno, pegava a fita, soprava o conector e dizia: *"Devolve na segunda, hein?"*

O **Modo Locadora** é uma homenagem a essa era dourada. Um sistema web que recria a experiência das videolocadoras brasileiras dos anos 90 — onde a escassez gerava valor, o tempo era um compromisso, e o conhecimento era compartilhado num caderninho de passwords grudento de guaraná.

---

## A Locadora

**O Balcão** — A tela de entrada. Você chega, mostra sua carteirinha e entra. Sócios em débito aparecem no Painel da Vergonha.

**A Prateleira** — Organizada por console (Mega Drive, SNES, NES...). Escolha a plataforma, navegue pelos cartuchos e veja os detalhes de cada fita — resumo, revista de origem, quantas vezes foi alugada e quem é o fã número 1.

**A Carteirinha de Sócio** — Cada membro recebe um número no formato `1991-XXX`. É digital, mas carrega o espírito daquele cartão plastificado com foto 3x4.

**O Caderno de Passwords** — Espaço pessoal para anotar senhas, códigos e mapas. Porque ninguém merece perder o progresso do Metroid.

**O Tio da Locadora** — O administrador. Abastece as prateleiras com capas brasileiras (TecToy, Playtronic), cuida do acervo e dá baixa nas devoluções.

**O Fiscal Automático** — Não devolveu no prazo? O sistema devolve a fita automaticamente e marca seu nome no Painel da Vergonha.

---

## O Espírito da Coisa

- **Escassez real** — Cada jogo tem cópias limitadas. Se todas estiverem alugadas, o jogo fica indisponível.
- **Visual 8-bit** — Interface com [NES.css](https://nostalgic-css.github.io/NES.css/) e fonte Press Start 2P. Cada pixel no lugar.
- **Sem JavaScript** — Renderização no servidor. Rápido, limpo, sem frescura.
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
| Deploy | Docker Compose (app + banco) |

---

## Começar a Jogar

### Pré-requisitos

- Docker e Docker Compose
- Credenciais da [API do IGDB](https://dev.twitch.tv/) (Twitch Developer)

### Início rápido

```bash
git clone https://github.com/cmellojr/modo-locadora.git
cd modo-locadora
cp .env.example .env        # preencha com seus valores
docker compose up -d --build # sobe tudo: app + banco
```

Acesse `http://localhost:8080` — a locadora está aberta.

Para desenvolvimento local sem Docker, migrações manuais e criação do primeiro sócio, veja **[docs/SETUP.md](docs/SETUP.md)**.

---

## Documentação

| Documento | Conteúdo |
|-----------|---------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | Entidades, fluxos e navegação |
| [docs/SETUP.md](docs/SETUP.md) | Guia completo de instalação |
| [docs/API.md](docs/API.md) | Referência de endpoints |
| [docs/SECURITY.md](docs/SECURITY.md) | Política de segurança |
| [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) | Como contribuir |
| [docs/CHANGELOG.md](docs/CHANGELOG.md) | Histórico de mudanças |

---

## Funcionalidades Futuras

- **Verso da Capa** — Deixe dicas públicas para os próximos jogadores ao devolver uma fita.
- **Regra da Sexta** — Alugou na sexta? Só precisa devolver na segunda!
- **Coleção Pessoal** — Marque os jogos que você já zerou.
- **Menções na Mídia** — Registre em quais podcasts, sites ou reportagens cada jogo foi mencionado.

---

## Licença

Distribuído sob a licença **GPL v3**. Veja [LICENSE](LICENSE) para mais informações.

---

*Desenvolvido com nostalgia pelo Tio da Locadora.*

*Em memória de todas as locadoras que fecharam, mas nunca foram esquecidas.*
