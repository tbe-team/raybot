<script setup lang="ts">
import type { CommandSource } from '@/types/command'
import { Badge } from '@/components/ui/badge'

interface Props {
  source: CommandSource
}
const props = defineProps<Props>()

const SOURCE_LABELS: Record<CommandSource, string> = {
  CLOUD: 'Cloud',
  APP: 'Application',
}

const SOURCE_CLASSES: Record<CommandSource, string> = {
  CLOUD: 'bg-purple-500/10 text-purple-500 hover:bg-purple-500/20',
  APP: 'bg-blue-500/10 text-blue-500 hover:bg-blue-500/20',
}

const label = SOURCE_LABELS[props.source]
const className = SOURCE_CLASSES[props.source]
</script>

<template>
  <Badge :class="className">
    {{ label }}
  </Badge>
</template>
