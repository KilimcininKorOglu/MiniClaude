<div align="center">

```
███╗   ███╗██╗███╗   ██╗██╗     ██████╗██╗      █████╗ ██╗   ██╗██████╗ ███████╗
████╗ ████║██║████╗  ██║██║    ██╔════╝██║     ██╔══██╗██║   ██║██╔══██╗██╔════╝
██╔████╔██║██║██╔██╗ ██║██║    ██║     ██║     ███████║██║   ██║██║  ██║█████╗  
██║╚██╔╝██║██║██║╚██╗██║██║    ██║     ██║     ██╔══██║██║   ██║██║  ██║██╔══╝  
██║ ╚═╝ ██║██║██║ ╚████║██║    ╚██████╗███████╗██║  ██║╚██████╔╝██████╔╝███████╗
╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚═╝     ╚═════╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝
```

### Lightweight Local AI Coding Assistant

[![Bun](https://img.shields.io/badge/Bun-1.3.11+-000000?style=for-the-badge&logo=bun&logoColor=white)](https://bun.sh)
[![TypeScript](https://img.shields.io/badge/TypeScript-6.0+-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![License](https://img.shields.io/badge/License-MIT-1ba784?style=for-the-badge)](LICENSE)

[Website](https://txl16095.github.io/MiniClaude/) · [Guide](https://txl16095.github.io/MiniClaude/#/guide/quick-start) · [Community](https://github.com/txl16095/MiniClaude/discussions)

</div>

---

## Overview

MiniClaude is a Bun and TypeScript command-line coding assistant derived from Claude Code. It keeps the local development workflow, file tools, Git integration, MCP support, plugins, skills, hooks, permissions, and multi-provider configuration while removing cloud services, telemetry, collaboration features, and other hosted-product dependencies.

The project is local-first by default. An optional self-hosted sync server is included for teams that want workspace login, Web UI management, device-code client linking, and settings synchronization under their own infrastructure.

---

## Screenshot

![MiniClaude Screenshot](assets/screenshot.png)

---

## Why MiniClaude?

- **Local-first runtime**: the CLI runs locally and does not require vendor-hosted telemetry or collaboration services.
- **Claude Code-style workflow**: interactive chat, code generation, file editing, shell commands, Git helpers, MCP, plugins, skills, hooks, and permissions remain available.
- **Multi-provider configuration**: switch between compatible Anthropic-style endpoints from `settings.json` without restarting the CLI.
- **Optional self-hosted sync**: link clients to the included Go + HTMX + PostgreSQL sync server when workspace-scoped settings sync is needed.
- **Simple build path**: install dependencies with Bun and build a local executable at `./cli`.

---

## Requirements

| Requirement | Version or note                  |
|-------------|----------------------------------|
| Bun         | `>= 1.3.11`                      |
| Git         | Required for repository workflow |
| Node/npm    | Optional, for some MCP examples  |
| PostgreSQL  | Optional, only for `sync-server` |

---

## Quick Start

Clone the repository, install dependencies, and build the CLI:

```bash
git clone https://github.com/txl16095/MiniClaude.git
cd MiniClaude
bun install
bun run build
```

Run the built CLI:

```bash
./cli
```

Development mode runs the TypeScript entrypoint directly:

```bash
bun run dev
```

The package exposes the built binary as both `miniclaude` and `mclaude` when installed or linked by a package manager.

---

## Configuration

MiniClaude uses Claude Code-style settings. The recommended global configuration file is:

```text
~/.claude/settings.json
```

Start from the provided template:

```bash
cp settings.example.json ~/.claude/settings.json
```

Windows users should place the file at:

```text
C:\Users\<username>\.claude\settings.json
```

A minimal provider configuration looks like this:

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "https://api.deepseek.com/anthropic",
    "ANTHROPIC_AUTH_TOKEN": "sk-your-api-key-here",
    "ANTHROPIC_MODEL": "deepseek-v4-pro[1m]"
  },
  "model": "deepseek-v4-pro[1m]"
}
```

See [`settings.example.json`](settings.example.json) for a complete template covering providers, default models, permissions, MCP servers, plugins, interface preferences, and optional sync notes.

For development or temporary overrides, a project-root `.env` file is also supported:

```bash
cp .env.example .env
bun run dev
```

---

## Multi-Provider Switching

Providers are configured under the `providers` key in `~/.claude/settings.json`. MiniClaude can switch providers globally or for the current session:

```bash
/provider
/provider deepseek
/provider kiro --session
```

| Command                      | Behavior                            |
|------------------------------|-------------------------------------|
| `/provider`                  | Lists configured providers          |
| `/provider <name>`           | Switches the global provider        |
| `/provider <name> --session` | Switches only the current CLI session |

Provider entries can define endpoint environment variables, the default model, model aliases for Sonnet/Haiku/Opus-style selection, and a short description.

---

## Permissions and Safety Controls

MiniClaude keeps the local permission model and hook system used by the CLI workflow.

| Setting or feature        | Purpose                                                             |
|---------------------------|---------------------------------------------------------------------|
| `permissions.defaultMode` | Controls ask, accept-edits, or bypass permission modes              |
| `autoMode.hard_deny`      | Blocks matching tool invocations unconditionally                    |
| Hooks                     | Runs user-configured scripts around lifecycle events                |
| MCP configuration         | Adds external tools and resources through MCP servers               |

`hard_deny` rules are intended for operations that should never be allowed, even when broader permission modes are active.

---

## Optional Sync Server

The `sync-server/` directory contains an optional self-hosted Go + HTMX + PostgreSQL service. It provides browser login, workspace membership, device-code MiniClaude linking, provider management, settings editing, client/session management, audit records, and WebSocket-based settings synchronization.

Run the service locally:

```bash
cd sync-server
cp .env.example .env
make migrate-up
make dev
```

Or start the local container stack:

```bash
cd sync-server
docker compose up --build
```

Link a MiniClaude client from inside the CLI:

```bash
/login http://localhost:8080
```

### Sync behavior

The sync client watches the real `userSettings` file, normally `~/.claude/settings.json`. It does not watch the repository's `settings.example.json` file.

Synchronization uses full settings snapshots:

- local `userSettings` edits are debounced and pushed to the server;
- WebSocket is used when available;
- HTTP `/api/settings/push` is used as a fallback when WebSocket is unavailable;
- stale pushes receive `version_reject` and the latest server snapshot is applied locally;
- conflict notices are written to stderr so the user can see that the server snapshot won;
- invalid or unreadable `userSettings` is not converted into an empty document for sync;
- remote snapshots overwrite the synced user settings document;
- project, local, policy, and flag settings are not synced by this client.

Because sync sends the full parsed `userSettings` document, do not store plaintext API tokens or other secrets in synced settings unless you are intentionally sending them to your own sync server. MiniClaude sync credentials are stored separately under the MiniClaude configuration directory.

### Sync server Web UI

The Web UI supports owner/admin mutations for provider, settings, client, session, and audit management. Auth uses `HttpOnly` browser cookies with CSRF protection. Provider secrets are encrypted at rest and are not written into settings documents, audit metadata, documentation examples, or browser storage by the intended flow.

See [`sync-server/README.md`](sync-server/README.md) for deployment, Coolify configuration, migration, rollback, API, WebSocket, and security details.

---

## Development Commands

| Command                                                      | Description                                                  |
|--------------------------------------------------------------|--------------------------------------------------------------|
| `bun install`                                                | Installs root dependencies                                   |
| `bun run dev`                                                | Runs `src/entrypoints/cli.tsx` directly                      |
| `bun run build`                                              | Bundles the CLI to `./cli`                                   |
| `bun run build:dev`                                          | Builds `./cli` with development defines and version metadata |
| `bun run compile`                                            | Creates a compiled Bun executable at `./dist/cli`            |
| `bun run scripts/test-p0-features.ts`                        | Runs the existing P0 smoke checks                            |
| `cd website && bun run build`                                | Builds the documentation website                             |
| `cd sync-server && go test ./...`                            | Runs sync-server tests                                       |
| `cd sync-server && go vet ./...`                             | Runs Go vet for sync-server                                  |
| `cd sync-server && go build -o bin/sync-server ./cmd/server` | Builds the sync-server binary                                |
| `cd sync-server && make migrate-up && make migrate-down`     | Verifies sync-server migrations                              |

The root `package.json` does not define package-level `test` or `lint` scripts. Run the smoke script or targeted scripts directly with `bun run`.

---

## Project Layout

| Path                         | Purpose                                                       |
|------------------------------|---------------------------------------------------------------|
| `src/entrypoints/cli.tsx`    | First runtime entrypoint and fast-path startup handling        |
| `src/main.tsx`               | Full CLI bootstrap for settings, tools, commands, and UI      |
| `src/query.ts`               | Agentic query loop and tool-use orchestration                 |
| `src/tools.ts`               | Built-in tool registry                                        |
| `src/commands.ts`            | Slash command registry                                        |
| `src/utils/settings/`        | Settings parsing, caching, schema, and source handling        |
| `src/services/sync/`         | MiniClaude sync client, credentials, API, and settings bridge |
| `sync-server/`               | Optional Go + HTMX + PostgreSQL sync server                   |
| `website/`                   | Documentation website                                         |
| `scripts/build.ts`           | Bun build and compile pipeline                                |
| `settings.example.json`      | Example global settings file                                  |

---

## Core Features

| Category        | Features                                                                     |
|-----------------|------------------------------------------------------------------------------|
| AI workflow     | Interactive chat, code generation, project analysis, and multi-model usage    |
| File tools      | Read, write, edit, search, globbing, and generated HTML reports               |
| Dev integration | Shell commands, Git workflow support, LSP integrations, and browser helpers   |
| Extensibility   | MCP servers, plugins, bundled skills, custom skills, and slash commands       |
| Configuration   | Claude Code-style settings, provider profiles, permissions, hooks, and themes |
| Sync            | Optional self-hosted workspace login, settings sync, client/session controls  |

---

## Difference from Claude Code

MiniClaude is derived from Claude Code but intentionally removes hosted-product behavior that is not needed for a local-first assistant.

| Area                  | Claude Code | MiniClaude                         |
|-----------------------|-------------|------------------------------------|
| Core coding tools     | Included    | Included                           |
| MCP/plugins/skills    | Included    | Included                           |
| Local settings        | Included    | Included                           |
| Cloud services        | Included    | Removed from the local-first core  |
| Telemetry             | Included    | Removed                            |
| Hosted collaboration  | Included    | Removed                            |
| Optional sync         | Hosted      | Self-hosted `sync-server/` package |

The included sync server is intentionally separate from the CLI core. If it is not configured and `/login` is not used, MiniClaude remains local-first.

---

## Documentation

| Resource                                                                  | Description                                  |
|---------------------------------------------------------------------------|----------------------------------------------|
| [Website guide](https://txl16095.github.io/MiniClaude/#/guide/quick-start) | User-facing guide and feature documentation  |
| [`settings.example.json`](settings.example.json)                           | Full settings template                       |
| [`sync-server/README.md`](sync-server/README.md)                           | Self-hosted sync server documentation        |
| [`CLAUDE.md`](CLAUDE.md)                                                   | Repository-specific development guidance     |

---

## Contributing

Use a topic branch for changes and keep commits focused:

```bash
git checkout -b feat/your-change
bun install
bun run build
bun run scripts/test-p0-features.ts
git commit -m 'feat: describe your change'
```

For sync-server changes, also run the relevant Go checks from `sync-server/` before opening a pull request.

---

## License

MIT License © 2026 [txl16095](https://github.com/txl16095)

Based on [free-code](https://github.com/paoloanzn/free-code). Original code copyright [Anthropic PBC](https://www.anthropic.com).

This is not an official Anthropic project. It is intended for learning, research, and local development workflows.

---

<div align="center">

[![Website](https://img.shields.io/badge/Website-1ba784?style=for-the-badge)](https://txl16095.github.io/MiniClaude/)
[![Claude Code](https://img.shields.io/badge/Claude_Code-orange?style=for-the-badge)](https://docs.anthropic.com/en/docs/claude-code)
[![Bun](https://img.shields.io/badge/Bun-000000?style=for-the-badge&logo=bun&logoColor=white)](https://bun.sh)

If this project helps you, please give it a star.

</div>
