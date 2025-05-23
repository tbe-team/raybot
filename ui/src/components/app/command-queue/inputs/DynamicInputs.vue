<script setup lang="ts">
import type { CommandType } from '@/types/command'
import CargoCheckQRInputs from './CargoCheckQRInputs.vue'
import CargoCloseInputs from './CargoCloseInputs.vue'
import CargoLiftInputs from './CargoLiftInputs.vue'
import CargoLowerInputs from './CargoLowerInputs.vue'
import CargoOpenInputs from './CargoOpenInputs.vue'
import MoveBackwardInputs from './MoveBackwardInputs.vue'
import MoveForwardInputs from './MoveForwardInputs.vue'
import MoveToInputs from './MoveToInputs.vue'
import WaitInputs from './WaitInputs.vue'

const props = defineProps<{
  commandType: CommandType
}>()

const componentMap: Record<CommandType, Component | null> = {
  MOVE_TO: MoveToInputs,
  MOVE_FORWARD: MoveForwardInputs,
  MOVE_BACKWARD: MoveBackwardInputs,
  CARGO_OPEN: CargoOpenInputs,
  CARGO_CLOSE: CargoCloseInputs,
  CARGO_LIFT: CargoLiftInputs,
  CARGO_LOWER: CargoLowerInputs,
  CARGO_CHECK_QR: CargoCheckQRInputs,
  WAIT: WaitInputs,
  STOP_MOVEMENT: null,
  SCAN_LOCATION: null,
}

const inputComponent = computed(() => componentMap[props.commandType])
</script>

<template>
  <component :is="inputComponent" v-if="inputComponent" :key="commandType" />
</template>
