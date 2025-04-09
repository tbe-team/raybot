import type { SerialPort } from '@/types/peripherals'
import type { AxiosRequestConfig } from 'axios'
import http from '@/lib/http'

const peripheralsAPI = {
  listAvailableSerialPorts: (axiosOpts?: Partial<AxiosRequestConfig>): Promise<{ items: SerialPort[] }> =>
    http.get('/peripherals/serials', axiosOpts),
}

export default peripheralsAPI
