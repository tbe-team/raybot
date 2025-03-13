<script setup lang="ts">
import type { BatteryState } from '@/types/state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface Props {
  battery: BatteryState
}
const props = defineProps<Props>()
function getBatteryColor(percent: number): string {
  if (percent < 20)
    return 'text-red-5500'
  if (percent < 40)
    return 'text-orange-500'
  if (percent < 60)
    return 'text-yellow-500'
  return 'text-green-500'
}
function getTemperatureColor(temp: number): string {
  if (temp < 15)
    return 'text-blue-500'
  if (temp < 35)
    return 'text-green-500'
  return 'text-red-500'
}
</script>

<template>
  <div>
    <Card class="h-full rounded-sm shadow-lg">
      <CardHeader>
        <CardTitle>Battery State</CardTitle>
      </CardHeader>

      <CardContent>
        <div class="mx-auto space-x-2">
          <div class="grid grid-cols-2 gap-2 text-sm">
            <p><span class="font-medium">Current: </span>{{ props.battery.current }} A</p>
            <p><span class="font-medium">Battery Level: </span><span :class="getBatteryColor(props.battery.percent)">{{ props.battery.percent }}%</span></p>
            <p><span class="font-medium">Voltage: </span>{{ props.battery.voltage }} V</p>
            <p><span class="font-medium">Health: </span>{{ props.battery.health }}%</p>
            <p><span class="font-medium">Cell Voltages: </span><span>{{ props.battery.cellVoltages.join(', ') }} V</span></p>
            <p><span class="font-medium">Fault Status: </span><span class="text-green-600">{{ props.battery.fault }}</span></p>
            <p><span class="font-medium">Temperature: </span><span :class="getTemperatureColor(props.battery.temp)">{{ props.battery.temp }}°C</span></p>
            <p><span class="font-medium">Last Updated: </span>{{ (props.battery.updatedAt).toISOString() }}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
