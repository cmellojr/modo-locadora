# Modo Locadora

> **"Sopre a fita, pegue o controle e respeite o tempo."**

![Go](https://img.shields.io/badge/go-1.24-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgresql-15-336791?logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/docker-compose-2496ED?logo=docker&logoColor=white)
![NES.css](https://img.shields.io/badge/nes.css-2.3.0-E76F51)
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

**As Turmas** — Crie ou entre numa turma — representando seu podcast favorito, canal do YouTube, grupo de WhatsApp ou qualquer comunidade gamer. Cada turma tem badge, descrição e URL. Múltiplos admins, participação livre em quantas turmas quiser. A carteirinha mostra suas turmas com cargo.

**O Fiscal Automático** — Não devolveu no prazo? O sistema devolve a fita automaticamente e marca seu nome no Painel da Vergonha.

---

## O Espírito da Coisa

- **Escassez real** — Cada jogo tem cópias limitadas. Se todas estiverem alugadas, o jogo fica indisponível.
- **Visual 8-bit** — Interface com [NES.css](https://nostalgic-css.github.io/NES.css/) e fonte Press Start 2P. Cada pixel no lugar.
- **Sem JavaScript** — Renderização no servidor. Rápido, limpo, sem frescura.
- **Copyleft** — Licenciado sob GPL v3. O código é livre, como deveria ser.

---

## Stack

| Camada | Tecnologia |
|--------|------------|
| Backend | Go 1.24+ com `net/http.ServeMux` |
| Banco de dados | PostgreSQL 15 + `pgx/v5` |
| Interface | SSR com `html/template` |
| Estilo | NES.css 2.3.0 + Press Start 2P |
| Dados de jogos | IGDB via Twitch OAuth2 |
| Deploy local | Docker Compose |
| Uploads | `internal/storage` (`LocalStorage` ou `GCSStorage`) |

## Início Rápido

```bash
git clone https://github.com/cmellojr/modo-locadora.git
cd modo-locadora
cp .env.example .env
docker compose up -d --build
docker exec modo_locadora_app /app/server --seed
```

Acesse `http://localhost:8080`.

Contas de seed:

| Perfil | Senha |
|--------|-------|
| `MegaDriveKid` | `sega1991` |
| `Devedor` | `atrasado123` |
| `Novato` | `novato2026` |
| `tio_da_locadora` | `sopre_a_fita` |

Para configuração completa, variáveis de ambiente, migrations e troubleshooting, veja [docs/setup.md](docs/setup.md).

## Documentação

| Documento | Responsabilidade |
|-----------|------------------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | Arquitetura, módulos, fluxo de domínio e templates |
| [ROADMAP.md](ROADMAP.md) | Versões entregues e próximas ideias |
| [AGENTS.md](AGENTS.md) | Instruções para agentes de IA trabalhando no repositório |
| [docs/setup.md](docs/setup.md) | Ambiente local, Docker, migrations, seed e storage |
| [docs/api.md](docs/api.md) | Rotas SSR, endpoints de formulário e JSON |
| [docs/security.md](docs/security.md) | Autenticação, autorização, uploads e checklist de deploy |
| [docs/contributing.md](docs/contributing.md) | Fluxo de contribuição e convenções |
| [docs/prd.md](docs/prd.md) | Contexto de produto e requisitos |
| [docs/changelog.md](docs/changelog.md) | Histórico de mudanças notáveis |

## Automação

O projeto usa [Task](https://taskfile.dev/) para comandos comuns:

```bash
task build
task vet
task lint
task check
task dev
task seed
task up
task down
task reset
task logs
task psql
```

Veja [docs/contributing.md](docs/contributing.md) para o fluxo de branches e validação antes de PR.

## Licença

Distribuído sob a licença **GPL v3**. Veja [LICENSE](LICENSE).

*Desenvolvido com nostalgia pelo Tio da Locadora.*

*Em memória de todas as locadoras que fecharam, mas nunca foram esquecidas.*
