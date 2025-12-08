<script setup lang="ts">
import { ChartContainer, ChartTooltip, ChartTooltipContent } from '@/components/ui/chart'

interface Props {
  data: Array<{ timestamp: string | number | Date; [key: string]: string | number | Date }>
  metrics: Array<{ key: string; label: string; color: string }>
  type?: 'line' | 'bar'
  height?: number
}

const props = withDefaults(defineProps<Props>(), {
  type: 'line',
  height: 300,
})

const chartConfig = computed(() => {
  const config = {}
  props.metrics.forEach((metric) => {
    config[metric.key] = {
      label: metric.label,
      color: metric.color,
    }
  })
  return config
})

const chartData = computed(() => {
  return props.data.map((item) => {
    const result = {
      timestamp: typeof item.timestamp === 'string' ? new Date(item.timestamp).getTime() : typeof item.timestamp === 'number' ? item.timestamp : item.timestamp.getTime(),
    }
    props.metrics.forEach((metric) => {
      result[metric.key] = item[metric.key] || 0
    })
    return result
  })
})
</script>

<template>
  <Card>
    <CardContent class="p-6">
      <ChartContainer
        :config="chartConfig"
        class="h-[300px] w-full"
      >
        <template #default="{ config }">
          <!-- Simple line chart using SVG -->
          <svg
            :width="'100%'"
            :height="height"
            class="overflow-visible"
          >
            <defs>
              <linearGradient
                v-for="metric in metrics"
                :key="metric.key"
                :id="`gradient-${metric.key}`"
                x1="0%"
                y1="0%"
                x2="0%"
                y2="100%"
              >
                <stop offset="0%" :stop-color="metric.color" stop-opacity="0.2" />
                <stop offset="100%" :stop-color="metric.color" stop-opacity="0" />
              </linearGradient>
            </defs>

            <!-- Grid lines -->
            <g v-for="i in 5" :key="`grid-${i}`" class="text-muted-foreground">
              <line
                :x1="0"
                :y1="(i * height) / 5"
                :x2="'100%'"
                :y2="(i * height) / 5"
                stroke="currentColor"
                stroke-width="0.5"
                stroke-dasharray="2,2"
                opacity="0.3"
              />
            </g>

            <!-- Data paths for each metric -->
            <g
              v-for="metric in metrics"
              :key="metric.key"
            >
              <path
                v-if="chartData.length > 1"
                :d="`M ${chartData.map((d, i) => `${(i / (chartData.length - 1)) * 100}%,${height - ((d[metric.key] as number) / Math.max(...chartData.map(dd => dd[metric.key] as number), 1)) * height * 0.9}`).join(' L ')}`"
                :stroke="metric.color"
                stroke-width="2"
                fill="none"
                class="transition-all"
              />
            </g>
          </svg>
        </template>
      </ChartContainer>
    </CardContent>
  </Card>
</template>


