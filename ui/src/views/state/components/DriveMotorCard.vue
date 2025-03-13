<script setup lang="ts">
import type { DriveMotorState } from '@/types/state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface Props {
  driveMotor: DriveMotorState
}

const props = defineProps<Props>()
</script>

<template>
  <div>
    <Card class="h-full rounded-sm shadow-lg ">
      <CardHeader>
        <CardTitle>Drive Motor State</CardTitle>
      </CardHeader>

      <CardContent>
        <div class="mx-auto space-x-2">
          <div class="grid grid-rows-1 gap-2 text-sm">
            <div class="grid grid-cols-2 gap-2">
              <p><span class="font-medium">Direction: </span>{{ props.driveMotor.direction }}</p>
              <p><span class="font-medium">Speed: </span>{{ props.driveMotor.speed }} m/s</p>
              <p>
                <span class="font-medium">Running: </span>
                <span :class="props.driveMotor.isRunning ? 'text-green-500' : 'text-red-500'">
                  {{ props.driveMotor.isRunning ? 'Yes' : 'No' }}
                </span>
              </p>
              <p>
                <span class="font-medium">Enabled: </span>
                <span :class="props.driveMotor.enabled ? 'text-green-500' : 'text-red-500'">
                  {{ props.driveMotor.enabled ? 'Yes' : 'No' }}
                </span>
              </p>
            </div>
            <p><span class="font-medium">Last Updated: </span>{{ new Date(props.driveMotor.updatedAt).toISOString() }}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
