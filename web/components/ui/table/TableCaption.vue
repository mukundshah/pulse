<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/utils/style'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <caption
    data-slot="table-caption"
    :class="cn('mt-4 text-sm text-muted-foreground', props.class)"
  >
    <slot></slot>
  </caption>
</template>
