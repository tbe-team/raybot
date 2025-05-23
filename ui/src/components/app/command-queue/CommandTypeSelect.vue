<script setup lang="ts">
import type { SelectRootEmits, SelectRootProps } from 'reka-ui'
import type { CommandType } from '@/types/command'
import { ArrowDown, ArrowUp, Clock, MapPin, Package, QrCode, Scan, StopCircle } from 'lucide-vue-next'
import { useForwardPropsEmits } from 'reka-ui'
import { FormControl } from '@/components/ui/form'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

const props = defineProps<SelectRootProps>()
const emits = defineEmits<SelectRootEmits>()

const forwarded = useForwardPropsEmits(props, emits)

const selectItemMap: Record<CommandType, { icon: Component, label: string }> = {
  STOP_MOVEMENT: { icon: StopCircle, label: 'Stop Movement' },
  MOVE_FORWARD: { icon: ArrowUp, label: 'Move Forward' },
  MOVE_BACKWARD: { icon: ArrowDown, label: 'Move Backward' },
  MOVE_TO: { icon: MapPin, label: 'Move To' },
  CARGO_OPEN: { icon: Package, label: 'Cargo Open' },
  CARGO_CLOSE: { icon: Package, label: 'Cargo Close' },
  CARGO_LIFT: { icon: Package, label: 'Cargo Lift' },
  CARGO_LOWER: { icon: Package, label: 'Cargo Lower' },
  CARGO_CHECK_QR: { icon: QrCode, label: 'Cargo Check QR' },
  SCAN_LOCATION: { icon: Scan, label: 'Scan Location' },
  WAIT: { icon: Clock, label: 'Wait' },
}
</script>

<template>
  <Select v-bind="forwarded">
    <FormControl>
      <SelectTrigger>
        <SelectValue placeholder="Select command type" />
      </SelectTrigger>
    </FormControl>
    <SelectContent>
      <SelectItem v-for="(val, key) in selectItemMap" :key="key" :value="key">
        <div class="flex items-center gap-2">
          <component :is="val.icon" class="w-4 h-4" />
          <span>{{ val.label }}</span>
        </div>
      </SelectItem>
    </SelectContent>
  </Select>
</template>
