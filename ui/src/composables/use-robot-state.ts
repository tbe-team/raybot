import type { AxiosRequestConfig } from 'axios'
import { useQuery } from '@tanstack/vue-query'
import robotStateAPI from '@/api/robot-state'

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
