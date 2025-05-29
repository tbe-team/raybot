<script setup lang="ts">
import type { Info } from '@/types/info'
import type { RobotState } from '@/types/robot-state'
import BatteryCard from './BatteryCard.vue'
import ConnectionsCard from './ConnectionsCard.vue'
import InfoCard from './InfoCard.vue'
import LocationCard from './LocationCard.vue'

const props = defineProps<{
  robotState: RobotState
  systemInfo: Info
}>()
</script>

<template>
  <div class="grid grid-cols-1 gap-4 mb-6 md:grid-cols-2 lg:grid-cols-4">
    <BatteryCard
      :battery="props.robotState.battery"
    />
    <InfoCard
      :info="props.systemInfo"
    />

    <!-- Location Card -->
    <LocationCard
      :location="props.robotState.location"
    />

    <!-- Connections Card -->
    <ConnectionsCard
      :app-connection="props.robotState.appConnection"
    />
  </div>
</template>
