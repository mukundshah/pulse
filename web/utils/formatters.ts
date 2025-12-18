/**
 * Formats a duration given in microseconds into a readable string.
 *
 * For example:
 *  - < 1,000 µs: formatted as microseconds
 *  - < 1,000,000 µs: formatted as milliseconds
 *  - >= 1,000,000 µs: formatted as seconds
 *
 * @param {number} duration - The duration to format, in microseconds.
 * @returns {string} The formatted duration.
 */
export const formatDuration = (duration: number): string => {
  if (duration < 1000) {
    return duration.toLocaleString('en-US', {
      style: 'unit',
      unit: 'microsecond',
    })
  }
  if (duration < 1000000) {
    return (duration / 1000).toLocaleString('en-US', {
      style: 'unit',
      unit: 'millisecond',
      maximumFractionDigits: 0,
    })
  }

  return (duration / 1000000).toLocaleString('en-US', {
    style: 'unit',
    unit: 'second',
    maximumFractionDigits: 2,
  })
}
