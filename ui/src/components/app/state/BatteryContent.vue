<script setup lang="ts">
import type { Outdated } from '@/components/app/state/outdated'
import type { BatteryState, ChargeState, DischargeState } from '@/types/robot-state'
import { TriangleAlert } from 'lucide-vue-next'
import { clearOutdatedTimeout, setOutdated } from '@/components/app/state/outdated'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface Props {
  battery: BatteryState
  charge: ChargeState
  discharge: DischargeState
}

const props = defineProps<Props>()
const expiredAfter = 10
const batteryOutdated = reactive<Outdated>({
  isOutdated: false,
  timeoutId: null,
})

function getBatteryColor(percent: number): string {
  if (percent < 20)
    return 'text-destructive'
  if (percent < 60)
    return 'text-warning'
  return 'text-success'
}

function convertMvToV(mv: number): number {
  return Number.parseFloat((mv / 1000).toFixed(2))
}

watch(() => props.battery.updatedAt, (newVal) => {
  setOutdated(newVal, expiredAfter, batteryOutdated)
}, { immediate: true })

onUnmounted(() => {
  clearOutdatedTimeout(batteryOutdated)
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle class="flex gap-3 items-center">
        Battery
        <div v-if="batteryOutdated.isOutdated" class="flex gap-1 items-center text-warning">
          <TriangleAlert class="size-4" />
          <span class="text-xs font-normal">{{ `Last updated ${batteryOutdated.timeAgo}` || 'Outdated' }}</span>
        </div>
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="gap-2 space-y-2 text-sm columns-1 md:columns-2">
        <div class="grid grid-cols-2 gap-2">
          <span class="font-medium text-muted-foreground">Battery SOC</span>
          <span :class="getBatteryColor(props.battery.percent)">
            {{ props.battery.percent }} %
          </span>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <span class="font-medium text-muted-foreground">Health</span>
          <span>
            {{ props.battery.health }} %
          </span>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <span class="font-medium text-muted-foreground">Temperature</span>
          <span>
            {{ props.battery.temp }} °C
          </span>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <span class="font-medium text-muted-foreground">Voltage</span>
          <span>
            {{ props.battery.voltage }} V
          </span>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <span class="font-medium text-muted-foreground">Current</span>
          <span>
            {{ props.battery.current }} mA
          </span>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <span class="font-medium text-muted-foreground">Cell Voltages</span>
          <span>
            {{ props.battery.cellVoltages.map(v => convertMvToV(v)).join(', ') }} V
          </span>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <span class="font-medium text-muted-foreground">Fault Code</span>
          <span class="flex gap-2 items-center" :class="props.battery.fault !== 0 ? 'text-destructive' : ''">
            <TriangleAlert v-if="props.battery.fault !== 0" class="size-4" />
            {{ props.battery.fault }}
          </span>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <span class="font-medium text-muted-foreground">Charging</span>
          <span>
            <Badge class="text-white" :class="props.charge.enabled ? '!bg-success ' : '!bg-destructive'">
              {{ props.charge.enabled ? 'Yes' : 'No' }}
            </Badge>
          </span>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <span class="font-medium text-muted-foreground">Charging Current Limit</span>
          <span>
            {{ props.charge.currentLimit }} mA
          </span>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <span class="font-medium text-muted-foreground">Discharging</span>
          <span>
            <Badge class="text-white" :class="props.discharge.enabled ? '!bg-success' : '!bg-destructive'">
              {{ props.discharge.enabled ? 'Yes' : 'No' }}
            </Badge>
          </span>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <span class="font-medium text-muted-foreground">Discharging Current Limit</span>
          <span>
            {{ props.discharge.currentLimit }} mA
          </span>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
