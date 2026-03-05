# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| main    | :white_check_mark: |

Only the latest version on `main` receives security updates.

## Reporting a Vulnerability

**Do not open a public issue for security vulnerabilities.**

Report privately via email to the project maintainer or use [GitHub's private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing-information-about-vulnerabilities/privately-reporting-a-security-vulnerability).

Include: description, steps to reproduce, potential impact, and any suggested fix.

## Authentication

- Passwords are hashed with **bcrypt** (default cost) before storage. Never logged or exposed in API responses.
- Login requires profile name + password, validated against the database.

## Session Management

Sessions use a signed cookie (`session_member`):

| Property | Value |
|----------|-------|
| Format | `{member_uuid}.{hmac_sha256_hex}` |
| HttpOnly | `true` |
| SameSite | `Strict` |
| MaxAge | 604800 (7 days) |
| Path | `/` |

The cookie is signed using HMAC-SHA256 with the `COOKIE_SECRET` environment variable. Forged or tampered cookies are rejected.

## Authorization

| Scope | Middleware | Check |
|-------|-----------|-------|
| Member routes (`/carteirinha`, `/rent`) | `RequireAuth` | Valid signed cookie |
| Admin routes (`/admin/*`) | `RequireAdmin` | Valid cookie + email matches `ADMIN_EMAIL` |

Unauthenticated requests redirect to `/`. Non-admin users receive `403 Forbidden`.

## Data Integrity

- Rental and return operations use **database transactions** for atomicity (find copy, update status, create/update rental record).
- Game purchases atomically create both the game and a copy in a single transaction.

## Environment & Secrets

- Credentials loaded from `.env` (never committed — listed in `.gitignore`).
- `.env.example` provided with placeholder values.

## Deployment Checklist

- Set a strong, random `COOKIE_SECRET` (minimum 32 characters).
- Set `ADMIN_EMAIL` to restrict admin access.
- Use **HTTPS** in production to protect cookies and form data.
- Restrict database access to the application server only.
- Rotate Twitch API credentials if compromised.
