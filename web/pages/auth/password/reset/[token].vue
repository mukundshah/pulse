<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import { z } from 'zod'

useHead({
  title: 'Reset password',
})

const route = useRoute()

const { resetPassword } = usePasswordManagement()

const resetPasswordSchema = z.object({
  password: z.string().regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/, 'Password must be at least 8 characters and contain at least one uppercase letter, one lowercase letter, one number, and one special character'),
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
        <div class="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-md bg-foreground text-background">
          <Icon class="h-6 w-6" name="lucide:lock" />
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
            <FormField v-slot="{ componentField }" name="password">
              <FormItem>
                <FormLabel>Password</FormLabel>
                <FormControl>
                  <Input
                    autocomplete="new-password"
                    placeholder="••••••••"
                    type="password"
                    v-bind="componentField"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="confirmPassword">
              <FormItem>
                <FormLabel>Confirm Password</FormLabel>
                <FormControl>
                  <Input
                    autocomplete="new-password"
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
              <span v-if="isSubmitting">Resetting password...</span>
              <span v-else>Reset password</span>
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
