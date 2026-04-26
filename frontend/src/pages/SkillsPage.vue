<script lang="ts" setup>
import PageHeader from '@/components/app/PageHeader.vue'
import PageSkeleton from '@/components/app/PageSkeleton.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { useClients } from '@/composables/useClients'
import { useSkills } from '@/composables/useSkills'
import { useAutoRefresh } from '@/composables/useAutoRefresh'

const { connectedCount } = useClients()
const { skillFilter, visibleSkills, hasVisibleSkills, skills } = useSkills()
const { initialLoaded } = useAutoRefresh()
</script>

<template>
  <div class="flex flex-col flex-1 min-h-0">
    <PageHeader :title="$t('skills.title')" :description="$t('skills.description')" />

    <PageSkeleton v-if="!initialLoaded" variant="list" />

    <div v-else class="flex-1 overflow-y-auto p-6 space-y-4">
      <Tabs
        :model-value="skillFilter"
        @update:model-value="(v) => { if (v === 'managed' || v === 'all' || v === 'system') skillFilter = v }"
      >
        <TabsList>
          <TabsTrigger value="managed">{{ $t('skills.user') }}</TabsTrigger>
          <TabsTrigger value="system">{{ $t('skills.builtIn') }}</TabsTrigger>
          <TabsTrigger value="all">{{ $t('skills.all') }}</TabsTrigger>
        </TabsList>
      </Tabs>

      <div v-if="hasVisibleSkills" class="grid gap-3">
          <Card v-for="skill in visibleSkills" :key="`${skill.clientID}-${skill.sourcePath}`">
            <CardContent class="py-4 space-y-2">
              <div class="flex items-center justify-between gap-3">
                <span class="font-medium text-sm truncate">{{ skill.name }}</span>
                <div class="flex items-center gap-1.5 shrink-0">
                  <Badge variant="secondary" class="text-xs">{{ skill.client }}</Badge>
                  <Badge variant="outline" class="text-xs">{{ skill.scope }}</Badge>
                </div>
              </div>
              <p class="text-xs text-muted-foreground line-clamp-2">
                {{ skill.description || $t('skills.noDescription') }}
              </p>
              <code class="block text-xs font-mono text-muted-foreground/70 truncate">
                {{ skill.sourcePath }}
              </code>
            </CardContent>
          </Card>
      </div>

      <div v-else class="flex flex-col items-center justify-center py-16 text-center">
        <p class="font-medium text-sm">
          {{ skills.length > 0 ? $t('skills.noUserSkills') : $t('skills.noSkillsYet') }}
        </p>
        <p class="text-sm text-muted-foreground mt-1 max-w-sm">
          {{
            skills.length > 0
              ? $t('skills.switchHint')
              : connectedCount > 0
                ? $t('skills.noSkillsConfigured')
                : $t('skills.connectHint')
          }}
        </p>
      </div>
    </div>
  </div>
</template>
