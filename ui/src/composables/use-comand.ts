import type { CommandSort } from '@/api/command'
import type { SortPrefix } from '@/lib/sort'
import type { AxiosRequestConfig } from 'axios'
import command from '@/api/command'
import { keepPreviousData, useQuery } from '@tanstack/vue-query'

export function useComandInProgress(opts?: { axiosOpts?: Partial<AxiosRequestConfig>, refetchInterval: number }) {
  return useQuery({
    queryKey: ['comandInProgress'],
    queryFn: () => command.getCommandInProgress(opts?.axiosOpts),
    refetchInterval: opts?.refetchInterval,
  })
}

export function useListComands(page: Ref<number>, pageSize: Ref<number>, sorts: Ref<SortPrefix<CommandSort>[]>) {
  return useQuery({
    queryKey: ['comands', page, pageSize, sorts],
    queryFn: () => command.listCommands({ page: page.value, pageSize: pageSize.value, sorts: sorts.value }),
    placeholderData: keepPreviousData,
  })
}
