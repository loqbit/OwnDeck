<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  SidebarSeparator,
} from '@/components/ui/sidebar'
import {
  LayoutDashboard,
  MonitorSmartphone,
  Server,
  Sparkles,
  Bot,
  Layers,
  Settings,
} from 'lucide-vue-next'
import ThemeToggle from '@/components/app/ThemeToggle.vue'

const route = useRoute()
const { t } = useI18n()

const navMain = [
  { key: 'nav.overview', url: '/overview', icon: LayoutDashboard },
  { key: 'nav.clients', url: '/clients', icon: MonitorSmartphone },
  { key: 'nav.mcpServers', url: '/mcp-servers', icon: Server },
  { key: 'nav.skills', url: '/skills', icon: Sparkles },
  { key: 'nav.agents', url: '/agents', icon: Bot },
]

const navSecondary = [
  { key: 'nav.profiles', url: '/profiles', icon: Layers },
  { key: 'nav.settings', url: '/settings', icon: Settings },
]

function isActive(url: string): boolean {
  return route.path === url
}
</script>

<template>
  <Sidebar collapsible="icon">
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child>
            <router-link to="/overview">
              <div class="flex aspect-square size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                <Layers class="size-4" />
              </div>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">{{ $t('app.name') }}</span>
                <span class="truncate text-xs text-muted-foreground">{{ $t('app.tagline') }}</span>
              </div>
            </router-link>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>

    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupLabel>{{ $t('nav.navigation') }}</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in navMain" :key="item.key">
              <SidebarMenuButton
                as-child
                :is-active="isActive(item.url)"
                :tooltip="t(item.key)"
              >
                <router-link :to="item.url">
                  <component :is="item.icon" />
                  <span>{{ t(item.key) }}</span>
                </router-link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>

      <SidebarSeparator />

      <SidebarGroup>
        <SidebarGroupLabel>{{ $t('nav.configuration') }}</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in navSecondary" :key="item.key">
              <SidebarMenuButton
                as-child
                :is-active="isActive(item.url)"
                :tooltip="t(item.key)"
              >
                <router-link :to="item.url">
                  <component :is="item.icon" />
                  <span>{{ t(item.key) }}</span>
                </router-link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>

    <SidebarFooter>
      <SidebarMenu>
        <SidebarMenuItem>
          <div class="flex items-center justify-between px-2">
            <span class="text-xs text-muted-foreground">{{ $t('app.version') }}</span>
            <ThemeToggle />
          </div>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarFooter>

    <SidebarRail />
  </Sidebar>
</template>
