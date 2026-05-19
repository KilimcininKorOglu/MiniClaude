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
  settingsChecksum?: string
  sessionID?: string
}

export type SessionStartedPayload = {
  client_id?: string
  session_id?: string
}

export type TerminateSessionPayload = {
  client_id?: string
  session_id?: string
  reason?: string
}

export type SyncMessage =
  | {
      type: 'session_started'
      workspace_id?: string
      payload?: SessionStartedPayload
    }
  | {
      type: 'hello_ack'
      workspace_id?: string
      payload?: SessionStartedPayload
    }
  | {
      type: 'snapshot' | 'settings_updated' | 'settings_applied' | 'version_reject'
      workspace_id: string
      version: number
      payload: SettingsSnapshot
    }
  | {
      type: 'terminate_session'
      workspace_id?: string
      payload?: TerminateSessionPayload
    }
  | {
      type: 'pong'
      workspace_id?: string
      payload?: { session_id?: string }
    }
  | {
      type: 'error'
      workspace_id?: string
      payload?: unknown
    }
