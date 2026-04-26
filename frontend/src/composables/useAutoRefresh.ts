import { ref, onMounted, onUnmounted } from 'vue'
import { useClients } from './useClients'
import { useMCPServers } from './useMCPServers'
import { useSkills } from './useSkills'

const POLL_INTERVAL = 5000 // 5 seconds

let pollTimer: ReturnType<typeof setInterval> | null = null
let activeSubscribers = 0
const lastRefresh = ref(0)

/** True after the first poll cycle completes. Pages use this to
 *  show a loading skeleton until data is available. */
const initialLoaded = ref(false)

async function poll() {
  const { refreshClients, refreshConnections } = useClients()
  const { refreshServers } = useMCPServers()
  const { refreshSkills } = useSkills()

  await refreshClients()
  await refreshConnections()
  await Promise.all([refreshServers(true), refreshSkills()])
  lastRefresh.value = Date.now()
  if (!initialLoaded.value) initialLoaded.value = true
}

function startPolling() {
  activeSubscribers++
  if (pollTimer === null) {
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

export function useAutoRefresh() {
  onMounted(startPolling)
  onUnmounted(stopPolling)

  return { lastRefresh, initialLoaded }
}
