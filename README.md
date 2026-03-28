# Modo Locadora

> **"Sopre a fita, pegue o controle e respeite o tempo."**

![Go](https://img.shields.io/badge/go-1.24-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgresql-15-336791?logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/docker-compose-2496ED?logo=docker&logoColor=white)
![NES.css](https://img.shields.io/badge/nes.css-2.3.0-E76F51)
![Retro](https://img.shields.io/badge/visual-8--bit-E76F51)
![License](https://img.shields.io/badge/license-GPL%20v3-blue)
![SSR](https://img.shields.io/badge/rendering-SSR-green)
![JS](https://img.shields.io/badge/javascript-zero-red)

---

Era sexta-feira. Você corria da escola, passava no balcão, mostrava a carteirinha e rezava para aquele cartucho ainda estar na prateleira. O tio conferia o caderno, pegava a fita, soprava o conector e dizia: *"Devolve na segunda, hein?"*

O **Modo Locadora** é um diário de bordo de jogatinas disfarçado de videolocadora. Na prática, funciona como um rastreador de backlog pessoal — você aluga jogos, registra se zerou ou desistiu, acumula títulos de reputação e acompanha seu progresso numa carteirinha de sócio — tudo com a estética e as regras de uma locadora brasileira dos anos 90, onde a escassez era real e o tempo era um compromisso.

O projeto é uma homenagem aos antigos **Projeto Jogatina** e **Fórum NES Archive** — iniciativas da comunidade retrogamer brasileira que combatiam a *Síndrome do Labirinto* (ter jogos demais e não terminar nenhum) através do registro e compartilhamento de jogatinas. Inspirado também pelo **[Backloggery](https://backloggery.com/)**, o Modo Locadora adiciona mecânicas de escassez e reputação social para transformar o simples ato de jogar em uma experiência coletiva e com consequências.

---

## A Locadora

**O Balcão** — A tela de entrada. Você chega, mostra sua carteirinha e entra. Sócios em débito aparecem no Painel da Vergonha.

**A Prateleira** — Organizada por console (Mega Drive, SNES, NES...). Escolha a plataforma, navegue pelos cartuchos e veja os detalhes de cada fita — resumo, revista de origem, quantas vezes foi alugada e quem é o fã número 1.

**A Carteirinha de Sócio** — Cada membro recebe um número no formato `1991-XXX`. É digital, mas carrega o espírito daquele cartão plastificado com foto 3x4. Conforme você aluga e devolve, conquista títulos: Sócio Novato, Sócio Prata, Sócio Ouro, até Dono da Calçada.

**O Caderno de Passwords** — Espaço pessoal para anotar senhas, códigos e mapas. Porque ninguém merece perder o progresso do Metroid.

**O Tio da Locadora** — O administrador. Abastece as prateleiras com capas brasileiras (TecToy, Playtronic), cuida do acervo e dá baixa nas devoluções. No inventário, cada fita tem um indicador de saúde (Cartucho Novo, Clássico Eterno, Precisa Soprar, Fita Gasta) baseado no histórico de vereditos.

**O Veredito** — Ao devolver uma fita, diga ao Tio se você zerou, jogou um pouco ou desistiu. Quem zerou ganha uma estrela dourada na prateleira.

**Aconteceu na Locadora** — Feed de atividades em tempo real. Quem alugou, quem zerou, quem foi pro Painel da Vergonha — tudo aparece no balcão.

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
docker exec modo_locadora_app /app/server --seed  # popula com dados de teste (opcional)
```

Acesse `http://localhost:8080` — a locadora está aberta.
Com seed: `MegaDriveKid` / `sega1991`, `Devedor` / `atrasado123`, `Novato` / `novato2026`.

Para desenvolvimento local sem Docker, migrações manuais e criação do primeiro sócio, veja **[docs/SETUP.md](docs/SETUP.md)**.

---

## Documentação

| Documento | Conteúdo |
|-----------|---------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | Arquitetura do sistema |
| [ROADMAP.md](ROADMAP.md) | Plano de evolução e versões |
| [docs/SETUP.md](docs/SETUP.md) | Configuração do ambiente |
| [docs/API.md](docs/API.md) | Referência da API |
| [docs/SECURITY.md](docs/SECURITY.md) | Política de segurança |
| [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) | Como contribuir |
| [docs/CHANGELOG.md](docs/CHANGELOG.md) | Histórico de mudanças |

---

## Automação

O projeto usa [Task](https://taskfile.dev/) para comandos comuns:

```bash
task check     # build + vet + lint
task seed      # aplica migrações + dados de teste
task reset     # reset completo (down + up + seed)
task logs      # logs do container
```

## Próximos Passos

Confira o [ROADMAP.md](ROADMAP.md) para o plano completo de evolução — versões lançadas, em andamento e futuras.

---

## Contribuindo

Pull requests são bem-vindos. Leia o [guia de contribuição](docs/CONTRIBUTING.md) antes de abrir um PR. Issues e discussões em português, por favor.

---

## Licença

Distribuído sob a licença **GPL v3**. Veja [LICENSE](LICENSE) para mais informações.

---

*Desenvolvido com nostalgia pelo Tio da Locadora.*

*Em memória de todas as locadoras que fecharam, mas nunca foram esquecidas.*
