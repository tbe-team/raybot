<script setup lang="ts">
import type { CommandSort } from '@/api/commands'
import type { SortPrefix } from '@/lib/sort'
import type { Command } from '@/types/command'
import type { Table } from '@tanstack/vue-table'
import CommandDetailSheet from '@/components/app/command-history/CommandDetailSheet.vue'
import { columns } from '@/components/app/command-history/commands-table'
import DataTable from '@/components/shared/DataTable.vue'
import DataTableColumnVisibility from '@/components/shared/DataTableColumnVisibility.vue'
import PageContainer from '@/components/shared/PageContainer.vue'
import { Button } from '@/components/ui/button'
import { useListComandsQuery } from '@/composables/use-command'
import { AlertCircle, Loader, RefreshCw } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const tableRef = useTemplateRef<{ table: Table<Command> } | null>('table')
const selectedCommand = ref<Command | null>(null)
const isCommandDetailOpen = ref(false)

const page = ref(Number(route.query.page) || 1)
const pageSize = ref(Number(route.query.pageSize) || 10)
const sorts = ref<SortPrefix<CommandSort>[]>(
  route.query.sorts ? (route.query.sorts as string).split(',') as SortPrefix<CommandSort>[] : [],
)

const { data, isPending, isFetching, isError, error, refetch } = useListComandsQuery(page, pageSize, sorts)

function handleSortingChange(s: SortPrefix<CommandSort>[]) {
  sorts.value = s
  router.replace({
    query: {
      ...route.query,
      sorts: s.length ? s.join(',') : undefined,
    },
  })
}

function handlePageChange(p: number) {
  page.value = p
  router.replace({ query: { ...route.query, page: p.toString() } })
}

function handlePageSizeChange(ps: number) {
  pageSize.value = ps
  page.value = 1
  router.replace({ query: { ...route.query, pageSize: ps.toString(), page: '1' } })
}

function handleSelectCommand(row: Command) {
  selectedCommand.value = row
  isCommandDetailOpen.value = true
}
</script>

<template>
  <PageContainer>
    <div v-if="isPending" class="flex flex-col items-center justify-center gap-4 pt-20">
      <div class="flex items-center gap-4">
        <Loader class="w-8 h-8 animate-spin text-muted-foreground" />
      </div>
      <p class="text-lg text-muted-foreground">
        Loading commands...
      </p>
    </div>

    <div v-else-if="isError" class="flex flex-col items-center justify-center gap-4 pt-20">
      <div class="flex flex-col items-center gap-4 p-6 text-destructive">
        <AlertCircle class="w-8 h-8" />
        <div class="space-y-2 text-center">
          <h2 class="text-lg font-semibold">
            Failed to load commands
          </h2>
          <p class="text-sm text-muted-foreground">
            {{ error?.message || 'An unexpected error occurred' }}
          </p>
        </div>
      </div>
    </div>

    <div v-else-if="!data" class="flex flex-col items-center justify-center gap-4 pt-20">
      <div class="flex flex-col items-center gap-4 p-6">
        <AlertCircle class="w-8 h-8 text-muted-foreground" />
        <div class="space-y-2 text-center">
          <h2 class="text-lg font-semibold">
            No commands found
          </h2>
          <p class="text-sm text-muted-foreground">
            There are no commands to display
          </p>
        </div>
      </div>
    </div>

    <div v-else class="flex flex-col w-full">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-xl font-semibold tracking-tight">
            Commands
          </h1>
          <p class="text-sm text-muted-foreground">
            View and manage robot commands
          </p>
        </div>
        <div class="flex items-center gap-2">
          <DataTableColumnVisibility
            v-if="tableRef?.table"
            :table="tableRef.table"
          />
          <Button
            variant="outline"
            :disabled="isFetching || !sorts.length"
            @click="() => tableRef?.table?.resetSorting()"
          >
            Clear Sort
          </Button>
          <Button
            variant="outline"
            :disabled="isFetching"
            @click="() => refetch()"
          >
            <RefreshCw
              class="w-4 h-4 mr-2"
              :class="{ 'animate-spin': isFetching }"
            />
            Refresh
          </Button>
        </div>
      </div>

      <DataTable
        ref="table"
        :initial-state="{ columnVisibility: { error: false } }"
        :page="page"
        :page-size="pageSize"
        :columns="columns"
        :data="data.items"
        :total-items="data.totalItems"
        :is-loading="isPending"
        :sorts="sorts"
        @sorts="handleSortingChange"
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
        @row-click="handleSelectCommand"
      />
    </div>
    <CommandDetailSheet
      v-if="selectedCommand"
      v-model:is-open="isCommandDetailOpen"
      :command="selectedCommand"
    />
  </PageContainer>
</template>
