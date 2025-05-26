import type { AxiosRequestConfig } from 'axios'
import type { CommandSort } from '@/api/commands'
import type { SortPrefix } from '@/lib/sort'
import { keepPreviousData, useMutation, useQuery } from '@tanstack/vue-query'
import commandsAPI from '@/api/commands'

export const COMMAND_QUEUE_QUERY_KEY = 'queuedCommand'
export const CURRENT_PROCESSING_COMMAND_QUERY_KEY = 'currentProcessingCommand'
export const COMMANDS_QUERY_KEY = 'commands'
export const COMMAND_QUERY_KEY = 'command'

export function useCurrentProcessingCommandQuery(
  opts?: { axiosOpts?: Partial<AxiosRequestConfig>, refetchInterval?: number },
) {
  return useQuery({
    queryKey: [CURRENT_PROCESSING_COMMAND_QUERY_KEY],
    queryFn: () => commandsAPI.getCurrentProcessingCommand(opts?.axiosOpts),
    refetchInterval: opts?.refetchInterval,
  })
}

export function useListQueuedCommandsQuery(
  page: Ref<number>,
  pageSize: Ref<number>,
  opts?: { axiosOpts?: Partial<AxiosRequestConfig>, refetchInterval?: number },
) {
  return useQuery({
    queryKey: [COMMAND_QUEUE_QUERY_KEY, page, pageSize],
    queryFn: () => commandsAPI.listCommands({
      page: page.value,
      pageSize: pageSize.value,
      sorts: ['created_at'],
      statuses: ['QUEUED'],
    }, opts?.axiosOpts),
    refetchInterval: opts?.refetchInterval,
  })
}

export function useListComandsQuery(
  page: Ref<number>,
  pageSize: Ref<number>,
  sorts: Ref<SortPrefix<CommandSort>[]>,
) {
  return useQuery({
    queryKey: [COMMANDS_QUERY_KEY, page, pageSize, sorts],
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
    queryKey: [COMMAND_QUERY_KEY, id],
    queryFn: () => commandsAPI.getCommand(id.value, opts?.axiosOpts),
  })
}

export function useCreateCommandMutation() {
  return useMutation({
    mutationFn: commandsAPI.createCommand,
  })
}

export function useCancelProcessingCommandMutation() {
  return useMutation({
    mutationFn: commandsAPI.cancelProcessingCommand,
  })
}
export function useDeleteCommandMutation() {
  return useMutation({
    mutationFn: commandsAPI.deleteCommand,
  })
}
