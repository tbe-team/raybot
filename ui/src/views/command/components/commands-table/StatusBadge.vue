<script setup lang="ts">
import type { CommandStatus } from '@/types/command'
import type { Component } from 'vue'
import { Badge } from '@/components/ui/badge'
import { AlertCircle, CheckCircle, Clock, Loader, XCircle } from 'lucide-vue-next'

interface Props {
  status: CommandStatus
}

const props = defineProps<Props>()

const STATUS_LABELS: Record<CommandStatus, { label: string, class: string, icon: Component }> = {
  QUEUED: { label: 'Queued', class: 'bg-yellow-500/10 text-yellow-500 hover:bg-yellow-500/20', icon: Clock },
  PROCESSING: { label: 'Processing', class: 'bg-blue-500/10 text-blue-500 hover:bg-blue-500/20', icon: Loader },
  SUCCEEDED: { label: 'Succeeded', class: 'bg-green-500/10 text-green-500 hover:bg-green-500/20', icon: CheckCircle },
  FAILED: { label: 'Failed', class: 'bg-red-500/10 text-red-500 hover:bg-red-500/20', icon: AlertCircle },
  CANCELED: { label: 'Canceled', class: 'bg-gray-500/10 text-gray-500 hover:bg-gray-500/20', icon: XCircle },
}

const { icon, label, class: className } = STATUS_LABELS[props.status]
</script>

<template>
  <Badge :class="`${className}`">
    <component :is="icon" class="w-4 h-4 mr-2" />
    {{ label }}
  </Badge>
</template>
