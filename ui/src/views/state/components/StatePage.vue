<script setup lang="ts">
import type { LocationState, RobotState } from '@/types/state'
import { CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import BatteryCard from '@/views/state/components/BatteryCard.vue'
import ChargeCard from '@/views/state/components/ChargeCard.vue'
import DischargeCard from '@/views/state/components/DischargeCard.vue'
import DistanceSensorCard from '@/views/state/components/DistanceSensorCard.vue'
import DriveMotorCard from '@/views/state/components/DriveMotorCard.vue'
import LiftMotorCard from '@/views/state/components/LiftMotorCard.vue'
import LocationCard from '@/views/state/components/LocationCard.vue'

interface Props {
  robotState: RobotState
  locationState: LocationState
}

const props = defineProps<Props>()
</script>

<template>
  <div>
    <CardHeader>
      <div class="space-y-2">
        <CardTitle>
          Robot state
        </CardTitle>
        <CardDescription>
          The current state of the robot is updated once per second.
        </CardDescription>
      </div>
    </CardHeader>

    <CardContent class="grid grid-cols-4 gap-4">
      <BatteryCard :battery="props.robotState.battery" class="col-span-2" />
      <ChargeCard :charge="props.robotState.charge" />
      <DischargeCard :discharge="props.robotState.discharge" />
      <LiftMotorCard :lift-motor="props.robotState.liftMotor" />
      <DriveMotorCard :drive-motor="props.robotState.driveMotor" />
      <DistanceSensorCard :distance-sensor="props.robotState.distanceSensor" />
      <LocationCard :location="props.locationState" />
    </CardContent>
  </div>
</template>
