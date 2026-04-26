import { computed, ref } from 'vue'
import { DiscoverSkillsForClients } from '../../wailsjs/go/main/App'
import type { discovery } from '../../wailsjs/go/models'
import { useClients } from './useClients'

const skills = ref<discovery.SkillAsset[]>([])
const skillFilter = ref<'managed' | 'all' | 'system'>('managed')

const managedSkills = computed(() =>
  skills.value.filter(s => s.scope === 'local' || s.scope === 'user'),
)
const systemSkills = computed(() =>
  skills.value.filter(s => s.scope === 'system' || s.scope === 'plugin'),
)
const visibleSkills = computed(() => {
  if (skillFilter.value === 'all') return skills.value
  if (skillFilter.value === 'system') return systemSkills.value
  return managedSkills.value
})
const hasVisibleSkills = computed(() => visibleSkills.value.length > 0)

async function refreshSkills() {
  const { connectedClientIDs } = useClients()

  if (connectedClientIDs.value.length === 0) {
    if (skills.value.length > 0) skills.value.splice(0)
    return
  }
  const fresh = (await DiscoverSkillsForClients(connectedClientIDs.value)) ?? []
  // Only replace when content actually changed to avoid re-render flicker
  if (!skillsEqual(skills.value, fresh)) {
    skills.value = fresh
  }
}

function skillsEqual(a: discovery.SkillAsset[], b: discovery.SkillAsset[]): boolean {
  if (a.length !== b.length) return false
  for (let i = 0; i < a.length; i++) {
    if (a[i].name !== b[i].name || a[i].clientID !== b[i].clientID || a[i].sourcePath !== b[i].sourcePath) return false
  }
  return true
}

export function useSkills() {
  return {
    skills,
    skillFilter,
    managedSkills,
    systemSkills,
    visibleSkills,
    hasVisibleSkills,
    refreshSkills,
  }
}
