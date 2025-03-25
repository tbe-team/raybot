import type { AxiosRequestConfig } from 'axios'
import robotStateAPI from '@/api/robot-state'
import { useQuery } from '@tanstack/vue-query'

export function useQueryRobotState(
  opts?: {
    axiosOpts?: Partial<AxiosRequestConfig>
    refetchInterval?: Ref<number>
  },
) {
  return useQuery({
    queryKey: ['robotState'],
    queryFn: () => robotStateAPI.getRobotState(opts?.axiosOpts),
    refetchInterval: opts?.refetchInterval,
  })
}
