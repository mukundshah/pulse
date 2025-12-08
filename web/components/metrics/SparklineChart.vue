<script setup lang="ts">
interface Props {
  data: Array<{ timestamp: string | number | Date; value: number }>
  width?: number
  height?: number
  color?: string
  showPoints?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  width: 100,
  height: 20,
  color: 'hsl(var(--primary))',
  showPoints: false,
})

const chartConfig = {
  value: {
    label: 'Value',
    color: props.color,
  },
}

const chartData = computed(() => {
  return props.data.map((item, index) => ({
    timestamp: typeof item.timestamp === 'string' ? new Date(item.timestamp).getTime() : typeof item.timestamp === 'number' ? item.timestamp : item.timestamp.getTime(),
    value: item.value,
    index,
  }))
})

const maxValue = computed(() => Math.max(...chartData.value.map(d => d.value), 1))
const minValue = computed(() => Math.min(...chartData.value.map(d => d.value), 0))

const normalizedData = computed(() => {
  const range = maxValue.value - minValue.value || 1
  return chartData.value.map(d => ({
    ...d,
    normalized: ((d.value - minValue.value) / range) * props.height,
  }))
})
</script>

<template>
  <div
    class="relative inline-block"
    :style="{ width: `${width}px`, height: `${height}px` }"
  >
    <svg
      :width="width"
      :height="height"
      class="overflow-visible"
    >
      <!-- Background area -->
      <defs>
        <linearGradient id="sparkline-gradient" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" :stop-color="color" stop-opacity="0.2" />
          <stop offset="100%" :stop-color="color" stop-opacity="0" />
        </linearGradient>
      </defs>

      <!-- Area path -->
      <path
        v-if="normalizedData.length > 1"
        :d="`M ${normalizedData.map((d, i) => `${(i / (normalizedData.length - 1)) * width},${height - d.normalized}`).join(' L ')} L ${width},${height} L 0,${height} Z`"
        fill="url(#sparkline-gradient)"
        class="transition-all"
      />

      <!-- Line path -->
      <path
        v-if="normalizedData.length > 1"
        :d="`M ${normalizedData.map((d, i) => `${(i / (normalizedData.length - 1)) * width},${height - d.normalized}`).join(' L ')}`"
        :stroke="color"
        stroke-width="1.5"
        fill="none"
        class="transition-all"
      />

      <!-- Points -->
      <circle
        v-for="(point, i) in normalizedData"
        :key="i"
        :cx="(i / (normalizedData.length - 1)) * width"
        :cy="height - point.normalized"
        :r="showPoints ? 2 : 0"
        :fill="color"
        class="transition-all"
      />
    </svg>
  </div>
</template>


