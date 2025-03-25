<script setup lang="ts" generic="TData">
import type { Table } from '@tanstack/vue-table'
import {
  Pagination,
  PaginationFirst,
  PaginationLast,
  PaginationList,
  PaginationNext,
  PaginationPrev,
} from '@/components/ui/pagination'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

interface Props {
  table: Table<TData>
  pageSizeOptions: number[]
}

const props = defineProps<Props>()

const pageSize = defineModel<number>('pageSize', { required: true })
const pageSizeStr = computed<string>({
  get: () => pageSize.value.toString(),
  set: value => pageSize.value = Number(value),
})

const currentPage = computed(() => props.table.getState().pagination.pageIndex + 1)
const totalPages = computed(() => props.table.getPageCount())
const totalItems = computed(() => props.table.getRowCount())
const recordRange = computed(() => {
  const from = (currentPage.value - 1) * pageSize.value + 1
  const to = Math.min(currentPage.value * pageSize.value, totalItems.value)
  return `${from} - ${to} of ${totalItems.value}`
})
</script>

<template>
  <div class="flex ml-auto space-x-6">
    <div class="flex items-center gap-1">
      <span class="text-sm sr-only sm:not-sr-only">
        Items per page:
      </span>
      <Select v-model="pageSizeStr">
        <SelectTrigger class="w-20">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectItem
              v-for="option in props.pageSizeOptions"
              :key="option"
              :value="option.toString()"
            >
              {{ option }}
            </SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>
    </div>

    <div class="flex items-center space-x-4">
      <span class="text-sm">
        {{ recordRange }}
      </span>
      <Pagination
        :items-per-page="pageSize"
        :total="totalItems"
      >
        <PaginationList class="flex items-center gap-1">
          <PaginationFirst
            :disabled="!props.table.getCanPreviousPage()"
            @click="props.table.setPageIndex(0)"
          />
          <PaginationPrev
            :disabled="!props.table.getCanPreviousPage()"
            @click="props.table.previousPage()"
          />
          <PaginationNext
            :disabled="!props.table.getCanNextPage()"
            @click="props.table.nextPage()"
          />
          <PaginationLast
            :disabled="!props.table.getCanNextPage()"
            @click="props.table.setPageIndex(totalPages - 1)"
          />
        </PaginationList>
      </Pagination>
    </div>
  </div>
</template>
