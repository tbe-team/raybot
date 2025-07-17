<script setup lang="ts">
import { useColorMode } from '@vueuse/core'
import { Moon, Pause, Sun } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { useSystemStopEmergencyMutation } from '@/composables/use-system'
import { useConfirmationStore } from '@/stores/confirmation-store'
import AlarmButton from './AlarmButton.vue'

const { store } = useColorMode()

const { mutate: stopEmergency } = useSystemStopEmergencyMutation()
const { openConfirmation } = useConfirmationStore()

function handleEmergencyStop() {
  openConfirmation({
    title: 'Stop Emergency',
    description:
      'Stopping emergency will stop all motors, canceling all commands. Are you sure you want to continue?',
    actionLabel: 'Confirm',
    cancelLabel: 'Cancel',
    onAction: () => {
      stopEmergency(undefined, {
        onSuccess: () => {
          notification.success('Stop emergency successfully')
        },
        onError: () => {
          notification.error('Failed to stop emergency')
        },
      })
    },
    onCancel: () => {},
  })
}
</script>

<template>
  <div class="flex gap-2 items-center">
    <Button class="text-destructive" variant="ghost" @click="handleEmergencyStop">
      <Pause class="w-4 h-4" />
      STOP EMERGENCY
    </Button>

    <AlarmButton />

    <Button
      variant="ghost"
      size="icon"
      class="rounded-lg bg-muted hover:bg-muted-hover"
      @click="store = store === 'light' ? 'dark' : 'light'"
    >
      <Sun v-if="store === 'light'" class="w-5 h-5" />
      <Moon v-else class="w-5 h-5" />
      <span class="sr-only">Toggle theme</span>
    </Button>
  </div>
</template>
