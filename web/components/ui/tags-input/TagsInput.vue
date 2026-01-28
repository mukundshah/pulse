<script setup lang="ts">
import type { TagsInputRootEmits, TagsInputRootProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'

import { reactiveOmit } from '@vueuse/core'
import { TagsInputRoot, useForwardPropsEmits } from 'reka-ui'

import { cn } from '@/utils/style'

const props = defineProps<TagsInputRootProps & { class?: HTMLAttributes['class'] }>()
const emits = defineEmits<TagsInputRootEmits>()

const delegatedProps = reactiveOmit(props, 'class')

const forwarded = useForwardPropsEmits(delegatedProps, emits)
</script>

<template>
  <TagsInputRoot v-slot="slotProps" v-bind="forwarded" :class="cn('flex flex-wrap items-center gap-2 rounded-md border border-input bg-background px-3 py-1.5 text-sm', props.class)">
    <slot v-bind="slotProps"></slot>
  </TagsInputRoot>
</template>
