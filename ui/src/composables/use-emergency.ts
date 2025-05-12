import emergencyAPI from '@/api/emergency'
import { useMutation, useQuery } from '@tanstack/vue-query'

export const EMERGENCY_STATE_QUERY_KEY = 'emergencyState'

export function useEmergencyStateQuery() {
  return useQuery({
    queryKey: [EMERGENCY_STATE_QUERY_KEY],
    queryFn: emergencyAPI.getEmergencyState,
  })
}

export function useEmergencyStopMutation() {
  return useMutation({
    mutationFn: emergencyAPI.stopEmergency,
  })
}

export function useEmergencyResumeMutation() {
  return useMutation({
    mutationFn: emergencyAPI.resumeEmergency,
  })
}
