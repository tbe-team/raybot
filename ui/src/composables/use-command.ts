import type { CommandSort, CreateCommandParams } from '@/api/commands'
import type { SortPrefix } from '@/lib/sort'
import type { AxiosRequestConfig } from 'axios'
import commandsAPI from '@/api/commands'
import { keepPreviousData, useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import type { CommandType } from '@/types/command'


export function useListProcessingCommandsQuery(opts?: { axiosOpts?: Partial<AxiosRequestConfig>, refetchInterval: number }) {
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

export function useGetCommandQuery(id:ComputedRef<number>, opts?: { axiosOpts?: Partial<AxiosRequestConfig>, refetchInterval: number}) {
  return useQuery({
    queryKey: ['command', id],
    queryFn: () => commandsAPI.getCommand(id.value, opts?.axiosOpts) ,
    refetchInterval: opts?.refetchInterval,
  })
}

export function useCreateCommandMutation() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (p : CreateCommandParams<CommandType>) => commandsAPI.createCommand(p),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['queuedComand'] })
    },
  })
}
