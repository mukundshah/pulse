<script setup lang="ts">
import { toast } from 'vue-sonner'

useHead({
  title: 'Verify email',
})

const route = useRoute()
const { verifyEmail } = useEmailManagement()

const isVerifying = ref(true)
const isSuccess = ref(false)
const error = ref<string | null>(null)

onMounted(async () => {
  const token = route.query.token as string | undefined

  if (!token) {
    error.value = 'Verification token is missing'
    isVerifying.value = false
    return
  }

  try {
    await verifyEmail({ token })
    isSuccess.value = true
    toast('Email verified successfully', {
      description: 'Your email has been verified. You can now sign in.',
    })
    // Redirect to login after a short delay
    setTimeout(() => {
      navigateTo('/auth/login')
    }, 2000)
  } catch (err: any) {
    error.value = err?.message || 'Failed to verify email. The token may be invalid or expired.'
    toast.error('Verification failed', {
      description: error.value,
    })
  } finally {
    isVerifying.value = false
  }
})
</script>

<template>
  <div class="-mt-[52px] flex min-h-screen items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="mb-8 text-center">
        <div class="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-md bg-foreground text-background">
          <Icon class="h-6 w-6" name="lucide:mail" />
        </div>
        <h1 class="mb-2 text-3xl font-semibold tracking-tight">
          Verify your email
        </h1>
        <p class="text-muted-foreground">
          <span v-if="isVerifying">Verifying your email address...</span>
          <span v-else-if="isSuccess">Your email has been verified successfully!</span>
          <span v-else>There was an issue verifying your email.</span>
        </p>
      </div>

      <Card>
        <CardContent class="pt-6">
          <div class="flex flex-col items-center justify-center gap-4 py-8">
            <div v-if="isVerifying" class="flex flex-col items-center gap-4">
              <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
              <p class="text-sm text-muted-foreground">
                Please wait while we verify your email...
              </p>
            </div>

            <div v-else-if="isSuccess" class="flex flex-col items-center gap-4">
              <div class="flex h-12 w-12 items-center justify-center rounded-full bg-green-100 dark:bg-green-900/20">
                <Icon class="h-6 w-6 text-green-600 dark:text-green-400" name="lucide:check" />
              </div>
              <p class="text-center text-sm text-muted-foreground">
                Your email has been verified successfully. Redirecting to sign in...
              </p>
            </div>

            <div v-else class="flex flex-col items-center gap-4">
              <div class="flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/20">
                <Icon class="h-6 w-6 text-red-600 dark:text-red-400" name="lucide:x" />
              </div>
              <p class="text-center text-sm text-muted-foreground">
                {{ error }}
              </p>
              <p class="text-center text-xs text-muted-foreground">
                The verification link may have expired. Please request a new verification email.
              </p>
            </div>
          </div>
        </CardContent>
        <CardFooter class="flex flex-col gap-4">
          <Button
            as-child
            class="w-full"
            variant="outline"
          >
            <NuxtLink to="/auth/login">
              <Icon class="mr-2 h-4 w-4" name="lucide:arrow-left" />
              Back to sign in
            </NuxtLink>
          </Button>
          <Button
            v-if="!isSuccess && !isVerifying"
            as-child
            class="w-full"
            variant="ghost"
          >
            <NuxtLink to="/auth/register">
              Create a new account
            </NuxtLink>
          </Button>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
