import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

export default function Differences() {
  return (
    <MarkdownDoc title="Differences from Upstream" description="How MiniClaude trims and extends the original Claude Code experience." content={<>
      <DocStyles />
      <h2 id="removed">Removed</h2>
      <table>
        <thead><tr><th>Module</th><th>Approx. Size</th><th>Description</th></tr></thead>
        <tbody>
          <tr><td>Cloud service integration</td><td>~7,173 lines</td><td>OAuth, telemetry, settings sync, and policy constraints.</td></tr>
          <tr><td>Collaboration features</td><td>~24,387 lines</td><td>Team collaboration, bridge mode, remote control, and multi-agent coordination.</td></tr>
          <tr><td>Experimental features</td><td>~1,950 lines</td><td>Voice mode, desktop integration, mobile support, and Buddy features.</td></tr>
          <tr><td>Complex integrations</td><td>~3,170 lines</td><td>Teleport, auto-updater, and Slack integration.</td></tr>
          <tr><td>Rate limiting system</td><td>~48,000 lines</td><td>Rate-limit simulation, message pipelines, processing chains, and usage analytics.</td></tr>
          <tr><td>Command cleanup</td><td>42 commands</td><td>Removed auth, experimental, internal, and stub commands.</td></tr>
        </tbody>
      </table>

      <h2 id="kept">Kept</h2>
      <ul>
        <li>AI chat and code generation</li>
        <li>File read, write, and edit workflows (<code>Read</code>/<code>Write</code>/<code>Edit</code>/<code>Glob</code>/<code>Grep</code>)</li>
        <li>Shell command execution (<code>Bash</code>/<code>PowerShell</code>)</li>
        <li>Git integration</li>
        <li>MCP protocol support</li>
        <li>LSP language services</li>
        <li>Plugin and skills systems</li>
        <li>Task management and permission controls</li>
      </ul>

      <h2 id="added">Added</h2>
      <table>
        <thead><tr><th>Feature</th><th>Description</th></tr></thead>
        <tbody>
          <tr><td><code>hard_deny</code></td><td>Unconditional deny rules with pipeline-level safety enforcement.</td></tr>
          <tr><td><code>html-output</code></td><td>Bundled skill for generating self-contained HTML reports.</td></tr>
          <tr><td><code>HUD status bar</code></td><td>Built-in status display for model, context, branch, permissions, and elapsed time.</td></tr>
          <tr><td>Third-party model optimization</td><td>Automatic request shaping for third-party API providers.</td></tr>
        </tbody>
      </table>
    </>} />
  )
}
