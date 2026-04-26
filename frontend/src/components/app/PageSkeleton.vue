<script lang="ts" setup>
import { Skeleton } from '@/components/ui/skeleton'

defineProps<{
  variant?: 'cards' | 'table' | 'list'
}>()
</script>

<template>
  <!-- Cards skeleton (Overview, Clients) -->
  <div v-if="variant === 'cards'" class="p-6 space-y-6">
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <div v-for="i in 4" :key="i" class="rounded-lg border p-6 space-y-3">
        <div class="flex items-center justify-between">
          <Skeleton class="h-4 w-24" />
          <Skeleton class="h-4 w-4 rounded" />
        </div>
        <Skeleton class="h-8 w-16" />
      </div>
    </div>
    <div class="grid gap-4 md:grid-cols-2">
      <div v-for="i in 2" :key="i" class="rounded-lg border p-6 space-y-3">
        <Skeleton class="h-4 w-32" />
        <div class="space-y-2">
          <Skeleton v-for="j in 3" :key="j" class="h-5 w-full" />
        </div>
      </div>
    </div>
  </div>

  <!-- Table skeleton (MCP Servers) -->
  <div v-else-if="variant === 'table'" class="p-6">
    <div class="rounded-lg border">
      <div class="border-b px-4 py-3 flex gap-4">
        <Skeleton v-for="i in 6" :key="i" class="h-4 flex-1" />
      </div>
      <div v-for="i in 5" :key="i" class="px-4 py-3 flex gap-4 border-b last:border-b-0">
        <Skeleton class="h-5 w-20 rounded-full" />
        <Skeleton class="h-5 flex-1" />
        <Skeleton class="h-5 w-16 rounded-full" />
        <Skeleton class="h-5 flex-[2]" />
        <Skeleton class="h-5 w-20 rounded-full" />
        <Skeleton class="h-5 w-16" />
      </div>
    </div>
  </div>

  <!-- List skeleton (Skills) -->
  <div v-else class="p-6 space-y-3">
    <Skeleton class="h-9 w-48 rounded-lg" />
    <div v-for="i in 4" :key="i" class="rounded-lg border p-4 space-y-2">
      <div class="flex items-center justify-between">
        <Skeleton class="h-4 w-40" />
        <div class="flex gap-1.5">
          <Skeleton class="h-5 w-16 rounded-full" />
          <Skeleton class="h-5 w-12 rounded-full" />
        </div>
      </div>
      <Skeleton class="h-3 w-3/4" />
      <Skeleton class="h-3 w-1/2" />
    </div>
  </div>
</template>
