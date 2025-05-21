<script setup lang="ts">
import type { CargoDoorMotorState } from '@/types/cargo'
import type { DriveMotorState, LiftMotorState } from '@/types/robot-state'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
import { Separator } from '@/components/ui/separator'
import { formatDate } from '@/lib/date'
import { ArrowLeft, ArrowRight, Lock, Unlock } from 'lucide-vue-next'

const props = defineProps<{
  liftMotor: LiftMotorState
  driveMotor: DriveMotorState
  cargoDoorMotor: CargoDoorMotorState
}>()

const liftMotorProgress = computed(() => {
  const percent = (props.liftMotor.currentPosition / props.liftMotor.targetPosition) * 100
  return percent > 100 ? 100 : percent
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Motor Status</CardTitle>
      <CardDescription>Lift and drive motor information</CardDescription>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
        <!-- Lift Motor -->
        <div class="space-y-4">
          <h3 class="font-medium">
            Lift Motor
          </h3>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium">Status</span>
              <Badge :variant="props.liftMotor.isRunning ? 'default' : 'outline'">
                {{ props.liftMotor.isRunning ? 'Running' : 'Idle' }}
              </Badge>
            </div>

            <div class="space-y-1">
              <div class="flex justify-between text-sm">
                <span>Current Position</span>
                <span>{{ props.liftMotor.currentPosition }}cm</span>
              </div>
            </div>

            <div class="space-y-1">
              <div class="flex justify-between text-sm">
                <span>Target Position</span>
                <span>{{ props.liftMotor.targetPosition }}cm</span>
              </div>
              <div class="relative h-2 rounded-full bg-muted">
                <div
                  class="absolute w-1 h-4 -translate-y-1/2 rounded-full bg-primary top-1/2"
                  :style="{ left: `${liftMotorProgress}%` }"
                />
              </div>
            </div>

            <div class="flex items-center justify-between">
              <span class="text-sm font-medium">Enabled</span>
              <Badge :variant="props.liftMotor.enabled ? 'default' : 'destructive'">
                {{ props.liftMotor.enabled ? 'Yes' : 'No' }}
              </Badge>
            </div>

            <div class="text-xs text-muted-foreground">
              Last updated: {{ formatDate(props.liftMotor.updatedAt) }}
            </div>
          </div>
        </div>

        <!-- Drive Motor -->
        <div class="space-y-4">
          <h3 class="font-medium">
            Drive Motor
          </h3>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium">Status</span>
              <Badge :variant="props.driveMotor.isRunning ? 'default' : 'outline'">
                {{ props.driveMotor.isRunning ? 'Running' : 'Idle' }}
              </Badge>
            </div>

            <div class="flex items-center justify-between">
              <span class="text-sm font-medium">Direction</span>
              <div class="flex items-center gap-2">
                <ArrowRight
                  v-if="props.driveMotor.direction === 'FORWARD'"
                  class="w-4 h-4 leading-5 text-success"
                />
                <ArrowLeft
                  v-else
                  class="w-4 h-4 leading-5 text-warning"
                />
                <span class="leading-5">{{ props.driveMotor.direction }}</span>
              </div>
            </div>

            <div class="space-y-1">
              <div class="flex justify-between text-sm">
                <span>Speed</span>
                <span>{{ props.driveMotor.speed }}%</span>
              </div>
              <Progress :model-value="props.driveMotor.speed" :max="100" class="h-2" />
            </div>

            <div class="flex items-center justify-between">
              <span class="text-sm font-medium">Enabled</span>
              <Badge :variant="props.driveMotor.enabled ? 'default' : 'destructive'">
                {{ props.driveMotor.enabled ? 'Yes' : 'No' }}
              </Badge>
            </div>

            <div class="text-xs text-muted-foreground">
              Last updated: {{ formatDate(props.driveMotor.updatedAt) }}
            </div>
          </div>
        </div>
      </div>

    </CardContent>
  </Card>
</template>
