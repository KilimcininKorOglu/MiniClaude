import type { Command } from '../../commands.js'

const config = {
  aliases: ['settings'],
  type: 'local-jsx',
  name: 'config',
  description: 'Open the configuration panel',
  load: () => import('./config.js'),
} satisfies Command

export default config
