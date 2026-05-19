import { logError } from '../../utils/log.js'
import { syncWebSocketURL } from './api.js'
import { loadSyncCredentials, updateSyncCredentials } from './credentials.js'
import {
  applySettingsSnapshot,
  localSettingsPushMessage,
  snapshotFromPayload,
} from './settingsBridge.js'
import type { SyncMessage } from './types.js'

let started = false
let activeSocket: WebSocket | null = null

export function startSyncAgent(): void {
  if (started) return
  started = true

  const credentials = loadSyncCredentials()
  if (!credentials) return

  try {
    const socket = new WebSocket(
      syncWebSocketURL(credentials.serverURL, credentials.accessToken),
    )
    activeSocket = socket

    socket.addEventListener('open', () => {
      socket.send(JSON.stringify({
        type: 'hello',
        client_id: credentials.clientID,
        workspace_id: credentials.workspaceID,
        session_id: credentials.sessionID,
        last_seen_version: credentials.settingsVersion ?? 0,
      }))
    })

    socket.addEventListener('message', event => {
      try {
        const message = JSON.parse(String(event.data)) as SyncMessage
        if (message.type === 'session_started') {
          rememberSession(message.payload)
          return
        }
        if (message.type === 'terminate_session') {
          terminateIfTargeted(message.payload)
          return
        }
        if (message.type === 'hello_ack' || message.type === 'pong') {
          return
        }
        if (
          message.type !== 'snapshot' &&
          message.type !== 'settings_updated' &&
          message.type !== 'settings_applied' &&
          message.type !== 'version_reject'
        ) {
          return
        }
        const snapshot = snapshotFromPayload(message.payload)
        if (!snapshot) return
        const error = applySettingsSnapshot(snapshot)
        if (error) logError(error)
      } catch (error) {
        logError(error)
      }
    })

    socket.addEventListener('error', event => {
      logError(new Error(`Sync WebSocket error: ${event.type}`))
    })
  } catch (error) {
    logError(error)
  }
}

export function pushLocalSettingsNow(): Error | null {
  if (!activeSocket || activeSocket.readyState !== WebSocket.OPEN) return null
  const message = localSettingsPushMessage()
  if (!message) return null
  try {
    activeSocket.send(JSON.stringify(message))
    return null
  } catch (error) {
    return error instanceof Error ? error : new Error(String(error))
  }
}

function rememberSession(payload: unknown): void {
  if (!payload || typeof payload !== 'object') return
  const value = payload as Record<string, unknown>
  if (typeof value.session_id !== 'string') return
  updateSyncCredentials(credentials => ({
    ...credentials,
    sessionID: value.session_id as string,
  }))
}

function terminateIfTargeted(payload: unknown): void {
  if (!payload || typeof payload !== 'object') return
  const value = payload as Record<string, unknown>
  const credentials = loadSyncCredentials()
  if (!credentials) return
  if (typeof value.client_id === 'string' && value.client_id !== credentials.clientID) {
    return
  }
  if (
    typeof value.session_id === 'string' &&
    credentials.sessionID &&
    value.session_id !== credentials.sessionID
  ) {
    return
  }
  activeSocket?.close(4000, 'terminated by sync server')
  activeSocket = null
}
