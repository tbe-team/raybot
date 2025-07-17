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
    accessorKey: 'time',
    header: () => h('div', { class: 'text-xs' }, 'Time'),
    cell: ({ row }) => h('div', { class: 'flex flex-col' }, [
      h('span', null, formatDate(row.original.activatedAt)),
      h('span', null, row.original.deactivatedAt ? formatDate(row.original.deactivatedAt) : '-'),
    ]),
  },
  {
    accessorKey: 'duration',
    header: () => h('div', { class: 'text-xs' }, 'Duration'),
    cell: ({ row }) => formatDuration(row.original.activatedAt),
  },
]
