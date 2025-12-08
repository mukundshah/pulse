<script setup lang="ts">
interface Props {
  status: 'passing' | 'degraded' | 'failing' | 'unknown'
  label?: string
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  label: undefined,
  size: 'md',
})

const statusConfig = {
  passing: {
    label: 'Passing',
    color: 'bg-green-500',
    textColor: 'text-green-600',
    dotColor: 'bg-green-500',
  },
  degraded: {
    label: 'Degraded',
    color: 'bg-yellow-500',
    textColor: 'text-yellow-600',
    dotColor: 'bg-yellow-500',
  },
  failing: {
    label: 'Failing',
    color: 'bg-red-500',
    textColor: 'text-red-600',
    dotColor: 'bg-red-500',
  },
  unknown: {
    label: 'Unknown',
    color: 'bg-gray-500',
    textColor: 'text-gray-600',
    dotColor: 'bg-gray-500',
  },
}

const config = statusConfig[props.status]
const displayLabel = props.label || config.label

const sizeClasses = {
  sm: 'h-1.5 w-1.5',
  md: 'h-2 w-2',
  lg: 'h-2.5 w-2.5',
}
</script>

<template>
  <div class="flex items-center gap-2">
    <div
      class="rounded-full"
      :class="[config.dotColor, sizeClasses[size]]"
    />
    <span
      v-if="displayLabel"
      class="text-sm font-medium"
      :class="config.textColor"
    >
      {{ displayLabel }}
    </span>
  </div>
</template>

