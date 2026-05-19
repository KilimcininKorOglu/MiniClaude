export type DeviceStartResponse = {
  device_code: string
  user_code: string
  verification_url: string
  expires_at: string
  interval_seconds: number
}

export type DevicePollResponse = {
  status: string
  client_id?: string
  workspace_id?: string
  access_token?: string
  expires_at?: string
}

export type SettingsSnapshot = {
  workspace_id: string
  document: Record<string, unknown>
  version: number
  checksum: string
}

export type SyncCredentials = {
  serverURL: string
  clientID: string
  workspaceID: string
  accessToken: string
  expiresAt: string
  settingsVersion?: number
  sessionID?: string
}

export type SyncMessage = {
  type:
    | 'session_started'
    | 'hello_ack'
    | 'snapshot'
    | 'settings_updated'
    | 'settings_applied'
    | 'version_reject'
    | 'terminate_session'
    | 'pong'
    | 'error'
    | string
  workspace_id?: string
  version?: number
  payload?: unknown
}
