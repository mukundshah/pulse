export const TIME_FORMAT = {
  second: {
    hour: 'numeric',
    minute: 'numeric',
    second: 'numeric',
  },
  minute: {
    hour: 'numeric',
    minute: 'numeric',
  },
  hour: {
    day: 'numeric',
    month: 'short',
    hour: 'numeric',
    minute: 'numeric',
  },
  day: {
    day: 'numeric',
    month: 'short',
  },
  week: {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  },
} satisfies Record<string, Intl.DateTimeFormatOptions>
