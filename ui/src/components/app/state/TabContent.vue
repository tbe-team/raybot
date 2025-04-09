<script setup lang="ts">
import type { RobotState } from '@/types/robot-state'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import BatteryTabContent from './BatteryTabContent.vue'
import CargoTabContent from './CargoTabContent.vue'
import ConnectionsTabContent from './ConnectionsTabContent.vue'
import DistanceSensorsTabContent from './DistanceSensorsTabContent.vue'
import MotorsTabContent from './MotorsTabContent.vue'

const props = defineProps<{
  robotState: RobotState
}>()

const route = useRoute()
const router = useRouter()
const tab = route.query.tab as string | undefined ?? 'battery'

function handleTabChange(value: string | number) {
  router.replace({ query: { tab: value } })
}
</script>

<template>
  <div class="w-full">
    <Tabs :default-value="tab" @update:model-value="handleTabChange">
      <TabsList class="mb-4">
        <TabsTrigger value="battery">
          Battery
        </TabsTrigger>
        <TabsTrigger value="motors">
          Motors
        </TabsTrigger>
        <TabsTrigger value="sensors">
          Sensors
        </TabsTrigger>
        <TabsTrigger value="cargo">
          Cargo
        </TabsTrigger>
        <TabsTrigger value="connections">
          Connections
        </TabsTrigger>
      </TabsList>

      <TabsContent value="battery">
        <BatteryTabContent
          :battery="props.robotState.battery"
          :charge="props.robotState.charge"
          :discharge="props.robotState.discharge"
        />
      </TabsContent>

      <TabsContent value="motors">
        <MotorsTabContent
          :lift-motor="props.robotState.liftMotor"
          :drive-motor="props.robotState.driveMotor"
          :cargo-door-motor="props.robotState.cargoDoorMotor"
        />
      </TabsContent>

      <TabsContent value="sensors">
        <DistanceSensorsTabContent
          :distance-sensor="props.robotState.distanceSensor"
        />
      </TabsContent>

      <TabsContent value="cargo">
        <CargoTabContent
          :cargo="props.robotState.cargo"
          :cargo-door-motor="props.robotState.cargoDoorMotor"
        />
      </TabsContent>

      <TabsContent value="connections">
        <ConnectionsTabContent
          :app-connection="props.robotState.appConnection"
        />
      </TabsContent>
    </Tabs>
  </div>
</template>
