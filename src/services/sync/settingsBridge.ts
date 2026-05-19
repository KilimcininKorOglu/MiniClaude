import { dirname } from 'path'
import { writeFileSyncAndFlush_DEPRECATED } from '../../utils/file.js'
import { getFsImplementation } from '../../utils/fsOperations.js'
import { jsonStringify } from '../../utils/slowOperations.js'
import { markInternalWrite } from '../../utils/settings/internalWrites.js'
import { getSettingsFilePathForSource } from '../../utils/settings/settings.js'
import { resetSettingsCache } from '../../utils/settings/settingsCache.js'
import { updateSyncCredentials } from './credentials.js'
import type { SettingsSnapshot } from './types.js'

export function applySettingsSnapshot(snapshot: SettingsSnapshot): Error | null {
  const filePath = getSettingsFilePathForSource('userSettings')
  if (!filePath) return null

  try {
    getFsImplementation().mkdirSync(dirname(filePath))
    markInternalWrite(filePath)
    writeFileSyncAndFlush_DEPRECATED(
      filePath,
      jsonStringify(snapshot.document, null, 2) + '\n',
    )
    resetSettingsCache()
    updateSyncCredentials(credentials => ({
      ...credentials,
      settingsVersion: snapshot.version,
    }))
    return null
  } catch (error) {
    return error instanceof Error ? error : new Error(String(error))
  }
}

export function snapshotFromPayload(payload: unknown): SettingsSnapshot | null {
  if (!payload || typeof payload !== 'object') return null
  const value = payload as Record<string, unknown>
  if (
    typeof value.workspace_id !== 'string' ||
    typeof value.version !== 'number' ||
    typeof value.checksum !== 'string' ||
    !value.document ||
    typeof value.document !== 'object' ||
    Array.isArray(value.document)
  ) {
    return null
  }
  return value as SettingsSnapshot
}
