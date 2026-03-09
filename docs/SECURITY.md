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

## Authentication & Sessions

- Passwords hashed with **bcrypt** (default cost). Never logged or exposed in API responses.
- Sessions use a signed cookie (`session_member`): `{member_uuid}.{hmac_sha256_hex}`
- Cookie flags: `HttpOnly`, `SameSite=Strict`, `MaxAge=604800` (7 days), `Path=/`
- `COOKIE_SECRET` must be at least 32 characters.

## Authorization

| Scope | Middleware | Check |
|-------|-----------|-------|
| Member routes | `RequireAuth` | Valid signed cookie |
| Admin routes (`/admin/*`) | `RequireAdmin` | Valid cookie + email matches `ADMIN_EMAIL` |

Unauthenticated requests redirect to `/`. Non-admin users receive `403 Forbidden`.

## Member Reputation

- Overdue rentals are auto-returned by a background job (5-minute interval).
- Offending members are marked `em_debito` with a permanent `late_count` increment.
- Members in debt cannot rent games until they redeem themselves via `/carteirinha/redeem`.
- The Wall of Shame on the landing page displays top offenders.

## Data Integrity

- Rental, return, and game acquisition operations use **database transactions**.
- All SQL queries use parameterized placeholders (no string interpolation).

## File Uploads

- Cover uploads are restricted to image files (`accept="image/*"`).
- Maximum form size: 10 MB.
- Files are saved with the game UUID as filename (prevents path traversal).

## Deployment Checklist

- Set a strong, random `COOKIE_SECRET` (minimum 32 characters).
- Set `ADMIN_EMAIL` to restrict admin access.
- Use **HTTPS** in production to protect cookies and form data.
- Restrict database access to the application server only.
- Rotate Twitch API credentials if compromised.
- Credentials loaded from `.env` (never committed — listed in `.gitignore`).
