import type { ColumnDef } from '@tanstack/vue-table'
import type { Alarm } from '@/types/alarm'
import { formatDate, formatDuration } from '@/lib/date'

export const columns: ColumnDef<Alarm>[] = [
  {
    accessorKey: 'name',
    header: () => h('div', { class: 'text-xs' }, 'Name'),
    cell: ({ row }) => row.original.type,
  },
  {
    accessorKey: 'message',
    header: () => h('div', { class: 'text-xs' }, 'Message'),
    cell: ({ row }) => h('div', { class: 'max-w-96' }, row.original.message),
  },
  {
    accessorKey: 'activatedAt',
    header: () => h('div', { class: 'text-xs' }, 'Activated At'),
    cell: ({ row }) => formatDate(row.original.activatedAt),
  },
  {
    accessorKey: 'duration',
    header: () => h('div', { class: 'text-xs' }, 'Duration'),
    cell: ({ row }) => formatDuration(row.original.activatedAt),
  },
]
