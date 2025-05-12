import emergencyAPI from "@/api/emergency"
import { useMutation, useQuery } from "@tanstack/vue-query"

export const EMERGENCY_STATE_QUERY_KEY = 'emergencyState'

export const useEmergencyStateQuery = () => {
  return useQuery({
    queryKey: [EMERGENCY_STATE_QUERY_KEY],
    queryFn: emergencyAPI.getEmergencyState,
  })
}

export const useEmergencyStopMutation = () => {
  return useMutation({
    mutationFn: emergencyAPI.stopEmergency,
  })
}

export const useEmergencyResumeMutation = () => {
  return useMutation({
    mutationFn: emergencyAPI.resumeEmergency,
  })
}
