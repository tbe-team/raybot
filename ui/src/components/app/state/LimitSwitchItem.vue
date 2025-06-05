<script setup lang="ts">
import type { LimitSwitchItemObject } from './LimitSwitchesTabContent.vue'
import { Power, PowerOff } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { formatDate } from '@/lib/date'

const props = defineProps<{
  switchData: LimitSwitchItemObject
}>()

function getStatusColor(): string {
  return props.switchData.pressed ? 'bg-success/10 border-success' : ''
}

function formatName(name: string): string {
  return name.charAt(0).toUpperCase() + name.slice(1).split(/(?=[A-Z0-9])/).join(' ')
}
</script>

<template>
  <div class="p-4 space-y-2 border rounded-md" :class="[getStatusColor()]">
    <div class="flex items-center justify-between">
      <span class="text-sm">{{ formatName(switchData.name) }}</span>
      <Power v-if="switchData.pressed" class="w-3 h-3 text-success" />
      <PowerOff v-else class="w-3 h-3 text-muted-foreground" />
    </div>
    <div class="flex items-center justify-between">
      <span class="text-xs text-muted-foreground">
        State:
      </span>
      <Badge :variant="switchData.pressed ? 'default' : 'outline'" class="py-0.5 px-1.5 text-xs rounded-full font-normal" :class="[switchData.pressed && '!bg-success !border-success']">
        {{ switchData.pressed ? 'Active' : 'Inactive' }}
      </Badge>
    </div>
    <div class="flex items-center justify-between">
      <span class="text-xs text-muted-foreground">
        Last updated:
      </span>
      <span class="text-xs text-muted-foreground">
        {{ formatDate(switchData.updatedAt) }}
      </span>
    </div>
  </div>
</template>
