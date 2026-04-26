import { computed, ref } from 'vue'
import {
  ConnectClient,
  DisconnectClient,
  DiscoverClients,
  GetConnectedClientIDs,
} from '../../wailsjs/go/main/App'
import type { discovery } from '../../wailsjs/go/models'

const clients = ref<discovery.ClientInfo[]>([])
const connectedClientIDs = ref<string[]>([])

const connectedCount = computed(() => connectedClientIDs.value.length)
const connectedClients = computed(() => new Set(connectedClientIDs.value))

/**
 * Refresh clients in-place. Only replace the array when actual
 * content changes to avoid triggering unnecessary Vue re-renders.
 */
async function refreshClients() {
  const fresh = (await DiscoverClients()) ?? []
  // Only replace if the client list actually changed
  if (!arraysShallowEqual(clients.value, fresh, c => c.id + c.detected + c.status)) {
    clients.value = fresh
  }
}

/**
 * Refresh connected IDs in-place. Only replace when content changes.
 */
async function refreshConnections() {
  const fresh = (await GetConnectedClientIDs()) ?? []
  // Only replace if the ID list actually changed
  if (!stringArraysEqual(connectedClientIDs.value, fresh)) {
    connectedClientIDs.value = fresh
  }
}

async function connectClient(clientID: string) {
  await ConnectClient(clientID)
  await refreshConnections()
}

async function disconnectClient(clientID: string) {
  await DisconnectClient(clientID)
  await refreshConnections()
}

/** Fast equality check for string arrays (order-sensitive). */
function stringArraysEqual(a: string[], b: string[]): boolean {
  if (a.length !== b.length) return false
  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) return false
  }
  return true
}

/** Shallow equality check using a key extractor. */
function arraysShallowEqual<T>(a: T[], b: T[], key: (item: T) => string): boolean {
  if (a.length !== b.length) return false
  for (let i = 0; i < a.length; i++) {
    if (key(a[i]) !== key(b[i])) return false
  }
  return true
}

export function useClients() {
  return {
    clients,
    connectedClientIDs,
    connectedCount,
    connectedClients,
    refreshClients,
    refreshConnections,
    connectClient,
    disconnectClient,
  }
}
