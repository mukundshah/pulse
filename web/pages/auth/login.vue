<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import { z } from 'zod'

const { $pulseAPI } = useNuxtApp()

useHead({
  title: 'Sign in',
})

const loginSchema = z.object({
  email: z.email('Please enter a valid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
})

const { handleSubmit, isSubmitting } = useForm({
  validationSchema: toTypedSchema(loginSchema),
})

const onSubmit = handleSubmit(async (data) => {
  await $pulseAPI('/v1/auth/login', {
    method: 'POST',
    body: data,
  })
  toast('Welcome back!', {
    description: 'You are now logged in',
  })
  await navigateTo('/dashboard')
})
</script>

<template>
  <div class="-mt-[52px] flex min-h-screen items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="mb-8 text-center">
        <div class="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-md bg-foreground text-background">
          <span class="text-xl font-semibold">P</span>
        </div>
        <h1 class="mb-2 text-3xl font-semibold tracking-tight">
          Welcome back
        </h1>
        <p class="text-muted-foreground">
          Sign in to your account to continue
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Sign in</CardTitle>
          <CardDescription>
            Enter your credentials to access your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form class="flex flex-col gap-y-6" @submit="onSubmit">
            <FormField v-slot="{ componentField }" name="email">
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input
                    placeholder="you@example.com"
                    type="email"
                    v-bind="componentField"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="password">
              <FormItem>
                <div class="flex items-center justify-between">
                  <FormLabel>Password</FormLabel>
                  <NuxtLink
                    class="text-sm text-muted-foreground hover:text-foreground"
                    to="/auth/password/forgot"
                  >
                    Forgot password?
                  </NuxtLink>
                </div>
                <FormControl>
                  <Input
                    placeholder="••••••••"
                    type="password"
                    v-bind="componentField"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <Button
              class="w-full"
              type="submit"
              :disabled="isSubmitting"
            >
              <span v-if="isSubmitting">Signing in...</span>
              <span v-else>Sign in</span>
            </Button>
          </form>
        </CardContent>
        <CardFooter class="flex flex-col gap-4">
          <div class="relative w-full text-center text-sm">
            <span class="bg-card px-2 text-muted-foreground">
              Don't have an account?
            </span>
          </div>
          <Button
            as-child
            class="w-full"
            variant="outline"
          >
            <NuxtLink to="/auth/register">
              Create an account
            </NuxtLink>
          </Button>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
