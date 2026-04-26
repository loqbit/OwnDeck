<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import PageHeader from '@/components/app/PageHeader.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { ConfigPath } from '../../wailsjs/go/main/App'
import { Settings, Globe } from 'lucide-vue-next'
import { setLocale, supportedLocales } from '@/i18n'

const { locale } = useI18n()
const configPath = ref('')

onMounted(async () => {
  try {
    configPath.value = await ConfigPath()
  } catch {
    configPath.value = 'N/A'
  }
})

function handleLocaleChange(val: any) {
  if (typeof val === 'string') setLocale(val)
}
</script>

<template>
  <div class="flex flex-col flex-1 min-h-0">
    <PageHeader :title="$t('settings.title')" :description="$t('settings.description')" />
    <div class="flex-1 overflow-y-auto p-6 space-y-4">
      <Card>
        <CardHeader>
          <CardTitle class="text-sm font-medium flex items-center gap-2">
            <Settings class="size-4" />
            {{ $t('settings.application') }}
          </CardTitle>
        </CardHeader>
        <CardContent class="space-y-3">
          <div class="space-y-1">
            <p class="text-xs font-medium text-muted-foreground uppercase tracking-wider">{{ $t('settings.configPath') }}</p>
            <code class="block text-sm font-mono text-foreground/80">{{ configPath || $t('settings.loading') }}</code>
          </div>
          <div class="space-y-1">
            <p class="text-xs font-medium text-muted-foreground uppercase tracking-wider">{{ $t('settings.version') }}</p>
            <code class="block text-sm font-mono text-foreground/80">{{ $t('app.version') }}</code>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="text-sm font-medium flex items-center gap-2">
            <Globe class="size-4" />
            {{ $t('settings.language') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Select :model-value="locale" @update:model-value="handleLocaleChange">
            <SelectTrigger class="w-[200px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="loc in supportedLocales" :key="loc.value" :value="loc.value">
                {{ loc.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
