import { computed, ref, reactive } from 'vue'
import { DiscoverMCPServersForClients } from '../../wailsjs/go/main/App'
import type { discovery } from '../../wailsjs/go/models'
import { useClients } from './useClients'

const servers = ref<discovery.MCPServer[]>([])
const isLoading = ref(false)
const errorMessage = ref('')

// Count of currently running inspections. When > 0, auto-refresh
// skips server polling to avoid reactive conflicts and flickering.
const activeInspections = ref(0)

// Cache introspection results so they survive auto-refresh cycles.
// Keyed by "client:name:sourcePath".
const introspectionCache = reactive<Record<string, {
  tools: any[]
  toolCount: number
  healthStatus: string
  healthMessage: string
  introspectedAt: string
  status: string
}>>({})

const hasServers = computed(() => servers.value.length > 0)

function serverKey(s: { client: string; name: string; sourcePath: string }): string {
  return `${s.client}:${s.name}:${s.sourcePath}`
}

/** Save introspection results into the cache. */
function cacheIntrospection(s: discovery.MCPServer) {
  introspectionCache[serverKey(s)] = {
    tools: s.tools ?? [],
    toolCount: s.toolCount ?? 0,
    healthStatus: s.healthStatus ?? '',
    healthMessage: s.healthMessage ?? '',
    introspectedAt: s.introspectedAt ?? '',
    status: s.status ?? '',
  }
}

/**
 * Patch a single server's fields in-place on the reactive proxy.
 * Only assigns when the value actually changed to avoid triggering
 * unnecessary Vue re-renders (see frontend-pitfalls.md).
 */
function patchServer(cur: discovery.MCPServer, src: discovery.MCPServer) {
  // Config fields (from backend)
  if (cur.command !== src.command) cur.command = src.command
  if (cur.args !== src.args) cur.args = src.args
  if (cur.url !== src.url) cur.url = src.url
  if (cur.env !== src.env) cur.env = src.env
  if (cur.cwd !== src.cwd) cur.cwd = src.cwd
  if (cur.transport !== src.transport) cur.transport = src.transport
  if (cur.auth !== src.auth) cur.auth = src.auth
  if (cur.origin !== src.origin) cur.origin = src.origin
  if (cur.originPath !== src.originPath) cur.originPath = src.originPath

  // Introspection fields — only update if source actually has data
  // (avoids overwriting cached results with backend defaults)
  const srcHealth = src.healthStatus ?? ''
  const srcStatus = src.status ?? ''
  if (cur.healthStatus !== srcHealth) cur.healthStatus = srcHealth
  if (cur.healthMessage !== (src.healthMessage ?? '')) cur.healthMessage = src.healthMessage ?? ''
  if (cur.status !== srcStatus) cur.status = srcStatus
  if (cur.toolCount !== (src.toolCount ?? 0)) cur.toolCount = src.toolCount ?? 0
  if (cur.introspectedAt !== (src.introspectedAt ?? '')) cur.introspectedAt = src.introspectedAt ?? ''

  // Tools array — only replace reference if content actually changed
  const srcTools = src.tools ?? []
  if (cur.tools !== srcTools) {
    // Quick length check before deep-comparing
    if (!cur.tools || cur.tools.length !== srcTools.length ||
        JSON.stringify(cur.tools) !== JSON.stringify(srcTools)) {
      cur.tools = srcTools
    }
  }
}

/** Update a server in-place in the reactive list after introspection. */
function updateServer(result: discovery.MCPServer) {
  const idx = servers.value.findIndex(
    s => s.name === result.name && s.client === result.client && s.sourcePath === result.sourcePath,
  )
  if (idx >= 0) {
    patchServer(servers.value[idx], result)
  }
  cacheIntrospection(result)
}

/**
 * Merge fresh server list into the existing reactive array in-place.
 * This avoids replacing the entire array reference which would trigger
 * a full Vue re-render and cause visible flickering (see frontend-pitfalls.md).
 *
 * Strategy:
 *  1. Apply cached introspection data to fresh items
 *  2. Patch existing items in-place (only changed fields)
 *  3. Remove stale items
 *  4. Append genuinely new items
 */
function mergeServers(fresh: discovery.MCPServer[]) {
  const freshByKey = new Map<string, discovery.MCPServer>()
  for (const s of fresh) {
    const cached = introspectionCache[serverKey(s)]
    if (cached) {
      s.tools = cached.tools
      s.toolCount = cached.toolCount
      s.healthStatus = cached.healthStatus
      s.healthMessage = cached.healthMessage
      s.introspectedAt = cached.introspectedAt
      s.status = cached.status
    }
    freshByKey.set(serverKey(s), s)
  }

  // Patch existing items & track which keys are still present
  const existingKeys = new Set<string>()
  for (let i = 0; i < servers.value.length; i++) {
    const key = serverKey(servers.value[i])
    existingKeys.add(key)
    const updated = freshByKey.get(key)
    if (updated) {
      patchServer(servers.value[i], updated)
    }
  }

  // Remove stale items (no longer in backend)
  for (let i = servers.value.length - 1; i >= 0; i--) {
    if (!freshByKey.has(serverKey(servers.value[i]))) {
      servers.value.splice(i, 1)
    }
  }

  // Add genuinely new items
  for (const [key, s] of freshByKey) {
    if (!existingKeys.has(key)) {
      servers.value.push(s)
    }
  }
}

async function refreshServers(silent = false) {
  // Skip server refresh entirely while inspections are in progress.
  // This prevents reactive conflicts and flickering during inspect.
  if (silent && activeInspections.value > 0) return

  const { connectedClientIDs } = useClients()

  // Only touch isLoading for non-silent (manual) refreshes.
  // Silent (auto-poll) refreshes must not touch loading state at all
  // to avoid flickering the refresh button (see frontend-pitfalls.md).
  if (!silent) isLoading.value = true

  try {
    if (connectedClientIDs.value.length === 0) {
      if (servers.value.length > 0) servers.value.splice(0)
      return
    }
    const fresh = (await DiscoverMCPServersForClients(connectedClientIDs.value)) ?? []
    mergeServers(fresh)
  } catch (error) {
    if (!silent) {
      errorMessage.value = error instanceof Error ? error.message : String(error)
    }
  } finally {
    if (!silent) isLoading.value = false
  }
}

export function useMCPServers() {
  return {
    servers,
    isLoading,
    errorMessage,
    hasServers,
    activeInspections,
    refreshServers,
    updateServer,
  }
}
