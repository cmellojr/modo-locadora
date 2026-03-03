# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| main    | :white_check_mark: |

Only the latest version on the `main` branch receives security updates.

## Reporting a Vulnerability

**Please do not open a public issue for security vulnerabilities.**

Instead, report them privately by emailing the project maintainer or using [GitHub's private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing-information-about-vulnerabilities/privately-reporting-a-security-vulnerability) feature on this repository.

When reporting, please include:

- A description of the vulnerability.
- Steps to reproduce or a proof of concept.
- The potential impact.
- Any suggested fix (optional but appreciated).

You should receive an acknowledgment within **48 hours**. We will work with you to understand the issue and coordinate a fix before any public disclosure.

## Security Measures

This project implements the following security practices:

### Authentication
- Passwords are hashed using **bcrypt** with the default cost factor before storage.
- Login requires both profile name and password, validated against the database.
- Passwords are never logged or exposed in API responses.

### Session Management
- Session cookies are **HMAC-SHA256 signed** using a server-side secret (`COOKIE_SECRET`).
- Cookies are set with `HttpOnly`, `SameSite=Strict`, and a `MaxAge` of 7 days.
- Forged or tampered cookies are automatically rejected.

### Authorization
- Admin routes (`/admin/*`) are protected by middleware that verifies both authentication and admin role (checked against the `ADMIN_EMAIL` environment variable).
- Unauthenticated requests to admin routes are redirected to the login page.

### Environment & Secrets
- Sensitive credentials are loaded from environment variables (`.env` file).
- The `.env` file is listed in `.gitignore` and is never committed to version control.
- A `.env.example` file with placeholder values is provided for reference.

## Best Practices for Deployment

- **Always set a strong, random `COOKIE_SECRET`** in production (at least 32 characters).
- **Rotate the Twitch API credentials** if you suspect they have been compromised.
- **Use HTTPS** in production to protect cookies and form data in transit.
- **Restrict database access** to the application server only (do not expose PostgreSQL publicly).
