import type { DevicePollResponse, DeviceStartResponse } from './types.js'

function normalizeServerURL(serverURL: string): string {
  return serverURL.replace(/\/+$/, '')
}

function endpoint(serverURL: string, path: string): string {
  return `${normalizeServerURL(serverURL)}${path}`
}

async function postJSON<T>(url: string, body: unknown): Promise<T> {
  const response = await fetch(url, {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!response.ok) {
    const message = await response.text()
    throw new Error(message || `Request failed with HTTP ${response.status}`)
  }
  return (await response.json()) as T
}

export async function startDeviceLogin(
  serverURL: string,
  clientName: string,
): Promise<DeviceStartResponse> {
  return postJSON<DeviceStartResponse>(endpoint(serverURL, '/api/device/start'), {
    client_name: clientName,
  })
}

export async function pollDeviceLogin(
  serverURL: string,
  deviceCode: string,
): Promise<DevicePollResponse> {
  return postJSON<DevicePollResponse>(endpoint(serverURL, '/api/device/poll'), {
    device_code: deviceCode,
  })
}

export function syncWebSocketURL(serverURL: string, accessToken: string): string {
  const url = new URL(endpoint(serverURL, '/api/sync/ws'))
  url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:'
  url.searchParams.set('access_token', accessToken)
  return url.toString()
}
