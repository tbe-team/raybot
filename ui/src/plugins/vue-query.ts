import type { VueQueryPluginOptions } from '@tanstack/vue-query'
import { QueryCache, QueryClient } from '@tanstack/vue-query'
import { HttpStatusCode } from 'axios'
import { HTTPError } from '@/types/error'

const RETRY_LIMIT = 3
const STALE_TIME = 1000 * 60 * 5 // 5 minutes
const TIMEOUT_STATUS_CODES = [
  HttpStatusCode.RequestTimeout,
  HttpStatusCode.GatewayTimeout,
]

function handleNumRetries(failureCount: number, error: Error) {
  if (error instanceof HTTPError) {
    const isUnderRetryLimit = failureCount < RETRY_LIMIT
    const isTimeout = TIMEOUT_STATUS_CODES.includes(error.status)
    return isUnderRetryLimit && isTimeout
  }

  return false
}

const queryClient = new QueryClient({
  queryCache: new QueryCache({}),
  defaultOptions: {
    queries: {
      retry: handleNumRetries,
      refetchOnWindowFocus: false,
      staleTime: STALE_TIME,
    },
    mutations: {
      retry: handleNumRetries,
    },
  },
})

export const queryPluginOpts: VueQueryPluginOptions = {
  queryClient,
}
