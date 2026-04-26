<script lang="ts" setup>
import { computed, ref, reactive } from 'vue'
import PageHeader from '@/components/app/PageHeader.vue'
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
import { RefreshCw, Search, Scan, ChevronDown, ChevronRight, CircleCheck, CircleAlert, CircleX, Circle } from 'lucide-vue-next'
import { useClients } from '@/composables/useClients'
import { useMCPServers } from '@/composables/useMCPServers'
import { useAutoRefresh } from '@/composables/useAutoRefresh'
import { IntrospectMCPServer } from '../../wailsjs/go/main/App'

const { connectedCount } = useClients()
const { servers, isLoading } = useMCPServers()

const searchQuery = ref('')
const inspecting = reactive<Record<string, boolean>>({})
const expandedRows = reactive<Record<string, boolean>>({})

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

const { manualRefresh } = useAutoRefresh()

function serverKey(s: any): string {
  return `${s.client}-${s.name}-${s.sourcePath}`
}

async function handleInspect(server: any, idx: number) {
  const key = serverKey(server)
  inspecting[key] = true
  try {
    const result = await IntrospectMCPServer(server)
    // Update the server in-place in the reactive array
    const serverIdx = servers.value.findIndex(
      (s: any) => s.name === server.name && s.client === server.client && s.sourcePath === server.sourcePath
    )
    if (serverIdx >= 0) {
      servers.value[serverIdx] = result
    }
    // Auto-expand if tools found
    if (result.toolCount > 0) {
      expandedRows[key] = true
    }
  } catch (err) {
    console.error('Introspect failed:', err)
  } finally {
    inspecting[key] = false
  }
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
        <Button variant="outline" size="sm" :disabled="isLoading" @click="manualRefresh">
          <RefreshCw class="mr-2 size-4" :class="{ 'animate-spin': isLoading }" />
          {{ $t('actions.refresh') }}
        </Button>
      </template>
    </PageHeader>

    <div class="flex-1 overflow-y-auto p-6">
      <div class="rounded-lg border overflow-x-hidden">
        <Table class="table-fixed">
          <TableHeader>
            <TableRow>
              <TableHead class="w-[120px]">{{ $t('mcpServers.client') }}</TableHead>
              <TableHead>{{ $t('mcpServers.name') }}</TableHead>
              <TableHead class="w-[100px]">{{ $t('mcpServers.transport') }}</TableHead>
              <TableHead>{{ $t('mcpServers.commandUrl') }}</TableHead>
              <TableHead class="w-[100px]">{{ $t('mcpServers.health') }}</TableHead>
              <TableHead class="w-[100px]">{{ $t('mcpServers.tools') }}</TableHead>
              <TableHead class="w-[90px]"></TableHead>
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
                    <Badge
                      :variant="healthVariant(server.healthStatus)"
                      class="text-xs gap-1"
                      :title="server.healthMessage || ''"
                    >
                      <component :is="healthIcon(server.healthStatus)" class="size-3" />
                      {{ $t(healthLabel(server.healthStatus)) }}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <span v-if="server.toolCount > 0" class="text-xs font-medium">
                      {{ $t('mcpServers.toolsCount', { count: server.toolCount }) }}
                    </span>
                    <span v-else-if="server.healthStatus === 'healthy'" class="text-xs text-muted-foreground">
                      {{ $t('mcpServers.noTools') }}
                    </span>
                    <span v-else class="text-xs text-muted-foreground">—</span>
                  </TableCell>
                  <TableCell>
                    <Button
                      v-if="server.transport === 'stdio' && server.command"
                      variant="ghost"
                      size="sm"
                      class="h-7 px-2 text-xs"
                      :disabled="inspecting[serverKey(server)]"
                      @click="handleInspect(server, idx)"
                    >
                      <Scan class="mr-1 size-3.5" :class="{ 'animate-pulse': inspecting[serverKey(server)] }" />
                      {{ inspecting[serverKey(server)] ? $t('actions.inspecting') : $t('actions.inspect') }}
                    </Button>
                  </TableCell>
                </TableRow>

                <!-- Expanded tools row -->
                <TableRow
                  v-if="expandedRows[serverKey(server)] && server.tools?.length > 0"
                  class="bg-muted/30"
                >
                  <TableCell :colspan="7" class="p-0">
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
            <TableEmpty v-else :colspan="7">
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
