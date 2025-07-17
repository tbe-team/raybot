import type { AxiosRequestConfig } from 'axios'
import type { AlarmStatus } from '@/api/alarm'
import { useMutation, useQuery } from '@tanstack/vue-query'
import alarmAPI from '@/api/alarm'

export const COUNT_ACTIVE_ALARMS_QUERY_KEY = 'countActiveAlarms'
export const LIST_ALARMS_QUERY_KEY = 'listAlarms'

export function useCountActiveAlarmsQuery(opts?: { axiosOpts?: Partial<AxiosRequestConfig>, refetchInterval?: number }) {
  return useQuery({
    queryKey: [COUNT_ACTIVE_ALARMS_QUERY_KEY],
    queryFn: () => alarmAPI.listAlarms({ page: 1, pageSize: 1, status: 'ACTIVE' }, opts?.axiosOpts),
    refetchInterval: opts?.refetchInterval,
    select: data => ({ count: data.totalItems }),
  })
}

export function useListAlarmsQuery(page: Ref<number>, pageSize: Ref<number>, status: Ref<AlarmStatus>, opts?: { axiosOpts?: Partial<AxiosRequestConfig>, refetchInterval?: number }) {
  return useQuery({
    queryKey: [LIST_ALARMS_QUERY_KEY, page, pageSize, status],
    queryFn: () => alarmAPI.listAlarms({
      page: page.value,
      pageSize: pageSize.value,
      status: status.value,
    }, opts?.axiosOpts),
    refetchInterval: opts?.refetchInterval,
  })
}

export function useDeleteDeactiveAlarmsMutation(opts?: { axiosOpts?: Partial<AxiosRequestConfig> }) {
  return useMutation({
    mutationFn: () => alarmAPI.deleteDeactiveAlarms(opts?.axiosOpts),
  })
}
