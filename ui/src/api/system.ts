import http from '@/lib/http'

const systemAPI = {
  reboot(): Promise<void> {
    return http.post('/system/reboot')
  },
}
export default systemAPI
