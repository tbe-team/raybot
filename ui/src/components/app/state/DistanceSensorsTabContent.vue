<script setup lang="ts">
import type { DistanceSensorState } from '@/types/robot-state'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
import { formatDate } from '@/lib/date'
import { ArrowDown, ArrowLeft, ArrowRight } from 'lucide-vue-next'

const props = defineProps<{
  distanceSensor: DistanceSensorState
}>()
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Distance Sensors</CardTitle>
      <CardDescription>Proximity sensor readings</CardDescription>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-1 gap-6 md:grid-cols-3">
        <div class="space-y-2">
          <div class="flex items-center gap-2">
            <ArrowLeft class="w-5 h-5 text-muted-foreground" />
            <h3 class="font-medium">
              Front Distance
            </h3>
          </div>
          <div class="text-3xl font-bold">
            {{ props.distanceSensor.frontDistance }} cm
          </div>
          <Progress
            :value="Math.min(100, (props.distanceSensor.frontDistance / 200) * 100)"
            class="h-2"
          />
        </div>

        <div class="space-y-2">
          <div class="flex items-center gap-2">
            <ArrowRight class="w-5 h-5 text-muted-foreground" />
            <h3 class="font-medium">
              Back Distance
            </h3>
          </div>
          <div class="text-3xl font-bold">
            {{ props.distanceSensor.backDistance }} cm
          </div>
          <Progress
            :value="Math.min(100, (props.distanceSensor.backDistance / 200) * 100)"
            class="h-2"
          />
        </div>

        <div class="space-y-2">
          <div class="flex items-center gap-2">
            <ArrowDown class="w-5 h-5 text-muted-foreground" />
            <h3 class="font-medium">
              Down Distance
            </h3>
          </div>
          <div class="text-3xl font-bold">
            {{ props.distanceSensor.downDistance }} cm
          </div>
          <Progress
            :value="Math.min(100, (props.distanceSensor.downDistance / 50) * 100)"
            class="h-2"
          />
        </div>
      </div>

      <div class="mt-6 text-xs text-muted-foreground">
        Last updated: {{ formatDate(props.distanceSensor.updatedAt) }}
      </div>
    </CardContent>
  </Card>
</template>
