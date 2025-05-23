import type { AxiosRequestConfig } from 'axios'
import type { CommandSort } from '@/api/commands'
import type { SortPrefix } from '@/lib/sort'
import { keepPreviousData, useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import commandsAPI from '@/api/commands'

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

export function useCancelProcessingCommandMutation() {
  return useMutation({
    mutationFn: commandsAPI.cancelProcessingCommand,
  })
}
export function useDeleteCommandMutation() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: commandsAPI.deleteCommand,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comands'] })
    },
  })
}
