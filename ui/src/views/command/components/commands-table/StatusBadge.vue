<script setup lang="ts">
import type { CommandStatus } from '@/types/command'
import { Badge } from '@/components/ui/badge'

interface Props {
  status: CommandStatus
}

const props = defineProps<Props>()

const STATUS_LABELS: Record<CommandStatus, { label: string, class: string }> = {
  SUCCEEDED: { label: 'Succeeded', class: 'bg-green-500/10 text-green-500 hover:bg-green-500/20' },
  FAILED: { label: 'Failed', class: 'bg-red-500/10 text-red-500 hover:bg-red-500/20' },
  IN_PROGRESS: { label: 'In Progress', class: 'bg-yellow-500/10 text-yellow-500 hover:bg-yellow-500/20' },
}

const { label, class: className } = STATUS_LABELS[props.status]
</script>

<template>
  <Badge :class="className">
    {{ label }}
  </Badge>
</template>
