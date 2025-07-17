<script setup lang="ts">
import type { Alarm } from '@/types/alarm'
import { AlertCircle, CircleMinus, Loader, RefreshCw } from 'lucide-vue-next'
import { columns } from '@/components/app/alarm/alarm-table/alarm-history/columns'
import DataTable from '@/components/shared/DataTable.vue'
import { Button } from '@/components/ui/button'
import { useDeleteDeactiveAlarmsMutation } from '@/composables/use-alarm'
import { useConfirmationStore } from '@/stores/confirmation-store'
import { RaybotError } from '@/types/error'

interface Props {
  page: number
  pageSize: number
  data: Alarm[]
  totalItems: number
  isFetching: boolean
  isPending: boolean
  isError: boolean
  error: Error | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:page', value: number): void
  (e: 'update:pageSize', value: number): void
  (e: 'update:data'): void
}>()

const { openConfirmation } = useConfirmationStore()
const { mutate: deleteAlarms, isPending: isDeletingAlarms } = useDeleteDeactiveAlarmsMutation()

function handleClearAlarmHistory() {
  openConfirmation({
    title: 'Clear alarm history',
    description: 'Are you sure you want to clear all alarm history?',
    actionLabel: 'Confirm',
    cancelLabel: 'Cancel',
    onAction: () => {
      deleteAlarms(undefined, {
        onSuccess: () => {
          emit('update:data')
        },
        onError: (error) => {
          if (error instanceof RaybotError) {
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
  <div v-if="isPending" class="flex flex-col gap-4 justify-center items-center pt-20">
    <div class="flex gap-4 items-center">
      <Loader class="animate-spin size-8 text-muted-foreground" />
    </div>
    <p class="text-lg text-muted-foreground">
      Loading alarms...
    </p>
  </div>

  <div v-else-if="props.isError" class="flex flex-col gap-4 justify-center items-center pt-20">
    <div class="flex flex-col gap-4 items-center p-6 text-red-500">
      <AlertCircle class="size-8" />
      <div class="space-y-2 text-center">
        <h2 class="text-lg font-semibold">
          Failed to load alarms
        </h2>
        <p class="text-sm text-muted-foreground">
          {{ props.error?.message || 'An unexpected error occurred' }}
        </p>
      </div>
    </div>
  </div>

  <div v-else-if="!props.data || props.data.length === 0" class="flex flex-col gap-4 justify-center items-center pt-20">
    <div class="flex flex-col gap-4 items-center p-6">
      <AlertCircle class="size-8 text-muted-foreground" />
      <div class="space-y-2 text-center">
        <h2 class="text-lg font-semibold">
          No alarms active found
        </h2>
        <p class="text-sm text-muted-foreground">
          There are no alarms to display
        </p>
      </div>
    </div>
  </div>

  <div v-else class="flex flex-col w-full">
    <div class="flex justify-end items-center mb-6">
      <div class="flex gap-2 items-center">
        <Button
          variant="outline"
          class="!text-destructive border-destructive"
          :disabled="isDeletingAlarms"
          @click="handleClearAlarmHistory"
        >
          <CircleMinus
            class="mr-2 size-4"
          />
          Clear History Alarms
        </Button>
        <Button
          variant="outline"
          :disabled="props.isFetching"
          @click="emit('update:data')"
        >
          <RefreshCw
            class="mr-2 size-4"
            :class="{ 'animate-spin': props.isFetching }"
          />
          Refresh
        </Button>
      </div>
    </div>

    <DataTable
      :page="props.page"
      :page-size="props.pageSize"
      :columns="columns"
      :data="props.data"
      :total-items="props.totalItems"
      :is-loading="props.isFetching"
      @update:page="emit('update:page', $event)"
      @update:page-size="emit('update:pageSize', $event)"
    />
  </div>
</template>
