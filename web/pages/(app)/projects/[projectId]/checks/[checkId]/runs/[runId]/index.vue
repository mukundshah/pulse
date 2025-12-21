<script lang="ts">
import { ASSERTION_PROPERTIES } from '@/constants/http'
import { formatDuration } from '@/utils/formatters'
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
  title: `${run.value?.check?.name} - #${run.value?.run_number}`,
})

useLayoutContext({
  breadcrumbOverrides: computed(() => [
    undefined, // Root
    undefined, // Projects
    {
      label: run.value?.check?.project?.name || 'Project',
      to: `/projects/${projectId}/checks`,
    }, // Project
    {
      label: run.value?.check?.name || 'Check',
      to: `/projects/${projectId}/checks/${checkId}`,
    }, // Check
    false, // false to hide the check id
    {
      label: `#${run.value?.run_number}`,
      to: `/projects/${projectId}/checks/${checkId}/runs/${runId}`,
      active: true,
    }, // Run
    false, // false to hide the current breadcrumb item
  ]),
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

interface DNSAnswer {
  name: string
  type: string
  TTL: number
  data: string | string[] | Record<string, unknown>
  preference?: number
  port?: number
  priority?: number
  weight?: number
}

interface DNSAuthority {
  name: string
  type: string
  TTL: number
}

interface DNSQuestion {
  Name: string
  Type: string
}

interface DNSJSONFormat {
  Status?: string
  TC?: boolean
  AD?: boolean
  CD?: boolean
  ID?: number
  Question?: DNSQuestion[]
  Answer?: DNSAnswer[]
  Authority?: DNSAuthority[]
  Additional?: DNSAuthority[]
}

const dnsResponse = computed(() => {
  const runData = run.value as Record<string, unknown> | undefined
  if (!runData?.response || typeof runData.response !== 'object') {
    return null
  }

  const response = runData.response as Record<string, unknown>

  // Check if it's a DNS response
  if (response.type !== 'dns') {
    return null
  }

  return {
    type: response.type,
    records: response.records,
    dns_server: response.dns_server,
    formats: response.formats as {
      raw?: string
      json?: DNSJSONFormat
    } | undefined,
  }
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

          <div class="flex items-center gap-2 flex-1">
            <h1 class="text-lg font-semibold text-foreground">
              {{ run?.check?.name || 'Unnamed Check' }}
            </h1>
            <Badge v-if="run?.region" variant="secondary">
              {{ run.region.flag }} {{ run.region.name }}
            </Badge>
          </div>
        </div>

        <div class="text-sm text-muted-foreground flex items-center gap-2">
          <div>
            Run #{{ run?.run_number || '—' }} • {{ STATUS_TEXT_MAP[run?.status as keyof typeof STATUS_TEXT_MAP] }} on
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
            <span>{{ run?.ip_version }}{{ ':' }} {{ run?.ip_address }}</span>
          </div>
        </div>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Badge class="text-xs font-mono" variant="secondary">
              {{ run?.check?.type.toUpperCase() }}
            </Badge>

            <div v-if="run?.check?.type === 'http'" class="text-sm text-muted-foreground font-mono flex items-center gap-2">
              <Badge class="text-xs" variant="secondary">
                {{ run?.check?.method }}
              </Badge>

              <Badge class="text-xs" variant="secondary">
                {{ constructURL({
                  host: run?.check?.host!,
                  port: run?.check?.port,
                  path: run?.check?.path,
                  queryParams: run?.check?.query_params as Record<string, string> | undefined,
                  secure: run?.check?.secure,
                }) }}
              </Badge>
            </div>
            <div v-else-if="run?.check?.type === 'tcp'" class="text-sm text-muted-foreground font-mono flex items-center gap-2">
              <Badge class="text-xs" variant="secondary">
                {{ run?.check?.host }}{{ ':' }}{{ run?.check?.port }}
              </Badge>
            </div>
            <div v-else-if="run?.check?.type === 'dns'" class="text-sm text-muted-foreground font-mono flex items-center gap-2">
              <Badge class="text-xs" variant="secondary">
                {{ run?.check?.dns_record_type?.toUpperCase() }}
              </Badge>
              <Badge class="text-xs" variant="secondary">
                {{ run?.check?.host }}
              </Badge>
            </div>
          </div>

          <div class="flex items-center gap-4">
            <div v-if="run?.check?.type === 'http'" class="flex items-center gap-2">
              <Badge
                class="px-3 text-base font-mono"
                :variant="(run as Record<string, unknown>)?.response_status_code && (run as Record<string, unknown>).response_status_code as number >= 200 && (run as Record<string, unknown>).response_status_code as number < 300 ? 'default' : 'destructive'"
              >
                {{ (run as Record<string, unknown>)?.response_status_code || '—' }}
              </Badge>
            </div>
            <div class="flex items-center gap-2 text-sm font-mono">
              <Icon class="w-4 h-4 text-muted-foreground" name="lucide:clock" />
              <FormattedNumber class="text-sm font-mono text-muted-foreground" :options="{ style: 'unit', unit: 'millisecond' }" :value="(run as Record<string, unknown>)?.total_time_ms as number || 0" />
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
                  <td class="py-2 px-3 font-mono">
                    <Badge v-if="assertion.property" variant="secondary">
                      {{ assertion.property }}
                    </Badge>
                    <span v-else>n/a</span>
                  </td>
                  <td class="py-2 px-3">
                    {{ assertion.comparison }}
                  </td>
                  <td class="py-2 px-3 font-mono">
                    <Badge variant="secondary">
                      {{ assertion.target }}
                    </Badge>
                  </td>
                  <td class="py-2 px-3 font-mono">
                    <Badge v-if="assertion.received" variant="secondary" :class="assertion.passed ? 'text-green-600 bg-green-100' : 'text-red-600 bg-red-100'">
                      {{ assertion.received }}
                    </Badge>
                    <span v-else>n/a</span>
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
                {{ formatDuration(timing.length) }}
              </div>
            </div>
            <div v-if="timelineData.length > 3" class="grid grid-cols-[160px_1fr_80px] gap-2 h-8 items-center bg-muted">
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
                {{ formatDuration(timing.length) }}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- DNS Response -->
      <Card v-if="run?.check?.type === 'dns' && dnsResponse">
        <CardHeader>
          <CardTitle>DNS Response</CardTitle>
        </CardHeader>
        <CardContent>
          <!-- DNS Server -->
          <div v-if="dnsResponse.dns_server" class="mb-6 flex items-center gap-3 rounded-lg border bg-muted/30 p-4">
            <Icon class="size-5 text-muted-foreground" name="lucide:server" />
            <div>
              <div class="text-xs font-medium uppercase text-muted-foreground">
                DNS Server
              </div>
              <div class="font-mono text-sm font-semibold">
                {{ dnsResponse.dns_server }}
              </div>
            </div>
          </div>

          <!-- Tabs for Structured, JSON, and Raw -->
          <Tabs default-value="structured">
            <TabsList class="mb-4">
              <TabsTrigger value="structured">
                <Icon class="size-4" name="lucide:table" />
                Structured
              </TabsTrigger>
              <TabsTrigger value="json">
                <Icon class="size-4" name="lucide:file-json" />
                JSON
              </TabsTrigger>
              <TabsTrigger value="raw">
                <Icon class="size-4" name="lucide:code" />
                Raw
              </TabsTrigger>
            </TabsList>

            <TabsContent class="space-y-6" value="structured">
              <!-- Records -->
              <div v-if="dnsResponse.formats?.json && dnsResponse.formats.json.Answer && dnsResponse.formats.json.Answer.length > 0" class="space-y-3">
                <div class="flex items-center gap-2">
                  <Icon class="size-4 text-muted-foreground" name="lucide:check-circle" />
                  <h3 class="text-sm font-semibold">
                    Records
                  </h3>
                  <Badge class="ml-auto" variant="secondary">
                    {{ dnsResponse.formats.json.Answer.length }}
                  </Badge>
                </div>
                <div class="overflow-x-auto rounded-lg border">
                  <table class="w-full text-sm">
                    <thead>
                      <tr class="border-b bg-muted/50">
                        <th class="text-left py-3 px-4 text-xs font-semibold uppercase text-muted-foreground">
                          Name
                        </th>
                        <th class="text-left py-3 px-4 text-xs font-semibold uppercase text-muted-foreground">
                          Type
                        </th>
                        <th class="text-left py-3 px-4 text-xs font-semibold uppercase text-muted-foreground">
                          TTL
                        </th>
                        <th class="text-left py-3 px-4 text-xs font-semibold uppercase text-muted-foreground">
                          Data
                        </th>
                        <th v-if="dnsResponse.formats.json.Answer.some(a => a.preference !== undefined || a.port !== undefined || a.priority !== undefined)" class="text-left py-3 px-4 text-xs font-semibold uppercase text-muted-foreground">
                          Additional
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="(answer, index) in dnsResponse.formats.json.Answer" :key="index" class="border-b border-border/50 transition-colors hover:bg-muted/30">
                        <td class="py-3 px-4 font-mono text-sm">
                          {{ answer.name }}
                        </td>
                        <td class="py-3 px-4">
                          <Badge variant="secondary">
                            {{ answer.type }}
                          </Badge>
                        </td>
                        <td class="py-3 px-4 font-mono text-sm">
                          {{ answer.TTL }}s
                        </td>
                        <td class="py-3 px-4">
                          <div class="flex flex-wrap gap-1">
                            <Badge
                              v-if="Array.isArray(answer.data)"
                              class="font-mono text-xs"
                              variant="secondary"
                            >
                              {{ answer.data.join(', ') }}
                            </Badge>
                            <Badge
                              v-else-if="typeof answer.data === 'object'"
                              class="font-mono text-xs"
                              variant="secondary"
                            >
                              {{ JSON.stringify(answer.data) }}
                            </Badge>
                            <Badge
                              v-else
                              class="font-mono text-xs"
                              variant="secondary"
                            >
                              {{ answer.data }}
                            </Badge>
                          </div>
                        </td>
                        <td v-if="answer.preference !== undefined || answer.port !== undefined || answer.priority !== undefined" class="py-3 px-4">
                          <div class="flex flex-wrap gap-2 text-xs">
                            <Badge v-if="answer.preference !== undefined" class="font-mono" variant="outline">
                              Pref: {{ answer.preference }}
                            </Badge>
                            <Badge v-if="answer.priority !== undefined" class="font-mono" variant="outline">
                              Priority: {{ answer.priority }}
                            </Badge>
                            <Badge v-if="answer.weight !== undefined" class="font-mono" variant="outline">
                              Weight: {{ answer.weight }}
                            </Badge>
                            <Badge v-if="answer.port !== undefined" class="font-mono" variant="outline">
                              Port: {{ answer.port }}
                            </Badge>
                          </div>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
              <div v-else class="rounded-lg border border-dashed p-8 text-center">
                <Icon class="mx-auto mb-2 size-8 text-muted-foreground" name="lucide:inbox" />
                <p class="text-sm text-muted-foreground">
                  No answer records available
                </p>
              </div>
            </TabsContent>

            <TabsContent value="json">
              <!-- JSON Format -->
              <div v-if="dnsResponse.formats?.json">
                <div class="rounded-lg border bg-muted/30 p-4">
                  <pre class="overflow-x-auto text-xs font-mono leading-relaxed whitespace-pre-wrap text-foreground">{{ JSON.stringify(dnsResponse.formats.json, null, 2) }}</pre>
                </div>
              </div>
              <div v-else class="rounded-lg border border-dashed p-8 text-center">
                <Icon class="mx-auto mb-2 size-8 text-muted-foreground" name="lucide:file-x" />
                <p class="text-sm text-muted-foreground">
                  No JSON response data available
                </p>
              </div>
            </TabsContent>

            <TabsContent value="raw">
              <!-- Raw Format -->
              <div v-if="dnsResponse.formats?.raw">
                <div class="rounded-lg border bg-muted/30 p-4">
                  <pre class="overflow-x-auto text-xs font-mono leading-relaxed whitespace-pre-wrap text-foreground">{{ dnsResponse.formats.raw }}</pre>
                </div>
              </div>
              <div v-else class="rounded-lg border border-dashed p-8 text-center">
                <Icon class="mx-auto mb-2 size-8 text-muted-foreground" name="lucide:file-x" />
                <p class="text-sm text-muted-foreground">
                  No raw response data available
                </p>
              </div>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </main>
  </div>
</template>
