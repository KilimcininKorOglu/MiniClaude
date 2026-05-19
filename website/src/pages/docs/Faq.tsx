import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

const toc = [
  { id: 'start-fail', text: 'Startup Failure', level: 2 as const },
  { id: 'permission', text: 'Permission Issues', level: 2 as const },
  { id: 'slow', text: 'Slow Model Responses', level: 2 as const },
  { id: 'build', text: 'Build Errors', level: 2 as const },
  { id: 'config', text: 'Settings Not Applied', level: 2 as const },
  { id: 'other', text: 'Other Problems', level: 2 as const },
]

export default function Faq() {
  return (
    <MarkdownDoc title="FAQ" description="Common MiniClaude problems and how to solve them." toc={toc} content={<>
      <DocStyles />
      <h2 id="start-fail">Startup Failure</h2>
      <p><strong>Symptom:</strong> Running <code>./cli</code> throws an error or exits immediately.</p>
      <p><strong>Fix:</strong></p>
      <ul>
        <li>Make sure you already ran <code>bun run build</code>.</li>
        <li>Confirm that <code>.env</code> exists and <code>ANTHROPIC_API_KEY</code> is configured.</li>
        <li>Check the Bun version with <code>bun --version</code>; you need 1.3.11 or newer.</li>
        <li>Try recovery mode: <code>CLAUDE_CODE_FORCE_RECOVERY_CLI=1 ./cli</code>.</li>
      </ul>

      <h2 id="permission">Permission Issues</h2>
      <p><strong>Symptom:</strong> Tool execution is denied or permission prompts appear too often.</p>
      <p><strong>Fix:</strong></p>
      <ul>
        <li>Configure permission rules in <code>settings.json</code>:<br/>
          <code>"permissions": &#123; "allow": ["Bash(git *)"], "defaultMode": "acceptEdits" &#125;</code></li>
        <li>Use <code>hard_deny</code> rules to block high-risk operations unconditionally.</li>
        <li>In a development environment, you can set <code>"defaultMode": "bypassPermissions"</code> to skip confirmations.</li>
      </ul>

      <h2 id="slow">Slow Model Responses</h2>
      <p><strong>Symptom:</strong> The model responds slowly or times out frequently.</p>
      <p><strong>Fix:</strong></p>
      <ul>
        <li>Switch to a faster model, such as <code>/model deepseek-v4-pro</code>.</li>
        <li>Reduce the context size with <code>/compact</code>.</li>
        <li>Use <code>/fast</code> mode to reduce thinking tokens.</li>
        <li>Check your network connection and proxy configuration.</li>
      </ul>

      <h2 id="build">Build Errors</h2>
      <p><strong>Symptom:</strong> <code>bun run build</code> fails.</p>
      <p><strong>Fix:</strong></p>
      <ul>
        <li>Run <code>bun install</code> to make sure dependencies are present.</li>
        <li>Clear the install and reinstall: <code>rm -rf node_modules && bun install</code>.</li>
        <li>Check whether disk space is sufficient.</li>
        <li>On Windows, make sure you are using Git Bash or WSL.</li>
      </ul>

      <h2 id="config">Settings Not Applied</h2>
      <p><strong>Symptom:</strong> Nothing changes after editing <code>settings.json</code>.</p>
      <p><strong>Fix:</strong></p>
      <ul>
        <li>The main config file is <code>~/.claude/settings.json</code>, not the project root.</li>
        <li>Verify the JSON syntax, for example with <code>cat settings.json | python -m json.tool</code>.</li>
        <li>Restart MiniClaude after saving the file.</li>
      </ul>

      <h2 id="other">Other Problems</h2>
      <ul>
        <li><strong>Windows Git Bash text encoding issues:</strong> set <code>LANG=en_US.UTF-8</code> or another UTF-8 locale.</li>
        <li><strong>MCP server connection failed:</strong> check that the server process is running and the configured path is correct.</li>
        <li><strong>Plugin installation failed:</strong> confirm that GitHub is reachable and review the marketplace configuration.</li>
      </ul>
    </>} />
  )
}
