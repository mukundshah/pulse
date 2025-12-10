<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { useVModel } from '@vueuse/core'
import { EyeIcon, EyeOffIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import { InputGroup, InputGroupAddon, InputGroupButton, InputGroupInput } from '@/components/ui/input-group'
import { cn } from '@/utils/style'

defineOptions({
  inheritAttrs: false,
})

const props = defineProps<{
  defaultValue?: string
  modelValue?: string
  class?: HTMLAttributes['class']
}>()

const emits = defineEmits<{
  'update:modelValue': [payload: string]
}>()

const modelValue = useVModel(props, 'modelValue', emits, {
  passive: true,
  defaultValue: props.defaultValue,
})

const isVisible = ref(false)

const toggleVisibility = () => {
  isVisible.value = !isVisible.value
}
</script>

<template>
  <InputGroup>
    <InputGroupInput
      v-model="modelValue"
      :class="cn('', props.class)"
      :type="isVisible ? 'text' : 'password'"
      v-bind="$attrs"
    />
    <InputGroupAddon align="inline-end" class="pr-2">
      <InputGroupButton
        aria-label="Toggle password visibility"
        class="rounded-sm rounded-l-none"
        size="icon-sm"
        title="Toggle password visibility"
        type="button"
        @click="toggleVisibility"
      >
        <EyeIcon v-if="isVisible" aria-hidden="true" :size="16" />
        <EyeOffIcon v-else aria-hidden="true" :size="16" />
      </InputGroupButton>
    </InputGroupAddon>
  </InputGroup>
</template>
