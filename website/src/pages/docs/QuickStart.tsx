import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

const toc = [
  { id: 'install-bun', text: 'Install Bun', level: 2 as const },
  { id: 'deps', text: 'Install Dependencies and Configure', level: 2 as const },
  { id: 'start', text: 'Start MiniClaude', level: 2 as const },
  { id: 'global', text: 'Global Usage', level: 2 as const },
  { id: 'recovery', text: 'Recovery Mode', level: 2 as const },
]

export default function QuickStart() {
  return (
    <MarkdownDoc
      title="Install and Run"
      description="Run MiniClaude from scratch in three steps."
      toc={toc}
      content={
        <>
          <DocStyles />
          <h2 id="install-bun">Install Bun</h2>
          <p>MiniClaude is built on the Bun runtime, so install Bun 1.3.11 or newer first.</p>
          <h3>macOS / Linux</h3>
          <pre><code>{`curl -fsSL https://bun.sh/install | bash`}</code></pre>
          <h3>Homebrew</h3>
          <pre><code>{`brew install bun`}</code></pre>
          <h3>Windows</h3>
          <p>Visit <a href="https://bun.sh" target="_blank">bun.sh</a> to download the installer, or use PowerShell:</p>
          <pre><code>{`powershell -c "irm bun.sh/install.ps1 | iex"`}</code></pre>

          <h2 id="deps">Install Dependencies and Configure</h2>
          <pre><code>{`git clone https://github.com/txl16095/MiniClaude.git
cd MiniClaude
bun install`}</code></pre>
          <p>Copy the environment template and add your API key:</p>
          <pre><code>{`cp .env.example .env
# Edit the .env file
ANTHROPIC_API_KEY=sk-ant-xxxxx`}</code></pre>
          <p>If you use a third-party API such as DeepSeek, also set <code>ANTHROPIC_BASE_URL</code>:</p>
          <pre><code>{`ANTHROPIC_API_KEY=sk-your-deepseek-key
ANTHROPIC_BASE_URL=https://api.deepseek.com`}</code></pre>

          <h2 id="start">Start MiniClaude</h2>
          <h3>macOS / Linux</h3>
          <pre><code>{`# Build and run
bun run build
./cli

# Headless mode (non-interactive)
./cli -p "your question"

# Show all options
./cli --help`}</code></pre>
          <h3>Windows</h3>
          <p>Use Git Bash or WSL:</p>
          <pre><code>{`bun run build
./cli`}</code></pre>
          <p>You can also start it directly with Bun, which works well in PowerShell:</p>
          <pre><code>{`bun --env-file=.env run build
bun --env-file=.env cli`}</code></pre>

          <h2 id="global">Global Usage</h2>
          <p>Add the <code>bin/</code> directory to PATH so you can start MiniClaude from any directory:</p>
          <pre><code>{`# macOS / Linux
export PATH="$PWD/bin:$PATH"
# Or add it to ~/.bashrc or ~/.zshrc

# Windows PowerShell
$env:Path += ";E:\\Product\\MiniClaude\\bin"`}</code></pre>

          <h2 id="recovery">Recovery Mode</h2>
          <p>If you hit Ink TUI rendering issues, force the plain-text recovery mode:</p>
          <pre><code>{`export CLAUDE_CODE_FORCE_RECOVERY_CLI=1
./cli`}</code></pre>
        </>
      }
    />
  )
}
