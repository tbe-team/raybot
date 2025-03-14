import robotStateAPI from '@/api/robot-state'
import { useQuery } from '@tanstack/vue-query'

export function useQueryRobotState() {
  return useQuery({
    queryKey: ['robotState'],
    queryFn: () => robotStateAPI.getRobotState(),
    refetchInterval: 1000,
  })
}
