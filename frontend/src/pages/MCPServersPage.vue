<script lang="ts" setup>
import { computed, ref } from 'vue'
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
import { RefreshCw, Search } from 'lucide-vue-next'
import { useClients } from '@/composables/useClients'
import { useMCPServers } from '@/composables/useMCPServers'
import { useAutoRefresh } from '@/composables/useAutoRefresh'

const { connectedCount } = useClients()
const { servers, isLoading } = useMCPServers()

const searchQuery = ref('')

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
              <TableHead class="w-[100px]">{{ $t('mcpServers.status') }}</TableHead>
              <TableHead>{{ $t('mcpServers.source') }}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <template v-if="filteredServers.length > 0">
              <TableRow
                v-for="server in filteredServers"
                :key="`${server.client}-${server.name}-${server.sourcePath}`"
                class="cursor-default"
              >
                <TableCell>
                  <Badge variant="secondary" class="text-xs font-medium">{{ server.client }}</Badge>
                </TableCell>
                <TableCell class="font-medium">{{ server.name }}</TableCell>
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
                    :variant="(server.status === 'configured' || server.status === 'active') ? 'default' : 'secondary'"
                    class="text-xs"
                  >
                    {{ server.status || 'unknown' }}
                  </Badge>
                </TableCell>
                <TableCell class="max-w-[220px]">
                  <code class="text-xs font-mono text-muted-foreground truncate block">
                    {{ server.sourcePath }}
                  </code>
                </TableCell>
              </TableRow>
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
