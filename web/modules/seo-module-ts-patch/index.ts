import { defineNuxtModule } from 'nuxt/kit'

export default defineNuxtModule({
  meta: {
    name: 'seo-module-ts-patch',
    configKey: undefined,
  },
  async setup(options, nuxt) {
    nuxt.hook('prepare:types', ({ nodeReferences }) => {
      nodeReferences.push(
        ...[
          'module/nuxt-site-config.d.ts',
          'module/nuxt-robots.d.ts',
          'module/nuxt-sitemap.d.ts',
          'module/nuxt-og-image.d.ts',
          'module/nuxt-schema-org.d.ts',
          'module/nuxt-seo-utils.assets.d.ts',
          'module/nuxt-seo-utils.d.ts',
        ].map(path => ({ path }))
      )
    })
  },
})
