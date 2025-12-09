<script setup lang="ts">
import type { NuxtError } from '#app'

const props = defineProps<{ error: NuxtError }>()

const {isAuthenticated} = useAuth()

const errorMessage = computed(() => {
  if (props.error.statusMessage) {
    return props.error.statusMessage
  }
  if (props.error.statusCode === 404) {
    return 'The page you are looking for does not exist'
  }
  if (props.error.statusCode === 500) {
    return 'An internal server error occurred'
  }
  return 'An unexpected error occurred'
})

const handleTryAgain = () => {
  clearError({ redirect: isAuthenticated.value ? '/dashboard' : '/' })
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center p-4">
    <div class="w-full max-w-lg text-center">
      <div class="mb-12">
        <h1 class="mb-3 text-6xl font-light tracking-tight text-foreground/60">
          {{ error.statusCode }}
        </h1>
        <p class="text-base text-muted-foreground">
          {{ errorMessage }}
        </p>
      </div>

      <div class="flex flex-col gap-3 sm:flex-row sm:justify-center items-center">
        <Button variant="ghost" @click="handleTryAgain">
          Try again
        </Button>
        <Button as-child variant="outline">
          <NuxtLink :to="isAuthenticated ? '/dashboard' : '/'">
            {{ isAuthenticated ? 'Go to dashboard' : 'Return home' }}
          </NuxtLink>
        </Button>
      </div>
    </div>
  </div>
</template>
