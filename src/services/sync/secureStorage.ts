import { existsSync, mkdirSync, readFileSync, unlinkSync, writeFileSync } from 'fs'
import { dirname, join } from 'path'
import { getClaudeConfigHomeDir } from '../../utils/envUtils.js'
import { safeParseJSON } from '../../utils/json.js'
import { jsonStringify } from '../../utils/slowOperations.js'
import type { SyncCredentials } from './types.js'

const credentialsPath = join(getClaudeConfigHomeDir(), 'sync-credentials.json')

function isSyncCredentials(value: unknown): value is SyncCredentials {
  if (!value || typeof value !== 'object') return false
  const candidate = value as Record<string, unknown>
  return (
    typeof candidate.serverURL === 'string' &&
    typeof candidate.clientID === 'string' &&
    typeof candidate.workspaceID === 'string' &&
    typeof candidate.accessToken === 'string' &&
    typeof candidate.expiresAt === 'string'
  )
}

export function getSyncCredentialsPath(): string {
  return credentialsPath
}

export function loadStoredSyncCredentials(): SyncCredentials | null {
  if (!existsSync(credentialsPath)) return null
  const parsed = safeParseJSON(readFileSync(credentialsPath, 'utf8'), false)
  return isSyncCredentials(parsed) ? parsed : null
}

export function saveStoredSyncCredentials(credentials: SyncCredentials): void {
  mkdirSync(dirname(credentialsPath), { recursive: true })
  writeFileSync(credentialsPath, jsonStringify(credentials, null, 2) + '\n', {
    mode: 0o600,
  })
}

export function deleteStoredSyncCredentials(): void {
  if (existsSync(credentialsPath)) unlinkSync(credentialsPath)
}
