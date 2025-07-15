<script setup lang="ts">
import type { LimitSwitch } from '@/types/limit-switch-state'
import { Badge } from '@/components/ui/badge'

export interface Props {
  name: string
  switchData: LimitSwitch
}

const props = defineProps<Props>()
</script>

<template>
  <span class="font-medium text-muted-foreground">{{ props.name }}</span>
  <span>
    <Badge :variant="props.switchData.pressed ? 'default' : 'outline'">
      {{ props.switchData.pressed ? 'Active' : 'Inactive' }}
    </Badge>
  </span>
</template>
