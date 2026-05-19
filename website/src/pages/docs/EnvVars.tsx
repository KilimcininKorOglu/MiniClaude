import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

const toc = [
  { id: 'required', text: 'Required Variables', level: 2 as const },
  { id: 'model', text: 'Model Configuration', level: 2 as const },
  { id: 'cache', text: 'Caching Configuration', level: 2 as const },
  { id: 'debug', text: 'Debug Options', level: 2 as const },
  { id: 'advanced', text: 'Advanced Configuration', level: 2 as const },
]

const rows = (vars: [string, string, string][]) =>
  vars.map(([name, required, desc]) => (
    <tr key={name}>
      <td><code>{name}</code></td>
      <td><span className={`badge ${required === 'Required' ? 'badge-red' : required === 'Recommended' ? 'badge-yellow' : 'badge-green'}`}>{required}</span></td>
      <td>{desc}</td>
    </tr>
  ))

export default function EnvVars() {
  return (
    <MarkdownDoc title="Environment Variables" description="Every supported MiniClaude environment variable and how to use it." toc={toc} content={<>
      <DocStyles />
      <h2 id="required">Required Variables</h2>
      <table><thead><tr><th>Variable</th><th>Required</th><th>Description</th></tr></thead>
      <tbody>{rows([
        ['ANTHROPIC_API_KEY', 'Required', 'API key. Use an sk-ant-xxx key for Anthropic, or the provider-specific key for DeepSeek and other third-party services.'],
        ['ANTHROPIC_BASE_URL', 'Optional', 'Custom API endpoint. Leave unset to use the official Anthropic API.'],
      ])}</tbody></table>

      <h2 id="model">Model Configuration</h2>
      <table><thead><tr><th>Variable</th><th>Required</th><th>Description</th></tr></thead>
      <tbody>{rows([
        ['ANTHROPIC_MODEL', 'Optional', 'Default model, such as claude-sonnet-4-6, claude-opus-4-6, or deepseek-v4-pro.'],
        ['ANTHROPIC_SMALL_FAST_MODEL', 'Optional', 'Fast-task model. Defaults to claude-haiku-4-5.'],
      ])}</tbody></table>

      <h2 id="cache">Caching Configuration</h2>
      <p>Prompt caching can reduce API cost significantly. MiniClaude enables a five-minute TTL cache by default, and you can extend it to one hour with an environment variable.</p>
      <table><thead><tr><th>Variable</th><th>Required</th><th>Description</th></tr></thead>
      <tbody>{rows([
        ['ENABLE_PROMPT_CACHING_1H', 'Recommended', 'Enable a one-hour cache TTL instead of the default five minutes. Strongly recommended when you want lower API cost.'],
        ['DISABLE_PROMPT_CACHING', 'Optional', 'Disable prompt caching entirely. Not recommended.'],
        ['DISABLE_PROMPT_CACHING_HAIKU', 'Optional', 'Disable caching only for Haiku models.'],
        ['DISABLE_PROMPT_CACHING_SONNET', 'Optional', 'Disable caching only for Sonnet models.'],
        ['DISABLE_PROMPT_CACHING_OPUS', 'Optional', 'Disable caching only for Opus models.'],
      ])}</tbody></table>

      <h2 id="debug">Debug Options</h2>
      <table><thead><tr><th>Variable</th><th>Required</th><th>Description</th></tr></thead>
      <tbody>{rows([
        ['DEBUG', 'Optional', 'Enable debug logs. Set it to * for everything, or choose modules such as DEBUG=api,cli.'],
        ['CLAUDE_CODE_FORCE_RECOVERY_CLI', 'Optional', 'Force the plain-text recovery mode to work around Ink TUI issues.'],
      ])}</tbody></table>

      <h2 id="advanced">Advanced Configuration</h2>
      <table><thead><tr><th>Variable</th><th>Required</th><th>Description</th></tr></thead>
      <tbody>{rows([
        ['ANTHROPIC_API_KEY_HELPER', 'Optional', 'Path to an external script that resolves the API key dynamically.'],
        ['ANTHROPIC_AUTH_TOKEN', 'Optional', 'OAuth bearer token when you are not using an API key.'],
        ['CLAUDE_CODE_ENABLE_XAA', 'Optional', 'Enable XAA (SEP-990) IdP integration.'],
        ['DISABLE_AUTOUPDATER', 'Optional', 'Set to 1 to disable automatic update checks.'],
        ['CLAUDE_CODE_SHELL', 'Optional', 'Custom shell path. Defaults to the system shell.'],
        ['CLAUDE_DISABLE_STREAM_WATCHDOG', 'Optional', 'Disable the stream idle watchdog. By default it helps recover after Mac sleep and wake cycles.'],
        ['CLAUDE_STREAM_IDLE_TIMEOUT_MS', 'Optional', 'Stream idle timeout in milliseconds. The default is 180000, which is three minutes.'],
        ['CLAUDE_CODE_USE_POWERSHELL_TOOL', 'Optional', 'Prefer PowerShell on Windows. Enabled by default.'],
      ])}</tbody></table>
    </>} />
  )
}
