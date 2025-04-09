import http from '@/lib/http'

const systemAPI = {
  restartSystem(): Promise<void> {
    return http.post('/system/restart')
  },
}
export default systemAPI
