<script setup lang="ts">
import type { LimitSwitch } from '@/types/limit-switch-state'
import { Badge } from '@/components/ui/badge'
import { formatDate } from '@/lib/date'

const props = defineProps<{
  name: string
  switchData: LimitSwitch
}>()

function getStatusColor(): string {
  return props.switchData.pressed ? 'bg-success/10 border-success' : ''
}
</script>

<template>
  <div class="p-4 space-y-2 border rounded-md" :class="[getStatusColor()]">
    <span class="text-sm">{{ props.name }}</span>
    <div class="flex items-center justify-between">
      <span class="text-xs text-muted-foreground">
        State:
      </span>
      <Badge :variant="props.switchData.pressed ? 'default' : 'outline'" class="py-0.5 px-1.5 text-xs rounded-full font-normal" :class="[props.switchData.pressed && '!bg-success !border-success']">
        {{ props.switchData.pressed ? 'Active' : 'Inactive' }}
      </Badge>
    </div>
    <div class="flex items-center justify-between">
      <span class="text-xs text-muted-foreground">
        Last updated:
      </span>
      <span class="text-xs text-muted-foreground">
        {{ formatDate(props.switchData.updatedAt) }}
      </span>
    </div>
  </div>
</template>
