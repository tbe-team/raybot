<script setup lang="ts">
import type { CommandStatus } from '@/types/command'
import type { LucideIcon } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Ban, CheckCircle, Clock, Loader2, XCircle } from 'lucide-vue-next'

const props = defineProps<{
  status: CommandStatus
}>()

const statusConfig: Record<CommandStatus, {
  icon: LucideIcon
  label: string
  class: string
}> = {
  QUEUED: {
    icon: Clock,
    label: 'Queued',
    class: 'bg-amber-500/10 text-amber-500 hover:bg-amber-500/20 hover:text-amber-500',
  },
  PROCESSING: {
    icon: Loader2,
    label: 'Processing',
    class: 'bg-blue-500/10 text-blue-500 hover:bg-blue-500/20 hover:text-blue-500',
  },
  SUCCEEDED: {
    icon: CheckCircle,
    label: 'Succeeded',
    class: 'bg-green-500/10 text-green-500 hover:bg-green-500/20 hover:text-green-500',
  },
  FAILED: {
    icon: XCircle,
    label: 'Failed',
    class: 'bg-red-500/10 text-red-500 hover:bg-red-500/20 hover:text-red-500',
  },
  CANCELED: {
    icon: Ban,
    label: 'Canceled',
    class: 'bg-slate-500/10 text-slate-500 hover:bg-slate-500/20 hover:text-slate-500',
  },
} as const
</script>

<template>
  <Badge
    variant="secondary"
    class="flex items-center gap-1" :class="[
      statusConfig[props.status].class,
    ]"
  >
    <component
      :is="statusConfig[props.status].icon"
      class="w-3 h-3"
      :class="{ 'animate-spin': props.status === 'PROCESSING' }"
    />
    {{ statusConfig[props.status].label }}
  </Badge>
</template>
