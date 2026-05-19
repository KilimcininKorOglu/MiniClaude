import type { Command } from '../../commands.js'

export default {
  type: 'local-jsx',
  name: 'provider',
  get description() {
    return 'Switch model provider (DeepSeek / Kiro, etc.)'
  },
  argumentHint: '[provider-name] [--session]',
  load: () => import('./provider.js'),
} satisfies Command
