import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

const toc = [
  { id: 'overview', text: 'Overview', level: 2 as const },
  { id: 'config', text: 'Configuration', level: 2 as const },
  { id: 'pipeline', text: 'Execution Pipeline', level: 2 as const },
  { id: 'rules', text: 'Rule Format', level: 2 as const },
  { id: 'examples', text: 'Examples', level: 2 as const },
]

export default function HardDeny() {
  return (
    <MarkdownDoc title="hard_deny Security Rules" description="Unconditional deny rules with pipeline-level safety enforcement." toc={toc} content={<>
      <DocStyles />
      <h2 id="overview">Overview</h2>
      <p><code>hard_deny</code> is MiniClaude&apos;s built-in safety mechanism. It checks every tool call at the <strong>front of the execution pipeline</strong>. When a rule matches, the call is rejected immediately: no prompt, no permission flow, and no mode can bypass it.</p>
      <blockquote>
        Unlike <code>soft_deny</code>, which relies on an AI classifier hint, <code>hard_deny</code> is enforced directly in code. It does not depend on model judgment or feature flags.
      </blockquote>

      <h2 id="config">Configuration</h2>
      <p>Configure it in <code>~/.claude/settings.json</code>:</p>
      <pre><code>{`{
  "autoMode": {
    "hard_deny": [
      "Bash(rm -rf)",
      "Bash(git push --force main)",
      "FileWrite(*.env)",
      "WebFetch"
    ]
  }
}`}</code></pre>

      <h2 id="pipeline">Execution Pipeline</h2>
      <pre><code>{`Tool call flow:
  checkHardDenyRules()    ← first line of defense
      ↓ match → reject immediately (no prompt / no bypass)
  Zod input validation
      ↓
  PreToolUse hooks
      ↓
  Permission checks (canUseTool / autoMode / bypassPermissions)
      ↓
  Tool execution`}</code></pre>
      <p><code>hard_deny</code> runs <strong>before everything else</strong>, so even <code>bypassPermissions</code> cannot skip it.</p>

      <h2 id="rules">Rule Format</h2>
      <table>
        <thead><tr><th>Format</th><th>Description</th><th>Example</th></tr></thead>
        <tbody>
          <tr><td><code>ToolName</code></td><td>Disable the entire tool.</td><td><code>WebFetch</code> — block all network fetches.</td></tr>
          <tr><td><code>ToolName(pattern)</code></td><td>Match a substring in the tool input.</td><td><code>Bash(rm -rf)</code> — block any command containing <code>rm -rf</code>.</td></tr>
        </tbody>
      </table>

      <h2 id="examples">Examples</h2>
      <h3>Prevent accidental deletion</h3>
      <pre><code>{`"Bash(rm -rf)"        // Block any rm -rf command
"Bash(rm -r)"         // Block any rm -r command`}</code></pre>
      <h3>Protect sensitive files</h3>
      <pre><code>{`"FileWrite(*.env)"    // Prevent overwriting .env files
"FileRead(*.pem)"     // Prevent reading private key files
"FileRead(*.key)"     // Prevent reading key files`}</code></pre>
      <h3>Block risky network access</h3>
      <pre><code>{`"Bash(curl 10.)"      // Block access to 10.x private networks
"Bash(curl 192.168.)" // Block access to 192.168.x private networks`}</code></pre>
      <h3>Enforce team policy</h3>
      <pre><code>{`"Bash(git push --force)"      // Block all force pushes
"Bash(npm publish)"            // Block direct npm publishes
"Bash(git commit --no-verify)" // Block skipping git hooks`}</code></pre>

      <blockquote>
        <strong>Tip:</strong> Combine precise and broad matches. For example, use both <code>"Bash(rm -rf /)"</code> for root-path protection and <code>"Bash(rm -rf)"</code> for general defense in depth.
      </blockquote>
    </>} />
  )
}
