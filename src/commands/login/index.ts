import type { Command } from '../../commands.js'

export default {
  type: 'local-jsx',
  name: 'login',
  get description() {
    return 'Link MiniClaude to a sync server account'
  },
  argumentHint: '[sync-server-url]',
  load: () => import('./login.js'),
} satisfies Command
