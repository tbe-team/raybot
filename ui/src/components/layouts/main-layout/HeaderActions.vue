<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { useEmergencyResumeMutation, useEmergencyStateQuery, useEmergencyStopMutation } from '@/composables/use-emergency'
import { useConfirmationStore } from '@/stores/confirmation-store'
import { useColorMode } from '@vueuse/core'
import { Moon, Pause, Play, Sun } from 'lucide-vue-next'

const { store } = useColorMode()

const { data: emergencyState, refetch: refetchEmergencyState } = useEmergencyStateQuery()

const { mutate: stopEmergency } = useEmergencyStopMutation()
const { mutate: resumeEmergency } = useEmergencyResumeMutation()

const { openConfirmation } = useConfirmationStore()

function handleEmergencyStop() {
  if (emergencyState.value?.locked) {
    resumeEmergency(undefined, {
      onSuccess: () => {
        refetchEmergencyState()
      },
      onError: () => {
        notification.error('Failed to resume emergency')
      },
    })
  }
  else {
    openConfirmation({
      title: 'Stop emergency',
      description: 'Are you sure you want to stop the emergency?',
      actionLabel: 'Confirm',
      cancelLabel: 'Cancel',
      onAction: () => {
        stopEmergency(undefined, {
          onSuccess: () => {
            refetchEmergencyState()
          },
          onError: () => {
            notification.error('Failed to stop emergency')
          },
        })
      },
      onCancel: () => {
      },
    })
  }
}
</script>

<template>
  <div class="flex items-center gap-2">
    <Button variant="ghost" @click="handleEmergencyStop">
      <div class="flex items-center gap-2">
        <span v-if="emergencyState?.locked" class="flex items-center gap-2 font-bold text-destructive"><Play class="w-4 h-4" /> RESUME EMERGENCY</span>
        <span v-else class="flex items-center gap-2 font-bold text-destructive"><Pause class="w-4 h-4" /> STOP EMERGENCY</span>
      </div>
    </Button>
    <Button
      class="rounded-lg bg-muted hover:bg-muted-hover"
      variant="ghost"
      size="icon"
      as-child
    >
      <a
        href="https://github.com/tbe-team/raybot"
        target="_blank"
        rel="noopener noreferrer"
        class="text-foreground"
      >
        <svg role="img" viewBox="0 0 24 24">
          <path fill="currentColor" d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12" />
        </svg>
        <span class="sr-only">GitHub</span>
      </a>
    </Button>

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
