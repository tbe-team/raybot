<script setup lang="ts">
import BatteryTabContent from '@/components/app/state/BatteryTabContent.vue'
import CargoTabContent from '@/components/app/state/CargoTabContent.vue'
import ConnectionsTabContent from '@/components/app/state/ConnectionsTabContent.vue'
import ControlCard from '@/components/app/state/ControlCard.vue'
import DistanceSensorsTabContent from '@/components/app/state/DistanceSensorsTabContent.vue'
import MotorsTabContent from '@/components/app/state/MotorsTabContent.vue'
import StateOverview from '@/components/app/state/StateOverview.vue'
import PageContainer from '@/components/shared/PageContainer.vue'
import { Card } from '@/components/ui/card'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useQueryRobotState } from '@/composables/use-robot-state'
import { AlertCircle, Loader } from 'lucide-vue-next'

const REFRESH_INTERVAL = 1000
const refetchInterval = ref(REFRESH_INTERVAL)

const { data: robotState, isPending, isError, error } = useQueryRobotState({
  axiosOpts: { doNotShowLoading: true },
  refetchInterval,
})
</script>

<template>
  <PageContainer>
    <div v-if="isPending" class="flex flex-col items-center justify-center gap-4 pt-20">
      <div class="flex items-center gap-4">
        <Loader class="w-8 h-8 animate-spin text-muted-foreground" />
      </div>
      <p class="text-lg text-muted-foreground">
        Loading state...
      </p>
    </div>
    <div v-else-if="isError" class="flex flex-col items-center justify-center gap-4 pt-20">
      <Card class="flex flex-col items-center gap-4 p-6 text-red-500">
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
    <div v-else-if="!robotState" class="flex flex-col items-center justify-center gap-4 pt-20">
      <Card class="flex flex-col items-center gap-4 p-6">
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
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-xl font-semibold tracking-tight">
            Robot Status Dashboard
          </h1>
          <p class="text-sm text-muted-foreground">
            The current state of the robot is continuously updated.
          </p>
        </div>
        <div class="flex items-center gap-2">
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

      <!-- Overview Cards -->
      <StateOverview v-if="robotState" :robot-state="robotState" />

      <!-- Content Cards -->
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <MotorsTabContent
          v-if="robotState"
          :lift-motor="robotState.liftMotor"
          :drive-motor="robotState.driveMotor"
          :cargo-door-motor="robotState.cargoDoorMotor"
        />
        <ControlCard />

        <DistanceSensorsTabContent
          v-if="robotState"
          :distance-sensor="robotState.distanceSensor"
        />

        <CargoTabContent
          v-if="robotState"
          :cargo="robotState.cargo"
          :cargo-door-motor="robotState.cargoDoorMotor"
        />
      </div>
    </div>
  </PageContainer>
</template>
