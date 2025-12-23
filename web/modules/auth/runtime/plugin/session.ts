import { defineNuxtPlugin } from '#imports'

export default defineNuxtPlugin({
  name: 'session',
  enforce: 'pre',
  async setup() {
    const route = useRoute()
    const authNamespace = route.meta?.auth?.namespace

    const { syncAuthenticationStatus } = useAuth({ namespace: authNamespace })

    await callOnce(async () => {
      await syncAuthenticationStatus()
    })
  },
})
