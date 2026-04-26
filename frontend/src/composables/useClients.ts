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

async function refreshClients() {
  clients.value = (await DiscoverClients()) ?? []
}

async function refreshConnections() {
  connectedClientIDs.value = (await GetConnectedClientIDs()) ?? []
}

async function connectClient(clientID: string) {
  await ConnectClient(clientID)
  await refreshConnections()
}

async function disconnectClient(clientID: string) {
  await DisconnectClient(clientID)
  await refreshConnections()
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
