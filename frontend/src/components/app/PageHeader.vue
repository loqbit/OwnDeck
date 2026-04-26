<script lang="ts" setup>
import { SidebarTrigger } from '@/components/ui/sidebar'
import { Separator } from '@/components/ui/separator'

defineProps<{
  title: string
  description?: string
}>()
</script>

<template>
  <header class="flex h-14 shrink-0 items-center gap-2 border-b border-border overflow-hidden">
    <div class="flex items-center gap-2 px-4 shrink-0">
      <SidebarTrigger class="-ml-1" />
      <Separator orientation="vertical" class="mr-2 h-4" />
      <div>
        <h1 class="text-base font-semibold leading-tight">{{ title }}</h1>
        <p v-if="description" class="text-xs text-muted-foreground">{{ description }}</p>
      </div>
    </div>
    <div class="ml-auto flex items-center gap-2 px-4 shrink-0">
      <slot name="actions" />
    </div>
  </header>
</template>
