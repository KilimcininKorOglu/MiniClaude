import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

export default function HtmlOutput() {
  return (
    <MarkdownDoc title="HTML Output" description="Generate self-contained HTML reports instead of long Markdown blocks." content={<>
      <DocStyles />
      <h2 id="overview">Overview</h2>
      <p>When you need rich output such as complex reports, code review results, or architecture proposals, MiniClaude can generate <strong>self-contained HTML files</strong> that are easier to read than long Markdown responses.</p>

      <h2 id="when">When to Use It</h2>
      <ul>
        <li><strong>Reports and analysis</strong>: data tables, multi-dimensional comparisons, and structured summaries.</li>
        <li><strong>Code review</strong>: multi-file findings, severity badges, and per-file breakdowns.</li>
        <li><strong>Architecture and planning</strong>: ASCII diagrams, decision trees, and dependency maps.</li>
        <li><strong>Search results</strong>: grouped output, clickable links, and metadata annotations.</li>
        <li><strong>Anything over 50 lines</strong>: readability drops quickly after that point in plain Markdown.</li>
      </ul>

      <h2 id="invoke">How to Invoke It</h2>
      <p>Use the <code>/html-output</code> command directly in the conversation:</p>
      <pre><code>{`/html-output Help me generate a code contribution summary for the last week`}</code></pre>

      <h2 id="features">HTML Features</h2>
      <table>
        <thead><tr><th>Feature</th><th>Description</th></tr></thead>
        <tbody>
          <tr><td>Self-contained</td><td>No external CSS, JavaScript, or font dependencies. A single file renders correctly on its own.</td></tr>
          <tr><td>Dark mode</td><td>Uses CSS variables and <code>prefers-color-scheme: dark</code> for automatic theme support.</td></tr>
          <tr><td>Responsive layout</td><td>Works in both desktop and mobile browsers.</td></tr>
          <tr><td>System font stack</td><td>No external font loading and no extra render delay.</td></tr>
          <tr><td>Tables and badges</td><td>Clear table styling and status badges for dense information.</td></tr>
          <tr><td>Portable</td><td>Open it with a double-click in a browser. No server is required.</td></tr>
        </tbody>
      </table>

      <h2 id="anti-patterns">Design Patterns to Avoid</h2>
      <ul>
        <li>Do not use gradient backgrounds or off-brand emoji.</li>
        <li>Do not load external resources such as Google Fonts, CDN CSS, or analytics scripts.</li>
        <li>Do not generate excessively large HTML files; keep them under roughly 100KB.</li>
        <li>Do not default to rounded containers with a left accent stripe in the generic AI-report style.</li>
      </ul>
    </>} />
  )
}
