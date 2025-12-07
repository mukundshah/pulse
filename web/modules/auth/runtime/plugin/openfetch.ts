import { useAuth } from '../composable'

export default defineNuxtPlugin({
  name: 'api-auth',
  enforce: 'pre',
  setup: (nuxtApp) => {

    nuxtApp.hook('openFetch:onRequest', (ctx) => {
      if (!ctx.options.headers.get('Authorization')) {
        const { token } = useAuth()
        if (token) {
          ctx.options.headers.set('Authorization', `Bearer ${token}`)
        }
      }
    })
  }
})
