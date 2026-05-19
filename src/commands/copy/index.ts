/**
 * Copy command - minimal metadata only.
 * Implementation is lazy-loaded from copy.tsx to reduce startup time.
 */
import type { Command } from '../../commands.js'

const copy = {
  type: 'local-jsx',
  name: 'copy',
  description:
    "Copy Claude's last response to the clipboard (or use /copy N to copy the Nth-last response)",
  load: () => import('./copy.js'),
} satisfies Command

export default copy
