<script setup lang="ts">
import type { CargoDoorMotorState } from '@/types/cargo'
import type { DriveMotorState, LiftMotorState } from '@/types/robot-state'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Gauge } from 'lucide-vue-next'

const props = defineProps<{
  liftMotor: LiftMotorState
  driveMotor: DriveMotorState
  cargoDoorMotor: CargoDoorMotorState
}>()

function getMotorStatus(isRunning: boolean): { variant: 'default' | 'secondary', label: string } {
  return {
    variant: isRunning ? 'default' : 'secondary',
    label: isRunning ? 'Active' : 'Idle',
  }
}
</script>

<template>
  <Card>
    <CardHeader class="pb-2">
      <CardTitle class="flex items-center gap-2 text-sm font-medium">
        <Gauge class="w-6 h-6" />
        Motors
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="space-y-2">
        <div class="flex items-center justify-between">
          <span class="text-sm">Lift motor</span>
          <Badge :variant="getMotorStatus(props.liftMotor.isRunning).variant">
            {{ getMotorStatus(props.liftMotor.isRunning).label }}
          </Badge>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm">Drive motor</span>
          <Badge :variant="getMotorStatus(props.driveMotor.isRunning).variant">
            {{ getMotorStatus(props.driveMotor.isRunning).label }}
          </Badge>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm">Cargo door motor</span>
          <Badge :variant="getMotorStatus(props.cargoDoorMotor.isRunning).variant">
            {{ getMotorStatus(props.cargoDoorMotor.isRunning).label }}
          </Badge>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
