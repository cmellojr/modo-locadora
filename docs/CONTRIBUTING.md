# Contributing to Modo Locadora

First off, thanks for taking the time to contribute! Every contribution matters, whether it's a bug report, a feature suggestion, or a pull request.

## Code of Conduct

This project follows a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold it.

## How Can I Contribute?

### Reporting Bugs

Before creating a bug report, please check the existing [Issues](https://github.com/cmellojr/modo-locadora/issues) to see if the problem has already been reported.

When filing a bug report, include:

- **A clear title** describing the problem.
- **Steps to reproduce** the issue.
- **Expected behavior** vs. **actual behavior**.
- **Environment details**: OS, Go version (`go version`), browser.
- **Logs or screenshots** if applicable.

### Suggesting Features

Feature requests are welcome. Open an issue with the `enhancement` label and describe:

- **What problem** the feature solves.
- **How you envision it** working.
- **Any alternatives** you've considered.

### Submitting Pull Requests

1. **Fork** the repository and clone your fork.
2. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Follow the conventions** listed below.
4. **Test your changes** — make sure `go build ./...` and `go vet ./...` pass.
5. **Commit** with a descriptive message (see commit conventions below).
6. **Push** to your fork and open a Pull Request against `main`.
7. In the PR description, tell us **which game you were playing while coding** (it's tradition).

## Conventions

### Language

- **Code, routes, database schemas, and documentation**: English.
- **UI templates** (`web/templates/`): Portuguese (BR).
- This follows the rule defined in [ARCHITECTURE.md](../ARCHITECTURE.md).

### Go Style

- Follow the [Google Go Style Guide](https://google.github.io/styleguide/go/guide.html).
- Run `go vet ./...` before committing.
- Keep functions short and focused.
- Use meaningful names — avoid abbreviations except for well-known ones (`ctx`, `err`, `req`).

### Project Structure

```
cmd/server/         -> Application entrypoint
internal/
  auth/             -> Authentication utilities (cookie signing)
  config/           -> Environment configuration loader
  database/         -> Store interface and PostgreSQL implementation
    migrations/     -> SQL migration files (001-003)
  handlers/         -> HTTP request handlers
  igdb/             -> IGDB API client
  middleware/        -> HTTP middleware (auth, admin)
  models/           -> Domain entities (Member, Game, GameCopy, Rental)
web/
  static/css/       -> Stylesheets
  templates/        -> Go HTML templates (PT-BR)
    index.html          Login page (Balcao)
    games.html          Game shelf with rental status
    carteirinha.html    Membership card
    admin_stock.html    IGDB search & game purchase
    admin_inventory.html Catalog listing with edit buttons
    admin_edit.html     Game edit form
    admin_returns.html  Active rentals check-in
```

### Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add game rental flow
fix: correct password validation on login
docs: update API reference
style: adjust shelf grid spacing
refactor: extract cookie signing to auth package
```

### Database Migrations

- Place new migrations in `internal/database/migrations/`.
- Name them with an incremental prefix: `004_description.sql`, `005_description.sql`.
- Migrations are applied manually — document what each one does.
- Current migrations: `001` (initial schema), `002` (games table update), `003` (membership and rental support).

## Development Setup

See [SETUP.md](SETUP.md) for instructions on setting up your local development environment.

## License

By contributing, you agree that your contributions will be licensed under the [GPL v3](../LICENSE) license.
