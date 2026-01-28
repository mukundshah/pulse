<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { Separator } from '@/components/ui/separator'
import { cn } from '@/utils/style'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()
</script>

<template>
  <Separator
    data-sidebar="separator"
    data-slot="sidebar-separator"
    :class="cn('mx-2 w-auto bg-sidebar-border', props.class)"
  >
    <slot></slot>
  </Separator>
</template>
