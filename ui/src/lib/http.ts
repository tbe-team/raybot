import type { InternalAxiosRequestConfig } from 'axios'
import { HTTPError, RaybotError } from '@/types/error'
import axios, { isAxiosError } from 'axios'
import { useNProgress } from './nprogress'

const nprogress = useNProgress()

const instance = axios.create({
  baseURL: '/api/v1',
  timeout: 20000,
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true,
})

instance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    if (!config.doNotShowLoading) {
      nprogress.start()
    }
    return config
  },
  (err) => {
    return Promise.reject(err)
  },
)

instance.interceptors.response.use(
  (response) => {
    nprogress.done()
    return Promise.resolve(response.data)
  },
  (error) => {
    nprogress.done()

    if (isAxiosError(error) && error.response && error.response.data) {
      const { status } = error.response
      const { message, code, details } = error.response.data

      // TODO: Should handle 401, 403 here
      if (status >= 400 && status < 500) {
        return Promise.reject(new RaybotError(message, status, code, details))
      }

      return Promise.reject(new HTTPError(message || 'Unexpected error', status))
    }

    return Promise.reject(error)
  },
)

export default instance
