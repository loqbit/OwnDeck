<script lang="ts" setup>
import AppSidebar from '@/components/app/AppSidebar.vue'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
</script>

<template>
  <SidebarProvider>
    <AppSidebar />
    <SidebarInset>
      <div class="flex h-svh flex-col">
        <slot />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
