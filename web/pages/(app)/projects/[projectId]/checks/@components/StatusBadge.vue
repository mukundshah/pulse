<script setup lang="ts">
interface Props {
  status: 'passing' | 'degraded' | 'failing' | 'unknown'
}

const props = defineProps<Props>()

const statusConfig = {
  passing: {
    label: 'Passing',
    icon: 'lucide:check',
    style: {
      '--color-status-secondary': 'oklch(79.2% 0.209 151.711)',
      '--color-status-primary': 'oklch(72.3% 0.219 149.579)',
    },
  },
  degraded: {
    label: 'Degraded',
    icon: 'lucide:octagon-alert',
    style: {
      '--color-status-secondary': 'oklch(82.8% 0.189 84.429)',
      '--color-status-primary': 'oklch(76.9% 0.188 70.08)',
    },
  },
  failing: {
    label: 'Failing',
    icon: 'lucide:x',
    style: {
      '--color-status-secondary': 'oklch(63.7% 0.237 25.331)',
      '--color-status-primary': 'oklch(63.7% 0.237 25.331)',
    },
  },
  unknown: {
    label: 'Unknown',
    icon: 'lucide:minus',
    style: {
      '--color-status-secondary': 'oklch(70.7% 0.022 261.325)',
      '--color-status-primary': 'oklch(55.1% 0.027 264.364)',
    },
  },
}

const config = statusConfig[props.status || 'unknown']
</script>

<template>
  <div :style="config.style">
    <span class="relative flex size-2.5">
      <span v-if="props.status !== 'unknown'" class="absolute inline-flex h-full w-full animate-ping rounded-full bg-(--color-status-secondary) opacity-75"></span>
      <span class="relative inline-flex size-2.5 rounded-full bg-(--color-status-primary)"></span>
    </span>
    <span class="sr-only">{{ config.label }}</span>
  </div>
</template>
