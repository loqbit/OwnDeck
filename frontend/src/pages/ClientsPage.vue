<script lang="ts" setup>
import PageHeader from '@/components/app/PageHeader.vue'
import PageSkeleton from '@/components/app/PageSkeleton.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'

import { useClients } from '@/composables/useClients'
import { useMCPServers } from '@/composables/useMCPServers'
import { useSkills } from '@/composables/useSkills'
import { useAutoRefresh } from '@/composables/useAutoRefresh'
import { toast } from 'vue-sonner'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const {
  clients,
  connectedClients,
  connectClient: rawConnect,
  disconnectClient: rawDisconnect,
} = useClients()
const { refreshServers } = useMCPServers()
const { refreshSkills } = useSkills()
const { initialLoaded } = useAutoRefresh()

async function handleConnect(clientID: string) {
  try {
    await rawConnect(clientID)
    await Promise.all([refreshServers(true), refreshSkills()])
    toast.success(t('clients.connected'), { description: t('clients.connectedMsg', { id: clientID }) })
  } catch (e) {
    toast.error(t('clients.connectionFailed'), { description: String(e) })
  }
}

async function handleDisconnect(clientID: string) {
  try {
    await rawDisconnect(clientID)
    await Promise.all([refreshServers(true), refreshSkills()])
    toast.info(t('actions.disconnect'), { description: t('clients.disconnectedMsg', { id: clientID }) })
  } catch (e) {
    toast.error(t('clients.disconnectFailed'), { description: String(e) })
  }
}
</script>

<template>
  <div class="flex flex-col flex-1 min-h-0">
    <PageHeader :title="$t('clients.title')" :description="$t('clients.description')" />

    <PageSkeleton v-if="!initialLoaded" variant="cards" />

    <div v-else class="flex-1 overflow-y-auto p-6">
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <Card v-for="client in clients" :key="client.id" class="flex flex-col">
          <CardHeader class="flex flex-row items-start justify-between space-y-0">
            <div class="space-y-1">
              <CardTitle class="text-base">{{ client.name }}</CardTitle>
              <Badge
                :variant="connectedClients.has(client.id) ? 'default' : client.detected ? 'secondary' : 'outline'"
                class="text-xs"
              >
                {{ connectedClients.has(client.id) ? $t('clients.connected') : client.status }}
              </Badge>
            </div>
            <Badge variant="outline" class="text-xs shrink-0">{{ $t('clients.readOnly') }}</Badge>
          </CardHeader>

          <CardContent class="flex-1 space-y-3">
            <div class="space-y-1">
              <p class="text-xs font-medium text-muted-foreground uppercase tracking-wider">{{ $t('clients.executable') }}</p>
              <code class="block text-xs font-mono text-foreground/80 truncate">
                {{ client.executable || $t('clients.notFound') }}
              </code>
            </div>
            <Separator />
            <div class="space-y-1">
              <p class="text-xs font-medium text-muted-foreground uppercase tracking-wider">{{ $t('clients.config') }}</p>
              <code class="block text-xs font-mono text-foreground/80 truncate">
                {{ client.configPaths.length ? client.configPaths.join(' | ') : $t('clients.noConfig') }}
              </code>
            </div>
          </CardContent>

          <CardFooter>
            <Button
              v-if="!connectedClients.has(client.id)"
              class="w-full"
              :disabled="!client.detected"
              @click="handleConnect(client.id)"
            >
              {{ $t('actions.connect') }}
            </Button>
            <Button
              v-else
              variant="outline"
              class="w-full"
              @click="handleDisconnect(client.id)"
            >
              {{ $t('actions.disconnect') }}
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  </div>
</template>
