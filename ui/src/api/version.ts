import type { Version } from '@/types/version'
import http from '@/lib/http'

const versionAPI = {
  getVersion(): Promise<Version> {
    return http.get('/version')
  },
}

export default versionAPI
