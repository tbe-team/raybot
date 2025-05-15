import http from '@/lib/http'

const systemAPI = {
  reboot(): Promise<void> {
    return http.post('/system/reboot')
  },
  stopEmergency(): Promise<void> {
    return http.post('/system/stop-emergency')
  },
}
export default systemAPI
