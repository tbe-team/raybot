<script setup lang="ts">
import type { LiftMotorState } from '@/types/state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface Props {
  liftMotor: LiftMotorState
}

const props = defineProps<Props>()
</script>

<script lang="ts">
</script>

<template>
  <div>
    <Card class="h-full rounded-sm shadow-lg">
      <CardHeader>
        <CardTitle>Lift Motor State</CardTitle>
      </CardHeader>

      <CardContent>
        <div class="mx-auto space-x-2">
          <div class="grid grid-rows-1 gap-2 text-sm">
            <div class="grid grid-cols-2 gap-2">
              <p><span class="font-medium">Current Position: </span>{{ props.liftMotor.currentPosition }}%</p>
              <p><span class="font-medium">Target Position: </span>{{ props.liftMotor.targetPosition }}%</p>
              <p>
                <span class="font-medium">Running: </span>
                <span :class="props.liftMotor.isRunning ? 'text-green-500' : 'text-red-500'">{{ props.liftMotor.isRunning ? 'Yes' : 'No' }}</span>
              </p>
              <p>
                <span class="font-medium">Enabled: </span>
                <span :class="props.liftMotor.enabled ? 'text-green-500' : 'text-red-500'">{{ props.liftMotor.enabled ? 'Yes' : 'No' }}</span>
              </p>
            </div>
            <p><span class="font-medium">Last Updated:</span> {{ (props.liftMotor.updatedAt).toISOString() }}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
