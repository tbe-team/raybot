import type { AxiosRequestConfig } from 'axios'
import { useQuery } from '@tanstack/vue-query'
import peripheralsAPI from '@/api/peripherals'

export function useListAvailableSerialPortsQuery(axiosOpts?: Partial<AxiosRequestConfig>) {
  return useQuery({
    queryKey: ['serial-ports'],
    queryFn: () => peripheralsAPI.listAvailableSerialPorts(axiosOpts),
    select: data => data.items,
  })
}
