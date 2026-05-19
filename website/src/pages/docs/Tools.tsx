import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

export default function Tools() {
  return (
    <MarkdownDoc title="Tool Reference" description="All built-in tools available to MiniClaude Agent." content={<>
      <DocStyles />
      <p>MiniClaude ships with 30+ built-in tools, and the agent can automatically choose the right ones for the current task.</p>

      <h2 id="file-tools">File Tools</h2>
      <table>
        <thead><tr><th>Tool</th><th>Description</th><th>Typical Use</th></tr></thead>
        <tbody>
          <tr><td><code>FileRead</code></td><td>Read file contents with syntax highlighting support.</td><td>Inspect source code or config files.</td></tr>
          <tr><td><code>FileWrite</code></td><td>Create or overwrite files.</td><td>Create new files or generate code.</td></tr>
          <tr><td><code>FileEdit</code></td><td>Edit files precisely with string replacement.</td><td>Patch snippets or refactor code.</td></tr>
          <tr><td><code>Glob</code></td><td>Search for files by filename pattern.</td><td>Find every <code>.ts</code> file.</td></tr>
          <tr><td><code>Grep</code></td><td>Search file contents with regular expressions via ripgrep.</td><td>Locate definitions or references.</td></tr>
          <tr><td><code>NotebookEdit</code></td><td>Edit Jupyter Notebook (<code>.ipynb</code>) files.</td><td>Update notebook cells.</td></tr>
        </tbody>
      </table>

      <h2 id="exec-tools">Execution Tools</h2>
      <table>
        <thead><tr><th>Tool</th><th>Description</th><th>Typical Use</th></tr></thead>
        <tbody>
          <tr><td><code>Bash</code></td><td>Run shell commands.</td><td>Git, npm, builds, and other CLI tasks.</td></tr>
          <tr><td><code>PowerShell</code></td><td>Run PowerShell commands on Windows.</td><td>Windows-specific administration.</td></tr>
        </tbody>
      </table>

      <h2 id="web-tools">Web Tools</h2>
      <table>
        <thead><tr><th>Tool</th><th>Description</th><th>Typical Use</th></tr></thead>
        <tbody>
          <tr><td><code>WebFetch</code></td><td>Fetch and analyze web page content.</td><td>Read documentation or inspect API responses.</td></tr>
          <tr><td><code>WebSearch</code></td><td>Search the web.</td><td>Look up recent information.</td></tr>
        </tbody>
      </table>

      <h2 id="task-tools">Task Management</h2>
      <table>
        <thead><tr><th>Tool</th><th>Description</th><th>Typical Use</th></tr></thead>
        <tbody>
          <tr><td><code>Task</code></td><td>Start a sub-agent for a complex task.</td><td>Parallelize multi-step work.</td></tr>
          <tr><td><code>TaskCreate</code></td><td>Create todo items.</td><td>Break work down and track it.</td></tr>
          <tr><td><code>TaskUpdate</code></td><td>Update task status.</td><td>Mark progress.</td></tr>
          <tr><td><code>TaskStop</code></td><td>Stop a running task.</td><td>Cancel long-running work.</td></tr>
        </tbody>
      </table>

      <h2 id="interaction-tools">Interaction Tools</h2>
      <table>
        <thead><tr><th>Tool</th><th>Description</th><th>Typical Use</th></tr></thead>
        <tbody>
          <tr><td><code>AskUserQuestion</code></td><td>Ask the user a question.</td><td>Clarify requirements or confirm an approach.</td></tr>
          <tr><td><code>Skill</code></td><td>Invoke an installed skill.</td><td>Run a skill-based workflow.</td></tr>
        </tbody>
      </table>

      <h2 id="plan-tools">Planning and Worktrees</h2>
      <table>
        <thead><tr><th>Tool</th><th>Description</th><th>Typical Use</th></tr></thead>
        <tbody>
          <tr><td><code>EnterPlanMode</code></td><td>Enter planning mode.</td><td>Design a complex feature before implementation.</td></tr>
          <tr><td><code>ExitPlanMode</code></td><td>Exit planning mode.</td><td>Start execution after the plan is approved.</td></tr>
          <tr><td><code>EnterWorktree</code></td><td>Create a Git worktree.</td><td>Isolate a development environment.</td></tr>
          <tr><td><code>ExitWorktree</code></td><td>Leave a worktree.</td><td>Clean up or keep the branch.</td></tr>
        </tbody>
      </table>

      <h2 id="mcp-tools">MCP Tools</h2>
      <table>
        <thead><tr><th>Tool</th><th>Description</th><th>Typical Use</th></tr></thead>
        <tbody>
          <tr><td><code>mcp__*</code></td><td>MCP protocol tools registered dynamically.</td><td>Browser control, APK reverse engineering, and more.</td></tr>
          <tr><td><code>ListMcpResourcesTool</code></td><td>List MCP resources.</td><td>Inspect available data sources.</td></tr>
          <tr><td><code>ReadMcpResourceTool</code></td><td>Read an MCP resource.</td><td>Fetch MCP-provided data.</td></tr>
        </tbody>
      </table>

      <blockquote>Tool permissions are controlled together by <code>permissions</code> and <code>autoMode.hard_deny</code> in <code>settings.json</code>. The <code>hard_deny</code> layer blocks matched tool calls unconditionally and has the highest priority.</blockquote>
    </>} />
  )
}
