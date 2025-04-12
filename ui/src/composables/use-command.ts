import type { CommandSort } from '@/api/commands'
import type { SortPrefix } from '@/lib/sort'
import type { AxiosRequestConfig } from 'axios'
import commandsAPI from '@/api/commands'
import { keepPreviousData, useMutation, useQuery } from '@tanstack/vue-query'

export function useCurrentProcessingCommandQuery(
  opts?: { axiosOpts?: Partial<AxiosRequestConfig> },
) {
  return useQuery({
    queryKey: ['currentProcessingCommand'],
    queryFn: () => commandsAPI.getCurrentProcessingCommand(opts?.axiosOpts),
  })
}

export function useListQueuedCommandsQuery(
  page: Ref<number>,
  pageSize: Ref<number>,
  opts?: { axiosOpts?: Partial<AxiosRequestConfig> },
) {
  return useQuery({
    queryKey: ['queuedComand', page, pageSize],
    queryFn: () => commandsAPI.listCommands({
      page: page.value,
      pageSize: pageSize.value,
      sorts: ['created_at'],
      statuses: ['QUEUED'],
    }, opts?.axiosOpts),
  })
}

export function useListComandsQuery(
  page: Ref<number>,
  pageSize: Ref<number>,
  sorts: Ref<SortPrefix<CommandSort>[]>,
) {
  return useQuery({
    queryKey: ['comands', page, pageSize, sorts],
    queryFn: () => commandsAPI.listCommands({
      page: page.value,
      pageSize: pageSize.value,
      sorts: sorts.value,
    }),
    placeholderData: keepPreviousData,
  })
}

export function useGetCommandQuery(
  id: Ref<number>,
  opts?: { axiosOpts?: Partial<AxiosRequestConfig> },
) {
  return useQuery({
    queryKey: ['command', id],
    queryFn: () => commandsAPI.getCommand(id.value, opts?.axiosOpts),
  })
}

export function useCreateCommandMutation() {
  return useMutation({
    mutationFn: commandsAPI.createCommand,
  })
}
