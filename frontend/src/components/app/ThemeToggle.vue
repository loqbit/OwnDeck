<script lang="ts" setup>
import { useDark, useToggle } from '@vueuse/core'
import { Button } from '@/components/ui/button'
import { Sun, Moon } from 'lucide-vue-next'

const isDark = useDark({
  storageKey: 'owndeck-theme',
})
const toggleDark = useToggle(isDark)
</script>

<template>
  <Button variant="ghost" size="icon" class="size-7" @click="toggleDark()">
    <Sun v-if="isDark" class="size-4" />
    <Moon v-else class="size-4" />
    <span class="sr-only">Toggle theme</span>
  </Button>
</template>
