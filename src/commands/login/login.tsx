import chalk from 'chalk'
import * as React from 'react'
import { Box, Text } from '../../ink.js'
import type { LocalJSXCommandCall } from '../../types/command.js'
import { openBrowser } from '../../utils/browser.js'
import { startDeviceLogin, pollDeviceLogin } from '../../services/sync/api.js'
import { saveSyncCredentials } from '../../services/sync/credentials.js'

const defaultServerURL = 'http://localhost:8080'
const maxPollAttempts = 300

function clientName(): string {
  return `${process.env.USER || process.env.USERNAME || 'MiniClaude'}@${process.env.HOSTNAME || 'local'}`
}

function serverURLFromArgs(args: string): string {
  return args.trim() || process.env.MINICLAUDE_SYNC_SERVER_URL || defaultServerURL
}

function LoginFlow({
  onDone,
  serverURL,
}: {
  onDone: (message: string) => void
  serverURL: string
}) {
  const [userCode, setUserCode] = React.useState<string | null>(null)
  const [verificationURL, setVerificationURL] = React.useState<string | null>(null)

  React.useEffect(() => {
    let cancelled = false

    async function run(): Promise<void> {
      try {
        const start = await startDeviceLogin(serverURL, clientName())
        const link = new URL(start.verification_url)
        link.searchParams.set('user_code', start.user_code)
        setUserCode(start.user_code)
        setVerificationURL(link.toString())
        const opened = await openBrowser(link.toString())
        const browserStatus = opened
          ? 'Opened the verification link in your browser.'
          : `Open this verification link: ${link.toString()}`

        for (let attempt = 0; attempt < maxPollAttempts && !cancelled; attempt++) {
          const poll = await pollDeviceLogin(serverURL, start.device_code)
          if (poll.status === 'approved') {
            if (!poll.client_id || !poll.workspace_id || !poll.access_token || !poll.expires_at) {
              onDone('Sync server approved the device but returned an incomplete token response.')
              return
            }
            saveSyncCredentials({
              serverURL,
              clientID: poll.client_id,
              workspaceID: poll.workspace_id,
              accessToken: poll.access_token,
              expiresAt: poll.expires_at,
            })
            onDone(
              `Linked MiniClaude to ${chalk.bold(serverURL)}. Restart MiniClaude to start settings sync.`,
            )
            return
          }
          if (poll.status !== 'pending') {
            onDone(`Device login failed with status: ${poll.status}`)
            return
          }
          await new Promise(resolve => setTimeout(resolve, start.interval_seconds * 1000))
        }

        if (!cancelled) {
          onDone(
            `${browserStatus}\nDevice code ${chalk.bold(start.user_code)} expired before approval.`,
          )
        }
      } catch (error) {
        onDone(error instanceof Error ? error.message : String(error))
      }
    }

    void run()
    return () => {
      cancelled = true
    }
  }, [onDone, serverURL])

  return (
    <Box flexDirection="column">
      <Text>Waiting for sync server approval...</Text>
      {userCode ? <Text>User code: {chalk.bold(userCode)}</Text> : null}
      {verificationURL ? <Text>Verification URL: {verificationURL}</Text> : null}
    </Box>
  )
}

export const call: LocalJSXCommandCall = async (onDone, _context, args) => {
  const serverURL = serverURLFromArgs(args)
  try {
    const parsed = new URL(serverURL)
    if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') {
      onDone('Sync server URL must use http:// or https://.')
      return
    }
  } catch {
    onDone(`Invalid sync server URL: ${serverURL}`)
    return
  }

  return <LoginFlow onDone={onDone} serverURL={serverURL.replace(/\/+$/, '')} />
}
