/**
 * Analytics metadata - placeholder file
 * Original functionality removed; keep this file for compatibility
 */

export function sanitizeToolNameForAnalytics(toolName: string): string {
  return toolName
}

export function extractMcpToolDetails(_toolName: string): Record<string, unknown> {
  return {}
}

export function extractSkillName(_toolName: string): string | null {
  return null
}

export function extractToolInputForTelemetry(_input: unknown): Record<string, unknown> {
  return {}
}

export function getFileExtensionForAnalytics(_filePath: string): string {
  return ''
}

export function getFileExtensionsFromBashCommand(_command: string): string[] {
  return []
}

export function isToolDetailsLoggingEnabled(): boolean {
  return false
}

export function mcpToolDetailsForAnalytics(_toolName: string): Record<string, unknown> {
  return {}
}

export type AnalyticsMetadata_I_VERIFIED_THIS_IS_NOT_CODE_OR_FILEPATHS = string
export type AnalyticsMetadata_I_VERIFIED_THIS_IS_PII_TAGGED = string
