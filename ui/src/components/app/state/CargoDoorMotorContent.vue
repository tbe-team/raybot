<script setup lang="ts">
import type { Outdated } from './outdated'
import type { CargoDoorMotorState } from '@/types/cargo'
import { TriangleAlert } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { clearOutdatedTimeout, setOutdated } from './outdated'

interface Props {
  cargoDoorMotor: CargoDoorMotorState
}

const props = defineProps<Props>()
const expiredAfter = 10
const cargoDoorMotorOutdated = reactive<Outdated>({
  isOutdated: false,
  timeoutId: null,
})

watch(() => props.cargoDoorMotor.updatedAt, (newVal) => {
  setOutdated(newVal, expiredAfter, cargoDoorMotorOutdated)
}, { immediate: true })

onUnmounted(() => {
  clearOutdatedTimeout(cargoDoorMotorOutdated)
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle class="flex gap-3 items-center">
        Cargo Door Motor
        <div v-if="cargoDoorMotorOutdated.isOutdated" class="flex gap-1 items-center text-warning">
          <TriangleAlert class="size-4" />
          <span class="text-xs font-normal">{{ `Last updated ${cargoDoorMotorOutdated.timeAgo}` || 'Outdated' }}</span>
        </div>
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-2 gap-2 text-sm">
        <span class="font-medium text-muted-foreground">Direction</span>
        <span class="font-semibold">
          <Badge :variant="props.cargoDoorMotor.direction === 'CLOSE' ? 'outline' : 'default'">
            {{ props.cargoDoorMotor.direction === 'CLOSE' ? 'Close' : 'Open' }}
          </Badge>
        </span>
        <span class="font-medium text-muted-foreground">Is Running</span>
        <span>
          <Badge class="text-white" :class="props.cargoDoorMotor.isRunning ? '!bg-success' : '!bg-destructive'">
            {{ props.cargoDoorMotor.isRunning ? 'Yes' : 'No' }}
          </Badge>
        </span>
        <span class="font-medium text-muted-foreground">Speed</span>
        <span>
          {{ props.cargoDoorMotor.speed }} %
        </span>
        <span class="font-medium text-muted-foreground">Enabled</span>
        <span>
          <Badge class="text-white" :class="props.cargoDoorMotor.enabled ? '!bg-success' : '!bg-destructive'">
            {{ props.cargoDoorMotor.enabled ? 'Yes' : 'No' }}
          </Badge>
        </span>
      </div>
    </CardContent>
  </Card>
</template>
