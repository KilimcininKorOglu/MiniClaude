import type { Command } from '../../commands.js'

const status = {
  type: 'local-jsx',
  name: 'status',
  description:
    'Show Claude Code status, including version, model, account, API connection, and tool status',
  immediate: true,
  load: () => import('./status.js'),
} satisfies Command

export default status
