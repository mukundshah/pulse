<script lang="ts">
import { ASSERTION_PROPERTIES, IP_VERSION_LABELS } from '@/constants/http'
import { constructURL } from '@/utils/url'

const STATUS_ICON_COLOR_MAP = {
  passing: {
    icon: 'lucide:circle-check',
    color: 'text-green-500',
  },
  degraded: {
    icon: 'lucide:circle-alert',
    color: 'text-amber-500',
  },
  failing: {
    icon: 'lucide:circle-x',
    color: 'text-red-500',
  },
  unknown: {
    icon: 'lucide:circle-minus',
    color: 'text-gray-500',
  },
} as const

const STATUS_TEXT_MAP = {
  passing: 'Passed',
  failing: 'Failed',
  degraded: 'Degraded',
  unknown: 'Unknown',
} as const

const durationFormatter = (duration: number) => {
  if (duration < 1000) {
    return `${duration} µs`
  }
  if (duration < 1000000) {
    return `${(duration / 1000).toFixed(0)} ms`
  }

  return `${(duration / 1000000).toFixed(2)} s`
}

const parseTimestamp = (ts: string): number => {
  return new Date(ts).getTime() * 1000
}
</script>

<script setup lang="ts">
const route = useRoute()
const { projectId, checkId, runId } = route.params as { projectId: string, checkId: string, runId: string }

const { data: run } = await usePulseAPI('/internal/projects/{projectId}/checks/{checkId}/runs/{runId}', {
  path: {
    projectId,
    checkId,
    runId,
  },
})

useHead({
  title: `Run #${run.value?.id.slice(0, 8)}`,
})

const timelineData = computed(() => {
  if (!run.value?.network_timings || typeof run.value.network_timings !== 'object') {
    return []
  }

  const timings = run.value.network_timings as Record<string, unknown>

  const phases: Array<{
    type: string
    startTime: number // in microseconds
    length: number // in microseconds
    color: string
  }> = []

  // DNS phase
  if (timings.dns_start && typeof timings.dns_duration_us === 'number') {
    phases.push({
      type: 'DNS',
      startTime: parseTimestamp(timings.dns_start as string),
      length: timings.dns_duration_us,
      color: 'oklch(62.7% 0.265 303.9)',
    })
  }

  // TCP phase
  if (timings.tcp_start && typeof timings.tcp_duration_us === 'number') {
    phases.push({
      type: 'TCP',
      startTime: parseTimestamp(timings.tcp_start as string),
      length: timings.tcp_duration_us,
      color: 'oklch(70.5% 0.213 47.604)',
    })
  }

  // TLS phase
  if (timings.tls_start && typeof timings.tls_duration_us === 'number') {
    phases.push({
      type: 'TLS',
      startTime: parseTimestamp(timings.tls_start as string),
      length: timings.tls_duration_us,
      color: 'oklch(79.5% 0.184 86.047)',
    })
  }

  // Request phase
  if (typeof timings.request_duration_us === 'number') {
    const requestStart = timings.tls_done || timings.request_start
    if (requestStart) {
      phases.push({
        type: 'Request',
        startTime: parseTimestamp(requestStart as string),
        length: timings.request_duration_us,
        color: 'oklch(62.3% 0.214 259.815)',
      })
    }
  }

  // TTFB phase
  if (timings.request_sent && typeof timings.ttfb_us === 'number') {
    phases.push({
      type: 'TTFB',
      startTime: parseTimestamp(timings.request_sent as string),
      length: timings.ttfb_us,
      color: 'oklch(68.5% 0.169 237.323)',
    })
  }

  // Download phase
  if (timings.first_byte && typeof timings.download_us === 'number') {
    phases.push({
      type: 'Download',
      startTime: parseTimestamp(timings.first_byte as string),
      length: timings.download_us,
      color: 'oklch(72.3% 0.219 149.579)',
    })
  }

  if (phases.length === 0) {
    return []
  }

  // Find the earliest start time and the latest end time
  const earliestTime = Math.min(...phases.map(p => p.startTime))
  const latestEndTime = Math.max(...phases.map(p => p.startTime + p.length))
  const totalDuration = latestEndTime - earliestTime

  // Calculate left offset and width for each phase
  return phases.map(phase => ({
    type: phase.type,
    color: phase.color,
    length: phase.length,
    left: ((phase.startTime - earliestTime) / totalDuration) * 100,
    width: (phase.length / totalDuration) * 100,
  }))
})
</script>

<template>
  <div class="flex">
    <main class="flex-1 p-4 md:p-6 overflow-y-auto h-[calc(100svh-61px)] flex flex-col gap-y-6">
      <!-- Header -->
      <div class="space-y-2.5">
        <div class="flex items-center gap-3">
          <Icon
            class="text-lg"
            :class="STATUS_ICON_COLOR_MAP[run?.status as keyof typeof STATUS_ICON_COLOR_MAP].color"
            :name="STATUS_ICON_COLOR_MAP[run?.status as keyof typeof STATUS_ICON_COLOR_MAP].icon"
          />

          <h1 class="text-lg font-semibold text-foreground">
            {{ run?.region?.name || 'Unknown region' }}
            <span v-if="run?.region?.flag" class="ml-2">{{ run.region.flag }}</span>
          </h1>
        </div>

        <div class="text-sm text-muted-foreground flex items-center gap-2">
          <div>
            {{ STATUS_TEXT_MAP[run?.status as keyof typeof STATUS_TEXT_MAP] }} on
            <NuxtTime
              :datetime="run?.created_at ?? new Date()"
              v-bind="{
                locale: 'en-US',
                dateStyle: 'long',
                timeStyle: 'short',
              }"
            />
          </div>

          <div v-if="run?.ip_version && run?.ip_address">
            <span class="mr-2">•</span>
            <span>{{ IP_VERSION_LABELS[run?.ip_version as keyof typeof IP_VERSION_LABELS] }}{{ ':' }} {{ run?.ip_address }}</span>
          </div>
        </div>
        <div class="flex items-center justify-between">
          <div class="text-sm font-mono text-foreground">
            <Badge class="mr-2" variant="secondary">
              {{ run?.check?.method || 'GET' }}
            </Badge>
            {{ constructURL({
              host: run?.check?.host!,
              port: run?.check?.port,
              path: run?.check?.path,
              queryParams: run?.check?.query_params as Record<string, string> | undefined,
              secure: run?.check?.secure,
            }) }}
          </div>
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2">
              <Badge
                class="px-3 text-base font-mono"
                :variant="(run as Record<string, unknown>)?.response_status_code && (run as Record<string, unknown>).response_status_code as number >= 200 && (run as Record<string, unknown>).response_status_code as number < 300 ? 'default' : 'destructive'"
              >
                {{ (run as Record<string, unknown>)?.response_status_code || '—' }}
              </Badge>
            </div>
            <div class="flex items-center gap-2 text-sm font-mono">
              <Icon class="w-4 h-4 text-muted-foreground" name="lucide:clock" />
              <span>{{ (run as Record<string, unknown>)?.total_time_ms || 0 }}ms</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Assertions -->
      <Card v-if="run?.assertion_results && run.assertion_results.length > 0">
        <CardHeader>
          <div class="flex items-center justify-between">
            <CardTitle>Assertions</CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-border">
                  <th class="text-left py-2 px-3 text-muted-foreground font-medium">
                    Source
                  </th>
                  <th class="text-left py-2 px-3 text-muted-foreground font-medium">
                    Property
                  </th>
                  <th class="text-left py-2 px-3 text-muted-foreground font-medium">
                    Comparison
                  </th>
                  <th class="text-left py-2 px-3 text-muted-foreground font-medium">
                    Target
                  </th>
                  <th class="text-left py-2 px-3 text-muted-foreground font-medium">
                    Actual
                  </th>
                  <th class="text-left py-2 px-3 text-muted-foreground font-medium sr-only">
                    Passed
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(assertion, index) in run?.assertion_results" :key="index" class="border-b border-border/50">
                  <td class="py-2 px-3">
                    {{ ASSERTION_PROPERTIES[assertion.source as keyof typeof ASSERTION_PROPERTIES].label }}
                  </td>
                  <td class="py-2 px-3">
                    {{ assertion.property || '—' }}
                  </td>
                  <td class="py-2 px-3">
                    {{ assertion.comparison || '—' }}
                  </td>
                  <td class="py-2 px-3">
                    <span v-if="assertion.target !== undefined" class="font-mono">{{ assertion.target }}</span>
                    <span v-else>—</span>
                  </td>
                  <td class="py-2 px-3">
                    <span v-if="assertion.received !== undefined" class="font-mono">{{ assertion.received }}</span>
                    <span v-else>—</span>
                  </td>
                  <td class="py-2 px-3">
                    <Icon
                      :class="assertion.passed ? 'text-green-500' : 'text-red-500'"
                      :name="assertion.passed ? 'lucide:check-circle' : 'lucide:x-circle'"
                    />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      <!-- Network Timings -->
      <Card v-if="timelineData.length > 0">
        <CardHeader>
          <CardTitle>Network Timings</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="divide-y divide-border">
            <div class="grid grid-cols-[160px_1fr_80px] gap-2 h-8 items-center bg-muted">
              <div class="text-xs font-semibold uppercase px-3">
                Connection Start
              </div>
              <div></div>
              <div class="text-right text-xs font-semibold uppercase px-3">
                Time
              </div>
            </div>
            <div v-for="timing in timelineData.slice(0, 3)" :key="timing.type" class="grid grid-cols-[160px_1fr_80px] gap-2 h-8 items-center">
              <div class="text-sm font-semibold px-3">
                {{ timing.type }}
              </div>
              <div class="relative h-4">
                <div class="absolute inset-0" :style="{ left: `${timing.left}%`, width: `${timing.width}%`, backgroundColor: timing.color }"></div>
              </div>
              <div class="text-right text-sm font-semibold px-3">
                {{ durationFormatter(timing.length) }}
              </div>
            </div>
            <div class="grid grid-cols-[160px_1fr_80px] gap-2 h-8 items-center bg-muted">
              <div class="text-xs font-semibold uppercase px-3">
                Request/Response
              </div>
              <div></div>
              <div class="text-right text-xs font-semibold uppercase px-3">
                Time
              </div>
            </div>
            <div v-for="timing in timelineData.slice(3)" :key="timing.type" class="grid grid-cols-[160px_1fr_80px] gap-2 h-8 items-center">
              <div class="text-sm font-semibold px-3">
                {{ timing.type }}
              </div>
              <div class="relative h-4">
                <div class="absolute inset-0" :style="{ left: `${timing.left}%`, width: `${timing.width}%`, backgroundColor: timing.color }"></div>
              </div>
              <div class="text-right text-sm font-semibold px-3">
                {{ durationFormatter(timing.length) }}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </main>
  </div>
</template>
