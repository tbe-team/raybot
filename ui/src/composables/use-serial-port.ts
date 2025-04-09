import type { AxiosRequestConfig } from 'axios'
import peripheralsAPI from '@/api/peripherals'
import { useQuery } from '@tanstack/vue-query'

export function useSerialPort(axiosOpts?: Partial<AxiosRequestConfig>) {
  return useQuery({
    queryKey: ['serial-ports'],
    queryFn: () => peripheralsAPI.getPorts(axiosOpts),
    select: data => data.items,
  })
}
