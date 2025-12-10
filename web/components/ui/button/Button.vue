<script setup lang="ts">
import type { PrimitiveProps } from 'reka-ui'
import type { ButtonHTMLAttributes, HTMLAttributes } from 'vue'
import type { ButtonVariants } from '.'

import { Primitive } from 'reka-ui'

import { Spinner } from '@/components/ui/spinner'
import { cn } from '@/utils/style'
import { buttonVariants } from '.'

interface Props extends PrimitiveProps {
  variant?: ButtonVariants['variant']
  size?: ButtonVariants['size']
  class?: HTMLAttributes['class']
  onClick?: (event: MouseEvent) => void
  loading?: boolean
  loadingAuto?: boolean
  type?: ButtonHTMLAttributes['type']
  disabled?: boolean
}

const props = withDefaults(
  defineProps<Props>(),
  {
    as: 'button',
    loading: false,
    type: 'button',
  },
)

const loadingAutoState = ref(false)

const isLoading = computed(() => {
  return props.loading || (props.loadingAuto && loadingAutoState.value)
})

const onClickWrapper = async (event: MouseEvent) => {
  loadingAutoState.value = true
  const callbacks = Array.isArray(props.onClick) ? props.onClick : [props.onClick]
  try {
    await Promise.all(callbacks.map(fn => fn?.(event)))
  } finally {
    loadingAutoState.value = false
  }
}
</script>

<template>
  <Primitive
    data-slot="button"
    :as="as"
    :as-child="asChild"
    :class="cn(buttonVariants({ variant, size }), props.class)"
    :disabled="props.disabled || isLoading"
    @click="onClickWrapper"
  >
    <slot name="loading">
      <Spinner v-if="isLoading" />
    </slot>
    <slot>
    </slot>
  </Primitive>
</template>
