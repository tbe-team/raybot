<script setup lang="ts">
import type { Alarm } from '@/types/alarm'
import { MoreHorizontal } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { useDeactivateAlarmMutation } from '@/composables/use-alarm'
import { useConfirmationStore } from '@/stores/confirmation-store'

const props = defineProps<{
  alarm: Alarm
}>()

const { openConfirmation } = useConfirmationStore()

const { mutate: deactivateAlarm } = useDeactivateAlarmMutation()

function handleDeactivateAlarm() {
  openConfirmation({
    title: 'Deactivate alarm',
    description: 'Are you sure you want to deactivate this alarm?',
    actionLabel: 'Confirm',
    cancelLabel: 'Cancel',
    onAction: () => {
      deactivateAlarm(props.alarm.id, {
        onSuccess: () => {
          notification.success('Alarm deactivated')
        },
        onError: (error) => {
          notification.error(error.message)
        },
      })
    },
    onCancel: () => {
    },
  })
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" size="icon" class="w-8 h-8" @click.stop>
        <MoreHorizontal class="w-4 h-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end">
      <DropdownMenuItem class="text-red-500" @click="handleDeactivateAlarm">
        Deactivate alarm
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
