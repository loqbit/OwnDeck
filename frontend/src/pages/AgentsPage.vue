<script lang="ts" setup>
import PageHeader from '@/components/app/PageHeader.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Bot } from 'lucide-vue-next'
</script>

<template>
  <div class="flex flex-col flex-1 min-h-0">
    <PageHeader :title="$t('agents.title')" :description="$t('agents.description')" />

    <div class="flex-1 overflow-y-auto p-6">
      <Card class="border-dashed">
        <CardContent class="flex flex-col items-center justify-center py-16 text-center">
          <div class="rounded-full bg-muted p-4 mb-4">
            <Bot class="size-8 text-muted-foreground" />
          </div>
          <p class="font-medium text-sm">{{ $t('agents.noAgents') }}</p>
          <p class="text-sm text-muted-foreground mt-1 max-w-sm">{{ $t('agents.noAgentsDesc') }}</p>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
