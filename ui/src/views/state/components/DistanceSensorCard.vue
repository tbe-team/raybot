<script setup lang="ts">
import type { DistanceSensorState } from '@/types/state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface Props {
  distanceSensor: DistanceSensorState
}

const props = defineProps<Props>()
</script>

<template>
  <div>
    <Card class="h-full rounded-sm shadow-lg">
      <CardHeader>
        <CardTitle>Distance Sensor State</CardTitle>
      </CardHeader>

      <CardContent>
        <div class="mx-auto space-x-2">
          <div class="grid grid-rows-2 gap-2 text-sm">
            <div class="grid grid-cols-3 gap-2">
              <p><span class="font-medium">Front: </span>{{ props.distanceSensor.frontDistance }} cm</p>
              <p><span class="font-medium">Back: </span>{{ props.distanceSensor.backDistance }} cm</p>
              <p><span class="font-medium">Down: </span>{{ props.distanceSensor.downDistance }} cm</p>
            </div>
            <p><span class="font-medium">Last Updated: </span>{{ (props.distanceSensor.updatedAt).toISOString() }}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
