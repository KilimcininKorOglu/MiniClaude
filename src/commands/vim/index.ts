import type { Command } from '../../commands.js'

const command = {
  name: 'vim',
  description: 'Switch between Vim and normal editing mode',
  supportsNonInteractive: false,
  type: 'local',
  load: () => import('./vim.js'),
} satisfies Command

export default command
