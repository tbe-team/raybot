<script setup lang="ts">
import type { RobotState } from '@/types/robot-state'
import BatteryCard from './BatteryCard.vue'
import ConnectionsCard from './ConnectionsCard.vue'
import LocationCard from './LocationCard.vue'
import SystemInfoCard from './SystemInfoCard.vue'

const props = defineProps<{
  robotState: RobotState
}>()
</script>

<template>
  <div class="grid grid-cols-1 gap-4 mb-6 md:grid-cols-2 lg:grid-cols-4">
    <!-- Battery Card -->
    <BatteryCard
      :battery="props.robotState.battery"
    />

    <!-- Info Card -->
    <SystemInfoCard />

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
