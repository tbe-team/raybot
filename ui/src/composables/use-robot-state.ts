import type { AxiosRequestConfig } from 'axios'
import robotStateAPI from '@/api/robot-state'
import { useQuery } from '@tanstack/vue-query'

export function useQueryRobotState(opt?: Partial<AxiosRequestConfig>, refetchInterval?: number) {
  return useQuery({
    queryKey: ['robotState'],
    queryFn: () => robotStateAPI.getRobotState(opt),
    refetchInterval,
  })
}
