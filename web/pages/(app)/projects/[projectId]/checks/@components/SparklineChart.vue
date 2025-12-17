<script setup lang="ts">
import type { PulseAPIResponse } from '#open-fetch'
import type { ChartConfig } from '@/components/ui/chart'
import { VisStackedBar, VisXYContainer } from '@unovis/vue'
import { ChartContainer, ChartCrosshair, ChartTooltip, ChartTooltipContent, componentToString } from '@/components/ui/chart'

type Run = PulseAPIResponse<'listProjectChecks'>[number]['last_24_runs'][number]

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
}
</script>

<template>
  <ChartContainer class="h-12 w-full" :config="chartConfig" :cursor="true">
    <VisXYContainer :data="props.runs" :padding="{ top: 10, right: 0, bottom: 0, left: 0 }">
      <VisStackedBar
        :bar-max-width="6"
        :bar-padding="0.4"
        :color="(d: Run) => STATUS_MAP[d.status as keyof typeof STATUS_MAP].color"
        :x="(d: Run) => new Date(d.timestamp as string)"
        :y="[(d: Run) => d.total_time_ms]"
      />

      <ChartTooltip />
      <ChartCrosshair
        :color="(d: Run) => STATUS_MAP[d.status as keyof typeof STATUS_MAP].color"
        :template="componentToString(chartConfig, ChartTooltipContent, {
          class: 'min-w-40',
          hideIndicator: true,
          labelFormatter: (d) => {
            return new Date(d).toLocaleString('en-US', {
              day: 'numeric',
              month: 'short',
              hour: 'numeric',
              minute: 'numeric',
            })
          },
          valueFormatter: (d: unknown) => {
            if (typeof d === 'number') {
              return d.toLocaleString('en-US', {
                style: 'unit',
                unit: 'millisecond',
              })
            }
            return STATUS_MAP[d as keyof typeof STATUS_MAP].label
          },
        })"
      />
    </VisXYContainer>
  </ChartContainer>
</template>
