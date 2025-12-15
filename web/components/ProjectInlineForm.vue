<script setup lang="ts">
import type { PulseAPIResponse } from '#open-fetch'

import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { z } from 'zod'

const emit = defineEmits<{
  success: [project: PulseAPIResponse<'createProject'>]
  error: [error: Error]
  cancel: []
}>()

const { $pulseAPI } = useNuxtApp()

const schema = z.object({
  name: z.string().min(1),
})

const { handleSubmit, isSubmitting, resetForm } = useForm({
  validationSchema: toTypedSchema(schema),
})

const onSubmit = handleSubmit(async (data) => {
  try {
    const res = await $pulseAPI('/internal/projects', {
      method: 'POST',
      body: data,
    })
    emit('success', res)
  } catch (error: unknown) {
    emit('error', error as Error)
  }
})

const onReset = () => {
  resetForm()
  emit('cancel')
}
</script>

<template>
  <form class="flex flex-row gap-y-6" @reset="onReset" @submit="onSubmit">
    <FormField v-slot="{ componentField }" name="name">
      <FormItem class="px-1.25 py-2 pb-2.5 w-full">
        <FormLabel class="sr-only">
          Project Name
        </FormLabel>
        <FormControl>
          <InputGroup>
            <InputGroupInput placeholder="Project XYZ" type="name" v-bind="componentField" />
            <InputGroupAddon align="inline-end">
              <InputGroupButton
                aria-label="Add project"
                class="rounded-full"
                size="icon-xs"
                title="Add project"
                type="submit"
              >
                <template v-if="!isSubmitting">
                  <Icon name="lucide:check" /> <span class="sr-only">Add project</span>
                </template>
                <template v-else>
                  <Icon class="animate-spin" name="lucide:loader-circle" /> <span class="sr-only">Adding project...</span>
                </template>
              </InputGroupButton>
              <InputGroupButton
                v-if="!isSubmitting"
                aria-label="Cancel"
                class="rounded-full"
                size="icon-xs"
                title="Cancel"
                type="reset"
              >
                <Icon name="lucide:x" /> <span class="sr-only">Cancel</span>
              </InputGroupButton>
            </InputGroupAddon>
          </InputGroup>
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>
  </form>
</template>
