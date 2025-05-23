import { useQuery } from '@tanstack/vue-query'
import versionAPI from '@/api/version'

export function useVersionQuery() {
  return useQuery({
    queryKey: ['version'],
    queryFn: versionAPI.getVersion,
  })
}
