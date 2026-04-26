import { computed, ref } from 'vue'
import { DiscoverMCPServersForClients } from '../../wailsjs/go/main/App'
import type { discovery } from '../../wailsjs/go/models'
import { useClients } from './useClients'

const servers = ref<discovery.MCPServer[]>([])
const isLoading = ref(false)
const errorMessage = ref('')

const hasServers = computed(() => servers.value.length > 0)

async function refreshServers(silent = false) {
  const { connectedClientIDs } = useClients()

  if (!silent) isLoading.value = true
  errorMessage.value = ''

  try {
    if (connectedClientIDs.value.length === 0) {
      servers.value = []
      return
    }
    servers.value = (await DiscoverMCPServersForClients(connectedClientIDs.value)) ?? []
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : String(error)
  } finally {
    isLoading.value = false
  }
}

export function useMCPServers() {
  return {
    servers,
    isLoading,
    errorMessage,
    hasServers,
    refreshServers,
  }
}
