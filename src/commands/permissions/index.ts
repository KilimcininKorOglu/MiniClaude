import type { Command } from '../../commands.js'

const permissions = {
  type: 'local-jsx',
  name: 'permissions',
  aliases: ['allowed-tools'],
  description: 'Manage allow and deny rules for tool permissions',
  load: () => import('./permissions.js'),
} satisfies Command

export default permissions
