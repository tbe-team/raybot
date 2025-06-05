import type { AxiosRequestConfig } from 'axios'
import { useQuery } from '@tanstack/vue-query'
import limitSwitchStateAPI from '@/api/limit-switch-state'

export function useLimitSwitchStateQuery(opts?: {
  axiosOpts?: Partial<AxiosRequestConfig>
  refetchInterval?: Ref<number>
}) {
  return useQuery({
    queryKey: ['limit-switch-state'],
    queryFn: () => limitSwitchStateAPI.getLimitSwitchState(opts?.axiosOpts),
    refetchInterval: opts?.refetchInterval,
  })
}
