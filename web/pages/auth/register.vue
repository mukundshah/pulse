<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm, Field as VeeField } from 'vee-validate'
import { toast } from 'vue-sonner'
import { z } from 'zod'

const { signup } = useAuth()

useHead({
  title: 'Create account',
})

const registerSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.email('Please enter a valid email address'),
  password: z.string().regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/, 'Password must be at least 8 characters and contain at least one uppercase letter, one lowercase letter, one number, and one special character'),
  confirmPassword: z.string(),
}).refine(data => data.password === data.confirmPassword, {
  message: 'Passwords don\'t match',
  path: ['confirmPassword'],
})

const { handleSubmit, isSubmitting } = useForm({
  validationSchema: toTypedSchema(registerSchema),
})

const onSubmit = handleSubmit(async (data) => {
  await signup(data)

  toast('Account created successfully')
  await navigateTo('/auth/login')
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
          Create an account
        </h1>
        <p class="text-muted-foreground">
          Get started with Pulse today
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Sign up</CardTitle>
          <CardDescription>
            Enter your information to create your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form class="flex flex-col gap-y-6" @submit="onSubmit">
            <FieldGroup>
              <VeeField v-slot="{ field, errors }" name="name">
                <Field :data-invalid="!!errors.length">
                  <FieldLabel for="name">
                    Name
                  </FieldLabel>
                  <Input
                    id="name"
                    v-bind="field"
                    autocomplete="name"
                    placeholder="John Doe"
                    type="text"
                    :aria-invalid="!!errors.length"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
              <VeeField v-slot="{ field, errors }" name="email">
                <Field :data-invalid="!!errors.length">
                  <FieldLabel for="email">
                    Email
                  </FieldLabel>
                  <Input
                    id="email"
                    v-bind="field"
                    autocomplete="email"
                    placeholder="you@example.com"
                    type="email"
                    :aria-invalid="!!errors.length"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
              <VeeField v-slot="{ field, errors }" name="password">
                <Field :data-invalid="!!errors.length">
                  <FieldLabel for="password">
                    Password
                  </FieldLabel>
                  <PasswordInput
                    id="password"
                    v-bind="field"
                    autocomplete="new-password"
                    placeholder="••••••••"
                    :aria-invalid="!!errors.length"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
              <VeeField v-slot="{ field, errors }" name="confirmPassword">
                <Field :data-invalid="!!errors.length">
                  <FieldLabel for="confirmPassword">
                    Confirm Password
                  </FieldLabel>
                  <PasswordInput
                    id="confirmPassword"
                    v-bind="field"
                    autocomplete="new-password"
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
              Create account
            </Button>
          </form>
        </CardContent>
        <CardFooter class="flex flex-col gap-4">
          <div class="relative w-full text-center text-sm">
            <span class="bg-card px-2 text-muted-foreground">
              Already have an account?
            </span>
          </div>
          <Button
            as-child
            class="w-full"
            variant="outline"
          >
            <NuxtLink to="/auth/login">
              Sign in
            </NuxtLink>
          </Button>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
