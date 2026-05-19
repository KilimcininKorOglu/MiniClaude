import type { Command } from '../../commands.js'

const agents = {
  type: 'local-jsx',
  name: 'agents',
  description: 'Manage agent configuration',
  load: () => import('./agents.js'),
} satisfies Command

export default agents
