import type { SystemConfig } from '@/types/system-config'
import systemAPI from '@/api/system'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'

export function useQuerySystemConfig() {
  return useQuery({
    queryKey: ['systemConfig'],
    queryFn: () => systemAPI.getSystemConfig(),
  })
}
export function useMutationSystemConfig() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (configData: SystemConfig) => systemAPI.updateSystemConfig(configData),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['systemConfig'] }),
  })
}

export function useMutationSystemRestart() {
  return useMutation({
    mutationFn: () => systemAPI.restartSystem(),
  })
}
