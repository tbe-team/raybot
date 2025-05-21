<script setup lang="ts">
import type { BatteryState } from '@/types/robot-state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
import { BatteryFull, HeartPulse, Zap } from 'lucide-vue-next'

const props = defineProps<{
  battery: BatteryState
}>()

function getBatteryColor(percent: number): string {
  if (percent < 20)
    return 'text-destructive'
  if (percent < 60)
    return 'text-warning'
  return 'text-success'
}

function getProgressVariant(percent: number): string {
  if (percent < 20)
    return 'destructive'
  if (percent < 60)
    return 'warning'
  return 'success'
}
</script>

<template>
  <Card>
    <CardHeader class="pb-2">
      <CardTitle class="flex items-center gap-2 text-base font-semibold">
        <BatteryFull class="w-5 h-5 text-primary" />
        Battery Status
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="space-y-4">
        <!-- Battery Percentage + Progress -->
        <div>
          <div class="text-2xl font-bold" :class="getBatteryColor(props.battery.percent)">
            {{ props.battery.percent }}%
          </div>
          <Progress
            :model-value="props.battery.percent" :variant="getProgressVariant(props.battery.percent)"
            class="mt-2"
          />
        </div>

        <!-- Voltage and Health Info -->
        <div class="grid grid-cols-2 gap-4">
          <div class="flex items-center gap-2">
            <Zap class="w-5 h-5 text-warning" />
            <div>
              <p class="text-sm font-medium">
                Voltage
              </p>
              <p class="text-xl font-bold">
                {{ battery.voltage }} V
              </p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <HeartPulse class="w-5 h-5 text-destructive" />
            <div>
              <p class="text-sm font-medium">
                Health
              </p>
              <p class="text-xl font-bold">
                {{ battery.health }}%
              </p>
            </div>
          </div>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
