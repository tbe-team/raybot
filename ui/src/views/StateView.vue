<script setup lang="ts">
import { AlertCircle, Loader } from 'lucide-vue-next'
import BatteryContent from '@/components/app/state/BatteryContent.vue'
import CargoContent from '@/components/app/state/CargoContent.vue'
import CargoDoorMotorContent from '@/components/app/state/CargoDoorMotorContent.vue'
import ConnectionsContent from '@/components/app/state/ConnectionsContent.vue'
import DistanceSensorContent from '@/components/app/state/DistanceSensorContent.vue'
import DriveMotorContent from '@/components/app/state/DriveMotorContent.vue'
import LedContent from '@/components/app/state/LedContent.vue'
import LiftMotorContent from '@/components/app/state/LiftMotorContent.vue'
import LimitSwitchContent from '@/components/app/state/LimitSwitchContent.vue'
import LocationContent from '@/components/app/state/LocationContent.vue'
import SystemInfoContent from '@/components/app/state/SystemInfoContent.vue'
import PageContainer from '@/components/shared/PageContainer.vue'
import { Card } from '@/components/ui/card'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useQueryRobotState } from '@/composables/use-robot-state'

const REFRESH_INTERVAL = 1000

const refetchInterval = ref(REFRESH_INTERVAL)

const { data: robotState, isPending, isError, error } = useQueryRobotState({
  axiosOpts: { doNotShowLoading: true },
  refetchInterval,
})
</script>

<template>
  <PageContainer>
    <div v-if="isPending" class="flex flex-col gap-4 justify-center items-center pt-20">
      <div class="flex gap-4 items-center">
        <Loader class="w-8 h-8 animate-spin text-muted-foreground" />
      </div>
      <p class="text-lg text-muted-foreground">
        Loading state...
      </p>
    </div>
    <div v-else-if="isError" class="flex flex-col gap-4 justify-center items-center pt-20">
      <Card class="flex flex-col gap-4 items-center p-6 text-red-500">
        <AlertCircle class="w-8 h-8" />
        <div class="space-y-2 text-center">
          <h2 class="text-lg font-semibold">
            Failed to load state
          </h2>
          <p class="text-sm text-muted-foreground">
            {{ error?.message || 'An unexpected error occurred' }}
          </p>
        </div>
      </Card>
    </div>
    <div v-else-if="!robotState" class="flex flex-col gap-4 justify-center items-center pt-20">
      <Card class="flex flex-col gap-4 items-center p-6">
        <AlertCircle class="w-8 h-8 text-muted-foreground" />
        <div class="space-y-2 text-center">
          <h2 class="text-lg font-semibold">
            Robot state not found
          </h2>
          <p class="text-sm text-muted-foreground">
            The robot state appears to be empty
          </p>
        </div>
      </Card>
    </div>
    <div v-else class="flex flex-col w-full">
      <div class="flex justify-between items-center mb-6">
        <div>
          <h1 class="text-xl font-semibold tracking-tight">
            State Dashboard
          </h1>
          <p class="text-sm text-muted-foreground">
            The current state of the robot is continuously updated.
          </p>
        </div>
        <div class="flex gap-2 items-center">
          <span class="whitespace-nowrap">Refresh rate: </span>
          <Select v-model="refetchInterval">
            <SelectTrigger>
              <SelectValue class="w-5" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem v-for="interval in [1000, 3000, 5000, 10000]" :key="interval" :value="interval">
                  <SelectValue>{{ interval / 1000 }}</SelectValue>
                </SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
          <span>seconds</span>
        </div>
      </div>
      <div class="flex flex-col gap-4">
        <BatteryContent class="break-inside-avoid" :battery="robotState.battery" :charge="robotState.charge" :discharge="robotState.discharge" />
        <div class="space-y-4 columns-1 md:columns-3">
          <DriveMotorContent class="break-inside-avoid" :drive-motor="robotState.driveMotor" />
          <LiftMotorContent class="break-inside-avoid" :lift-motor="robotState.liftMotor" />
          <DistanceSensorContent class="break-inside-avoid" :distance-sensor="robotState.distanceSensor" />
        </div>

        <div class="space-y-4 columns-1 md:columns-2">
          <LimitSwitchContent v-model:refetch-interval="refetchInterval" class="break-inside-avoid" />
          <CargoContent class="break-inside-avoid" :cargo="robotState.cargo" />
          <LocationContent class="break-inside-avoid" :location="robotState.location" />
          <CargoDoorMotorContent class="break-inside-avoid" :cargo-door-motor="robotState.cargoDoorMotor" />
        </div>

        <div class="space-y-4 columns-1 md:columns-2">
          <SystemInfoContent class="break-inside-avoid" />
          <LedContent class="break-inside-avoid" :led="robotState.leds" />
          <ConnectionsContent class="break-inside-avoid" :app-connection="robotState.appConnection" />
        </div>
      </div>
    </div>
  </PageContainer>
</template>
