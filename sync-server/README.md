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

The local Docker stack exposes PostgreSQL on `localhost:5432` and the server on `http://localhost:8080`.

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

## Web UI

The server-rendered Web UI exposes:

- `GET /login` and `GET /signup` for browser auth.
- `GET /dashboard` for the authenticated workspace landing page.
- `GET /providers` and `GET /clients` as authenticated management pages.
- `GET /sessions` for workspace-wide MiniClaude session review and termination.
- `GET /audit` for recent workspace audit events.
- `GET /device` and `POST /device` for browser approval of MiniClaude device-code login requests.

Form submissions use same-origin POST routes, CSRF tokens, and security headers. The JWT is stored only in the `miniclaude_sync_session` cookie with `HttpOnly` and `SameSite=Lax`; set `SYNC_SERVER_COOKIE_SECURE=true` in HTTPS environments.

## MiniClaude client login

Start the server, sign in through the Web UI, then link a local MiniClaude client:

```text
/login http://localhost:8080
```

The command starts a device-code request, opens the browser approval page, polls for approval, and stores the returned client token under the MiniClaude config directory. The approved device code is consumed on the first successful poll and cannot be exchanged for another token. Subsequent MiniClaude startup connects to `/api/sync/ws`, sends `hello`, applies workspace snapshots to `userSettings`, pushes local settings changes as `settings_push`, handles `settings_applied` and `version_reject`, and accepts targeted `terminate_session` messages.

## Sync API

The sync API exposes:

- `POST /api/device/start` to create a device-code login request.
- `POST /api/device/approve` to approve a user code from an authenticated web session.
- `POST /api/device/poll` to exchange an approved device code for a client access token. The first successful exchange consumes the request.
- `GET /api/settings/snapshot` to read the workspace settings document.
- `POST /api/settings/push` to update settings with base-version checking; stale writes return `409` with `version_reject`.
- `GET /api/sync/ws` for client WebSocket sync using a bearer client token. The WebSocket protocol supports `hello`, `snapshot`, `settings_push`, `settings_applied`, `settings_updated`, `version_reject`, `ping`, `pong`, and targeted `terminate_session`.

## Auth API

The initial JSON auth API exposes:

- `POST /auth/signup` with `email`, `username`, `password`, and `workspace_name`.
- `POST /auth/login` with `login` and `password`.
- `POST /auth/logout` to revoke the current web session.
- `GET /me` to read the authenticated user and workspace context.

Successful signup and login responses set the JWT in the `miniclaude_sync_session` cookie with `HttpOnly` and `SameSite=Lax`. The token is not returned in the JSON body.

## Verification checklist

Before shipping sync-server changes, verify:

- Signup, login, logout, and dashboard access through the Web UI.
- CSRF rejection for Web UI POST routes when the hidden token is missing or invalid.
- Provider create, update, delete, and encrypted secret storage.
- Settings save success and stale-version conflict handling.
- Audit and sessions pages at `GET /audit` and `GET /sessions`.
- Device login start, browser approval, first poll token exchange, and second poll no-token behavior.
- WebSocket `hello`, initial snapshot, `settings_push`, `settings_applied`, `settings_updated`, and `version_reject` behavior with two MiniClaude clients.
- Targeted session termination from the Web UI.
- Browser storage contains no auth token in localStorage or sessionStorage.
- `SYNC_SERVER_COOKIE_SECURE=true` is enabled for HTTPS deployments.
- Docker Compose and Coolify health checks pass on `/healthz`.

## Migration policy

Migration files use paired `*.up.sql` and `*.down.sql` files. Every schema change must include a rollback file and should run inside the migration transaction.

## Coolify

Use `Dockerfile` as the build source and configure environment variables from `coolify.example.env`. PostgreSQL can be attached as a Coolify database resource or provided through `DATABASE_URL`.

Recommended Coolify settings:

- Build pack: Dockerfile.
- Exposed port: `8080`.
- Health check path: `/healthz`.
- Set `SYNC_SERVER_APP_URL` to the public HTTPS URL.
- Set `SYNC_SERVER_COOKIE_SECURE=true` for HTTPS deployments.
- Generate unique production values for `SYNC_SERVER_JWT_SECRET` and `SYNC_SERVER_PROVIDER_SECRET_KEY`; both must be at least 32 bytes.

Run migrations before serving a fresh database:

```bash
sync-server migrate up
```

Rollback is available through:

```bash
sync-server migrate down
```
