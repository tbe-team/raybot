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

export function useComand(page: Ref<number>, pageSize: Ref<number>) {
  return useQuery({
    queryKey: ['comand', page, pageSize],
    queryFn: () => command.listCommands({ page: page.value, pageSize: pageSize.value }),
    placeholderData: keepPreviousData,
  })
}
