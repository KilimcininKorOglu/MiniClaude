# MiniClaude Sync Server

Go + HTMX + PostgreSQL service for centralized MiniClaude login, workspace-scoped settings sync, provider management, and client session control.

## Local development

```bash
cp .env.example .env
make migrate-up
make dev
```

For a full local stack:

```bash
docker compose up --build
```

## Commands

```bash
make build
make test
make vet
make migrate-up
make migrate-down
```

## Configuration

| Variable | Description |
| -------- | ----------- |
| `SYNC_SERVER_ADDR` | HTTP listen address. |
| `SYNC_SERVER_APP_URL` | Public URL used for browser device-code links. |
| `SYNC_SERVER_COOKIE_SECURE` | Enables Secure cookies in HTTPS environments. |
| `SYNC_SERVER_JWT_SECRET` | Secret used to sign web session JWTs; must be at least 32 bytes. |
| `SYNC_SERVER_PROVIDER_SECRET_KEY` | Master key used to encrypt provider secrets; must be at least 32 bytes. |
| `SYNC_SERVER_SESSION_TTL` | Web session lifetime, such as `24h`. |
| `SYNC_SERVER_SHUTDOWN_TIMEOUT` | Graceful shutdown timeout. |
| `SYNC_SERVER_MIGRATIONS_DIR` | Directory containing SQL migration files. |
| `DATABASE_URL` | PostgreSQL connection URL. |

## Auth API

The initial JSON auth API exposes:

- `POST /auth/signup` with `email`, `username`, `password`, and `workspace_name`.
- `POST /auth/login` with `login` and `password`.
- `POST /auth/logout` to revoke the current web session.
- `GET /me` to read the authenticated user and workspace context.

Successful signup and login responses set the JWT in the `miniclaude_sync_session` cookie with `HttpOnly` and `SameSite=Lax`. The token is not returned in the JSON body.

## Migration policy

Migration files use paired `*.up.sql` and `*.down.sql` files. Every schema change must include a rollback file and should run inside the migration transaction.

## Coolify

Use `Dockerfile` as the build source and configure environment variables from `coolify.example.env`. PostgreSQL can be attached as a Coolify database resource or provided through `DATABASE_URL`.
