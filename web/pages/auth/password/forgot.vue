<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import { z } from 'zod'

const { $pulseAPI } = useNuxtApp()

useHead({
  title: 'Forgot password',
})

const forgotPasswordSchema = z.object({
  email: z.email(),
})

const { handleSubmit, isSubmitting } = useForm({
  validationSchema: toTypedSchema(forgotPasswordSchema),
})

const onSubmit = handleSubmit(async (data) => {
  await $pulseAPI('/v1/auth/password/reset', {
    method: 'POST',
    body: data,
  })
  toast('Password reset link sent to your email')
})
</script>

<template>
  <div class="-mt-[52px] flex min-h-screen items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="mb-8 text-center">
        <div class="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-md bg-foreground text-background">
          <Icon class="h-6 w-6" name="lucide:key" />
        </div>
        <h1 class="mb-2 text-3xl font-semibold tracking-tight">
          Forgot password?
        </h1>
        <p class="text-muted-foreground">
          No worries, we'll send you reset instructions
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Reset password</CardTitle>
          <CardDescription>
            Enter your email address and we'll send you a link to reset your password
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

            <Button
              class="w-full"
              type="submit"
              :disabled="isSubmitting"
            >
              <span v-if="isSubmitting">Sending...</span>
              <span v-else>Send reset link</span>
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
