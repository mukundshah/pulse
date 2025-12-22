<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm, Field as VeeField } from 'vee-validate'
import { toast } from 'vue-sonner'
import { z } from 'zod'

useHead({
  title: 'Reset password',
})

const route = useRoute()

const { resetPassword } = usePasswordManagement()

const resetPasswordSchema = z.object({
  password: z.string().regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9])\S{8,}$/, 'Password must be at least 8 characters, contain uppercase, lowercase, number, and a special character (no spaces).'),
  confirmPassword: z.string(),
}).refine(data => data.password === data.confirmPassword, {
  message: 'Passwords don\'t match',
  path: ['confirmPassword'],
})

const { handleSubmit, isSubmitting } = useForm({
  validationSchema: toTypedSchema(resetPasswordSchema),
})

const onSubmit = handleSubmit(async (data) => {
  await resetPassword({ ...data, token: route.params.token?.toString() ?? '' })
  toast('Password reset successfully')
  await navigateTo('/auth/login')
})
</script>

<template>
  <div class="-mt-[52px] flex min-h-screen items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="mb-8 text-center">
        <div class="mb-4 inline-flex items-center justify-center">
          <AppIcon class="text-foreground size-12" />
        </div>
        <h1 class="mb-2 text-3xl font-semibold tracking-tight">
          Reset your password
        </h1>
        <p class="text-muted-foreground">
          Enter your new password below
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>New password</CardTitle>
          <CardDescription>
            Choose a strong password for your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form class="flex flex-col gap-y-6" @submit="onSubmit">
            <VeeField v-slot="{ componentField, errors }" name="password">
              <Field :data-invalid="!!errors.length">
                <FieldLabel for="password">
                  Password
                </FieldLabel>
                <PasswordInput
                  id="password"
                  v-bind="componentField"
                  autocomplete="new-password"
                  placeholder="••••••••"
                  :aria-invalid="!!errors.length"
                />
                <FieldError v-if="errors.length" :errors="errors" />
              </Field>
            </VeeField>

            <VeeField v-slot="{ componentField, errors }" name="confirmPassword">
              <Field :data-invalid="!!errors.length">
                <FieldLabel for="confirmPassword">
                  Confirm Password
                </FieldLabel>
                <PasswordInput
                  id="confirmPassword"
                  v-bind="componentField"
                  autocomplete="new-password"
                  placeholder="••••••••"
                  :aria-invalid="!!errors.length"
                />
                <FieldError v-if="errors.length" :errors="errors" />
              </Field>
            </VeeField>

            <Button
              class="w-full"
              type="submit"
              :loading="isSubmitting"
            >
              Reset password
            </Button>
          </form>
        </CardContent>
        <CardFooter>
          <Button
            as-child
            class="w-full"
            variant="ghost"
          >
            <NuxtLink to="/auth/login">
              <Icon class="mr-2 h-4 w-4" name="lucide:arrow-left" />
              Back to sign in
            </NuxtLink>
          </Button>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
