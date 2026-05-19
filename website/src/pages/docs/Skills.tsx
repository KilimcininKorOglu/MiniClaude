import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

export default function Skills() {
  return (
    <MarkdownDoc title="Skills System" description="Skills are reusable, composable agent workflows." content={<>
      <DocStyles />
      <h2 id="overview">Overview</h2>
      <p>Skills are reusable workflow modules in MiniClaude. Each skill defines a <strong>standard operating procedure</strong>, and the agent loads the matching prompt automatically when the trigger conditions are met.</p>

      <h2 id="builtin">Bundled Skills</h2>
      <table>
        <thead><tr><th>Skill</th><th>Description</th><th>Trigger</th></tr></thead>
        <tbody>
          <tr><td><code>html-output</code></td><td>Generate HTML reports.</td><td><code>/html-output</code></td></tr>
          <tr><td><code>simplify</code></td><td>Review and clean up code.</td><td><code>/simplify</code></td></tr>
          <tr><td><code>debug</code></td><td>Debug the system.</td><td><code>/debug</code></td></tr>
          <tr><td><code>batch</code></td><td>Run batch operations.</td><td><code>/batch</code></td></tr>
          <tr><td><code>stuck</code></td><td>Diagnose blocked work.</td><td><code>/stuck</code></td></tr>
          <tr><td><code>verify</code></td><td>Verify code behavior.</td><td><code>/verify</code></td></tr>
          <tr><td><code>update-config</code></td><td>Update configuration.</td><td><code>/update-config</code></td></tr>
          <tr><td><code>remember</code></td><td>Manage memory.</td><td><code>/remember</code></td></tr>
        </tbody>
      </table>

      <h2 id="custom">Custom Skills</h2>
      <p>Create a <code>SKILL.md</code> file under <code>~/.claude/skills/</code>:</p>
      <pre><code>{`~/.claude/skills/
└── my-skill/
    └── SKILL.md`}</code></pre>
      <p>Example <code>SKILL.md</code> format:</p>
      <pre><code>{`---
name: my-skill
description: Skill description used for automatic matching
---
# Skill Title

Detailed prompt instructions for the skill...`}</code></pre>
      <p>The agent can <strong>match the skill automatically</strong> from conversation content, so manual invocation is not always required.</p>

      <h2 id="code">Code Registration</h2>
      <p>Bundled skills are registered with <code>registerBundledSkill</code>:</p>
      <pre><code>{`import { registerBundledSkill } from '../bundledSkills.js'

registerBundledSkill({
  name: 'html-output',
  description: 'Generate HTML reports...',
  userInvocable: true,
  async getPromptForCommand(args) {
    return [{ type: 'text', text: skillPrompt }]
  },
})`}</code></pre>
    </>} />
  )
}
