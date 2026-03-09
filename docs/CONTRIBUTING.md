# Contributing to Modo Locadora

Thanks for contributing! Every contribution matters, whether it's a bug report, feature suggestion, or pull request.

## Reporting Bugs

Check existing [Issues](https://github.com/cmellojr/modo-locadora/issues) first.

Include: clear title, steps to reproduce, expected vs. actual behavior, environment details (OS, Go version, browser), and logs or screenshots if applicable.

## Suggesting Features

Open an issue with the `enhancement` label. Describe the problem it solves, how you envision it working, and any alternatives you've considered.

## Pull Requests

1. Fork the repository and clone your fork.
2. Create a branch from `develop`: `git checkout -b feature/your-feature-name`
3. Follow the conventions below.
4. Verify with `go build ./...` and `go vet ./...`.
5. Commit with a descriptive message (see conventions below).
6. Push to your fork and open a PR against `develop`.
7. In the PR description, tell us **which game you were playing while coding** (it's tradition).

## Conventions

### Language

- **Code, routes, database, and documentation**: English.
- **UI templates** (`web/templates/`): Portuguese (BR).

### Go Style

- Follow the [Google Go Style Guide](https://google.github.io/styleguide/go/guide.html).
- Run `go vet ./...` before committing.
- Keep functions short and focused.

### Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/): `feat:`, `fix:`, `docs:`, `style:`, `refactor:`.

### Branching

- `main` — stable releases
- `develop` — active development (PR target)
- Feature branches: `feature/*`, `fix/*`, `hotfix/*`, `docs/*`

### CSS & Templates

- NES.css classes with dark theme overrides in `retro.css`.
- Shared utility classes go in `retro.css`. Page-specific styles go in the template's inline `<style>`.

### Database Migrations

- Place new migrations in `internal/database/migrations/`.
- Use incremental numbering: `006_description.sql`, `007_description.sql`.
- Document what each migration does in the file header.

## Development Setup

See [SETUP.md](SETUP.md) for local environment instructions.

## License

By contributing, you agree that your contributions will be licensed under the [GPL v3](../LICENSE).
