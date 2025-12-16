<script setup lang="ts">
import type { PulseAPIResponse } from '#open-fetch'
import type { ChartConfig } from '@/components/ui/chart'
import { VisAxis, VisStackedBar, VisXYContainer } from '@unovis/vue'
import { ChartContainer, ChartCrosshair, ChartTooltip, ChartTooltipContent, componentToString } from '@/components/ui/chart'

const props = withDefaults(defineProps<{
  projectId: string
  checkId: string
  period?: 'today' | '1hr' | '3hr' | '24hr' | '7d' | '30d'
}>(), {
  period: '7d',
})

type Data = PulseAPIResponse<'getCheckUptime'>['data'][number]

const chartConfig = {
  passing: {
    label: 'Passing',
    color: 'oklch(79.2% 0.209 151.711)',
  },
  degraded: {
    label: 'Degraded',
    color: 'oklch(82.8% 0.189 84.429)',
  },
  failing: {
    label: 'Failing',
    color: 'oklch(63.7% 0.237 25.331)',
  },
} satisfies ChartConfig

const TIME_FORMAT = {
  minute: {
    hour: 'numeric',
    minute: 'numeric',
  },
  hour: {
    day: 'numeric',
    month: 'short',
    hour: 'numeric',
  },
  day: {
    day: 'numeric',
    month: 'short',
  },
} satisfies Record<string, Intl.DateTimeFormatOptions>

// Fetch uptime data
const { data: response, pending, error, refresh } = useLazyPulseAPI('/internal/projects/{projectId}/checks/{checkId}/uptime', {
  path: {
    projectId: props.projectId,
    checkId: props.checkId,
  },
  query: {
    period: props.period,
  },
})
</script>

<template>
  <div class="h-48 w-full">
    <Skeleton v-if="pending" class="h-full w-full" />

    <div v-else-if="error" class="h-full flex flex-col items-center justify-center text-sm text-destructive gap-2">
      <p>Failed to load uptime data</p>
      <Button size="sm" variant="outline" @click="refresh()">
        Retry
      </Button>
    </div>

    <div v-else class="h-full w-full">
      <ChartContainer class="h-48 w-full" :config="chartConfig" :cursor="true">
        <VisXYContainer :data="response?.data">
          <VisStackedBar
            :bar-padding="0.5"
            :color="[chartConfig.passing.color, chartConfig.degraded.color, chartConfig.failing.color]"
            :x="(d: Data) => new Date(d.timestamp)"
            :y="[(d: Data) => d.passing, (d: Data) => d.degraded, (d: Data) => d.failing]"
          />

          <VisAxis
            type="x"
            :domain-line="false"
            :grid-line="false"
            :tick-format="(d: number) => {
              const date = new Date(d)
              return date.toLocaleString('en-US', TIME_FORMAT[response?.time_bucket as keyof typeof TIME_FORMAT])
            }"
            :tick-line="false"
            :x="(d: Data) => new Date(d.timestamp)"
          />
          <VisAxis
            type="y"
            :domain-line="false"
            :grid-line="true"
            :tick-format="(d: number) => d.toLocaleString()"
            :tick-line="false"
            :y="(d: Data) => d.total_runs"
          />

          <ChartTooltip />
          <ChartCrosshair
            :color="[chartConfig.passing.color, chartConfig.degraded.color, chartConfig.failing.color]"
            :template=" componentToString(chartConfig, ChartTooltipContent, {
              labelFormatter(d) {
                return new Date(d).toLocaleString('en-US', TIME_FORMAT[response?.time_bucket as keyof typeof TIME_FORMAT])
              },
            })
            "
          />
        </VisXYContainer>
      </ChartContainer>
    </div>
  </div>
</template>
