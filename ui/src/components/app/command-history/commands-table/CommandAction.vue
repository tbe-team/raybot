<script setup lang="ts">
import type { Command } from '@/types/command'
import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { useDeleteCommandMutation } from '@/composables/use-command'
import { RaybotError } from '@/types/error'
import { MoreHorizontal } from 'lucide-vue-next'

const props = defineProps<{
  command: Command
}>()

const { mutate: deleteCommand } = useDeleteCommandMutation()

function handleRemoveFromHistory() {
  deleteCommand(props.command.id, {
    onSuccess: () => {
      notification.success('Command removed from history')
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
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" size="icon" class="w-8 h-8" @click.stop>
        <MoreHorizontal class-name="h-4 w-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end" @click.stop>
      <DropdownMenuItem class="text-red-500" :disabled="command.status === 'PROCESSING'" @click="handleRemoveFromHistory">
        Remove from History
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
