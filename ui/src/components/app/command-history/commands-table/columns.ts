import type { Command } from '@/types/command'
import type { ColumnDef } from '@tanstack/vue-table'
import DataTableSortableHeader from '@/components/shared/DataTableSortableHeader.vue'
import { formatDate } from '@/lib/date'
import { h } from 'vue'
import { getCommandName } from '../../command-queue/utils'
import CommandActionDropdownMenu from './CommandActionDropdownMenu.vue'
import SourceBadge from './SourceBadge.vue'
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
    cell: ({ row }) => getCommandName(row.original.type),
  },
  {
    accessorKey: 'status',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Status' }),
    cell: ({ row }) => h(StatusBadge, { key: row.original.id + row.original.status, status: row.original.status }),
  },
  {
    accessorKey: 'source',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Source' }),
    cell: ({ row }) => h(SourceBadge, { key: row.original.id + row.original.source, source: row.original.source }),
  },
  {
    accessorKey: 'inputs',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Inputs' }),
    cell: ({ row }) => {
      const inputs = row.original.inputs
      return h('pre', {
        class: 'text-xs overflow-auto max-h-32 whitespace-pre-wrap',
        style: 'max-width: 300px;',
      }, JSON.stringify(inputs, null, 2))
    },
    enableSorting: false,
  },
  {
    accessorKey: 'error',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Error' }),
    cell: ({ row }) => {
      const error = row.original.error
      if (!error) {
        return '-'
      }
      else {
        return h('div', {
          class: 'text-red-500 truncate',
          style: 'max-width: 300px;',
          title: error,
        }, error)
      }
    },
    enableSorting: false,
  },
  {
    accessorKey: 'createdAt',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Created At' }),
    cell: ({ row }) => formatDate(row.original.createdAt),
  },
  {
    accessorKey: 'completedAt',
    header: ({ column }) => h(DataTableSortableHeader<Command>, { column, title: 'Completed At' }),
    cell: ({ row }) => {
      const completedAt = row.original.completedAt
      return completedAt ? formatDate(completedAt) : '-'
    },
  },
  {
    accessorKey: 'action',
    header: 'Action',
    cell: ({ row }) => h(CommandActionDropdownMenu, { command: row.original }),
  },
]
