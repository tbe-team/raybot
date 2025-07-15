<script setup lang="ts">
import type { Outdated } from './outdated'
import type { LiftMotorState } from '@/types/robot-state'
import { TriangleAlert } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { clearOutdatedTimeout, setOutdated } from './outdated'

interface Props {
  liftMotor: LiftMotorState
}

const props = defineProps<Props>()
const expiredAfter = 10
const liftMotorOutdated = reactive<Outdated>({
  isOutdated: false,
  timeoutId: null,
})

watch(() => props.liftMotor.updatedAt, (newVal) => {
  setOutdated(newVal, expiredAfter, liftMotorOutdated)
}, { immediate: true })

onUnmounted(() => {
  clearOutdatedTimeout(liftMotorOutdated)
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle class="flex gap-3 items-center">
        Lift Motor
        <div v-if="liftMotorOutdated.isOutdated" class="flex gap-1 items-center text-warning">
          <TriangleAlert class="size-4" />
          <span class="text-xs font-normal">{{ `Last updated ${liftMotorOutdated.timeAgo}` || 'Outdated' }}</span>
        </div>
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-2 gap-2 text-sm">
        <span class="font-medium text-muted-foreground">Status</span>
        <span>
          <Badge :variant="props.liftMotor.isRunning ? 'default' : 'outline'">
            {{ props.liftMotor.isRunning ? 'Running' : 'Idle' }}
          </Badge>
        </span>
        <span class="font-medium text-muted-foreground">Current Position</span>
        <span>
          {{ props.liftMotor.currentPosition }} cm
        </span>
        <span class="font-medium text-muted-foreground">Target Position</span>
        <span>
          {{ props.liftMotor.targetPosition }} cm
        </span>
        <span class="font-medium text-muted-foreground">Enabled</span>
        <span>
          <Badge class="text-white" :class="props.liftMotor.enabled ? '!bg-success' : '!bg-destructive'">
            {{ props.liftMotor.enabled ? 'Yes' : 'No' }}
          </Badge>
        </span>
      </div>
    </CardContent>
  </Card>
</template>
