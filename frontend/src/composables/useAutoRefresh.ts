import { ref, onMounted, onUnmounted } from 'vue'
import { useClients } from './useClients'
import { useMCPServers } from './useMCPServers'
import { useSkills } from './useSkills'

const POLL_INTERVAL = 5000 // 5 seconds

let pollTimer: ReturnType<typeof setInterval> | null = null
let activeSubscribers = 0
const lastRefresh = ref(0)

async function poll() {
  const { refreshClients, refreshConnections } = useClients()
  const { refreshServers } = useMCPServers()
  const { refreshSkills } = useSkills()

  await refreshClients()
  await refreshConnections()
  await Promise.all([refreshServers(true), refreshSkills()])
  lastRefresh.value = Date.now()
}

function startPolling() {
  activeSubscribers++
  if (pollTimer === null) {
    // Initial load
    poll()
    pollTimer = setInterval(poll, POLL_INTERVAL)
  }
}

function stopPolling() {
  activeSubscribers--
  if (activeSubscribers <= 0) {
    activeSubscribers = 0
    if (pollTimer !== null) {
      clearInterval(pollTimer)
      pollTimer = null
    }
  }
}

/** Manual refresh with visible loading indicator */
async function manualRefresh() {
  const { refreshClients, refreshConnections } = useClients()
  const { refreshServers } = useMCPServers()
  const { refreshSkills } = useSkills()

  const { isLoading } = useMCPServers()
  isLoading.value = true

  await refreshClients()
  await refreshConnections()
  await Promise.all([refreshServers(false), refreshSkills()])
  lastRefresh.value = Date.now()
}

export function useAutoRefresh() {
  onMounted(startPolling)
  onUnmounted(stopPolling)

  return {
    lastRefresh,
    manualRefresh,
  }
}
