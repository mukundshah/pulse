<script setup lang="ts">
useHead({
  title: 'Dashboard',
})

useLayoutContext({
  actions: [],
  breadcrumbOverrides: computed(() => [
    undefined, // Root
    undefined, // Dashboard
  ]),
})

const { me } = useAuth()

const { data: user, pending: userPending } = useAsyncData('me', (_nuxtApp, { signal }) => me({ signal }), { dedupe: 'cancel' })

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Good morning'
  if (hour < 18) return 'Good afternoon'
  return 'Good evening'
})
</script>

<template>
  <div class="flex flex-1 flex-col gap-6 p-4 md:p-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="flex items-center gap-2 text-3xl font-semibold tracking-tight">
          <span>{{ greeting }}, </span>
          <Skeleton v-if="userPending" class="inline-block h-6 w-24" />
          <span v-else class="inline-block">{{ user?.name }}</span>
        </h1>
        <p class="text-muted-foreground">
          Overview of your monitoring infrastructure
        </p>
      </div>
    </div>
  </div>
</template>
