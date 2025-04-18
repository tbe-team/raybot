<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import { useCancelProcessingCommandMutation, useDeleteCommandMutation, useGetCommandQuery } from '@/composables/use-command'
import { formatDate, formatUptime } from '@/lib/date'
import { useConfirmationStore } from '@/stores/confirmation-store'
import { RaybotError } from '@/types/error'
import { useIntervalFn } from '@vueuse/core'
import { Loader2 } from 'lucide-vue-next'
import SourceBadge from './SourceBadge.vue'
import StatusBadge from './StatusBadge.vue'
import { getCommandIcon, getCommandName } from './utils'

const props = defineProps<{
  commandId: number
}>()

const isOpen = defineModel<boolean>('isOpen', { required: true })

const commandId = toRef(props, 'commandId')
const { data: command, isPending, isError, refetch } = useGetCommandQuery(commandId, {
  axiosOpts: {
    doNotShowLoading: true,
  },
})
const { mutate: deleteCommand } = useDeleteCommandMutation()
const { mutate: cancelProcessingCommand } = useCancelProcessingCommandMutation()
const { openConfirmation } = useConfirmationStore()

const REFRESH_INTERVAL = 1000
const { pause, resume } = useIntervalFn(() => {
  refetch()
}, REFRESH_INTERVAL, { immediate: false })

watch(command, (cmd) => {
  if (cmd && ['SUCCEEDED', 'FAILED', 'CANCELED'].includes(cmd.status)) {
    pause()
  }
  else {
    resume()
  }
}, { immediate: true })

watch(isOpen, (open) => {
  if (!open) {
    pause()
  }
  else {
    resume()
  }
}, { immediate: true })

function handleRemoveFromQueue() {
  openConfirmation({
    title: 'Remove command',
    description: 'Are you sure you want to remove this command from queue?',
    actionLabel: 'Confirm',
    cancelLabel: 'Cancel',
    onAction: () => {
      deleteCommand(props.commandId, {
        onSuccess: () => {
          notification.success('Command removed from queue')
        },
        onError: (error) => {
          if (error instanceof RaybotError) {
            if (error.errorCode === 'command.inProcessingCanNotBeDeleted') {
              notification.error('Command is being processed and cannot be deleted')
            }
            else {
              notification.error(error.message)
            }
          }
          else {
            notification.error(error.message)
          }
        },
      })
    },
    onCancel: () => {
    },
  })
}

function handleStopExecution() {
  openConfirmation({
    title: 'Cancel command',
    description: 'Are you sure you want to cancel this command?',
    actionLabel: 'Confirm',
    cancelLabel: 'Cancel',
    onAction: () => {
      cancelProcessingCommand(undefined, {
        onSuccess: () => {
          notification.success('Command cancelled successfully')
        },
        onError: (error) => {
          if (error instanceof RaybotError) {
            if (error.errorCode === 'command.noCommandBeingProcessed') {
              notification.error('No command is being processed')
            }
            else {
              notification.error(error.message)
            }
          }
          else {
            notification.error(error.message)
          }
        },
      })
    },
    onCancel: () => {
    },
  })
}
</script>

<template>
  <Sheet v-model:open="isOpen">
    <SheetContent class="max-h-screen overflow-y-auto sm:max-w-md">
      <SheetHeader>
        <SheetTitle>
          Command detail
        </SheetTitle>
        <SheetDescription>
          Real-time information about the selected command
        </SheetDescription>
      </SheetHeader>
      <div v-if="isPending">
        <div class="flex items-center justify-center py-8">
          <Loader2 class="w-8 h-8 animate-spin" />
        </div>
      </div>
      <div v-else-if="isError">
        <div class="flex items-center justify-center py-8">
          <p class="text-sm text-muted-foreground">
            Failed to load command details
          </p>
        </div>
      </div>
      <div v-else-if="command" class="mt-6 space-y-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 ">
            <component :is="getCommandIcon(command.type)" />
            <span class="text-sm font-medium">
              {{ getCommandName(command.type) }}
            </span>
          </div>
          <StatusBadge :status="command.status" />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <p class="text-sm text-muted-foreground">
              ID
            </p>
            <p class="font-medium">
              {{ command.id }}
            </p>
          </div>
          <div class="flex flex-col items-start gap-2">
            <p class="text-sm text-muted-foreground">
              Source
            </p>
            <SourceBadge :source="command.source" />
          </div>
          <div v-if="command.completedAt && command.startedAt">
            <p class="text-sm text-muted-foreground">
              Duration
            </p>
            <p class="font-medium">
              {{ formatUptime((new Date(command.completedAt).getTime() - new Date(command.startedAt).getTime()) / 1000) }}
            </p>
          </div>
        </div>

        <div class="space-y-2">
          <p class="text-sm font-medium">
            Command Inputs
          </p>
          <div class="p-4 bg-gray-100 rounded-xl dark:bg-gray-800">
            <span class="font-mono text-sm text-gray-800 break-words whitespace-pre-wrap dark:text-gray-200">
              {{ JSON.stringify(command.inputs, null, 4) }}
            </span>
          </div>
        </div>

        <div v-if="command.status !== 'PROCESSING' && command.status !== 'QUEUED'" class="space-y-2">
          <p class="text-sm font-medium">
            Command Outputs
          </p>
          <div class="p-4 bg-gray-100 rounded-xl dark:bg-gray-800">
            <span class="font-mono text-sm text-gray-800 break-words whitespace-pre-wrap dark:text-gray-200">
              {{ JSON.stringify(command.outputs, null, 4) }}
            </span>
          </div>
        </div>

        <div v-if="command.error" class="space-y-2">
          <p class="text-sm font-medium text-red-500">
            Error
          </p>
          <div class="p-3 text-red-500 rounded-md bg-red-500/10">
            {{ command.error }}
          </div>
        </div>

        <div class="space-y-2">
          <p class="text-sm font-medium">
            Timeline
          </p>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-muted-foreground">Created</span>
              <span>{{ formatDate(command.createdAt) }}</span>
            </div>
            <div v-if="command.startedAt" class="flex justify-between text-sm">
              <span class="text-muted-foreground">Started</span>
              <span>{{ formatDate(command.startedAt) }}</span>
            </div>
            <div v-if="command.completedAt" class="flex justify-between text-sm">
              <span class="text-muted-foreground">Completed</span>
              <span>{{ formatDate(command.completedAt) }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-muted-foreground">Last Updated</span>
              <span>{{ formatDate(command.updatedAt) }}</span>
            </div>
          </div>
        </div>

        <div class="flex gap-2 pt-4">
          <template v-if="command.status === 'QUEUED'">
            <Button variant="destructive" class="flex-1" @click="handleRemoveFromQueue">
              Remove from queue
            </Button>
          </template>
          <template v-else-if="command.status === 'PROCESSING'">
            <Button variant="destructive" class="flex-1" @click="handleStopExecution">
              Stop execution
            </Button>
          </template>
          <template v-else>
            <Button variant="outline" class="flex-1" disabled>
              Re-run
            </Button>
          </template>
        </div>
      </div>
    </SheetContent>
  </Sheet>
</template>
