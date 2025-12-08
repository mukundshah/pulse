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

const { $api, $pulseAPI } = useNuxtApp()

const schema = z.object({
  name: z.string().min(1),
})

const { handleSubmit, isSubmitting, resetForm } = useForm({
  validationSchema: toTypedSchema(schema),
})

const onSubmit = handleSubmit(async (data) => {
  console.log(data)
  try {
    const res = await $pulseAPI('/v1/projects', {
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
  <form class="flex flex-row gap-y-6" @submit="onSubmit" @reset="onReset">
    <FormField v-slot="{ componentField }" name="name">
      <FormItem class="px-1.25 py-2 pb-2.5 w-full">
        <FormLabel class="sr-only">Project Name</FormLabel>
        <FormControl>
          <InputGroup>
            <InputGroupInput placeholder="Project XYZ" type="name" v-bind="componentField"/>
            <InputGroupAddon align="inline-end">
              <InputGroupButton type="submit" aria-label="Add project" title="Add project" size="icon-xs" class="rounded-full">
                <template v-if="!isSubmitting">
                  <Icon name="lucide:check" /> <span class="sr-only">Add project</span>
                </template>
                <template v-else>
                  <Icon name="lucide:loader-circle" class="animate-spin" /> <span class="sr-only">Adding project...</span>
                </template>
              </InputGroupButton>
              <InputGroupButton v-if="!isSubmitting" type="reset" aria-label="Cancel" title="Cancel" size="icon-xs" class="rounded-full">
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
