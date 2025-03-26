<script setup lang="ts" generic="TData, TValue, TSort extends string">
import type { SortPrefix } from '@/lib/sort'
import type { ColumnDef, SortingState, VisibilityState } from '@tanstack/vue-table'

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

import { convertParamsToSorting, convertSortingToParams } from '@/lib/sort'
import { valueUpdater } from '@/lib/utils'
import {
  FlexRender,
  getCoreRowModel,
  useVueTable,
} from '@tanstack/vue-table'

import { LoaderCircleIcon } from 'lucide-vue-next'
import DataTablePagination from './DataTablePagination.vue'

interface Props {
  isLoading: boolean
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
  totalItems: number
  pageSizeOptions?: number[]
  sorts?: SortPrefix<TSort>[]
}
const props = withDefaults(defineProps<Props>(), {
  pageSizeOptions: () => [10, 20, 30],
  sorts: () => [],
})
const emit = defineEmits<{
  (e: 'sorts', sorts: SortPrefix<TSort>[]): void
}>()
const columnVisibility = ref<VisibilityState>({})
const page = defineModel<number>('page', { required: true })
const pageSize = defineModel<number>('pageSize', { required: true })

const sorting = computed<SortingState>({
  get: () => convertParamsToSorting(props.sorts ?? []),
  set: value => emit('sorts', convertSortingToParams(value)),
})

const table = useVueTable({
  get data() { return props.data },
  get columns() { return props.columns },
  get rowCount() { return props.totalItems },
  state: {
    get pagination() { return { pageIndex: page.value - 1, pageSize: pageSize.value } },
    get sorting() { return sorting.value },
    get columnVisibility() { return columnVisibility.value },
  },
  getCoreRowModel: getCoreRowModel(),

  // Server-side pagination
  manualPagination: true,
  onColumnVisibilityChange: updaterOrValue => valueUpdater(updaterOrValue, columnVisibility),
  onPaginationChange: (updaterOrValue) => {
    const newState = typeof updaterOrValue === 'function'
      ? updaterOrValue(table.getState().pagination)
      : updaterOrValue

    page.value = newState.pageIndex + 1
    pageSize.value = newState.pageSize
  },

  // Server-side sorting
  manualSorting: true,
  enableMultiSort: true,
  isMultiSortEvent: (event: unknown) => (event as MouseEvent).shiftKey,
  onSortingChange: (updaterOrValue) => {
    const newState = typeof updaterOrValue === 'function'
      ? updaterOrValue(table.getState().sorting)
      : updaterOrValue

    sorting.value = newState
  },
})

defineExpose({
  table,
})
</script>

<template>
  <div class="flex flex-col justify-end pb-1 space-y-4">
    <div class="h-full overflow-y-auto border rounded-md">
      <Table>
        <TableHeader class="bg-muted-foreground/10">
          <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
            <TableHead v-for="header in headerGroup.headers" :key="header.id">
              <FlexRender
                v-if="!header.isPlaceholder"
                :render="header.column.columnDef.header"
                :props="header.getContext()"
              />
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-if="props.isLoading">
            <TableRow>
              <TableCell
                :colspan="columns.length"
                class="relative h-24"
              >
                <div class="absolute inset-0 flex items-center justify-center">
                  <LoaderCircleIcon
                    class="w-8 h-8 animate-spin text-primary"
                    aria-label="Loading..."
                  />
                </div>
              </TableCell>
            </TableRow>
          </template>
          <template v-else-if="table.getRowModel().rows?.length">
            <TableRow
              v-for="row in table.getRowModel().rows"
              :key="row.id"
              :data-state="row.getIsSelected() ? 'selected' : 'undefined'"
            >
              <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
                <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
              </TableCell>
            </TableRow>
          </template>
          <template v-else>
            <TableRow>
              <TableCell :colspan="columns.length" class="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          </template>
        </TableBody>
      </Table>
    </div>

    <DataTablePagination
      v-if="!props.isLoading"
      v-model:page-size="pageSize"
      :total-items="props.totalItems ?? 0"
      :page-size-options="props.pageSizeOptions"
      :table="table"
    />
  </div>
</template>
