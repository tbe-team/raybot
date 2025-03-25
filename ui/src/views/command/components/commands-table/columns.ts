import type { Command } from '@/types/command'
import type { ColumnDef } from '@tanstack/vue-table'
import DataTableSortableHeader from '@/components/shared/DataTableSortableHeader.vue'
import { formatDate } from '@/lib/date'
import { h } from 'vue'
import StatusBadge from './StatusBadge.vue'

export const columns: ColumnDef<Command>[] = [
  {
    accessorKey: 'id',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'ID' }),
    enableSorting: false,
  },
  {
    accessorKey: 'type',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Type' }),
  },
  {
    accessorKey: 'status',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Status' }),
    cell: ({ row }) => h(StatusBadge, { status: row.original.status }),
  },
  {
    accessorKey: 'source',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Source' }),
  },
  {
    accessorKey: 'inputs',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Inputs' }),
    cell: ({ row }) => {
      const inputs = row.getValue('inputs') as Record<string, unknown>
      return JSON.stringify(inputs)
    },
    enableSorting: false,
  },
  {
    accessorKey: 'error',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Error' }),
    cell: ({ row }) => {
      const error = row.getValue('error') as string | null
      return error || '-'
    },
    enableSorting: false,
  },
  {
    accessorKey: 'createdAt',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Created At' }),
    cell: ({ row }) => formatDate(row.getValue('createdAt')),
  },
  {
    accessorKey: 'completedAt',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Completed At' }),
    cell: ({ row }) => {
      const completedAt = row.getValue('completedAt') as string | null
      return completedAt ? formatDate(completedAt) : '-'
    },
  },
]
