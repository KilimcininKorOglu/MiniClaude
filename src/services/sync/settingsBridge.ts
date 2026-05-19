import { createHash } from 'crypto'
import { dirname } from 'path'
import { writeFileSyncAndFlush_DEPRECATED } from '../../utils/file.js'
import { getFsImplementation } from '../../utils/fsOperations.js'
import { jsonStringify } from '../../utils/slowOperations.js'
import { markInternalWrite } from '../../utils/settings/internalWrites.js'
import {
  getSettingsFilePathForSource,
  getSettingsForSource,
} from '../../utils/settings/settings.js'
import { resetSettingsCache } from '../../utils/settings/settingsCache.js'
import { loadSyncCredentials, updateSyncCredentials } from './credentials.js'
import type { SettingsSnapshot } from './types.js'

type LocalSettingsPush = {
  checksum: string
  message: Record<string, unknown>
}

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
      settingsChecksum: checksumDocument(snapshot.document),
    }))
    return null
  } catch (error) {
    return error instanceof Error ? error : new Error(String(error))
  }
}

export function localSettingsPushMessage(): LocalSettingsPush | null {
  const credentials = loadSyncCredentials()
  if (!credentials) return null
  const document = getSettingsForSource('userSettings') ?? {}
  return {
    checksum: checksumDocument(document),
    message: {
      type: 'settings_push',
      base_version: credentials.settingsVersion ?? 0,
      document,
    },
  }
}

export function localSettingsChecksum(): string | null {
  const document = getSettingsForSource('userSettings') ?? {}
  return checksumDocument(document)
}

export function settingsFilePath(): string | null {
  return getSettingsFilePathForSource('userSettings') ?? null
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

function checksumDocument(document: Record<string, unknown>): string {
  return createHash('sha256').update(stableStringify(document)).digest('hex')
}

function stableStringify(value: unknown): string {
  return JSON.stringify(sortValue(value)) ?? 'null'
}

function sortValue(value: unknown): unknown {
  if (Array.isArray(value)) {
    return value.map(sortValue)
  }
  if (!value || typeof value !== 'object') {
    return value
  }
  const sorted: Record<string, unknown> = {}
  for (const key of Object.keys(value as Record<string, unknown>).sort()) {
    sorted[key] = sortValue((value as Record<string, unknown>)[key])
  }
  return sorted
}
