<script setup lang="ts">
import type { CommandSource } from '@/types/command'
import type { LucideIcon } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Cloud, Smartphone } from 'lucide-vue-next'

const props = defineProps<{
  source: CommandSource
}>()

const sourceConfig: Record<CommandSource, {
  icon: LucideIcon
  label: string
}> = {
  CLOUD: {
    icon: Cloud,
    label: 'Cloud',
  },
  APP: {
    icon: Smartphone,
    label: 'App',
  },
}
</script>

<template>
  <Badge
    variant="outline"
    class="flex items-center gap-1"
  >
    <component
      :is="sourceConfig[props.source].icon"
      class="w-3 h-3"
    />
    {{ sourceConfig[props.source].label }}
  </Badge>
</template>
