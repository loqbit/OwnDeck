<script lang="ts" setup>
import { ref } from 'vue'
import PageHeader from '@/components/app/PageHeader.vue'
import PageSkeleton from '@/components/app/PageSkeleton.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  MonitorSmartphone,
  Server,
  Sparkles,
  Bot,
  AlertTriangle,
  Radar,
  Loader2,
} from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useClients } from '@/composables/useClients'
import { useMCPServers } from '@/composables/useMCPServers'
import { useSkills } from '@/composables/useSkills'
import { useAutoRefresh } from '@/composables/useAutoRefresh'
import { RescanAgents } from '../../wailsjs/go/main/App'
import { toast } from 'vue-sonner'

const { t } = useI18n()
const { connectedCount } = useClients()
const { servers } = useMCPServers()
const { skills } = useSkills()
const { initialLoaded } = useAutoRefresh()

const agentCount = 0

const metrics = [
  { key: 'overview.connectedClients', icon: MonitorSmartphone, getValue: () => connectedCount.value },
  { key: 'overview.mcpServers', icon: Server, getValue: () => servers.value.length },
  { key: 'overview.skillsDiscovered', icon: Sparkles, getValue: () => skills.value.length },
  { key: 'overview.agents', icon: Bot, getValue: () => agentCount },
]

const isScanning = ref(false)

async function handleRescan() {
  isScanning.value = true
  try {
    const agents = await RescanAgents()
    toast.success(t('actions.rescan'), {
      description: `${agents.length} agents found`,
    })
  } catch (err) {
    toast.error('Rescan failed', { description: String(err) })
  } finally {
    isScanning.value = false
  }
}
</script>

<template>
  <div class="flex flex-col flex-1 min-h-0">
    <PageHeader :title="$t('overview.title')" :description="$t('overview.description')">
      <template #actions>
        <Button variant="outline" size="sm" :disabled="isScanning" @click="handleRescan">
          <Radar v-if="!isScanning" class="mr-1.5 size-4" />
          <Loader2 v-else class="mr-1.5 size-4 animate-spin" />
          {{ isScanning ? $t('actions.rescanning') : $t('actions.rescan') }}
        </Button>
      </template>
    </PageHeader>

    <PageSkeleton v-if="!initialLoaded" variant="cards" />

    <div v-else class="flex-1 overflow-y-auto p-6 space-y-6">
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Card v-for="metric in metrics" :key="metric.key">
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium text-muted-foreground">
              {{ t(metric.key) }}
            </CardTitle>
            <component :is="metric.icon" class="size-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div class="text-2xl font-bold">{{ metric.getValue() }}</div>
          </CardContent>
        </Card>
      </div>

      <Card v-if="connectedCount === 0" class="border-dashed">
        <CardContent class="flex items-center gap-3 py-4">
          <AlertTriangle class="size-5 text-amber-500" />
          <div class="text-sm">
            <p class="font-medium">{{ $t('overview.noClients') }}</p>
            <p class="text-muted-foreground">{{ $t('overview.noClientsDesc') }}</p>
          </div>
          <Button variant="outline" size="sm" class="ml-auto" as-child>
            <router-link to="/clients">{{ $t('actions.connectClients') }}</router-link>
          </Button>
        </CardContent>
      </Card>

      <div v-if="connectedCount > 0" class="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle class="text-sm font-medium">{{ $t('overview.recentServers') }}</CardTitle>
          </CardHeader>
          <CardContent>
            <div v-if="servers.length > 0" class="space-y-2">
              <div
                v-for="server in servers.slice(0, 5)"
                :key="`${server.client}-${server.name}`"
                class="flex items-center justify-between text-sm"
              >
                <div class="flex items-center gap-2 min-w-0">
                  <span class="font-medium truncate">{{ server.name }}</span>
                  <Badge variant="secondary" class="text-xs shrink-0">{{ server.client }}</Badge>
                </div>
                <Badge variant="outline" class="text-xs shrink-0">
                  {{ server.transport || 'stdio' }}
                </Badge>
              </div>
              <router-link
                v-if="servers.length > 5"
                to="/mcp-servers"
                class="block text-xs text-muted-foreground hover:text-foreground pt-2"
              >
                {{ $t('overview.viewAllServers', { count: servers.length }) }}
              </router-link>
            </div>
            <p v-else class="text-sm text-muted-foreground">{{ $t('overview.noServers') }}</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle class="text-sm font-medium">{{ $t('overview.skillsSummary') }}</CardTitle>
          </CardHeader>
          <CardContent>
            <div v-if="skills.length > 0" class="space-y-2">
              <div
                v-for="skill in skills.slice(0, 5)"
                :key="`${skill.clientID}-${skill.sourcePath}`"
                class="flex items-center justify-between text-sm"
              >
                <span class="font-medium truncate">{{ skill.name }}</span>
                <Badge variant="secondary" class="text-xs shrink-0">{{ skill.scope }}</Badge>
              </div>
              <router-link
                v-if="skills.length > 5"
                to="/skills"
                class="block text-xs text-muted-foreground hover:text-foreground pt-2"
              >
                {{ $t('overview.viewAllSkills', { count: skills.length }) }}
              </router-link>
            </div>
            <p v-else class="text-sm text-muted-foreground">{{ $t('overview.noSkills') }}</p>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>
