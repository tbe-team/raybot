<script setup lang="ts">
import type { Alarm } from '@/types/alarm'
import { AlertCircle, Loader, RefreshCw, Settings } from 'lucide-vue-next'
import { columns } from '@/components/app/alarm/alarm-table/alarm-active/columns'
import DataTable from '@/components/shared/DataTable.vue'
import { Button } from '@/components/ui/button'

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

const router = useRouter()
</script>

<template>
  <div v-if="props.isPending" class="flex flex-col gap-4 justify-center items-center pt-20">
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

  <div v-else-if="!props.data" class="flex flex-col gap-4 justify-center items-center pt-20">
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
          @click="() => router.push('/monitoring')"
        >
          <Settings class="size-4" />
          Settings
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
