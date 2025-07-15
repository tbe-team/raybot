import type { AxiosRequestConfig } from 'axios'
import { useMutation, useQuery } from '@tanstack/vue-query'
import systemAPI from '@/api/system'

export const SYSTEM_INFO_QUERY_KEY = 'system-info'
export const SYSTEM_STATUS_QUERY_KEY = 'system-status'

export function useSystemRebootMutation() {
  return useMutation({
    mutationFn: systemAPI.reboot,
  })
}

export function useSystemStopEmergencyMutation() {
  return useMutation({
    mutationFn: systemAPI.stopEmergency,
  })
}

export function useSystemGetInfoQuery(opts?: {
  axiosOpts?: Partial<AxiosRequestConfig>
  refetchInterval?: number
}) {
  return useQuery({
    queryKey: [SYSTEM_INFO_QUERY_KEY],
    queryFn: () => systemAPI.getInfo(opts?.axiosOpts),
    refetchInterval: opts?.refetchInterval,
  })
}

export function useSystemStatusQuery(opts?: {
  axiosOpts?: Partial<AxiosRequestConfig>
  refetchInterval?: number
}) {
  return useQuery({
    queryKey: [SYSTEM_STATUS_QUERY_KEY],
    queryFn: () => systemAPI.getStatus({
      ...opts?.axiosOpts,

    }),
    refetchInterval: opts?.refetchInterval,

  })
}
