import type { AxiosRequestConfig } from 'axios'
import peripheralsAPI from '@/api/peripherals'
import { useQuery } from '@tanstack/vue-query'

export function useListAvailableSerialPorts(axiosOpts?: Partial<AxiosRequestConfig>) {
  return useQuery({
    queryKey: ['serial-ports'],
    queryFn: () => peripheralsAPI.listAvailableSerialPorts(axiosOpts),
    select: data => data.items,
  })
}
