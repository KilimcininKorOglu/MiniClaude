import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

export default function ProjectStructure() {
  return (
    <MarkdownDoc title="Project Structure" description="MiniClaude source layout and architecture notes." content={<>
      <DocStyles />
      <pre><code>{`MiniClaude/
├── src/
│   ├── cli/              # CLI entrypoint and argument parsing
│   ├── commands/         # Slash command implementations
│   ├── components/       # Ink terminal UI components
│   ├── hooks/            # React hooks
│   ├── services/         # Core service layer
│   │   ├── api/          # AI API clients
│   │   ├── mcp/          # MCP protocol implementation
│   │   └── lsp/          # LSP language services
│   ├── tools/            # Agent tool implementations
│   ├── utils/            # Shared utility functions
│   │   ├── settings/     # Configuration system
│   │   └── permissions/  # Permission system
│   ├── skills/           # Skills system
│   │   └── bundled/      # Bundled skills
│   └── plugins/          # Plugin system
├── scripts/
│   └── build.ts          # Build script
├── website/              # Website source code
├── .env.example          # Environment template
└── README.md`}</code></pre>
    </>} />
  )
}
