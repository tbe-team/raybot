<script setup lang="ts">
import type { BatteryState } from '@/types/robot-state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
import { BatteryFull } from 'lucide-vue-next'

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
      <CardTitle class="flex items-center gap-2 text-sm font-medium">
        <BatteryFull class="w-8 h-8" />
        Battery
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="text-2xl font-bold" :class="getBatteryColor(props.battery.percent)">
        {{ props.battery.percent }}%
      </div>
      <Progress
        :model-value="props.battery.percent"
        :variant="getProgressVariant(props.battery.percent)"
        class="mt-2"
      />
      <p class="mt-2 text-xs text-muted-foreground">
        Health: {{ props.battery.health }}%
      </p>
    </CardContent>
  </Card>
</template>
