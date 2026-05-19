import { basename, dirname } from 'path'
import { watch, type FSWatcher } from 'fs'
import { logError } from '../../utils/log.js'
import { getFsImplementation } from '../../utils/fsOperations.js'
import { resetSettingsCache } from '../../utils/settings/settingsCache.js'
import { syncWebSocketURL } from './api.js'
import {
  deleteSyncCredentials,
  loadSyncCredentials,
  updateSyncCredentials,
} from './credentials.js'
import {
  applySettingsSnapshot,
  localSettingsPushMessage,
  settingsFilePath,
  snapshotFromPayload,
} from './settingsBridge.js'
import type { SyncCredentials } from './types.js'

const settingsPushDebounceMs = 400
const reconnectDelayMs = 2_000

let started = false
let activeSocket: WebSocket | null = null
let settingsWatcher: FSWatcher | null = null
let settingsPushTimer: ReturnType<typeof setTimeout> | null = null
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
let stopped = false
let lastSentSettingsChecksum: string | null = null

export function startSyncAgent(): void {
  if (started) return

  const credentials = loadSyncCredentials()
  if (!credentials) return

  started = true
  stopped = false
  startSettingsWatcher()
  connect(credentials)
}

export function pushLocalSettingsNow(): Error | null {
  const push = localSettingsPushMessage()
  if (!push) return null
  if (!activeSocket || activeSocket.readyState !== WebSocket.OPEN) {
    scheduleLocalSettingsPush()
    return null
  }
  return sendLocalSettingsPush(push)
}

function connect(credentials: SyncCredentials): void {
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
      scheduleLocalSettingsPush()
    })

    socket.addEventListener('message', event => {
      try {
        const message = JSON.parse(String(event.data)) as { type?: unknown; payload?: unknown }
        if (typeof message.type !== 'string') return
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
        lastSentSettingsChecksum = null
        const error = applySettingsSnapshot(snapshot)
        if (error) logError(error)
      } catch (error) {
        logError(error)
      }
    })

    socket.addEventListener('error', event => {
      logError(new Error(`Sync WebSocket error: ${event.type}`))
    })

    socket.addEventListener('close', () => {
      if (activeSocket === socket) activeSocket = null
      if (!stopped) scheduleReconnect()
    })
  } catch (error) {
    logError(error)
    scheduleReconnect()
  }
}

function startSettingsWatcher(): void {
  if (settingsWatcher) return
  const filePath = settingsFilePath()
  if (!filePath) return
  const parent = dirname(filePath)
  const targetName = basename(filePath)

  try {
    getFsImplementation().mkdirSync(parent)
    settingsWatcher = watch(parent, { persistent: false }, (_event, filename) => {
      if (filename && basename(String(filename)) !== targetName) return
      resetSettingsCache()
      scheduleLocalSettingsPush()
    })
  } catch (error) {
    logError(error)
  }
}

function scheduleLocalSettingsPush(): void {
  if (settingsPushTimer) clearTimeout(settingsPushTimer)
  settingsPushTimer = setTimeout(() => {
    settingsPushTimer = null
    const credentials = loadSyncCredentials()
    const push = localSettingsPushMessage()
    if (!credentials || !push) return
    if (credentials.settingsChecksum === push.checksum) return
    if (lastSentSettingsChecksum === push.checksum) return
    if (!activeSocket || activeSocket.readyState !== WebSocket.OPEN) return
    const error = sendLocalSettingsPush(push)
    if (error) logError(error)
  }, settingsPushDebounceMs)
}

function sendLocalSettingsPush(push: ReturnType<typeof localSettingsPushMessage>): Error | null {
  if (!push || !activeSocket || activeSocket.readyState !== WebSocket.OPEN) return null
  try {
    activeSocket.send(JSON.stringify(push.message))
    lastSentSettingsChecksum = push.checksum
    return null
  } catch (error) {
    return error instanceof Error ? error : new Error(String(error))
  }
}

function scheduleReconnect(): void {
  if (reconnectTimer || stopped) return
  reconnectTimer = setTimeout(() => {
    reconnectTimer = null
    const credentials = loadSyncCredentials()
    if (!credentials) return
    connect(credentials)
  }, reconnectDelayMs)
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
  stopped = true
  if (settingsPushTimer) clearTimeout(settingsPushTimer)
  if (reconnectTimer) clearTimeout(reconnectTimer)
  settingsWatcher?.close()
  settingsWatcher = null
  activeSocket?.close(4000, 'terminated by sync server')
  activeSocket = null

  const clientRevoked = typeof value.client_id === 'string' && typeof value.session_id !== 'string'
  if (clientRevoked) {
    deleteSyncCredentials()
  } else {
    updateSyncCredentials(current => {
      const { sessionID: _sessionID, ...remaining } = current
      return remaining
    })
  }
  notifyTermination(value.reason, clientRevoked)
}

function notifyTermination(reason: unknown, clientRevoked: boolean): void {
  const defaultMessage = clientRevoked
    ? 'MiniClaude sync client was revoked by the server. Run /login again to relink.'
    : 'MiniClaude sync session was terminated by the server.'
  if (reason === 'client_revoked' || reason === 'session_terminated') {
    process.stderr.write(`\n${defaultMessage}\n`)
    return
  }
  const message = typeof reason === 'string' && reason.trim() ? reason.trim() : defaultMessage
  process.stderr.write(`\n${message}\n`)
}
