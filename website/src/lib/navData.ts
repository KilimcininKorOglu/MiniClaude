export interface NavItem {
  title: string
  path: string
}

export interface NavSection {
  title: string
  items: NavItem[]
}

export const sidebarNav: NavSection[] = [
  {
    title: 'Quick Start',
    items: [
      { title: 'Install and Run', path: '/guide/quick-start' },
      { title: 'Environment Variables', path: '/guide/env-vars' },
      { title: 'Third-Party Models', path: '/guide/third-party-models' },
      { title: 'Global Usage', path: '/guide/global-usage' },
      { title: 'FAQ', path: '/guide/faq' },
    ],
  },
  {
    title: 'Features',
    items: [
      { title: 'Command Reference', path: '/features/commands' },
      { title: 'Tool Reference', path: '/features/tools' },
      { title: 'hard_deny Security Rules', path: '/features/hard-deny' },
      { title: 'HTML Output', path: '/features/html-output' },
      { title: 'Skills System', path: '/features/skills' },
      { title: 'MCP Protocol', path: '/features/mcp' },
    ],
  },
  {
    title: 'Reference',
    items: [
      { title: 'Project Structure', path: '/reference/project-structure' },
      { title: 'Differences from Upstream', path: '/reference/differences' },
    ],
  },
]

// Build-time: flatten all pages into searchable entries.
export interface SearchEntry {
  title: string
  path: string
  section: string
  content: string
}

export const searchIndex: SearchEntry[] = [
  {
    title: 'Install and Run',
    path: '/guide/quick-start',
    section: 'Quick Start',
    content:
      'Install Bun on macOS Linux or Windows, run bun install, configure .env.example with your API key, start with ./cli or bun, configure PATH for global usage, and use recovery mode with CLAUDE_CODE_FORCE_RECOVERY_CLI.',
  },
  {
    title: 'Environment Variables',
    path: '/guide/env-vars',
    section: 'Quick Start',
    content:
      'ANTHROPIC_API_KEY is required, ANTHROPIC_BASE_URL sets a custom endpoint, ANTHROPIC_MODEL chooses the default model, ENABLE_PROMPT_CACHING_1H enables one-hour prompt caching, DISABLE_PROMPT_CACHING disables caching, CLAUDE_DISABLE_STREAM_WATCHDOG controls stream idle watchdog behavior, and HTTP proxy variables configure networking.',
  },
  {
    title: 'Third-Party Models',
    path: '/guide/third-party-models',
    section: 'Quick Start',
    content:
      'Use DeepSeek, OpenAI-compatible APIs, OpenRouter, SiliconFlow, or other Anthropic-compatible providers by setting ANTHROPIC_BASE_URL, ANTHROPIC_API_KEY, and ANTHROPIC_MODEL for multi-model support.',
  },
  {
    title: 'Global Usage',
    path: '/guide/global-usage',
    section: 'Quick Start',
    content:
      'Add the bin directory to PATH, create a symlink or shell export, start MiniClaude from any directory, and verify the global command with version and help checks.',
  },
  {
    title: 'FAQ',
    path: '/guide/faq',
    section: 'Quick Start',
    content:
      'Troubleshooting for startup failures, permission prompts, slow model responses, build errors, settings that do not apply, MCP connection issues, and plugin installation problems.',
  },
  {
    title: 'Command Reference',
    path: '/features/commands',
    section: 'Features',
    content:
      'Slash commands including /help, /clear, /exit, /config, /model, /provider, /theme, /files, /add-dir, /mcp, /skills, /tasks, /hooks, /permissions, /vim, /compact, /review, /stats, /status, /fast, /effort, /copy, /doctor, /diff, /init, /pr_comments, /plan, /export, and /rename.',
  },
  {
    title: 'Tool Reference',
    path: '/features/tools',
    section: 'Features',
    content:
      'Built-in tools including FileRead, FileWrite, FileEdit, Glob, Grep, Bash, PowerShell, WebFetch, WebSearch, AgentTool, SkillTool, MCP tools, TaskCreate, TaskStop, NotebookEdit, and AskUserQuestion for file work, search, command execution, networking, and automation.',
  },
  {
    title: 'hard_deny Security Rules',
    path: '/features/hard-deny',
    section: 'Features',
    content:
      'hard_deny blocks matched tool calls unconditionally, enforces safety at the pipeline level, protects sensitive files, prevents destructive shell commands, and works alongside settings.json permissions and auto mode rules.',
  },
  {
    title: 'HTML Output',
    path: '/features/html-output',
    section: 'Features',
    content:
      'Generate self-contained HTML reports with dark mode, responsive layout, system fonts, table styling, and browser-friendly output for code reviews, architecture notes, and analysis summaries.',
  },
  {
    title: 'Skills System',
    path: '/features/skills',
    section: 'Features',
    content:
      'Skills are reusable workflows powered by SKILL.md files, bundled registrations, automatic prompt loading, and slash-command entry points such as simplify, html-output, debug, batch, stuck, and verify.',
  },
  {
    title: 'MCP Protocol',
    path: '/features/mcp',
    section: 'Features',
    content:
      'MCP stands for Model Context Protocol and supports stdio, SSE, HTTP, WebSocket, alwaysLoad preloading, tool discovery, resource access, and external servers like chrome-devtools or jadx-mcp.',
  },
  {
    title: 'Project Structure',
    path: '/reference/project-structure',
    section: 'Reference',
    content:
      'Source layout covering CLI entrypoints, commands, tools, components, services, utilities, skills, plugins, scripts, and the website source tree.',
  },
  {
    title: 'Differences from Upstream',
    path: '/reference/differences',
    section: 'Reference',
    content:
      'MiniClaude removes cloud services, OAuth, telemetry, settings sync, collaboration features, experiments, and rate limiting while keeping core coding workflows and adding hard_deny plus HTML output.',
  },
]
