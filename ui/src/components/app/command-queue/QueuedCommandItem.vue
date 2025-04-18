<script setup lang="ts">
import type { CargoCheckQRInputs, Command, MoveToInputs } from '@/types/command'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useDeleteCommandMutation } from '@/composables/use-command'
import { useConfirmationStore } from '@/stores/confirmation-store'
import { RaybotError } from '@/types/error'
import { Clock, MoreHorizontal } from 'lucide-vue-next'
import SourceBadge from './SourceBadge.vue'
import StatusBadge from './StatusBadge.vue'
import { getCommandIcon, getCommandName, timeSince } from './utils'

const props = defineProps<{
  command: Command
}>()
const emit = defineEmits<{
  (e: 'viewDetails', commandId: number): void
}>()

const { openConfirmation } = useConfirmationStore()

const { mutate: deleteCommand } = useDeleteCommandMutation()

function handleRemoveFromQueue() {
  openConfirmation({
    title: 'Remove command',
    description: 'Are you sure you want to remove this command from queue?',
    actionLabel: 'Confirm',
    cancelLabel: 'Cancel',
    onAction: () => {
      deleteCommand(props.command.id, {
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
</script>

<template>
  <div
    class="p-4 space-y-3 transition-colors border rounded-lg cursor-pointer bg-muted/30 border-border/50 hover:border-border"
    @click="emit('viewDetails', props.command.id)"
  >
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2 font-medium">
        <component :is="getCommandIcon(props.command.type)" class="w-5 h-5" />
        <span>{{ getCommandName(props.command.type) }}</span>
      </div>
      <div class="flex items-center gap-2" @click.stop>
        <StatusBadge :status="props.command.status" />
        <SourceBadge :source="props.command.source" />
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button variant="ghost" size="icon" class="w-8 h-8">
              <MoreHorizontal class-name="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" @click.stop>
            <DropdownMenuItem @click="emit('viewDetails', props.command.id)">
              View details
            </DropdownMenuItem>
            <DropdownMenuItem>Edit command</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem class="text-red-500" @click="handleRemoveFromQueue">
              Remove from queue
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
    <div class="flex items-center gap-2 text-sm text-muted-foreground">
      <Clock class="w-4 h-4" />
      <span>Created {{ timeSince(props.command.createdAt) }}</span>
    </div>

    <template v-if="props.command.type === 'MOVE_TO'">
      <div class="text-sm">
        <span class="font-medium">Location: </span> {{ (command.inputs as MoveToInputs).location }}
      </div>
    </template>
    <template v-else-if="props.command.type === 'CARGO_CHECK_QR'">
      <div class="text-sm">
        <span class="font-medium">QR Code: </span> {{ (command.inputs as CargoCheckQRInputs).qrCode }}
      </div>
    </template>
  </div>
</template>
