import { defineNuxtPlugin } from '#imports'

export default defineNuxtPlugin({
  name: 'session',
  enforce: 'pre',
  async setup(nuxtApp) {
    const route = useRoute()
    const authNamespace = route.meta?.auth?.namespace

    const { syncAuthenticationStatus } = useAuth({ namespace: authNamespace })
    useAsyncData('auth-state', () => syncAuthenticationStatus())
  },
})
