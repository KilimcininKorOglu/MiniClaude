import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

export default function Mcp() {
  return (
    <MarkdownDoc title="MCP Protocol" description="Model Context Protocol extends the agent with external tools." content={<>
      <DocStyles />
      <h2 id="overview">Overview</h2>
      <p>MCP, short for Model Context Protocol, is an open protocol that lets MiniClaude connect to external tool servers and expand what the agent can do.</p>

      <h2 id="config">Server Configuration</h2>
      <p>Configure servers in <code>~/.claude/settings.json</code> or <code>.mcp.json</code>:</p>
      <pre><code>{`{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["-y", "chrome-devtools-mcp@latest"],
      "alwaysLoad": true
    },
    "my-server": {
      "command": "python",
      "args": ["path/to/server.py"]
    },
    "remote-server": {
      "type": "sse",
      "url": "https://example.com/mcp/sse"
    }
  }
}`}</code></pre>
      <p>When <code>"alwaysLoad": true</code> is enabled, every tool from that server is loaded into the agent tool list immediately instead of waiting for delayed discovery. This works well for frequently used servers.</p>

      <h2 id="transports">Transport Types</h2>
      <table>
        <thead><tr><th>Type</th><th>Description</th><th>Typical Use</th></tr></thead>
        <tbody>
          <tr><td><code>stdio</code></td><td>Standard input and output transport.</td><td>Local processes and the most common setup.</td></tr>
          <tr><td><code>sse</code></td><td>Server-Sent Events.</td><td>Remote HTTP servers.</td></tr>
          <tr><td><code>http</code></td><td>HTTP streaming transport.</td><td>Web service integration.</td></tr>
          <tr><td><code>ws</code></td><td>WebSocket transport.</td><td>Bidirectional real-time communication.</td></tr>
          <tr><td><code>sdk</code></td><td>VS Code SDK transport.</td><td>IDE extension integration.</td></tr>
        </tbody>
      </table>

      <h2 id="tools">Available Tools</h2>
      <p>After an MCP server is registered, its tools appear in the agent tool list with the format <code>mcp__serverName__toolName</code>. For example:</p>
      <ul>
        <li><code>mcp__chrome-devtools__navigate_page</code> — browser page navigation.</li>
        <li><code>mcp__chrome-devtools__take_screenshot</code> — web page screenshots.</li>
        <li><code>mcp__jadx-mcp__get_class_source</code> — APK decompiled source lookup.</li>
      </ul>

      <h2 id="management">Management Commands</h2>
      <pre><code>{`/mcp              # Show MCP server status
/mcp add          # Add a new server
/mcp remove       # Remove a server`}</code></pre>

      <blockquote>MCP tools are also subject to <code>hard_deny</code> rules. Configure <code>"mcp__*"</code> if you need to disable every MCP tool.</blockquote>
    </>} />
  )
}
