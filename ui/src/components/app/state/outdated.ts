import type { Reactive } from 'vue'
import { getSecondsFromNow } from '@/lib/date'

export interface Outdated {
  isOutdated: boolean
  timeoutId: number | null
}

export function setOutdated(
  updatedAt: string,
  threshold: number,
  outdated: Reactive<Outdated>,
) {
  const diff = getSecondsFromNow(updatedAt)
  if (diff > threshold) {
    outdated.isOutdated = true
  }
  else {
    if (outdated.timeoutId) {
      clearTimeout(outdated.timeoutId)
    }

    outdated.timeoutId = setTimeout(() => {
      outdated.isOutdated = true
    }, (threshold - diff) * 1000)

    outdated.isOutdated = false
  }
}

export function clearOutdatedTimeout(outdated: Reactive<Outdated>) {
  if (outdated.timeoutId) {
    clearTimeout(outdated.timeoutId)
    outdated.timeoutId = null
  }
}
