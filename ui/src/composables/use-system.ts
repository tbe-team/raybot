import systemAPI from '@/api/system'
import { useMutation } from '@tanstack/vue-query'

export function useSystemRestartMutation() {
  return useMutation({
    mutationFn: systemAPI.restartSystem,
  })
}
