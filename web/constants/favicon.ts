import type { ResolvableArray, ResolvableLink } from '@unhead/vue'

export const FAVICONS: Record<string, ResolvableArray<ResolvableLink>> = {
  passing: [
    {
      rel: 'icon',
      href: '/favicons/icon-passed-light.svg',
      type: 'image/svg+xml',
      sizes: 'any',
      media: '(prefers-color-scheme: light)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-passed-light.png',
      sizes: '48x48',
      type: 'image/png',
      media: '(prefers-color-scheme: light)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-passed-dark.svg',
      type: 'image/svg+xml',
      sizes: 'any',
      media: '(prefers-color-scheme: dark)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-passed-dark.png',
      sizes: '48x48',
      type: 'image/png',
      media: '(prefers-color-scheme: dark)',
    },
  ],
  degraded: [
    {
      rel: 'icon',
      href: '/favicons/icon-degraded-light.svg',
      type: 'image/svg+xml',
      sizes: 'any',
      media: '(prefers-color-scheme: light)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-degraded-light.png',
      sizes: '48x48',
      type: 'image/png',
      media: '(prefers-color-scheme: light)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-degraded-dark.svg',
      type: 'image/svg+xml',
      sizes: 'any',
      media: '(prefers-color-scheme: dark)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-degraded-dark.png',
      sizes: '48x48',
      type: 'image/png',
      media: '(prefers-color-scheme: dark)',
    },
  ],
  failing: [
    {
      rel: 'icon',
      href: '/favicons/icon-failed-light.svg',
      type: 'image/svg+xml',
      sizes: 'any',
      media: '(prefers-color-scheme: light)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-failed-light.png',
      sizes: '48x48',
      type: 'image/png',
      media: '(prefers-color-scheme: light)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-failed-dark.svg',
      type: 'image/svg+xml',
      sizes: 'any',
      media: '(prefers-color-scheme: dark)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-failed-dark.png',
      sizes: '48x48',
      type: 'image/png',
      media: '(prefers-color-scheme: dark)',
    },
  ],
  unknown: [
    {
      rel: 'icon',
      href: '/favicons/icon-unknown-light.svg',
      type: 'image/svg+xml',
      sizes: 'any',
      media: '(prefers-color-scheme: light)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-unknown-light.png',
      sizes: '48x48',
      type: 'image/png',
      media: '(prefers-color-scheme: light)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-unknown-dark.svg',
      type: 'image/svg+xml',
      media: '(prefers-color-scheme: dark)',
    },
    {
      rel: 'icon',
      href: '/favicons/icon-unknown-dark.png',
      sizes: '48x48',
      type: 'image/png',
      media: '(prefers-color-scheme: dark)',
    },
  ],
} as const
