import type { Command } from '../../commands.js'

const hooks = {
  type: 'local-jsx',
  name: 'hooks',
  description: 'View hook configuration for tool events',
  immediate: true,
  load: () => import('./hooks.js'),
} satisfies Command

export default hooks
