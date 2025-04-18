<script setup lang="ts">
import type { CargoCheckQRInputs, MoveToInputs } from '@/types/command'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useCancelProcessingCommandMutation, useCurrentProcessingCommandQuery } from '@/composables/use-command'
import { formatDate } from '@/lib/date'
import { useConfirmationStore } from '@/stores/confirmation-store'
import { RaybotError } from '@/types/error'
import { Clock, Loader, MoreHorizontal } from 'lucide-vue-next'
import SourceBadge from './SourceBadge.vue'
import StatusBadge from './StatusBadge.vue'
import { getCommandIcon, getCommandName } from './utils'

const emit = defineEmits<{
  (e: 'viewDetails', commandId: number): void
}>()

const { openConfirmation } = useConfirmationStore()

const { data: command, refetch, isError } = useCurrentProcessingCommandQuery({ axiosOpts: { doNotShowLoading: true } })
const { mutate: cancelProcessingCommand, isPending: isCancellingCommand } = useCancelProcessingCommandMutation()
const REFRESH_INTERVAL = 1000
const interval = setInterval(refetch, REFRESH_INTERVAL)

onUnmounted(() => {
  clearInterval(interval)
})

function handleCancelCommand() {
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
  <Card v-if="command && !isError">
    <CardHeader class="pb-3">
      <CardTitle class="flex items-center justify-between">
        <span>Current Command</span>
        <StatusBadge :status="command.status" />
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div v-if="isCancellingCommand">
        <div class="flex items-center justify-center gap-2 p-9">
          <Loader class="w-4 h-4 animate-spin" />
          <span>Cancelling command...</span>
        </div>
      </div>
      <div
        v-else
        class="p-4 space-y-3 border rounded-lg cursor-pointer bg-muted/50"
        @click="emit('viewDetails', command.id)"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 font-medium">
            <component :is="getCommandIcon(command.type)" class="w-5 h-5" />
            <span>{{ getCommandName(command.type) }}</span>
          </div>
          <div class="flex items-center gap-2" @click.stop>
            <SourceBadge :source="command.source" />
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <Button variant="ghost" size="icon" class="w-8 h-8">
                  <MoreHorizontal class="w-4 h-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem class="cursor-pointer" @click.stop="emit('viewDetails', command.id)">
                  View details
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem class="text-red-500 cursor-pointer" @click="handleCancelCommand">
                  Cancel command
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
        <div class="flex items-center gap-2 text-sm text-muted-foreground">
          <Clock class="w-4 h-4" />
          <span>Started at: {{ command.startedAt ? formatDate(command.startedAt) : 'N/A' }}</span>
        </div>

        <template v-if="command.type === 'MOVE_TO'">
          <div class="text-sm">
            <span class="font-medium">Location: </span> {{ (command.inputs as MoveToInputs).location }}
          </div>
        </template>
        <template v-else-if="command.type === 'CARGO_CHECK_QR'">
          <div class="text-sm">
            <span class="font-medium">QR Code: </span> {{ (command.inputs as CargoCheckQRInputs).qrCode }}
          </div>
        </template>
      </div>
    </CardContent>
  </Card>
</template>
