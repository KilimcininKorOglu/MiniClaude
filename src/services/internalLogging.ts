/**
 * Internal logging - placeholder file
 * Original functionality removed; keep this file for compatibility
 */

import type { ToolPermissionContext } from '../Tool.js'

export async function getContainerId(): Promise<string | null> {
  return null
}

export async function logPermissionContextForAnts(
  _toolPermissionContext: ToolPermissionContext | null,
  _moment: 'summary' | 'initialization',
): Promise<void> {
  // Empty implementation
}
