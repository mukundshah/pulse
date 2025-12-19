<script setup lang="ts">
import type { PulseAPIResponse } from '#open-fetch'
import type { ChartConfig } from '@/components/ui/chart'
import { VisStackedBar, VisXYContainer } from '@unovis/vue'
import { ChartContainer, ChartCrosshair, ChartTooltip, ChartTooltipContent, componentToString } from '@/components/ui/chart'

type Run = PulseAPIResponse<'listProjectChecks'>[number]['last_24_runs'][number]
type FormattedRun = Run & { index: number }

const props = defineProps<{
  runs: Run[]
}>()

const chartConfig = {
  status: {
    label: 'Status',
  },
  total_time_ms: {
    label: 'Response Time',
  },
} satisfies ChartConfig

const STATUS_MAP = {
  passing: {
    color: 'oklch(79.2% 0.209 151.711)',
    label: 'Passed',
  },
  degraded: {
    color: 'oklch(82.8% 0.189 84.429)',
    label: 'Degraded',
  },
  failing: {
    color: 'oklch(63.7% 0.237 25.331)',
    label: 'Failed',
  },
  unknown: {
    color: 'oklch(55.1% 0.027 264.364)',
    label: 'Unknown',
  },
}

const formattedRuns = computed(() => {
  return props.runs.map((d, i) => ({
    ...d,
    index: i,
    timestamp: new Date(d.timestamp as string).toLocaleString('en-US', {
      day: 'numeric',
      month: 'short',
      hour: 'numeric',
      minute: 'numeric',
      second: 'numeric',
    }),
  }))
})
</script>

<template>
  <ChartContainer class="h-10 w-full max-w-56" :config="chartConfig" :cursor="true">
    <VisXYContainer :data="formattedRuns" :padding="{ top: 4, right: 6, bottom: 0, left: 6 }">
      <VisStackedBar
        bar-min-height-1-px
        :bar-padding="0.4"
        :color="(d: FormattedRun) => STATUS_MAP[d.status as keyof typeof STATUS_MAP].color"
        :x="(d: FormattedRun) => d.index"
        :y="[(d: FormattedRun) => d.total_time_ms]"
      />

      <ChartTooltip />
      <ChartCrosshair
        :color="(d: FormattedRun) => STATUS_MAP[d.status as keyof typeof STATUS_MAP].color"
        :template="componentToString(chartConfig, ChartTooltipContent, {
          class: 'min-w-44',
          hideIndicator: true,
          labelKey: 'timestamp',
          valueFormatter: (d: unknown) => {
            if (typeof d === 'number') {
              return formatDuration(d * 1000)
            }
            return STATUS_MAP[d as keyof typeof STATUS_MAP].label
          },
        })"
      />
    </VisXYContainer>
  </ChartContainer>
</template>
