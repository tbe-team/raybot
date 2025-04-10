<script setup lang="ts">
import type { CommandSource } from '@/types/command'
import { Badge } from '@/components/ui/badge'

interface Props {
  source: CommandSource
}

const props = defineProps<Props>()

const SOURCE_LABELS: Record<CommandSource, { label: string, class: string }> = {
  CLOUD: { label: 'Cloud', class: 'bg-purple-500/10 text-purple-500 hover:bg-purple-500/20' },
  APP: { label: 'Application', class: 'bg-blue-500/10 text-blue-500 hover:bg-blue-500/20' },
}

const { label, class: className } = SOURCE_LABELS[props.source]
</script>

<template>
  <Badge :class="className">
    {{ label }}
  </Badge>
</template>
