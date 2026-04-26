<script lang="ts" setup>
import { computed, ref, reactive } from 'vue'
import PageHeader from '@/components/app/PageHeader.vue'
import PageSkeleton from '@/components/app/PageSkeleton.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
  TableEmpty,
} from '@/components/ui/table'
import { Search, Scan, ChevronDown, ChevronRight, CircleCheck, CircleAlert, CircleX, Circle, Loader2 } from 'lucide-vue-next'
import { useClients } from '@/composables/useClients'
import { useMCPServers } from '@/composables/useMCPServers'
import { useAutoRefresh } from '@/composables/useAutoRefresh'
import { IntrospectMCPServer } from '../../wailsjs/go/main/App'

const { connectedCount } = useClients()
const { servers, updateServer, activeInspections } = useMCPServers()

const searchQuery = ref('')
const inspecting = reactive<Record<string, boolean>>({})
const expandedRows = reactive<Record<string, boolean>>({})
const isInspectingAll = ref(false)
const inspectProgress = ref({ done: 0, total: 0 })

const filteredServers = computed(() => {
  if (!searchQuery.value) return servers.value
  const q = searchQuery.value.toLowerCase()
  return servers.value.filter(
    s =>
      s.name.toLowerCase().includes(q) ||
      s.client.toLowerCase().includes(q) ||
      (s.command || '').toLowerCase().includes(q) ||
      (s.url || '').toLowerCase().includes(q),
  )
})

// Servers that can be introspected (stdio with a command)
const inspectableServers = computed(() =>
  servers.value.filter(s => s.transport === 'stdio' && s.command)
)

const { initialLoaded } = useAutoRefresh()

function serverKey(s: any): string {
  return `${s.client}-${s.name}-${s.sourcePath}`
}

const MAX_CONCURRENCY = 10

/**
 * Inspect a single server. Updates the row in-place when done.
 */
async function inspectOne(server: any) {
  const key = serverKey(server)
  inspecting[key] = true
  try {
    const result = await IntrospectMCPServer(server)
    updateServer(result)
    if (result.toolCount > 0) {
      expandedRows[key] = true
    }
  } catch (err) {
    console.error('Introspect failed:', err)
  } finally {
    inspecting[key] = false
    inspectProgress.value.done++
  }
}

/**
 * Batch-inspect all stdio servers with concurrency limit.
 * Like Clash Verge's connectivity test: one button, all servers,
 * results stream in as each completes.
 */
async function handleInspectAll() {
  const targets = inspectableServers.value
  if (targets.length === 0) return

  isInspectingAll.value = true
  activeInspections.value++
  inspectProgress.value = { done: 0, total: targets.length }

  // Semaphore-based concurrency limiter
  const queue = [...targets]
  const workers: Promise<void>[] = []

  for (let i = 0; i < Math.min(MAX_CONCURRENCY, queue.length); i++) {
    workers.push((async () => {
      while (queue.length > 0) {
        const server = queue.shift()!
        await inspectOne(server)
      }
    })())
  }

  await Promise.all(workers)

  isInspectingAll.value = false
  activeInspections.value--
}

function toggleExpand(key: string) {
  expandedRows[key] = !expandedRows[key]
}

function healthVariant(status: string): 'default' | 'secondary' | 'destructive' | 'outline' {
  switch (status) {
    case 'healthy': return 'default'
    case 'degraded': return 'secondary'
    case 'error': return 'destructive'
    default: return 'outline'
  }
}

function healthIcon(status: string) {
  switch (status) {
    case 'healthy': return CircleCheck
    case 'degraded': return CircleAlert
    case 'error': return CircleX
    default: return Circle
  }
}

function healthLabel(status: string): string {
  switch (status) {
    case 'healthy': return 'mcpServers.healthHealthy'
    case 'degraded': return 'mcpServers.healthDegraded'
    case 'error': return 'mcpServers.healthError'
    default: return 'mcpServers.healthUnknown'
  }
}
</script>

<template>
  <div class="flex flex-col flex-1 min-h-0">
    <PageHeader :title="$t('mcpServers.title')" :description="$t('mcpServers.description')">
      <template #actions>
        <div class="relative">
          <Search class="absolute left-2.5 top-2.5 size-4 text-muted-foreground" />
          <Input
            v-model="searchQuery"
            :placeholder="$t('actions.searchServers')"
            class="pl-8 w-[200px] h-9"
          />
        </div>
        <Button
          variant="outline"
          size="sm"
          :disabled="isInspectingAll || inspectableServers.length === 0"
          @click="handleInspectAll"
        >
          <Scan v-if="!isInspectingAll" class="mr-1.5 size-4" />
          <Loader2 v-else class="mr-1.5 size-4 animate-spin" />
          <template v-if="isInspectingAll">
            {{ inspectProgress.done }}/{{ inspectProgress.total }}
          </template>
          <template v-else>
            {{ $t('actions.inspect') }}
          </template>
        </Button>
      </template>
    </PageHeader>

    <PageSkeleton v-if="!initialLoaded" variant="table" />

    <div v-else class="flex-1 overflow-y-auto p-6">
      <div class="rounded-lg border overflow-x-hidden">
        <Table class="table-fixed">
          <TableHeader>
            <TableRow>
              <TableHead class="w-[120px]">{{ $t('mcpServers.client') }}</TableHead>
              <TableHead>{{ $t('mcpServers.name') }}</TableHead>
              <TableHead class="w-[100px]">{{ $t('mcpServers.transport') }}</TableHead>
              <TableHead>{{ $t('mcpServers.commandUrl') }}</TableHead>
              <TableHead class="w-[120px]">{{ $t('mcpServers.health') }}</TableHead>
              <TableHead class="w-[100px]">{{ $t('mcpServers.tools') }}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <template v-if="filteredServers.length > 0">
              <template
                v-for="(server, idx) in filteredServers"
                :key="serverKey(server)"
              >
                <!-- Main row -->
                <TableRow class="cursor-default">
                  <TableCell>
                    <Badge variant="secondary" class="text-xs font-medium">{{ server.client }}</Badge>
                  </TableCell>
                  <TableCell class="font-medium">
                    <div class="flex items-center gap-1.5">
                      <button
                        v-if="server.toolCount > 0"
                        class="p-0.5 rounded hover:bg-muted transition-colors"
                        @click="toggleExpand(serverKey(server))"
                      >
                        <component
                          :is="expandedRows[serverKey(server)] ? ChevronDown : ChevronRight"
                          class="size-3.5 text-muted-foreground"
                        />
                      </button>
                      <span>{{ server.name }}</span>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant="outline" class="text-xs">{{ server.transport || 'stdio' }}</Badge>
                  </TableCell>
                  <TableCell class="max-w-[300px]">
                    <code class="text-xs font-mono text-muted-foreground truncate block">
                      {{ server.url || server.command || '-' }}
                    </code>
                  </TableCell>
                  <TableCell>
                    <!-- Inspecting: show spinner -->
                    <div v-if="inspecting[serverKey(server)]" class="flex items-center gap-1.5 text-xs text-muted-foreground">
                      <Loader2 class="size-3.5 animate-spin" />
                      <span>{{ $t('actions.inspecting') }}</span>
                    </div>
                    <!-- Result: show health badge -->
                    <Badge
                      v-else
                      :variant="healthVariant(server.healthStatus)"
                      class="text-xs gap-1 transition-all duration-300"
                      :title="server.healthMessage || ''"
                    >
                      <component :is="healthIcon(server.healthStatus)" class="size-3" />
                      {{ $t(healthLabel(server.healthStatus)) }}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <span v-if="inspecting[serverKey(server)]" class="text-xs text-muted-foreground">...</span>
                    <span v-else-if="server.toolCount > 0" class="text-xs font-medium">
                      {{ $t('mcpServers.toolsCount', { count: server.toolCount }) }}
                    </span>
                    <span v-else-if="server.healthStatus === 'healthy'" class="text-xs text-muted-foreground">
                      {{ $t('mcpServers.noTools') }}
                    </span>
                    <span v-else class="text-xs text-muted-foreground">—</span>
                  </TableCell>
                </TableRow>

                <!-- Expanded tools row -->
                <TableRow
                  v-if="expandedRows[serverKey(server)] && server.tools?.length > 0"
                  class="bg-muted/30"
                >
                  <TableCell :colspan="6" class="p-0">
                    <div class="px-6 py-3">
                      <p class="text-xs font-medium text-muted-foreground mb-2 uppercase tracking-wider">
                        {{ $t('mcpServers.tools') }} ({{ server.tools.length }})
                      </p>
                      <div class="grid gap-1.5">
                        <div
                          v-for="tool in server.tools"
                          :key="tool.name"
                          class="flex items-start gap-3 text-xs py-1.5 px-3 rounded bg-background border"
                        >
                          <code class="font-mono font-medium text-foreground shrink-0">{{ tool.name }}</code>
                          <span class="text-muted-foreground">{{ tool.description || '—' }}</span>
                        </div>
                      </div>
                    </div>
                  </TableCell>
                </TableRow>
              </template>
            </template>
            <TableEmpty v-else :colspan="6">
              <div class="py-8 text-center">
                <p class="font-medium">
                  {{ connectedCount > 0 ? $t('mcpServers.noServersConfigured') : $t('mcpServers.noServersLoaded') }}
                </p>
                <p class="text-sm text-muted-foreground mt-1">
                  {{ connectedCount > 0 ? $t('mcpServers.noServersConfiguredDesc') : $t('mcpServers.noServersLoadedDesc') }}
                </p>
              </div>
            </TableEmpty>
          </TableBody>
        </Table>
      </div>
    </div>
  </div>
</template>
