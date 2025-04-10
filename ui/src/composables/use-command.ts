import type { CommandSort } from '@/api/commands'
import type { SortPrefix } from '@/lib/sort'
import type { AxiosRequestConfig } from 'axios'
import commandsAPI from '@/api/commands'
import { keepPreviousData, useQuery } from '@tanstack/vue-query'

export function useProcessingComandQuery(opts?: { axiosOpts?: Partial<AxiosRequestConfig>, refetchInterval: number }) {
  return useQuery({
    queryKey: ['processingComand'],
    queryFn: () => commandsAPI.listCommands({ statuses: ['PROCESSING'] }, opts?.axiosOpts),
    refetchInterval: opts?.refetchInterval,
  })
}

export function useQueuedComandQuery(page: Ref<number>, pageSize: Ref<number>, opts?: { axiosOpts?: Partial<AxiosRequestConfig> }) {
  return useQuery({
    queryKey: ['queuedComand', page, pageSize],
    queryFn: () => commandsAPI.listCommands({ page: page.value, pageSize: pageSize.value, sorts: ['created_at'], statuses: ['QUEUED'] }, opts?.axiosOpts),
  })
}

export function useListComandsQuery(page: Ref<number>, pageSize: Ref<number>, sorts: Ref<SortPrefix<CommandSort>[]>) {
  return useQuery({
    queryKey: ['comands', page, pageSize, sorts],
    queryFn: () => commandsAPI.listCommands({ page: page.value, pageSize: pageSize.value, sorts: sorts.value }),
    placeholderData: keepPreviousData,
  })
}
