<script setup lang="ts">
import type { PulseAPIResponse } from '#open-fetch'
import type { ChartConfig } from '@/components/ui/chart'
import { Scale } from '@unovis/ts'
import { VisAxis, VisLine, VisXYContainer } from '@unovis/vue'

import { ChartContainer, ChartCrosshair, ChartTooltip, ChartTooltipContent, componentToString } from '@/components/ui/chart'
import { TIME_FORMAT } from '@/constants/intl'
import { formatDuration } from '@/utils/formatters'

const props = withDefaults(defineProps<{
  projectId: string
  checkId: string
  period?: 'today' | '1hr' | '3hr' | '24hr' | '7d' | '30d'
}>(), {
  period: '7d',
})

type TimingData = PulseAPIResponse<'getCheckTimings'>['data'][number]

const chartConfig = {
  dns: {
    label: 'DNS',
    color: 'oklch(62.7% 0.265 303.9)',
  },
  tcp: {
    label: 'TCP',
    color: 'oklch(70.5% 0.213 47.604)',
  },
  tls: {
    label: 'TLS',
    color: 'oklch(79.5% 0.184 86.047)',
  },
  request: {
    label: 'Request',
    color: 'oklch(62.3% 0.214 259.815)',
  },
  ttfb: {
    label: 'TTFB',
    color: 'oklch(68.5% 0.169 237.323)',
  },
  download: {
    label: 'Download',
    color: 'oklch(72.3% 0.219 149.579)',
  },
} satisfies ChartConfig

const { data: response, pending, error, refresh } = useLazyPulseAPI('/internal/projects/{projectId}/checks/{checkId}/timings', {
  path: {
    projectId: props.projectId,
    checkId: props.checkId,
  },
  query: {
    period: props.period,
  },
})

interface ChartDataPoint {
  timestamp: Date
  dns: number | null
  tcp: number | null
  tls: number | null
  request: number | null
  ttfb: number | null
  download: number | null
}

const chartData = computed(() => {
  if (!response.value?.data) return []

  return response.value.data.map((point: TimingData) => {
    const timings = point.network_timings || {}

    return {
      timestamp: new Date(point.timestamp),
      dns: typeof timings.dns_duration_us === 'number' ? timings.dns_duration_us : null,
      tcp: typeof timings.tcp_duration_us === 'number' ? timings.tcp_duration_us : null,
      tls: typeof timings.tls_duration_us === 'number' ? timings.tls_duration_us : null,
      request: typeof timings.request_duration_us === 'number' ? timings.request_duration_us : null,
      ttfb: typeof timings.ttfb_us === 'number' ? timings.ttfb_us : null,
      download: typeof timings.download_us === 'number' ? timings.download_us : null,
    } as ChartDataPoint
  }).filter((d: ChartDataPoint) =>
    d.dns !== null
    || d.tcp !== null
    || d.tls !== null
    || d.request !== null
    || d.ttfb !== null
    || d.download !== null,
  )
})

// Determine time bucket for formatting based on period
const timeBucket = computed(() => {
  switch (props.period) {
    case 'today':
    case '1hr':
    case '3hr':
      return 'minute'
    case '24hr':
    case '7d':
      return 'hour'
    case '30d':
      return 'day'
    default:
      return 'hour'
  }
})
</script>

<template>
  <div class="h-74 w-full">
    <Skeleton v-if="pending" class="h-full w-full" />

    <div v-else-if="error" class="h-full flex flex-col items-center justify-center text-sm text-destructive gap-2">
      <p>Failed to load performance data</p>
      <Button size="sm" variant="outline" @click="refresh()">
        Retry
      </Button>
    </div>

    <div v-else-if="chartData.length === 0" class="h-full flex items-center justify-center text-sm text-muted-foreground">
      No performance data available
    </div>

    <ChartContainer
      v-else
      class="h-72 w-full"
      :config="chartConfig"
      :cursor="true"
    >
      <VisXYContainer :data="chartData" :margin="{ top: 0, right: 0, bottom: 10, left: 0 }" :y-scale="Scale.scalePow().exponent(0.1)">
        <VisLine
          v-if="chartData.some(d => d.dns !== null)"
          :color="chartConfig.dns.color"
          :x="(d: ChartDataPoint) => d.timestamp"
          :y="(d: ChartDataPoint) => d.dns ?? null"
        />
        <VisLine
          v-if="chartData.some(d => d.tcp !== null)"
          :color="chartConfig.tcp.color"
          :x="(d: ChartDataPoint) => d.timestamp"
          :y="(d: ChartDataPoint) => d.tcp ?? null"
        />
        <VisLine
          v-if="chartData.some(d => d.tls !== null)"
          :color="chartConfig.tls.color"
          :x="(d: ChartDataPoint) => d.timestamp"
          :y="(d: ChartDataPoint) => d.tls ?? null"
        />
        <VisLine
          v-if="chartData.some(d => d.request !== null)"
          :color="chartConfig.request.color"
          :x="(d: ChartDataPoint) => d.timestamp"
          :y="(d: ChartDataPoint) => d.request ?? null"
        />
        <VisLine
          v-if="chartData.some(d => d.ttfb !== null)"
          :color="chartConfig.ttfb.color"
          :x="(d: ChartDataPoint) => d.timestamp"
          :y="(d: ChartDataPoint) => d.ttfb ?? null"
        />
        <VisLine
          v-if="chartData.some(d => d.download !== null)"
          :color="chartConfig.download.color"
          :x="(d: ChartDataPoint) => d.timestamp"
          :y="(d: ChartDataPoint) => d.download ?? null"
        />

        <VisAxis
          type="x"
          :domain-line="false"
          :grid-line="false"
          :tick-format="(d: number) => {
            const date = new Date(d)
            return date.toLocaleString('en-US', TIME_FORMAT[timeBucket as keyof typeof TIME_FORMAT])
          }"
          :tick-line="true"
          :x="(d: ChartDataPoint) => d.timestamp"
        />
        <VisAxis
          type="y"
          :domain-line="false"
          :grid-line="true"
          :tick-format="(d: number) => formatDuration(d)"
          :tick-line="false"
          :tick-values="[10, 100, 1000, 10000, 100000, 200000, 300000, 500000, 1000000, 2500000, 5000000, 10000000]"
        />

        <ChartTooltip />
        <ChartCrosshair
          :color="[
            chartConfig.dns.color,
            chartConfig.tcp.color,
            chartConfig.tls.color,
            chartConfig.request.color,
            chartConfig.ttfb.color,
            chartConfig.download.color,
          ]"
          :template="componentToString(chartConfig, ChartTooltipContent, {
            class: 'min-w-40',
            labelFormatter: (d) => {
              return new Date(d).toLocaleString('en-US', TIME_FORMAT[timeBucket as keyof typeof TIME_FORMAT])
            },
            valueFormatter: (d: unknown) => {
              if (typeof d === 'number') {
                return formatDuration(d)
              }
              return String(d)
            },
          })"
        />
      </VisXYContainer>
      <ChartLegendContent vertical-align="bottom" />
    </ChartContainer>
  </div>
</template>
