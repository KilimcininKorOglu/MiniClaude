import {
  getSyncCredentialsPath,
  loadStoredSyncCredentials,
  saveStoredSyncCredentials,
} from './secureStorage.js'
import type { SyncCredentials } from './types.js'

export { getSyncCredentialsPath }

export function loadSyncCredentials(): SyncCredentials | null {
  return loadStoredSyncCredentials()
}

export function saveSyncCredentials(credentials: SyncCredentials): void {
  saveStoredSyncCredentials(credentials)
}

export function updateSyncCredentials(
  updater: (credentials: SyncCredentials) => SyncCredentials,
): void {
  const credentials = loadSyncCredentials()
  if (!credentials) return
  saveSyncCredentials(updater(credentials))
}
