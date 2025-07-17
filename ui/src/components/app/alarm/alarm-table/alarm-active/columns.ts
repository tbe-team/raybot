import type { ColumnDef } from '@tanstack/vue-table'
import type { Alarm } from '@/types/alarm'
import dayjs from 'dayjs'
import { formatDate, formatUptimeShort } from '@/lib/date'

function formatDuration(dateString: string): string {
  const duration = dayjs().diff(dayjs(dateString), 'second')
  return formatUptimeShort(duration)
}

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
