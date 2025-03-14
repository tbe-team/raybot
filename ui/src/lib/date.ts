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
