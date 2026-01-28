<script setup lang="ts">
import type { TagsInputItemDeleteProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'

import { reactiveOmit } from '@vueuse/core'
import { X } from 'lucide-vue-next'
import { TagsInputItemDelete, useForwardProps } from 'reka-ui'

import { cn } from '@/utils/style'

const props = defineProps<TagsInputItemDeleteProps & { class?: HTMLAttributes['class'] }>()

const delegatedProps = reactiveOmit(props, 'class')

const forwardedProps = useForwardProps(delegatedProps)
</script>

<template>
  <TagsInputItemDelete v-bind="forwardedProps" :class="cn('mr-1 flex rounded-sm bg-transparent', props.class)">
    <slot>
      <X class="size-4" />
    </slot>
  </TagsInputItemDelete>
</template>
