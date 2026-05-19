/**
 * GrowthBook integration - placeholder file
 * Original functionality removed; keep this file for compatibility
 */

export function getFeatureValue_CACHED_MAY_BE_STALE<T>(
  _featureName: string,
  defaultValue: T,
): T {
  return defaultValue
}

export function getFeatureValue_DEPRECATED<T>(
  _featureName: string,
  defaultValue: T,
): T {
  return defaultValue
}

export function getFeatureValue_CACHED_WITH_REFRESH<T>(
  _featureName: string,
  defaultValue: T,
): T {
  return defaultValue
}

export function getDynamicConfig_CACHED_MAY_BE_STALE<T>(
  _configName: string,
  defaultValue: T,
): T {
  return defaultValue
}

export async function getDynamicConfig_BLOCKS_ON_INIT<T>(
  _configName: string,
  defaultValue: T,
): Promise<T> {
  return defaultValue
}

export function checkStatsigFeatureGate_CACHED_MAY_BE_STALE(
  _gateName: string,
): boolean {
  return false
}

export async function checkGate_CACHED_OR_BLOCKING(
  _gateName: string,
): Promise<boolean> {
  return false
}

export async function checkSecurityRestrictionGate(
  _gateName: string,
): Promise<boolean> {
  return false
}

export function onGrowthBookRefresh(_callback: () => void): () => void {
  // Return an empty cleanup function
  return () => {}
}

export function initializeGrowthBook(): void {}

export function hasGrowthBookEnvOverride(): boolean {
  return false
}

export async function refreshGrowthBookAfterAuthChange(): Promise<void> {}

export function resetGrowthBook(): void {}

export function shutdownGrowthBook(): void {}
