<script setup lang="ts">
import type { FetchError } from 'ofetch'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm, Field as VeeField } from 'vee-validate'
import { toast } from 'vue-sonner'
import { z } from 'zod'

const route = useRoute()

const { login } = useAuth()
const { requestEmailVerification } = useEmailManagement()

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
  try {
    await login(data)
    toast('Welcome back!', {
      description: 'You are now logged in',
    })
    await navigateTo(route.query.next as string || '/dashboard')
  } catch (error: unknown) {
    // Handle 403 - email not verified
    if ((error as FetchError)?.status === 403) {
      toast.error('Email not verified', {
        description: 'Please verify your email address before signing in.',
        action: {
          label: 'Resend verification email',
          onClick: async () => {
            try {
              await requestEmailVerification({ email: data.email })
              toast.success('Verification email sent', {
                description: 'Please check your inbox for the verification link.',
              })
            } catch (resendError: unknown) {
              toast.error('Failed to send verification email', {
                description: (resendError as FetchError)?.message || 'Please try again later.',
              })
            }
          },
        },
      })
    } else {
      // Handle other errors
      toast.error('Sign in failed', {
        description: (error as FetchError)?.message || 'Invalid credentials. Please try again.',
      })
    }
  }
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
            <FieldGroup>
              <VeeField v-slot="{ componentField, errors }" name="email">
                <Field :data-invalid="!!errors.length">
                  <FieldLabel for="email">
                    Email
                  </FieldLabel>
                  <Input
                    id="email"
                    v-bind="componentField"
                    autocomplete="email"
                    placeholder="you@example.com"
                    :aria-invalid="!!errors.length"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
              <VeeField v-slot="{ componentField, errors }" name="password">
                <Field :data-invalid="!!errors.length">
                  <FieldLabel for="password">
                    Password
                  </FieldLabel>
                  <PasswordInput
                    id="password"
                    v-bind="componentField"
                    autocomplete="current-password"
                    placeholder="••••••••"
                    :aria-invalid="!!errors.length"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
            </FieldGroup>
            <Button
              class="w-full"
              type="submit"
              :loading="isSubmitting"
            >
              Sign in
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
