<script setup lang="ts">
import type { AlarmStatus } from '@/api/alarm'
import { AlertCircle, Loader, RefreshCw, Settings } from 'lucide-vue-next'
import DataTable from '@/components/shared/DataTable.vue'
import { Button } from '@/components/ui/button'
import { useAlarmsQuery } from '@/composables/use-alarm'
import { columns } from './alarm-table'

const route = useRoute()
const router = useRouter()

const page = defineModel<number>('page', { required: true })
const pageSize = defineModel<number>('pageSize', { required: true })
const status = ref<AlarmStatus>('ACTIVE')

const { data, isPending, isFetching, isError, error, refetch } = useAlarmsQuery(page, pageSize, status)

function handlePageChange(p: number) {
  page.value = p
  router.replace({ query: { ...route.query, page: p.toString() } })
}

function handlePageSizeChange(ps: number) {
  pageSize.value = ps
  page.value = 1
  router.replace({ query: { ...route.query, pageSize: ps.toString(), page: '1' } })
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

  <div v-else-if="isError" class="flex flex-col gap-4 justify-center items-center pt-20">
    <div class="flex flex-col gap-4 items-center p-6 text-red-500">
      <AlertCircle class="size-8" />
      <div class="space-y-2 text-center">
        <h2 class="text-lg font-semibold">
          Failed to load alarms
        </h2>
        <p class="text-sm text-muted-foreground">
          {{ error?.message || 'An unexpected error occurred' }}
        </p>
      </div>
    </div>
  </div>

  <div v-else-if="!data" class="flex flex-col gap-4 justify-center items-center pt-20">
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
          :disabled="isFetching"
          @click="() => refetch()"
        >
          <RefreshCw
            class="mr-2 size-4"
            :class="{ 'animate-spin': isFetching }"
          />
          Refresh
        </Button>
      </div>
    </div>

    <DataTable
      :page="page"
      :page-size="pageSize"
      :columns="columns"
      :data="data.items"
      :total-items="data.totalItems"
      :is-loading="isPending"
      @update:page="handlePageChange"
      @update:page-size="handlePageSizeChange"
    />
  </div>
</template>
