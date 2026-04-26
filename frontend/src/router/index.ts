import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      redirect: '/overview',
    },
    {
      path: '/overview',
      name: 'overview',
      component: () => import('@/pages/OverviewPage.vue'),
    },
    {
      path: '/clients',
      name: 'clients',
      component: () => import('@/pages/ClientsPage.vue'),
    },
    {
      path: '/mcp-servers',
      name: 'mcp-servers',
      component: () => import('@/pages/MCPServersPage.vue'),
    },
    {
      path: '/skills',
      name: 'skills',
      component: () => import('@/pages/SkillsPage.vue'),
    },
    {
      path: '/agents',
      name: 'agents',
      component: () => import('@/pages/AgentsPage.vue'),
    },
    {
      path: '/profiles',
      name: 'profiles',
      component: () => import('@/pages/ProfilesPage.vue'),
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('@/pages/SettingsPage.vue'),
    },
  ],
})

export default router
