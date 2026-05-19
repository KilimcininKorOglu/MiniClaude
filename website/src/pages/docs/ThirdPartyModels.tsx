import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

const toc = [
  { id: 'deepseek', text: 'DeepSeek', level: 2 as const },
  { id: 'openai', text: 'OpenAI-Compatible APIs', level: 2 as const },
  { id: 'custom', text: 'Other Compatible APIs', level: 2 as const },
  { id: 'model-switch', text: 'Switch Models', level: 2 as const },
]

export default function ThirdPartyModels() {
  return (
    <MarkdownDoc title="Third-Party Models" description="Connect DeepSeek, OpenAI-compatible services, and other third-party models." toc={toc} content={<>
      <DocStyles />
      <p>MiniClaude can connect to <strong>any Anthropic-compatible API</strong>, so you are not limited to a paid Anthropic account.</p>

      <h2 id="deepseek">DeepSeek</h2>
      <p>Configure these variables in <code>.env</code>:</p>
      <pre><code>{`ANTHROPIC_API_KEY=sk-your-deepseek-api-key
ANTHROPIC_BASE_URL=https://api.deepseek.com
ANTHROPIC_MODEL=deepseek-v4-pro`}</code></pre>
      <p>Common DeepSeek models:</p>
      <ul>
        <li><code>deepseek-v4-pro</code> — flagship model for demanding coding tasks.</li>
        <li><code>deepseek-r1</code> — reasoning model for deeper analysis.</li>
        <li><code>deepseek-v3</code> — general-purpose coding and chat.</li>
      </ul>

      <h2 id="openai">OpenAI-Compatible APIs</h2>
      <p>Any provider that supports an OpenAI-compatible Messages API can be used:</p>
      <pre><code>{`ANTHROPIC_API_KEY=sk-your-openai-key
ANTHROPIC_BASE_URL=https://api.openai.com
ANTHROPIC_MODEL=gpt-5-1-codex`}</code></pre>

      <h2 id="custom">Other Compatible APIs</h2>
      <p>MiniClaude also works with services that expose Anthropic-compatible APIs, including but not limited to:</p>
      <ul>
        <li><strong>Novita AI</strong> — lower-cost API gateway.</li>
        <li><strong>SiliconFlow</strong> — large-model platform with Anthropic-compatible endpoints.</li>
        <li><strong>OpenRouter</strong> — multi-model routing service.</li>
        <li><strong>Self-hosted services</strong> — gateways such as LiteLLM or One API.</li>
      </ul>
      <pre><code>{`# OpenRouter example
ANTHROPIC_API_KEY=sk-or-v1-xxx
ANTHROPIC_BASE_URL=https://openrouter.ai/api
ANTHROPIC_MODEL=anthropic/claude-sonnet-4-6`}</code></pre>

      <h2 id="model-switch">Switch Models</h2>
      <p>You can switch models at any time after startup:</p>
      <pre><code>{`# Use the /model command
/model deepseek-r1

# Or update the setting directly
# settings.json → "model": "deepseek-v4-pro"`}</code></pre>
      <blockquote>Model pricing and speed vary widely. DeepSeek V4 Pro works well for everyday use, while Claude Opus 4.6 is a better fit for the heaviest tasks.</blockquote>
    </>} />
  )
}
