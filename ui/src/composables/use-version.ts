import versionAPI from '@/api/version'
import { useQuery } from '@tanstack/vue-query'

export function useVersionQuery() {
  return useQuery({
    queryKey: ['version'],
    queryFn: versionAPI.getVersion,
  })
}
