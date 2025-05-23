import { useMutation } from '@tanstack/vue-query'
import systemAPI from '@/api/system'

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
