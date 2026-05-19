import type { Command } from '../../commands.js'

const rewind = {
  description: `Restore code and/or conversation to an earlier point in time`,
  name: 'rewind',
  aliases: ['checkpoint'],
  argumentHint: '',
  type: 'local',
  supportsNonInteractive: false,
  load: () => import('./rewind.js'),
} satisfies Command

export default rewind
