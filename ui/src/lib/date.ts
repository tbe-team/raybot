import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)
dayjs.extend(duration)
/**
 * Formats a date string into a human-readable string.
 * @param dateString - The date string to format.
 * @returns A string representing the formatted date.
 */
export function formatDate(dateString: string): string {
  if (!dateString)
    return 'N/A'

  const date = new Date(dateString)
  if (Number.isNaN(date.getTime()))
    return 'Invalid date'

  return new Intl.DateTimeFormat('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  }).format(date)
}

export function formatUptime(seconds: number): string {
  const days = Math.floor(seconds / (24 * 60 * 60))
  const hours = Math.floor((seconds % (24 * 60 * 60)) / (60 * 60))
  const minutes = Math.floor((seconds % (60 * 60)) / 60)
  const remainingSeconds = seconds % 60

  const parts = []
  if (days > 0)
    parts.push(`${days} day${days > 1 ? 's' : ''}`)
  if (hours > 0)
    parts.push(`${hours} hour${hours > 1 ? 's' : ''}`)
  if (minutes > 0)
    parts.push(`${minutes} minute${minutes > 1 ? 's' : ''}`)
  if (remainingSeconds > 0 || parts.length === 0)
    parts.push(`${remainingSeconds.toFixed(0)} second${remainingSeconds !== 1 ? 's' : ''}`)

  return parts.join(' ')
}
