import tailwindcss from '@tailwindcss/vite'
import { createResolver } from 'nuxt/kit'
import { env } from 'std-env'
import { joinURL } from 'ufo'

const { resolve } = createResolver(import.meta.url)

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  modules: [
    '@nuxt/eslint',
    '@nuxt/fonts',
    '@nuxt/icon',
    '@nuxtjs/seo',
    '@nuxtjs/color-mode',
    'shadcn-nuxt',
    'nuxt-gtag',
    'nuxt-open-fetch',
  ],

  srcDir: 'web',

  dir: {
    modules: resolve('./web/modules'),
  },

  app: {
    rootAttrs: {
      id: 'app',
    },
    head: {
      templateParams: {
        separator: '-',
      },
    },
  },

  site: {
    name: 'Pulse',
    url: 'https://pulse.app',
  },

  routeRules: {
    '/**': { appLayout: 'app', robots: false, auth: { required: true } },
    '/auth/login': { robots: true },
    '/auth/**': { appLayout: 'auth', robots: false, auth: { required: false, redirectIfLoggedIn: '/dashboard' } },
    '/': { appLayout: 'site', robots: true, auth: { required: false } },
    '/terms': { appLayout: 'site', robots: true, auth: { required: false } },
    '/privacy': { appLayout: 'site', robots: true, auth: { required: false } },
    '/contact': { appLayout: 'site', robots: true, auth: { required: false } },
    '/about': { appLayout: 'site', robots: true, auth: { required: false } },
  },

  css: ['assets/css/tailwind.css'],

  colorMode: {
    classSuffix: '',
  },

  eslint: {
    config: {
      standalone: false,
    },
  },

  shadcn: {
    prefix: '',
  },

  openFetch: {
    clients: {
      pulseAPI: {
        baseURL: 'http://localhost:8080/api',
        schema: joinURL(env.NUXT_PUBLIC_OPEN_FETCH_PULSE_API_BASE_URL ?? 'http://localhost:8080', '/docs/internal/openapi.json'),
      },
    },
  },

  vite: {
    plugins: [
      tailwindcss(),
    ],
    optimizeDeps: {
      include: [
        '@vee-validate/zod',
        '@vueuse/core',
        'class-variance-authority',
        'clsx',
        'lucide-vue-next',
        'reka-ui',
        'tailwind-merge',
        'vee-validate',
        'vue-sonner',
        'zod',
      ],
    },
  },

  nitro: {
    handlers: [
      {
        route: '/health',
        handler: '~/health.ts',
      },
    ],
  },
})
