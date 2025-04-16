import systemAPI from '@/api/system'
import { useMutation } from '@tanstack/vue-query'

export function useSystemRebootMutation() {
  return useMutation({
    mutationFn: systemAPI.reboot,
  })
}
