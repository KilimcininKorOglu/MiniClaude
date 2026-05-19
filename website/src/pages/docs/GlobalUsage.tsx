import MarkdownDoc, { DocStyles } from '../../components/MarkdownDoc'

export default function GlobalUsage() {
  return (
    <MarkdownDoc title="Global Usage" description="Run MiniClaude from any directory." content={<>
      <DocStyles />
      <h2 id="unix">macOS / Linux</h2>
      <p>Add MiniClaude&apos;s <code>bin/</code> directory to PATH:</p>
      <pre><code>{`# Add this to your shell config file (~/.bashrc, ~/.zshrc, and so on)
export PATH="$HOME/MiniClaude/bin:$PATH"

# Reload the shell config
source ~/.zshrc

# You can now start MiniClaude from any directory
cd ~/my-project
claude`}</code></pre>

      <h2 id="windows">Windows</h2>
      <p>Option 1: configure it in Git Bash:</p>
      <pre><code>{`echo 'export PATH="/e/Product/MiniClaude/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc`}</code></pre>
      <p>Option 2: add it temporarily in PowerShell:</p>
      <pre><code>{`$env:Path += ";E:\\Product\\MiniClaude\\bin"
claude`}</code></pre>
      <p>Option 3: set a permanent system environment variable:</p>
      <p>Open System Properties → Advanced → Environment Variables, then add <code>E:\Product\MiniClaude\bin</code> to Path.</p>

      <h2 id="verify">Verify the Installation</h2>
      <pre><code>{`# Show the version
claude --version

# Show help
claude --help

# Start chatting
claude`}</code></pre>

      <blockquote>Make sure the <code>.env</code> file lives in the MiniClaude project root, or promote the required variables to system-level environment variables.</blockquote>
    </>} />
  )
}
