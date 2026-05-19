import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

const toc = [
  { id: 'base', text: 'Basic Commands', level: 2 as const },
  { id: 'config', text: 'Configuration Commands', level: 2 as const },
  { id: 'file', text: 'File Commands', level: 2 as const },
  { id: 'dev', text: 'Development Commands', level: 2 as const },
  { id: 'adv', text: 'Advanced Commands', level: 2 as const },
  { id: 'cli-flags', text: 'CLI Flags', level: 2 as const },
]

export default function Commands() {
  return (
    <MarkdownDoc title="Command Reference" description="All available MiniClaude commands and how to use them." toc={toc} content={<>
      <DocStyles />
      <h2 id="base">Basic Commands</h2>
      <table>
        <thead><tr><th>Command</th><th>Description</th><th>Example</th></tr></thead>
        <tbody>
          <tr><td><code>/help</code></td><td>Show the help screen.</td><td><code>/help</code></td></tr>
          <tr><td><code>/clear</code></td><td>Clear the current conversation.</td><td><code>/clear</code></td></tr>
          <tr><td><code>/exit</code></td><td>Exit the program.</td><td><code>/exit</code></td></tr>
          <tr><td><code>/status</code></td><td>Show the current status, including model, branch, and permissions.</td><td><code>/status</code></td></tr>
          <tr><td><code>/stats</code></td><td>View usage statistics.</td><td><code>/stats</code></td></tr>
          <tr><td><code>/doctor</code></td><td>Run system diagnostics.</td><td><code>/doctor</code></td></tr>
        </tbody>
      </table>

      <h2 id="config">Configuration Commands</h2>
      <table>
        <thead><tr><th>Command</th><th>Description</th><th>Example</th></tr></thead>
        <tbody>
          <tr><td><code>/config</code></td><td>Open the configuration file.</td><td><code>/config</code></td></tr>
          <tr><td><code>/model</code></td><td>Switch the active AI model.</td><td><code>/model deepseek-v4-pro</code></td></tr>
          <tr><td><code>/theme</code></td><td>Switch the terminal theme.</td><td><code>/theme dark</code></td></tr>
          <tr><td><code>/provider</code></td><td>Switch the model provider with hot swapping.</td><td><code>/provider kiro</code></td></tr>
          <tr><td><code>/permissions</code></td><td>Manage permission rules.</td><td><code>/permissions</code></td></tr>
          <tr><td><code>/hooks</code></td><td>Manage lifecycle hooks.</td><td><code>/hooks</code></td></tr>
          <tr><td><code>/output-style</code></td><td>Set the response style.</td><td><code>/output-style concise</code></td></tr>
        </tbody>
      </table>

      <h2 id="file">File Commands</h2>
      <table>
        <thead><tr><th>Command</th><th>Description</th><th>Example</th></tr></thead>
        <tbody>
          <tr><td><code>/files</code></td><td>Show the files currently in context.</td><td><code>/files</code></td></tr>
          <tr><td><code>/add-dir</code></td><td>Add a directory to the working context.</td><td><code>/add-dir src/</code></td></tr>
          <tr><td><code>/diff</code></td><td>Show the current code changes.</td><td><code>/diff</code></td></tr>
          <tr><td><code>/copy</code></td><td>Copy the last response to the clipboard.</td><td><code>/copy</code></td></tr>
          <tr><td><code>/export</code></td><td>Export the conversation history.</td><td><code>/export</code></td></tr>
        </tbody>
      </table>

      <h2 id="dev">Development Commands</h2>
      <table>
        <thead><tr><th>Command</th><th>Description</th><th>Example</th></tr></thead>
        <tbody>
          <tr><td><code>/init</code></td><td>Initialize the project CLAUDE.md file.</td><td><code>/init</code></td></tr>
          <tr><td><code>/compact</code></td><td>Compact the conversation context.</td><td><code>/compact</code></td></tr>
          <tr><td><code>/review</code></td><td>Review the current code changes.</td><td><code>/review</code></td></tr>
          <tr><td><code>/pr_comments</code></td><td>View PR comments.</td><td><code>/pr_comments</code></td></tr>
          <tr><td><code>/plan</code></td><td>Enter planning mode.</td><td><code>/plan</code></td></tr>
          <tr><td><code>/rename</code></td><td>Rename the current conversation.</td><td><code>/rename "New Name"</code></td></tr>
        </tbody>
      </table>

      <h2 id="adv">Advanced Commands</h2>
      <table>
        <thead><tr><th>Command</th><th>Description</th><th>Example</th></tr></thead>
        <tbody>
          <tr><td><code>/mcp</code></td><td>Manage MCP servers.</td><td><code>/mcp</code></td></tr>
          <tr><td><code>/skills</code></td><td>Manage skills.</td><td><code>/skills</code></td></tr>
          <tr><td><code>/tasks</code></td><td>View background tasks.</td><td><code>/tasks</code></td></tr>
          <tr><td><code>/vim</code></td><td>Toggle Vim mode.</td><td><code>/vim</code></td></tr>
          <tr><td><code>/fast</code></td><td>Toggle fast mode.</td><td><code>/fast</code></td></tr>
          <tr><td><code>/effort</code></td><td>Set the reasoning depth.</td><td><code>/effort high</code></td></tr>
          <tr><td><code>/html-output</code></td><td>Generate an HTML report.</td><td><code>/html-output analysis report</code></td></tr>
        </tbody>
      </table>

      <h2 id="cli-flags">CLI Flags</h2>
      <table>
        <thead><tr><th>Flag</th><th>Description</th></tr></thead>
        <tbody>
          <tr><td><code>-p, --print</code></td><td>Run in non-interactive mode, answer once, then exit.</td></tr>
          <tr><td><code>--model</code></td><td>Choose the startup model.</td></tr>
          <tr><td><code>--permission-mode</code></td><td>Set the permission mode: default / acceptEdits / bypassPermissions / plan.</td></tr>
          <tr><td><code>--output-format</code></td><td>Select the output format: text / json / stream-json.</td></tr>
          <tr><td><code>--continue</code></td><td>Continue the most recent conversation.</td></tr>
          <tr><td><code>--resume</code></td><td>Resume a session by ID.</td></tr>
          <tr><td><code>--version</code></td><td>Show the version number.</td></tr>
          <tr><td><code>--dir</code></td><td>Set the startup working directory.</td></tr>
        </tbody>
      </table>
    </>} />
  )
}
