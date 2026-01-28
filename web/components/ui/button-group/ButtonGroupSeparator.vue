<script setup lang="ts">
import type { SeparatorProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'

import { reactiveOmit } from '@vueuse/core'

import { Separator } from '@/components/ui/separator'
import { cn } from '@/utils/style'

const props = withDefaults(
  defineProps<SeparatorProps & { class?: HTMLAttributes['class'] }>(),
  {
    orientation: 'vertical',
  },
)
const delegatedProps = reactiveOmit(props, 'class')
</script>

<template>
  <Separator
    data-slot="button-group-separator"
    v-bind="delegatedProps"
    :class="cn(
      'relative m-0! self-stretch bg-input data-[orientation=vertical]:h-auto',
      props.class,
    )"
    :orientation="props.orientation"
  />
</template>
