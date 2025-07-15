<script setup lang="ts">
import type { Outdated } from './outdated'
import type { DriveMotorState } from '@/types/robot-state'
import { TriangleAlert } from 'lucide-vue-next'
import { reactive } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { setOutdated } from './outdated'

interface Props {
  driveMotor: DriveMotorState
}

const props = defineProps<Props>()
const expiredAfter = 10
const driveMotorOutdated = reactive<Outdated>({
  isOutdated: false,
  timeoutId: null,
})

watch(() => props.driveMotor.updatedAt, (newVal) => {
  setOutdated(newVal, expiredAfter, driveMotorOutdated)
}, { immediate: true })
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle class="flex gap-3 items-center">
        Drive Motor
        <div v-if="driveMotorOutdated.isOutdated" class="flex gap-1 items-center text-warning">
          <TriangleAlert class="size-4" />
          <span class="text-xs font-normal">{{ `Last updated ${driveMotorOutdated.timeAgo}` || 'Outdated' }}</span>
        </div>
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-2 gap-2 text-sm">
        <span class="font-medium text-muted-foreground">Status</span>
        <span>
          <Badge :variant="props.driveMotor.isRunning ? 'default' : 'outline'">
            {{ props.driveMotor.isRunning ? 'Running' : 'Idle' }}
          </Badge>
        </span>
        <span class="font-medium text-muted-foreground">Direction</span>
        <span>
          {{ props.driveMotor.direction || 'N/A' }}
        </span>
        <span class="font-medium text-muted-foreground">Speed</span>
        <span>
          {{ props.driveMotor.speed }} %
        </span>
        <span class="font-medium text-muted-foreground">Enabled</span>
        <span>
          <Badge class="text-white" :class="props.driveMotor.enabled ? '!bg-success' : '!bg-destructive'">
            {{ props.driveMotor.enabled ? 'Yes' : 'No' }}
          </Badge>
        </span>
      </div>
    </CardContent>
  </Card>
</template>
