<script setup lang="ts">
import type { Outdated } from './outdated'
import type { DistanceSensorState } from '@/types/robot-state'
import { TriangleAlert } from 'lucide-vue-next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { clearOutdatedTimeout, setOutdated } from './outdated'

interface Props {
  distanceSensor: DistanceSensorState
}

const props = defineProps<Props>()
const expiredAfter = 10
const distanceSensorOutdated = reactive<Outdated>({
  isOutdated: false,
  timeoutId: null,
})

watch(() => props.distanceSensor.updatedAt, (newVal) => {
  setOutdated(newVal, expiredAfter, distanceSensorOutdated)
}, { immediate: true })

onUnmounted(() => {
  clearOutdatedTimeout(distanceSensorOutdated)
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle class="flex gap-3 items-center">
        Distance Sensor
        <div v-if="distanceSensorOutdated.isOutdated" class="flex gap-1 items-center text-warning">
          <TriangleAlert class="size-4" />
          <span class="text-xs font-normal">{{ `Last updated ${distanceSensorOutdated.timeAgo}` || 'Outdated' }}</span>
        </div>
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-2 gap-2 text-sm">
        <span class="font-medium text-muted-foreground">Front Distance</span>
        <span>
          {{ props.distanceSensor.frontDistance }} cm
        </span>
        <span class="font-medium text-muted-foreground">Back Distance</span>
        <span>
          {{ props.distanceSensor.backDistance }} cm
        </span>
        <span class="font-medium text-muted-foreground">Down Distance</span>
        <span>
          {{ props.distanceSensor.downDistance }} cm
        </span>
      </div>
    </CardContent>
  </Card>
</template>
